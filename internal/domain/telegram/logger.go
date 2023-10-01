package telegram

import (
	"fmt"

	"github.com/rs/zerolog"
)

// LoggerAdapter implements tgbotapi.BotLogger for zerolog.Logger.
// All logs are printed in debug level.
type LoggerAdapter struct {
	zerolog.Logger
}

// Println implements tgbotapi.BotLogger.
func (a LoggerAdapter) Println(v ...interface{}) {
	a.Logger.Debug().Msg(fmt.Sprint(v...))
}

// Printf implements tgbotapi.BotLogger.
func (a LoggerAdapter) Printf(format string, v ...interface{}) {
	a.Logger.Debug().Msgf(format, v...)
}
