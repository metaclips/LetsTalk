package test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"os"
	"testing"

	apihelpers "github.com/metaclips/LetsTalk/backend/apihelper"
	v1controller "github.com/metaclips/LetsTalk/backend/controllers/api/v1"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	appcontext "github.com/metaclips/LetsTalk/backend/middlewares/context"
	"github.com/metaclips/LetsTalk/backend/middlewares/cookies"
	"github.com/metaclips/LetsTalk/backend/middlewares/database"
	"github.com/metaclips/LetsTalk/backend/models"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	v1router "github.com/metaclips/LetsTalk/backend/router/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type handler struct {
	db         *mongo.Database
	ctx        context.Context
	cancelfunc context.CancelFunc

	usersModel  models.User
	roomRequest models.RoomRequest
	room        models.Room
	messages    models.Messages

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

	homeController := v1controller.NewHomepageController(usersModel, roomModel, roomRequestModel, usersCookie)
	loginController := v1controller.NewLoginController(usersModel, usersCookie)
	usersController := v1controller.NewRegistrationController(usersModel)

	router := v1router.NewRouter(homeController, loginController, usersController)

	testHandler.server = httptest.NewTLSServer(router)
	testHandler.client = testHandler.server.Client()

	testHandler.client.Jar, err = cookiejar.New(nil)
	panicOnError(err)
}

func (c *handler) seedOneUser(user *v1resources.UserData) error {
	collection := c.db.Collection(constants.UserCollection)
	_, err := collection.InsertOne(c.ctx, user)
	return err
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

func (c *handler) makeNormalCookieRequest(method, urlVal string, body interface{}, cookie *http.Cookie) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, urlVal, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.AddCookie(cookie)

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
