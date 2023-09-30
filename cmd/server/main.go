package main

import (
	"os"

	"github.com/hedhyw/telegram-pictionary-backend/internal/config"
	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/core"
	"github.com/hedhyw/telegram-pictionary-backend/internal/transport/httpserver"
	"github.com/hedhyw/telegram-pictionary-backend/internal/transport/websocketserver"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	var config config.Config

	if err := env.Parse(&config); err != nil {
		logger.Fatal().Err(err).Msg("failed to parse config")

		return
	}

	if config.DebugEnabled {
		logger = logger.Level(zerolog.DebugLevel)
	} else {
		logger = logger.Level(zerolog.InfoLevel)
	}

	logger.Debug().Any("config", config.Sanitized()).Msg("read config")

	core := core.New(core.Essentials{
		Logger: logger,
		Config: &config,
	})

	webSocketHandler := websocketserver.New(websocketserver.Essentials{
		Logger: logger,
		Core:   core,
		Config: &config,
	})

	httpServer := httpserver.New(httpserver.Essentials{
		Logger:           logger,
		Config:           config,
		WebSocketHandler: webSocketHandler,
	})

	err := httpServer.ListenAndServe()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse config")

		return
	}
}
