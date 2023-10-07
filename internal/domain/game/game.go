package game

import (
	"github.com/hedhyw/semerr/pkg/v1/semerr"
	"github.com/rs/zerolog"

	"github.com/hedhyw/telegram-pictionary-it-backend/internal/config"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/transport/clientshub"
)

const (
	errGameInProgress      semerr.Error = "the game is in progress"
	errNotEnoughPlayers    semerr.Error = "not enough players"
	errGameIsNotInProgress semerr.Error = "the game has not been started yet"
	errWordNotMatch        semerr.Error = "word is not matched"
	errPlayerNotFound      semerr.Error = "player is not found"
	errPlayerNotLeader     semerr.Error = "player is not leader"
	errPlayerLeader        semerr.Error = "player is leader"
)

// Game is a facade for the view, the model and the controller.
type Game struct {
	*Model
	view *view

	*Controller
}

// Essentials contains the required arguments for New.
type Essentials struct {
	ClientsHub *clientshub.Hub
	Logger     zerolog.Logger
	ChatID     string
	Config     *config.Config
}

// New creates a new game facade.
func New(es Essentials) *Game {
	es.Logger = es.Logger.With().Str("chat", es.ChatID).Logger()

	model := newModel(es)
	view := newView(es, model)
	controller := newController(model)

	return &Game{
		view:       view,
		Model:      model,
		Controller: controller,
	}
}
