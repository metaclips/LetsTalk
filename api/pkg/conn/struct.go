package conn

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Websocket users.
type wsUsers struct {
	users map[string]map[*connection]bool // Users connected to websocket
	mutex *sync.RWMutex
}

// Hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	users      wsUsers           // Registered connections.
	broadcast  chan wsMessage    // Inbound messages from the connections.
	register   chan subscription // Register requests from the connections.
	unRegister chan subscription // Unregister requests from connections.
}

type subscription struct {
	conn *connection
	user string
}

// Connection is an middleman between the websocket connection and the hub.
type connection struct {
	ws   *websocket.Conn
	send chan []byte
}

// Websocket containing message and user who sent.
type wsMessage struct {
	data []byte
	user string
}
