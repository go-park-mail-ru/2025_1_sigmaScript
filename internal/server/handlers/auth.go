package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/ds"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/messages"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/validation/auth"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/session"
	synccredmap "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/sync_cred_map"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
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
	// username --> hashedPass
	users synccredmap.SyncCredentialsMap
	// sessionID --> username
	sessions synccredmap.SyncCredentialsMap
	cfg      *config.Cookie
}

func NewAuthHandler(ctx context.Context) AuthHandlerInterface {
	res := &AuthHandler{
		users:    *synccredmap.NewSyncCredentialsMap(),
		sessions: *synccredmap.NewSyncCredentialsMap(),
		cfg:      config.FromCookieContext(ctx),
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
	fmt.Println(reg)

	if _, exists := a.users.Load(reg.Username); exists {
		msg := "user with that name already exists"
		logger.Error().Err(errors.New(errs.ErrAlreadyExists)).Msg(errs.ErrAlreadyExists)
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errors.New(errs.ErrAlreadyExistsShort).Error(), msg)
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

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrBcrypt)).Msg(errors.Wrap(err, errs.ErrBcrypt).Error())
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errors.Wrap(err, errs.ErrInvalidPasswordShort).Error(),
			errors.Wrap(err, errs.ErrInvalidPassword).Error())
		return
	}

	a.users.Store(reg.Username, string(hashedPass))

	logger.Info().Msg("User registered successfully")

	newSessionID, errSession := a.createNewSessionWithCookie(w, r)
	if errSession != nil {
		return
	}

	a.sessions.Store(newSessionID, reg.Username)
	logger.Info().Msg("Session created")

	http.SetCookie(w, preparedNewCookie(a.cfg, newSessionID))

	if err = jsonutil.SendJSON(r.Context(), w, ds.Response{Message: messages.SuccessfulRegister}); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSomethingWentWrong).Error())
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
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errors.Wrap(err, errs.ErrParseJSONShort).Error(),
			errors.Wrap(err, errs.ErrSomethingWentWrong).Error())
		return
	}

	// check if user with that name and pass exists
	hashedPass, exists := a.users.Load(login.Username)
	if !exists {
		errMsg := errors.New(errs.ErrIncorrectLogin)

		logger.Error().Err(errors.Wrap(errMsg, errs.ErrIncorrectLoginOrPassword)).Msg(errMsg.Error())

		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errors.Wrap(errMsg, errs.ErrIncorrectLoginOrPasswordShort).Error(),
			errors.Wrap(errMsg, errs.ErrIncorrectLoginOrPassword).Error())
		return
	} else if err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(login.Password)); err != nil {
		errMsg := errors.New(errs.ErrIncorrectPassword)

		logger.Error().Err(errors.Wrap(err, errs.ErrIncorrectLoginOrPassword)).Msg(errMsg.Error())

		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errors.Wrap(errMsg, errs.ErrIncorrectLoginOrPasswordShort).Error(),
			errors.Wrap(errMsg, errs.ErrIncorrectLoginOrPassword).Error())
		return
	}

	newSessionID, errSession := a.createNewSessionWithCookie(w, r)
	if errSession != nil {
		return
	}

	a.sessions.Store(newSessionID, login.Username)
	logger.Info().Msg("Session created")

	http.SetCookie(w, preparedNewCookie(a.cfg, newSessionID))

	err := jsonutil.SendJSON(r.Context(), w, ds.Response{Message: messages.SuccessfulLogin})
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSendJSON).Error())
		return
	}
}

func (a *AuthHandler) createNewSessionWithCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	logger := log.Ctx(r.Context())

	oldSessionCookie, err := r.Cookie("session_id")
	if errors.Is(err, http.ErrNoCookie) {
		logger.Info().Msg("user dont have old cookie")
	}

	// create new session for user
	newSessionID, err := session.GenerateSessionID(a.cfg.SessionLength)
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrGenerateSession)).Msg(errors.Wrap(err, errs.ErrGenerateSession).Error())
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errors.Wrap(err, errs.ErrGenerateSessionShort).Error(),
			errors.Wrap(err, errs.ErrGenerateSession).Error())
		return "", nil
	}

	if oldSessionCookie != nil {
		http.SetCookie(w, preparedExpiredCookie(a.cfg))
		a.sessions.Delete(oldSessionCookie.Value)
		logger.Info().Msg("successfully expire old sesssion cookie")
	}

	return newSessionID, nil
}

// Logout http handler method
func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	logger.Info().Msg("Logouting user")

	cookie, err := r.Cookie("session_id")
	if err != nil {
		logger.Warn().Msg(errors.Wrap(err, errs.ErrUnauthorized).Error())
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errors.Wrap(err, errs.ErrUnauthorizedShort).Error(),
			errors.Wrap(err, errs.ErrUnauthorized).Error())
		return
	}

	sessionID := cookie.Value
	if _, exists := a.sessions.Load(sessionID); !exists {
		err := errors.New("session does not exist")
		logger.Warn().Msg(errs.ErrSessionNotExists)
		jsonutil.SendError(r.Context(), w, http.StatusNotFound, errors.Wrap(err, errs.ErrSessionNotExistsShort).Error(), errs.ErrSessionNotExists)
		return
	}

	a.sessions.Delete(sessionID)
	http.SetCookie(w, preparedExpiredCookie(a.cfg))
	logger.Info().Msg("Session deleted")

	err = jsonutil.SendJSON(r.Context(), w, ds.Response{Message: messages.SuccessfulLogout})
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSendJSON).Error())
		return
	}
}

// Session http handler method
func (a *AuthHandler) Session(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	logger.Info().Msg("Checking session")

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		logger.Warn().Msg(errors.Wrap(err, errs.ErrUnauthorized).Error())
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errors.Wrap(err, errs.ErrUnauthorizedShort).Error(),
			errors.Wrap(err, errs.ErrUnauthorized).Error())
		return
	}

	username, ok := a.sessions.Load(sessionCookie.Value)
	if !ok {
		err := errors.New("failed to get session")
		logger.Error().Err(errors.Wrap(err, errs.ErrSessionNotExists)).Msg(errors.Wrap(err, errs.ErrSessionNotExists).Error())
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errors.Wrap(err, errs.ErrSessionNotExists).Error(),
			errors.Wrap(err, errs.ErrSessionNotExists).Error())
		return
	}

	logger.Info().Interface("session username", username).Msg("getSession success")

	err = jsonutil.SendJSON(r.Context(), w, ds.User{Username: username})
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSendJSON).Error())
		return
	}
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
