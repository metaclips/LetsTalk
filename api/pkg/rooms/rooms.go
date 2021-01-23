package rooms

import (
	"context"
	"errors"

	"github.com/metaclips/LetsTalk/api/pkg/database"
	"github.com/metaclips/LetsTalk/api/values"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	ctx context.Context
	db  database.Database
}

func NewRoomHandler(ctx context.Context, db database.Database) *RoomHandler {
	return &RoomHandler{
		ctx: ctx,
		db:  db,
	}
}

func (r *RoomHandler) FindRoom(roomID string) (Room, error) {
	room := Room{}

	err := r.db.FindOne(r.ctx, bson.M{"_id": roomID}, nil, &room)

	return room, err
}

func (r *RoomHandler) CreateNewRoom(roomRequest *NewRoomRequest) (string, error) {
	var chats Room

	chats.RoomID = primitive.NewObjectID().Hex()
	chats.RoomName = roomRequest.RoomName
	chats.RegisteredUsers = append(chats.RegisteredUsers, roomRequest.Email)

	if err := r.db.Save(r.ctx, chats, nil); err != nil {
		return "", err
	}

	return chats.RoomID, nil
}

func (r *RoomHandler) IsJoinRoomRequesterLegit(userToJoinEmail string, req *JoinRequest) error {
	room, err := r.FindRoom(req.RoomID)
	if err != nil {
		return err
	}

	var requesterLegit bool
	for _, registeredUser := range room.RegisteredUsers {
		if registeredUser == req.RequestingUserID {
			requesterLegit = true
		} else if registeredUser == userToJoinEmail {
			return values.ErrUserExistInRoom
		}
	}

	if !requesterLegit {
		return errors.New(("invalid user made a join request to room"))
	}

	return nil
}
