package v1controller

import (
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/julienschmidt/httprouter"
	apihelpers "github.com/metaclips/LetsTalk/backend/apihelper"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	api_errors "github.com/metaclips/LetsTalk/backend/helpers/errors"
	"github.com/metaclips/LetsTalk/backend/helpers/paging"
	"github.com/metaclips/LetsTalk/backend/helpers/validation"
	app_context "github.com/metaclips/LetsTalk/backend/middlewares/context"
	"github.com/metaclips/LetsTalk/backend/middlewares/cookies"
	"github.com/metaclips/LetsTalk/backend/models"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminLoginController struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
	cookie     cookies.Cookies
}

func NewAdminLoginController(user models.User, cookies cookies.Cookies) *AdminLoginController {
	ctx, cancelFunc := app_context.CreateChildCancelContext()
	return &AdminLoginController{
		ctx, cancelFunc, cookies,
	}
}

// TODO: We should delete this
func (l *LoginController) AdminLogin(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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

	cookie, uuid, err := l.cookie.GenerateCookie(l.ctx, login.Email, constants.AdminCookieName)
	if err != nil {
		log.Errorln("error generating cookies for admin: ", err)
		status, errVal := getStatusAndErr(err)
		apihelpers.Respond(w, status, nil, "", errVal)
		return
	}

	id, err := primitive.ObjectIDFromHex(uuid)
	if err != nil {
		log.Errorln("error creating admin ID:", err)
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", err.Error())
		return
	}

	apihelpers.Respond(w, http.StatusOK, v1resources.LoginResponse{
		Token:  cookie,
		UserID: id,
	}, "successful login", "")
}

type banUser struct {
	Email        string `json:"email"`
	ReasonBanned string `json:"reasonBanned"`
}

// TODO(metaclips): send reason banned to all the groups the user is in
func (l *LoginController) AdminBanUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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
	ReasonUnbanned string `json:"reasonBanned"`
}

// TODO(metaclips): send reason unbanned to all the groups the user is in
func (l *LoginController) AdminUnbanUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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

	user.IsBanned = false
	err = l.user.UpdateUser(l.ctx, user)
	if err != nil {
		log.WithError(err).Error("failed to unban user")
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	apihelpers.Respond(w, http.StatusOK, user, "unbanned user successfully", "")
}

func (l *LoginController) AdminGetBannedUsers(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	paging := &paging.Paging{}
	if err := json.NewDecoder(req.Body).Decode(paging); err != nil {
		log.WithError(err).Error("error decoding unban data")
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	if paging.Limit == 0 {
		paging.Limit = 20
	}

	bannedUsers, err := l.user.GetUsers(l.ctx, paging, bson.M{"isBanned": true})
	if err != nil {
		log.WithError(err).Error("failed to fetch banned users")
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	apihelpers.Respond(w, http.StatusOK, bannedUsers, "fetched banned users successfully", "")
}
