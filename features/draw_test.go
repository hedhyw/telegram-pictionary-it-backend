package features_test

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hedhyw/telegram-pictionary-it-backend/features"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/game"
)

func TestDraw(t *testing.T) {
	t.Parallel()

	/*
		As a leader player, I want to be able to draw the pictionary word
		so that other players can guess.
	*/

	t.Run("The leaders draws on the canvas", func(t *testing.T) {
		t.Parallel()

		// Given there is a started game with two players.
		ctx := features.Context(t)
		th := features.NewTestHelper(t)

		clientIDLeaderPlayer := uuid.NewString()
		clientIDGuesserPlayer := uuid.NewString()

		th.AddPlayer(t, ctx, clientIDLeaderPlayer)
		eventsGuesser := th.AddPlayer(t, ctx, clientIDGuesserPlayer)

		err := th.Game.Start(ctx)
		require.NoError(t, err)

		// When the leader draws the picture.
		image := getFakeImageBase64(t)

		err = th.Game.ChangeCanvas(ctx, clientIDLeaderPlayer, image)
		require.NoError(t, err)

		// Then the guesser player sees this picture.
		features.AwaitResponse(ctx, eventsGuesser, func(name string, scan features.EventScanFunc) bool {
			var value game.ResponseEventCanvasChanged

			if name != value.String() {
				return false
			}

			scan(t, &value)

			if value.ImageBase64 == "" {
				return false
			}

			assert.Equal(t, image, value.ImageBase64)

			return true
		})
	})
}

func getFakeImageBase64(tb testing.TB) string {
	tb.Helper()

	imgDot := image.NewGray(image.Rect(0, 0, 1, 1))

	var pngData bytes.Buffer

	err := png.Encode(&pngData, imgDot)
	require.NoError(tb, err)

	return base64.StdEncoding.EncodeToString(pngData.Bytes())
}
