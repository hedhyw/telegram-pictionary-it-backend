package game

import (
	"fmt"
	"time"

	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/player"

	"github.com/samber/lo"
)

// ResponseEventCanvasChanged implements asyncmodel.ResponseEvent.
// It notifies other players about new drawing.
type ResponseEventCanvasChanged struct {
	UnixNano int64 `json:"unixNano"`

	Players       []player.Model `json:"-"`
	ActorClientID string         `json:"actorClientId"`

	ImageBase64 string `json:"imageBase64"`
}

// String implements fmt.Stringer and asyncmodel.ResponseEvent.
func (e ResponseEventCanvasChanged) String() string { return fmt.Sprintf("%T", e) }

// TargetClientIDs implements asyncmodel.ResponseEvent.
func (e ResponseEventCanvasChanged) TargetClientIDs() []string {
	return lo.Filter(getPlayerClientIDs(e.Players), func(item string, index int) bool {
		return item != e.ActorClientID
	})
}

// IsResponseEvent implements asyncmodel.ResponseEvent.
func (e ResponseEventCanvasChanged) IsResponseEvent() {}

// ResponseEventGameStateChanged implements asyncmodel.ResponseEvent.
// It notifies all players about a new game state.
type ResponseEventGameStateChanged struct {
	UnixNano int64 `json:"unixNano"`

	Players  []player.Model `json:"players"`
	FinishAt time.Time      `json:"finishAt"`

	State string `json:"state"`

	Word string `json:"word"`
	Hint string `json:"hint"`
}

// String implements fmt.Stringer and asyncmodel.ResponseEvent.
func (e ResponseEventGameStateChanged) String() string { return fmt.Sprintf("%T", e) }

// TargetClientIDs implements asyncmodel.ResponseEvent.
func (e ResponseEventGameStateChanged) TargetClientIDs() []string {
	return getPlayerClientIDs(e.Players)
}

// IsResponseEvent implements asyncmodel.ResponseEvent.
func (e ResponseEventGameStateChanged) IsResponseEvent() {}

// ResponseEventPlayerGuessed implements asyncmodel.ResponseEvent.
// It notifies all players about successful guess.
type ResponseEventPlayerGuessed struct {
	Players []player.Model `json:"-"`

	ClientID string `json:"clientId"`
}

// String implements fmt.Stringer and asyncmodel.ResponseEvent.
func (e ResponseEventPlayerGuessed) String() string { return fmt.Sprintf("%T", e) }

// TargetClientIDs implements asyncmodel.ResponseEvent.
func (e ResponseEventPlayerGuessed) TargetClientIDs() []string {
	return getPlayerClientIDs(e.Players)
}

// IsResponseEvent implements asyncmodel.ResponseEvent.
func (e ResponseEventPlayerGuessed) IsResponseEvent() {}

// ResponseEventGameStarted implements asyncmodel.ResponseEvent.
// It notifies all players about start of the game.
type ResponseEventGameStarted struct {
	Players []player.Model `json:"-"`
}

// String implements fmt.Stringer and asyncmodel.ResponseEvent.
func (e ResponseEventGameStarted) String() string { return fmt.Sprintf("%T", e) }

// IsResponseEvent implements asyncmodel.ResponseEvent.
func (e ResponseEventGameStarted) IsResponseEvent() {}

// TargetClientIDs implements asyncmodel.ResponseEvent.
func (e ResponseEventGameStarted) TargetClientIDs() []string {
	return getPlayerClientIDs(e.Players)
}

// ResponseEventPlayerGuessFailed implements asyncmodel.ResponseEvent.
// It notifies all players about failed guess.
type ResponseEventPlayerGuessFailed struct {
	Players []player.Model `json:"-"`

	ActorClientID string `json:"clientId"`
	Word          string `json:"word"`
}

// TargetClientIDs implements asyncmodel.ResponseEvent.
func (e ResponseEventPlayerGuessFailed) TargetClientIDs() []string {
	return getPlayerClientIDs(e.Players)
}

// String implements fmt.Stringer and asyncmodel.ResponseEvent.
func (e ResponseEventPlayerGuessFailed) String() string { return fmt.Sprintf("%T", e) }

// IsResponseEvent implements asyncmodel.ResponseEvent.
func (e *ResponseEventPlayerGuessFailed) IsResponseEvent() {}

// ResponseEventPlayerHello implements asyncmodel.ResponseEvent.
// It sends a player model to the player who has just joined.
type ResponseEventPlayerHello struct {
	Player player.Model `json:"player"`
}

// TargetClientIDs implements asyncmodel.ResponseEvent.
func (e ResponseEventPlayerHello) TargetClientIDs() []string {
	return []string{e.Player.ClientID}
}

// String implements fmt.Stringer and asyncmodel.ResponseEvent.
func (e ResponseEventPlayerHello) String() string { return fmt.Sprintf("%T", e) }

// IsResponseEvent implements asyncmodel.ResponseEvent.
func (e *ResponseEventPlayerHello) IsResponseEvent() {}

func getPlayerClientIDs(players []player.Model) []string {
	return lo.Map(players, func(player player.Model, _ int) string {
		return player.ClientID
	})
}

// ResponseEventLeadHello implements asyncmodel.ResponseEvent.
// It sends the word to draw to the current round leader.
type ResponseEventLeadHello struct {
	ClientID string `json:"-"`

	Word string `json:"word"`
}

// TargetClientIDs implements asyncmodel.ResponseEvent.
func (e ResponseEventLeadHello) TargetClientIDs() []string {
	return []string{e.ClientID}
}

// String implements fmt.Stringer and asyncmodel.ResponseEvent.
func (e ResponseEventLeadHello) String() string { return fmt.Sprintf("%T", e) }

// IsResponseEvent implements asyncmodel.ResponseEvent.
func (e *ResponseEventLeadHello) IsResponseEvent() {}
