package handlers

import (
  "context"
  "net/http"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/config"
  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
  "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
  "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/session"
  "github.com/rs/zerolog/log"
  "golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
  // username --> hashedPass
  users map[string]string
  // sessionID --> username
  sessions map[string]string
  Config   *config.Cookie
}

func NewAuthHandler(ctx context.Context) *AuthHandler {
  return &AuthHandler{
    users:    make(map[string]string),
    sessions: make(map[string]string),
    Config:   config.FromCookieContext(ctx),
  }
}

func (a *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
  var reg models.RegisterData
  log.Info().Msg("Registering user")

  if err := jsonutil.ReadJSON(r, &reg); err != nil {
    log.Error().Err(err).Msg("Error parsing JSON")
    jsonutil.SendError(w, http.StatusBadRequest, "incorrect_data", "Incorrect data")
    return
  }

  if _, exists := a.users[reg.Username]; exists {
    log.Error().Msg("User already registered")
    jsonutil.SendError(w, http.StatusBadRequest, "already_exists", "User already exists")
    return
  }

  if reg.Password != reg.RepeatedPassword {
    log.Info().Msg("Password mismatch")
    jsonutil.SendError(w, http.StatusBadRequest, "password_mismatch", "Password mismatch")
    return
  }

  hashedPass, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
  if err != nil {
    log.Error().Err(err).Msg("Error hashing password")
    jsonutil.SendError(w, http.StatusInternalServerError, "internal_error", "Failed to hash password")
    return
  }

  a.users[reg.Username] = string(hashedPass)
  if err = jsonutil.SendJSON(w, map[string]string{"message": "Successfully register"}); err != nil {
    log.Error().Err(err).Msg("Error sending JSON")
    return
  }
  log.Info().Msg("User registered successfully")
}

func (a *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
  var login models.LoginData
  log.Info().Msg("Logining user")

  if err := jsonutil.ReadJSON(r, &login); err != nil {
    log.Error().Err(err).Msg("Error parsing JSON")
    jsonutil.SendError(w, http.StatusBadRequest, "incorrect_data", "Incorrect data")
    return
  }

  hashedPass, exists := a.users[login.Username]
  if err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(login.Password)); err != nil || !exists {
    log.Error().Err(err).Msg("Login or password incorrect")
    jsonutil.SendError(w, http.StatusUnauthorized, "not_found", "Login or password incorrect")
    return
  }

  sessionID, err := session.GenerateSessionID(a.Config.SessionLength)
  if err != nil {
    log.Error().Err(err).Msg("Error generating session ID")
    jsonutil.SendError(w, http.StatusInternalServerError, "internal_error", "Failed to generate session ID")
    return
  }

  a.sessions[sessionID] = login.Username
  log.Info().Msg("Session created")

  http.SetCookie(w, &http.Cookie{
    Name:     a.Config.SessionName,
    Value:    sessionID,
    HttpOnly: a.Config.HTTPOnly,
    Secure:   a.Config.Secure,
    SameSite: a.Config.SameSite,
    Path:     a.Config.Path,
  })

  err = jsonutil.SendJSON(w, map[string]string{"message": "Successfully logged in"})
  if err != nil {
    log.Error().Err(err).Msg("Error sending JSON")
    return
  }
}

func (a *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
  log.Info().Msg("Logouting user")

  cookie, err := r.Cookie("session_id")
  if err != nil {
    log.Warn().Msg("Unauthorized")
    jsonutil.SendError(w, http.StatusUnauthorized, "unauthorized", "Unauthorized")
    return
  }

  sessionID := cookie.Value
  if _, exists := a.sessions[sessionID]; !exists {
    log.Warn().Msg("Session does not exist")
    jsonutil.SendError(w, http.StatusNotFound, "not_exists", "Session does not exist")
    return
  }

  delete(a.sessions, sessionID)
  http.SetCookie(w, &http.Cookie{
    Name:     a.Config.SessionName,
    Value:    "",
    HttpOnly: a.Config.HTTPOnly,
    Secure:   a.Config.Secure,
    SameSite: a.Config.SameSite,
    Path:     a.Config.Path,
    MaxAge:   a.Config.ResetMaxAge,
  })
  log.Info().Msg("Session deleted")

  err = jsonutil.SendJSON(w, map[string]string{"message": "Successfully logged out"})
  if err != nil {
    log.Error().Err(err).Msg("Error sending JSON")
    return
  }
}
