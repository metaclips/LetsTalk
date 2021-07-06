package v1controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	apihelpers "github.com/metaclips/LetsTalk/backend/apihelper"
	api_errors "github.com/metaclips/LetsTalk/backend/helpers/errors"

	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	"github.com/metaclips/LetsTalk/backend/helpers/maps"
	"github.com/metaclips/LetsTalk/backend/helpers/paging"
	"github.com/metaclips/LetsTalk/backend/helpers/validation"
	app_context "github.com/metaclips/LetsTalk/backend/middlewares/context"
	"github.com/metaclips/LetsTalk/backend/middlewares/cookies"
	"github.com/metaclips/LetsTalk/backend/models"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WebsocketController struct {
	messages    models.Messages
	user        models.User
	room        models.Room
	roomRequest models.RoomRequest
	ctx         context.Context
	cancelFunc  context.CancelFunc
	cookie      cookies.Cookies

	// Websocket upgrader.
	upgrader websocket.Upgrader
	// Websocket Hub.
	hub *v1resources.Hub
}

func NewWebsocketController(hub *v1resources.Hub, message models.Messages, user models.User, room models.Room, roomRequest models.RoomRequest, cookies cookies.Cookies) *WebsocketController {
	ctx, cancelFunc := app_context.CreateChildCancelContext()
	return &WebsocketController{
		message,
		user, room, roomRequest, ctx, cancelFunc, cookies, websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		hub,
	}
}

