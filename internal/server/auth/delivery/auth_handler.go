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
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/validation/auth"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	COOKIE_DAYS_LIMIT        = 3
	COOKIE_EXPIRED_LAST_YEAR = -1
)

type AuthServiceInterface interface {
	Register(ctx context.Context, regUser models.RegisterData) (string, error)
	GetSession(ctx context.Context, sessionID string) (string, error)
	DeleteSession(ctx context.Context, sessionID string) error
	Login(ctx context.Context, login models.LoginData) (string, error)
	Logout(ctx context.Context, sessionID string) error
}

type AuthHandler struct {
	authService AuthServiceInterface
	cookie      *config.Cookie
}

func NewAuthHandler(ctx context.Context, authService AuthServiceInterface) *AuthHandler {
	return &AuthHandler{
		cookie:      config.FromCookieContext(ctx),
		authService: authService,
	}
}

// Register http handler method
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
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

	newSessionID, errSession := h.authService.Register(r.Context(), reg)
	if errSession != nil {
		logger.Error().Err(errSession).Msgf("error happened: %v", errSession.Error)

		switch errSession.Error() {
		case errs.ErrInvalidPassword:
			jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errors.Wrap(errSession, errs.ErrInvalidPasswordShort).Error(),
				errors.Wrap(errSession, errs.ErrInvalidPassword).Error())
			return
		case errs.ErrAlreadyExists:
			logger.Error().Err(errors.New(errs.ErrAlreadyExists)).Msg(errs.ErrAlreadyExists)
			jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errors.New(errs.ErrAlreadyExistsShort).Error(),
				common.MsgUserWithNameAlreadyExists)
			return
		default:
			jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errors.New(errs.ErrSomethingWentWrong).Error(), errs.ErrSomethingWentWrong)
			return
		}
	}

	errOldSession := h.expireOldSessionCookie(w, r)
	if errOldSession != nil {
		logger.Error().Err(errOldSession).Msg(errOldSession.Error())
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrSomethingWentWrong,
			errs.ErrSomethingWentWrong)
		return
	}

	logger.Info().Msg("User registered successfully")

	http.SetCookie(w, preparedNewCookie(h.cookie, newSessionID))

	if err := jsonutil.SendJSON(r.Context(), w, ds.Response{Message: messages.SuccessfulRegister}); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSomethingWentWrong).Error())
		return
	}
}

// Expires old cookie if it exists
func (h *AuthHandler) expireOldSessionCookie(w http.ResponseWriter, r *http.Request) error {
	logger := log.Ctx(r.Context())

	oldSessionCookie, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		logger.Info().Msg("user dont have old cookie")
	}

	if oldSessionCookie != nil {
		http.SetCookie(w, preparedExpiredCookie(h.cookie))
		err := h.authService.DeleteSession(r.Context(), oldSessionCookie.Value)
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
		Expires:  time.Now().AddDate(0, 0, COOKIE_DAYS_LIMIT),
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
		Expires:  time.Now().AddDate(COOKIE_EXPIRED_LAST_YEAR, 0, 0),
	}
}

// Session http handler method
func (h *AuthHandler) Session(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	logger.Info().Msg("Checking session")

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		logger.Warn().Msg(errors.Wrap(err, errs.ErrUnauthorized).Error())
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errs.ErrUnauthorized,
			errs.ErrUnauthorized)
		return
	}

	username, errSession := h.authService.GetSession(r.Context(), sessionCookie.Value)
	if errSession != nil {
		logger.Error().Err(errors.Wrap(errSession, errs.ErrSessionNotExists)).Msg(errs.ErrMsgFailedToGetSession)
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errs.ErrSessionNotExists, errs.ErrMsgFailedToGetSession)
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
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	var login models.LoginData
	logger.Info().Msg("Logining user")

	// get user credentials from request body
	if err := jsonutil.ReadJSON(r, &login); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrParseJSON)).Msg(errors.Wrap(err, errs.ErrParseJSON).Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errs.ErrBadPayload,
			errors.Wrap(err, errs.ErrSomethingWentWrong).Error())
		return
	}

	newSessionID, errSession := h.authService.Login(r.Context(), login)
	if errSession != nil {
		switch errSession.Error() {
		case errs.ErrIncorrectLogin:
			logger.Error().Err(errors.Wrap(errSession, errs.ErrIncorrectLoginOrPassword)).Msg(errSession.Error())

			jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errors.Wrap(errSession, errs.ErrIncorrectLoginOrPasswordShort).Error(),
				errors.Wrap(errSession, errs.ErrIncorrectLoginOrPassword).Error())
			return
		case errs.ErrIncorrectPassword:
			logger.Error().Err(errors.Wrap(errSession, errs.ErrIncorrectLoginOrPassword)).Msg(errSession.Error())

			jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errors.Wrap(errSession, errs.ErrIncorrectLoginOrPasswordShort).Error(),
				errors.Wrap(errSession, errs.ErrIncorrectLoginOrPassword).Error())
			return
		default:
			logger.Error().Err(errSession).Msgf("error happened: %v", errSession.Error)
			jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrSomethingWentWrong, errs.ErrSomethingWentWrong)
			return
		}
	}

	// expire old cookie if it exists
	errOldSession := h.expireOldSessionCookie(w, r)
	if errOldSession != nil {
		logger.Error().Err(errOldSession).Msg(errOldSession.Error())
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrSomethingWentWrong,
			errs.ErrSomethingWentWrong)
		return
	}

	logger.Info().Msg("User logged in successfully")

	http.SetCookie(w, preparedNewCookie(h.cookie, newSessionID))

	err := jsonutil.SendJSON(r.Context(), w, ds.Response{Message: messages.SuccessfulLogin})
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSendJSON).Error())
		return
	}
}

// Logout http handler method
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	logger.Info().Msg("Logouting user")

	cookie, err := r.Cookie("session_id")
	if err != nil {
		logger.Warn().Msg(errors.Wrap(err, errs.ErrUnauthorized).Error())
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errors.Wrap(err, errs.ErrUnauthorizedShort).Error(),
			errs.ErrUnauthorized)
		return
	}

	sessionID := cookie.Value

	errSession := h.authService.Logout(r.Context(), sessionID)
	if errSession != nil {
		logger.Err(errSession).Msgf("error happened: %v", errSession)

		jsonutil.SendError(r.Context(), w, http.StatusNotFound, errors.Wrap(errSession, errs.ErrSessionNotExistsShort).Error(),
			errs.ErrSessionNotExists)
		return
	}

	http.SetCookie(w, preparedExpiredCookie(h.cookie))
	logger.Info().Msg("Session deleted")

	err = jsonutil.SendJSON(r.Context(), w, ds.Response{Message: messages.SuccessfulLogout})
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSendJSON).Error())
		return
	}
}
