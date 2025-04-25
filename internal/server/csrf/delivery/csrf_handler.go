package delivery

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/common"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/ds"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/messages"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/delivery/interfaces"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/cookie"
	csrftoken "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/csrf_token"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/pkg/errors"

	"github.com/rs/zerolog/log"
)

type CSRFHandler struct {
	cookieData     *config.Cookie
	sessionService interfaces.SessionServiceInterface
}

func NewCSRFHandler(ctx context.Context, sessionService interfaces.SessionServiceInterface) *CSRFHandler {
	newCookieData := (*config.FromCookieContext(ctx))
	newCookieData.SessionName = common.CSRF_TOKEN_NAME
	return &CSRFHandler{
		cookieData:     &newCookieData,
		sessionService: sessionService,
	}
}

// CreateCSRFTokenHandler создает CSRF-токен и сохраняет его в куки
func (h *CSRFHandler) CreateCSRFTokenHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	logger.Info().Msg("Checking session")
	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		logger.Warn().Msg(errors.Wrap(err, errs.ErrUnauthorized).Error())
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errs.ErrUnauthorizedShort,
			errs.ErrUnauthorized)
		return
	}

	// check if user session is valid
	_, errSession := h.sessionService.GetSession(r.Context(), sessionCookie.Value)
	if errSession != nil {
		logger.Error().Err(errors.Wrap(errSession, errs.ErrMsgSessionNotExists)).Msg(errs.ErrMsgFailedToGetSession)
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errs.ErrMsgSessionNotExists, errs.ErrMsgFailedToGetSession)
		return
	}

	token, err := csrftoken.GenerateCSRFToken()
	if err != nil {
		http.Error(w, "Failed to generate CSRF token", http.StatusInternalServerError)
		logger.Info().Msg("Error generating CSRF token")
		return
	}

	logger.Info().Msg("Successfully created CSRF token")
	// пишем хедер с токеном
	w.Header().Set("X-CSRF-Token", token)

	// пишем токен в куки
	http.SetCookie(w, cookie.PreparedNewCookie(h.cookieData, token))

	if err := jsonutil.SendJSON(r.Context(), w, ds.Response{Message: messages.SuccessfulNewCSRFToken}); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrSendJSON)).Msg(errors.Wrap(err, errs.ErrSomethingWentWrong).Error())
		return
	}
}
