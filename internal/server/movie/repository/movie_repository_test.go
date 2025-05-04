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

// func TestGetMovieFromPostgres(t *testing.T) {
// 	// TODO fix config: it`s test database test password
// 	postgres := config.Postgres{
// 		Host:            "127.0.0.1",
// 		Port:            5433,
// 		User:            "filmlk_user",
// 		Password:        "filmlk_password",
// 		Name:            "filmlk",
// 		MaxOpenConns:    100,
// 		MaxIdleConns:    30,
// 		ConnMaxLifetime: 30,
// 		ConnMaxIdleTime: 5,
// 	}

// 	avatarLocalStorage := config.LocalAvatarsStorage{
// 		UserAvatarsFullPath:     "",
// 		UserAvatarsRelativePath: "",
// 	}

// 	pgDatabase := config.Databases{
// 		Postgres:     postgres,
// 		LocalStorage: avatarLocalStorage,
// 	}

// 	pgListener := config.Listener{
// 		Port: "5433",
// 	}

// 	cfgDB := config.ConfigPgDB{
// 		Listener:  pgListener,
// 		Databases: pgDatabase,
// 	}

// 	ctxDb := config.WrapPgDatabaseContext(context.Background(), cfgDB)
// 	ctxDb, cancel := context.WithTimeout(ctxDb, time.Second*30)
// 	defer cancel()

// 	pgdb, err := db.SetupDatabase(ctxDb, cancel)
// 	assert.NoError(t, err)

// 	movieRepo := NewMoviePostgresRepository(pgdb)

// 	resCollections, err := movieRepo.GetMovieFromRepoByID(t.Context(), 2)
// 	assert.NoError(t, err)

// 	resByteData, err := json.Marshal(resCollections)
// 	assert.NoError(t, err)

// 	assert.NoError(t, err)
// 	assert.NotEqual(t, nil, string(resByteData), "result Collections must be not nil")
// }

// func TestMovieReviewsFromPostgres(t *testing.T) {
// 	// TODO fix config: it`s test database test password
// 	postgres := config.Postgres{
// 		Host:            "127.0.0.1",
// 		Port:            5433,
// 		User:            "filmlk_user",
// 		Password:        "filmlk_password",
// 		Name:            "filmlk",
// 		MaxOpenConns:    100,
// 		MaxIdleConns:    30,
// 		ConnMaxLifetime: 30,
// 		ConnMaxIdleTime: 5,
// 	}

// 	avatarLocalStorage := config.LocalAvatarsStorage{
// 		UserAvatarsFullPath:     "",
// 		UserAvatarsRelativePath: "",
// 	}

// 	pgDatabase := config.Databases{
// 		Postgres:     postgres,
// 		LocalStorage: avatarLocalStorage,
// 	}

// 	pgListener := config.Listener{
// 		Port: "5433",
// 	}

// 	cfgDB := config.ConfigPgDB{
// 		Listener:  pgListener,
// 		Databases: pgDatabase,
// 	}

// 	ctxDb := config.WrapPgDatabaseContext(context.Background(), cfgDB)
// 	ctxDb, cancel := context.WithTimeout(ctxDb, time.Second*30)
// 	defer cancel()

// 	pgdb, err := db.SetupDatabase(ctxDb, cancel)
// 	assert.NoError(t, err)

// 	movieRepo := NewMoviePostgresRepository(pgdb)

// 	resCollections, err := movieRepo.GetAllReviewsOfMovieFromRepoByID(t.Context(), 1)
// 	assert.NoError(t, err)

// 	resByteData, err := json.Marshal(resCollections)
// 	assert.NoError(t, err)

// 	assert.NoError(t, err)
// 	assert.NotEqual(t, nil, string(resByteData), "result Collections must be not nil")
// }

// func TestCreateMovieReviewInPostgres(t *testing.T) {
// 	// TODO fix config: it`s test database test password
// 	postgres := config.Postgres{
// 		Host:            "127.0.0.1",
// 		Port:            5433,
// 		User:            "filmlk_user",
// 		Password:        "filmlk_password",
// 		Name:            "filmlk",
// 		MaxOpenConns:    100,
// 		MaxIdleConns:    30,
// 		ConnMaxLifetime: 30,
// 		ConnMaxIdleTime: 5,
// 	}

// 	avatarLocalStorage := config.LocalAvatarsStorage{
// 		UserAvatarsFullPath:     "",
// 		UserAvatarsRelativePath: "",
// 	}

// 	pgDatabase := config.Databases{
// 		Postgres:     postgres,
// 		LocalStorage: avatarLocalStorage,
// 	}

// 	pgListener := config.Listener{
// 		Port: "5433",
// 	}

// 	cfgDB := config.ConfigPgDB{
// 		Listener:  pgListener,
// 		Databases: pgDatabase,
// 	}

// 	ctxDb := config.WrapPgDatabaseContext(context.Background(), cfgDB)
// 	ctxDb, cancel := context.WithTimeout(ctxDb, time.Second*30)
// 	defer cancel()

// 	pgdb, err := db.SetupDatabase(ctxDb, cancel)
// 	assert.NoError(t, err)

// 	movieRepo := NewMoviePostgresRepository(pgdb)

// 	newReview := mocks.ReviewJSON{
// 		ID:         -1,
// 		Score:      8.0,
// 		ReviewText: "blo",
// 		CreatedAt:  time.Now().String(),
// 		User: mocks.ReviewUserDataJSON{
// 			Login:  "KinoLooker",
// 			Avatar: "",
// 		},
// 	}

// 	err = movieRepo.CreateNewMovieReviewInRepo(t.Context(), 4, newReview)
// 	assert.NoError(t, err)

// }
