package v1controller

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	apihelpers "github.com/metaclips/LetsTalk/backend/apiHelpers"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	api_errors "github.com/metaclips/LetsTalk/backend/helpers/errors"
	app_context "github.com/metaclips/LetsTalk/backend/middlewares/context"
	"github.com/metaclips/LetsTalk/backend/middlewares/cookies"
	"github.com/metaclips/LetsTalk/backend/models"
	log "github.com/sirupsen/logrus"
)

type HomepageController struct {
	user       models.User
	ctx        context.Context
	cancelFunc context.CancelFunc
	cookie     cookies.Cookies
}

func NewHomepageController(user models.User, cookies cookies.Cookies) *LoginController {
	ctx, cancelFunc := app_context.CreateChildCancelContext()
	return &LoginController{
		user, ctx, cancelFunc, cookies,
	}
}

func (h HomepageController) HomePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, err := h.cookie.CheckCookie(h.ctx, r, constants.UserCookieName)
	if err != nil {
		log.Errorln("error checking users cookies: ", err)
		apihelpers.Respond(w, http.StatusBadRequest, "", api_errors.ErrValidatingCredentials.Error())
		return
	}

}
