package game

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hedhyw/semerr/pkg/v1/semerr"

	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/asyncmodel"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/entities"
)

// StateInProgress represents on-going game. The leader draws a picture,
// and guessers are trying to guess the word. It might be finished after
// some timeout.
type StateInProgress struct {
	model *Model
}

// HandleRequestEvent implements asyncmodel.State.
func (s *StateInProgress) HandleRequestEvent(
	ctx context.Context,
	event asyncmodel.RequestEvent,
) error {
	switch event := event.(type) {
	case *RequestEventGameStarted:
		return semerr.NewBadRequestError(errGameInProgress)
	case *RequestEventPlayerJoined:
		return s.handleEventPlayerJoined(ctx, event)
	case *RequestEventWordGuessAttempted:
		return s.handleWordGuessAttempted(ctx, event)
	case *RequestEventCanvasChanged:
		return s.handleCanvasChanged(ctx, event)
	case *RequestEventPlayerRemoved:
		return s.handlePlayerRemoved(ctx, event)
	default:
		return entities.NewUnsupportedEventError(event)
	}
}

func (s StateInProgress) handleCanvasChanged(
	ctx context.Context,
	event *RequestEventCanvasChanged,
) error {
	if s.model.getLeader().ClientID != event.ClientID {
		return errPlayerNotLeader
	}

	return s.model.EmitResponses(ctx, &ResponseEventCanvasChanged{
		UnixNano:      time.Now().UnixNano(),
		Players:       s.model.getPlayers(),
		ActorClientID: event.ClientID,
		ImageBase64:   event.ImageBase64,
	})
}

func (s StateInProgress) handleWordGuessAttempted(
	ctx context.Context,
	event *RequestEventWordGuessAttempted,
) error {
	if s.model.getLeader().ClientID == event.ClientID {
		return errPlayerLeader
	}

	logger := s.model.essentials.Logger.With().Str("client", event.ClientID).Logger()

	actualWord := normalizeWord(event.Word)
	expectedWord := normalizeWord(s.model.word)

	if actualWord != expectedWord {
		errEvent := s.model.EmitResponses(ctx, &ResponseEventPlayerGuessFailed{
			Players: s.model.getPlayers(),

			ActorClientID: event.ClientID,
			Word:          event.Word,
		})

		logger.Debug().Msgf("client faield to guess with the word %s", event.Word)

		return errors.Join(errEvent, errWordNotMatch)
	}

	for _, p := range s.model.players {
		if p.ClientID != event.ClientID {
			continue
		}

		if p.RoundWordMatched {
			logger.Debug().Msgf("round word is already guessed by %s", event.Word)

			return nil
		}

		logger.Debug().Msgf("client guessed the word %s", event.Word)

		p.SetRoundWordMatched()

		if s.model.isEveryoneGuessed() {
			err := s.model.finishGame(ctx)
			if err != nil {
				return fmt.Errorf("finishing game: %w", err)
			}
		}

		return s.model.EmitResponses(ctx,
			&ResponseEventPlayerGuessed{
				Players:  s.model.getPlayers(),
				ClientID: event.ClientID,
			},
			s.model.responseEventGameStateChanged(),
		)
	}

	return semerr.NewNotFoundError(errPlayerNotFound)
}

func (s StateInProgress) handlePlayerRemoved(
	ctx context.Context,
	event *RequestEventPlayerRemoved,
) error {
	logger := s.model.essentials.Logger.With().Str("client", event.ClientID).Logger()

	if s.model.getLeader().ClientID == event.ClientID {
		logger.Debug().Msg("leader disconnected, finishing game")

		err := s.model.finishGame(ctx)
		if err != nil {
			return fmt.Errorf("finishing game: %w", err)
		}
	} else if len(s.model.players) == 2 {
		logger.Debug().Msgf("not enough players (%d), finishing game", len(s.model.players))

		err := s.model.finishGame(ctx)
		if err != nil {
			return fmt.Errorf("finishing game: %w", err)
		}
	}

	return s.model.removePlayer(ctx, event.ClientID)
}

func (s StateInProgress) String() string {
	return fmt.Sprintf("%T", s)
}

// String implements fmt.Stringer and asyncmodel.State.
func (s StateInProgress) MarshalText() (text []byte, err error) {
	return []byte(s.String()), nil
}

// MarshalText implements encoding.TextMarshaler and asyncmodel.State.
func (s StateInProgress) handleEventPlayerJoined(
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
