package jsonutil

import (
  "encoding/json"
  "fmt"
  "net/http"

  "github.com/rs/zerolog/log"
)

type ErrorResponse struct {
  Error string `json:"error"`
  Msg   string `json:"msg"`
}

func SendError(w http.ResponseWriter, errCode int, errResp, msg string) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(errCode)

  errResponse := ErrorResponse{
    Error: errResp,
    Msg:   msg,
  }
  if err := json.NewEncoder(w).Encode(errResponse); err != nil {
    log.Error().Err(err).Msg("Failed to encode error response")
  }
}

func ReadJSON(r *http.Request, data interface{}) error {
  defer func() {
    err := r.Body.Close()
    if err != nil {
      log.Error().Err(err).Msg("Failed to close body")
    }
  }()
  if err := json.NewDecoder(r.Body).Decode(data); err != nil {
    return fmt.Errorf("error reading json: %w", err)
  }
  return nil
}

func SendJSON(w http.ResponseWriter, data interface{}) error {
  w.Header().Set("Content-Type", "application/json")
  if err := json.NewEncoder(w).Encode(data); err != nil {
    log.Error().Err(err).Msg("Error encoding JSON")
    SendError(w, http.StatusInternalServerError, "encode_error", "Error encoding JSON")
    return fmt.Errorf("error encoding JSON: %w", err)
  }
  return nil
}
