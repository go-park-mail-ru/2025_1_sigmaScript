package handlers

import (
  "encoding/json"
  "net/http"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
  "github.com/rs/zerolog/log"
)

func GetGenres(w http.ResponseWriter, r *http.Request) {
  log.Info().
    Str("package", "handlers").
    Str("func", "GetGenres").
    Str("method", r.Method).
    Str("path", r.URL.Path).
    Msg("GetGenres")
  if err := json.NewEncoder(w).Encode(mocks.Genres); err != nil {
    log.Error().
      Err(err).
      Str("package", "handlers").
      Str("func", "GetGenres").
      Str("method", r.Method).
      Str("path", r.URL.Path).
      Msg("Encode error")
    http.Error(w, "Encode error", http.StatusInternalServerError)
    return
  }
}
