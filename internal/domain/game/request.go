package game

import (
	"fmt"

	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/telegram"
)

type RequestEventCanvasChanged struct {
	ImageBase64 string
	ClientID    string
}

func (e RequestEventCanvasChanged) String() string { return fmt.Sprintf("%T", e) }

func (e *RequestEventCanvasChanged) IsRequestEvent() {}

type RequestEventGameStarted struct{}

func (e RequestEventGameStarted) String() string { return fmt.Sprintf("%T", e) }

func (e *RequestEventGameStarted) IsRequestEvent() {}

type RequestEventPlayerJoined struct {
	ClientID string
	Meta     *telegram.InitDataMeta
}

func (e RequestEventPlayerJoined) String() string { return fmt.Sprintf("%T", e) }

func (e *RequestEventPlayerJoined) IsRequestEvent() {}

type RequestEventPlayerRemoved struct {
	ClientID string
}

func (e RequestEventPlayerRemoved) String() string { return fmt.Sprintf("%T", e) }

func (e *RequestEventPlayerRemoved) IsRequestEvent() {}

type RequestEventWordGuessAttempted struct {
	ClientID string `json:"clientId"`
	Word     string `json:"word"`
}

func (e RequestEventWordGuessAttempted) String() string { return fmt.Sprintf("%T", e) }

func (e *RequestEventWordGuessAttempted) IsRequestEvent() {}
