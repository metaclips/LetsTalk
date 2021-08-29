package models

import (
	"context"

	"github.com/metaclips/LetsTalk/backend/middlewares/database"
	v1resources "github.com/metaclips/LetsTalk/backend/resources/api/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Messages struct {
	db database.Database
}

func NewMessages(db database.Database) Messages {
	return Messages{
		db,
	}
}

func (m Messages) NewMessage(ctx context.Context, message *v1resources.Message) error {
	return m.db.Save(ctx, message, nil)
}

func (m Messages) RetrieveMessagesForRoom(ctx context.Context, roomID primitive.ObjectID, limit, skip int) (*v1resources.MessagesList, error) {
	var messages []v1resources.Message

	err := m.db.FindMany(ctx, bson.M{"roomID": roomID}, nil, bson.M{"_id": -1}, int64(limit), int64(skip), &messages)
	if err != nil {
		return nil, err
	}

	count, err := m.db.Count(ctx, bson.M{"roomID": roomID})
	if err != nil {
		return nil, err
	}

	return &v1resources.MessagesList{
		Data:  messages,
		Count: count,
	}, nil
}
