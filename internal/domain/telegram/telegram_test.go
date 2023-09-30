package telegram_test

import (
	"net/url"
	"testing"

	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/telegram"

	"github.com/stretchr/testify/assert"
)

func TestDecodeInitData(t *testing.T) {
	const (
		chatInstance = "1"
		chatType     = "private"
		user         = `{"id":2,"first_name":"F","last_name":"L","username":"U","language_code":"en","allows_write_to_pm":true}`
		expectedHash = "6d8a82ff4a2279b0927c87f3f4afc9171fe6ea1bae046ba329d777033a7a5303"
	)

	initDataRaw := url.Values{
		"chat_instance": []string{chatInstance},
		"chat_type":     []string{chatType},
		"user":          []string{user},
		"hash":          []string{expectedHash},
	}.Encode()

	decoder := telegram.NewDecoder("fake_bot_token")

	meta, err := decoder.DecodeInitData(initDataRaw)
	if assert.NoError(t, err) {
		assert.Equal(t, chatInstance, meta.ChatInstance)
		assert.Equal(t, chatType, meta.ChatType)

		if assert.NotNil(t, meta.Raw) {
			assert.Equal(t, expectedHash, meta.Raw.Get("hash"))
		}

		if assert.NotNil(t, meta.UserJSON) {
			assert.Equal(t, user, *meta.UserJSON)
		}
	}
}
