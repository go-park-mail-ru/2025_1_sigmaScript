package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/rs/zerolog/log"
)

func (s *UserService) CreateUser(ctx context.Context, user *models.User) error {
	logger := log.Ctx(ctx)

	err := s.repo.CreateUser(ctx, user)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		return err
	}

	return nil
}
