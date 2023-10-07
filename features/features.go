package features

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	"github.com/hedhyw/telegram-pictionary-it-backend/internal/config"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/entities"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/game"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/domain/telegram"
	"github.com/hedhyw/telegram-pictionary-it-backend/internal/transport/clientshub"
)

// TestHelper is a helper for feature tests. It also creates a game in its
// initial state.
type TestHelper struct {
	Game       *game.Game
	ClientsHub *clientshub.Hub
}

func NewTestHelper(tb testing.TB) *TestHelper {
	tb.Helper()

	clientsHub := clientshub.New()

	game := game.New(game.Essentials{
		ClientsHub: clientsHub,
		Logger:     zerolog.New(os.Stdout).Level(zerolog.TraceLevel),
		ChatID:     uuid.NewString(),
		Config: &config.Config{
			ServerTimeout:    10 * time.Second,
			DebugEnabled:     true,
			WorkersCount:     8,
			GameRoundTimeout: time.Minute,
		},
	})

	return &TestHelper{
		Game:       game,
		ClientsHub: clientsHub,
	}
}

// AddPlayer registers a new player in the game.
//
// nolint: revive // tb also should be first.
func (th *TestHelper) AddPlayer(
	tb testing.TB,
	ctx context.Context,
	clientID string,
	setters ...func(metadata *telegram.InitDataMeta),
) <-chan json.RawMessage {
	tb.Helper()

	eventsCh := make(chan json.RawMessage, 8)

	th.ClientsHub.AddClient(clientID, eventsCh)

	metadata := getFakePlayerMetadata()

	for _, s := range setters {
		s(metadata)
	}

	err := th.Game.AddPlayer(ctx, clientID, metadata)
	require.NoError(tb, err)

	return eventsCh
}

// Type definitions for callback functions.
type (
	EventScanFunc   func(tb testing.TB, target any)
	EventHandleFunc func(name string, scan EventScanFunc) bool
)

// AwaitEvent scans eventsCh until it finds the event T.
func AwaitEvent[T fmt.Stringer](
	tb testing.TB,
	eventsCh <-chan json.RawMessage,
) func(ctx context.Context) T {
	tb.Helper()

	return func(ctx context.Context) T {
		var value T

		AwaitResponse(ctx, eventsCh, func(name string, scan EventScanFunc) bool {
			tb.Logf("got event %s", name)

			if value.String() == name {
				scan(tb, &value)

				return true
			}

			return false
		})

		return value
	}
}

// AwaitResponse scans eventsCh and calls the callback function `handleFunc`
// until the last returns `true`.
func AwaitResponse(
	ctx context.Context,
	eventsCh <-chan json.RawMessage,
	handleFunc EventHandleFunc,
) {
	for {
		select {
		case <-ctx.Done():
		case event := <-eventsCh:
			var resp entities.SocketResponse[json.RawMessage]

			err := json.Unmarshal(event, &resp)
			if err != nil {
				continue
			}

			shouldStop := handleFunc(resp.Name, func(tb testing.TB, target any) {
				err := json.Unmarshal(resp.Payload, target)
				require.NoError(tb, err, string(resp.Payload))
			})
			if shouldStop {
				return
			}
		}
	}
}

func getFakePlayerMetadata() *telegram.InitDataMeta {
	return &telegram.InitDataMeta{
		UserJSON:     nil,
		ChatInstance: uuid.NewString(),
		ChatType:     "",
		Raw:          url.Values{},
	}
}

const defaultTimeout = 30 * time.Second

// Context returns a test context with a timeout set to defaultTimeout.
func Context(tb testing.TB) context.Context {
	tb.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	tb.Cleanup(cancel)

	return ctx
}
