package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	apihelpers "github.com/metaclips/LetsTalk/backend/apihelper"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateNewRoom(t *testing.T) {
	email := "michaelnew2@students.unilag.edu.ng"

	userData, err := testHandler.seedOneUser(email)
	assert.NoError(t, err)

	ws, err := createWSConn(email)
	assert.NoError(t, err)
	if ws == nil {
		t.FailNow()
	}

	msg := struct {
		MsgType string `json:"msgType"`
		v1resources.Room
	}{
		MsgType: constants.CreateRoomMsgType,
		Room: v1resources.Room{
			CreateByUserID: userData.ID,
			RoomName:       "My New Room",
		},
	}

	err = ws.WriteJSON(msg)
	assert.NoError(t, err)

	result := v1resources.Room{}
	err = ws.ReadJSON(&result)
	assert.NoError(t, err)

	assert.Equal(t, msg.Room.RoomName, result.RoomName)

	createRoomMsg := &v1resources.Message{}
	err = ws.ReadJSON(createRoomMsg)
	assert.NoError(t, err)
	assert.NotNil(t, createRoomMsg)

	joinedRooms, err := testHandler.roomRequest.RetrieveAcceptedJoinRequest(testHandler.ctx, userData.ID, 0, 0)
	assert.NoError(t, err)
	if assert.NotNil(t, joinedRooms) {
		assert.GreaterOrEqual(t, 1, len(joinedRooms.Data))

		hasCreatedRoom := false
		for _, joinedRoom := range joinedRooms.Data {
			if joinedRoom.RoomName == msg.RoomName {
				hasCreatedRoom = true
			}
		}

		assert.True(t, hasCreatedRoom)
	}
}

func TestAcceptRoomRequest(t *testing.T) {
	email := "michaelnew3@students.unilag.edu.ng"

	userData, err := testHandler.seedOneUser(email)
	assert.NoError(t, err)

	ws, err := createWSConn(email)
	assert.NoError(t, err)
	if ws == nil {
		t.FailNow()
	}

	email2 := "michaelnew4@students.unilag.edu.ng"

	userData2, err := testHandler.seedOneUser(email2)
	assert.NoError(t, err)

	ws2, err := createWSConn(email2)
	assert.NoError(t, err)
	if ws == nil {
		t.FailNow()
	}

	msg := struct {
		MsgType string `json:"msgType"`
		v1resources.Room
	}{
		MsgType: constants.CreateRoomMsgType,
		Room: v1resources.Room{
			CreateByUserID: userData.ID,
			RoomName:       "My New Room",
		},
	}

	err = ws.WriteJSON(msg)
	assert.NoError(t, err)

	result := v1resources.Room{}
	err = ws.ReadJSON(&result)
	assert.NoError(t, err)

	// Send room request to user
	msg2 := struct {
		MsgType string `json:"msgType"`
		v1resources.JoinRoomRequest
	}{
		MsgType: constants.RequestUsersToJoinRoomMsgType,
		JoinRoomRequest: v1resources.JoinRoomRequest{
			UserID:       userData.ID,
			UserToJoinID: userData2.ID,
			RoomID:       result.ID,
			RoomName:     "My New Room",
		},
	}
	// Send room request to user.
	err = ws.WriteJSON(msg2)
	assert.NoError(t, err)

	requestNotif := v1resources.JoinRoomRequest{}
	err = ws2.ReadJSON(&requestNotif)
	assert.NoError(t, err)

	msg3 := struct {
		MsgType string `json:"msgType"`
		v1resources.JoinRoomRequest
	}{
		MsgType:         constants.JoinRoomMsgType,
		JoinRoomRequest: requestNotif,
	}

	err = ws2.WriteJSON(msg3)
	assert.NoError(t, err)

	newRoomNotif := &v1resources.Room{}
	err = ws2.ReadJSON(newRoomNotif)
	assert.NoError(t, err)

	assert.NotNil(t, newRoomNotif)
	assert.Equal(t, "My New Room", newRoomNotif.RoomName)

	roomMessage := &v1resources.Message{}
	ws2.ReadJSON(roomMessage)
	assert.NotNil(t, roomMessage)

	assert.Equal(t, fmt.Sprintf("%s %s (%s) has has joined the room.", userData2.FirstName, userData2.LastName, userData2.Email), roomMessage.Message)
}

