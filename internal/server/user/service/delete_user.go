package service

import (
	"context"

	"github.com/rs/zerolog/log"
)

func (s *UserService) DeleteUser(ctx context.Context, login string) error {
	err := s.repo.DeleteUser(ctx, login)
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
		return err
	}

	return nil
}
