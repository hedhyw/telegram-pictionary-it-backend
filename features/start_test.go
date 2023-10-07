package features_test

import (
	"testing"

	"github.com/hedhyw/telegram-pictionary-it-backend/features"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/game"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	t.Parallel()

	/*
		As a player, I want to be able to start the game
		so all members can start the gameplay.
	*/

	t.Run("Start the game", func(t *testing.T) {
		t.Parallel()

		ctx := features.Context(t)
		th := features.NewTestHelper(t)

		// Given there are two players in the game.
		eventsLeader := th.AddPlayer(t, ctx, uuid.NewString())
		th.AddPlayer(t, ctx, uuid.NewString())

		// When the user starts the game.
		err := th.Game.Start(ctx)
		require.NoError(t, err)

		// Then the game is started.
		var (
			leadHello game.ResponseEventLeadHello

			leadHelloFound   bool
			gameStartedFound bool
		)

		features.AwaitResponse(ctx, eventsLeader, func(name string, scan features.EventScanFunc) bool {
			if name == (game.ResponseEventLeadHello{}).String() {
				leadHelloFound = true
				scan(t, &leadHello)
			}

			if name == (game.ResponseEventGameStarted{}).String() {
				gameStartedFound = true
			}

			return leadHelloFound && gameStartedFound
		})

		_, ok := th.Game.State().(*game.StateInProgress)
		assert.Truef(t, ok, "%T", th.Game.State())

		// And the leader player receives a random word.
		assert.NotEmpty(t, leadHello.Word)
	})
}
