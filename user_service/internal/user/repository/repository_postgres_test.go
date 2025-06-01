package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/DATA-DOG/go-sqlmock"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/internal/models"
	"github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestRepo(t *testing.T) (*UserRepository, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	require.NoError(t, err, "Failed to create sqlmock")

	repo := &UserRepository{
		pgdb: db,
	}

	return repo, mock
}

func setupTestContext() context.Context {
	disabledLogger := zerolog.Nop()
	return disabledLogger.WithContext(context.Background())
}

// Структура ошибки для имитации SQLSTATE
type SqlStateError struct {
	code    string
	message string
}

// Изменение приёма аргумента с переменной на указатель
func (se SqlStateError) Error() string {
	return se.message
}

// Дополнительно реализуем получение кода ошибки через метод Code
func (se SqlStateError) Code() string {
	return se.code
}

func TestUserRepository_CreateUserPostgres(t *testing.T) {
	ctx := setupTestContext()
	repo, mock := setupTestRepo(t)

	now := time.Now().Truncate(time.Microsecond).String()
	userToCreate := &models.User{
		Username:       "testuser",
		HashedPassword: "hashedpassword",
		Avatar:         "avatar.png",
	}

	expectedUser := models.User{
		Username:  userToCreate.Username,
		CreatedAt: now,
	}

	insertQueryRegex := regexp.QuoteMeta(insertUserQuery)

	t.Run("Success", func(t *testing.T) {
		prep := mock.ExpectPrepare(insertQueryRegex)
		prep.ExpectQuery().
			WithArgs(userToCreate.Username, userToCreate.HashedPassword, userToCreate.Avatar).
			WillReturnRows(sqlmock.NewRows([]string{"login", "created_at"}).
				AddRow(expectedUser.Username, expectedUser.CreatedAt))

		err := repo.CreateUserPostgres(ctx, userToCreate)

		assert.NoError(t, err, "CreateUserPostgres should not return an error on success")
		assert.NoError(t, mock.ExpectationsWereMet(), "Sqlmock expectations were not met")
	})

	t.Run("Prepare_Failure", func(t *testing.T) {
		prepErr := errors.New("prepare failed")
		mock.ExpectPrepare(insertQueryRegex).WillReturnError(prepErr)

		err := repo.CreateUserPostgres(ctx, userToCreate)

		assert.Error(t, err, "CreateUserPostgres should return an error on prepare failure")
		assert.Contains(t, err.Error(), "prepare statement", "Error message should indicate prepare failure")
		assert.NoError(t, mock.ExpectationsWereMet(), "Sqlmock expectations were not met")
	})

	t.Run("Exec_Failure_AlreadyExists", func(t *testing.T) {
		// execErr := errors.New("unique constraint violation")
		pqErr := pq.Error{
			Code:    uniqueViolationCode,
			Message: "this row already exists",
		}

		prep := mock.ExpectPrepare(insertQueryRegex)
		prep.ExpectQuery().
			WithArgs(userToCreate.Username, userToCreate.HashedPassword, userToCreate.Avatar).
			WillReturnError(&pqErr)

		err := repo.CreateUserPostgres(ctx, userToCreate)

		assert.Error(t, err, "CreateUserPostgres should return an error on exec failure")
		assert.EqualError(t, err, errs.ErrAlreadyExists, "Error should be ErrAlreadyExists")
		assert.NoError(t, mock.ExpectationsWereMet(), "Sqlmock expectations were not met")
	})

	// SQL_Injection_Attempt_On_Create
	t.Run("SQL_Injection_Attempt_On_Create", func(t *testing.T) {
		injectionUsername := "' OR '1'='1; --"
		userWithInjection := &models.User{
			Username:       injectionUsername,
			HashedPassword: "somepassword",
			Avatar:         "injection_avatar.png",
		}
		expectedReturnedUser := models.User{
			Username:  injectionUsername,
			CreatedAt: now,
		}

		prep := mock.ExpectPrepare(insertQueryRegex)

		prep.ExpectQuery().
			WithArgs(userWithInjection.Username, userWithInjection.HashedPassword, userWithInjection.Avatar).
			WillReturnRows(sqlmock.NewRows([]string{"login", "created_at"}).
				AddRow(expectedReturnedUser.Username, expectedReturnedUser.CreatedAt))

		err := repo.CreateUserPostgres(ctx, userWithInjection)

		assert.NoError(t, err, "CreateUserPostgres should not return an error, injection should be treated as string")
		assert.NoError(t, mock.ExpectationsWereMet(), "Sqlmock expectations were not met")
	})
}

