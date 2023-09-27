package game

import "fmt"

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
}

// TODO: rename to hello.
func (e RequestEventPlayerJoined) String() string { return fmt.Sprintf("%T", e) }

func (e *RequestEventPlayerJoined) IsRequestEvent() {}

type RequestEventWordGuessAttempted struct {
	ClientID string `json:"clientId"`
	Word     string `json:"word"`
}

func (e RequestEventWordGuessAttempted) String() string { return fmt.Sprintf("%T", e) }

func (e *RequestEventWordGuessAttempted) IsRequestEvent() {}
