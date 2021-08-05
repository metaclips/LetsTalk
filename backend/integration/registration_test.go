package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	apihelpers "github.com/metaclips/LetsTalk/backend/apihelper"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestUserRegistration(t *testing.T) {
	testHandler.mocker.Mailer = func(email, uuid string) error {
		return nil
	}

	for _, user := range []struct {
		user   v1resources.UserRegistration
		name   string
		err    string
		status int
	}{
		{
			user: v1resources.UserRegistration{
				Email:         "utimichael@gmail.com",
				FirstName:     "Michael",
				LastName:      "Uti",
				Department:    constants.ComputerEngineering,
				FirebaseToken: "Firebase",
				Faculty:       constants.EngineeringFaculty,
				Password:      "password",
				Year:          2015,
				IsStaff:       false,
				DOB:           time.Now(),
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
				Email:         "utimichael@unilag.edu.ng",
				FirstName:     "Michael",
				LastName:      "Uti",
				Department:    "Null",
				Faculty:       constants.EngineeringFaculty,
				FirebaseToken: "Firebase",
				Password:      "password",
				Year:          2015,
				IsStaff:       false,
				DOB:           time.Now(),
			},
			name:   "invalid department",
			err:    "department not recognized",
			status: http.StatusBadRequest,
		},

		{
			user: v1resources.UserRegistration{
				Email:         "utimichael0J000@student.unilag.edu.ng",
				FirstName:     "Michael",
				LastName:      "Uti",
				Department:    constants.ComputerEngineering,
				Faculty:       constants.EngineeringFaculty,
				FirebaseToken: "Firebase",
				Password:      "password",
				Year:          2015,
				IsStaff:       false,
				DOB:           time.Now(),
			},
			name:   "Accurate",
			err:    "",
			status: http.StatusCreated,
		},
	} {
		resp, err := testHandler.makeNormalRequest("POST", fmt.Sprintf("%s/v1/register", testHandler.server.URL), user.user)
		assert.NoError(t, err)
		assert.Equal(t, user.status, resp.StatusCode, user.name)

		body := apihelpers.Response{}
		err = testHandler.getResponseData(resp.Body, &body)
		assert.NoError(t, err)

		assert.Equal(t, body.Error, user.err, user.name)

		assert.Equal(t, user.err, body.Error, user.name)
		assert.Equal(t, user.status, body.StatusCode, user.name)
	}
}

func TestDuplicateRegistration(t *testing.T) {
	matric := "150408002"

	user := v1resources.UserData{
		ID:            primitive.NewObjectID(),
		Email:         "utimichael1111@student.unilag.edu.ng",
		FirstName:     "Michael",
		LastName:      "Uti",
		MatricNumber:  &matric,
		Department:    constants.ComputerEngineering,
		Faculty:       constants.EngineeringFaculty,
		Password:      []byte("password"),
		FirebaseToken: "Firebase",
		Year:          2015,
		IsStaff:       false,
		DOB:           time.Now(),
	}

	err := testHandler.usersModel.CreateUser(testHandler.ctx, &user)
	assert.NoError(t, err)

	user2 := v1resources.UserData{
		ID:            primitive.NewObjectID(),
		Email:         "utimichael1111@student.unilag.edu.ng",
		FirstName:     "Michael",
		LastName:      "Uti",
		Department:    constants.ComputerEngineering,
		Faculty:       constants.EngineeringFaculty,
		FirebaseToken: "Firebase",
		Password:      []byte("password"),
		Year:          2015,
		IsStaff:       false,
		DOB:           time.Now(),
	}

	err = testHandler.usersModel.CreateUser(testHandler.ctx, &user2)
	assert.Error(t, err)

	user3 := v1resources.UserData{
		ID:            primitive.NewObjectID(),
		Email:         "utimichael1112@student.unilag.edu.ng",
		FirstName:     "Michael",
		LastName:      "Uti",
		Department:    constants.ComputerEngineering,
		Faculty:       constants.EngineeringFaculty,
		FirebaseToken: "Firebase",
		Password:      []byte("password"),
		Year:          2015,
		IsStaff:       false,
		DOB:           time.Now(),
	}

	err = testHandler.usersModel.CreateUser(testHandler.ctx, &user3)
	assert.NoError(t, err)

	user4 := v1resources.UserData{
		ID:            primitive.NewObjectID(),
		Email:         "utimichael1113@student.unilag.edu.ng",
		FirstName:     "Michael",
		LastName:      "Uti",
		FirebaseToken: "Firebase",
		Department:    constants.ComputerEngineering,
		Faculty:       constants.EngineeringFaculty,
		Password:      []byte("password"),
		Year:          2015,
		IsStaff:       false,
		DOB:           time.Now(),
	}

	err = testHandler.usersModel.CreateUser(testHandler.ctx, &user4)
	assert.NoError(t, err)

	user5 := v1resources.UserData{
		ID:            primitive.NewObjectID(),
		Email:         "utimichael1114@student.unilag.edu.ng",
		FirstName:     "Michael",
		LastName:      "Uti",
		FirebaseToken: "Firebase",
		MatricNumber:  &matric,
		Department:    constants.ComputerEngineering,
		Faculty:       constants.EngineeringFaculty,
		Password:      []byte("password"),
		Year:          2015,
		IsStaff:       false,
		DOB:           time.Now(),
	}

	err = testHandler.usersModel.CreateUser(testHandler.ctx, &user5)
	assert.Error(t, err)
}
