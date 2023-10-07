package features_test

import (
	"testing"

	"github.com/hedhyw/telegram-pictionary-it-backend/features"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/game"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuess(t *testing.T) {
	t.Parallel()

	/*
		As a player, I want to be able to guess the word
		so that I can win the game.
	*/

	t.Run("All players guess the word correctly", func(t *testing.T) {
		t.Parallel()

		// Given there is a started game with two players.
		ctx := features.Context(t)
		th := features.NewTestHelper(t)

		clientIDLeaderPlayer := uuid.NewString()
		clientIDGuesserPlayer := uuid.NewString()

		eventsLeader := th.AddPlayer(t, ctx, clientIDLeaderPlayer)
		th.AddPlayer(t, ctx, clientIDGuesserPlayer)

		err := th.Game.Start(ctx)
		require.NoError(t, err)

		// When the guesser player guesses the word correctly.
		leadHello := features.AwaitEvent[game.ResponseEventLeadHello](t, eventsLeader)(ctx)

		err = th.Game.GuessWord(ctx, clientIDGuesserPlayer, leadHello.Word)
		require.NoError(t, err)

		features.AwaitEvent[game.ResponseEventPlayerGuessed](t, eventsLeader)(ctx)

		// Then the game is finished.
		_, ok := th.Game.State().(*game.StateFinished)
		assert.Truef(t, ok, "%T", th.Game.State())
	})

	t.Run("Guess the word incorrectly", func(t *testing.T) {
		t.Parallel()

		// Given there is a started game with two players.
		ctx := features.Context(t)
		th := features.NewTestHelper(t)

		clientIDLeaderPlayer := uuid.NewString()
		clientIDGuesserPlayer := uuid.NewString()

		eventsLeader := th.AddPlayer(t, ctx, clientIDLeaderPlayer)
		th.AddPlayer(t, ctx, clientIDGuesserPlayer)

		err := th.Game.Start(ctx)
		require.NoError(t, err)

		// When the guesser player guesses the word incorrectly.
		err = th.Game.GuessWord(ctx, clientIDGuesserPlayer, uuid.NewString())
		require.NoError(t, err)

		features.AwaitEvent[game.ResponseEventPlayerGuessFailed](t, eventsLeader)(ctx)

		// Then the game is still in progress.
		_, ok := th.Game.State().(*game.StateInProgress)
		assert.Truef(t, ok, "%T", th.Game.State())
	})

	t.Run("Some players guess the word correctly", func(t *testing.T) {
		t.Parallel()

		// Given there is a started game with three players.
		ctx := features.Context(t)
		th := features.NewTestHelper(t)

		const (
			clientIDLeaderPlayer  = "leader"
			clientIDGuesserFirst  = "guesser_first"
			clientIDGuesserSecond = "guesser_second"
		)

		eventsLeader := th.AddPlayer(t, ctx, clientIDLeaderPlayer)
		th.AddPlayer(t, ctx, clientIDGuesserFirst)
		th.AddPlayer(t, ctx, clientIDGuesserSecond)

		err := th.Game.Start(ctx)
		require.NoError(t, err)

		// When one of the guesser players guesses the word correctly.
		leadHello := features.AwaitEvent[game.ResponseEventLeadHello](t, eventsLeader)(ctx)

		err = th.Game.GuessWord(ctx, clientIDGuesserFirst, leadHello.Word)
		require.NoError(t, err)

		features.AwaitEvent[game.ResponseEventPlayerGuessed](t, eventsLeader)(ctx)

		// Then the game is still in progress.
		_, ok := th.Game.State().(*game.StateInProgress)
		assert.Truef(t, ok, "%T", th.Game.State())
	})
}
