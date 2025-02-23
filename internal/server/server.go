package server

import (
  "context"
  "fmt"
  "net/http"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/config"
  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/handlers"
  "github.com/gorilla/mux"
  "github.com/rs/zerolog/log"
)

type Server struct {
  Router     *mux.Router
  httpServer *http.Server
}

func (s *Server) configureRoutes() {
  log.Info().
    Str("package", "server").
    Str("func", "configureRoutes").
    Msg("Configuring routes")
  s.Router.HandleFunc("/film/{id}", handlers.GetFilm).Methods("GET")
  s.Router.HandleFunc("/actor/{id}", handlers.GetActor).Methods("GET")
  s.Router.HandleFunc("/genres/", handlers.GetGenres).Methods("GET")
  log.Info().
    Str("package", "server").
    Str("func", "configureRoutes").
    Msg("Routes configured successfully")
}

func (s *Server) Run() error {
  log.Info().
    Str("package", "server").
    Str("func", "Run").
    Msg("Running server")
  return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
  log.Info().
    Str("package", "server").
    Str("func", "Shutdown").
    Msg("Shutting down server")
  return s.httpServer.Shutdown(ctx)
}

func New(srv *config.Server) *Server {
  log.Info().
    Str("package", "server").
    Str("func", "New").
    Msg("Initializing server")
  router := mux.NewRouter()
  s := &Server{
    Router: router,
    httpServer: &http.Server{
      Addr:         fmt.Sprintf("%s:%d", srv.Address, srv.Port),
      ReadTimeout:  srv.ReadTimeout,
      WriteTimeout: srv.WriteTimeout,
      IdleTimeout:  srv.IdleTimeout,
      Handler:      router,
    },
  }
  s.configureRoutes()
  log.Info().
    Str("package", "server").
    Str("func", "New").
    Msg("Server initialized successfully")
  return s
}
