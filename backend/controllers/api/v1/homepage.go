package v1controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	apihelpers "github.com/metaclips/LetsTalk/backend/apihelper"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	api_errors "github.com/metaclips/LetsTalk/backend/helpers/errors"
	app_context "github.com/metaclips/LetsTalk/backend/middlewares/context"
	"github.com/metaclips/LetsTalk/backend/middlewares/cookies"
	"github.com/metaclips/LetsTalk/backend/models"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HomepageController struct {
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

func NewHomepageController(hub *v1resources.Hub, message models.Messages, user models.User, room models.Room, roomRequest models.RoomRequest, cookies cookies.Cookies) *HomepageController {
	ctx, cancelFunc := app_context.CreateChildCancelContext()
	return &HomepageController{
		message,
		user, room, roomRequest, ctx, cancelFunc, cookies, websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		hub,
	}
}

func (h HomepageController) HomePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ID, err := h.cookie.CheckCookie(h.ctx, r, constants.UserCookieName)
	if err != nil {
		log.Errorln("error checking users cookies: ", err)
		apihelpers.Respond(w, http.StatusBadRequest, nil, "", api_errors.ErrValidatingCredentials.Error())
		return
	}

	joinedRooms, err := h.roomRequest.RetrieveAcceptedJoinRequest(h.ctx, *ID, 0, 0)
	if err != nil {
		log.Errorln("error checking joined rooms: ", err)
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	toJoinRequest, err := h.roomRequest.RetrieveUnacceptedRooms(h.ctx, *ID, 0, 0)
	if err != nil {
		log.Errorln("error checking room requests: ", err)
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	joinedRoomList, err := h.room.RetrieveJoinedRooms(h.ctx, joinedRooms.Data, 0, 0)
	if err != nil {
		log.Errorln("error checking joined room: ", err)
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	response := v1resources.HomeResponse{
		JoinedRoomList:  *joinedRoomList,
		RoomRequestList: *toJoinRequest,
	}

	apihelpers.Respond(w, http.StatusOK, response, "successful login", "")
}

func (h HomepageController) sendMessage(data interface{}, user primitive.ObjectID) error {
	val, err := json.Marshal(data)
	if err != nil {
		return err
	}

	m := v1resources.WSMessage{val, user}
	h.hub.Broadcast <- m

	return nil
}

// ReadPump pumps messages from the websocket connection to the hub.
func (h HomepageController) readPump(userID primitive.ObjectID, session v1resources.Subscription) {
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
			msg.handleRequestMessages(user)

		case constants.NewMessageMsgType:
			h.newMessage(userID)

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
