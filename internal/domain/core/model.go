package core

import (
	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/asyncmodel"
	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/game"
	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/telegram"
	"github.com/hedhyw/telegram-pictionary-backend/internal/transport/clientshub"

	"github.com/hedhyw/semerr/pkg/v1/semerr"
)

type Model struct {
	telegramDecoder *telegram.Decoder
	essentials      Essentials

	clientsHub *clientshub.Hub

	*asyncmodel.Model

	clientIDToChatID map[string]string
	chatIDToGame     map[string]*game.Game
}

func newModel(es Essentials) *Model {
	model := &Model{
		telegramDecoder: telegram.NewDecoder(es.Config.TelegramBotToken),
		essentials:      es,

		Model: nil,

		clientsHub: clientshub.New(),

		clientIDToChatID: map[string]string{},
		chatIDToGame:     map[string]*game.Game{},
	}

	model.Model = asyncmodel.New(
		&stateInitial{model: model},
		asyncmodel.DefaultLogRequestErrorHandler(es.Logger),
		es.Config.ServerTimeout,
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
