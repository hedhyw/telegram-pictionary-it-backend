package core

import (
	"github.com/hedhyw/telegram-pictionary-backend/internal/config"

	"github.com/hedhyw/semerr/pkg/v1/semerr"
	"github.com/rs/zerolog"
)

const (
	errClientConflict semerr.Error = "client id conflict"
)

type Core struct {
	*Controller

	model *Model
}

type Essentials struct {
	Logger zerolog.Logger
	Config *config.Config
}

func New(es Essentials) *Core {
	model := newModel(es)
	controller := newController(model)

	return &Core{
		Controller: controller,
		model:      model,
	}
}
