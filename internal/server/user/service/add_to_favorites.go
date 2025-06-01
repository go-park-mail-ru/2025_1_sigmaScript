package service

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (s *UserService) AddFavoriteMovie(ctx context.Context, login string, movieID string) error {
	logger := log.Ctx(ctx)

	err := s.repo.AddFavoriteMovie(ctx, login, movieID)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		return err
	}

	return nil
}
func (s *UserService) AddFavoriteActor(ctx context.Context, login string, actorID string) error {
	logger := log.Ctx(ctx)

	err := s.repo.AddFavoriteActor(ctx, login, actorID)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		return err
	}

	return nil
}
