package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/rs/zerolog/log"
)

func (s *UserService) UpdateUser(ctx context.Context, login string, newUser *models.User) error {
	logger := log.Ctx(ctx)

	// if err := s.DeleteUser(ctx, login); err != nil {
	// 	logger.Error().Err(err).Msg(err.Error())
	// 	return err
	// }

	// if err := s.repo.CreateUserPostgres(ctx, newUser); err != nil {
	// 	logger.Error().Err(err).Msg(err.Error())
	// 	return err
	// }

	if _, err := s.repo.UpdateUserPostgres(ctx, login, newUser); err != nil {
		logger.Error().Err(err).Msg(err.Error())
		return err
	}

	return nil
}
