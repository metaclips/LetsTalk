package rooms

import "github.com/metaclips/LetsTalk/api/pkg/messages"

type RoomsJoined struct {
	RoomID   string `bson:"roomID" json:"roomID"`
	RoomName string `bson:"roomName" json:"roomName"`
}

type JoinRequest struct {
	RoomID             string   `bson:"_id" json:"roomID"`
	RoomName           string   `bson:"roomName" json:"roomName"`
	RequestingUserName string   `bson:"requestingUserName" json:"requestingUserName"`
	RequestingUserID   string   `bson:"requestingUserID" json:"requestingUserID"`
	Users              []string `bson:"-" json:"users"` // Users whom join requests are to be sent to, if client wants to send request to multiple users.
}

// Room struct defineds content details for users room.
// Message count is used to track amount of messages sent by users in room,
// this helps with partitioning messages on retrieval, messages are retrieved in 20s on request.
// FirstLoad is specified if user initially wants to retrieve room messages from frontend.
type Room struct {
	RoomID          string             `bson:"_id" json:"roomID,omitempty"`
	RoomName        string             `bson:"roomName" json:"roomName,omitempty"`
	RoomIcon        string             `bson:"roomIcon" json:"roomIcon"`
	RegisteredUsers []string           `bson:"registeredUsers" json:"registeredUsers,omitempty"`
	MessageCount    int                `bson:"messageCount" json:"messageCount,omitempty"`
	Messages        []messages.Message `bson:"-" json:"messages,omitempty"`
	FirstLoad       bool               `bson:"-" json:"firstLoad,omitempty"`
}

type NewRoomRequest struct {
	Email       string `json:"userID"`
	RoomName    string `json:"roomName"`
	MessageType string `bson:"-" json:"msgType"`
}
