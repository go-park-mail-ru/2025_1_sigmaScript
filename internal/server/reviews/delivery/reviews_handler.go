package delivery

import (
	"context"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/common"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/ds"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/messages"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/validation/auth"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type ReviewHandlerInterface interface {
	GetReview(w http.ResponseWriter, r *http.Request)
	GetReviewsOfMovie(w http.ResponseWriter, r *http.Request)
	CreateReview(w http.ResponseWriter, r *http.Request)
	UpdateReview(w http.ResponseWriter, r *http.Request)
	DeleteReview(w http.ResponseWriter, r *http.Request)
}

/////////////////////////////////////////////////////////////////////////////////

type UserServiceInterface interface {
	Register(ctx context.Context, regUser models.RegisterData) error
	Login(ctx context.Context, login models.LoginData) error
}

type SessionServiceInterface interface {
	GetSession(ctx context.Context, sessionID string) (string, error)
	DeleteSession(ctx context.Context, sessionID string) error
	CreateSession(ctx context.Context, username string) (string, error)
}

type ReviewServiceInterface interface {
	GetReview(ctx context.Context, movieID, userID int) (mocks.ReviewJSON, error)
	GetReviewsOfMovie(ctx context.Context, movieID int, paginatorPageNumber ...int) ([]mocks.ReviewJSON, error)
	CreateReview(ctx context.Context, newReview mocks.ReviewJSON) error
	UpdateReview(ctx context.Context, updatedReview mocks.ReviewJSON) error
	DeleteReview(ctx context.Context, reviewID int) error
}

type ReviewHandler struct {
	userService    UserServiceInterface
	sessionService SessionServiceInterface
	reviewService  ReviewServiceInterface
	cookieData     *config.Cookie
}

func NewReviewHandler(ctx context.Context, userService UserServiceInterface,
	sessionService SessionServiceInterface, reviewService ReviewServiceInterface) *ReviewHandler {
	return &ReviewHandler{
		cookieData:     config.FromCookieContext(ctx),
		userService:    userService,
		sessionService: sessionService,
		reviewService:  reviewService,
	}
}

// GetReview http handler method
func (h *ReviewHandler) Register(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	var reg models.RegisterData
	logger.Info().Msg("Registering user")

	if err := jsonutil.ReadJSON(r, &reg); err != nil {
		msg := errs.ErrBadPayload
		logger.Error().Err(errors.Wrap(err, errs.ErrParseJSON)).Msg(errors.Wrap(err, errs.ErrParseJSON).Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errors.Wrap(err, errs.ErrParseJSONShort).Error(), msg)
		return
	}

	if reg.Password != reg.RepeatedPassword {
		logger.Info().Msg("Passwords mismatch")
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errors.New(errs.ErrPasswordsMismatchShort).Error(), errs.ErrPasswordsMismatch)
		return
	}

	if err := auth.IsValidPassword(reg.Password); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrInvalidPassword)).Msg(errors.Wrap(err, errs.ErrInvalidPassword).Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errors.Wrap(err, errs.ErrInvalidPasswordShort).Error(),
			errors.Wrap(err, errs.ErrInvalidPassword).Error())
		return
	}

	if err := auth.IsValidLogin(reg.Username); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrInvalidLogin)).Msg(errors.Wrap(err, errs.ErrInvalidLogin).Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errors.Wrap(err, errs.ErrInvalidLoginShort).Error(),
			errors.Wrap(err, errs.ErrInvalidLogin).Error())
		return
	}

	err := h.userService.Register(r.Context(), reg)
	if err != nil {
		logger.Error().Err(err).Msgf("error happened: %v", err.Error)

		switch err.Error() {
		case errs.ErrInvalidPassword:
			jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errors.Wrap(err, errs.ErrInvalidPasswordShort).Error(),
				errors.Wrap(err, errs.ErrInvalidPassword).Error())
			return
		case errs.ErrAlreadyExists:
			jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errors.New(errs.ErrAlreadyExistsShort).Error(),
				common.MsgUserWithNameAlreadyExists)
			return
		default:
			jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errors.New(errs.ErrSomethingWentWrong).Error(), errs.ErrSomethingWentWrong)
			return
		}
	}
	logger.Info().Msg("User registered successfully")

	// expire old session cookie if it exists
	errOldSession := h.expireOldSessionCookie(w, r)
	if errOldSession != nil {
		logger.Warn().Err(errOldSession).Msg(errOldSession.Error())
	}

	newSessionID, err := h.sessionService.CreateSession(r.Context(), reg.Username)
	if err != nil {
		logger.Error().Err(err).Msgf("error happened: %v", err.Error)

		if errors.Is(err, errs.ErrGenerateSession) {
			jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrMsgGenerateSessionShort,
				errs.ErrMsgGenerateSession)
			return
		}
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrSomethingWentWrong,
			errs.ErrSomethingWentWrong)
		return
	}

	http.SetCookie(w, preparedNewCookie(h.cookieData, newSessionID))

	if err := jsonutil.SendJSON(r.Context(), w, ds.Response{Message: messages.SuccessfulRegister}); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSomethingWentWrong).Error())
		return
	}
}

