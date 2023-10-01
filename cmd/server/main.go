package main

import (
	"context"
	"os"

	"github.com/hedhyw/telegram-pictionary-it-backend/internal/config"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/core"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/telegram"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/transport/httpserver"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/transport/telegrambot"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/transport/websocketserver"

	"github.com/caarlos0/env/v6"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	err := tgbotapi.SetLogger(telegram.LoggerAdapter{
		Logger: logger,
	})
	if err != nil {
		logger.Err(err).Msg("setting telegram logger")
	}

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

	telegramBot := telegrambot.New(telegrambot.Essentials{
		Config: config,
		Logger: logger,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := telegramBot.Run(ctx); err != nil {
			logger.Err(err).Msg("running telegram bot")
		}
	}()

	err = httpServer.ListenAndServe()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to parse config")

		return
	}
}
