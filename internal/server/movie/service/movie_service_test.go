package service

import (
	"context"
	"testing"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/service/service_mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMovieService_GetMovieByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := service_mocks.NewMockMovieRepositoryInterface(ctrl)
	testMovie := mocks.ExistingMovies[0] // "Бойцовский клуб"

	tests := []struct {
		name          string                             // Название теста
		movieID       int                                // Входной параметр
		mockSetup     func()                             // Настройка моков
		expectedMovie *mocks.MovieJSON                   // Ожидаемый фильм
		expectedErr   error                              // Ожидаемая ошибка
		checkDetails  func(*testing.T, *mocks.MovieJSON) // Доп. проверки
	}{
		// Успешное получение фильма
		{
			name:    "Success - Get Fight Club",
			movieID: 0,
			mockSetup: func() {
				mockRepo.EXPECT().
					GetMovieFromRepoByID(gomock.Any(), 0).
					Return(&testMovie, nil)
			},
			expectedMovie: &testMovie,
			expectedErr:   nil,
			checkDetails: func(t *testing.T, m *mocks.MovieJSON) {
				assert.Equal(t, "Бойцовский клуб", m.Name)
				assert.Equal(t, 8.8, m.Rating)
				assert.Len(t, m.Reviews, 5)
			},
		},
		// Фильм не найден
		{
			name:    "Fail - Movie Not Found",
			movieID: 999,
			mockSetup: func() {
				mockRepo.EXPECT().
					GetMovieFromRepoByID(gomock.Any(), 999).
					Return(nil, errs.ErrMovieNotFound)
			},
			expectedMovie: nil,
			expectedErr:   errs.ErrMovieNotFound,
			checkDetails:  nil,
		},
		// Проверка структуры отзывов
		{
			name:    "Check Reviews Structure",
			movieID: 0,
			mockSetup: func() {
				mockRepo.EXPECT().
					GetMovieFromRepoByID(gomock.Any(), 0).
					Return(&testMovie, nil)
			},
			expectedMovie: &testMovie,
			expectedErr:   nil,
			checkDetails: func(t *testing.T, m *mocks.MovieJSON) {
				firstReview := m.Reviews[0]
				assert.Equal(t, 1, firstReview.ID)
				assert.Equal(t, "KinoKritik77", firstReview.User.Login)
				assert.Contains(t, firstReview.ReviewText, "Абсолютный шедевр")
			},
		},
		// Проверка данных о персонале
		{
			name:    "Check Staff Data",
			movieID: 0,
			mockSetup: func() {
				mockRepo.EXPECT().
					GetMovieFromRepoByID(gomock.Any(), 0).
					Return(&testMovie, nil)
			},
			expectedMovie: &testMovie,
			expectedErr:   nil,
			checkDetails: func(t *testing.T, m *mocks.MovieJSON) {
				assert.Equal(t, "Брэд Питт", m.Staff[0].FullName)
				assert.Equal(t, "/static/img/brad_pitt.webp", m.Staff[0].Photo)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Настраиваем моки
			tt.mockSetup()

			// 2. Вызываем метод сервиса
			service := NewMovieService(mockRepo)
			movie, err := service.GetMovieByID(context.Background(), tt.movieID)

			// 3. Базовые проверки
			assert.Equal(t, tt.expectedMovie, movie)
			assert.ErrorIs(t, err, tt.expectedErr)

			// 4. Дополнительные проверки (если есть)
			if tt.checkDetails != nil && movie != nil {
				tt.checkDetails(t, movie)
			}
		})
	}
}
