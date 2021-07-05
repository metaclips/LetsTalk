package models

import (
	"context"

	"github.com/metaclips/LetsTalk/backend/middlewares/database"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomRequest struct {
	db database.Database
}

func NewRoomRequest(db database.Database) RoomRequest {
	return RoomRequest{
		db,
	}
}

// Create new room instance. Note should be called right after storing
// success message.
func (b RoomRequest) NewRoomRequest(ctx context.Context, room *v1resources.JoinRoomRequest) error {
	room.ID = primitive.NewObjectID()

	if err := b.db.Save(ctx, room, nil); err != nil {
		return err
	}

	return nil
}

func (b RoomRequest) AcceptRequest(ctx context.Context, ID primitive.ObjectID) error {
	return b.db.UpdateByID(ctx, ID, bson.M{"accepted": true, "exitedRoom": false})
}

func (b RoomRequest) ExitRoom(ctx context.Context, ID primitive.ObjectID) error {
	return b.db.UpdateByID(ctx, ID, bson.M{"exitedRoom": true})
}

func (b RoomRequest) GetUsersInRoom(ctx context.Context, roomID primitive.ObjectID) ([]v1resources.JoinRoomRequest, error) {
	var rooms []v1resources.JoinRoomRequest

	if err := b.db.FindMany(ctx, bson.M{"roomID": roomID}, nil, nil, 0, 0, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (b RoomRequest) IsUserInRoom(ctx context.Context, roomID, userID primitive.ObjectID) (bool, error) {
	var room *v1resources.JoinRoomRequest

	if err := b.db.FindOne(ctx, bson.M{"roomID": roomID, "userID": userID}, nil, room); err != nil {
		return false, err
	}

	return room != nil, nil
}

func (b RoomRequest) RetrieveUnacceptedRooms(ctx context.Context, userID primitive.ObjectID, limit, skip int64) (*v1resources.JoinRoomRequestList, error) {
	var rooms []v1resources.JoinRoomRequest

	if err := b.db.FindMany(ctx, bson.M{"userToJoinID": userID, "accepted": false}, nil, nil, limit, skip, &rooms); err != nil {
		return nil, err
	}

	count, err := b.db.Count(ctx, bson.M{"userToJoinID": userID, "accepted": false})
	if err != nil {
		return nil, err
	}

	return &v1resources.JoinRoomRequestList{
		Data:  rooms,
		Count: count,
	}, nil
}

func (b RoomRequest) RetrieveAcceptedJoinRequest(ctx context.Context, userID primitive.ObjectID, limit, skip int64) (*v1resources.GetJoinedRoomRequestList, error) {
	var rooms []v1resources.GetJoinedRoomRequest

	err := b.db.FindMany(ctx, bson.M{"userToJoinID": userID, "accepted": true}, nil, bson.M{"_id": -1}, limit, skip, &rooms)
	if err != nil {
		return nil, err
	}

	count, err := b.db.Count(ctx, bson.M{"userToJoinID": userID, "accepted": true})
	if err != nil {
		return nil, err
	}

	return &v1resources.GetJoinedRoomRequestList{
		Data:  rooms,
		Count: count,
	}, nil
}
