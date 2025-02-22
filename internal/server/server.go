package server

import (
  "context"
  "net/http"

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

func New() *Server {
  router := mux.NewRouter()
  s := &Server{
    Router: router,
    httpServer: &http.Server{
      Addr:    ":8080",
      Handler: router,
    },
  }
  s.configureRoutes()
  return s
}
