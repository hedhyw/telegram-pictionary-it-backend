package core

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hedhyw/semerr/pkg/v1/semerr"

	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/asyncmodel"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/entities"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/game"
)

// StateInitial implements the main state of the core.
type StateInitial struct {
	model *Model
}

// HandleRequestEvent implements asyncmodel.State.
func (s *StateInitial) HandleRequestEvent(ctx context.Context, event asyncmodel.RequestEvent) error {
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
	case *RequestEventPlayerInitialized:
		return s.handlePlayerInitialized(ctx, event)
	default:
		return entities.NewUnsupportedEventError(event)
	}
}

func (s *StateInitial) handleGameStarted(ctx context.Context, event *RequestEventGameStarted) (err error) {
	game, err := s.model.getGameByClient(event.ClientID)
	if err != nil {
		return fmt.Errorf("getting game: %w", err)
	}

	return game.Start(ctx)
}

func (s *StateInitial) handleClientConnected(_ context.Context, event *RequestEventClientConnected) (err error) {
	_, err = uuid.Parse(event.ClientID)
	if err != nil {
		return semerr.NewBadRequestError(fmt.Errorf("parsing uuid: %w", err))
	}

	if s.model.clientsHub.HasClient(event.ClientID) {
		return semerr.NewConflictError(errClientConflict)
	}

	s.model.clientsHub.AddClient(event.ClientID, event.EventsCh)

	return nil
}

func (s *StateInitial) createGameIfNotExists(
	_ context.Context,
	chatID string,
) *game.Game {
	foundGame, ok := s.model.chatIDToGame[chatID]
	if ok {
		return foundGame
	}

	s.model.essentials.Logger.Debug().Str("chat", chatID).Msgf("created a new game in the chat %s", chatID)

	// nolint: contextcheck // It is a constructor.
	createdGame := game.New(game.Essentials{
		Logger:     s.model.essentials.Logger,
		ChatID:     chatID,
		ClientsHub: s.model.clientsHub,
		Config:     s.model.essentials.Config,
	})

	s.model.chatIDToGame[chatID] = createdGame

	return createdGame
}

func (s *StateInitial) handleClientDisconnected(
	ctx context.Context,
	event *RequestEventClientDisconnnected,
) (err error) {
	logger := s.model.essentials.Logger

	game, err := s.model.getGameByClient(event.ClientID)
	if err != nil {
		return fmt.Errorf("getting game: %w", err)
	}

	delete(s.model.clientIDToChatID, event.ClientID)
	s.model.clientsHub.RemoveClient(event.ClientID)

	logger.Debug().
		Str("client", event.ClientID).
		Msgf("client %s disconnected", event.ClientID)

	return game.RemovePlayer(ctx, event.ClientID)
}

func (s *StateInitial) handleCanvasChanged(ctx context.Context, event *RequestEventCanvasChanged) (err error) {
	game, err := s.model.getGameByClient(event.ClientID)
	if err != nil {
		return fmt.Errorf("getting game: %w", err)
	}

	return game.ChangeCanvas(ctx, event.ClientID, event.ImageBase64)
}

func (s *StateInitial) handleGuessAttempted(ctx context.Context, event *RequestEventWordGuessAttempted) (err error) {
	game, err := s.model.getGameByClient(event.ClientID)
	if err != nil {
		return fmt.Errorf("getting game: %w", err)
	}

	return game.GuessWord(ctx, event.ClientID, event.Word)
}

func (s *StateInitial) handlePlayerInitialized(ctx context.Context, event *RequestEventPlayerInitialized) (err error) {
	logger := s.model.essentials.Logger.With().Str("client", event.ClientID).Logger()

	meta, err := s.model.telegramDecoder.DecodeInitData(event.InitDataRaw)
	if err != nil {
		return fmt.Errorf("decoding init data: %w", err)
	}

	chatID := meta.ChatInstance

	game := s.createGameIfNotExists(ctx, chatID)
	s.model.clientIDToChatID[event.ClientID] = chatID

	logger.Debug().
		Str("chat", chatID).
		Msgf("client connected to the chat %s", chatID)

	return game.AddPlayer(ctx, event.ClientID, meta)
}

// String implements fmt.Stringer and asyncmodel.State.
func (s StateInitial) String() string {
	return fmt.Sprintf("%T", s)
}

// MarshalText implements encoding.TextMarshaler and asyncmodel.State.
func (s StateInitial) MarshalText() (text []byte, err error) {
	return []byte(s.String()), nil
}
