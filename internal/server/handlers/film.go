package handlers

import (
  "net/http"
  "strconv"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
  "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
  "github.com/gorilla/mux"
  "github.com/rs/zerolog/log"
)

func GetFilm(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  log.Info().
    Str("package", "handlers").
    Str("func", "GetFilm").
    Str("method", r.Method).
    Str("path", r.URL.Path).
    Str("id", vars["id"]).
    Msg("GetFilm")

  id, err := strconv.Atoi(vars["id"])
  if err != nil {
    log.Error().
      Err(err).
      Str("package", "handlers").
      Str("func", "GetFilm").
      Str("method", r.Method).
      Str("path", r.URL.Path).
      Msg("Invalid film id")
    jsonutil.SendError(w, http.StatusBadRequest, "invalid_id", "Invalid film id")
    return
  }

  film, exists := mocks.Films[id]
  if !exists {
    log.Info().
      Str("package", "handlers").
      Str("func", "GetFilm").
      Str("method", r.Method).
      Str("path", r.URL.Path).
      Str("id", vars["id"]).
      Msg("Film not found")
    jsonutil.SendError(w, http.StatusNotFound, "not_found", "Film not found")
    return
  }

  if err = jsonutil.SendJSON(w, film); err != nil {
    log.Error().
      Err(err).
      Str("package", "handlers").
      Str("func", "GetFilm").
      Str("method", r.Method).
      Str("path", r.URL.Path).
      Str("id", vars["id"]).
      Msg("Error sending JSON")
    return
  }
}
