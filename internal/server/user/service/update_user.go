package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/rs/zerolog/log"
)

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
