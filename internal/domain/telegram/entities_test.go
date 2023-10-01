package telegram_test

import (
	"testing"

	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/telegram"

	"github.com/hedhyw/semerr/pkg/v1/semerr"
	"github.com/stretchr/testify/assert"
)

func TestInitDataMetaUser(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		userJSON := webAppUserJSON

		meta := telegram.InitDataMeta{UserJSON: &userJSON}

		user, err := meta.User()
		if assert.NoError(t, err) {
			assert.NotEmpty(t, user.ID)
			assert.NotEmpty(t, user.FirstName)
			assert.NotEmpty(t, user.LastName)
			assert.NotEmpty(t, user.Username)
		}
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		meta := telegram.InitDataMeta{
			UserJSON: nil,
		}

		_, err := meta.User()
		if assert.Error(t, err) {
			assert.ErrorAs(t, err, &semerr.NotFoundError{})
		}
	})

	t.Run("invalid_json", func(t *testing.T) {
		t.Parallel()

		userJSON := "-"

		meta := telegram.InitDataMeta{UserJSON: &userJSON}

		_, err := meta.User()
		if assert.Error(t, err) {
			assert.ErrorAs(t, err, &semerr.BadRequestError{})
		}
	})
}
