package jsonutil

import (
  "encoding/json"
  "net/http"

  "github.com/rs/zerolog/log"
)

func SendJSON(w http.ResponseWriter, data interface{}) error {
  log.Info().
    Str("package", "jsonutil").
    Str("method", "SendJSON").
    Msg("SendJSON")

  w.Header().Set("Content-Type", "application/json")
  if err := json.NewEncoder(w).Encode(data); err != nil {
    log.Error().
      Err(err).
      Str("package", "jsonutil").
      Str("method", "SendJSON").
      Msg("Error encoding JSON")
    http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
  }
  return nil
}
