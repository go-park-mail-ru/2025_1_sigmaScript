package handlers

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/rs/zerolog/log"
)

func GetCollections(w http.ResponseWriter, r *http.Request) {
	log.Info().Msg("GetCollections")

	collection := mocks.Collections

	if err := jsonutil.SendJSON(w, collection); err != nil {
		log.Error().Err(err).Msg("Error sending JSON")
		return
	}
}
