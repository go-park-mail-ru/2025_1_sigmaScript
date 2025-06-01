package delivery

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mocks.go -package=delivery_mocks GenreServiceInterface
type GenreServiceInterface interface {
	GetGenreByID(ctx context.Context, genreID string) (*mocks.Genre, error)
	GetAllGenres(ctx context.Context) (*[]mocks.Genre, error)
}

type GenreHandler struct {
	genreService GenreServiceInterface
}

func NewGenreHandler(movieService GenreServiceInterface) *GenreHandler {
	return &GenreHandler{
		genreService: movieService,
	}
}

func (h *GenreHandler) GetGenreByID(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	genreIDStr, ok := mux.Vars(r)["genre_id"]
	if !ok {
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errs.ErrBadPayload, "Missing genre_id parameter")
		return
	}

	logger.Info().Msgf("getting genre by id: %s", genreIDStr)
	genreJSON, err := h.genreService.GetGenreByID(r.Context(), genreIDStr)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		if errors.Is(err, errs.ErrMovieNotFound) {
			jsonutil.SendError(r.Context(), w, http.StatusNotFound, errors.Wrap(err, errs.ErrNotFoundShort).Error(), err.Error())
			return
		}
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrSomethingWentWrong, errs.ErrSomethingWentWrong)
		return
	}
	logger.Info().Msgf("successfully got genre data by id: %s", genreIDStr)

	if err := jsonutil.SendJSON(r.Context(), w, genreJSON); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSomethingWentWrong).Error())
		return
	}
}

func (h *GenreHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	logger.Info().Msgf("getting all genres")
	genreJSON, err := h.genreService.GetAllGenres(r.Context())
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		if errors.Is(err, errs.ErrMovieNotFound) {
			jsonutil.SendError(r.Context(), w, http.StatusNotFound, errors.Wrap(err, errs.ErrNotFoundShort).Error(), err.Error())
			return
		}
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrSomethingWentWrong, errs.ErrSomethingWentWrong)
		return
	}

	logger.Info().Msgf("successfully got all genres")

	if err := jsonutil.SendJSON(r.Context(), w, genreJSON); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSomethingWentWrong).Error())
		return
	}
}
