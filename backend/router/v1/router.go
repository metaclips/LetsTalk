package v1router

import (
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	apihelper "github.com/metaclips/LetsTalk/backend/apihelpers"
	v1controller "github.com/metaclips/LetsTalk/backend/controllers/api/v1"
)

func NewRouter(home *v1controller.HomepageController, login *v1controller.LoginController,
	registration *v1controller.RegistrationController) *httprouter.Router {

	router := httprouter.New()

	router.POST("/v1/login", login.Login)
	router.POST("/v1/register", registration.RegisterUser)
	router.GET("/v1/home", home.HomePage)

	return router
}

// Starts httprouter on designated port.
// This is a blocking function.
func StartRouter(router *httprouter.Router) error {
	port := apihelper.Config.Port
	if port == "" {
		port = os.Getenv("PORT")
	}

	if port == "" {
		port = "8080"
	}

	// Optional use of TLS due to Heroku serving TLS at low level.
	if apihelper.Config.TLS.CertPath != "" && apihelper.Config.TLS.KeyPath != "" {
		if err := http.ListenAndServeTLS(":"+port, apihelper.Config.TLS.CertPath, apihelper.Config.TLS.KeyPath, router); err != nil {
			return err
		}
	}

	// Note: without HTTPS users wont be able to login as SetCookie uses Secure flag.
	if err := http.ListenAndServe(":"+port, router); err != nil {
		return err
	}

	return nil
}
