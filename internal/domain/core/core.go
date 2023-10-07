package core

import (
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/config"

	"github.com/hedhyw/semerr/pkg/v1/semerr"
	"github.com/rs/zerolog"
)

const (
	errClientConflict semerr.Error = "client id conflict"
)

// Core is a facade for the model and the controller.
type Core struct {
	*Controller

	model *Model
}

// Essentials contains the required arguments for New.
type Essentials struct {
	Logger zerolog.Logger
	Config *config.Config
}

// New creates a new core facade.
func New(es Essentials) *Core {
	model := newModel(es)
	controller := newController(model)

	return &Core{
		Controller: controller,
		model:      model,
	}
}
