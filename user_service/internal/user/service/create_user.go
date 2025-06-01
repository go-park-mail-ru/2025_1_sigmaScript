package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/internal/models"
	"github.com/rs/zerolog/log"
)

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	logger := log.Ctx(ctx)

	err := s.repo.CreateUserPostgres(ctx, user)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		return err
	}

	return nil
}
