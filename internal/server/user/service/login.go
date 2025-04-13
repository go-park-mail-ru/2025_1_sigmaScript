package service

import (
	"context"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

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
