package asyncmodel

import (
	"github.com/rs/zerolog"
)

type RequestErrorHandlerFunc func(err error, event RequestEvent)

func DefaultLogRequestErrorHandler(logger zerolog.Logger) RequestErrorHandlerFunc {
	return func(err error, event RequestEvent) {
		logger.Err(err).Msgf("failed to handle request %s", event)
	}
}
