package service

import (
	"context"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, login, hashedPass string) error
	GetUser(ctx context.Context, login string) (hashedPass string, errRepo error)
}

type UserService struct {
	userRepo UserRepositoryInterface
}

func NewUserService(ctx context.Context, userRepo UserRepositoryInterface) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// Register method registers user with given parameters
func (s *UserService) Register(ctx context.Context, regUser models.RegisterData) error {
	logger := log.Ctx(ctx)

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(regUser.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrBcrypt)).Msg(errs.ErrInvalidPassword)
		return errors.New(errs.ErrInvalidPassword)
	}

	errRepo := s.userRepo.CreateUser(ctx, regUser.Username, string(hashedPass))
	if errRepo != nil {
		logger.Error().Err(errRepo).Msg(errRepo.Error())
		return errRepo
	}
	return nil
}

// Login method checks if user with given credentials exists
func (s *UserService) Login(ctx context.Context, loginData models.LoginData) error {
	logger := log.Ctx(ctx)

	hashedPass, errRepo := s.userRepo.GetUser(ctx, loginData.Username)
	if errRepo != nil {
		logger.Error().Err(errors.Wrap(errRepo, errs.ErrIncorrectLoginOrPassword)).Msg(errRepo.Error())
		return errRepo
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(loginData.Password)); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrIncorrectLoginOrPassword)).Msg(errs.ErrIncorrectPassword)
		return errors.New(errs.ErrIncorrectPassword)
	}

	return nil
}
