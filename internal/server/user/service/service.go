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
	GetUser(ctx context.Context, login string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, login string) error
}

type UserService struct {
	repo UserRepositoryInterface
}

func NewUserService(repo UserRepositoryInterface) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetUser(ctx context.Context, login string) (*models.User, error) {
	user, err := s.repo.GetUser(ctx, login)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return nil, err
	}

	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return err
	}

	return nil
}

func (s *UserService) Login(ctx context.Context, loginData models.LoginData) error {
	logger := log.Ctx(ctx)

	user, err := s.repo.GetUser(ctx, loginData.Username)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(loginData.Password)); err != nil {
		logger.Error().Err(errors.Wrap(err, errs.ErrIncorrectLoginOrPassword)).Msg(errs.ErrIncorrectPassword)
		return errors.New(errs.ErrIncorrectPassword)
	}

	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, login string) error {
	err := s.repo.DeleteUser(ctx, login)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return err
	}

	return nil
}

func (s *UserService) UpdateUser(ctx context.Context, login string, newUser *models.User) error {
	if err := s.DeleteUser(ctx, login); err != nil {
		log.Error().Err(err).Msg(err.Error())
		return err
	}

	if err := s.repo.CreateUser(ctx, newUser); err != nil {
		log.Error().Err(err).Msg(err.Error())
		return err
	}

	return nil
}
