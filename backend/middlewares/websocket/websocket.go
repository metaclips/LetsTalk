package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	"github.com/metaclips/LetsTalk/backend/helpers/maps"
	"github.com/metaclips/LetsTalk/backend/helpers/paging"
	"github.com/metaclips/LetsTalk/backend/helpers/validation"
	appContext "github.com/metaclips/LetsTalk/backend/middlewares/context"
	"github.com/metaclips/LetsTalk/backend/middlewares/peerconnection"
	"github.com/metaclips/LetsTalk/backend/models"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type gorrilaWebsocket struct {
	// Websocket Hub.
	hub            *hub
	peerconnection peerconnection.PeerInterface

	ctx        context.Context
	cancelFunc context.CancelFunc

	roomRequest models.RoomRequest
	messages    models.Messages
	user        models.User
	room        models.Room
}

var _ WSInterface = &gorrilaWebsocket{}

func NewGorillaWebsocket(peerconnection peerconnection.PeerInterface, roomRequest models.RoomRequest, messages models.Messages, user models.User, room models.Room) WSInterface {
	hub := &hub{
		broadcast:  make(chan wsMessage),
		register:   make(chan Subscription),
		unRegister: make(chan Subscription),
		users: wsUsers{
			roomUsers: make(map[primitive.ObjectID]*Connection),
			mutex:     &sync.RWMutex{},
		},
	}

	ctx, cancelFunc := appContext.CreateChildCancelContext()
	ws := &gorrilaWebsocket{
		hub,
		peerconnection,
		ctx, cancelFunc,
		roomRequest,
		messages,
		user,
		room,
	}

	go ws.run()

	return ws
}

func (h *gorrilaWebsocket) RegisterToHub(sub Subscription) {
	h.hub.register <- sub
}

func (h *gorrilaWebsocket) run() {
	for {
		select {
		case s := <-h.hub.register:
			h.hub.users.mutex.Lock()

			// Close recent connection if exists.
			// There's only one supported websocket connection on the hub.
			if conn, exists := h.hub.users.roomUsers[s.UserID]; exists {
				if err := conn.WS.Close(); err != nil {
					log.Errorln("could not close websocket channel:", err)
				}
			}

			h.hub.users.roomUsers[s.UserID] = s.Conn

			log.Infoln(s.UserID, "registered to hub.")
			go h.broadcastOnlineStatusToAllUserRoom(s.UserID, true)

			h.hub.users.mutex.Unlock()

		case s := <-h.hub.unRegister:
			h.hub.users.mutex.Lock()
			if _, exists := h.hub.users.roomUsers[s.UserID]; exists {
				delete(h.hub.users.roomUsers, s.UserID)
				close(s.Conn.Send)
				go h.broadcastOnlineStatusToAllUserRoom(s.UserID, false)
				log.Infoln(s.UserID, "offline")
			}
			h.hub.users.mutex.Unlock()

		case m := <-h.hub.broadcast:
			h.hub.users.mutex.RLock()
			connection, exist := h.hub.users.roomUsers[m.user]
			if !exist {
				log.Errorln("ws conn closed and message broadcasted:", m.user)
				continue
			}

			select {
			case connection.Send <- m.data:
			default:
				connection.WS.Close()
				delete(h.hub.users.roomUsers, m.user)
				go h.broadcastOnlineStatusToAllUserRoom(m.user, false)
			}
			h.hub.users.mutex.RUnlock()
		}
	}
}

func (h *gorrilaWebsocket) WritePump(conn *Connection) {
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
			// if err := h.write(websocket.PingMessage, []byte{}, conn.WS); err != nil {
			// 	return
			// }
		}
	}
}

// write writes a message with the given message type and payload.
func (h *gorrilaWebsocket) write(mt int, payload []byte, ws *websocket.Conn) error {
	if err := ws.SetWriteDeadline(time.Now().Add(constants.WriteWait)); err != nil {
		return err
	}

	return ws.WriteMessage(mt, payload)
}