func (h WebsocketController) Websocket(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
	ticket := p.ByName("ticket")
	userID, err := primitive.ObjectIDFromHex(p.ByName("userID"))
	if err != nil {
		log.Errorln("invalid userID:", err)
		apihelpers.Respond(w, http.StatusBadRequest, nil, "", "invalid user ID")
		return
	}

	user, err := h.user.GetUserByUUID(h.ctx, ticket)
	if err != nil {
		log.Errorln("error getting user details using UUID:", err)
		status, errVal := getStatusAndErr(err)
		apihelpers.Respond(w, status, nil, "", errVal)
		return
	}

	if user.ID != userID {
		apihelpers.Respond(w, http.StatusBadRequest, nil, "", "invalid user")
		return
	}

	conn, err := h.upgrader.Upgrade(w, req, nil)
	if err != nil {
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	subscribstion := v1resources.Subscription{
		Conn: &v1resources.Connection{
			WS:   conn,
			Send: make(chan []byte, 256),
		},
		UserID: userID,
	}

	h.hub.Register <- subscribstion

	go h.readPump(userID, subscribstion)
	go h.writePump(subscribstion.Conn)
}

func (h WebsocketController) Run() {
	for {
		select {
		case s := <-h.hub.Register:
			h.hub.Users.Mutex.Lock()

			// Close recent connection if exists.
			// There's only one supported websocket connection on the hub.
			if conn, exists := h.hub.Users.RoomUsers[s.UserID]; exists {
				if err := conn.WS.Close(); err != nil {
					log.Errorln("could not close websocket channel:", err)
				}
			}

			h.hub.Users.RoomUsers[s.UserID] = s.Conn

			log.Infoln(s.UserID, "registered to hub.")
			go h.broadcastOnlineStatusToAllUserRoom(s.UserID, true)

			h.hub.Users.Mutex.Unlock()

		case s := <-h.hub.UnRegister:
			h.hub.Users.Mutex.Lock()
			if _, exists := h.hub.Users.RoomUsers[s.UserID]; exists {
				delete(h.hub.Users.RoomUsers, s.UserID)
				close(s.Conn.Send)
				go h.broadcastOnlineStatusToAllUserRoom(s.UserID, false)
				log.Infoln(s.UserID, "offline")
			}
			h.hub.Users.Mutex.Unlock()

		case m := <-h.hub.Broadcast:
			h.hub.Users.Mutex.RLock()
			connection, exist := h.hub.Users.RoomUsers[m.User]
			if !exist {
				log.Errorln("ws conn closed and message broadcasted:", m.User)
				continue
			}

			select {
			case connection.Send <- m.Data:
			default:
				connection.WS.Close()
				delete(h.hub.Users.RoomUsers, m.User)
				go h.broadcastOnlineStatusToAllUserRoom(m.User, false)
			}
			h.hub.Users.Mutex.RUnlock()
		}
	}
}

func (h WebsocketController) writePump(conn *v1resources.Connection) {
	ticker := time.NewTicker(constants.PingPeriod)
	defer func() {
		ticker.Stop()
		conn.WS.Close()
	}()

	for {
		select {
		case message, ok := <-conn.Send:
			// Send channel has been closed
			if !ok {
				h.write(websocket.CloseMessage, []byte{}, conn.WS)
				return
			}

			if err := h.write(websocket.TextMessage, message, conn.WS); err != nil {
				return
			}

		case <-ticker.C:
			if err := h.write(websocket.PingMessage, []byte{}, conn.WS); err != nil {
				return
			}
		}
	}
}

// write writes a message with the given message type and payload.
func (h WebsocketController) write(mt int, payload []byte, ws *websocket.Conn) error {
	if err := ws.SetWriteDeadline(time.Now().Add(constants.WriteWait)); err != nil {
		return err
	}

	return ws.WriteMessage(mt, payload)
}

// ReadPump pumps messages from the websocket connection to the hub.
func (h WebsocketController) readPump(userID primitive.ObjectID, session v1resources.Subscription) {
	c := session.Conn

	defer func() {
		h.hub.UnRegister <- session
		c.WS.Close()
	}()

	c.WS.SetReadDeadline(time.Now().Add(constants.PongWait))

	c.WS.SetPongHandler(
		func(string) error {
			return c.WS.SetReadDeadline(time.Now().Add(constants.PongWait))
		},
	)

	for {
		var err error
		_, msg, err := c.WS.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Errorln("error: ", err)
			}

			return
		}

		data := struct {
			MsgType    string `json:"msgType"`
			SearchText string `json:"searchText"`
		}{}

		err = json.Unmarshal(msg, &data)
		if err != nil {
			log.Println("could not unmarshal json")
		}

		switch data.MsgType {
		// TODO: add support to remove message.
		case constants.RequestMessages:
			h.requestMessages(userID, msg)

		case constants.NewMessageMsgType:
			h.newMessage(userID, msg)

		case constants.CreateRoomMsgType:
			h.createNewRoom(userID, msg)

		case constants.JoinRoomMsgType:
			h.acceptRoomRequest(userID, msg)

		case constants.RequestUsersToJoinRoomMsgType:
			h.joinRoomRequest(session.UserID, msg)

		// case constants.ExitRoomMsgType:
		// 	msg.handleExitRoom(user)

		case constants.SearchUserMsgType:
			h.searchUsers(data.SearchText, session.UserID)

		// case constants.NewFileUploadMsgType:
		// 	msg.handleNewFileUpload()

		// case constants.UploadFileChunkMsgType:
		// 	msg.handleUploadFileChunk()

		// case constants.UploadFileSuccessMsgType:
		// 	msg.handleUploadFileUploadComplete()

		// case constants.RequestDownloadMsgType:
		// 	msg.handleRequestDownload(user)

		// case constants.DownloadFileChunkMsgType:
		// 	msg.handleFileDownload(user)

		// case constants.StartClassSession:
		// 	classSessions.startClassSession(msg, user)

		// case constants.JoinClassSession:
		// 	classSessions.joinClassSession(msg, user)

		// case constants.EndClassSession:
		// 	classSessions.endClassSession(user)

		// case constants.RenegotiateSDP:
		// 	sdpConstruct{}.acceptRenegotiation(msg)

		default:
			log.Println("Could not convert required type", data.MsgType)
		}
	}
}

func (h WebsocketController) searchUsers(key string, user primitive.ObjectID) {
	if key == "" {
		return
	}

	pager := paging.Paging{
		Limit:  0,
		Skip:   0,
		Search: key,
	}

	usersList, err := h.user.GetUsersPrunedData(h.ctx, &pager)
	if err != nil {
		log.Errorln("error retrieving users data on `searchUsers`: ", err)
		return
	}

	if err := h.sendMessage(usersList, user); err != nil {
		log.Errorln("error sending users data on `searchUsers`: ", err)
		return
	}
}

