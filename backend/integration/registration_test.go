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
)

func TestUserRegistration(t *testing.T) {
	for _, user := range []struct {
		user   v1resources.UserRegistration
		name   string
		err    string
		status int
	}{
		{
			user: v1resources.UserRegistration{
				Email:      "utimichael@gmail.com",
				FirstName:  "Michael",
				LastName:   "Uti",
				Department: constants.ComputerEngineering,
				Faculty:    constants.EngineeringFaculty,
				Password:   "password",
				Year:       2015,
				IsStaff:    false,
				DOB:        time.Now(),
			},
			name:   "unsupported email address",
			err:    "invalid student email",
			status: http.StatusBadRequest,
		},

		{
			user:   v1resources.UserRegistration{},
			name:   "empty body",
			err:    "email is required",
			status: http.StatusBadRequest,
		},

		{
			user: v1resources.UserRegistration{
				Email:      "utimichael@unilag.edu.ng",
				FirstName:  "Michael",
				LastName:   "Uti",
				Department: constants.ComputerEngineering,
				Faculty:    constants.EngineeringFaculty,
				Password:   "password",
				Year:       2015,
				IsStaff:    false,
				DOB:        time.Now(),
			},
			name:   "unsupported email address",
			err:    "invalid student email",
			status: http.StatusBadRequest,
		},

		{
			user: v1resources.UserRegistration{
				Email:      "utimichael@student.unilag.edu.ng",
				FirstName:  "Michael",
				LastName:   "Uti",
				Department: constants.ComputerEngineering,
				Faculty:    constants.EngineeringFaculty,
				Password:   "password",
				Year:       2015,
				IsStaff:    false,
				DOB:        time.Now(),
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
	}

}
