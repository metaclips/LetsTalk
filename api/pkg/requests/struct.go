package requests

type newRoomRequest struct {
	Email       string `json:"userID"`
	RoomName    string `json:"roomName"`
	MessageType string `bson:"-" json:"msgType"`
}

type joined struct {
	RoomID      string `json:"roomID"`
	RoomName    string `json:"roomName"`
	Email       string `json:"userID"`
	Name        string `json:"name"`
	RequesterID string `json:"requesterID"`
	Joined      bool   `json:"joined"`
	MessageType string `bson:"-" json:"msgType"`
}

type joinRequest struct {
	RoomID             string   `bson:"_id" json:"roomID"`
	RoomName           string   `bson:"roomName" json:"roomName"`
	RequestingUserName string   `bson:"requestingUserName" json:"requestingUserName"`
	RequestingUserID   string   `bson:"requestingUserID" json:"requestingUserID"`
	Users              []string `bson:"-" json:"users"` // Users whom join requests are to be sent to, if client wants to send request to multiple users.
}
