package websocketserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hedhyw/telegram-pictionary-backend/internal/config"
	"github.com/hedhyw/telegram-pictionary-backend/internal/domain/core"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
)

type WebSocketServer struct {
	essentials Essentials

	upgrader *websocket.Upgrader
}

type Essentials struct {
	Logger zerolog.Logger
	Core   *core.Core
	Config *config.Config
}

func New(es Essentials) *WebSocketServer {
	upgrader := &websocket.Upgrader{
		HandshakeTimeout: es.Config.ServerTimeout,
		ReadBufferSize:   es.Config.ServerReadBufferSize,
		WriteBufferSize:  es.Config.ServerWriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			actualOrigin := r.Header["Origin"]
			if len(actualOrigin) == 0 {
				return true
			}

			expectedOrigin := es.Config.ServerCheckOrigin
			ok := actualOrigin[0] == expectedOrigin

			es.Logger.Debug().Msgf("checking origin for %s: actual %q, expected %q, ok %t", r.URL, actualOrigin[0], expectedOrigin, ok)

			return ok
		},
	}

	return &WebSocketServer{
		essentials: es,
		upgrader:   upgrader,
	}
}

func (s WebSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	clientID := r.URL.Query().Get("client_id")

	logger := s.essentials.Logger.With().Str("client", clientID).Logger()

	connection, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Err(err).Msg("failed to upgrade connection")

		return
	}

	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	eventsCh := make(chan json.RawMessage)

	err = s.essentials.Core.RegisterClient(ctx, clientID, eventsCh)
	if err != nil {
		logger.Err(err).Msg("failed to register client")

		return
	}

	defer func() {
		err := s.essentials.Core.UnregisterClient(ctx, clientID)
		if err != nil {
			logger.Err(err).Msg("failed to unregister client")
		}
	}()

	go s.runConnectionWriter(ctx, connection, eventsCh)

	logger.Debug().Msg("a new websocket client connected")

	s.runConnectionReader(ctx, clientID, connection)

	if err = connection.Close(); err != nil {
		logger.Err(err).Msg("failed to close websocket connection")
	}
}

func (s WebSocketServer) runConnectionReader(
	ctx context.Context,
	clientID string,
	connection *websocket.Conn,
) {
	logger := s.essentials.Logger

	defer logger.Debug().Err(ctx.Err()).Msgf("closed websocket reader")

	for ctx.Err() == nil {
		var errClose *websocket.CloseError

		var event json.RawMessage

		err := connection.ReadJSON(&event)
		if err != nil {
			if errors.As(err, &errClose) {
				logger.Debug().Msg("websocket connection closed")

				return
			}

			logger.Err(err).Msg("failed to read json")

			continue
		}

		err = s.essentials.Core.EmitClientEvent(ctx, clientID, event)
		if err != nil {
			logger.Err(err).Interface("event", event).Msg("failed to handle event")

			return
		}

		logger.Debug().Interface("event", event).Msgf("received websocket event: %s", event)
	}
}

func (s WebSocketServer) runConnectionWriter(
	ctx context.Context,
	connection *websocket.Conn,
	eventsCh <-chan json.RawMessage,
) {
	logger := s.essentials.Logger

	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
			logger.Debug().Err(ctx.Err()).Msgf("closed websocket writer")

			return
		case event := <-eventsCh:
			err := connection.WriteJSON(event)
			if err != nil {
				logger.Warn().Interface("event", event).Err(err).Msg("failed to write websocket message")

				continue
			}

			logger.Debug().Interface("event", event).Msg("sent websocket message")
		}
	}
}
