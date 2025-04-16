package middleware

import (
	"net/http"
	"strings"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/common"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// CsrfMiddleware checks CSRF token in http Header and cookie
func CsrfTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.URL.Path == "/csrf-token" || strings.HasPrefix(r.URL.Path, "/auth") {
			next.ServeHTTP(w, r)
			return
		}

		logger := log.Ctx(r.Context())

		logger.Info().Msg("Checking CSRF")
		cookieCSRFtoken, err := r.Cookie(common.CSRF_TOKEN_NAME)
		if err != nil {
			logger.Error().Err(err).Msg(errors.Wrap(err, errs.ErrMsgBadCSRFToken).Error())
			jsonutil.SendError(r.Context(), w, http.StatusForbidden, errs.ErrUnauthorizedShort,
				errs.ErrMsgBadCSRFToken)
			return
		}

		token := cookieCSRFtoken.Value
		logger.Info().Msgf("got CSRF from cookie %v", token)

		headerCSRFtoken := r.Header.Get("X-CSRF-Token")
		logger.Info().Msgf("got CSRF from HEADER /%v/.", headerCSRFtoken)

		// compare tokens
		if headerCSRFtoken == "" || strings.Compare(headerCSRFtoken, token) != 0 {
			logger.Error().Msg(errors.New(errs.ErrMsgBadCSRFToken).Error())
			jsonutil.SendError(r.Context(), w, http.StatusForbidden, errs.ErrUnauthorizedShort,
				errs.ErrMsgBadCSRFToken)
			return
		}

		next.ServeHTTP(w, r)
	})
}
