package jsonutil

import (
	"context"
	"encoding/json"
	"net/http"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func SendError(ctx context.Context, w http.ResponseWriter, errCode int, errResp, msg string) {
	logger := log.Ctx(ctx)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errCode)

	errResponse := ErrorResponse{
		Error:   errResp,
		Message: msg,
	}
	if err := json.NewEncoder(w).Encode(errResponse); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrEncodeJSON)).Msg(errors.Wrap(err, errs.ErrEncodeJSON).Error())
	}
}

func ReadJSON(r *http.Request, data interface{}) error {
	logger := log.Ctx(r.Context())

	defer func() {
		if err := r.Body.Close(); err != nil {
			logger.Error().Err(errors.Wrap(err, errs.ErrCloseBody)).Msg(errors.Wrap(err, errs.ErrCloseBody).Error())
		}
	}()
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return errors.Wrap(err, errs.ErrParseJSON)
	}
	return nil
}

func SendJSON(ctx context.Context, w http.ResponseWriter, data interface{}) error {
	logger := log.Ctx(ctx)

	w.Header().Set("Content-Type", "application/json")

	code := http.StatusOK
	w.WriteHeader(code)

	if data == nil {
		return nil
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrEncodeJSON)).Msg(errors.Wrap(err, errs.ErrEncodeJSON).Error())
		SendError(ctx, w, http.StatusInternalServerError, errors.Wrap(err, errs.ErrEncodeJSONShort).Error(), errors.Wrap(err, errs.ErrEncodeJSON).Error())
		return errors.Wrap(err, errs.ErrParseJSON)
	}
	return nil
}
