package controller

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"

	"github.com/metaclips/LetsTalk/backend/model"
	"github.com/metaclips/LetsTalk/backend/values"
)

var templ *template.Template

func init() {
	var terr error
	templ, terr = template.New("home.html").ParseFiles(os.ExpandEnv("./dist/home.html"),
		os.ExpandEnv("./dist/login.html"), os.ExpandEnv("./dist/register.html"))

	if terr != nil {
		log.Fatalln("error parsing homepage templates", terr)
	}
}

func HomePageLoginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := r.ParseForm(); err != nil {
		log.Errorln("error parsing form in homepage login post request, err: ", err)
		return
	}

	email := strings.ToLower(r.FormValue("username"))
	password := r.FormValue("password")

	err := model.User{
		Email: email,
	}.CreateUserLogin(password, w)

	if err != nil {
		if err == values.ErrRetrieveUser {
			http.Error(w, "invalid login credentials", http.StatusBadRequest)
			return
		}

		http.Error(w, "error creating login", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HomePageLoginGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie := model.CookieDetail{CookieName: values.UserCookieName, Collection: values.UsersCollectionName}

	err := cookie.CheckCookie(r, values.UserCookieName)
	if err != nil {
		if err := templ.ExecuteTemplate(w, "login.html", nil); err != nil {
			log.Fatalln(err)
		}
		return
	}

	// Redirect to homepage if cookie already exists
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func sendResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Errorf("error sending response, err: %s \n", err.Error())
	}

}
