package router

import (
  "net/http"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/handlers"
  "github.com/gorilla/mux"
  "github.com/rs/zerolog/log"
)

func New() *mux.Router {
  log.Info().Msg("Configuring routes")

  router := mux.NewRouter()
  router.HandleFunc("/film/{id}", handlers.GetFilm).Methods(http.MethodGet)
  router.HandleFunc("/actor/{id}", handlers.GetActor).Methods(http.MethodGet)
  router.HandleFunc("/genres/", handlers.GetGenres).Methods(http.MethodGet)
  router.HandleFunc("/auth/login/", handlers.LoginHandler).Methods(http.MethodPost)
  router.HandleFunc("/auth/logout/", handlers.LogoutHandler).Methods(http.MethodPost)
  router.HandleFunc("/auth/register/", handlers.RegisterHandler).Methods(http.MethodPost)

  log.Info().Msg("Routes configured successfully")
  return router
}
