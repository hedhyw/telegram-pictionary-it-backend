package core

import (
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/asyncmodel"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/game"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/telegram"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/transport/clientshub"

	"github.com/hedhyw/semerr/pkg/v1/semerr"
)

// Model holds active games.
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

	model.Model = asyncmodel.New(asyncmodel.Essentials{
		InitialState:           &StateInitial{model: model},
		HandleRequestErrorFunc: asyncmodel.DefaultLogRequestErrorHandler(es.Logger),
		RequestTimeout:         es.Config.ServerTimeout,
		Logger:                 es.Logger,
		ChannelSize:            es.Config.WorkersCount,
	})

	return model
}

func (c *Model) getGameByClient(clientID string) (*game.Game, error) {
	game, ok := c.chatIDToGame[c.clientIDToChatID[clientID]]

	if !ok {
		return nil, semerr.NewNotFoundError(semerr.Error("game is not found"))
	}

	return game, nil
}
