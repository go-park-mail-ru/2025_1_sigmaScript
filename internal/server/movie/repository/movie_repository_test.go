package repository

import (
	"context"
	"testing"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMovieRepository_GetMovieFromRepoByID(t *testing.T) {
	// Подготавливаем тестовые данные
	fightClub := mocks.ExistingMovies[0] // "Бойцовский клуб"
	emptyDB := make(mocks.Movies)
	fullDB := mocks.ExistingMovies // Все фильмы из моков

	tests := []struct {
		name      string
		dbSetup   mocks.Movies     // Исходное состояние "базы"
		movieID   int              // Входной параметр
		wantMovie *mocks.MovieJSON // Ожидаемый результат
		wantErr   error            // Ожидаемая ошибка
	}{
		{
			name:      "Success - Get existing movie",
			dbSetup:   fullDB,
			movieID:   0,
			wantMovie: &fightClub,
			wantErr:   nil,
		},
		{
			name:      "Fail - Movie not found",
			dbSetup:   fullDB,
			movieID:   999,
			wantMovie: nil,
			wantErr:   errs.ErrMovieNotFound,
		},
		{
			name:      "Fail - Empty database",
			dbSetup:   emptyDB,
			movieID:   0,
			wantMovie: nil,
			wantErr:   errs.ErrMovieNotFound,
		},
		{
			name:      "Fail - Negative ID",
			dbSetup:   fullDB,
			movieID:   -1,
			wantMovie: nil,
			wantErr:   errs.ErrMovieNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 1. Подготавливаем репозиторий с тестовой "БД"
			repo := NewMovieRepository(&tt.dbSetup)
			ctx := context.Background()

			// 2. Вызываем метод
			movie, err := repo.GetMovieFromRepoByID(ctx, tt.movieID)

			// 3. Проверяем результаты
			assert.Equal(t, tt.wantMovie, movie)
			assert.ErrorIs(t, err, tt.wantErr)

			// 4. Дополнительные проверки для успешных случаев
			if tt.wantMovie != nil {
				assert.Equal(t, tt.wantMovie.Name, movie.Name)
				assert.Equal(t, tt.wantMovie.Rating, movie.Rating)
			}
		})
	}
}
