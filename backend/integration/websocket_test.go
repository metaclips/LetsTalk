package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	apihelpers "github.com/metaclips/LetsTalk/backend/apihelper"
	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	"github.com/stretchr/testify/assert"
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
