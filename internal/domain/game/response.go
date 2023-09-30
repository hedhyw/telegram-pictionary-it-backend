package game

import (
	"fmt"
	"time"

	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/asyncmodel"
	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/player"

	"github.com/samber/lo"
)

type ResponseEventCanvasChanged struct {
	Players       []*player.Model `json:"-"`
	ActorClientID string          `json:"actorClientId"`

	ImageBase64 string `json:"imageBase64"`
}

func (e ResponseEventCanvasChanged) String() string { return fmt.Sprintf("%T", e) }

func (e ResponseEventCanvasChanged) TargetClientIDs() []string {
	return lo.Filter(getPlayerClientIDs(e.Players), func(item string, index int) bool {
		return item != e.ActorClientID
	})
}

func (e ResponseEventCanvasChanged) IsResponseEvent() {}

type ResponseEventGameStateChanged struct {
	Players  []*player.Model `json:"players"`
	FinishAt time.Time       `json:"finishAt"`

	State asyncmodel.State `json:"state"`
}

func (e ResponseEventGameStateChanged) String() string { return fmt.Sprintf("%T", e) }

func (e ResponseEventGameStateChanged) TargetClientIDs() []string {
	return getPlayerClientIDs(e.Players)
}

func (e ResponseEventGameStateChanged) IsResponseEvent() {}

type ResponseEventPlayerGuessed struct {
	Players []*player.Model `json:"-"`

	ClientID string `json:"clientId"`
}

func (e ResponseEventPlayerGuessed) String() string { return fmt.Sprintf("%T", e) }

func (e ResponseEventPlayerGuessed) TargetClientIDs() []string {
	return getPlayerClientIDs(e.Players)
}

func (e ResponseEventPlayerGuessed) IsResponseEvent() {}

type ResponseEventGameStarted struct {
	Players []*player.Model `json:"-"`
}

func (e ResponseEventGameStarted) String() string { return fmt.Sprintf("%T", e) }

func (e ResponseEventGameStarted) IsResponseEvent() {}

func (e ResponseEventGameStarted) TargetClientIDs() []string {
	return getPlayerClientIDs(e.Players)
}

type ResponseEventPlayerGuessFailed struct {
	Players []*player.Model `json:"-"`

	ActorClientID string `json:"clientId"`
	Word          string `json:"word"`
}

func (e ResponseEventPlayerGuessFailed) TargetClientIDs() []string {
	return getPlayerClientIDs(e.Players)
}

func (e ResponseEventPlayerGuessFailed) String() string { return fmt.Sprintf("%T", e) }

func (e *ResponseEventPlayerGuessFailed) IsResponseEvent() {}

type ResponseEventPlayerHello struct {
	Player *player.Model `json:"player"`
}

func (e ResponseEventPlayerHello) TargetClientIDs() []string {
	return []string{e.Player.ClientID}
}

func (e ResponseEventPlayerHello) String() string { return fmt.Sprintf("%T", e) }

func (e *ResponseEventPlayerHello) IsResponseEvent() {}

func getPlayerClientIDs(players []*player.Model) []string {
	return lo.Map(players, func(player *player.Model, _ int) string {
		return player.ClientID
	})
}

type ResponseEventLeadHello struct {
	ClientID string `json:"-"`

	Word string `json:"word"`
}

func (e ResponseEventLeadHello) TargetClientIDs() []string {
	return []string{e.ClientID}
}

func (e ResponseEventLeadHello) String() string { return fmt.Sprintf("%T", e) }

func (e *ResponseEventLeadHello) IsResponseEvent() {}
