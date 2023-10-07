// nolint: dupl // Different states, different responsibilities.
package game

import (
	"context"
	"fmt"

	"github.com/hedhyw/semerr/pkg/v1/semerr"

	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/asyncmodel"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/entities"
)

// StateFinished is a state of the game after finishing the round. Players
// can start a new game at this state.
type StateFinished struct {
	model *Model
}

// HandleRequestEvent implements asyncmodel.State.
func (s *StateFinished) HandleRequestEvent(
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
	case *RequestEventPlayerRemoved:
		return s.model.removePlayer(ctx, event.ClientID)
	default:
		return entities.NewUnsupportedEventError(event)
	}
}

func (s StateFinished) handleEventPlayerJoined(
	ctx context.Context,
	event *RequestEventPlayerJoined,
) error {
	player := s.model.addPlayer(event.ClientID, event.Meta)

	return s.model.EmitResponses(ctx,
		&ResponseEventPlayerHello{
			Player: *player,
		},
		s.model.responseEventGameStateChanged(),
	)
}

func (s StateFinished) handleEventGameStarted(ctx context.Context) error {
	return s.model.startGame(ctx)
}

// String implements fmt.Stringer and asyncmodel.State.
func (s StateFinished) String() string {
	return fmt.Sprintf("%T", s)
}

// MarshalText implements encoding.TextMarshaler and asyncmodel.State.
func (s StateFinished) MarshalText() (text []byte, err error) {
	return []byte(s.String()), nil
}
