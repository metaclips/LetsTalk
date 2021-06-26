package v1controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	apihelpers "github.com/metaclips/LetsTalk/backend/apihelpers"
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

	if _, ok := maps.Departments[registration.Department]; !ok {
		apihelpers.Respond(w, http.StatusBadRequest, nil, "", "department not recognized")
		return
	}

	if _, ok := maps.Faculties[registration.Faculty]; !ok {
		apihelpers.Respond(w, http.StatusBadRequest, nil, "", "faculty not recognized")
		return
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
		DOB:        registration.DOB,
		Password:   hashedPassword,
		Department: registration.Department,
		Faculty:    registration.Faculty,
		IsStaff:    registration.IsStaff,
		Year:       registration.Year,
		CreatedAt:  time.Now(),
	}

	if err := r.user.CreateUser(r.ctx, &userData); err != nil {
		log.Errorln("error creating users:", err)
		apihelpers.Respond(w, http.StatusInternalServerError, nil, "", api_errors.ErrGeneric.Error())
		return
	}

	// ToDo: Send an email to users.

	apihelpers.Respond(w, http.StatusCreated, nil, "user successfully registered", "")
}