func TestUserRepository_GetUserPostgres(t *testing.T) {
	ctx := setupTestContext()
	repo, mock := setupTestRepo(t)
	defer repo.pgdb.Close()

	login := "existinguser"
	now := time.Now().Truncate(time.Microsecond)
	expectedUser := &models.User{
		Username:       login,
		HashedPassword: "existingpassword",
		Avatar:         "existing_avatar.jpg",
		CreatedAt:      now.Add(-time.Hour).String(),
		UpdatedAt:      now.String(),
	}

	getUserQueryRegex := regexp.QuoteMeta(getUserByUsernameQuery)

	t.Run("Success", func(t *testing.T) {
		prep := mock.ExpectPrepare(getUserQueryRegex)
		prep.ExpectQuery().
			WithArgs(login).
			WillReturnRows(sqlmock.NewRows([]string{"id", "login", "hashed_password", "avatar", "created_at", "updated_at"}).
				AddRow(expectedUser.ID, expectedUser.Username, expectedUser.HashedPassword, expectedUser.Avatar, expectedUser.CreatedAt, expectedUser.UpdatedAt))

		user, err := repo.GetUserPostgres(ctx, login)

		assert.NoError(t, err, "GetUserPostgres should not return error on success")
		assert.NotNil(t, user, "GetUserPostgres should return a user")

		assert.Equal(t, expectedUser.Username, user.Username)
		assert.Equal(t, expectedUser.HashedPassword, user.HashedPassword)
		assert.Equal(t, expectedUser.Avatar, user.Avatar)
		assert.Equal(t, expectedUser.CreatedAt, user.CreatedAt, "CreatedAt should match")
		assert.Equal(t, expectedUser.UpdatedAt, user.UpdatedAt, "UpdatedAt should match")
		assert.NoError(t, mock.ExpectationsWereMet(), "Sqlmock expectations were not met")
	})

	t.Run("Prepare_Failure", func(t *testing.T) {
		prepErr := errors.New("prepare failed")
		mock.ExpectPrepare(getUserQueryRegex).WillReturnError(prepErr)

		user, err := repo.GetUserPostgres(ctx, login)

		assert.Error(t, err, "GetUserPostgres should return an error on prepare failure")
		assert.Nil(t, user, "User should be nil on prepare failure")
		assert.Contains(t, err.Error(), "prepare statement", "Error message should indicate prepare failure")
		assert.NoError(t, mock.ExpectationsWereMet(), "Sqlmock expectations were not met")
	})

	t.Run("Not_Found", func(t *testing.T) {
		nonExistentLogin := "nouser"

		prep := mock.ExpectPrepare(getUserQueryRegex)

		prep.ExpectQuery().
			WithArgs(nonExistentLogin).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserPostgres(ctx, nonExistentLogin)

		assert.Error(t, err, "GetUserPostgres should return an error when user not found")

		assert.EqualError(t, err, errs.ErrIncorrectLogin, "Error should be ErrIncorrectLogin")
		assert.Nil(t, user, "User should be nil when not found")
		assert.NoError(t, mock.ExpectationsWereMet(), "Sqlmock expectations were not met")
	})

	t.Run("Query_Scan_Failure", func(t *testing.T) {
		scanErr := errors.New("scan error")

		prep := mock.ExpectPrepare(getUserQueryRegex)

		prep.ExpectQuery().
			WithArgs(login).
			WillReturnError(scanErr)

		user, err := repo.GetUserPostgres(ctx, login)

		assert.Error(t, err, "GetUserPostgres should return an error on scan failure")
		assert.Nil(t, user, "User should be nil on scan failure")

		assert.Contains(t, err.Error(), "failed to get user", "Error message should indicate select failure")
		assert.ErrorIs(t, err, scanErr, "Original scan error should be wrapped")
		assert.NoError(t, mock.ExpectationsWereMet(), "Sqlmock expectations were not met")
	})

	// SQL_Injection_Attempt_On_Get
	t.Run("SQL_Injection_Attempt_On_Get", func(t *testing.T) {
		injectionLogin := "' OR '1'='1"
		prep := mock.ExpectPrepare(getUserQueryRegex)

		prep.ExpectQuery().
			WithArgs(injectionLogin).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserPostgres(ctx, injectionLogin)

		assert.Error(t, err, "GetUserPostgres should return an error on SQL injection attempt (treated as not found)")
		assert.EqualError(t, err, errs.ErrIncorrectLogin, "Error should be ErrIncorrectLogin, indicating user not found")
		assert.Nil(t, user, "User should be nil")
		assert.NoError(t, mock.ExpectationsWereMet(), "Sqlmock expectations were not met")
	})
}

