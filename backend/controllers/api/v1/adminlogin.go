package v1controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/julienschmidt/httprouter"
	apihelpers "github.com/metaclips/LetsTalk/backend/apihelper"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	api_errors "github.com/metaclips/LetsTalk/backend/helpers/errors"
	"github.com/metaclips/LetsTalk/backend/helpers/paging"
	"github.com/metaclips/LetsTalk/backend/helpers/validation"
	app_context "github.com/metaclips/LetsTalk/backend/middlewares/context"
	"github.com/metaclips/LetsTalk/backend/models"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminLoginController struct {
	ctx         context.Context
	cancelFunc  context.CancelFunc
	user        models.User
	roomRequest models.RoomRequest
	messages    models.Messages
}

func NewAdminLoginController(user models.User, roomRequest models.RoomRequest, messages models.Messages) *AdminLoginController {
	ctx, cancelFunc := app_context.CreateChildCancelContext()
	return &AdminLoginController{
		ctx, cancelFunc, user, roomRequest, messages,
	}
}

var adminUUID primitive.ObjectID

// TODO: We should delete this
func (l *AdminLoginController) AdminLogin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	login := v1resources.UserLogin{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&login); err != nil {
		log.Errorln("error decoding login data: ", err)
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	if err := validation.ValidateInputs(login); err != nil {
		apihelpers.Respond(w, http.StatusBadRequest, nil, "", err.Error())
		return
	}

	if login.Email != apihelpers.Config.AdminEmail || login.Password != apihelpers.Config.AdminPassword {
		apihelpers.Respond(w, http.StatusBadRequest, nil, "", "invalid username or password")
		return
	}

	adminUUID = primitive.NewObjectID()

	apihelpers.Respond(w, http.StatusOK, v1resources.LoginResponse{
		Token:  adminUUID.Hex(),
		UserID: adminUUID,
	}, "successful login", "")
}

type banUser struct {
	Email        string `json:"email"`
	ReasonBanned string `json:"reasonBanned"`
}

