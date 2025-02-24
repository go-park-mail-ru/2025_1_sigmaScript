package router

import (
  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/handlers"
  "github.com/gorilla/mux"
  "github.com/rs/zerolog/log"
)

func New() *mux.Router {
  log.Info().Msg("Configuring routes")

  router := mux.NewRouter()
  router.HandleFunc("/film/{id}", handlers.GetFilm).Methods("GET")
  router.HandleFunc("/actor/{id}", handlers.GetActor).Methods("GET")
  router.HandleFunc("/genres/", handlers.GetGenres).Methods("GET")

  log.Info().Msg("Routes configured successfully")
  return router
}
