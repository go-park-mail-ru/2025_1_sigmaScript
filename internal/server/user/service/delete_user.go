package service

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (s *UserService) DeleteUser(ctx context.Context, login string) error {
	logger := log.Ctx(ctx)

	err := s.repo.DeleteUserPostgres(ctx, login)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		return err
	}

	return nil
}
