package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/metaclips/LetsTalk/api/model"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterUserPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		log.Errorln("error parsing form in register user post request, err: ", err)
	}

	email := strings.ToLower(r.FormValue("email"))
	password := r.FormValue("password")
	name := r.FormValue("name")
	DOB := r.FormValue("DOB")

	err := model.User{
		Email:            email,
		PasswordInString: password,
		Name:             name,
		DOB:              DOB,
	}.CreateUser(r)

	if err != nil {
		if isMongoDup(err) {
			http.Error(w, "Email already registered", http.StatusUnauthorized)
		} else {
			http.Error(w, "Internal error", http.StatusUnauthorized)
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}

func isMongoDup(err error) bool {
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
