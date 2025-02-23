package handlers

import (
  "encoding/json"
  "net/http"
  "strconv"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
  "github.com/gorilla/mux"
  "github.com/rs/zerolog/log"
)

func GetActor(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  log.Info().
    Str("package", "handlers").
    Str("func", "GetActor").
    Str("method", r.Method).
    Str("path", r.URL.Path).
    Str("id", vars["id"]).
    Msg("GetActor")

  id, err := strconv.Atoi(vars["id"])
  if err != nil {
    log.Error().
      Err(err).
      Str("package", "handlers").
      Str("func", "GetActor").
      Str("method", r.Method).
      Str("path", r.URL.Path).
      Msg("Invalid actor id")
    http.Error(w, "Invalid actor ID", http.StatusBadRequest)
    return
  }

  actor, exists := mocks.Actors[id]
  if !exists {
    log.Info().
      Str("package", "handlers").
      Str("func", "GetActor").
      Str("method", r.Method).
      Str("path", r.URL.Path).
      Str("id", vars["id"]).
      Msg("Actor not found")
    http.Error(w, "Actor not found", http.StatusNotFound)
    return
  }

  if err = json.NewEncoder(w).Encode(actor); err != nil {
    log.Error().
      Err(err).
      Str("package", "handlers").
      Str("func", "GetActor").
      Str("method", r.Method).
      Str("path", r.URL.Path).
      Str("id", vars["id"]).
      Msg("Encode error")
    http.Error(w, "Encode error", http.StatusInternalServerError)
    return
  }
}
