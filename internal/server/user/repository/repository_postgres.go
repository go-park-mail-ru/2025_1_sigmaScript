package repository

import (
	"context"
	"database/sql"
	"fmt"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	insertUserQuery = `
		INSERT INTO "user" (login, hashed_password, avatar)
		VALUES ($1, $2, $3, $4)
		RETURNING login, created_at, updated_at;
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
	).Scan(&newUser.Username, &newUser.CreatedAt, &newUser.UpdatedAt)

	if err != nil {
		errPg := fmt.Errorf("postgres: error while creating user - %w", err)
		logger.Error().Err(errPg).Msg(errs.ErrSomethingWentWrong)
		return errors.New(errs.ErrAlreadyExists)
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
		return nil, fmt.Errorf("failed to select user: %w", err)
	}
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

	return nil
}
