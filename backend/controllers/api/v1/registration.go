package v1controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	apihelpers "github.com/metaclips/LetsTalk/backend/apihelper"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	api_errors "github.com/metaclips/LetsTalk/backend/helpers/errors"
	"github.com/metaclips/LetsTalk/backend/helpers/maps"
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
	user   models.User
	mailer v1resources.Mailer

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewRegistrationController(user models.User, mailer v1resources.Mailer) *RegistrationController {
	ctx, cancelFunc := app_context.CreateChildCancelContext()
	return &RegistrationController{
		user, mailer, ctx, cancelFunc,
	}
}

func (r *RegistrationController) RegisterUser(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	registration := v1resources.UserRegistration{}
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&registration); err != nil {
		log.Errorln("error decoding registration data: ", err)
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	if err := validation.ValidateInputs(registration); err != nil {
		apihelpers.Respond(w, http.StatusBadRequest, nil, "", err.Error())
		return
	}

	// Ensure email address is a school email.
	words := strings.Split(registration.Email, "@")
	if len(words) < 2 {
		apihelpers.Respond(w, http.StatusBadRequest, nil, "", "invalid student email")
		return
	}

	eduContainer := strings.Split(words[1], ".")
	// School email should contain school name @ edu . country.
	if len(eduContainer) < 3 || eduContainer[len(eduContainer)-2] != "edu" {
		apihelpers.Respond(w, http.StatusBadRequest, nil, "", "invalid student email")
		return
	}

	if !registration.IsStaff {
		if _, ok := maps.Departments[registration.Department]; !ok {
			apihelpers.Respond(w, http.StatusBadRequest, nil, "", "department not recognized")
			return
		}

		if _, ok := maps.Faculties[registration.Faculty]; !ok {
			apihelpers.Respond(w, http.StatusBadRequest, nil, "", "faculty not recognized")
			return
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registration.Password), constants.BcryptCost)
	if err != nil {
		log.Errorln("error hashing users password:", err)
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	userData := v1resources.UserData{
		ID:         primitive.NewObjectID(),
		Email:      registration.Email,
		FirstName:  registration.FirstName,
		LastName:   registration.LastName,
		DOB:        time.Time(registration.DOB),
		Password:   hashedPassword,
		Department: registration.Department,
		Faculty:    registration.Faculty,
		IsStaff:    registration.IsStaff,
		Year:       registration.Year,
		UUID:       uuid.New().String(),
		CreatedAt:  time.Now(),
	}

	if registration.MatricNumber != "" {
		userData.MatricNumber = &registration.MatricNumber
	}

	if err := r.mailer.MailTo(userData.Email, userData.UUID); err != nil {
		log.Errorln("error sending verification email:", err)
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", "error sending verification email")
		return
	}

	if err := r.user.CreateUser(r.ctx, &userData); err != nil {
		log.Errorln("error creating users:", err)
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	apihelpers.Respond(w, http.StatusCreated, nil, "user successfully registered", "")
}

func (r *RegistrationController) VerifyUser(w http.ResponseWriter, req *http.Request, param httprouter.Params) {
	userUUID := param.ByName("uuid")
	email := param.ByName("email")
	userData, err := r.user.GetUserByUUID(r.ctx, userUUID)
	if err != nil {
		log.Errorln("error getting users details on verify user:", err)
		status, err := getStatusAndErr(err)
		apihelpers.Respond(w, status, nil, "", err)
		return
	}

	if userData.Email != email {
		apihelpers.Respond(w, http.StatusBadRequest, nil, "", "user not associated")
		return
	}

	userData.Verified = true
	if err := r.user.UpdateUser(r.ctx, userData); err != nil {
		log.Errorln("error updating users details on verify user:", err)
		status, err := getStatusAndErr(err)
		apihelpers.Respond(w, status, nil, "", err)
		return
	}

	apihelpers.Respond(w, http.StatusOK, nil, "user successfully verified", "")
}
