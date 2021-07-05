package v1controller

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/metaclips/LetsTalk/backend/helpers/constants"
	"github.com/metaclips/LetsTalk/backend/helpers/maps"
	"github.com/metaclips/LetsTalk/backend/helpers/paging"
	"github.com/metaclips/LetsTalk/backend/helpers/validation"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h HomepageController) searchUsers(key string, user primitive.ObjectID) {
	if key == "" {
		return
	}

	pager := paging.Paging{
		Limit:  0,
		Skip:   0,
		Search: key,
	}

	usersList, err := h.user.GetUsersPrunedData(h.ctx, &pager)
	if err != nil {
		log.Errorln("error retrieving users data on `searchUsers`: ", err)
		return
	}

	if err := h.sendMessage(usersList, user); err != nil {
		log.Errorln("error sending users data on `searchUsers`: ", err)
		return
	}
}

func (h HomepageController) newMessage(userID primitive.ObjectID, msg []byte) {
	var message v1resources.Message

	if err := unmarshalTo(msg, &message); err != nil {
		// ToDo: Return error message.
		return
	}

	if err := validation.ValidateInputs(message); err != nil {
		// ToDo: Return error message.
		return
	}

	if message.UserID != userID {
		log.Warning("user tried using a fake userID")
		// ToDo: Return error message.
		return
	}

	if _, ok := maps.FileTypes[message.Type]; !ok {
		// ToDo: Return error message.
		return
	}

	users, err := h.roomRequest.GetUsersInRoom(h.ctx, message.RoomID)
	if err != nil {
		log.Errorln("error geting users in room:", err)
		// ToDo: Return error message.
		return
	}

	if message.Type == constants.FileTypeFile && (message.Size == "" || message.Hash == "") {
		// ToDo: Return error message.
		return
	}

	message.Time = time.Now()
	message.ID = primitive.NewObjectID()

	if err := h.messages.NewMessage(h.ctx, &message); err != nil {
		log.Errorln("error saving message to DB:", err)
		// ToDo: Return error message.
		return
	}

	for _, user := range users {
		// ToDo: Handle this error properly.
		if err := h.sendMessage(message, user.UserID); err != nil {
			log.Errorln("error sending message to user:", err)
			// ToDo: Return error message.
			continue
		}
	}
}

func (h HomepageController) joinRoomRequest(userID primitive.ObjectID, msg []byte) {
	var request = v1resources.JoinRoomRequest{}
	if err := unmarshalTo(msg, &request); err != nil {
		// ToDo: Return error message.
		return
	}

	if err := validation.ValidateInputs(request); err != nil {
		// ToDo: Return error message.
		return
	}

	isUserInRoom, err := h.roomRequest.IsUserInRoom(h.ctx, request.RoomID, request.UserID)
	if err != nil {
		// ToDo: Return error message.
		return
	}

	if _, err := h.user.GetUserByID(h.ctx, request.UserToJoinID); err != nil {
		if err != mongo.ErrNoDocuments {
			log.Errorln("user tried inviting an unregistered user")
		}

		return
	}

	if !isUserInRoom {
		// ToDo: Return error message.
		return
	}

	request.Accepted = false
	request.ExitedRoom = false
	request.ID = primitive.NewObjectID()

	if err := h.roomRequest.NewRoomRequest(h.ctx, &request); err != nil {
		// ToDo: Return error message.
		return
	}

	requester, err := h.user.GetUserByID(h.ctx, request.UserID)
	if err != nil {
		// ToDo: Return error message.
		return
	}

	userToJoin, err := h.user.GetUserByID(h.ctx, request.UserToJoinID)
	if err != nil {
		// ToDo: Return error message.
		return
	}

	h.saveMessageAndNotifyRoom(request.RoomID, request.UserID, fmt.Sprintf("%s %s (%s) has requested %s %s (%s) to join the room.", requester.FirstName, requester.LastName, requester.Email, userToJoin.FirstName, userToJoin.LastName, userToJoin.Email))
}

