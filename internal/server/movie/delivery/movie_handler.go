package delivery

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mocks.go -package=delivery_mocks MovieServiceInterface
type MovieServiceInterface interface {
	GetMovieByID(ctx context.Context, movieID int) (*mocks.MovieJSON, error)
	GetAllReviewsOfMovieByID(ctx context.Context, movieID int) (*[]mocks.ReviewJSON, error)
	CreateNewMovieReview(ctx context.Context,
		userID string,
		movieID string,
		newReview mocks.NewReviewDataJSON) (*mocks.NewReviewDataJSON, error)
	UpdateMovieReview(ctx context.Context, userID string, movieID string, newReview mocks.NewReviewDataJSON) (*mocks.NewReviewDataJSON, error)
}

type MovieHandler struct {
	movieService MovieServiceInterface
}

func NewMovieHandler(movieService MovieServiceInterface) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
	}
}

func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	movieIDStr, ok := mux.Vars(r)["movie_id"]
	if !ok {
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errs.ErrBadPayload, "Missing movie_id parameter")
		return
	}

	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		errMsg := errors.Wrapf(err, "getMovieByID action: bad request")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errs.ErrBadPayload, errs.ErrBadPayload)
		return
	}

	logger.Info().Msgf("getting movie by id: %d", movieID)
	movieJSON, err := h.movieService.GetMovieByID(r.Context(), movieID)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		if errors.Is(err, errs.ErrMovieNotFound) {
			jsonutil.SendError(r.Context(), w, http.StatusNotFound, errors.Wrap(err, errs.ErrNotFoundShort).Error(), err.Error())
			return
		}
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrSomethingWentWrong, errs.ErrSomethingWentWrong)
		return
	}
	logger.Info().Msgf("successfully got movie data by id: %d", movieID)

	if err := jsonutil.SendJSON(r.Context(), w, movieJSON); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSomethingWentWrong).Error())
		return
	}
}
