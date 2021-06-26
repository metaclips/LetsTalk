package errors

import (
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrGeneric            = errors.New("something went wrong, please try again")
	ErrRecordExists       = errors.New("record already exists")
	ErrInvalidCredentials = errors.New("invalid credentials provided")
)

func IsMongoDup(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}

	return false
}

func ErrorResponse(err error) (int, string) {
	status := http.StatusInternalServerError
	errVal := ErrGeneric.Error()

	if IsMongoDup(err) {
		status = http.StatusConflict
		errVal = ErrRecordExists.Error()
	}

	if err == mongo.ErrNoDocuments {
		status = http.StatusBadRequest
		errVal = "record does not exit"
	}

	return status, errVal
}
