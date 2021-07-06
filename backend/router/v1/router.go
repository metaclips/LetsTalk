package v1router

import (
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/julienschmidt/httprouter"
	apihelper "github.com/metaclips/LetsTalk/backend/apihelper"
	v1controller "github.com/metaclips/LetsTalk/backend/controllers/api/v1"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	appcontext "github.com/metaclips/LetsTalk/backend/middlewares/context"
	"github.com/metaclips/LetsTalk/backend/middlewares/cookies"
	"github.com/metaclips/LetsTalk/backend/middlewares/database"
	"github.com/metaclips/LetsTalk/backend/models"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func StartV1() {
	ctx, _ := appcontext.CreateChildCancelContext()

	clientOptions := options.Client().ApplyURI(apihelper.Config.DbHost)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalln("error connecting mongo instance:", err)
	}

	db := client.Database(apihelper.Config.DbName)

	database.InitDatabaseIndexing(ctx, db)

	usersCollection := database.NewMongoStore(db, constants.UserCollection)
	roomRequestCollection := database.NewMongoStore(db, constants.RoomRequestCollection)
	roomCollection := database.NewMongoStore(db, constants.RoomCollection)
	messageCollection := database.NewMongoStore(db, constants.MessagesCollection)

	messageModel := models.NewMessages(messageCollection)
	roomModel := models.NewRoom(roomCollection)
	roomRequestModel := models.NewRoomRequest(roomRequestCollection)
	usersModel := models.NewUser(usersCollection)

	usersCookie := cookies.NewSecureCookie(usersCollection)

	wsHub := v1resources.Hub{
		Broadcast:  make(chan v1resources.WSMessage),
		Register:   make(chan v1resources.Subscription),
		UnRegister: make(chan v1resources.Subscription),
		Users: v1resources.WSUsers{
			RoomUsers: make(map[primitive.ObjectID]*v1resources.Connection),
			Mutex:     &sync.RWMutex{},
		},
	}

	homeController := v1controller.NewHomepageController(messageModel, usersModel, roomModel, roomRequestModel, usersCookie)
	loginController := v1controller.NewLoginController(usersModel, usersCookie)
	usersController := v1controller.NewRegistrationController(usersModel)
	wsController := v1controller.NewWebsocketController(&wsHub, messageModel, usersModel, roomModel, roomRequestModel, usersCookie)

	go wsController.Run()

	router := NewRouter(homeController, loginController, usersController, wsController)
	srv := http.Server{}

	go func() {
		if err := StartRouter(router); err != nil {
			panic(err)
		}

		log.Infoln("Webserver exited successfully")
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Infoln("Shutdown Server ...")

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// Signal all children contexts to close.
	// This stops all database operations and queries.
	if err := appcontext.CancelContext(); err != nil {
		panic(err)
	}

	log.Infoln("Server exiting")
}

func NewRouter(home *v1controller.HomepageController, login *v1controller.LoginController,
	registration *v1controller.RegistrationController, ws *v1controller.WebsocketController) *httprouter.Router {

	router := httprouter.New()

	router.POST("/v1/login", login.Login)
	router.POST("/v1/register", registration.RegisterUser)
	router.GET("/v1/home", home.HomePage)
	router.GET("/v1/ws/:userID/:ticket", ws.Websocket)

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