// ReadPump pumps messages from the websocket connection to the hub.
func (h *gorrilaWebsocket) ReadPump(userID primitive.ObjectID, session Subscription) {
	c := session.Conn

	defer func() {
		h.hub.unRegister <- session
		c.WS.Close()
	}()

	c.WS.SetReadDeadline(time.Now().Add(constants.PongWait))

	c.WS.SetPongHandler(
		func(string) error {
			return c.WS.SetReadDeadline(time.Now().Add(constants.PongWait))
		},
	)

	for {
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
			continue
		}

		switch data.MsgType {
		// TODO: add support to remove message.
		case constants.WebsocketOpenMsgType:
			h.loadUserContent(userID)

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

		case constants.ExitRoomMsgType:
			h.handleExitRoom(userID, msg)

		case constants.SearchUserMsgType:
			h.searchUsers(data.SearchText, session.UserID)

		// case constants.NewFileUploadMsgType:
		// 	h.newFileUpload(userID, msg)

		// case constants.UploadFileChunkMsgType:
		// 	msg.handleUploadFileChunk()

		// case constants.UploadFileSuccessMsgType:
		// 	msg.handleUploadFileUploadComplete()

		// case constants.RequestDownloadMsgType:
		// 	msg.handleRequestDownload(user)

		// case constants.DownloadFileChunkMsgType:
		// 	msg.handleFileDownload(user)

		case constants.StartClassSession:
			h.peerconnection.NewPublisherPeerConnection(userID, msg, h.onPC)

		case constants.JoinClassSession:
			h.peerconnection.NewSubscriberPeerConnection(userID, msg, h.onPC)

		// case constants.EndClassSession:
		// 	classSessions.endClassSession(user)

		// case constants.RenegotiateSDP:
		// 	sdpConstruct{}.acceptRenegotiation(msg)

		default:
			log.Println("Could not convert required type", data.MsgType)
		}
	}
}

func (h *gorrilaWebsocket) onPC(data interface{}, roomID *primitive.ObjectID, userID *primitive.ObjectID) error {
	if roomID != nil {
		msg, err := json.Marshal(data)
		if err != nil {
			return err
		}

		h.newMessage(*userID, msg)
		return nil
	}

	if userID != nil {
		h.sendMessage(data, *userID)
		return nil
	}

	return nil
}

// func (h Websocket) newFileUpload(user primitive.ObjectID, msg []byte) {
// 	var fileData = v1resources.FileDetail{}
// 	if err := unmarshalTo(msg, &fileData); err != nil {
// 		// ToDo: Return error message.
// 		return
// 	}

// 	if err := validation.ValidateInputs(fileData); err != nil {
// 		// ToDo: Return error message.
// 		return
// 	}

// 	if fileData.UserID != user {
// 		// ToDo: Return error message.
// 		return
// 	}

// 	// This tracks file status.
// 	if err := h.file.CreateNewFile(h.ctx, &fileData); err != nil {
// 		// ToDo: Return error message.
// 		return
// 	}

// 	data := struct {
// 		MsgType  string `json:"msgType"`
// 		FileName string `json:"fileName,omitempty"`
// 		FileHash string `json:"fileHash"`
// 		Chunk    int    `json:"nextChunk"`
// 	}{
// 		constants.UploadFileChunkMsgType,
// 		fileData.Name,
// 		fileData.Hash,
// 		1,
// 	}

// 	h.sendMessage(data, user)
// }

func (g *gorrilaWebsocket) loadUserContent(user primitive.ObjectID) {
	acceptedRooms, err := g.roomRequest.RetrieveAcceptedJoinRequest(g.ctx, user, 0, 0)
	if err != nil {
		log.Errorln("error retrieving accepted rooms by user:", err)
		return
	}

	userAssociates, err := g.roomRequest.GetUserAssociates(g.ctx, user)
	if err != nil {
		log.Errorln("error retrieving users associates:", err)
		return
	}

	unacceptedRooms, err := g.roomRequest.RetrieveUnacceptedRooms(g.ctx, user, 0, 0)
	if err != nil {
		log.Errorln("error retrieving unaccepted rooms by users:", err)
		return
	}

	// Broadcast to user that is to join, that there's a request.
	g.sendMessage(map[string]interface{}{
		"msgType":      constants.WebsocketOpenMsgType,
		"joinedRooms":  acceptedRooms,
		"joinRequests": unacceptedRooms,
		"associates":   userAssociates,
	}, user)
}

