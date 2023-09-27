package core

import (
	"encoding/json"
	"fmt"
)

type RequestEventWordGuessAttempted struct {
	ClientID string `json:"clientId"`
	Word     string `json:"word"`
}

func (e RequestEventWordGuessAttempted) String() string { return fmt.Sprintf("%T", e) }

func (e *RequestEventWordGuessAttempted) IsRequestEvent() {}

type RequestEventCanvasChanged struct {
	ImageBase64 string `json:"imageBase64"`

	ClientID string `json:"clientId"`
}

func (e RequestEventCanvasChanged) String() string { return fmt.Sprintf("%T", e) }

func (e *RequestEventCanvasChanged) IsRequestEvent() {}

type RequestEventGameStarted struct {
	ClientID string `json:"clientId"`
}

func (e RequestEventGameStarted) String() string { return fmt.Sprintf("%T", e) }

func (e *RequestEventGameStarted) IsRequestEvent() {}

type RequestEventClientConnected struct {
	ChatID   string                 `json:"chatId"`
	ClientID string                 `json:"clientId"`
	EventsCh chan<- json.RawMessage `json:"-"`
}

func (e RequestEventClientConnected) String() string { return fmt.Sprintf("%T", e) }

func (e *RequestEventClientConnected) IsRequestEvent() {}

type RequestEventClientDisconnnected struct {
	ClientID string
}

func (e RequestEventClientDisconnnected) String() string { return fmt.Sprintf("%T", e) }

func (e *RequestEventClientDisconnnected) IsRequestEvent() {}
