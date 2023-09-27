package core

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hedhyw/semerr/pkg/v1/semerr"

	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/asyncmodel"
	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/entities"
	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/game"
)

type stateInitial struct {
	model *Model
}

func (s *stateInitial) HandleRequestEvent(ctx context.Context, event asyncmodel.RequestEvent) error {
	switch event := event.(type) {
	case *RequestEventGameStarted:
		return s.handleGameStarted(ctx, event)
	case *RequestEventClientConnected:
		return s.handleClientConnected(ctx, event)
	case *RequestEventClientDisconnnected:
		return s.handleClientDisconnected(ctx, event)
	case *RequestEventCanvasChanged:
		return s.handleCanvasChanged(ctx, event)
	case *RequestEventWordGuessAttempted:
		return s.handleGuessAttempted(ctx, event)
	default:
		return entities.NewUnsupportedEventError(event)
	}
}

func (s *stateInitial) handleGameStarted(ctx context.Context, event *RequestEventGameStarted) (err error) {
	game, err := s.model.getGameByClient(event.ClientID)
	if err != nil {
		return fmt.Errorf("getting game: %w", err)
	}

	return game.Start(ctx)
}

func (s *stateInitial) handleClientConnected(ctx context.Context, event *RequestEventClientConnected) (err error) {
	logger := s.model.logger

	_, err = uuid.Parse(event.ClientID)
	if err != nil {
		return semerr.NewBadRequestError(fmt.Errorf("parsing uuid: %w", err))
	}

	game := s.createGameIfNotExists(ctx, event.ChatID)

	_, ok := s.model.clientIDToChatID[event.ClientID]
	if ok {
		return semerr.NewConflictError(errClientConflict)
	}

	s.model.clientIDToChatID[event.ClientID] = event.ChatID

	connection := s.model.clientsHub.AddClient(event.ClientID, event.EventsCh)

	logger.Debug().
		Str("client", event.ClientID).
		Str("chat", event.ChatID).
		Msgf("client connected to the chat %s", event.ChatID)

	return game.AddPlayer(ctx, connection.ClientID())
}

func (s *stateInitial) createGameIfNotExists(
	_ context.Context,
	chatID string,
) *game.Game {
	foundGame, ok := s.model.chatIDToGame[chatID]
	if ok {
		return foundGame
	}

	s.model.logger.Debug().Str("chat", chatID).Msgf("created a new game in the chat %s", chatID)

	// nolint: contextcheck // It is a constructor.
	createdGame := game.New(game.Essentials{
		Logger:     s.model.logger,
		ChatID:     chatID,
		ClientsHub: s.model.clientsHub,
	})

	s.model.chatIDToGame[chatID] = createdGame

	return createdGame
}

func (s *stateInitial) handleClientDisconnected(
	_ context.Context,
	event *RequestEventClientDisconnnected,
) (err error) {
	logger := s.model.logger

	delete(s.model.clientIDToChatID, event.ClientID)
	s.model.clientsHub.RemoveClient(event.ClientID)

	logger.Debug().
		Str("client", event.ClientID).
		Msgf("client %s disconnected", event.ClientID)

	return nil
}

func (s stateInitial) String() string {
	return fmt.Sprintf("%T", s)
}

func (s stateInitial) MarshalText() (text []byte, err error) {
	return []byte(s.String()), nil
}

func (s *stateInitial) handleCanvasChanged(ctx context.Context, event *RequestEventCanvasChanged) (err error) {
	game, err := s.model.getGameByClient(event.ClientID)
	if err != nil {
		return fmt.Errorf("getting game: %w", err)
	}

	return game.ChangeCanvas(ctx, event.ClientID, event.ImageBase64)
}

func (s *stateInitial) handleGuessAttempted(ctx context.Context, event *RequestEventWordGuessAttempted) (err error) {
	game, err := s.model.getGameByClient(event.ClientID)
	if err != nil {
		return fmt.Errorf("getting game: %w", err)
	}

	return game.GuessWord(ctx, event.ClientID, event.Word)
}
