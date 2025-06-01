package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/rs/zerolog/log"
)

func (s *UserService) GetProfile(ctx context.Context, login string) (*models.Profile, error) {
	logger := log.Ctx(ctx)

	profile, err := s.repo.GetUserProfilePostgres(ctx, login)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		return nil, err
	}

	return profile, nil
}
