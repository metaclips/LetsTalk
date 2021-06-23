package v1resources

import "go.mongodb.org/mongo-driver/bson/primitive"

// Message struct defines user message contents, size and hash is defined if user is sending files.
// Index is used to track message count as to Rooms messages, this should help with partitioning if we
// are to retrieve message of a particular count.
type Message struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	RoomID  primitive.ObjectID `bson:"roomID" json:"roomID,omitempty"`
	UserID  primitive.ObjectID `bson:"userID" json:"userID,omitempty"`
	Message string             `bson:"message" json:"message,omitempty"`
	Time    string             `bson:"time" json:"time,omitempty"`
	Type    string             `bson:"type" json:"type,omitempty"`
	Size    string             `bson:"size,omitempty" json:"size,omitempty"`
	Hash    string             `bson:"hash,omitempty" json:"hash,omitempty"`
	IsFile  bool               `json:"isFile" bson:"isFile"`
	Index   int                `bson:"index" json:"index,omitempty"`
}

type MessagesList struct {
	Data  []Message `json:"data"`
	Count int64     `json:"count"`
}
