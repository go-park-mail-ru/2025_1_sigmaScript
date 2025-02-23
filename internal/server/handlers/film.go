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
    http.Error(w, "Invalid film ID", http.StatusBadRequest)
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
    http.Error(w, "Film not found", http.StatusNotFound)
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
