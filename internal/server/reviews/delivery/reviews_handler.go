package delivery

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/ds"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	authInterfaces "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/delivery/interfaces"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	movieInterfaces "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/delivery"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/reviews/delivery/dto"
	escapingutil "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/escaping_util"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	NEW_REVIEW_PLACEHOLDER_ID = -1
	REVIEWS_PER_PAGE          = 20
)

type ReviewHandler struct {
	userService    authInterfaces.UserServiceInterface
	sessionService authInterfaces.SessionServiceInterface
	movieService   movieInterfaces.MovieServiceInterface
}

func NewReviewHandler(userService authInterfaces.UserServiceInterface,
	sessionService authInterfaces.SessionServiceInterface, movieService movieInterfaces.MovieServiceInterface) *ReviewHandler {
	return &ReviewHandler{
		userService:    userService,
		sessionService: sessionService,
		movieService:   movieService,
	}
}

// GetPerson handles GET request to obtain person info by id
func (h *ReviewHandler) GetAllReviewsOfMovie(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	vars := mux.Vars(r)
	movieIDStr, ok := vars["movie_id"]
	if !ok {
		errMsg := errors.New("movie_id not found in path variables")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errs.ErrBadPayload, "Missing movie_id parameter")
		return
	}

	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		errMsg := errors.Wrapf(err, "getMovie action: bad request: %v", err)
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errs.ErrBadPayload, errs.ErrBadPayload)
		return
	}

	logger.Info().Msgf("getting all reviews of movie by it`s id: %d", movieID)
	reviewsJSON, err := h.movieService.GetAllReviewsOfMovieByID(r.Context(), movieID)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		if errors.Is(err, errs.ErrMovieNotFound) {
			jsonutil.SendError(r.Context(), w, http.StatusNotFound, errors.Wrap(err, errs.ErrNotFoundShort).Error(), err.Error())
			return
		}
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrSomethingWentWrong, errs.ErrSomethingWentWrong)
		return
	}
	logger.Info().Msgf("successfully got all movie reviews data by it`s id: %d", movieID)

	if err := jsonutil.SendJSON(r.Context(), w, reviewsJSON); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSomethingWentWrong).Error())
		return
	}
}

func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())
	var newReviewDataJSON *dto.NewReviewDataJSON

	vars := mux.Vars(r)
	movieIDStr, ok := vars["movie_id"]
	if !ok {
		errMsg := errors.New("movie_id not found in path variables")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errs.ErrBadPayload, "Missing movie_id parameter")
		return
	}
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		errMsg := errors.Wrapf(err, "getMovie action: bad request: %v", err)
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errs.ErrBadPayload, errs.ErrBadPayload)
		return
	}

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		logger.Warn().Msg(errors.Wrap(err, errs.ErrUnauthorized).Error())
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errs.ErrUnauthorizedShort,
			errs.ErrUnauthorized)
		return
	}

	username, errSession := h.sessionService.GetSession(r.Context(), sessionCookie.Value)
	if errSession != nil {
		logger.Error().Err(errors.Wrap(errSession, errs.ErrMsgSessionNotExists)).Msg(errs.ErrMsgFailedToGetSession)
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errs.ErrMsgSessionNotExists, errs.ErrMsgFailedToGetSession)
		return
	}

	if err = jsonutil.ReadJSON(r, &newReviewDataJSON); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrParseJSON)).Msg(errors.Wrap(err, errs.ErrParseJSON).Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errors.Wrap(err, errs.ErrParseJSONShort).Error(), errs.ErrBadPayload)
		return
	}

	if newReviewDataJSON.Score < 1 || newReviewDataJSON.Score > 10 {
		logger.Error().Err(errors.New(errs.ErrBadPayload)).Msg(fmt.Sprintf("bad score of new review: %d", newReviewDataJSON.Score))
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errs.ErrBadPayload,
			errs.ErrBadPayload)
		return
	}

	user, err := h.userService.GetUser(r.Context(), username)
	if err != nil {
		wrapped := errors.Wrap(err, "error getting user")
		logger.Error().Err(wrapped).Msg(wrapped.Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, wrapped.Error(), wrapped.Error())
		return
	}

	validatedReviewText := ""
	var errEscaping error
	if len(newReviewDataJSON.ReviewText) > 0 {
		validatedReviewText, errEscaping = escapingutil.ValidateInputTextData(newReviewDataJSON.ReviewText)
		if errEscaping != nil {
			logger.Error().Err(errEscaping).Msg(errEscaping.Error())
			jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errs.ErrBadPayload, errEscaping.Error())
			return
		}
	}

	newReview := mocks.ReviewJSON{
		ID:         NEW_REVIEW_PLACEHOLDER_ID,
		Score:      newReviewDataJSON.Score,
		ReviewText: validatedReviewText,
		CreatedAt:  "",
		User: mocks.ReviewUserDataJSON{
			Login:  user.Username,
			Avatar: user.Avatar,
		},
	}

	h.movieService.CreateNewMovieReview(r.Context(), movieID, newReview)

	if err = jsonutil.SendJSON(r.Context(), w, ds.Response{Message: "successfully created new review"}); err != nil {
		logger.Error().Err(err).Msg(errs.ErrSendJSON)
		return
	}
}