func TestUserRepository_DeleteUserPostgres(t *testing.T) {
	ctx := setupTestContext()
	repo, mock := setupTestRepo(t)
	defer repo.pgdb.Close()

	loginToDelete := "usertodelete"
	expectedDeletedID := 123

	deleteQueryRegex := regexp.QuoteMeta(deleteUserByUsernameQuery)

	t.Run("Success", func(t *testing.T) {
		mock.ExpectQuery(deleteQueryRegex).
			WithArgs(loginToDelete).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedDeletedID))

		err := repo.DeleteUserPostgres(ctx, loginToDelete)

		assert.NoError(t, err, "DeleteUserPostgres should not return error on success")
		assert.NoError(t, mock.ExpectationsWereMet(), "Sqlmock expectations were not met")
	})

	t.Run("Not_Found", func(t *testing.T) {
		nonExistentLogin := "nonexistent"
		mock.ExpectQuery(deleteQueryRegex).
			WithArgs(nonExistentLogin).
			WillReturnError(sql.ErrNoRows)

		err := repo.DeleteUserPostgres(ctx, nonExistentLogin)

		assert.Error(t, err, "DeleteUserPostgres should return error when user not found")
		assert.EqualError(t, err, errs.ErrIncorrectLogin, "Error should be ErrIncorrectLogin")
		assert.NoError(t, mock.ExpectationsWereMet(), "Sqlmock expectations were not met")
	})

	t.Run("Query_Scan_Failure", func(t *testing.T) {
		scanErr := errors.New("db connection lost")
		mock.ExpectQuery(deleteQueryRegex).
			WithArgs(loginToDelete).
			WillReturnError(scanErr)

		err := repo.DeleteUserPostgres(ctx, loginToDelete)

		assert.Error(t, err, "DeleteUserPostgres should return error on query/scan failure")
		assert.Contains(t, err.Error(), errs.ErrSomethingWentWrong, "Error message should indicate something went wrong")
		assert.ErrorIs(t, err, scanErr, "Original db error should be wrapped")
		assert.NoError(t, mock.ExpectationsWereMet(), "Sqlmock expectations were not met")
	})

	t.Run("SQL_Injection_Attempt_On_Delete", func(t *testing.T) {
		injectionLogin := "'; DROP TABLE public.\"user\"; -- "

		mock.ExpectQuery(deleteQueryRegex).
			WithArgs(injectionLogin).
			WillReturnError(sql.ErrNoRows)

		err := repo.DeleteUserPostgres(ctx, injectionLogin)

		assert.Error(t, err, "DeleteUserPostgres should return an error on SQL injection attempt (treated as not found)")
		assert.EqualError(t, err, errs.ErrIncorrectLogin, "Error should be ErrIncorrectLogin, indicating user not found")
		assert.NoError(t, mock.ExpectationsWereMet(), "Sqlmock expectations were not met")
	})
}

