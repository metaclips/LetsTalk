package test

import (
	"fmt"
	"net/http"
	"testing"

	apihelpers "github.com/metaclips/LetsTalk/backend/apihelper"
	api_errors "github.com/metaclips/LetsTalk/backend/helpers/errors"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	email := "michaelnew1@students.unilag.edu.ng"

	_, err := testHandler.seedOneUser(email)
	assert.NoError(t, err)

	for _, user := range []struct {
		user   v1resources.UserLogin
		name   string
		err    string
		status int
	}{
		{
			user: v1resources.UserLogin{
				Email:    email,
				Password: "password",
			},
			name:   "correct impl",
			err:    "",
			status: http.StatusOK,
		},

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
			err:    "email is required",
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
	} {
		resp, err := testHandler.makeNormalRequest("POST", fmt.Sprintf("%s/v1/login", testHandler.server.URL), user.user)
		assert.NoError(t, err, user.name)
		assert.Equal(t, user.status, resp.StatusCode, user.name)

		body := apihelpers.Response{}

		err = testHandler.getResponseData(resp.Body, &body)
		assert.NoError(t, err, user.name, user.name)

		assert.Equal(t, user.err, body.Error, user.name)
		assert.Equal(t, user.status, body.StatusCode, user.name)
		if body.StatusCode == http.StatusOK {
			assert.NotNil(t, body.Data, user.name)
		}
	}
}
