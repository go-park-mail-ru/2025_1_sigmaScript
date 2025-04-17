package repository

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
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

func TestUserRepository_CreateUserPostgres(t *testing.T) {
	ctx := setupTestContext()
	repo, mock := setupTestRepo(t)
	defer repo.pgdb.Close()
	now := time.Now().Truncate(time.Microsecond)
	userToCreate := &models.User{
		Username:       "testuser",
		HashedPassword: "hashedpassword",
		Avatar:         "avatar.png",
	}

	expectedUser := models.User{
		Username:  userToCreate.Username,
		CreatedAt: now,
		UpdatedAt: now,
	}

	insertQueryRegex := regexp.QuoteMeta(insertUserQuery)

	t.Run("Success", func(t *testing.T) {
		prep := mock.ExpectPrepare(insertQueryRegex)
		prep.ExpectQuery().
			WithArgs(userToCreate.Username, userToCreate.HashedPassword, userToCreate.Avatar).
			WillReturnRows(sqlmock.NewRows([]string{"login", "created_at", "updated_at"}).
				AddRow(expectedUser.Username, expectedUser.CreatedAt, expectedUser.UpdatedAt))

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
		execErr := errors.New("unique constraint violation")
		prep := mock.ExpectPrepare(insertQueryRegex)
		prep.ExpectQuery().
			WithArgs(userToCreate.Username, userToCreate.HashedPassword, userToCreate.Avatar).
			WillReturnError(execErr)

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
			UpdatedAt: now,
		}

		prep := mock.ExpectPrepare(insertQueryRegex)

		prep.ExpectQuery().
			WithArgs(userWithInjection.Username, userWithInjection.HashedPassword, userWithInjection.Avatar).
			WillReturnRows(sqlmock.NewRows([]string{"login", "created_at", "updated_at"}).
				AddRow(expectedReturnedUser.Username, expectedReturnedUser.CreatedAt, expectedReturnedUser.UpdatedAt))

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
		CreatedAt:      now.Add(-time.Hour),
		UpdatedAt:      now,
	}

	getUserQueryRegex := regexp.QuoteMeta(getUserByUsernameQuery)

	t.Run("Success", func(t *testing.T) {
		prep := mock.ExpectPrepare(getUserQueryRegex)
		prep.ExpectQuery().
			WithArgs(login).
			WillReturnRows(sqlmock.NewRows([]string{"login", "hashed_password", "avatar", "created_at", "updated_at"}).
				AddRow(expectedUser.Username, expectedUser.HashedPassword, expectedUser.Avatar, expectedUser.CreatedAt, expectedUser.UpdatedAt))

		user, err := repo.GetUserPostgres(ctx, login)

		assert.NoError(t, err, "GetUserPostgres should not return error on success")
		assert.NotNil(t, user, "GetUserPostgres should return a user")

		assert.Equal(t, expectedUser.Username, user.Username)
		assert.Equal(t, expectedUser.HashedPassword, user.HashedPassword)
		assert.Equal(t, expectedUser.Avatar, user.Avatar)
		assert.WithinDuration(t, expectedUser.CreatedAt, user.CreatedAt, time.Second, "CreatedAt should match")
		assert.WithinDuration(t, expectedUser.UpdatedAt, user.UpdatedAt, time.Second, "UpdatedAt should match")
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

		assert.Contains(t, err.Error(), "failed to select user", "Error message should indicate select failure")
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