func TestUserRepository_UpdateUserPostgres(t *testing.T) {
	ctx := setupTestContext()
	repo, mock := setupTestRepo(t)
	defer repo.pgdb.Close()

	targetUsername := "user_to_update"
	now := time.Now().Truncate(time.Microsecond)
	createdAt := now.Add(-2 * time.Hour)

	// Success_Update_Avatar
	t.Run("Success_Update_Avatar", func(t *testing.T) {
		login := targetUsername
		userToUpdate := &models.User{
			Avatar: "new_avatar.webp",
		}
		expectedUpdatedUser := &models.User{
			Username:  targetUsername,
			Avatar:    userToUpdate.Avatar,
			CreatedAt: createdAt.String(),
			UpdatedAt: now.String(),
		}

		expectedQuery := `UPDATE "user" SET avatar = $1, updated_at = CURRENT_TIMESTAMP WHERE login = $2 RETURNING login, avatar, created_at`
		mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WithArgs(userToUpdate.Avatar, login).
			WillReturnRows(sqlmock.NewRows([]string{"login", "avatar", "created_at"}).
				AddRow(expectedUpdatedUser.Username, expectedUpdatedUser.Avatar, expectedUpdatedUser.CreatedAt))

		updatedUser, err := repo.UpdateUserPostgres(ctx, login, userToUpdate)

		assert.NoError(t, err)
		assert.NotNil(t, updatedUser)
		assert.Equal(t, expectedUpdatedUser.Username, updatedUser.Username)
		assert.Equal(t, expectedUpdatedUser.Avatar, updatedUser.Avatar)
		assert.Equal(t, expectedUpdatedUser.CreatedAt, updatedUser.CreatedAt)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// No_Fields_To_Update
	t.Run("No_Fields_To_Update", func(t *testing.T) {
		login := targetUsername
		userToUpdate := &models.User{}

		updatedUser, err := repo.UpdateUserPostgres(ctx, login, userToUpdate)

		assert.Error(t, err, "Should return error if no fields are provided for update")
		assert.EqualError(t, err, errs.ErrEmptyLogin)
		assert.Nil(t, updatedUser)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("User_Not_Found_For_Update", func(t *testing.T) {
		login := "ghost_user"
		nonExistentUser := &models.User{
			Avatar: "some_avatar.png",
		}

		expectedQuery := `UPDATE "user" SET avatar = $1, updated_at = CURRENT_TIMESTAMP WHERE login = $2 RETURNING login, avatar, created_at`
		mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WithArgs(nonExistentUser.Avatar, login).
			WillReturnError(sql.ErrNoRows)

		updatedUser, err := repo.UpdateUserPostgres(ctx, login, nonExistentUser)

		assert.Error(t, err, "Should return error when user to update is not found")
		assert.Nil(t, updatedUser)
		assert.ErrorIs(t, err, sql.ErrNoRows, "Error should wrap sql.ErrNoRows")
		assert.Contains(t, err.Error(), "failed to select updated user", "Error message should indicate scan failure")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// SQL_Injection_Attempt_In_Updated_Field
	t.Run("SQL_Injection_Attempt_In_Updated_Field", func(t *testing.T) {
		injectionString := "'; DROP TABLE users; --"

		login := targetUsername
		userWithInjection := &models.User{
			Avatar: injectionString,
		}
		expectedInjectedUser := &models.User{
			Username:  targetUsername,
			Avatar:    injectionString,
			CreatedAt: createdAt.String(),
			UpdatedAt: now.String(),
		}

		expectedQuery := `UPDATE "user" SET avatar = $1, updated_at = CURRENT_TIMESTAMP WHERE login = $2 RETURNING login, avatar, created_at`
		mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WithArgs(injectionString, login).
			WillReturnRows(sqlmock.NewRows([]string{"login", "avatar", "created_at"}).
				AddRow(expectedInjectedUser.Username, expectedInjectedUser.Avatar, expectedInjectedUser.CreatedAt))

		updatedUser, err := repo.UpdateUserPostgres(ctx, login, userWithInjection)

		assert.NoError(t, err, "UpdateUserPostgres should not return an error, injection should be treated as string")
		assert.NotNil(t, updatedUser)
		assert.Equal(t, expectedInjectedUser.Username, updatedUser.Username)
		assert.Equal(t, injectionString, updatedUser.Avatar, "Avatar field should contain the literal injection string")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// SQL_Injection_Attempt_In_Where
	t.Run("SQL_Injection_Attempt_In_Where", func(t *testing.T) {
		injectionUsername := "' OR '1'='1"

		login := injectionUsername
		userWithInjection := &models.User{
			Avatar: "update_attempt.jpg",
		}

		expectedQuery := `UPDATE "user" SET avatar = $1, updated_at = CURRENT_TIMESTAMP WHERE login = $2 RETURNING login, avatar, created_at`
		mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
			WithArgs(userWithInjection.Avatar, injectionUsername).
			WillReturnError(sql.ErrNoRows)

		updatedUser, err := repo.UpdateUserPostgres(ctx, login, userWithInjection)

		assert.Error(t, err, "UpdateUserPostgres should return an error on SQL injection attempt in WHERE")
		assert.Nil(t, updatedUser)
		assert.ErrorIs(t, err, sql.ErrNoRows, "Error should wrap sql.ErrNoRows")
		assert.Contains(t, err.Error(), "failed to select updated user", "Error message should indicate scan failure")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	// Empty_Login_Select
	t.Run("Empty_Login_Select", func(t *testing.T) {
		login := ""
		userWithNoLogin := &models.User{
			Avatar: "update_attempt.jpg",
		}

		updatedUser, err := repo.UpdateUserPostgres(ctx, login, userWithNoLogin)

		assert.Error(t, err, "UpdateUserPostgres should return an error because login is empty")
		assert.Nil(t, updatedUser)
		assert.Contains(t, err.Error(), errs.ErrIncorrectLogin, "Error message should indicate no login")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

// INTEGRATIONAL TESTS

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

// 	staffRepo := NewUserRepository(pgdb)

// 	resCollections, err := staffRepo.GetUserProfilePostgres(t.Context(), "KinoLooker")
// 	assert.NoError(t, err)

// 	resByteData, err := json.Marshal(resCollections)
// 	assert.NoError(t, err)

// 	assert.NoError(t, err)
// 	assert.NotEqual(t, nil, string(resByteData), "result Collections must be not nil")
// }

// func TestAddActorToFavoritesPostgres(t *testing.T) {
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
// 	staffRepo := NewUserRepository(pgdb)
// 	err = staffRepo.AddFavoriteActor(t.Context(), "KinoLooker", "7")
// 	assert.NoError(t, err)
// }

// func TestRemoveActorFromFavoritesPostgres(t *testing.T) {
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
// 	staffRepo := NewUserRepository(pgdb)
// 	err = staffRepo.RemoveFavoriteActor(t.Context(), "KinoLooker", "7")
// 	assert.NoError(t, err)
// }

// func TestAddMovieToFavoritesPostgres(t *testing.T) {
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
// 	staffRepo := NewUserRepository(pgdb)
// 	err = staffRepo.AddFavoriteMovie(t.Context(), "KinoLooker", "7")
// 	assert.NoError(t, err)
// }

// func TestRemoveMovieFromFavoritesPostgres(t *testing.T) {
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
// 	staffRepo := NewUserRepository(pgdb)
// 	err = staffRepo.RemoveFavoriteMovie(t.Context(), "KinoLooker", "7")
// 	assert.NoError(t, err)
// }
