package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	apihelpers "github.com/metaclips/LetsTalk/backend/apihelper"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	api_errors "github.com/metaclips/LetsTalk/backend/helpers/errors"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	password, err := bcrypt.GenerateFromPassword([]byte("password"), constants.BcryptCost)
	assert.NoError(t, err)

	email := "michaelnew1@students.unilag.edu.ng"
	userData := v1resources.UserData{
		ID:         primitive.NewObjectID(),
		Email:      email,
		FirstName:  "Michael",
		LastName:   "Uti",
		DOB:        time.Now(),
		Password:   password,
		Department: constants.ComputerEngineering,
		Faculty:    constants.EngineeringFaculty,
		IsStaff:    false,
		Year:       2015,
		CreatedAt:  time.Now(),
	}

	err = testHandler.seedOneUser(&userData)
	assert.NoError(t, err)

	for _, user := range []struct {
		user   v1resources.UserLogin
		name   string
		err    string
		status int
	}{
		{
			user: v1resources.UserLogin{
				Email:    "test",
				Password: "password",
			},
			name:   "unsupported email address",
			err:    "email should be a valid email",
			status: http.StatusBadRequest,
		},

		{
			user:   v1resources.UserLogin{},
			name:   "empty body",
			err:    "email should be a valid email",
			status: http.StatusBadRequest,
		},

		{
			user: v1resources.UserLogin{
				Email:    email,
				Password: "wrongPassword",
			},
			name:   "wrong password",
			err:    api_errors.ErrInvalidCredentials.Error(),
			status: http.StatusBadRequest,
		},

		{
			user: v1resources.UserLogin{
				Email:    email,
				Password: "password",
			},
			name:   "correct impl",
			err:    "",
			status: http.StatusOK,
		},
	} {
		resp, err := testHandler.makeNormalRequest("POST", fmt.Sprintf("%s/v1/login", testHandler.server.URL), user.user)
		assert.NoError(t, err, user.name)
		assert.Equal(t, user.status, resp.StatusCode, user.name)

		body := apihelpers.Response{}

		err = testHandler.getResponseData(resp.Body, &body)
		assert.NoError(t, err, user.name, user.name)

		assert.Equal(t, user.err, body.Error, user.name)
		assert.Equal(t, user.status, body.StatusCode, user.name)
		assert.NotNil(t, body.Data)
	}
}
