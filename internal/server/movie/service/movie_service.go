package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mocks.go -package=service_mocks MovieRepositoryInterface
type MovieRepositoryInterface interface {
	GetMovieFromRepoByID(ctx context.Context, movieID int) (*mocks.MovieJSON, error)
	GetAllReviewsOfMovieFromRepoByID(ctx context.Context, movieID int) (*[]mocks.ReviewJSON, error)
	CreateNewMovieReviewInRepo(ctx context.Context, movieID int, newReview mocks.ReviewJSON) error
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

func (s *MovieService) GetAllReviewsOfMovieByID(ctx context.Context, movieID int) (*[]mocks.ReviewJSON, error) {
	logger := log.Ctx(ctx)

	movieReviews, err := s.movieRepo.GetAllReviewsOfMovieFromRepoByID(ctx, movieID)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		return nil, err
	}

	return movieReviews, nil
}

func (s *MovieService) CreateNewMovieReview(ctx context.Context, movieID int, newReview mocks.ReviewJSON) error {
	logger := log.Ctx(ctx)

	err := s.movieRepo.CreateNewMovieReviewInRepo(ctx, movieID, newReview)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		return err
	}

	return nil
}
