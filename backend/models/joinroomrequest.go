package models

import (
	"context"

	"github.com/metaclips/LetsTalk/backend/middlewares/database"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	if err := b.db.Save(ctx, room, nil); err != nil {
		return err
	}

	return nil
}

func (b RoomRequest) AcceptRequest(ctx context.Context, ID primitive.ObjectID) error {
	return b.db.UpdateByID(ctx, ID, bson.M{"accepted": true, "exitedRoom": false, "rejectRoom": false})
}

func (b RoomRequest) ExitRoom(ctx context.Context, ID primitive.ObjectID) error {
	return b.db.UpdateByID(ctx, ID, bson.M{"accepted": false, "exitedRoom": true, "rejectRoom": false})
}

func (b RoomRequest) RejectRoom(ctx context.Context, ID primitive.ObjectID) error {
	return b.db.UpdateByID(ctx, ID, bson.M{"accepted": false, "exitedRoom": false, "rejectRoom": true})
}

func (b RoomRequest) GetUsersInRoom(ctx context.Context, roomID primitive.ObjectID) ([]v1resources.JoinRoomRequest, error) {
	var rooms []v1resources.JoinRoomRequest

	if err := b.db.FindMany(ctx, bson.M{"roomID": roomID, "accepted": true, "exitedRoom": false, "rejectRoom": false}, nil, nil, 0, 0, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (b RoomRequest) GetUserAssociates(ctx context.Context, userID primitive.ObjectID) ([]primitive.ObjectID, error) {
	var rooms []v1resources.JoinRoomRequest

	if err := b.db.FindMany(ctx, bson.M{"userID": userID, "accepted": true, "exitedRoom": false, "rejectRoom": false}, nil, nil, 0, 0, &rooms); err != nil {
		return nil, err
	}

	associates := make(map[primitive.ObjectID]bool)
	roomsBSON := make([]bson.M, len(rooms))

	for i, room := range rooms {
		roomsBSON[i] = bson.M{"roomID": room.RoomID, "accepted": true, "exitedRoom": false, "rejectRoom": false}
	}

	var users []v1resources.JoinRoomRequest
	if err := b.db.FindMany(ctx, bson.M{"$or": roomsBSON}, nil, nil, 0, 0, &users); err != nil {
		return nil, err
	}

	for _, user := range users {
		associates[user.UserToJoinID] = true
	}

	userIDs := make([]primitive.ObjectID, 0)
	for associate := range associates {
		userIDs = append(userIDs, associate)
	}

	return userIDs, nil
}

func (b RoomRequest) IsRequesterAllowedInRoom(ctx context.Context, roomID, userID primitive.ObjectID) (bool, error) {
	var room *v1resources.JoinRoomRequest

	if err := b.db.FindOne(ctx, bson.M{"roomID": roomID, "userID": userID}, nil, nil, &room); err != nil && err != mongo.ErrNoDocuments {
		return false, err
	}

	return (room != nil && room.Accepted), nil
}

func (b RoomRequest) IsUserAllowedInRoom(ctx context.Context, roomID, userToJoinID primitive.ObjectID) (bool, error) {
	var room *v1resources.JoinRoomRequest

	if err := b.db.FindOne(ctx, bson.M{"roomID": roomID, "userToJoinID": userToJoinID}, nil, nil, &room); err != nil && err != mongo.ErrNoDocuments {
		return false, err
	}

	return room != nil, nil
}

func (b RoomRequest) HasUserAcceptedJoinRequest(ctx context.Context, userID, roomID primitive.ObjectID) (bool, error) {
	var room *v1resources.JoinRoomRequest

	if err := b.db.FindOne(ctx, bson.M{"roomID": roomID, "userToJoinID": userID}, nil, nil, &room); err != nil && err != mongo.ErrNoDocuments {
		return false, err
	}

	return (room != nil && room.Accepted), nil
}

func (b RoomRequest) RetrieveUnacceptedRooms(ctx context.Context, userID primitive.ObjectID, limit, skip int) (*v1resources.JoinRoomRequestList, error) {
	var rooms []v1resources.JoinRoomRequest

	if err := b.db.FindMany(ctx, bson.M{"userToJoinID": userID, "accepted": false, "exitedRoom": false, "rejectRoom": false}, nil, nil, int64(limit), int64(skip), &rooms); err != nil {
		return nil, err
	}

	count, err := b.db.Count(ctx, bson.M{"userToJoinID": userID, "accepted": false, "exitedRoom": false, "rejectRoom": false})
	if err != nil {
		return nil, err
	}

	return &v1resources.JoinRoomRequestList{
		Data:  rooms,
		Count: count,
	}, nil
}

func (b RoomRequest) RetrieveAcceptedJoinRequest(ctx context.Context, userID primitive.ObjectID, limit, skip int) (*v1resources.GetJoinedRoomRequestList, error) {
	var rooms []v1resources.GetJoinedRoomRequest

	err := b.db.FindMany(ctx, bson.M{"userToJoinID": userID, "accepted": true, "exitedRoom": false, "rejectRoom": false}, nil, bson.M{"_id": -1}, int64(limit), int64(skip), &rooms)
	if err != nil {
		return nil, err
	}

	count, err := b.db.Count(ctx, bson.M{"userToJoinID": userID, "accepted": true, "exitedRoom": false, "rejectRoom": false})
	if err != nil {
		return nil, err
	}

	return &v1resources.GetJoinedRoomRequestList{
		Data:  rooms,
		Count: count,
	}, nil
}

func (b RoomRequest) RetrieveRejectedRoomRequest(ctx context.Context, userID primitive.ObjectID, limit, skip int) (*v1resources.GetJoinedRoomRequestList, error) {
	var rooms []v1resources.GetJoinedRoomRequest

	err := b.db.FindMany(ctx, bson.M{"userToJoinID": userID, "accepted": false, "rejectRoom": true, "exitedRoom": false}, nil, bson.M{"_id": -1}, int64(limit), int64(skip), &rooms)
	if err != nil {
		return nil, err
	}

	count, err := b.db.Count(ctx, bson.M{"userToJoinID": userID, "accepted": false, "exitedRoom": false, "rejectRoom": false})
	if err != nil {
		return nil, err
	}

	return &v1resources.GetJoinedRoomRequestList{
		Data:  rooms,
		Count: count,
	}, nil
}

func (b RoomRequest) GetUsersRoom(ctx context.Context, users []v1resources.PrunedUser) ([]primitive.ObjectID, error) {
	var rooms []v1resources.GetJoinedRoomRequest

	filters := []bson.M{}
	for _, user := range users {
		filters = append(filters, bson.M{
			"userToJoinID": user.ID,
			"accepted":     true,
			"exitedRoom":   false,
			"rejectRoom":   false,
		})
	}

	if err := b.db.FindMany(ctx, bson.M{"$or": filters}, nil, bson.M{"_id": -1}, 0, 0, &rooms); err != nil {
		return nil, err
	}

	filtered := map[primitive.ObjectID]bool{}

	for _, room := range rooms {
		filtered[room.RoomID] = true
	}

	var filteredRoom []primitive.ObjectID
	for key := range filtered {
		filteredRoom = append(filteredRoom, key)
	}

	return filteredRoom, nil
}
