package core

import (
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

func New(logger zerolog.Logger) *Core {
	model := newModel(logger)
	controller := newController(model)

	return &Core{
		Controller: controller,
		model:      model,
	}
}
