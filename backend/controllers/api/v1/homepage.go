package controller

import (
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var templ *template.Template

func HomePageLoginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// if err := r.ParseForm(); err != nil {
	// 	log.Errorln("error parsing form in homepage login post request, err: ", err)
	// 	return
	// }

	// email := strings.ToLower(r.FormValue("email"))
	// password := r.FormValue("password")

	// err := model.User{
	// 	Email: email,
	// }.CreateUserLogin(password, w)

	// if err != nil {
	// 	if err := templ.ExecuteTemplate(w, "login.html", nil); err != nil {
	// 		log.Errorln("error loading login template", err)
	// 	}

	// 	return
	// }

	// http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func HomePageGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// cookie := model.CookieDetail{CookieName: values.UserCookieName, Collection: values.UsersCollectionName}

	// err := cookie.CheckCookie(r, values.UserCookieName)
	// if err != nil {
	// 	if err := templ.ExecuteTemplate(w, "login.html", nil); err != nil {
	// 		log.Errorln("error loading login template", err)
	// 	}
	// 	return
	// }

	// if err := templ.ExecuteTemplate(w, "home.html", nil); err != nil {
	// 	log.Errorln("error loading home template", err)
	// 	return
	// }
}
