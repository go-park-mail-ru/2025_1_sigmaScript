package staff_repo

import (
	"context"
	"testing"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/stretchr/testify/assert"
)

func TestStaffPersonRepository_GetPersonFromRepoByID(t *testing.T) {
	db := mocks.ExistingActors
	repo := NewStaffPersonRepository(&db)

	tests := []struct {
		name        string
		personID    int
		expected    *mocks.PersonJSON
		expectedErr error
	}{
		{
			name:     "OK. Get Keanu Reeves",
			personID: 11,
			expected: func() *mocks.PersonJSON {
				person := mocks.ExistingActors[11]
				return &person
			}(),
			expectedErr: nil,
		},
		{
			name:        "Fail. Person not found",
			personID:    999,
			expected:    nil,
			expectedErr: errs.ErrPersonNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			person, err := repo.GetPersonFromRepoByID(context.Background(), tt.personID)

			assert.Equal(t, tt.expected, person)
			assert.ErrorIs(t, err, tt.expectedErr)
		})
	}
}

// func TestGetCollectionFromPostgres(t *testing.T) {
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

// 	staffRepo := NewStaffPersonPostgresRepository(pgdb)

// 	resCollections, err := staffRepo.GetPersonFromRepoByID(t.Context(), 1)
// 	assert.NoError(t, err)

// 	resByteData, err := json.Marshal(resCollections)
// 	assert.NoError(t, err)

// 	assert.NoError(t, err)
// 	assert.NotEqual(t, nil, string(resByteData), "result Collections must be not nil")
// }
