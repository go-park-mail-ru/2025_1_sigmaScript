package handlers

import (
  "context"
  "net/http"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/config"
  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/ds"
  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/validation/auth"
  "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
  "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/session"
  "github.com/pkg/errors"
  "github.com/rs/zerolog/log"
  "golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
  // username --> hashedPass
  users map[string]string
  // sessionID --> username
  sessions map[string]string
  config   *config.Cookie
}

func NewAuthHandler(ctx context.Context) *AuthHandler {
  return &AuthHandler{
    users:    make(map[string]string),
    sessions: make(map[string]string),
    config:   config.FromCookieContext(ctx),
  }
}

func (a *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
  var reg models.RegisterData
  log.Info().Msg("Registering user")

  if err := jsonutil.ReadJSON(r, &reg); err != nil {
    log.Error().Err(errors.Wrap(err, errs.ErrParseJSON)).Msg(errors.Wrap(err, errs.ErrParseJSON).Error())
    jsonutil.SendError(w, http.StatusBadRequest, errors.Wrap(err, errs.ErrParseJSONShort).Error(),
      errors.Wrap(err, errs.ErrParseJSON).Error())
    return
  }

  if _, exists := a.users[reg.Username]; exists {
    log.Error().Err(errors.New(errs.ErrAlreadyExists)).Msg(errs.ErrAlreadyExists)
    jsonutil.SendError(w, http.StatusBadRequest, errors.New(errs.ErrAlreadyExistsShort).Error(), errs.ErrAlreadyExists)
    return
  }

  if reg.Password != reg.RepeatedPassword {
    log.Info().Msg("Passwords mismatch")
    jsonutil.SendError(w, http.StatusBadRequest, errors.New(errs.ErrPasswordsMismatchShort).Error(), errs.ErrPasswordsMismatch)
    return
  }

  if err := auth.IsValidPassword(reg.Password); err != nil {
    log.Error().Err(errors.Wrap(err, errs.ErrInvalidPassword)).Msg(errors.Wrap(err, errs.ErrInvalidPassword).Error())
    jsonutil.SendError(w, http.StatusBadRequest, errors.Wrap(err, errs.ErrInvalidPasswordShort).Error(), errors.Wrap(err, errs.ErrInvalidPassword).Error())
    return
  }

  hashedPass, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
  if err != nil {
    log.Error().Err(errors.Wrap(err, errs.ErrBcrypt)).Msg(errors.Wrap(err, errs.ErrBcrypt).Error())
    jsonutil.SendError(w, http.StatusInternalServerError, errors.Wrap(err, errs.ErrBcryptShort).Error(),
      errors.Wrap(err, errs.ErrBcrypt).Error())
    return
  }

  a.users[reg.Username] = string(hashedPass)
  if err = jsonutil.SendJSON(w, map[string]string{"message": ds.SuccessfulRegister}); err != nil {
    log.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSendJSON).Error())
    return
  }
  log.Info().Msg("User registered successfully")
}

func (a *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
  var login models.LoginData
  log.Info().Msg("Logining user")

  if err := jsonutil.ReadJSON(r, &login); err != nil {
    log.Error().Err(errors.Wrap(err, errs.ErrParseJSON)).Msg(errors.Wrap(err, errs.ErrParseJSON).Error())
    jsonutil.SendError(w, http.StatusBadRequest, errors.Wrap(err, errs.ErrParseJSONShort).Error(),
      errors.Wrap(err, errs.ErrParseJSON).Error())
    return
  }

  hashedPass, exists := a.users[login.Username]
  if err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(login.Password)); err != nil || !exists {
    log.Error().Err(errors.Wrap(err, errs.ErrIncorrectLoginOrPassword)).Msg(errors.Wrap(err, errs.ErrIncorrectLoginOrPassword).Error())
    jsonutil.SendError(w, http.StatusUnauthorized, errors.Wrap(err, errs.ErrIncorrectLoginOrPasswordShort).Error(),
      errors.Wrap(err, errs.ErrIncorrectLoginOrPassword).Error())
    return
  }

  sessionID, err := session.GenerateSessionID(a.config.SessionLength)
  if err != nil {
    log.Error().Err(errors.Wrap(err, errs.ErrGenerateSession)).Msg(errors.Wrap(err, errs.ErrGenerateSession).Error())
    jsonutil.SendError(w, http.StatusInternalServerError, errors.Wrap(err, errs.ErrGenerateSessionShort).Error(),
      errors.Wrap(err, errs.ErrGenerateSession).Error())
    return
  }

  a.sessions[sessionID] = login.Username
  log.Info().Msg("Session created")

  http.SetCookie(w, &http.Cookie{
    Name:     a.config.SessionName,
    Value:    sessionID,
    HttpOnly: a.config.HTTPOnly,
    Secure:   a.config.Secure,
    SameSite: a.config.SameSite,
    Path:     a.config.Path,
  })

  err = jsonutil.SendJSON(w, map[string]string{"message": ds.SuccessfulLogin})
  if err != nil {
    log.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSendJSON).Error())
    return
  }
}

func (a *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
  log.Info().Msg("Logouting user")

  cookie, err := r.Cookie("session_id")
  if err != nil {
    log.Warn().Msg(errors.Wrap(err, errs.ErrUnauthorized).Error())
    jsonutil.SendError(w, http.StatusUnauthorized, errors.Wrap(err, errs.ErrUnauthorizedShort).Error(),
      errors.Wrap(err, errs.ErrUnauthorized).Error())
    return
  }

  sessionID := cookie.Value
  if _, exists := a.sessions[sessionID]; !exists {
    log.Warn().Msg(errs.ErrSessionNotExists)
    jsonutil.SendError(w, http.StatusNotFound, errors.Wrap(err, errs.ErrSessionNotExistsShort).Error(), errs.ErrSessionNotExists)
    return
  }

  delete(a.sessions, sessionID)
  http.SetCookie(w, &http.Cookie{
    Name:     a.config.SessionName,
    Value:    "",
    HttpOnly: a.config.HTTPOnly,
    Secure:   a.config.Secure,
    SameSite: a.config.SameSite,
    Path:     a.config.Path,
    MaxAge:   a.config.ResetMaxAge,
  })
  log.Info().Msg("Session deleted")

  err = jsonutil.SendJSON(w, map[string]string{"message": ds.SuccessfulLogout})
  if err != nil {
    log.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSendJSON).Error())
    return
  }
}
