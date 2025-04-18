package repository

import (
	"context"
	"testing"
	"time"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMovieRepository_GetMovieFromRepoByID(t *testing.T) {
	db := mocks.ExistingMovies
	repo := NewMovieRepository(&db)
	fightClub := mocks.ExistingMovies[1]
	matrix := mocks.ExistingMovies[2]

	tests := []struct {
		name        string
		movieID     int
		expected    *mocks.MovieJSON
		expectedErr error
	}{
		{
			name:        "OK. Get Fight Club",
			movieID:     1,
			expected:    &fightClub,
			expectedErr: nil,
		},
		{
			name:        "OK. Get Matrix",
			movieID:     2,
			expected:    &matrix,
			expectedErr: nil,
		},
		{
			name:        "Fail. Movie not found",
			movieID:     999,
			expected:    nil,
			expectedErr: errs.ErrMovieNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			movie, err := repo.GetMovieFromRepoByID(context.Background(), tt.movieID)

			assert.Equal(t, tt.expected, movie)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestMovieRepository_GetAllReviewsOfMovieFromRepoByID(t *testing.T) {
	db := mocks.ExistingMovies
	repo := NewMovieRepository(&db)
	fightClub := mocks.ExistingMovies[1]
	matrix := mocks.ExistingMovies[2]

	tests := []struct {
		name        string
		movieID     int
		expected    *[]mocks.ReviewJSON
		expectedErr error
	}{
		{
			name:        "OK. Get Fight Club reviews",
			movieID:     1,
			expected:    &fightClub.Reviews,
			expectedErr: nil,
		},
		{
			name:        "OK. Get Matrix reviews (empty)",
			movieID:     2,
			expected:    &matrix.Reviews,
			expectedErr: nil,
		},
		{
			name:        "Fail. Movie not found",
			movieID:     999,
			expected:    nil,
			expectedErr: errs.ErrMovieNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reviews, err := repo.GetAllReviewsOfMovieFromRepoByID(context.Background(), tt.movieID)

			assert.Equal(t, tt.expected, reviews)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestMovieRepository_CreateNewMovieReviewInRepo(t *testing.T) {
	db := mocks.ExistingMovies
	repo := NewMovieRepository(&db)

	newReview := mocks.ReviewJSON{
		User:       mocks.ReviewUserDataJSON{Login: "test_user"},
		ReviewText: "Great movie!",
		Score:      9,
		CreatedAt:  time.Now().String(),
	}

	tests := []struct {
		name        string
		movieID     int
		review      mocks.ReviewJSON
		expectedErr error
	}{
		{
			name:        "OK. Create new review for existing movie",
			movieID:     1,
			review:      newReview,
			expectedErr: nil,
		},
		{
			name:        "Fail. Movie not found",
			movieID:     999,
			review:      newReview,
			expectedErr: errs.ErrMovieNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.CreateNewMovieReviewInRepo(context.Background(), tt.movieID, tt.review)

			assert.ErrorIs(t, err, tt.expectedErr)

			if tt.expectedErr == nil {
				movie := db[tt.movieID]
				found := false
				for _, r := range movie.Reviews {
					if r.User.Login == tt.review.User.Login {
						found = true
						break
					}
				}
				assert.True(t, found, "Review should be added to movie")
			}
		})
	}
}
