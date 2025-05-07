package movie_repo

import (
	"context"
	"database/sql"
	"math"
	"strconv"
	"time"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/rs/zerolog/log"
)

const (
	DEFAULT_MOVIE_SCORE = 5
)

type MoviePostgresRepository struct {
	pgdb *sql.DB
}

func NewMoviePostgresRepository(movieDB *sql.DB) *MoviePostgresRepository {
	return &MoviePostgresRepository{pgdb: movieDB}
}

type MovieRepository struct {
	db *mocks.Movies
}

func NewMovieRepository(movieDB *mocks.Movies) *MovieRepository {
	return &MovieRepository{db: movieDB}
}

func (r *MovieRepository) GetMovieFromRepoByID(ctx context.Context, movieID int) (*mocks.MovieJSON, error) {
	logger := log.Ctx(ctx)

	movie, exists := (*r.db)[movieID]
	if !exists {
		logger.Err(errs.ErrMovieNotFound).Msg(errs.ErrMovieNotFound.Error())
		return nil, errs.ErrMovieNotFound
	}

	return &movie, nil
}

func (r *MovieRepository) GetAllReviewsOfMovieFromRepoByID(ctx context.Context, movieID int) (*[]mocks.ReviewJSON, error) {
	logger := log.Ctx(ctx)

	movie, exists := (*r.db)[movieID]
	if !exists {
		logger.Err(errs.ErrMovieNotFound).Msg(errs.ErrMovieNotFound.Error())
		return nil, errs.ErrMovieNotFound
	}

	return &movie.Reviews, nil
}

func (r *MovieRepository) CreateNewMovieReviewInRepo(ctx context.Context, movieIDstr string, newReview mocks.ReviewJSON) error {
	logger := log.Ctx(ctx)

	movieID, _ := strconv.Atoi(movieIDstr)

	movie, exists := (*r.db)[movieID]
	if !exists {
		logger.Err(errs.ErrMovieNotFound).Msg(errs.ErrMovieNotFound.Error())
		return errs.ErrMovieNotFound
	}

	reviewsCount := len(movie.Reviews)
	for key, currentReview := range movie.Reviews {
		if currentReview.User.Login == newReview.User.Login {
			oldValue := movie.Reviews[key].Score
			movie.Reviews[key] = newReview
			movie.Reviews[key].CreatedAt = time.Now().String()

			newRating := movie.Rating + (float64(newReview.Score)-float64(oldValue))/float64(reviewsCount)
			movie.Rating = float64(math.Trunc((newRating)*100)) / 100

			(*r.db)[movieID] = movie
			return nil
		}
	}

	movie.Rating = float64(math.Trunc((movie.Rating+(float64(newReview.Score)-movie.Rating)/float64(reviewsCount+1))*100)) / 100

	newReviewID := reviewsCount + 1
	movie.Reviews = append(movie.Reviews, mocks.ReviewJSON{
		ID:         newReviewID,
		Score:      newReview.Score,
		ReviewText: newReview.ReviewText,
		CreatedAt:  time.Now().String(),
		User:       newReview.User,
	})

	(*r.db)[movieID] = movie
	return nil
}

func (r *MovieRepository) UpdateMovieReviewInRepo(
	ctx context.Context,
	userID string,
	movieID string,
	newReview mocks.NewReviewDataJSON) error {
	return nil
}
