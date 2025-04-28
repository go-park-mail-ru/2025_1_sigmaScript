package cookie

import (
	"context"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	CookieDaysLimit       = 3
	CookieExpiredLastYear = -1
)

type SessionServiceInterface interface {
	DeleteSession(ctx context.Context, sessionID string) error
}

// this function returns result of expiring old session cookie
// if something happend while parsing old cookie function returns error.
func ExpireOldSessionCookie(w http.ResponseWriter, r *http.Request, cookie *config.Cookie, sessionSrv SessionServiceInterface) error {
	logger := log.Ctx(r.Context())

	oldSessionCookie, err := r.Cookie("session_id")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			logger.Info().Msg("user dont have old cookie")
			return nil
		}
		return err
	}

	if oldSessionCookie != nil {
		http.SetCookie(w, PreparedExpiredCookie(cookie))
		err := sessionSrv.DeleteSession(r.Context(), oldSessionCookie.Value)
		if err != nil {
			logger.Warn().Err(err).Msg(errors.Wrap(err, "error happend while deleting old session from repo").Error())
			return nil
		}
	}

	logger.Info().Msg("successfully expired old sesssion cookie")
	return nil
}

func PreparedNewCookie(cookie *config.Cookie, newSessionID string) *http.Cookie {
	return &http.Cookie{
		Name:     cookie.SessionName,
		Value:    newSessionID,
		HttpOnly: cookie.HTTPOnly,
		Secure:   cookie.Secure,
		SameSite: cookie.SameSite,
		Path:     cookie.Path,
		Expires:  time.Now().AddDate(0, 0, CookieDaysLimit),
	}
}

func PreparedExpiredCookie(cookie *config.Cookie) *http.Cookie {
	return &http.Cookie{
		Name:     cookie.SessionName,
		Value:    "",
		HttpOnly: cookie.HTTPOnly,
		Secure:   cookie.Secure,
		SameSite: cookie.SameSite,
		Path:     cookie.Path,
		Expires:  time.Now().AddDate(CookieExpiredLastYear, 0, 0),
	}
}
