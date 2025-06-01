package search_repo

import (
	_ "github.com/lib/pq"
)

// func TestGetAllGenresFromPostgres(t *testing.T) {
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
// 		UserAvatarsFullPath:   "",
// 		UserAvatarsStaticPath: "",
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

// 	genresRepo := NewGenreRepository(pgdb)
// 	log.Println(genresRepo)

// 	resCollections, err := genresRepo.GetAllGenresFromRepo(t.Context())
// 	assert.NoError(t, err)
// 	assert.NotEqual(t, nil, resCollections, "result Genres must be not nil")
// }

// func TestGetGenreByIDFromPostgres(t *testing.T) {
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
// 		UserAvatarsFullPath:   "",
// 		UserAvatarsStaticPath: "",
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

// 	genresRepo := NewGenreRepository(pgdb)
// 	log.Println(genresRepo)

// 	resCollections, err := genresRepo.GetGenreFromRepoByID(t.Context(), "4")
// 	assert.NoError(t, err)
// 	assert.NotEqual(t, nil, resCollections, "result Genres must be not nil")
// }
