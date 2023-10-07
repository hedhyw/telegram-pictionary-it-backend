package features_test

import (
	"encoding/json"
	"testing"

	"github.com/hedhyw/telegram-pictionary-it-backend/features"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/game"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/telegram"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistration(t *testing.T) {
	t.Parallel()

	/*
		As a user, I want to be able to join the room
		so that I can participate in the gameplay.
	*/

	t.Run("A user joins the room for the first time", func(t *testing.T) {
		t.Parallel()

		ctx := features.Context(t)
		th := features.NewTestHelper(t)

		// When the user joins the room.
		clientID := uuid.NewString()
		username := uuid.NewString()

		events := th.AddPlayer(t, ctx, clientID, func(metadata *telegram.InitDataMeta) {
			webAppUserJSON, err := json.Marshal(telegram.WebAppUser{
				Username: username,
			})
			require.NoError(t, err)

			webAppUserText := string(webAppUserJSON)
			metadata.UserJSON = &webAppUserText
		})

		// Then they are showed in the list of players.
		gameStateChanged := features.AwaitEvent[game.ResponseEventGameStateChanged](t, events)(ctx)
		require.Len(t, gameStateChanged.Players, 1)

		player := gameStateChanged.Players[0]
		assert.Equal(t, clientID, player.ClientID)

		// And they have the correct username.
		assert.Equal(t, username, player.Username)

		// And they have a zero score.
		assert.Zero(t, player.Score)
	})
}
