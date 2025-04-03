package handlers

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/rs/zerolog/log"
)

func GetCollections(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	logger.Info().Msg("GetCollections")

	collections := mocks.MainPageCollections

	if err := jsonutil.SendJSON(w, collections); err != nil {
		logger.Error().Err(err).Msg("Error sending JSON")
		return
	}
}