func (h WebsocketController) newMessage(userID primitive.ObjectID, msg []byte) {
	var message v1resources.Message

	if err := unmarshalTo(msg, &message); err != nil {
		// ToDo: Return error message.
		return
	}

	if err := validation.ValidateInputs(message); err != nil {
		// ToDo: Return error message.
		return
	}

	if message.UserID != userID {
		log.Warning("user tried using a fake userID")
		// ToDo: Return error message.
		return
	}

	if _, ok := maps.FileTypes[message.Type]; !ok {
		// ToDo: Return error message.
		return
	}

	users, err := h.roomRequest.GetUsersInRoom(h.ctx, message.RoomID)
	if err != nil {
		log.Errorln("error getting users in room:", err)
		// ToDo: Return error message.
		return
	}

	if message.Type == constants.FileTypeFile && (message.Size == "" || message.Hash == "") {
		// ToDo: Return error message.
		return
	}

	message.Time = time.Now()
	message.ID = primitive.NewObjectID()

	if err := h.messages.NewMessage(h.ctx, &message); err != nil {
		log.Errorln("error saving message to DB:", err)
		// ToDo: Return error message.
		return
	}

	for _, user := range users {
		// ToDo: Handle this error properly.
		if err := h.sendMessage(message, user.UserToJoinID); err != nil {
			log.Errorln("error sending message to user:", err)
			// ToDo: Return error message.
			continue
		}
	}
}

func (h WebsocketController) requestMessages(userID primitive.ObjectID, msg []byte) {
	request := v1resources.RequestMessage{}
	if err := unmarshalTo(msg, &request); err != nil {
		// ToDo: Return error message.
		return
	}

	userIsInRoom, err := h.roomRequest.IsUserAllowedInRoom(h.ctx, request.RoomID, userID)
	if err != nil {
		log.Errorln("error checking if user is in room:", err)
		// ToDo: Return error message.
		return
	}

	if !userIsInRoom {
		// ToDo: Return error message.
		return
	}

	messages, err := h.messages.RetrieveMessagesForRoom(h.ctx, request.RoomID, int64(request.Limit), int64(request.Skip))
	if err != nil {
		log.Errorln("error retrieving requested message to user:", err)

		// ToDo: Return error message.
		return
	}

	if err := h.sendMessage(messages, userID); err != nil {
		log.Errorln("error sending requested message to user:", err)
	}
}

func (h WebsocketController) joinRoomRequest(userID primitive.ObjectID, msg []byte) {
	var request = v1resources.JoinRoomRequest{}
	if err := unmarshalTo(msg, &request); err != nil {
		// ToDo: Return error message.
		return
	}

	if err := validation.ValidateInputs(request); err != nil {
		// ToDo: Return error message.
		return
	}

	// Ensure the person making the request is a member of the room.
	isUserInRoom, err := h.roomRequest.IsRequesterAllowedInRoom(h.ctx, request.RoomID, userID)
	if err != nil {
		// ToDo: Return error message.
		return
	}

	if !isUserInRoom {
		// ToDo: Return error message.
		return
	}

	userToJoin, err := h.user.GetUserByID(h.ctx, request.UserToJoinID)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Errorln("user tried inviting an unregistered user")
		}

		return
	}

	request.Accepted = false
	request.ExitedRoom = false
	request.ID = primitive.NewObjectID()

	if err := h.roomRequest.NewRoomRequest(h.ctx, &request); err != nil {
		// ToDo: Return error message.
		return
	}

	requester, err := h.user.GetUserByID(h.ctx, request.UserID)
	if err != nil {
		// ToDo: Return error message.
		return
	}

	h.saveMessageAndNotifyRoom(request.RoomID, request.UserID, fmt.Sprintf("%s %s (%s) has requested %s %s (%s) to join the room.", requester.FirstName, requester.LastName, requester.Email, userToJoin.FirstName, userToJoin.LastName, userToJoin.Email))
}

