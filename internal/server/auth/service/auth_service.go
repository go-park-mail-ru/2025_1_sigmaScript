package service

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/repository"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/validation/auth"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/session"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

const (
	noData = ""
)

type AuthService struct {
	authRepo repository.AuthRepositoryInterface
	cfg      *config.Cookie
}

type AuthServiceInterface interface {
	Register(ctx context.Context, regUser models.RegisterData) (string, *errs.ServiceError)
	GetSession(ctx context.Context, sessionID string) (string, *errs.ServiceError)
	DeleteSession(ctx context.Context, sessionID string) *errs.ServiceError
	Login(ctx context.Context, login models.LoginData) (string, *errs.ServiceError)
	Logout(ctx context.Context, sessionID string) *errs.ServiceError
}

func NewAuthService(ctx context.Context, authRepo repository.AuthRepositoryInterface) AuthServiceInterface {
	res := &AuthService{
		cfg:      config.FromCookieContext(ctx),
		authRepo: authRepo,
	}

	return res
}

// Register user with given parameters
func (a *AuthService) Register(ctx context.Context, regUser models.RegisterData) (string, *errs.ServiceError) {
	logger := log.Ctx(ctx)

	logger.Info().Msg("Registering user")

	if regUser.Password != regUser.RepeatedPassword {
		logger.Err(errors.New(errs.ErrPasswordsMismatch)).Msg(errs.ErrPasswordsMismatch)
		return noData, &errs.ServiceError{
			Code:  http.StatusBadRequest,
			Error: errors.New(errs.ErrPasswordsMismatch),
		}
	}

	if err := auth.IsValidPassword(regUser.Password); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrInvalidPassword)).Msg(errors.Wrap(err, errs.ErrInvalidPassword).Error())
		return noData, &errs.ServiceError{
			Code:  http.StatusBadRequest,
			Error: errors.Wrap(err, errs.ErrInvalidPassword),
		}
	}

	if err := auth.IsValidLogin(regUser.Username); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrInvalidLogin)).Msg(errors.Wrap(err, errs.ErrInvalidLogin).Error())
		return noData, &errs.ServiceError{
			Code:  http.StatusBadRequest,
			Error: errors.Wrap(err, errs.ErrInvalidLogin),
		}
	}

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(regUser.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrBcrypt)).Msg(errors.Wrap(err, errs.ErrBcrypt).Error())
		return noData, &errs.ServiceError{
			Code:  http.StatusInternalServerError,
			Error: errors.Wrap(err, errs.ErrInvalidPassword),
		}
	}

	// create user
	errRepo := a.authRepo.CreateUser(ctx, regUser.Username, string(hashedPass))
	if errRepo != nil {
		msg := errRepo.Msg
		logger.Error().Err(errRepo.Error).Msg(msg)
		return noData, &errs.ServiceError{
			Code:  http.StatusBadRequest,
			Error: errRepo.Error,
		}
	}

	logger.Info().Msg("User stored successfully")

	newSessionID, errSession := a.createNewSession(ctx)
	if errSession != nil {
		return noData, &errs.ServiceError{
			Code:  http.StatusInternalServerError,
			Error: errors.New(errs.ErrSomethingWentWrong),
		}
	}

	errRepo = a.authRepo.StoreSession(ctx, newSessionID, regUser.Username)
	if errRepo != nil {
		msg := errRepo.Msg
		logger.Error().Err(errRepo.Error).Msg(msg)
		return noData, &errs.ServiceError{
			Code:  http.StatusBadRequest,
			Error: errRepo.Error,
		}
	}

	logger.Info().Msg("Session created")

	return newSessionID, nil
}

func (a *AuthService) createNewSession(ctx context.Context) (string, *errs.ServiceError) {
	logger := log.Ctx(ctx)

	// create new session for user
	newSessionID, err := session.GenerateSessionID(a.cfg.SessionLength)
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrGenerateSession)).Msg(errors.Wrap(err, errs.ErrGenerateSession).Error())
		return noData, &errs.ServiceError{
			Code:  http.StatusInternalServerError,
			Error: errors.Wrap(err, errs.ErrSomethingWentWrong),
		}
	}

	return newSessionID, nil
}

