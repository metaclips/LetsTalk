package v1controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	apihelpers "github.com/metaclips/LetsTalk/backend/apihelper"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	api_errors "github.com/metaclips/LetsTalk/backend/helpers/errors"
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

func (l *LoginController) AdminBanUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	userEmail := &v1resources.UserLogin{}
	if err := json.NewDecoder(req.Body).Decode(userEmail); err != nil {
		log.Errorln("error decoding login data: ", err)
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

	user.IsBanned = true
	err = l.user.UpdateUser(l.ctx, user)
	if err != nil {
		log.WithError(err).Error("failed to ban user")
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	apihelpers.Respond(w, http.StatusOK, nil, "banned user successfuly", "")
}
