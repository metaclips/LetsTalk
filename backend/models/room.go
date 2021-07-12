package models

import (
	"context"

	"github.com/metaclips/LetsTalk/backend/middlewares/database"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	db database.Database
}

func NewRoom(db database.Database) Room {
	return Room{
		db,
	}
}

// Create new Room instance. Note should be called right after storing
// success message.
func (b Room) CreateNewRoom(ctx context.Context, room *v1resources.Room) error {
	if err := b.db.Save(ctx, room, nil); err != nil {
		return err
	}

	return nil
}

func (b Room) IncRoomMessageCount(ctx context.Context, roomID primitive.ObjectID) error {
	return b.db.IncCount(ctx, bson.M{"_id": roomID}, bson.M{"messageCount": 1})
}

func (b Room) UpdateRoom(ctx context.Context, room *v1resources.Room) error {
	return b.db.UpdateByID(ctx, room.ID, room)
}

func (b Room) RetrieveJoinedRooms(ctx context.Context, IDs []v1resources.GetJoinedRoomRequest, limit, skip int) (*v1resources.RoomList, error) {
	rooms := make([]v1resources.Room, len(IDs))

	for i, id := range IDs {
		err := b.db.FindByID(ctx, id.RoomID, nil, &rooms[i])
		if err != nil {
			return nil, err
		}
	}

	return &v1resources.RoomList{
		Data:  rooms,
		Count: len(IDs),
	}, nil
}
