package server

import (
  "context"
  "fmt"
  "net/http"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/config"
  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/handlers"
  "github.com/gorilla/mux"
)

type Server struct {
  Router     *mux.Router
  httpServer *http.Server
}

func (s *Server) configureRoutes() {
  s.Router.HandleFunc("/film/{id}", handlers.GetFilm).Methods("GET")
  s.Router.HandleFunc("/actor/{id}", handlers.GetActor).Methods("GET")
  s.Router.HandleFunc("/genres/", handlers.GetGenres).Methods("GET")
}

func (s *Server) Run() error {
  return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
  return s.httpServer.Shutdown(ctx)
}

func New(srv *config.Server) *Server {
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
  return s
}
