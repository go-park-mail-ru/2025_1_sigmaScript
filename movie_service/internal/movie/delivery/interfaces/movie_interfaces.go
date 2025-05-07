package interfaces

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
)

type MovieServiceInterface interface {
	GetMovieByID(ctx context.Context, movieID int) (*mocks.MovieJSON, error)
	GetAllReviewsOfMovieByID(ctx context.Context, movieID int) (*[]mocks.ReviewJSON, error)
	CreateNewMovieReview(ctx context.Context,
		userID string,
		movieID string,
		newReview mocks.NewReviewDataJSON) (*mocks.NewReviewDataJSON, error)
	UpdateMovieReview(ctx context.Context, userID string, movieID string, newReview mocks.NewReviewDataJSON) (*mocks.NewReviewDataJSON, error)
}

type GenreServiceInterface interface {
	GetGenreByID(ctx context.Context, genreID string) (*mocks.Genre, error)
	GetAllGenres(ctx context.Context) (*[]mocks.Genre, error)
}

type SearchServiceInterface interface {
	SearchActorsAndMovies(ctx context.Context, searchStr string) (*models.SearchResponseJSON, error)
}

type StaffPersonServiceInterface interface {
	GetPersonByID(ctx context.Context, personID int) (*mocks.PersonJSON, error)
}

type CollectionServiceInterface interface {
	GetMainPageCollections(ctx context.Context) (mocks.Collections, error)
}
