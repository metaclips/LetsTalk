package messages

// Message struct defines user message contents, size and hash is defined if user is sending files.
// Index is used to track message count as to Rooms messages, this should help with partitioning if we
// are to retrieve message of a particular count.
type Message struct {
	RoomID      string `bson:"roomID" json:"roomID,omitempty"`
	Message     string `bson:"message" json:"message,omitempty"`
	UserID      string `bson:"userID" json:"userID,omitempty"`
	Name        string `bson:"name" json:"name,omitempty"`
	Time        string `bson:"time" json:"time,omitempty"`
	Type        string `bson:"type" json:"type,omitempty"`
	Size        string `bson:"size,omitempty" json:"size,omitempty"`
	Hash        string `bson:"hash,omitempty" json:"hash,omitempty"`
	Index       int    `bson:"index" json:"index,omitempty"`
	MessageType string `bson:"-" json:"msgType,omitempty"`
}
