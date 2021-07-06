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
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type LoginController struct {
	user       models.User
	ctx        context.Context
	cancelFunc context.CancelFunc
	cookie     cookies.Cookies
}

func NewLoginController(user models.User, cookies cookies.Cookies) *LoginController {
	ctx, cancelFunc := app_context.CreateChildCancelContext()
	return &LoginController{
		user, ctx, cancelFunc, cookies,
	}
}

func (l LoginController) Login(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
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

	userData, err := l.user.GetUserByEmail(l.ctx, login.Email)
	if err != nil {
		status, err := api_errors.ErrorResponse(err)
		log.Errorln("error decoding login data: ", err)
		apihelpers.Respond(w, status, nil, "", err)
		return
	}

	if err := bcrypt.CompareHashAndPassword(userData.Password, []byte(login.Password)); err != nil {
		log.Errorln("error decoding login data: ", err)
		apihelpers.Respond(w, http.StatusBadRequest, nil, "", api_errors.ErrInvalidCredentials.Error())
		return
	}

	cookie, uuid, err := l.cookie.GenerateCookie(l.ctx, login.Email, constants.UserCookieName)

	if err != nil {
		log.Errorln("error generating cookies for user: ", err)
		status, errVal := getStatusAndErr(err)
		apihelpers.Respond(w, status, nil, "", errVal)
		return
	}

	apihelpers.Respond(w, http.StatusOK, v1resources.LoginResponse{
		Encoded:  cookie,
		WSTicket: uuid,
	}, "successful login", "")
}

func getStatusAndErr(err error) (status int, errVal string) {
	if api_errors.IsMongoDup(err) {
		return http.StatusBadRequest, "content already registered"
	}
	if err == mongo.ErrNoDocuments {
		return http.StatusBadRequest, "content does not exist"
	}

	return http.StatusInternalServerError, api_errors.ErrGeneric.Error()
}
