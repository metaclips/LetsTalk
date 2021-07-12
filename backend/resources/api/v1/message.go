package v1resources

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message struct defines user message contents, size and hash is defined if user is sending files.
// Index is used to track message count as to Rooms messages, this should help with partitioning if we
// are to retrieve message of a particular count.
type Message struct {
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	RoomID    primitive.ObjectID `bson:"roomID" json:"roomID,omitempty" validate:"required"`
	UserID    primitive.ObjectID `bson:"userID" json:"userID,omitempty" validate:"required"`
	Message   string             `bson:"message" json:"message,omitempty" validate:"required"`
	Type      string             `bson:"type" json:"type,omitempty" validate:"required"`
	Size      string             `bson:"size,omitempty" json:"size,omitempty"`
	Hash      string             `bson:"hash,omitempty" json:"hash,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	DeletedAt *time.Time         `bson:"deletedAt" json:"-"`
}

type MessagesList struct {
	Data  []Message `json:"data"`
	Count int       `json:"count"`
}

type RequestMessage struct {
	RoomID primitive.ObjectID `bson:"roomID" json:"roomID,omitempty" validate:"required"`
	Skip   int                `json:"skip"`
	Limit  int                `json:"limit"`
}
