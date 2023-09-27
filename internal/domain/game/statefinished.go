// nolint: dupl // Different states, different responsibilities.
package game

import (
	"context"
	"fmt"

	"github.com/hedhyw/semerr/pkg/v1/semerr"

	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/asyncmodel"
	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/entities"
)

type stateFinished struct {
	model *Model
}

func (s *stateFinished) HandleRequestEvent(
	ctx context.Context,
	event asyncmodel.RequestEvent,
) error {
	switch event := event.(type) {
	case *RequestEventGameStarted:
		return s.handleEventGameStarted(ctx)
	case *RequestEventPlayerJoined:
		return s.handleEventPlayerJoined(ctx, event)
	case *RequestEventWordGuessAttempted:
		return semerr.NewBadRequestError(errGameIsNotInProgress)
	case *RequestEventCanvasChanged:
		return semerr.NewBadRequestError(errGameIsNotInProgress)
	default:
		return entities.NewUnsupportedEventError(event)
	}
}

func (s stateFinished) handleEventPlayerJoined(
	ctx context.Context,
	event *RequestEventPlayerJoined,
) error {
	player := s.model.addPlayer(event.ClientID)

	return s.model.EmitResponses(ctx, &ResponseEventPlayerHello{
		Player: player,
	}, s.model.responseEventGameStateChanged())
}

func (s stateFinished) handleEventGameStarted(ctx context.Context) error {
	return s.model.startGame(ctx)
}

func (s stateFinished) String() string {
	return fmt.Sprintf("%T", s)
}

func (s stateFinished) MarshalText() (text []byte, err error) {
	return []byte(s.String()), nil
}
