package v1resources

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JoinRoomRequest struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	UserID       primitive.ObjectID `bson:"userID" json:"userID" validate:"required"`
	UserToJoinID primitive.ObjectID `bson:"userToJoinID" json:"userToJoinID" validate:"required"`
	RoomID       primitive.ObjectID `bson:"roomID" json:"roomID" validate:"required"`
	RoomName     string             `bson:"roomName" json:"roomName" validate:"required"`
	Accepted     bool               `json:"accepted" bson:"accepted"`
	RejectRoom   bool               `json:"rejectRoom" bson:"rejectRoom"`
	ExitedRoom   bool               `json:"exitedRoom" bson:"exitedRoom"`
	CreatedAt    time.Time          `bson:"createdAt" json:"-"`
	DeletedAt    *time.Time         `bson:"deletedAt" json:"-"`
}

type GetJoinedRoomRequest struct {
	RoomID   primitive.ObjectID `bson:"roomID" json:"roomID"`
	RoomName string             `bson:"roomName" json:"roomName" validate:"required"`
}

type JoinRoomRequestList struct {
	Data  []JoinRoomRequest `json:"data"`
	Count int               `json:"count"`
}

type ExitRoomRequest struct {
	JoinRoomID primitive.ObjectID `bson:"joinRoomID" json:"joinRoomID" validate:"required"`
	UserID     primitive.ObjectID `bson:"userID" json:"userID" validate:"required"`
	RoomID     primitive.ObjectID `bson:"roomID" json:"roomID" validate:"required"`
}

type GetJoinedRoomRequestList struct {
	Data  []GetJoinedRoomRequest `json:"data"`
	Count int                    `json:"count"`
}

// Room struct defineds content details for users room.
// Message count is used to track amount of messages sent by users in room,
// this helps with partitioning messages on retrieval, messages are retrieved in 20s on request.
type RoomOutput struct {
	RoomID          string                            `bson:"_id" json:"roomID,omitempty"`
	RoomName        string                            `bson:"roomName" json:"roomName,omitempty"`
	RoomIcon        string                            `bson:"roomIcon" json:"roomIcon"`
	RegisteredUsers map[primitive.ObjectID]UserStatus `bson:"-" json:"registeredUsers,omitempty"`
	MessageCount    int                               `bson:"messageCount" json:"messageCount,omitempty"`
}

// NewRoomRequest struct defines contents to create a new room.
type Room struct {
	ID             primitive.ObjectID `bson:"_id" json:"_id"`
	CreateByUserID primitive.ObjectID `json:"createByUserID" bson:"createByUserID"`
	RoomName       string             `bson:"roomName" json:"roomName,omitempty" validate:"required"`
	RoomIcon       string             `bson:"roomIcon" json:"roomIcon"`
	CreatedAt      time.Time          `bson:"createdAt" json:"-"`
	DeletedAt      *time.Time         `bson:"deletedAt" json:"-"`
}

type RoomList struct {
	Data  []Room `json:"data"`
	Count int    `json:"count"`
}
