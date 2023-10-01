package telegram_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestLoggerAdapter(t *testing.T) {
	var logBuf bytes.Buffer

	var logger tgbotapi.BotLogger = telegram.LoggerAdapter{
		Logger: zerolog.New(&logBuf),
	}

	expected := fmt.Sprintf("logger.Printf: %s", t.Name())
	logger.Printf("logger.Printf: %s", t.Name())
	assert.Contains(t, logBuf.String(), expected)

	expected = fmt.Sprint("logger.Println", t.Name())
	logger.Println("logger.Println", t.Name())
	assert.Contains(t, logBuf.String(), expected)
}
