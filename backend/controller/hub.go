package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"

	"github.com/metaclips/LetsTalk/backend/model"
	"github.com/metaclips/LetsTalk/backend/values"
)

// ServeWs handles websocket requests from the peer, ensuring user is registered.
func ServeWs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		log.Errorln("error parsing form on websocket serve, err:", err)
	}

	cookie := model.CookieDetail{CookieName: values.UserCookieName, Collection: values.UsersCollectionName}

	if err := cookie.CheckCookie(r, values.UserCookieName); err != nil {
		ws, err := model.Upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Errorln("error upgrading websocket while sending error message, err: ", err)
			return
		}

		if err := ws.WriteJSON(struct {
			MessageType string `json:"msgType"`
		}{
			values.UnauthorizedAcces,
		}); err != nil {
			log.Errorln("could not send unauthorized message to user, err:", err)
		}

		if err := ws.Close(); err != nil {
			log.Errorln("error closing websocket on unauthorized access, err:", err)
		}

		return
	}

	ws, err := model.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorln("error upgrading websocket, err:", err)
		return
	}
	model.HubConstruct.RegisterWS(ws, cookie.Email)
}