// TODO(metaclips): send reason banned to all the groups the user is in
func (l *AdminLoginController) AdminBanUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	banUser := &banUser{}
	if err := json.NewDecoder(req.Body).Decode(banUser); err != nil {
		log.Errorln("error decoding ban data: ", err)
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	if banUser.Email == "" {
		apihelpers.Respond(w, http.StatusBadRequest, nil, "", "provide a user email")
		return
	}

	user, err := l.user.GetUserByEmail(l.ctx, banUser.Email)
	if err != nil {
		log.WithError(err).Error("failed to fetch user")
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	rooms, err := l.roomRequest.RetrieveAcceptedJoinRequest(l.ctx, user.ID, 0, 0)
	if err != nil {
		log.Errorln("error getting users associates:", err)
		return
	}

	for _, room := range rooms.Data {
		data := v1resources.Message{
			ID:        primitive.NewObjectID(),
			RoomID:    room.RoomID,
			UserID:    user.ID,
			Message:   fmt.Sprintf("I %s %s [%s] was banned from the platform: Reason %s", user.FirstName, user.LastName, user.Email, banUser.ReasonBanned),
			Type:      constants.MessageTypeMessage,
			CreatedAt: time.Now(),
		}

		if err := l.messages.NewMessage(l.ctx, &data); err != nil {
			log.Errorln("error sending banned message to room:", err)
		}
	}

	user.IsBanned = true
	err = l.user.UpdateUser(l.ctx, user)
	if err != nil {
		log.WithError(err).Error("failed to ban user")
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	apihelpers.Respond(w, http.StatusOK, user, "banned user successfully", "")
}

type unbanUser struct {
	Email          string `json:"email"`
	ReasonUnbanned string `json:"reasonUnBanned"`
}

// TODO(metaclips): send reason unbanned to all the groups the user is in
func (l *AdminLoginController) AdminUnbanUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	unbanUser := &unbanUser{}
	if err := json.NewDecoder(req.Body).Decode(unbanUser); err != nil {
		log.WithError(err).Error("error decoding unban data")
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	if unbanUser.Email == "" || unbanUser.ReasonUnbanned == "" {
		apihelpers.Respond(w, http.StatusBadRequest, nil, "", "user email and reason unbanned are required")
		return
	}

	user, err := l.user.GetUserByEmail(l.ctx, unbanUser.Email)
	if err != nil {
		log.WithError(err).Error("failed to fetch user")
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	rooms, err := l.roomRequest.RetrieveAcceptedJoinRequest(l.ctx, user.ID, 0, 0)
	if err != nil {
		log.Errorln("error getting users associates:", err)
		return
	}

	for _, room := range rooms.Data {
		data := v1resources.Message{
			ID:        primitive.NewObjectID(),
			RoomID:    room.RoomID,
			UserID:    user.ID,
			Message:   fmt.Sprintf("I %s %s [%s] have been unbanned from the platform: Reason %s", user.FirstName, user.LastName, user.Email, unbanUser.ReasonUnbanned),
			Type:      constants.MessageTypeMessage,
			CreatedAt: time.Now(),
		}

		if err := l.messages.NewMessage(l.ctx, &data); err != nil {
			log.Errorln("error sending banned message to room:", err)
		}
	}

	user.IsBanned = false
	err = l.user.UpdateUser(l.ctx, user)
	if err != nil {
		log.WithError(err).Error("failed to unban user")
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	apihelpers.Respond(w, http.StatusOK, user, "unbanned user successfully", "")
}

func (l *AdminLoginController) VerifyUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	unbanUser := &unbanUser{}
	if err := json.NewDecoder(req.Body).Decode(unbanUser); err != nil {
		log.WithError(err).Error("error decoding unban data")
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	user, err := l.user.GetUserByEmail(l.ctx, unbanUser.Email)
	if err != nil {
		log.WithError(err).Error("failed to fetch user")
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	if user.Verified {
		apihelpers.Respond(w, http.StatusOK, nil, "user already verified", "")
		return
	}

	user.Verified = true
	err = l.user.UpdateUser(l.ctx, user)
	if err != nil {
		log.WithError(err).Error("failed to unban user")
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	apihelpers.Respond(w, http.StatusOK, user, "user verified successfully", "")
}

func (l *AdminLoginController) AdminGetBannedUsers(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	page := &paging.Paging{}
	if err := json.NewDecoder(req.Body).Decode(page); err != nil {
		log.WithError(err).Error("error decoding page data")
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	if page.Limit == 0 {
		page.Limit = 20
	}

	bannedUsers, err := l.user.GetUsers(l.ctx, page, bson.M{"isBanned": true})
	if err != nil {
		log.WithError(err).Error("failed to fetch banned users")
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	apihelpers.Respond(w, http.StatusOK, bannedUsers, "fetched banned users successfully", "")
}

func (l *AdminLoginController) AdminGetUserByEmail(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	userEmail := &unbanUser{}
	if err := json.NewDecoder(req.Body).Decode(userEmail); err != nil {
		log.WithError(err).Error("error decoding user email data")
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	if userEmail.Email == "" {
		apihelpers.Respond(w, http.StatusBadRequest, nil, "", "provide a user email")
		return
	}

	user, err := l.user.GetUserByEmail(l.ctx, userEmail.Email)
	if err != nil {
		log.WithError(err).Error("failed to fetch user")
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	apihelpers.Respond(w, http.StatusOK, user, "fetched user successfully", "")
}

func (l *AdminLoginController) AdminGetMinorUsers(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	page := &paging.Paging{}
	if err := json.NewDecoder(req.Body).Decode(page); err != nil {
		log.WithError(err).Error("error decoding page data")
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	if page.Limit == 0 {
		page.Limit = 20
	}

	bannedUsers, err := l.user.GetRegisteredMinors(l.ctx, page)
	if err != nil {
		log.WithError(err).Error("failed to fetch minor users")
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	apihelpers.Respond(w, http.StatusOK, bannedUsers, "fetched minor users successfully", "")
}
