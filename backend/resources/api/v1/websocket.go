package v1resources

import (
	"sync"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Hub struct {
	// Registered connections.
	Users WSUsers
	// Inbound messages from the connections.
	Broadcast chan WSMessage
	// Register requests from the connections.
	Register chan Subscription
	// Unregister requests from connections.
	UnRegister chan Subscription
}

type Subscription struct {
	Conn   *Connection
	UserID primitive.ObjectID
}

// Connection is an middleman between the websocket connection and the hub.
type Connection struct {
	// The websocket connection.
	WS *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte
}

type WSUsers struct {
	// Registered users on a WS connections.
	RoomUsers map[primitive.ObjectID]*Connection
	// Lock.
	Mutex *sync.RWMutex
}

// Websocket message structure.
type WSMessage struct {
	// Websocket message.
	Data []byte
	// User who sent message.
	User primitive.ObjectID
}
