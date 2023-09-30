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

func NewSocketResponse(payload fmt.Stringer) *SocketResponse {
	return &SocketResponse{
		Name:    payload.String(),
		Payload: payload,
	}
}
