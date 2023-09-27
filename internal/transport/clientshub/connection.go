package clientshub

import "encoding/json"

type Connection struct {
	clientID string
	events   chan<- json.RawMessage
}

func (conn Connection) ClientID() string {
	return conn.clientID
}

func (h *Hub) newConnection(
	clientID string,
	eventsCh chan<- json.RawMessage,
) *Connection {
	return &Connection{
		clientID: clientID,
		events:   eventsCh,
	}
}
