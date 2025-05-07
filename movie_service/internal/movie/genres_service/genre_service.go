package genres_service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mocks.go -package=service_mocks GenreRepositoryInterface
type GenreRepositoryInterface interface {
	GetGenreFromRepoByID(ctx context.Context, genreID string) (*mocks.Genre, error)
	GetAllGenresFromRepo(ctx context.Context) (*[]mocks.Genre, error)
}

type GenreService struct {
	movieRepo GenreRepositoryInterface
}

func NewGenreService(movieRepo GenreRepositoryInterface) *GenreService {
	return &GenreService{
		movieRepo: movieRepo,
	}
}

func (s *GenreService) GetGenreByID(ctx context.Context, genreID string) (*mocks.Genre, error) {
	logger := log.Ctx(ctx)

	genre, err := s.movieRepo.GetGenreFromRepoByID(ctx, genreID)
	if err != nil {
		logger.Error().Err(err).Msgf("error happened while getting genre by id %s: %v", genreID, err.Error())
		return nil, err
	}

	return genre, nil
}

func (s *GenreService) GetAllGenres(ctx context.Context) (*[]mocks.Genre, error) {
	logger := log.Ctx(ctx)

	genres, err := s.movieRepo.GetAllGenresFromRepo(ctx)
	if err != nil {
		logger.Error().Err(err).Msgf("error happened while getting all genres: %v", err.Error())
		return nil, err
	}

	return genres, nil
}
