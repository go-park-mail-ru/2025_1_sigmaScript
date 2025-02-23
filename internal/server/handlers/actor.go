package handlers

import (
  "net/http"
  "strconv"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
  "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
  "github.com/gorilla/mux"
  "github.com/rs/zerolog/log"
)

func GetActor(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  log.Info().Msg("GetActor")

  id, err := strconv.Atoi(vars["id"])
  if err != nil {
    log.Error().Err(err).Msg("Invalid actor id")
    jsonutil.SendError(w, http.StatusBadRequest, "invalid_id", "Invalid actor id")
    return
  }

  actor, exists := mocks.Actors[id]
  if !exists {
    log.Info().Msg("Actor not found")
    jsonutil.SendError(w, http.StatusNotFound, "not_found", "Actor not found")
    return
  }

  if err = jsonutil.SendJSON(w, actor); err != nil {
    log.Error().Err(err).Msg("Error sending JSON")
    return
  }
}
