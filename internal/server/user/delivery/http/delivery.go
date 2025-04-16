package http

import (
	"context"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/ds"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/delivery/interfaces"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/user/delivery/http/dto"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/validation/auth"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/cookie"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/jsonutil"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	cookieData *config.Cookie
	userSvc    interfaces.UserServiceInterface
	sessionSvc interfaces.SessionServiceInterface
}

func NewUserHandler(ctx context.Context, userSvc interfaces.UserServiceInterface, sessionSvc interfaces.SessionServiceInterface) *UserHandler {
	return &UserHandler{
		cookieData: config.FromCookieContext(ctx),
		userSvc:    userSvc,
		sessionSvc: sessionSvc,
	}
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	logger := log.Ctx(r.Context())

	sessionCookie, err := r.Cookie("session_id")
	if err != nil {
		logger.Warn().Msg(errors.Wrap(err, errs.ErrUnauthorized).Error())
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errs.ErrUnauthorizedShort,
			errs.ErrUnauthorized)
		return
	}

	username, errSession := h.sessionSvc.GetSession(r.Context(), sessionCookie.Value)
	if errSession != nil {
		logger.Error().Err(errors.Wrap(errSession, errs.ErrMsgSessionNotExists)).Msg(errs.ErrMsgFailedToGetSession)
		jsonutil.SendError(r.Context(), w, http.StatusUnauthorized, errs.ErrMsgSessionNotExists, errs.ErrMsgFailedToGetSession)
		return
	}

	var userReq *dto.UpdateUserRequest
	if err = jsonutil.ReadJSON(r, &userReq); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrParseJSON)).Msg(errors.Wrap(err, errs.ErrParseJSON).Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errors.Wrap(err, errs.ErrParseJSONShort).Error(), errs.ErrBadPayload)
		return
	}

	if userReq.NewPassword != userReq.RepeatedNewPassword {
		logger.Info().Msg("Passwords mismatch")
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errors.New(errs.ErrPasswordsMismatchShort).Error(), errs.ErrPasswordsMismatch)
		return
	}

	if err = auth.IsValidPassword(userReq.NewPassword); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrInvalidPassword)).Msg(errors.Wrap(err, errs.ErrInvalidPassword).Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errors.Wrap(err, errs.ErrInvalidPasswordShort).Error(),
			errors.Wrap(err, errs.ErrInvalidPassword).Error())
		return
	}

	if err = auth.IsValidLogin(userReq.Username); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrInvalidLogin)).Msg(errors.Wrap(err, errs.ErrInvalidLogin).Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errors.Wrap(err, errs.ErrInvalidLoginShort).Error(),
			errors.Wrap(err, errs.ErrInvalidLogin).Error())
		return
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(userReq.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrBcrypt)).Msg(errors.Wrap(err, errs.ErrBcrypt).Error())
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errors.Wrap(err, errs.ErrInvalidPasswordShort).Error(),
			errors.Wrap(err, errs.ErrInvalidPassword).Error())
		return
	}

	user, err := h.userSvc.GetUser(r.Context(), username)
	if err != nil {
		wrapped := errors.Wrap(err, "error getting user")
		logger.Error().Err(wrapped).Msg(wrapped.Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, wrapped.Error(), wrapped.Error())
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(userReq.OldPassword)) != nil {
		err = errors.New(errs.ErrInvalidPassword)
		logger.Error().Err(err).Msg(err.Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, errs.ErrInvalidPasswordShort, err.Error())
		return
	}

	newUser := &models.User{
		Username:       userReq.Username,
		HashedPassword: string(hashedPass),
		Avatar:         userReq.Avatar,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      time.Now(),
	}

	if err = h.userSvc.UpdateUser(r.Context(), username, newUser); err != nil {
		wrapped := errors.Wrap(err, "error updating user")
		logger.Error().Err(wrapped).Msg(wrapped.Error())
		jsonutil.SendError(r.Context(), w, http.StatusBadRequest, wrapped.Error(), wrapped.Error())
		return
	}

	// expire old session cookie if it exists
	errOldSession := cookie.ExpireOldSessionCookie(w, r, h.cookieData, h.sessionSvc)
	if errOldSession != nil {
		logger.Warn().Err(errOldSession).Msg(errOldSession.Error())
	}

	newSessionID, err := h.sessionSvc.CreateSession(r.Context(), newUser.Username)
	if err != nil {
		logger.Error().Err(err).Msgf("error happened: %v", err.Error())

		if errors.Is(err, errs.ErrGenerateSession) {
			jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrMsgGenerateSessionShort,
				errs.ErrMsgGenerateSession)
			return
		}
		jsonutil.SendError(r.Context(), w, http.StatusInternalServerError, errs.ErrSomethingWentWrong,
			errs.ErrSomethingWentWrong)
		return
	}

	http.SetCookie(w, cookie.PreparedNewCookie(h.cookieData, newSessionID))

	if err = jsonutil.SendJSON(r.Context(), w, ds.Response{Message: "successfully updated user"}); err != nil {
		logger.Error().Err(err).Msg(errs.ErrSendJSON)
		return
	}
}
