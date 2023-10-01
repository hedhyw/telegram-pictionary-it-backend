package telegrambot

import (
	"context"
	"fmt"

	"github.com/hedhyw/telegram-pictionary-it-backend/internal/assets"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
)

// TelegramBot is a controller for Telegram Bot.
type TelegramBot struct {
	essentials Essentials
}

// Essentials contains required arguments for TelegramBot.
type Essentials struct {
	Config config.Config
	Logger zerolog.Logger
}

// New creates a new TelegramBot controller.
func New(es Essentials) *TelegramBot {
	return &TelegramBot{
		essentials: es,
	}
}

// Run start the bot's updates listener.
func (b TelegramBot) Run(ctx context.Context) error {
	logger := b.essentials.Logger

	const offset = 0

	bot, err := tgbotapi.NewBotAPI(b.essentials.Config.TelegramBotToken)
	if err != nil {
		return fmt.Errorf("creating bot api: %w", err)
	}

	go func() {
		<-ctx.Done()

		bot.StopReceivingUpdates()
	}()

	bot.Debug = b.essentials.Config.DebugEnabled

	updateConfig := tgbotapi.NewUpdate(offset)
	updateConfig.Timeout = int(b.essentials.Config.ServerTimeout.Seconds())

	updates := bot.GetUpdatesChan(updateConfig)

	logger.Info().Msg("telegram bot is running")

	for update := range updates {
		err = b.handleUpdate(ctx, bot, update)
		if err != nil {
			logger.Err(err).Int("update", update.UpdateID).Msg("failed to handle an update")
		}
	}

	return nil
}

func (b TelegramBot) handleUpdate(
	_ context.Context,
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
) error {
	logger := b.essentials.Logger.With().Int("update", update.UpdateID).Logger()

	if update.Message == nil {
		logger.Debug().Msg("skipping update, not a message")

		return nil
	}

	if update.Message.Chat == nil {
		logger.Warn().Msgf("skipping update %d, chat is nil", update.UpdateID)

		return nil
	}

	message := tgbotapi.NewMessage(
		update.Message.Chat.ID,
		assets.Hello(),
	)

	message.ParseMode = tgbotapi.ModeHTML
	message.DisableWebPagePreview = true

	_, err := bot.Send(message)
	if err != nil {
		return fmt.Errorf("sending plain message: %w", err)
	}

	_, err = bot.Send(tgbotapi.NewAnimation(
		update.Message.Chat.ID,
		tgbotapi.FileURL(b.essentials.Config.AnimationURL)),
	)
	if err != nil {
		return fmt.Errorf("sending animation: %w", err)
	}

	logger.Debug().Msgf("handled update to %d", update.Message.Chat.ID)

	return nil
}