func (h *gorrilaWebsocket) handleExitRoom(user primitive.ObjectID, msg []byte) {
	var request = v1resources.ExitRoomRequest{}
	if err := unmarshalTo(msg, &request); err != nil {
		// ToDo: Return error message.
		return
	}

	if err := validation.ValidateInputs(request); err != nil {
		// ToDo: Return error message.
		return
	}

	if request.UserID != user {
		// ToDo: Return error message.
		return
	}

	// Ensure the person making the request is a member of the room.
	isRequesterInRoom, err := h.roomRequest.IsRequesterAllowedInRoom(h.ctx, request.RoomID, user)
	if err != nil {
		// ToDo: Return error message.
		return
	}

	if !isRequesterInRoom {
		// ToDo: Return error message.
		return
	}

	userLeaving, err := h.user.GetUserByID(h.ctx, request.UserID)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			log.Errorln("user tried inviting an unregistered user")
		}

		return
	}

	if err := h.roomRequest.ExitRoom(h.ctx, request.JoinRoomID); err != nil {
		log.Errorln("error exiting room for user", err)
		return
	}

	// Broadcast to user that is to join, that there's a request.
	h.sendMessage(map[string]interface{}{
		"msgType": constants.ExitRoomMsgType,
		"data":    request,
	}, user)

	// Send to the other user the room request.
	h.saveMessageAndNotifyRoom(request.RoomID, request.UserID, fmt.Sprintf("%s %s (%s) has exited the room.", userLeaving.FirstName, userLeaving.LastName, userLeaving.Email))
}

func (h *gorrilaWebsocket) searchUsers(key string, user primitive.ObjectID) {
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

	if err := h.sendMessage(map[string]interface{}{
		"msgType": constants.SearchUserMsgType,
		"data":    usersList,
	}, user); err != nil {
		log.Errorln("error sending users data on `searchUsers`: ", err)
		return
	}
}

func (h *gorrilaWebsocket) newMessage(userID primitive.ObjectID, msg []byte) {
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

	if _, ok := maps.MessageTypes[message.Type]; !ok {
		// ToDo: Return error message.
		return
	}

	users, err := h.roomRequest.GetUsersInRoom(h.ctx, message.RoomID)
	if err != nil {
		log.Errorln("error getting users in room:", err)
		// ToDo: Return error message.
		return
	}

	// if message.Type == constants.MessageTypeFile && (message.SizeInKB == nil || message.Hash == "") {
	// 	// ToDo: Return error message.
	// 	return
	// }

	message.ID = primitive.NewObjectID()
	message.CreatedAt = time.Now()

	if err := h.messages.NewMessage(h.ctx, &message); err != nil {
		log.Errorln("error saving message to DB:", err)
		// ToDo: Return error message.
		return
	}

	for _, user := range users {
		// ToDo: Handle this error properly.
		if err := h.sendMessage(map[string]interface{}{
			"msgType": constants.NewMessageMsgType,
			"data":    message,
		}, user.UserToJoinID); err != nil {
			log.Errorln("error sending message to user:", err)
			// ToDo: Return error message.
			continue
		}
	}
}

func (h *gorrilaWebsocket) requestMessages(userID primitive.ObjectID, msg []byte) {
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

	messages, err := h.messages.RetrieveMessagesForRoom(h.ctx, request.RoomID, request.Limit, request.Skip)
	if err != nil {
		log.Errorln("error retrieving requested message to user:", err)

		// ToDo: Return error message.
		return
	}

	if err := h.sendMessage(map[string]interface{}{
		"msgType": constants.RequestMessages,
		"data":    messages,
	}, userID); err != nil {
		log.Errorln("error sending requested message to user:", err)
	}
}

