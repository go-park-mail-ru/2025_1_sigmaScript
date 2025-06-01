package jsonutil

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	easyjson "github.com/mailru/easyjson"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func SendResponse(ctx context.Context, w http.ResponseWriter, code int, data interface{}) {
	if data == nil {
		w.WriteHeader(code)
		return
	}

	logger := log.Ctx(ctx)

	marshaler, ok := data.(easyjson.Marshaler)
	if !ok {
		logger.Error().Msg("object does not implement easyjson.Marshaler")

		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(data); err != nil {
			logger.Error().Err(err).Msg("error while encoding response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(code)
		_, wErr := w.Write(buf.Bytes())
		if wErr != nil {
			logger.Error().Err(wErr).Msg("error while writing response")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}

	var buf bytes.Buffer
	if _, err := easyjson.MarshalToWriter(marshaler, &buf); err != nil {
		logger.Error().Err(err).Msg("error while encoding success response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	_, writeErr := w.Write(buf.Bytes())
	if writeErr != nil {
		logger.Error().Err(writeErr).Msg("error while writing response to client")
	}
}

func DecodeBody(w http.ResponseWriter, r *http.Request, data interface{}) bool {
	logger := log.Ctx(r.Context())

	unmarshaler, ok := data.(easyjson.Unmarshaler)
	if !ok {
		logger.Error().Msg("object does not implement easyjson.Unmarshaler")
		if err := json.NewDecoder(r.Body).Decode(data); err != nil {
			logger.Error().Err(err).Msg("cannot parse request")
			SendResponse(r.Context(), w, http.StatusBadRequest, fmt.Errorf("cannot parse request: %w", err))
		}
		return false
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error().Err(err).Msg("cannot read request body")
		SendResponse(r.Context(), w, http.StatusInternalServerError, fmt.Errorf("cannot read request body: %w", err))
		return false
	}

	if err := easyjson.Unmarshal(body, unmarshaler); err != nil {
		logger.Error().Err(err).Msg("cannot parse request")
		SendResponse(r.Context(), w, http.StatusBadRequest, fmt.Errorf("cannot parse request: %w", err))
		return false
	}

	return true
}

func SendRequestError(ctx context.Context, w http.ResponseWriter, code string, status int, err error) {
	logger := log.Ctx(ctx)
	logger.Error().Err(err).Msg("request error")

	SendResponse(ctx, w, status, errors.Wrap(err, errs.ErrParseJSON))
}
