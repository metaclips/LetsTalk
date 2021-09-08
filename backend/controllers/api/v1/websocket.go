package v1controller

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	apihelpers "github.com/metaclips/LetsTalk/backend/apihelper"
	apiErrors "github.com/metaclips/LetsTalk/backend/helpers/errors"

	appContext "github.com/metaclips/LetsTalk/backend/middlewares/context"
	appWebsocket "github.com/metaclips/LetsTalk/backend/middlewares/websocket"
	"github.com/metaclips/LetsTalk/backend/models"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebsocketController struct {
	user models.User

	apiWebsocket appWebsocket.WSInterface

	ctx        context.Context
	cancelFunc context.CancelFunc

	// Websocket upgrader.
	upgrader websocket.Upgrader
}

func NewWebsocketController(user models.User, apiWebsocket appWebsocket.WSInterface) *WebsocketController {
	ctx, cancelFunc := appContext.CreateChildCancelContext()

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	return &WebsocketController{
		user, apiWebsocket, ctx, cancelFunc, upgrader,
	}
}

func (h *WebsocketController) Websocket(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
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
		log.Errorln("error upgrading websocket:", err)
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", apiErrors.ErrGeneric.Error())
		return
	}

	subscribstion := appWebsocket.Subscription{
		Conn: &appWebsocket.Connection{
			WS:   conn,
			Send: make(chan []byte, 256),
		},
		UserID: userID,
	}

	h.apiWebsocket.RegisterToHub(subscribstion)

	go h.apiWebsocket.ReadPump(userID, subscribstion)
	go h.apiWebsocket.WritePump(subscribstion.Conn)
	if err := apihelpers.Respond(w, http.StatusOK, nil, "Websocket Connected", ""); err != nil {
		log.Errorln("error responding:", err)
	}
}
