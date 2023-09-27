package entities

import (
	"encoding/json"
	"fmt"
)

type SocketRequest struct {
	Name    string          `json:"name"`
	Payload json.RawMessage `json:"payload"`
}

type SocketResponse struct {
	Name    string `json:"name"`
	Payload any    `json:"payload"`
}

func NewSocketEventPayload(event fmt.Stringer) *SocketResponse {
	return &SocketResponse{
		Name:    event.String(),
		Payload: event,
	}
}
