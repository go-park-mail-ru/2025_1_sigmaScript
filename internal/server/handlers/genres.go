package handlers

import (
  "net/http"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
  "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
  "github.com/rs/zerolog/log"
)

func GetGenres(w http.ResponseWriter, r *http.Request) {
  log.Info().Msg("GetGenres")
  if err := jsonutil.SendJSON(w, mocks.Genres); err != nil {
    log.Error().Err(err).Msg("Error sending JSON")
    return
  }
}
