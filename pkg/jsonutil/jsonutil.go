package jsonutil

import (
  "encoding/json"
  "net/http"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
  "github.com/pkg/errors"
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
    log.Error().Err(errors.Wrap(err, errs.ErrEncodeJSON)).Msg(errors.Wrap(err, errs.ErrEncodeJSON).Error())
  }
}

func ReadJSON(r *http.Request, data interface{}) error {
  defer func() {
    if err := r.Body.Close(); err != nil {
      log.Error().Err(errors.Wrap(err, errs.ErrCloseBody)).Msg(errors.Wrap(err, errs.ErrCloseBody).Error())
    }
  }()
  if err := json.NewDecoder(r.Body).Decode(data); err != nil {
    return errors.Wrap(err, errs.ErrParseJSON)
  }
  return nil
}

func SendJSON(w http.ResponseWriter, data interface{}) error {
  w.Header().Set("Content-Type", "application/json")
  if err := json.NewEncoder(w).Encode(data); err != nil {
    log.Error().Err(errors.Wrap(err, errs.ErrEncodeJSON)).Msg(errors.Wrap(err, errs.ErrEncodeJSON).Error())
    SendError(w, http.StatusInternalServerError, errors.Wrap(err, errs.ErrEncodeJSONShort).Error(), errors.Wrap(err, errs.ErrEncodeJSON).Error())
    return errors.Wrap(err, errs.ErrParseJSON)
  }
  return nil
}
