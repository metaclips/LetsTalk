package test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	apihelpers "github.com/metaclips/LetsTalk/backend/apihelper"
	v1controller "github.com/metaclips/LetsTalk/backend/controllers/api/v1"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	"github.com/metaclips/LetsTalk/backend/integration/mocker"
	appcontext "github.com/metaclips/LetsTalk/backend/middlewares/context"
	"github.com/metaclips/LetsTalk/backend/middlewares/cookies"
	"github.com/metaclips/LetsTalk/backend/middlewares/database"
	appWebsocket "github.com/metaclips/LetsTalk/backend/middlewares/websocket"
	"github.com/metaclips/LetsTalk/backend/models"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	v1router "github.com/metaclips/LetsTalk/backend/router/v1"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type handler struct {
	db         *mongo.Database
	ctx        context.Context
	cancelfunc context.CancelFunc

	usersModel  models.User
	roomRequest models.RoomRequest
	room        models.Room
	messages    models.Messages

	mocker *mocker.Mocker

	server *httptest.Server
	client *http.Client
}

var testHandler handler

func TestMain(m *testing.M) {
	integration()
	code := m.Run()
	os.Exit(code)
}

func integration() {
	err := apihelpers.LoadConfiguration("../../config_test.json")
	panicOnError(err)

	ctx, cancel := appcontext.CreateChildCancelContext()

	clientOptions := options.Client().ApplyURI(apihelpers.Config.DbHost)

	client, err := mongo.Connect(ctx, clientOptions)
	panicOnError(err)

	testHandler.db = client.Database(apihelpers.Config.DbName)
	testHandler.cancelfunc = cancel
	testHandler.ctx = ctx

	// drop the database from previous run
	testHandler.db.Drop(context.TODO())
	panicOnError(err)

	database.InitDatabaseIndexing(ctx, testHandler.db)

	usersCollection := database.NewMongoStore(testHandler.db, constants.UserCollection)
	roomRequestCollection := database.NewMongoStore(testHandler.db, constants.RoomRequestCollection)
	roomCollection := database.NewMongoStore(testHandler.db, constants.RoomCollection)
	messageCollection := database.NewMongoStore(testHandler.db, constants.MessagesCollection)

	messagesModel := models.NewMessages(messageCollection)
	roomModel := models.NewRoom(roomCollection)
	roomRequestModel := models.NewRoomRequest(roomRequestCollection)
	usersModel := models.NewUser(usersCollection)

	testHandler.messages = messagesModel
	testHandler.room = roomModel
	testHandler.roomRequest = roomRequestModel
	testHandler.usersModel = usersModel

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

	mocker := &mocker.Mocker{}
	apoWebsocket := appWebsocket.NewWebsocket(&wsHub, roomRequestModel, messagesModel, usersModel, roomModel)
	go apoWebsocket.Run()

	homeController := v1controller.NewHomepageController(messagesModel, usersModel, roomModel, roomRequestModel, usersCookie)
	loginController := v1controller.NewLoginController(usersModel, usersCookie)
	usersController := v1controller.NewRegistrationController(usersModel, mocker)
	wsController := v1controller.NewWebsocketController(usersModel, apoWebsocket)

	router := v1router.NewRouter(homeController, loginController, usersController, wsController)

	testHandler.server = httptest.NewServer(router)
	testHandler.client = testHandler.server.Client()

	testHandler.mocker = mocker

	panicOnError(err)
}

func (c *handler) seedOneUser(email string) (*v1resources.UserData, error) {
	password, err := bcrypt.GenerateFromPassword([]byte("password"), constants.BcryptCost)
	if err != nil {
		return nil, err
	}

	user := v1resources.UserData{
		ID:         primitive.NewObjectID(),
		Email:      email,
		FirstName:  "Michael",
		LastName:   "Uti",
		DOB:        time.Now(),
		Password:   password,
		Department: constants.ComputerEngineering,
		Faculty:    constants.EngineeringFaculty,
		IsStaff:    false,
		Verified:   true,
		Year:       2015,
		CreatedAt:  time.Now(),
	}

	collection := c.db.Collection(constants.UserCollection)
	_, err = collection.InsertOne(c.ctx, &user)
	return &user, err
}

func (c *handler) makeNormalRequest(method, urlVal string, body interface{}) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, urlVal, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func (c *handler) makeAuthRequest(method, urlVal, token string, body interface{}) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, urlVal, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "BEARER"+" "+token)

	return c.client.Do(req)
}

func (c *handler) getResponseData(respBody io.ReadCloser, data interface{}) error {
	responseData, err := ioutil.ReadAll(respBody)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(responseData, &data); err != nil {
		return err
	}

	return nil
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
