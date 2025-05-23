package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/rs/zerolog/log"
)

func (s *UserService) GetUser(ctx context.Context, login string) (*models.User, error) {
	user, err := s.repo.GetUser(ctx, login)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return nil, err
	}

	return user, nil
}
