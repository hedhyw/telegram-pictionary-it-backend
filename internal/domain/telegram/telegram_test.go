package telegram_test

import (
	"net/url"
	"testing"

	"github.com/hedhyw/semerr/pkg/v1/semerr"
	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/telegram"

	"github.com/stretchr/testify/assert"
)

const (
	webAppUserJSON = `{"id":2,"first_name":"F","last_name":"L","username":"U","language_code":"en","allows_write_to_pm":true}`
	chatInstance   = "1"
	chatType       = "private"
	expectedHash   = "6d8a82ff4a2279b0927c87f3f4afc9171fe6ea1bae046ba329d777033a7a5303"
	testBotToken   = "fake_bot_token"
)

func TestDecodeInitDataSuccess(t *testing.T) {
	t.Parallel()

	initDataRaw := url.Values{
		"chat_instance": []string{chatInstance},
		"chat_type":     []string{chatType},
		"user":          []string{webAppUserJSON},
		"hash":          []string{expectedHash},
	}.Encode()

	decoder := telegram.NewDecoder(testBotToken)

	meta, err := decoder.DecodeInitData(initDataRaw)
	if assert.NoError(t, err) {
		assert.Equal(t, chatInstance, meta.ChatInstance)
		assert.Equal(t, chatType, meta.ChatType)

		if assert.NotNil(t, meta.Raw) {
			assert.Equal(t, expectedHash, meta.Raw.Get("hash"))
		}

		if assert.NotNil(t, meta.UserJSON) {
			assert.Equal(t, webAppUserJSON, *meta.UserJSON)
		}
	}
}

func TestDecodeInitDataEmptyHash(t *testing.T) {
	t.Parallel()

	initDataRaw := url.Values{
		"chat_instance": []string{"1"},
		"chat_type":     []string{"private"},
		"user":          []string{webAppUserJSON},
	}.Encode()

	decoder := telegram.NewDecoder(testBotToken)

	_, err := decoder.DecodeInitData(initDataRaw)
	if assert.Error(t, err) {
		assert.ErrorAs(t, err, &semerr.BadRequestError{})
	}
}

func TestDecodeInitDataInvaliHash(t *testing.T) {
	t.Parallel()

	initDataRaw := url.Values{
		"chat_instance": []string{"1"},
		"chat_type":     []string{"private"},
		"user":          []string{webAppUserJSON},
		"hash":          []string{"invalid"},
	}.Encode()

	decoder := telegram.NewDecoder(testBotToken)

	_, err := decoder.DecodeInitData(initDataRaw)
	if assert.Error(t, err) {
		assert.ErrorAs(t, err, &semerr.BadRequestError{})
	}
}

func TestDecodeInitDataEmptyBotToken(t *testing.T) {
	t.Parallel()

	initDataRaw := url.Values{
		"chat_instance": []string{"1"},
		"chat_type":     []string{"private"},
		"user":          []string{webAppUserJSON},
		"hash":          []string{expectedHash},
	}.Encode()

	decoder := telegram.NewDecoder("")

	_, err := decoder.DecodeInitData(initDataRaw)
	if assert.Error(t, err) {
		assert.ErrorAs(t, err, &semerr.BadRequestError{})
	}
}

func TestDecodeInitDataInvalidURLQuery(t *testing.T) {
	t.Parallel()

	decoder := telegram.NewDecoder(testBotToken)

	_, err := decoder.DecodeInitData("%&")
	if assert.Error(t, err) {
		assert.ErrorAs(t, err, &semerr.BadRequestError{})
	}
}
