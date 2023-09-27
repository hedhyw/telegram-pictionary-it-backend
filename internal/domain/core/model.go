package core

import (
	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/asyncmodel"
	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/game"
	"github.com/hedhyw/telegram-pictionary-backend/internal/transport/clientshub"

	"github.com/hedhyw/semerr/pkg/v1/semerr"
	"github.com/rs/zerolog"
)

type Model struct {
	logger zerolog.Logger

	clientsHub *clientshub.Hub

	*asyncmodel.Model

	clientIDToChatID map[string]string
	chatIDToGame     map[string]*game.Game
}

func newModel(logger zerolog.Logger) *Model {
	model := &Model{
		logger: logger,

		Model: nil,

		clientsHub: clientshub.New(),

		clientIDToChatID: map[string]string{},
		chatIDToGame:     map[string]*game.Game{},
	}

	model.Model = asyncmodel.New(
		&stateInitial{model: model},
		asyncmodel.DefaultLogRequestErrorHandler(logger),
	)

	return model
}

func (c *Model) getGameByClient(clientID string) (*game.Game, error) {
	game, ok := c.chatIDToGame[c.clientIDToChatID[clientID]]

	if !ok {
		return nil, semerr.NewNotFoundError(semerr.Error("game is not found"))
	}

	return game, nil
}
