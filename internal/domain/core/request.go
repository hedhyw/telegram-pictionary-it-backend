package core

import (
	"encoding/json"
	"fmt"
)

// RequestEventWordGuessAttempted implements asyncmodel.RequestEvent.
// It indicates an attempt to guess the word by a user.
type RequestEventWordGuessAttempted struct {
	ClientID string `json:"clientId"`
	Word     string `json:"word"`
}

// String implements fmt.Stringer and asyncmodel.RequestEvent.
func (e RequestEventWordGuessAttempted) String() string { return fmt.Sprintf("%T", e) }

// IsRequestEvent implements asyncmodel.RequestEvent.
func (e *RequestEventWordGuessAttempted) IsRequestEvent() {}

// RequestEventCanvasChanged implements asyncmodel.RequestEvent.
// It holds the current drawing by the leader player.
type RequestEventCanvasChanged struct {
	ImageBase64 string `json:"imageBase64"`

	ClientID string `json:"clientId"`
}

// String implements fmt.Stringer and asyncmodel.RequestEvent.
func (e RequestEventCanvasChanged) String() string { return fmt.Sprintf("%T", e) }

// IsRequestEvent implements asyncmodel.RequestEvent.
func (e *RequestEventCanvasChanged) IsRequestEvent() {}

// RequestEventCanvasChanged implements asyncmodel.RequestEvent.
// It indicates that the client started a game.
type RequestEventGameStarted struct {
	ClientID string `json:"clientId"`
}

// String implements fmt.Stringer and asyncmodel.RequestEvent.
func (e RequestEventGameStarted) String() string { return fmt.Sprintf("%T", e) }

// IsRequestEvent implements asyncmodel.RequestEvent.
func (e *RequestEventGameStarted) IsRequestEvent() {}

// RequestEventClientConnected implements asyncmodel.RequestEvent.
// It indicates that a new websocket client is connected.
type RequestEventClientConnected struct {
	ClientID string                 `json:"clientId"`
	EventsCh chan<- json.RawMessage `json:"-"`
}

// String implements fmt.Stringer and asyncmodel.RequestEvent.
func (e RequestEventClientConnected) String() string { return fmt.Sprintf("%T", e) }

// IsRequestEvent implements asyncmodel.RequestEvent.
func (e *RequestEventClientConnected) IsRequestEvent() {}

// RequestEventPlayerInitialized implements asyncmodel.RequestEvent.
// It holds InitData for the player from Telegram.
type RequestEventPlayerInitialized struct {
	ClientID    string `json:"clientId"`
	InitDataRaw string `json:"initDataRaw"`
}

// String implements fmt.Stringer and asyncmodel.RequestEvent.
func (e RequestEventPlayerInitialized) String() string { return fmt.Sprintf("%T", e) }

// IsRequestEvent implements asyncmodel.RequestEvent.
func (e *RequestEventPlayerInitialized) IsRequestEvent() {}

// RequestEventClientDisconnnected implements asyncmodel.RequestEvent.
// It indicates, that the player left the room.
type RequestEventClientDisconnnected struct {
	ClientID string
}

// String implements fmt.Stringer and asyncmodel.RequestEvent.
func (e RequestEventClientDisconnnected) String() string { return fmt.Sprintf("%T", e) }

// IsRequestEvent implements asyncmodel.RequestEvent.
func (e *RequestEventClientDisconnnected) IsRequestEvent() {}
