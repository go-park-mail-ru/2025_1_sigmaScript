package repository

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCollectionRepository_GetMainPageCollectionsFromRepo(t *testing.T) {
	tests := []struct {
		name        string
		dbSetup     mocks.Collections
		expected    mocks.Collections
		expectedErr error
	}{
		{
			name:        "OK. Get collections",
			dbSetup:     mocks.MainPageCollections,
			expected:    mocks.MainPageCollections,
			expectedErr: nil,
		},
		{
			name:        "OK. Empty collections",
			dbSetup:     make(mocks.Collections),
			expected:    make(mocks.Collections),
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewCollectionRepository(&tt.dbSetup)
			collections, err := repo.GetMainPageCollectionsFromRepo(context.Background())

			assert.Equal(t, tt.expected, collections)
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

// 	collectionRepo := NewCollectionPostgresRepository(pgdb)
// 	log.Println(collectionRepo)

// 	resCollections, err := collectionRepo.GetMainPageCollectionsFromRepo(t.Context())
// 	assert.NoError(t, err)
// 	assert.NotEqual(t, nil, resCollections, "result Collections must be not nil")
// }
