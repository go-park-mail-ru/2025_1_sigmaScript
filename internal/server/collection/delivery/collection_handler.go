package delivery

import (
	"context"
	"net/http"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=$GOFILE -destination=delivery_mocks/mocks.go -package=delivery_mocks CollectionServiceInterface
type CollectionServiceInterface interface {
	GetMainPageCollections(ctx context.Context) (mocks.Collections, error)
}

type CollectionHandler struct {
	collectionService CollectionServiceInterface
}

func NewCollectionHandler(collectionService CollectionServiceInterface) *CollectionHandler {
	return &CollectionHandler{
		collectionService: collectionService,
	}
}

func (h *CollectionHandler) GetMainPageCollections(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	logger.Info().Msg("GetCollections")

	collections, err := h.collectionService.GetMainPageCollections(r.Context())
	if err != nil {
		logger.Err(err).Msg(err.Error())
		jsonutil.SendError(r.Context(), w, http.StatusNotFound, errors.Wrap(err, errs.ErrCollectionNotExist.Error()).Error(),
			errs.ErrCollectionNotExist.Error())
		return
	}

	if err := jsonutil.SendJSON(r.Context(), w, collections); err != nil {
		logger.Error().Err(err).Msg("Error sending JSON")
		return
	}
}