func (a *AuthService) DeleteSession(ctx context.Context, sessionID string) *errs.ServiceError {
	logger := log.Ctx(ctx)

	errRepo := a.authRepo.DeleteSession(ctx, sessionID)
	if errRepo != nil {
		msg := errRepo.Msg
		logger.Error().Err(errRepo.Error).Msg(msg)
		return &errs.ServiceError{
			Code:  http.StatusBadRequest,
			Error: errRepo.Error,
		}
	}

	logger.Info().Msg("session successfully deleted")

	return nil
}

// Session http handler method
func (a *AuthService) GetSession(ctx context.Context, sessionID string) (string, *errs.ServiceError) {
	logger := log.Ctx(ctx)

	username, errRepo := a.authRepo.GetSession(ctx, sessionID)
	if errRepo != nil {
		logger.Error().Err(errors.Wrap(errRepo.Error, errs.ErrSessionNotExists)).Msg(errRepo.Msg)
		return noData, &errs.ServiceError{
			Code:  http.StatusUnauthorized,
			Error: errors.New(errs.ErrSessionNotExists),
		}
	}

	return username, nil
}

// Login http handler method
func (a *AuthService) Login(ctx context.Context, loginData models.LoginData) (string, *errs.ServiceError) {
	logger := log.Ctx(ctx)

	hashedPass, errRepo := a.authRepo.GetUser(ctx, loginData.Username)
	if errRepo != nil {
		errMsg := errors.New(errs.ErrIncorrectLogin)
		logger.Error().Err(errors.Wrap(errMsg, errs.ErrIncorrectLoginOrPassword)).Msg(errMsg.Error())
		return noData, &errs.ServiceError{
			Code:  http.StatusUnauthorized,
			Error: errors.Wrap(errMsg, errs.ErrIncorrectLoginOrPassword),
		}
	} else if err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(loginData.Password)); err != nil {
		errMsg := errors.New(errs.ErrIncorrectPassword)
		logger.Error().Err(errors.Wrap(err, errs.ErrIncorrectLoginOrPassword)).Msg(errMsg.Error())
		return noData, &errs.ServiceError{
			Code:  http.StatusUnauthorized,
			Error: errors.Wrap(errMsg, errs.ErrIncorrectLoginOrPassword),
		}
	}

	newSessionID, errSession := a.createNewSession(ctx)
	if errSession != nil {
		return noData, &errs.ServiceError{
			Code:  http.StatusInternalServerError,
			Error: errors.New(errs.ErrSomethingWentWrong),
		}
	}

	// repo session
	errRepo = a.authRepo.StoreSession(ctx, newSessionID, loginData.Username)
	if errRepo != nil {
		msg := errRepo.Msg
		logger.Error().Err(errRepo.Error).Msg(msg)
		return noData, &errs.ServiceError{
			Code:  http.StatusBadRequest,
			Error: errRepo.Error,
		}
	}

	logger.Info().Msg("Session created")

	return newSessionID, nil
}

// Logout http handler method
func (a *AuthService) Logout(ctx context.Context, sessionID string) *errs.ServiceError {
	logger := log.Ctx(ctx)

	if _, errRepo := a.authRepo.GetSession(ctx, sessionID); errRepo != nil {
		err := errors.New("session does not exist")
		logger.Err(errors.Wrap(errRepo.Error, err.Error())).Msg(errs.ErrSessionNotExists)
		return &errs.ServiceError{
			Code:  http.StatusNotFound,
			Error: errors.Wrap(err, errs.ErrSessionNotExists),
		}
	}

	errRepo := a.authRepo.DeleteSession(ctx, sessionID)
	if errRepo != nil {
		msg := errRepo.Msg
		logger.Error().Err(errRepo.Error).Msg(msg)
		return &errs.ServiceError{
			Code:  http.StatusBadRequest,
			Error: errRepo.Error,
		}
	}

	return nil
}
