package service

import (
	"context"
	"testing"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	service_mocks "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMovieService_GetMovieByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := service_mocks.NewMockMovieRepositoryInterface(ctrl)
	service := NewMovieService(mockRepo)

	fightClub := mocks.ExistingMovies[1]
	matrix := mocks.ExistingMovies[2]

	tests := []struct {
		name        string
		movieID     int
		mockSetup   func()
		expected    *mocks.MovieJSON
		expectedErr error
	}{
		{
			name:    "OK. Get Fight Club",
			movieID: 1,
			mockSetup: func() {
				mockRepo.EXPECT().
					GetMovieFromRepoByID(gomock.Any(), 1).
					Return(&fightClub, nil)
			},
			expected:    &fightClub,
			expectedErr: nil,
		},
		{
			name:    "OK. Get Matrix",
			movieID: 2,
			mockSetup: func() {
				mockRepo.EXPECT().
					GetMovieFromRepoByID(gomock.Any(), 2).
					Return(&matrix, nil)
			},
			expected:    &matrix,
			expectedErr: nil,
		},
		{
			name:    "Fail. Movie not found",
			movieID: 999,
			mockSetup: func() {
				mockRepo.EXPECT().
					GetMovieFromRepoByID(gomock.Any(), 999).
					Return(nil, errs.ErrMovieNotFound)
			},
			expected:    nil,
			expectedErr: errs.ErrMovieNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			movie, err := service.GetMovieByID(context.Background(), tt.movieID)

			assert.Equal(t, tt.expected, movie)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestMovieService_GetAllReviewsOfMovieByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := service_mocks.NewMockMovieRepositoryInterface(ctrl)
	service := NewMovieService(mockRepo)

	fightClub := mocks.ExistingMovies[1]

	tests := []struct {
		name        string
		movieID     int
		mockSetup   func()
		expected    *[]mocks.ReviewJSON
		expectedErr error
	}{
		{
			name:    "OK. Get Fight Club reviews",
			movieID: 1,
			mockSetup: func() {
				mockRepo.EXPECT().
					GetAllReviewsOfMovieFromRepoByID(gomock.Any(), 1).
					Return(&fightClub.Reviews, nil)
			},
			expected:    &fightClub.Reviews,
			expectedErr: nil,
		},
		{
			name:    "Fail. Movie not found",
			movieID: 999,
			mockSetup: func() {
				mockRepo.EXPECT().
					GetAllReviewsOfMovieFromRepoByID(gomock.Any(), 999).
					Return(nil, errs.ErrMovieNotFound)
			},
			expected:    nil,
			expectedErr: errs.ErrMovieNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			reviews, err := service.GetAllReviewsOfMovieByID(context.Background(), tt.movieID)

			assert.Equal(t, tt.expected, reviews)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

func TestMovieService_CreateNewMovieReview(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := service_mocks.NewMockMovieRepositoryInterface(ctrl)
	service := NewMovieService(mockRepo)

	newReview := mocks.ReviewJSON{
		User:       mocks.ReviewUserDataJSON{Login: "test_user"},
		ReviewText: "Great movie!",
		Score:      9,
	}

	tests := []struct {
		name        string
		movieID     int
		review      mocks.ReviewJSON
		mockSetup   func()
		expectedErr error
	}{
		{
			name:    "OK. Create new review",
			movieID: 1,
			review:  newReview,
			mockSetup: func() {
				mockRepo.EXPECT().
					CreateNewMovieReviewInRepo(gomock.Any(), 1, newReview).
					Return(nil)
			},
			expectedErr: nil,
		},
		{
			name:    "Fail. Movie not found",
			movieID: 999,
			review:  newReview,
			mockSetup: func() {
				mockRepo.EXPECT().
					CreateNewMovieReviewInRepo(gomock.Any(), 999, newReview).
					Return(errs.ErrMovieNotFound)
			},
			expectedErr: errs.ErrMovieNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			err := service.CreateNewMovieReview(context.Background(), tt.movieID, tt.review)

			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}
