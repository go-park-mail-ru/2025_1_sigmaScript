package delivery

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/ds"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/csat/delivery/dto"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	authInterfaces "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/delivery/interfaces"
	escapingutil "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/escaping_util"
)

const (
	NEW_REVIEW_PLACEHOLDER_ID = -1
	REVIEWS_PER_PAGE          = 20
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mocks.go -package=delivery_mocks MovieServiceInterface
type CSATServiceInterface interface {
	GetAllCSATReviews(ctx context.Context) (*[]dto.CSATReviewDataJSON, error)
	CreateNewCSATReview(ctx context.Context, newReview dto.CSATReviewDataJSON) error
	GetCSATStatistic(ctx context.Context) (*dto.CSATStatisticDataJSON, error)
}

type CSATHandler struct {
	userService    authInterfaces.UserServiceInterface
	sessionService authInterfaces.SessionServiceInterface
	CSATService    CSATServiceInterface
}

func NewCSATHandler(userService authInterfaces.UserServiceInterface,
	sessionService authInterfaces.SessionServiceInterface, csatService CSATServiceInterface) *CSATHandler {
	return &CSATHandler{
		userService:    userService,
		sessionService: sessionService,
		CSATService:    csatService,
	}
}

func (h *CSATHandler) GetAllCSATReviews(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	logger.Info().Msgf("getting all CSAT reviews")

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		logger.Warn().Msg(errors.Wrap(err, errs.ErrUnauthorized).Error())
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errs.ErrUnauthorizedShort,
			errs.ErrUnauthorized)
		return
	}

	_, errSession := h.sessionService.GetSession(r.Context(), sessionCookie.Value)
	if errSession != nil {
		logger.Error().Err(errors.Wrap(errSession, errs.ErrMsgSessionNotExists)).Msg(errs.ErrMsgFailedToGetSession)
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errs.ErrMsgSessionNotExists, errs.ErrMsgFailedToGetSession)
		return
	}

	reviewsJSON, err := h.CSATService.GetAllCSATReviews(r.Context())
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		if errors.Is(err, errs.ErrMovieNotFound) {
			jsonutil.SendError(r.Context(), w, http.StatusNotFound, errors.Wrap(err, errs.ErrNotFoundShort).Error(), err.Error())
			return
		}
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrSomethingWentWrong, errs.ErrSomethingWentWrong)
		return
	}
	logger.Info().Msgf("successfully got all CSAT reviews")

	if err := jsonutil.SendJSON(r.Context(), w, reviewsJSON); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSomethingWentWrong).Error())
		return
	}
}

func (h *CSATHandler) CreateCSATReview(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())
	var newReviewDataJSON *dto.NewCSATReviewDataJSON

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

	newReview := dto.CSATReviewDataJSON{
		ID:        NEW_REVIEW_PLACEHOLDER_ID,
		Score:     newReviewDataJSON.Score,
		CSATText:  validatedReviewText,
		CreatedAt: "",
		User: dto.ReviewUserDataJSON{
			Login:  user.Username,
			Avatar: user.Avatar,
		},
	}

	errCreate := h.CSATService.CreateNewCSATReview(r.Context(), newReview)
	if errCreate != nil {
		wrapped := errors.Wrap(errCreate, "error creating review")
		logger.Error().Err(wrapped).Msg(wrapped.Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, wrapped.Error(), wrapped.Error())
		return
	}

	if err = jsonutil.SendJSON(r.Context(), w, ds.Response{Message: "successfully created new review"}); err != nil {
		logger.Error().Err(err).Msg(errs.ErrSendJSON)
		return
	}
}

func (h *CSATHandler) GetCSATStatistic(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	logger.Info().Msgf("getting CSAT statistic")

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		logger.Warn().Msg(errors.Wrap(err, errs.ErrUnauthorized).Error())
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errs.ErrUnauthorizedShort,
			errs.ErrUnauthorized)
		return
	}

	_, errSession := h.sessionService.GetSession(r.Context(), sessionCookie.Value)
	if errSession != nil {
		logger.Error().Err(errors.Wrap(errSession, errs.ErrMsgSessionNotExists)).Msg(errs.ErrMsgFailedToGetSession)
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errs.ErrMsgSessionNotExists, errs.ErrMsgFailedToGetSession)
		return
	}

	statisticJSON, err := h.CSATService.GetCSATStatistic(r.Context())
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		if errors.Is(err, errs.ErrCSATReviewsNotFound) {
			jsonutil.SendError(r.Context(), w, http.StatusNotFound, errors.Wrap(err, errs.ErrNotFoundShort).Error(), err.Error())
			return
		}
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrSomethingWentWrong, errs.ErrSomethingWentWrong)
		return
	}

	logger.Info().Msgf("successfully got CSAT statistic")

	if err := jsonutil.SendJSON(r.Context(), w, statisticJSON); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSomethingWentWrong).Error())
		return
	}
}