func TestGetRoomRequestOnLogin(t *testing.T) {
	email := "michaelnew52@students.unilag.edu.ng"

	userData, err := testHandler.seedOneUser(email)
	assert.NoError(t, err)

	acceptedRoomID := primitive.NewObjectID()
	unacceptedRoomID := primitive.NewObjectID()
	roomID := primitive.NewObjectID()

	err = testHandler.roomRequest.NewRoomRequest(testHandler.ctx, &v1resources.JoinRoomRequest{
		ID:           primitive.NewObjectID(),
		UserID:       primitive.NewObjectID(),
		UserToJoinID: userData.ID,
		RoomID:       unacceptedRoomID,
		RoomName:     "Test Room",
		CreatedAt:    time.Now(),
	})
	assert.NoError(t, err)

	err = testHandler.room.CreateNewRoom(testHandler.ctx, &v1resources.Room{
		ID:             roomID,
		CreateByUserID: primitive.NewObjectID(),
		RoomName:       "Test Room",
		RoomIcon:       "",
		CreatedAt:      time.Now(),
	})
	assert.NoError(t, err)

	err = testHandler.roomRequest.NewRoomRequest(testHandler.ctx, &v1resources.JoinRoomRequest{
		ID:           acceptedRoomID,
		UserID:       primitive.NewObjectID(),
		UserToJoinID: userData.ID,
		RoomID:       roomID,
		RoomName:     "Test Room Accepted",
		Accepted:     true,
		CreatedAt:    time.Now(),
	})
	assert.NoError(t, err)

	login := v1resources.UserLogin{
		Email:    userData.Email,
		Password: "password",
	}

	resp, err := testHandler.makeNormalRequest("POST", fmt.Sprintf("%s/v1/login", testHandler.server.URL), login)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	loginResponse := struct {
		Data       v1resources.LoginResponse `json:"data"`
		Msg        string                    `json:"msg"`
		Error      string                    `json:"error"`
		StatusCode int                       `json:"statusCode"`
	}{}

	err = testHandler.getResponseData(resp.Body, &loginResponse)
	assert.NoError(t, err)

	resp, err = testHandler.makeAuthRequest("GET", fmt.Sprintf("%s/v1/home", testHandler.server.URL), loginResponse.Data.Token, nil)
	assert.NoError(t, err)

	homeResponse := struct {
		Data       v1resources.HomeResponse `json:"data"`
		Msg        string                   `json:"msg"`
		Error      string                   `json:"error"`
		StatusCode int                      `json:"statusCode"`
	}{}

	err = testHandler.getResponseData(resp.Body, &homeResponse)
	assert.NoError(t, err)

	if assert.Equal(t, 1, homeResponse.Data.RoomRequestList.Count) {
		assert.Equal(t, unacceptedRoomID, homeResponse.Data.RoomRequestList.Data[0].RoomID)
	}

	if assert.Equal(t, 1, homeResponse.Data.JoinedRoomList.Count) {
		assert.Equal(t, roomID, homeResponse.Data.JoinedRoomList.Data[0].ID)
	}
}

func TestPredatorScan(t *testing.T) {
	roomID := primitive.NewObjectID()
	predator := primitive.NewObjectID()
	normalUser := primitive.NewObjectID()

	minorsDOB, err := time.Parse(constants.DateOnly, "2006-01-02")
	assert.NoError(t, err)

	err = testHandler.usersModel.CreateUser(testHandler.ctx, &v1resources.UserData{
		ID:         normalUser,
		Email:      "anormaluser@students.unilag.edu.ng",
		FirstName:  "Michael",
		LastName:   "Uti",
		Department: constants.ComputerEngineering,
		Faculty:    constants.EngineeringFaculty,
		Year:       5,
		DOB:        minorsDOB,
	})
	assert.NoError(t, err)

	err = testHandler.room.CreateNewRoom(testHandler.ctx, &v1resources.Room{
		ID:             roomID,
		CreateByUserID: predator,
		RoomName:       "Just a random room to be predated on",
		RoomIcon:       "",
		CreatedAt:      time.Now(),
	})
	assert.NoError(t, err)

	testHandler.roomRequest.NewRoomRequest(testHandler.ctx, &v1resources.JoinRoomRequest{
		ID:           primitive.NewObjectID(),
		UserID:       predator,
		UserToJoinID: normalUser,
		RoomID:       roomID,
		RoomName:     "Just a random room to be predated on",
		Accepted:     true,
		CreatedAt:    time.Now(),
	})

	err = testHandler.messages.NewMessage(testHandler.ctx, &v1resources.Message{
		ID:        primitive.NewObjectID(),
		RoomID:    roomID,
		UserID:    predator,
		Message:   "How are you doing.",
		Type:      constants.NewMessageMsgType,
		CreatedAt: time.Now(),
	})
	assert.NoError(t, err)

	err = testHandler.messages.NewMessage(testHandler.ctx, &v1resources.Message{
		ID:        primitive.NewObjectID(),
		RoomID:    roomID,
		UserID:    normalUser,
		Message:   "Hi, I am fine thank you.",
		Type:      constants.NewMessageMsgType,
		CreatedAt: time.Now(),
	})
	assert.NoError(t, err)

	err = testHandler.messages.NewMessage(testHandler.ctx, &v1resources.Message{
		ID:        primitive.NewObjectID(),
		RoomID:    roomID,
		UserID:    predator,
		Message:   "I think I like you really.",
		Type:      constants.NewMessageMsgType,
		CreatedAt: time.Now(),
	})
	assert.NoError(t, err)

	err = testHandler.messages.NewMessage(testHandler.ctx, &v1resources.Message{
		ID:        primitive.NewObjectID(),
		RoomID:    roomID,
		UserID:    predator,
		Message:   "Can we meet any time soon, please, don't tell your parents as they might not approve.",
		Type:      constants.NewMessageMsgType,
		CreatedAt: time.Now(),
	})
	assert.NoError(t, err)
}

func createWSConn(email string) (*websocket.Conn, error) {
	user := v1resources.UserLogin{
		Email:    email,
		Password: "password",
	}

	resp, err := testHandler.makeNormalRequest("POST", fmt.Sprintf("%s/v1/login", testHandler.server.URL), user)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid status code")
	}

	body := apihelpers.Response{}
	err = testHandler.getResponseData(resp.Body, &body)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(body.Data)
	if err != nil {
		return nil, err
	}

	data := v1resources.LoginResponse{}

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("ws%s/v1/ws/%s/%s", strings.TrimPrefix(testHandler.server.URL, "http"), data.UserID.Hex(), data.WSTicket)
	websocket.DefaultDialer.TLSClientConfig = testHandler.server.TLS
	ws, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusSwitchingProtocols {
		return nil, errors.New("invalid status code")
	}
	return ws, nil
}
