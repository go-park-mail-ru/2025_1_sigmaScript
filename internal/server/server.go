package server

import (
  "context"
  "fmt"
  "net/http"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/config"
  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/router"
  "github.com/gorilla/mux"
  "github.com/rs/zerolog/log"
)

type Server struct {
  Router     *mux.Router
  httpServer *http.Server
}

func (s *Server) Run() error {
  log.Info().Msg("Running server")
  return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
  log.Info().Msg("Shutting down server")
  return s.httpServer.Shutdown(ctx)
}

func New(srv *config.Server) *Server {
  log.Info().Msg("Initializing server")
  mx := router.New()
  s := &Server{
    Router: mx,
    httpServer: &http.Server{
      Addr:         fmt.Sprintf("%s:%d", srv.Address, srv.Port),
      ReadTimeout:  srv.ReadTimeout,
      WriteTimeout: srv.WriteTimeout,
      IdleTimeout:  srv.IdleTimeout,
      Handler:      mx,
    },
  }

  log.Info().Msg("Server initialized successfully")
  return s
}
