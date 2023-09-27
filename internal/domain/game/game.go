package game

import (
	"github.com/hedhyw/semerr/pkg/v1/semerr"
	"github.com/rs/zerolog"

	"github.com/hedhyw/telegram-pictionary-backend/internal/transport/clientshub"
)

const (
	errGameInProgress      semerr.Error = "the game is in progress"
	errNotEnoughPlayers    semerr.Error = "not enough players"
	errGameIsNotInProgress semerr.Error = "the game has not been started yet"
	errWordNotMatch        semerr.Error = "word is not matched"
	errPlayerNotFound      semerr.Error = "player is not found"
)

type Game struct {
	model *Model
	view  *view

	*Controller
}

type Essentials struct {
	ClientsHub *clientshub.Hub
	Logger     zerolog.Logger
	ChatID     string
}

func New(es Essentials) *Game {
	es.Logger = es.Logger.With().Str("chat", es.ChatID).Logger()

	model := newModel(es)
	view := newView(es, model)
	controller := newController(model)

	return &Game{
		view:       view,
		model:      model,
		Controller: controller,
	}
}
