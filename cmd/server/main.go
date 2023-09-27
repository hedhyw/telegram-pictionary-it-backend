package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/hedhyw/telegram-pictionary-backend/internal/config"
	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/core"
	"github.com/hedhyw/telegram-pictionary-backend/internal/transport/httpserver"
	"github.com/hedhyw/telegram-pictionary-backend/internal/transport/websocketserver"

	"github.com/caarlos0/env/v6"
	"github.com/rs/zerolog"
)

func main() {
	rand.Seed(time.Now().Unix())

	// TODO: level from config.
	logger := zerolog.New(os.Stdout).Level(zerolog.DebugLevel).With().Timestamp().Logger()

	core := core.New(logger)

	var config config.Config

	if err := env.Parse(&config); err != nil {
		logger.Fatal().Err(err).Msg("failed to parse config")

		return
	}

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
