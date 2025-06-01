package delivery

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	escapingutil "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/escaping_util"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mocks.go -package=delivery_mocks SearchServiceInterface
type SearchServiceInterface interface {
	SearchActorsAndMovies(ctx context.Context, searchStr string) (*models.SearchResponseJSON, error)
}

type SearchHandler struct {
	searchService SearchServiceInterface
}

func NewSearchHandler(searchService SearchServiceInterface) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

func (h *SearchHandler) SearchActorsAndMovies(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	var searchReq *models.SearchRequestJSON
	if err := jsonutil.ReadJSON(r, &searchReq); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrParseJSON)).Msg(errors.Wrap(err, errs.ErrParseJSON).Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errors.Wrap(err, errs.ErrParseJSONShort).Error(), errs.ErrBadPayload)
		return
	}

	validatedReviewText, errEscaping := escapingutil.ValidateInputTextData(searchReq.SearchString)
	if errEscaping != nil {
		logger.Error().Err(errEscaping).Msg(errEscaping.Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errs.ErrBadPayload, errEscaping.Error())
		return
	}

	logger.Info().Msgf("getting search results")
	searchResJSON, err := h.searchService.SearchActorsAndMovies(r.Context(), validatedReviewText)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		if errors.Is(err, errs.ErrMovieNotFound) {
			jsonutil.SendError(r.Context(), w, http.StatusNotFound, errors.Wrap(err, errs.ErrNotFoundShort).Error(), err.Error())
			return
		}
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrSomethingWentWrong, errs.ErrSomethingWentWrong)
		return
	}

	logger.Info().Msgf("successfully got search results")

	if err := jsonutil.SendJSON(r.Context(), w, searchResJSON); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSomethingWentWrong).Error())
		return
	}
}
