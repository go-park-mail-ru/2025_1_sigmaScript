package delivery

import (
	"context"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/ds"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/messages"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/service"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	COOKIE_DAYS_LIMIT        = 3
	COOKIE_EXPIRED_LAST_YEAR = -1
)

type AuthHandlerInterface interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	Session(w http.ResponseWriter, r *http.Request)
}

type AuthHandler struct {
	authService service.AuthServiceInterface
	cfg         *config.Cookie
}

func NewAuthHandler(ctx context.Context, authService service.AuthServiceInterface) AuthHandlerInterface {
	res := &AuthHandler{
		cfg:         config.FromCookieContext(ctx),
		authService: authService,
	}

	return res
}

// Register http handler method
func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	var reg models.RegisterData
	logger.Info().Msg("Registering user")

	if err := jsonutil.ReadJSON(r, &reg); err != nil {
		msg := errs.ErrBadPayload
		logger.Error().Err(errors.Wrap(err, errs.ErrParseJSON)).Msg(errors.Wrap(err, errs.ErrParseJSON).Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errors.Wrap(err, errs.ErrParseJSONShort).Error(), msg)
		return
	}

	newSessionID, errSession := a.authService.Register(r.Context(), reg)
	if errSession != nil {
		logger.Error().Err(errSession.Error).Msgf("error happened %v, with code %d", errSession.Error, errSession.Code)
		jsonutil.SendError(r.Context(), w, errSession.Code, errors.New(errs.ErrAlreadyExistsShort).Error(), errs.ErrSomethingWentWrong)
		return
	}

	errOldSession := a.expireOldSessionCookie(w, r)
	if errOldSession != nil {
		logger.Error().Err(errOldSession.Error).Msg(errOldSession.Error.Error())
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrSomethingWentWrong,
			errs.ErrSomethingWentWrong)
		return
	}

	logger.Info().Msg("User registered successfully")

	http.SetCookie(w, preparedNewCookie(a.cfg, newSessionID))

	if err := jsonutil.SendJSON(r.Context(), w, ds.Response{Message: messages.SuccessfulRegister}); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSomethingWentWrong).Error())
		return
	}
}

// Expires old cookie if it exists
func (a *AuthHandler) expireOldSessionCookie(w http.ResponseWriter, r *http.Request) *errs.ServiceError {
	logger := log.Ctx(r.Context())

	oldSessionCookie, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		logger.Info().Msg("user dont have old cookie")
	}

	if oldSessionCookie != nil {
		http.SetCookie(w, preparedExpiredCookie(a.cfg))
		err := a.authService.DeleteSession(r.Context(), oldSessionCookie.Value)
		if err != nil {
			return err
		}
		logger.Info().Msg("successfully expired old sesssion cookie")
	}

	return nil
}

func preparedNewCookie(cfg *config.Cookie, newSessionID string) *http.Cookie {
	return &http.Cookie{
		Name:     cfg.SessionName,
		Value:    newSessionID,
		HttpOnly: cfg.HTTPOnly,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
		Path:     cfg.Path,
		Expires:  time.Now().AddDate(0, 0, COOKIE_DAYS_LIMIT),
	}
}

func preparedExpiredCookie(cfg *config.Cookie) *http.Cookie {
	return &http.Cookie{
		Name:     cfg.SessionName,
		Value:    "",
		HttpOnly: cfg.HTTPOnly,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
		Path:     cfg.Path,
		Expires:  time.Now().AddDate(COOKIE_EXPIRED_LAST_YEAR, 0, 0),
	}
}

// Session http handler method
func (a *AuthHandler) Session(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	logger.Info().Msg("Checking session")

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		logger.Warn().Msg(errors.Wrap(err, errs.ErrUnauthorized).Error())
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errs.ErrUnauthorized,
			errs.ErrUnauthorized)
		return
	}

	username, errSession := a.authService.GetSession(r.Context(), sessionCookie.Value)
	if errSession != nil {
		errMsg := "failed to get session"
		logger.Error().Err(errors.Wrap(errSession.Error, errs.ErrSessionNotExists)).Msg(errMsg)
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errs.ErrSessionNotExists, errMsg)
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
func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
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

	newSessionID, errSession := a.authService.Login(r.Context(), login)
	if errSession != nil {
		logger.Error().Err(errSession.Error).Msgf("error happened %v, with code %d", errSession.Error, errSession.Code)
		jsonutil.SendError(r.Context(), w, errSession.Code, errSession.Error.Error(), errs.ErrSomethingWentWrong)
		return
	}

	// expire old cookie if it exists
	errOldSession := a.expireOldSessionCookie(w, r)
	if errOldSession != nil {
		logger.Error().Err(errOldSession.Error).Msg(errOldSession.Error.Error())
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrSomethingWentWrong,
			errs.ErrSomethingWentWrong)
		return
	}

	logger.Info().Msg("User logged in successfully")

	http.SetCookie(w, preparedNewCookie(a.cfg, newSessionID))

	err := jsonutil.SendJSON(r.Context(), w, ds.Response{Message: messages.SuccessfulLogin})
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSendJSON).Error())
		return
	}
}

// Logout http handler method
func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	logger.Info().Msg("Logouting user")

	cookie, err := r.Cookie("session_id")
	if err != nil {
		logger.Warn().Msg(errors.Wrap(err, errs.ErrUnauthorized).Error())
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errs.ErrUnauthorized,
			errs.ErrUnauthorized)
		return
	}

	sessionID := cookie.Value

	errSession := a.authService.Logout(r.Context(), sessionID)
	if errSession != nil {
		logger.Error().Err(errSession.Error).Msgf("error happened %v, with code %d", errSession.Error, errSession.Code)
		jsonutil.SendError(r.Context(), w, errSession.Code, errs.ErrSessionNotExists, errs.ErrSomethingWentWrong)
		return
	}

	http.SetCookie(w, preparedExpiredCookie(a.cfg))
	logger.Info().Msg("Session deleted")

	err = jsonutil.SendJSON(r.Context(), w, ds.Response{Message: messages.SuccessfulLogout})
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSendJSON).Error())
		return
	}
}
