package jsonutil

import (
  "encoding/json"
  "net/http"

  "github.com/rs/zerolog/log"
)

type ErrorResponse struct {
  Error string `json:"error"`
  Msg   string `json:"msg"`
}

func SendError(w http.ResponseWriter, errCode int, errResp, msg string) {
  log.Info().
    Str("packet", "jsonutil").
    Str("func", "SendError").
    Msg("SendError")
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(errCode)

  errResponse := ErrorResponse{
    Error: errResp,
    Msg:   msg,
  }
  if err := json.NewEncoder(w).Encode(errResponse); err != nil {
    log.Error().
      Err(err).
      Str("packet", "jsonutil").
      Str("func", "SendError").
      Msg("Failed to encode error response")
  }
}

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
