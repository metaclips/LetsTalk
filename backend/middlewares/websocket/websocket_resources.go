package websocket

import (
	"sync"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WSInterface interface {
	WritePump(conn *Connection)
	ReadPump(userID primitive.ObjectID, session Subscription)
	RegisterToHub(sub Subscription)
}

type hub struct {
	// Registered connections.
	users wsUsers
	// Inbound messages from the connections.
	broadcast chan wsMessage
	// Register requests from the connections.
	register chan Subscription
	// Unregister requests from connections.
	unRegister chan Subscription
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

type wsUsers struct {
	// Registered users on a WS connections.
	roomUsers map[primitive.ObjectID]*Connection
	// Lock.
	mutex *sync.RWMutex
}

// Websocket message structure.
type wsMessage struct {
	// Websocket message.
	data []byte
	// User who sent message.
	user primitive.ObjectID
}
