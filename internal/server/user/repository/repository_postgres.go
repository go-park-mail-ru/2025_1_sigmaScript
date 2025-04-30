package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	insertUserQuery = `
		INSERT INTO "user" (login, hashed_password, avatar)
		VALUES ($1, $2, $3, $4)
		RETURNING login, created_at;
	`

	getUserByUsernameQuery = `
	SELECT login, hashed_password, avatar, created_at, updated_at
	FROM "user"
	WHERE login = $1;
`

	deleteUserByUsernameQuery = `
	DELETE FROM "user"
	WHERE login = $1
	RETURNING id;
`

	uniqueViolationCode = "23505"
)

func (r *UserRepository) CreateUserPostgres(ctx context.Context, user *models.User) error {
	// пишем логи БД
	logger := log.Ctx(ctx)

	newUser := models.User{}
	execRow, err := r.pgdb.Prepare(insertUserQuery)
	if err != nil {
		return errors.Wrapf(err, "prepare statement in CreateUserPostgres")
	}

	defer func() {
		if clErr := execRow.Close(); clErr != nil {
			logger.Error().Err(clErr).Msg("failed_to_close_statement")
		}
	}()

	err = execRow.QueryRowContext(
		ctx,
		user.Username, user.HashedPassword, user.Avatar,
	).Scan(&newUser.Username, &newUser.CreatedAt)

	if err != nil {
		errPg := fmt.Errorf("postgres: error while creating user - %w", err)
		logger.Error().Err(errPg).Msg(errs.ErrSomethingWentWrong)

		sqlErr, ok := err.(interface {
			Code() string
		})
		if ok && sqlErr.Code() == uniqueViolationCode {
			return errors.New(errs.ErrAlreadyExists)
		}
		return errors.New(errs.ErrSomethingWentWrong)
	}

	logger.Info().Msgf("successfully created new user: %s", newUser.Username)

	return nil
}

func (r *UserRepository) GetUserPostgres(ctx context.Context, login string) (*models.User, error) {
	// пишем логи БД
	logger := log.Ctx(ctx)
	var user models.User

	execRow, err := r.pgdb.Prepare(getUserByUsernameQuery)
	if err != nil {
		return nil, errors.Wrapf(err, "prepare statement in GetUserPostgres")
	}

	defer func() {
		if errClose := execRow.Close(); errClose != nil {
			logger.Error().Err(errClose).Msg(errClose.Error())
		}
	}()

	err = execRow.QueryRowContext(ctx, login).Scan(
		&user.Username,
		&user.HashedPassword,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errs.ErrIncorrectLogin)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	logger.Info().Msgf("successfully got user")
	return &user, nil
}

func (r *UserRepository) DeleteUserPostgres(ctx context.Context, login string) error {
	// пишем логи БД
	logger := log.Ctx(ctx)

	row := r.pgdb.QueryRowContext(ctx, deleteUserByUsernameQuery, login)

	var deletedID int
	err := row.Scan(&deletedID)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		if err == sql.ErrNoRows {
			return errors.New(errs.ErrIncorrectLogin)
		}
		return errors.Wrap(err, errs.ErrSomethingWentWrong)
	}

	logger.Info().Msgf("successfully deleted user")
	return nil
}

func (r *UserRepository) UpdateUserPostgres(ctx context.Context, login string, user *models.User) (*models.User, error) {
	// пишем логи БД
	logger := log.Ctx(ctx)

	if login == "" {
		errMsg := "login is empty"
		logger.Error().Msg(errMsg)
		return nil, errors.New(errs.ErrIncorrectLogin)
	}

	query := `UPDATE "user" SET `
	updatesQueryString := make([]string, 0)
	args := make([]interface{}, 0)
	argIndex := 1

	if user.Username != "" {
		updatesQueryString = append(updatesQueryString, fmt.Sprintf("login = $%d, ", argIndex))
		args = append(args, user.Username)
		argIndex++
	}
	if user.HashedPassword != "" {
		updatesQueryString = append(updatesQueryString, fmt.Sprintf("hashed_password = $%d, ", argIndex))
		args = append(args, user.HashedPassword)
		argIndex++
	}
	if user.Avatar != "" {
		updatesQueryString = append(updatesQueryString, fmt.Sprintf("avatar = $%d", argIndex))
		args = append(args, user.Avatar)
		argIndex++
	}
	if len(updatesQueryString) == 0 {
		return nil, errors.New(errs.ErrEmptyLogin)
	}

	query += fmt.Sprintf("%s, updated_at = CURRENT_TIMESTAMP WHERE login = $%d RETURNING login, avatar, created_at, updated_at",
		strings.Join(updatesQueryString, ", "), argIndex)
	args = append(args, login)

	row := r.pgdb.QueryRowContext(ctx, query, args...)
	var updatedUser models.User
	err := row.Scan(
		&updatedUser.Username,
		&updatedUser.Avatar,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to select updated user: %w", err)
	}

	logger.Info().Msgf("successfully updated user")
	return &updatedUser, nil
}
