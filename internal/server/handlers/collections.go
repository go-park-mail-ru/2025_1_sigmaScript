package handlers

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func GetCollection(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	log.Info().Msg("GetCollection")

	id := vars["id"]

	collection, exists := mocks.Collections[id]
	if !exists {
		log.Info().Msg("Collection not found")
		jsonutil.SendError(w, http.StatusNotFound, "not_found", "Collection not found")
		return
	}

	if err := jsonutil.SendJSON(w, collection); err != nil {
		log.Error().Err(err).Msg("Error sending JSON")
		return
	}
}
