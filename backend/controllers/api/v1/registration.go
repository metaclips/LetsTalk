package v1controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	apihelpers "github.com/metaclips/LetsTalk/backend/apiHelpers"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	api_errors "github.com/metaclips/LetsTalk/backend/helpers/errors"
	"github.com/metaclips/LetsTalk/backend/helpers/validation"
	app_context "github.com/metaclips/LetsTalk/backend/middlewares/context"
	"github.com/metaclips/LetsTalk/backend/models"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"

	"github.com/julienschmidt/httprouter"
)

type RegistrationController struct {
	user       models.User
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewRegistrationController(user models.User) *RegistrationController {
	ctx, cancelFunc := app_context.CreateChildCancelContext()
	return &RegistrationController{
		user, ctx, cancelFunc,
	}
}

func (r RegistrationController) RegisterUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	registration := v1resources.UserRegistration{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&registration); err != nil {
		log.Errorln("error decoding registration data: ", err)
		apihelpers.Respond(w, http.StatusInternalServerError, "", api_errors.ErrGeneric.Error())
		return
	}

	if err := validation.ValidateInputs(registration); err != nil {
		apihelpers.Respond(w, http.StatusBadRequest, "", err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registration.Password), constants.BcryptCost)
	if err != nil {
		log.Errorln("error hashing users password:", err)
		apihelpers.Respond(w, http.StatusInternalServerError, "", api_errors.ErrGeneric.Error())
		return
	}

	userData := v1resources.UserData{
		ID:        primitive.NewObjectID(),
		Email:     registration.Email,
		FirstName: registration.FirstName,
		LastName:  registration.LastName,
		DOB:       registration.DOB,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
	}

	if err := r.user.CreateUser(r.ctx, &userData); err != nil {
		log.Errorln("error creating users:", err)
		apihelpers.Respond(w, http.StatusInternalServerError, "", api_errors.ErrGeneric.Error())
		return
	}
}
