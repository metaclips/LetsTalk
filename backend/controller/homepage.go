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
	templ, terr = template.New("templates").ParseFiles(os.ExpandEnv("./dist/home.html"),
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

	email := strings.ToLower(r.FormValue("email"))
	password := r.FormValue("password")

	err := model.User{
		Email: email,
	}.CreateUserLogin(password, w)

	if err != nil {
		if err := templ.ExecuteTemplate(w, "login.html", nil); err != nil {
			log.Errorln("error loading login template", err)
		}

		return
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func LoginPageGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie := model.CookieDetail{CookieName: values.UserCookieName, Collection: values.UsersCollectionName}

	err := cookie.CheckCookie(r, values.UserCookieName)
	if err != nil {
		if err := templ.ExecuteTemplate(w, "login.html", nil); err != nil {
			log.Errorln("error loading login template", err)
		}
		return
	}

	// Redirect to homepage if cookie already exists
	http.Redirect(w, r, "/", http.StatusPermanentRedirect)
}

func RegisterPageGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if err := templ.ExecuteTemplate(w, "register.html", nil); err != nil {
		log.Errorln("error loading registration template", err)
	}
	return
}

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
		if err := templ.ExecuteTemplate(w, "register.html", nil); err != nil {
			log.Errorln("error loading registration template", err)
		}
		return
	}

	if err := templ.ExecuteTemplate(w, "login.html", nil); err != nil {
		log.Errorln("error loading registration template", err)
	}
}

func HomePageGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	cookie := model.CookieDetail{CookieName: values.UserCookieName, Collection: values.UsersCollectionName}

	err := cookie.CheckCookie(r, values.UserCookieName)
	if err != nil {
		if err := templ.ExecuteTemplate(w, "login.html", nil); err != nil {
			log.Errorln("error loading login template", err)
		}
		return
	}

	if err := templ.ExecuteTemplate(w, "home.html", nil); err != nil {
		log.Errorln("error loading home template", err)
		return
	}
}

func sendResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Errorf("error sending response, err: %s \n", err.Error())
	}
}
