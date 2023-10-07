package httpserver

import (
	"net/http"

	"github.com/rs/zerolog"

	"github.com/hedhyw/telegram-pictionary-it-backend/internal/config"
)

// Server serves HTTP requests.
type Server struct {
	essentials Essentials

	server *http.Server
}

// Essentials contains the required arguments for New.
type Essentials struct {
	Logger           zerolog.Logger
	Config           config.Config
	WebSocketHandler http.Handler
}

func New(es Essentials) *Server {
	return &Server{
		essentials: es,
	}
}

func (s *Server) ListenAndServe() error {
	logger := s.essentials.Logger

	mux := http.NewServeMux()

	mux.Handle("/api/websocket", s.essentials.WebSocketHandler)
	mux.HandleFunc("/debug/health", s.handleHealth)

	s.server = &http.Server{
		Addr:         s.essentials.Config.ServerAddress,
		Handler:      mux,
		WriteTimeout: s.essentials.Config.ServerTimeout,
		ReadTimeout:  s.essentials.Config.ServerTimeout,
	}

	logger.Info().Msgf("the server is started, try http://%s/debug/health", s.server.Addr)

	return s.server.ListenAndServe()
}

func (s Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	logger := s.essentials.Logger

	_, err := w.Write([]byte("OK"))
	if err != nil {
		logger.Err(err).Msg("faield to write health response")
	}
}
