package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	apihelpers "github.com/metaclips/LetsTalk/backend/apihelpers"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
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
			err:    "invalid student email",
			status: http.StatusBadRequest,
		},

		{
			user:   v1resources.UserLogin{},
			name:   "empty body",
			err:    "email is required",
			status: http.StatusBadRequest,
		},

		{
			user: v1resources.UserLogin{
				Email:    email,
				Password: "wrongPassword",
			},
			name:   "unsupported email address",
			err:    "invalid student email",
			status: http.StatusBadRequest,
		},

		{
			user: v1resources.UserLogin{
				Email:    email,
				Password: "password",
			},
			name:   "unsupported email address",
			err:    "invalid student email",
			status: http.StatusBadRequest,
		},
	} {
		resp, err := testHandler.makeNormalRequest("POST", fmt.Sprintf("%s/v1/register", testHandler.server.URL), user.user)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		body := apihelpers.Response{}

		err = testHandler.getResponseData(resp.Body, &body)
		assert.NoError(t, err)

		assert.Equal(t, user.err, body.Error)
		assert.Equal(t, user.status, body.StatusCode)

		if assert.Equal(t, 1, len(resp.Cookies())) {
			resp, err = testHandler.makeNormalCookieRequest("GET", fmt.Sprintf("%s/v1/home", testHandler.server.URL), nil, resp.Cookies()[0])
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		}
	}
}