// expires old session cookie if it exists
func (h *ReviewHandler) expireOldSessionCookie(w http.ResponseWriter, r *http.Request) error {
	logger := log.Ctx(r.Context())

	oldSessionCookie, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		logger.Info().Msg("user dont have old cookie")
		return nil
	}

	if oldSessionCookie != nil {
		http.SetCookie(w, preparedExpiredCookie(h.cookieData))
		err := h.sessionService.DeleteSession(r.Context(), oldSessionCookie.Value)
		if err != nil {
			return err
		}
		logger.Info().Msg("successfully expired old sesssion cookie")
	}

	return nil
}

func preparedNewCookie(cookie *config.Cookie, newSessionID string) *http.Cookie {
	return &http.Cookie{
		Name:     cookie.SessionName,
		Value:    newSessionID,
		HttpOnly: cookie.HTTPOnly,
		Secure:   cookie.Secure,
		SameSite: cookie.SameSite,
		Path:     cookie.Path,
		Expires:  time.Now().AddDate(0, 0, common.COOKIE_DAYS_LIMIT),
	}
}

func preparedExpiredCookie(cookie *config.Cookie) *http.Cookie {
	return &http.Cookie{
		Name:     cookie.SessionName,
		Value:    "",
		HttpOnly: cookie.HTTPOnly,
		Secure:   cookie.Secure,
		SameSite: cookie.SameSite,
		Path:     cookie.Path,
		Expires:  time.Now().AddDate(common.COOKIE_EXPIRED_LAST_YEAR, 0, 0),
	}
}

// Session http handler method gets user data by session
func (h *ReviewHandler) Session(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	logger.Info().Msg("Checking session")
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
	logger.Info().Interface("session username", username).Msg("getSession success")

	err = jsonutil.SendJSON(r.Context(), w, ds.User{Username: username})
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSendJSON).Error())
		return
	}
}

// Login http handler method
func (h *ReviewHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	var login models.LoginData
	logger.Info().Msg("Logining user")

	// get user credentials from request body
	err := jsonutil.ReadJSON(r, &login)
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrParseJSON)).Msg(errors.Wrap(err, errs.ErrParseJSON).Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errs.ErrBadPayload,
			errors.Wrap(err, errs.ErrSomethingWentWrong).Error())
		return
	}

	err = h.userService.Login(r.Context(), login)
	if err != nil {
		switch err.Error() {
		case errs.ErrIncorrectLogin:
			logger.Error().Err(errors.Wrap(err, errs.ErrIncorrectLoginOrPassword)).Msg(err.Error())
			jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errors.Wrap(err, errs.ErrIncorrectLoginOrPasswordShort).Error(),
				errors.Wrap(err, errs.ErrIncorrectLoginOrPassword).Error())
			return
		case errs.ErrIncorrectPassword:
			logger.Error().Err(errors.Wrap(err, errs.ErrIncorrectLoginOrPassword)).Msg(err.Error())
			jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errors.Wrap(err, errs.ErrIncorrectLoginOrPasswordShort).Error(),
				errors.Wrap(err, errs.ErrIncorrectLoginOrPassword).Error())
			return
		default:
			logger.Error().Err(err).Msgf("error happened: %v", err.Error)
			jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrSomethingWentWrong, errs.ErrSomethingWentWrong)
			return
		}
	}
	logger.Info().Msg("User logged in successfully")

	// expire old session cookie if it exists
	errOldSession := h.expireOldSessionCookie(w, r)
	if errOldSession != nil {
		logger.Warn().Err(errOldSession).Msg(errOldSession.Error())
	}

	newSessionID, err := h.sessionService.CreateSession(r.Context(), login.Username)
	if err != nil {
		logger.Error().Err(err).Msgf("error happened: %v", err.Error)

		if errors.Is(err, errs.ErrGenerateSession) {
			jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrMsgGenerateSessionShort,
				errs.ErrMsgGenerateSession)
			return
		}
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrSomethingWentWrong,
			errs.ErrSomethingWentWrong)
		return
	}

	http.SetCookie(w, preparedNewCookie(h.cookieData, newSessionID))

	err = jsonutil.SendJSON(r.Context(), w, ds.Response{Message: messages.SuccessfulLogin})
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSendJSON).Error())
		return
	}
}

// Logout http handler method
func (h *ReviewHandler) Logout(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	logger.Info().Msg("Logouting user")
	cookie, err := r.Cookie("session_id")
	if err != nil {
		logger.Warn().Msg(errors.Wrap(err, errs.ErrUnauthorized).Error())
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errors.Wrap(err, errs.ErrUnauthorizedShort).Error(),
			errs.ErrUnauthorized)
		return
	}

	errSession := h.sessionService.DeleteSession(r.Context(), cookie.Value)
	if errSession != nil {
		logger.Err(errSession).Msgf("error happened: %v", errSession)
		jsonutil.SendError(r.Context(), w, http.StatusNotFound, errors.Wrap(errSession, errs.ErrMsgSessionNotExistsShort).Error(),
			errs.ErrMsgSessionNotExists)
		return
	}

	http.SetCookie(w, preparedExpiredCookie(h.cookieData))
	logger.Info().Msg("Session deleted")

	err = jsonutil.SendJSON(r.Context(), w, ds.Response{Message: messages.SuccessfulLogout})
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSendJSON).Error())
		return
	}
}
