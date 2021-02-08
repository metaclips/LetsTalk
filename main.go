package main

import (
	"encoding/gob"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
	"github.com/metaclips/LetsTalk/backend/controller"
	"github.com/metaclips/LetsTalk/backend/model"
	"github.com/metaclips/LetsTalk/backend/values"
)

func main() {
	if err := values.LoadConfiguration("./config.json"); err != nil {
		log.Fatalln("unable to load config", err)
	}

	gob.Register(time.Time{})
	model.InitDB()
	go model.HubConstruct.Run()

	router := httprouter.New()

	router.GET("/ws", controller.ServeWs)
	router.GET("/login", controller.HomePageLoginGet)

	router.POST("/login", controller.HomePageLoginPost)
	router.POST("/register", controller.RegisterUserPost)

	port := values.Config.Port
	if port == "" {
		port = os.Getenv("PORT")
	}
	if port == "" {
		port = "8081"
	}

	router.ServeFiles("/img/*filepath", http.Dir("./dist/img"))
	router.ServeFiles("/js/*filepath", http.Dir("./dist/js"))
	router.ServeFiles("/css/*filepath", http.Dir("./dist/css"))
	router.ServeFiles("/favicon.ico/*filepath", http.Dir("./dist"))

	log.Infoln("Webserver UP")

	// Optional use of TLS due to Heroku serving TLS at low level.
	if values.Config.TLS.CertPath != "" && values.Config.TLS.KeyPath != "" {
		if err := http.ListenAndServeTLS(":"+port, values.Config.TLS.CertPath, values.Config.TLS.KeyPath, router); err != nil {
			log.Fatalln(err)
		}
		return
	}

	// Note: without HTTPS users wont be able to login as SetCookie uses Secure flag.
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalln(err)
	}
}
