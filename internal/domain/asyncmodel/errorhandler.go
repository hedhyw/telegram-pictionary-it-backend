package asyncmodel

import (
	"github.com/rs/zerolog"
)

// RequestErrorHandlerFunc is a error callback function.
type RequestErrorHandlerFunc func(err error, event RequestEvent)

// DefaultLogRequestErrorHandler is RequestErrorHandlerFunc
// that prints all errors to the logger.
func DefaultLogRequestErrorHandler(logger zerolog.Logger) RequestErrorHandlerFunc {
	return func(err error, event RequestEvent) {
		logger.Err(err).Msgf("failed to handle request %s", event)
	}
}
