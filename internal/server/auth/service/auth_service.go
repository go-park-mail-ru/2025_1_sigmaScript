package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/session"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

const (
	noData = ""
)

type AuthRepositoryInterface interface {
	CreateUser(ctx context.Context, login, hashedPass string) error
	GetUser(ctx context.Context, login string) (hashedPass string, errRepo error)
	StoreSession(ctx context.Context, newSessionID, login string) error
	DeleteSession(ctx context.Context, sessionID string) error
	GetSession(ctx context.Context, sessionID string) (string, error)
}

type AuthService struct {
	authRepo      AuthRepositoryInterface
	sessionLength int
}

func NewAuthService(ctx context.Context, authRepo AuthRepositoryInterface) *AuthService {
	return &AuthService{
		sessionLength: config.FromCookieContext(ctx).SessionLength,
		authRepo:      authRepo,
	}
}

// Register user with given parameters
func (s *AuthService) Register(ctx context.Context, regUser models.RegisterData) (string, error) {
	logger := log.Ctx(ctx)

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(regUser.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrBcrypt)).Msg(errs.ErrInvalidPassword)
		return noData, errors.New(errs.ErrInvalidPassword)
	}

	// create user
	errRepo := s.authRepo.CreateUser(ctx, regUser.Username, string(hashedPass))
	if errRepo != nil {
		logger.Error().Err(errRepo).Msg(errRepo.Error())
		return noData, errRepo
	}

	// session service
	newSessionID, errSession := s.createNewSession(ctx)
	if errSession != nil {
		return noData, errors.New(errs.ErrSomethingWentWrong)
	}

	errRepo = s.authRepo.StoreSession(ctx, newSessionID, regUser.Username)
	if errRepo != nil {
		logger.Error().Err(errRepo).Msg(errRepo.Error())
		return noData, errRepo
	}
	logger.Info().Msg("Session created")
	/////////////////////////////

	return newSessionID, nil
}

func (s *AuthService) createNewSession(ctx context.Context) (string, error) {
	logger := log.Ctx(ctx)

	// create new session for user
	newSessionID, err := session.GenerateSessionID(s.sessionLength)
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrGenerateSession)).Msg(errors.Wrap(err, errs.ErrGenerateSession).Error())
		return noData, errors.Wrap(err, errs.ErrSomethingWentWrong)
	}

	return newSessionID, nil
}

func (s *AuthService) DeleteSession(ctx context.Context, sessionID string) error {
	logger := log.Ctx(ctx)

	errRepo := s.authRepo.DeleteSession(ctx, sessionID)
	if errRepo != nil {
		logger.Error().Err(errRepo).Msg(errRepo.Error())
		return errRepo
	}

	logger.Info().Msg("session successfully deleted")

	return nil
}

// Session http handler method
func (s *AuthService) GetSession(ctx context.Context, sessionID string) (string, error) {
	logger := log.Ctx(ctx)

	username, errRepo := s.authRepo.GetSession(ctx, sessionID)
	if errRepo != nil {
		logger.Error().Err(errors.Wrap(errRepo, errs.ErrSessionNotExists)).Msg(errRepo.Error())
		return noData, errRepo
	}

	return username, nil
}

// Login http handler method
func (s *AuthService) Login(ctx context.Context, loginData models.LoginData) (string, error) {
	logger := log.Ctx(ctx)

	hashedPass, errRepo := s.authRepo.GetUser(ctx, loginData.Username)
	if errRepo != nil {
		logger.Error().Err(errors.Wrap(errRepo, errs.ErrIncorrectLoginOrPassword)).Msg(errRepo.Error())
		return noData, errRepo
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(loginData.Password)); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrIncorrectLoginOrPassword)).Msg(errs.ErrIncorrectPassword)
		return noData, errors.New(errs.ErrIncorrectPassword)
	}

	newSessionID, errSession := s.createNewSession(ctx)
	if errSession != nil {
		return noData, errors.New(errs.ErrSomethingWentWrong)
	}

	// repo session
	errRepo = s.authRepo.StoreSession(ctx, newSessionID, loginData.Username)
	if errRepo != nil {

		logger.Error().Err(errRepo).Msg(errRepo.Error())
		return noData, errRepo
	}

	logger.Info().Msg("Session created")

	return newSessionID, nil
}

// Logout http handler method
func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	logger := log.Ctx(ctx)

	if _, errRepo := s.authRepo.GetSession(ctx, sessionID); errRepo != nil {
		logger.Err(errRepo).Msg(errs.ErrSessionNotExists)
		return errRepo
	}

	errRepo := s.authRepo.DeleteSession(ctx, sessionID)
	if errRepo != nil {
		logger.Error().Err(errRepo).Msg(errRepo.Error())
		return errRepo
	}

	return nil
}