func (h HomepageController) acceptRoomRequest(userID primitive.ObjectID, msg []byte) {
	var request = v1resources.JoinRoomRequest{}
	if err := unmarshalTo(msg, &request); err != nil {
		// ToDo: Return error message.
		return
	}

	if err := validation.ValidateInputs(request); err != nil {
		// ToDo: Return error message.
		return
	}

	userIsAllowedToRoom, err := h.roomRequest.IsUserInRoom(h.ctx, request.RoomID, request.UserID)
	if err != nil {
		// ToDo: Return error message.
		return
	}
	if !userIsAllowedToRoom {
		// ToDo: Return error message.
		return
	}

	if request.UserToJoinID != userID {
		// ToDo: Return error message.
		return
	}

	if request.ID.IsZero() {
		// ToDo: Return error message.
		return
	}

	request.Accepted = true
	request.ExitedRoom = false

	if err := h.roomRequest.AcceptRequest(h.ctx, request.ID); err != nil {
		log.Errorln("error accepting room request")
		// ToDo: Return error message.
		return
	}

	userToJoin, err := h.user.GetUserByID(h.ctx, userID)
	if err != nil {
		// ToDo: Return error message.
		return
	}

	h.saveMessageAndNotifyRoom(request.RoomID, request.UserID, fmt.Sprintf("%s %s (%s) has has joined the room.", userToJoin.FirstName, userToJoin.LastName, userToJoin.Email))
}

func (h HomepageController) createNewRoom(userID primitive.ObjectID, msg []byte) {
	var room v1resources.Room
	if err := unmarshalTo(msg, &room); err != nil {
		// ToDo: Return error message.
		return
	}

	if err := validation.ValidateInputs(room); err != nil {
		// ToDo: Return error message.
		return
	}

	if room.CreateByUserID != userID {
		// ToDo: Return error message.
		return
	}

	room.ID = primitive.NewObjectID()
	if err := h.room.CreateNewRoom(h.ctx, &room); err != nil {
		// ToDo: Return error message.
		log.Errorln("error creating new room:", err)
		return
	}

	err := h.roomRequest.NewRoomRequest(h.ctx, &v1resources.JoinRoomRequest{
		ID:           primitive.NewObjectID(),
		UserID:       userID,
		UserToJoinID: userID,
		RoomID:       room.ID,
		RoomName:     room.RoomName,
		Accepted:     true,
	})

	if err != nil {
		log.Errorln("error creating new room request:", err)
		return
	}

	userToJoin, err := h.user.GetUserByID(h.ctx, userID)
	if err != nil {
		// ToDo: Return error message.
		return
	}

	h.saveMessageAndNotifyRoom(room.ID, userID, fmt.Sprintf("%s %s (%s) has has joined the room.", userToJoin.FirstName, userToJoin.LastName, userToJoin.Email))
}

func (h HomepageController) saveMessageAndNotifyRoom(roomID, userID primitive.ObjectID, message string) {
	data := v1resources.Message{
		ID:      primitive.NewObjectID(),
		RoomID:  roomID,
		UserID:  userID,
		Message: message,
		Type:    constants.NewMessageMsgType,
		Time:    time.Now(),
	}

	err := h.messages.NewMessage(h.ctx, &data)
	if err != nil {
		// ToDo: Return error message.
		return
	}

	// Broadcast message to all room users that a `join room request` has been sent.
	if err := h.broadcastToAllRoomUsers(roomID, data); err != nil {
		// ToDo: Return error message.
		return
	}
}

func (h HomepageController) broadcastToAllRoomUsers(roomID primitive.ObjectID, data interface{}) error {
	joinedUsers, err := h.roomRequest.GetUsersInRoom(h.ctx, roomID)
	if err != nil {
		return err
	}

	for _, user := range joinedUsers {
		if err := h.sendMessage(data, user.UserID); err != nil {
			return err
		}
	}

	return nil
}

func unmarshalTo(data []byte, to interface{}) error {
	return json.Unmarshal(data, &to)
}
