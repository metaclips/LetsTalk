package v1router

import (
	"net/http"
	"os"
	"os/signal"

	"github.com/julienschmidt/httprouter"
	apihelper "github.com/metaclips/LetsTalk/backend/apihelper"
	v1controller "github.com/metaclips/LetsTalk/backend/controllers/api/v1"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	appcontext "github.com/metaclips/LetsTalk/backend/middlewares/context"
	"github.com/metaclips/LetsTalk/backend/middlewares/cookies"
	"github.com/metaclips/LetsTalk/backend/middlewares/database"
	"github.com/metaclips/LetsTalk/backend/middlewares/peerconnection"
	appWebsocket "github.com/metaclips/LetsTalk/backend/middlewares/websocket"
	"github.com/metaclips/LetsTalk/backend/models"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func StartV1() {
	err := apihelper.LoadConfiguration("config.json")
	if err != nil {
		log.Fatalln("error loading config:", err)
	}

	ctx, _ := appcontext.CreateChildCancelContext()

	clientOptions := options.Client().ApplyURI(apihelper.Config.DbHost)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalln("error connecting mongo URL instance:", err, apihelper.Config.DbHost)
	}

	db := client.Database(apihelper.Config.DbName)

	go database.InitDatabaseIndexing(ctx, db)

	usersCollection := database.NewMongoStore(db, constants.UserCollection)
	roomRequestCollection := database.NewMongoStore(db, constants.RoomRequestCollection)
	roomCollection := database.NewMongoStore(db, constants.RoomCollection)
	messageCollection := database.NewMongoStore(db, constants.MessagesCollection)

	messageModel := models.NewMessages(messageCollection)
	roomModel := models.NewRoom(roomCollection)
	roomRequestModel := models.NewRoomRequest(roomRequestCollection)
	usersModel := models.NewUser(usersCollection)

	usersCookie := cookies.NewSecureCookie(usersCollection)

	pc := peerconnection.NewPionPeerConnection()
	mailer := v1resources.GoogleMailer{}

	apoWebsocket := appWebsocket.NewGorillaWebsocket(pc, roomRequestModel, messageModel, usersModel, roomModel)

	homeController := v1controller.NewHomepageController(messageModel, usersModel, roomModel, roomRequestModel, usersCookie)
	loginController := v1controller.NewLoginController(usersModel, usersCookie)
	usersController := v1controller.NewRegistrationController(usersModel, &mailer)
	wsController := v1controller.NewWebsocketController(usersModel, apoWebsocket)

	router := NewRouter(homeController, loginController, usersController, wsController)

	port := apihelper.Config.Port
	if port == "" {
		port = os.Getenv("PORT")
	}

	if port == "" {
		port = "8080"
	}

	handler := cors.New(cors.Options{
		Debug:          true,
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST"},
		AllowedHeaders: []string{"Content-Type", "Access-Control-Allow-Origin", "Authorization"},
	}).Handler(router)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	log.Infoln("Started Webserver On Port", port)
	go func() {
		if err := StartRouter(&srv); err != nil {
			panic(err)
		}

		log.Infoln("Webserver exited successfully")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Infoln("Shuting down Server ...")

	// Signal all children contexts to close.
	// This stops all database operations and queries.
	if err := appcontext.CancelContext(); err != nil {
		panic(err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Infoln("Server exiting")
}

func NewRouter(home *v1controller.HomepageController, login *v1controller.LoginController,
	registration *v1controller.RegistrationController, ws *v1controller.WebsocketController) *httprouter.Router {

	router := httprouter.New()

	router.POST("/v1/login", login.Login)
	router.POST("/v1/register", registration.RegisterUser)
	router.POST("/v1/verify/:email/:uuid", registration.VerifyUser)
	router.GET("/v1/home", home.HomePage)
	router.GET("/v1/ws/:userID/:ticket", ws.Websocket)

	router.GET("/v1/helpers/departments", v1controller.GetDepartments)
	router.GET("/v1/helpers/faculties", v1controller.GetFaculties)

	return router
}

// Starts httprouter on designated port.
// This is a blocking function.
func StartRouter(server *http.Server) error {
	log.Infoln(apihelper.Config.TLS.CertPath != "" && apihelper.Config.TLS.KeyPath != "")
	// Optional use of TLS due to Heroku serving TLS at low level.
	if apihelper.Config.TLS.CertPath != "" && apihelper.Config.TLS.KeyPath != "" {
		if err := server.ListenAndServeTLS(apihelper.Config.TLS.CertPath, apihelper.Config.TLS.KeyPath); err != nil {
			log.Fatalln(err)
		}
	}

	log.Infoln(apihelper.Config.TLS.CertPath != "" && apihelper.Config.TLS.KeyPath != "")
	// Note: without HTTPS users wont be able to login as SetCookie uses Secure flag.
	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}
