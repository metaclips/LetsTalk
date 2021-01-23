package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"

	"github.com/metaclips/LetsTalk/api/model"
)

func HomePageLoginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		log.Errorln("error parsing form in homepage login post request, err: ", err)
		return
	}

	email := strings.ToLower(r.FormValue("username"))
	password := r.FormValue("password")

	token, err := model.User{
		Email: email,
	}.CreateUserLogin(password, w)

	if err != nil {
		http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	sendResponse(w, token)
}

func sendResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Errorf("error sending response, err: %s \n", err.Error())
	}

}
