package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=$GOFILE -destination=service_mocks/mock_repository.go -package=service_mocks MovieRepositoryInterface
type MovieRepositoryInterface interface {
	GetMovieFromRepoByID(ctx context.Context, movieID int) (*mocks.MovieJSON, error)
}

type MovieService struct {
	movieRepo MovieRepositoryInterface
}

func NewMovieService(movieRepo MovieRepositoryInterface) *MovieService {
	return &MovieService{
		movieRepo: movieRepo,
	}
}

func (s *MovieService) GetMovieByID(ctx context.Context, movieID int) (*mocks.MovieJSON, error) {
	logger := log.Ctx(ctx)

	movieJSON, err := s.movieRepo.GetMovieFromRepoByID(ctx, movieID)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		return nil, err
	}

	return movieJSON, nil
}
