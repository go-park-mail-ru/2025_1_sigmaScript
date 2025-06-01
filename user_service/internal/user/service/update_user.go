package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/internal/models"
	"github.com/rs/zerolog/log"
)

func (s *UserService) UpdateUser(ctx context.Context, login string, newUser *models.User) error {
	logger := log.Ctx(ctx)

	if _, err := s.repo.UpdateUserPostgres(ctx, login, newUser); err != nil {
		logger.Error().Err(err).Msg(err.Error())
		return err
	}

	return nil
}