func (h WebsocketController) acceptRoomRequest(userID primitive.ObjectID, msg []byte) {
	var request = v1resources.JoinRoomRequest{}
	if err := unmarshalTo(msg, &request); err != nil {
		// ToDo: Return error message.
		return
	}

	if err := validation.ValidateInputs(request); err != nil {
		// ToDo: Return error message.
		return
	}

	userIsAllowedToRoom, err := h.roomRequest.IsUserAllowedInRoom(h.ctx, request.RoomID, request.UserToJoinID)
	if err != nil {
		// ToDo: Return error message.
		return
	}

	if !userIsAllowedToRoom {
		// ToDo: Return error message.
		return
	}

	if request.UserToJoinID != userID {
		// ToDo: Return error message.
		return
	}

	if request.ID.IsZero() {
		// ToDo: Return error message.
		return
	}

	request.Accepted = true
	request.ExitedRoom = false

	if err := h.roomRequest.AcceptRequest(h.ctx, request.ID); err != nil {
		log.Errorln("error accepting room request")
		// ToDo: Return error message.
		return
	}

	userToJoin, err := h.user.GetUserByID(h.ctx, userID)
	if err != nil {
		// ToDo: Return error message.
		return
	}

	h.saveMessageAndNotifyRoom(request.RoomID, request.UserID, fmt.Sprintf("%s %s (%s) has has joined the room.", userToJoin.FirstName, userToJoin.LastName, userToJoin.Email))
}

func (h WebsocketController) createNewRoom(userID primitive.ObjectID, msg []byte) {
	var room v1resources.Room
	if err := unmarshalTo(msg, &room); err != nil {
		// ToDo: Return error message.
		return
	}

	if err := validation.ValidateInputs(room); err != nil {
		// ToDo: Return error message.
		return
	}

	if room.CreateByUserID != userID {
		// ToDo: Return error message.
		return
	}

	room.ID = primitive.NewObjectID()
	if err := h.room.CreateNewRoom(h.ctx, &room); err != nil {
		// ToDo: Return error message.
		log.Errorln("error creating new room:", err)
		return
	}

	err := h.roomRequest.NewRoomRequest(h.ctx, &v1resources.JoinRoomRequest{
		ID:           primitive.NewObjectID(),
		UserID:       userID,
		UserToJoinID: userID,
		RoomID:       room.ID,
		RoomName:     room.RoomName,
		Accepted:     true,
	})

	if err != nil {
		log.Errorln("error creating new room request:", err)
		return
	}

	userToJoin, err := h.user.GetUserByID(h.ctx, userID)
	if err != nil {
		// ToDo: Return error message.
		return
	}

	h.saveMessageAndNotifyRoom(room.ID, userID, fmt.Sprintf("%s %s (%s) has has joined the room.", userToJoin.FirstName, userToJoin.LastName, userToJoin.Email))
}

func (h WebsocketController) broadcastOnlineStatusToAllUserRoom(userID primitive.ObjectID, isOnline bool) {
	associates, err := h.roomRequest.GetUserAssociates(h.ctx, userID)
	if err != nil {
		log.Errorln("error getting users associates:", err)
		return
	}

	status := v1resources.UserStatus{
		UserID:   userID,
		IsOnline: isOnline,
	}

	for _, assassociate := range associates {
		if err := h.sendMessage(status, assassociate); err != nil {
			log.Errorln("error sending online status message:", err)
		}
	}
}

func (h WebsocketController) saveMessageAndNotifyRoom(roomID, userID primitive.ObjectID, message string) {
	data := v1resources.Message{
		ID:      primitive.NewObjectID(),
		RoomID:  roomID,
		UserID:  userID,
		Message: message,
		Type:    constants.NewMessageMsgType,
		Time:    time.Now(),
	}

	err := h.messages.NewMessage(h.ctx, &data)
	if err != nil {
		// ToDo: Return error message.
		return
	}

	// Broadcast message to all room users that a `join room request` has been sent.
	if err := h.broadcastToAllRoomUsers(roomID, data); err != nil {
		// ToDo: Return error message.
		return
	}
}

func (h WebsocketController) broadcastToAllRoomUsers(roomID primitive.ObjectID, data interface{}) error {
	joinedUsers, err := h.roomRequest.GetUsersInRoom(h.ctx, roomID)
	if err != nil {
		return err
	}

	for _, user := range joinedUsers {
		if err := h.sendMessage(data, user.UserToJoinID); err != nil {
			return err
		}
	}

	return nil
}

func (h WebsocketController) sendMessage(data interface{}, user primitive.ObjectID) error {
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}

	m := v1resources.WSMessage{Data: val, User: user}
	h.hub.Broadcast <- m

	return nil
}

func unmarshalTo(data []byte, to interface{}) error {
	return json.Unmarshal(data, &to)
}