func (h *gorrilaWebsocket) joinRoomRequest(userID primitive.ObjectID, msg []byte) {
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
	isRequesterInRoom, err := h.roomRequest.IsRequesterAllowedInRoom(h.ctx, request.RoomID, userID)
	if err != nil {
		// ToDo: Return error message.
		return
	}

	if !isRequesterInRoom {
		// ToDo: Return error message.
		return
	}

	hasUserAccepted, err := h.roomRequest.HasUserAcceptedJoinRequest(h.ctx, request.UserToJoinID, request.RoomID)
	if err != nil {
		// ToDo: Return error message.
		return
	}

	if hasUserAccepted {
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
	request.CreatedAt = time.Now()

	if err := h.roomRequest.NewRoomRequest(h.ctx, &request); err != nil {
		log.Errorln("error creating room request:", err)
		// ToDo: Return error message.
		return
	}

	requester, err := h.user.GetUserByID(h.ctx, request.UserID)
	if err != nil {
		// ToDo: Return error message.
		return
	}

	// Broadcast to user that is to join, that there's a request.
	h.sendMessage(map[string]interface{}{
		"msgType": constants.RequestUsersToJoinRoomMsgType,
		"data":    request,
	}, request.UserToJoinID)

	// Send to the other user the room request.
	h.saveMessageAndNotifyRoom(request.RoomID, request.UserToJoinID, fmt.Sprintf("%s %s (%s) has requested %s %s (%s) to join the room.", requester.FirstName, requester.LastName, requester.Email, userToJoin.FirstName, userToJoin.LastName, userToJoin.Email))
}

func (h *gorrilaWebsocket) acceptRoomRequest(userID primitive.ObjectID, msg []byte) {
	var request = v1resources.JoinRoomRequest{}
	if err := unmarshalTo(msg, &request); err != nil {
		// ToDo: Return error message.
		return
	}

	if err := validation.ValidateInputs(request); err != nil {
		// ToDo: Return error message.
		return
	}

	hasUserAcceptedJoinRequest, err := h.roomRequest.HasUserAcceptedJoinRequest(h.ctx, request.RoomID, request.UserToJoinID)
	if err != nil {
		// ToDo: Return error message.
		return
	}

	if hasUserAcceptedJoinRequest {
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

	room := v1resources.Room{
		RoomName: request.RoomName,
		ID:       request.RoomID,
	}

	// Broadcast to user to join that there's a new room.
	h.sendMessage(map[string]interface{}{
		"msgType": constants.JoinRoomMsgType,
		"data":    room,
	}, userID)

	h.saveMessageAndNotifyRoom(request.RoomID, request.UserToJoinID, fmt.Sprintf("%s %s (%s) has has joined the room.", userToJoin.FirstName, userToJoin.LastName, userToJoin.Email))
}

func (h *gorrilaWebsocket) createNewRoom(userID primitive.ObjectID, msg []byte) {
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
	room.CreatedAt = time.Now()

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
		CreatedAt:    time.Now(),
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

	if err := h.sendMessage(map[string]interface{}{
		"msgType": constants.JoinRoomMsgType,
		"data":    room,
	}, userID); err != nil {
		// ToDo: Return error message.
		return
	}

	h.saveMessageAndNotifyRoom(room.ID, userID, fmt.Sprintf("%s %s (%s) has has joined the room.", userToJoin.FirstName, userToJoin.LastName, userToJoin.Email))
}

func (h *gorrilaWebsocket) broadcastOnlineStatusToAllUserRoom(userID primitive.ObjectID, isOnline bool) {
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
		if err := h.sendMessage(map[string]interface{}{
			"msgType": constants.OnlineStatusMsgType,
			"data":    status,
		}, assassociate); err != nil {
			log.Errorln("error sending online status message:", err)
		}
	}
}

func (h *gorrilaWebsocket) saveMessageAndNotifyRoom(roomID, userID primitive.ObjectID, message string) {
	data := v1resources.Message{
		ID:        primitive.NewObjectID(),
		RoomID:    roomID,
		UserID:    userID,
		Message:   message,
		Type:      constants.NewMessageMsgType,
		CreatedAt: time.Now(),
	}

	err := h.messages.NewMessage(h.ctx, &data)
	if err != nil {
		// ToDo: Return error message.
		return
	}

	// Broadcast message to all room users that a `join room request` has been sent.
	if err := h.broadcastToAllRoomUsers(roomID, map[string]interface{}{
		"data":    data,
		"msgType": "NewMessage",
	}); err != nil {
		// ToDo: Return error message.
		return
	}
}

func (h *gorrilaWebsocket) broadcastToAllRoomUsers(roomID primitive.ObjectID, data interface{}) error {
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

func (h *gorrilaWebsocket) sendMessage(data interface{}, user primitive.ObjectID) error {
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}

	m := wsMessage{data: val, user: user}
	h.hub.broadcast <- m

	return nil
}

func unmarshalTo(data []byte, to interface{}) error {
	return json.Unmarshal(data, &to)
}
