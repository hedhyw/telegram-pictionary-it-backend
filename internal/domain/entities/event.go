package entities

import (
	"encoding/json"
	"fmt"
)

type SocketRequest struct {
	Name    string          `json:"name"`
	Payload json.RawMessage `json:"payload"`
}

type SocketResponse[T any] struct {
	Name    string `json:"name"`
	Payload T      `json:"payload"`
}

func NewSocketResponse[T fmt.Stringer](payload T) *SocketResponse[T] {
	return &SocketResponse[T]{
		Name:    payload.String(),
		Payload: payload,
	}
}
