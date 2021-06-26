package v1controller

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	apihelpers "github.com/metaclips/LetsTalk/backend/apihelpers"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	api_errors "github.com/metaclips/LetsTalk/backend/helpers/errors"
	app_context "github.com/metaclips/LetsTalk/backend/middlewares/context"
	"github.com/metaclips/LetsTalk/backend/middlewares/cookies"
	"github.com/metaclips/LetsTalk/backend/models"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	log "github.com/sirupsen/logrus"
)

type HomepageController struct {
	user        models.User
	room        models.Room
	roomRequest models.RoomRequest
	ctx         context.Context
	cancelFunc  context.CancelFunc
	cookie      cookies.Cookies
}

func NewHomepageController(user models.User, room models.Room, roomRequest models.RoomRequest, cookies cookies.Cookies) *HomepageController {
	ctx, cancelFunc := app_context.CreateChildCancelContext()
	return &HomepageController{
		user, room, roomRequest, ctx, cancelFunc, cookies,
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
