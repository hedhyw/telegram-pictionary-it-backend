package game

import (
	"fmt"

	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/telegram"
)

// RequestEventCanvasChanged implements asyncmodel.RequestEvent.
// It indicates that the leader player changed their drawing.
type RequestEventCanvasChanged struct {
	ImageBase64 string
	ClientID    string
}

// String implements fmt.Stringer and asyncmodel.RequestEvent.
func (e RequestEventCanvasChanged) String() string { return fmt.Sprintf("%T", e) }

// IsRequestEvent implements asyncmodel.RequestEvent.
func (e *RequestEventCanvasChanged) IsRequestEvent() {}

// RequestEventGameStarted implements asyncmodel.RequestEvent.
// It indicates that a player started the game.
type RequestEventGameStarted struct{}

// String implements fmt.Stringer and asyncmodel.RequestEvent.
func (e RequestEventGameStarted) String() string { return fmt.Sprintf("%T", e) }

// IsRequestEvent implements asyncmodel.RequestEvent.
func (e *RequestEventGameStarted) IsRequestEvent() {}

// RequestEventPlayerJoined implements asyncmodel.RequestEvent.
// It indicates that the player joined the game.
type RequestEventPlayerJoined struct {
	ClientID string
	Meta     *telegram.InitDataMeta
}

// String implements fmt.Stringer and asyncmodel.RequestEvent.
func (e RequestEventPlayerJoined) String() string { return fmt.Sprintf("%T", e) }

// IsRequestEvent implements asyncmodel.RequestEvent.
func (e *RequestEventPlayerJoined) IsRequestEvent() {}

// RequestEventPlayerRemoved implements asyncmodel.RequestEvent.
// It indicates that the player left the game.
type RequestEventPlayerRemoved struct {
	ClientID string
}

// String implements fmt.Stringer and asyncmodel.RequestEvent.
func (e RequestEventPlayerRemoved) String() string { return fmt.Sprintf("%T", e) }

// IsRequestEvent implements asyncmodel.RequestEvent.
func (e *RequestEventPlayerRemoved) IsRequestEvent() {}

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
