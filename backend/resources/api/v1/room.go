package v1resources

import "go.mongodb.org/mongo-driver/bson/primitive"

type JoinRoomRequest struct {
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	UserID       primitive.ObjectID `bson:"userID" json:"userID"`
	UserToJoinID primitive.ObjectID `bson:"userToJoinID" json:"userToJoinID"`
	RoomID       string             `bson:"roomID" json:"roomID"`
	RoomName     string             `bson:"roomName" json:"roomName"`
	Accepted     bool               `json:"accepted" bson:"accepted"`
	ExitedRoom   bool               `json:"exitedRoom" bson:"exitedRoom"`
}

type GetJoinedRoomRequest struct {
	ID primitive.ObjectID `bson:"_id" json:"_id"`
}

type JoinRoomRequestList struct {
	Data  []JoinRoomRequest `json:"data"`
	Count int64             `json:"count"`
}

type GetJoinedRoomRequestList struct {
	Data  []GetJoinedRoomRequest `json:"data"`
	Count int64                  `json:"count"`
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
	ID           primitive.ObjectID `bson:"_id" json:"_id"`
	RoomName     string             `bson:"roomName" json:"roomName,omitempty"`
	RoomIcon     string             `bson:"roomIcon" json:"roomIcon"`
	MessageCount int                `bson:"messageCount" json:"messageCount"`
}

type RoomList struct {
	Data  []Room `json:"data"`
	Count int64  `json:"count"`
}
