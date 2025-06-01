package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/internal/models"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	uniqueViolationCode = "23505"

	insertUserQuery = `
	INSERT INTO "user" (login, hashed_password, avatar)
	VALUES ($1, $2, $3)
	RETURNING login, created_at;
`

	getUserByUsernameQuery = `
	SELECT id, login, hashed_password, avatar, created_at, updated_at
	FROM "user"
	WHERE login = $1;
`

	deleteUserByUsernameQuery = `
	DELETE FROM "user"
	WHERE login = $1
	RETURNING id;
`
	getFavoriteUserActorsQuery = `
	SELECT
		u.id as user_id,
		p.id as id,
		p.full_name,
		p.photo
	FROM "user" u
	LEFT JOIN user_person_favorite upf ON u.id = upf.user_id
	LEFT JOIN person p ON p.id = upf.person_id
	WHERE u.id = $1;
`

	getFavoriteUserMoviesQuery = `
	SELECT
		u.id as user_id,
		m.id as id,
		m.name as title,
		m.poster as preview_url
	FROM "user" u
	LEFT JOIN user_movie_favorite umf ON u.id = umf.user_id
	LEFT JOIN movie m ON m.id = umf.movie_id
	WHERE u.id = $1;
`

	getUserReviewsQuery = `
SELECT
	r.id,
	r.score,
	r.review_text,
	r.updated_at as created_at,
	u.login,
	u.avatar
FROM "user" u
LEFT JOIN review r ON r.user_id = u.id
WHERE u.id = $1
ORDER BY u.id;
`

	addFavoriteActorQuery = `
	WITH user_id_res AS (
		SELECT id
		FROM public."user"
		WHERE login = $1
	)
	insert into user_person_favorite (user_id, person_id)
	SELECT (SELECT id FROM user_id_res), $2;
`

	addFavoriteMovieQuery = `
	WITH user_id_res AS (
		SELECT id
		FROM public."user"
		WHERE login = $1
	)
	insert into user_movie_favorite (user_id, movie_id)
	SELECT (SELECT id FROM user_id_res), $2;
`

	removeFavoriteActorQuery = `
WITH user_id_res AS (
	SELECT id
	FROM public."user"
	WHERE login = $1
)
delete from user_person_favorite
WHERE user_id = (SELECT id FROM user_id_res) AND person_id = $2;
`

	removeFavoriteMovieQuery = `
WITH user_id_res AS (
	SELECT id
	FROM public."user"
	WHERE login = $1
)
delete from user_movie_favorite
WHERE user_id = (SELECT id FROM user_id_res) AND movie_id = $2;
`
)

func (r *UserRepository) CreateUserPostgres(ctx context.Context, user *models.User) error {
	// пишем логи БД
	logger := log.Ctx(ctx)

	newUser := models.User{}
	execRow, err := r.pgdb.Prepare(insertUserQuery)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "prepare statement in CreateUserPostgres").Error())
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
		logger.Error().Err(errPg).Msg(errors.Wrap(errPg, errs.ErrSomethingWentWrong).Error())

		sqlErr, ok := err.(*pq.Error)
		if ok && sqlErr.Code == uniqueViolationCode {
			return errors.New(errs.ErrAlreadyExists)
		}
		return errors.Wrap(err, errs.ErrSomethingWentWrong)
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
		logger.Error().Err(err).Msg(errors.Wrapf(err, "prepare statement in GetUserPostgres").Error())
		return nil, errors.Wrapf(err, "prepare statement in GetUserPostgres")
	}

	defer func() {
		if errClose := execRow.Close(); errClose != nil {
			logger.Error().Err(errClose).Msg(errClose.Error())
		}
	}()

	err = execRow.QueryRowContext(ctx, login).Scan(
		&user.ID,
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
		updatesQueryString = append(updatesQueryString, fmt.Sprintf("login = $%d", argIndex))
		args = append(args, user.Username)
		argIndex++
	}
	if user.HashedPassword != "" {
		updatesQueryString = append(updatesQueryString, fmt.Sprintf("hashed_password = $%d", argIndex))
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

	query += fmt.Sprintf("%s, updated_at = CURRENT_TIMESTAMP WHERE login = $%d RETURNING login, avatar, created_at",
		strings.Join(updatesQueryString, ", "), argIndex)
	args = append(args, login)

	row := r.pgdb.QueryRow(query, args...)
	var updatedUser models.User
	err := row.Scan(
		&updatedUser.Username,
		&updatedUser.Avatar,
		&updatedUser.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to select updated user: %w", err)
	}

	logger.Info().Msgf("successfully updated user")
	return &updatedUser, nil
}

func (r *UserRepository) GetUserProfilePostgres(ctx context.Context, login string) (*models.Profile, error) {
	logger := log.Ctx(ctx)

	user, err := r.GetUserPostgres(ctx, login)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "prepare statement in GetUserPostgres").Error())
		return nil, errors.Wrapf(err, "prepare statement in GetUserPostgres")
	}

	// get staff
	resStaff := []models.PersonJSON{}
	execRowStaff, err := r.pgdb.Query(getFavoriteUserActorsQuery, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "error in query statement in GetMovieFromRepoByID").Error())
		return nil, errors.Wrap(err, "error in prepare query statement in GetMovieFromRepoByID")
	}
	defer func() {
		if closeErr := execRowStaff.Close(); closeErr != nil {
			logger.Error().Err(closeErr).Msg("failed_to_close_statement")
			return
		}
	}()

	for execRowStaff.Next() {
		var userID string

		var personID sql.NullInt64
		var personFullName sql.NullString
		var personPhoto sql.NullString

		if err := execRowStaff.Scan(
			&userID,
			&personID,
			&personFullName,
			&personPhoto,
		); err != nil {
			errMsg := errors.Wrap(err, "error in query scan in GetMovieFromRepoByID")
			logger.Error().Err(errMsg).Msg(errMsg.Error())
			return nil, errMsg
		}

		// skip if not found
		if !personID.Valid {
			continue
		}

		person := models.PersonJSON{
			ID:       int(personID.Int64),
			FullName: personFullName.String,
			Photo:    personPhoto.String,
		}

		resStaff = append(resStaff, person)
	}
	if execErr := execRowStaff.Err(); execErr != nil {
		errMsg := errors.Wrap(execErr, "error in query next in GetMovieFromRepoByID")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		return nil, errMsg
	}

	// get favorite movies
	resMovies := []models.Movie{}
	execRowMovie, err := r.pgdb.Query(getFavoriteUserMoviesQuery, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "error in query statement in GetMovieFromRepoByID").Error())
		return nil, errors.Wrap(err, "error in prepare query statement in GetMovieFromRepoByID")
	}
	defer func() {
		if closeErr := execRowMovie.Close(); closeErr != nil {
			logger.Error().Err(closeErr).Msg("failed_to_close_statement")
			return
		}
	}()

	for execRowMovie.Next() {
		var userID string

		var movieID sql.NullInt64
		var movieTitle sql.NullString
		var moViePreviewURL sql.NullString

		if err := execRowMovie.Scan(
			&userID,
			&movieID,
			&movieTitle,
			&moViePreviewURL,
		); err != nil {
			errMsg := errors.Wrap(err, "error in query scan in GetMovieFromRepoByID")
			logger.Error().Err(errMsg).Msg(errMsg.Error())
			return nil, errMsg
		}

		// skip if not found
		if !movieID.Valid {
			continue
		}

		movie := models.Movie{
			ID:         int(movieID.Int64),
			Title:      movieTitle.String,
			PreviewURL: moViePreviewURL.String,
		}

		resMovies = append(resMovies, movie)
	}
	if execErr := execRowMovie.Err(); execErr != nil {
		errMsg := errors.Wrap(execErr, "error in query next in GetMovieFromRepoByID")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		return nil, errMsg
	}

	// get reviews
	resReviews := []models.ReviewJSON{}
	execRowReviews, err := r.pgdb.Query(getUserReviewsQuery, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "error in query statement in GetMovieFromRepoByID").Error())
		return nil, errors.Wrap(err, "error in prepare query statement in GetMovieFromRepoByID")
	}
	defer func() {
		if closeErr := execRowReviews.Close(); closeErr != nil {
			logger.Error().Err(closeErr).Msg("failed_to_close_statement")
			return
		}
	}()

	for execRowReviews.Next() {
		var reviewID sql.NullInt64
		var reviewScore sql.NullFloat64
		var reviewText sql.NullString
		var reviewCreatedAt sql.NullString

		var userLogin sql.NullString
		var userAvatar sql.NullString

		if err := execRowReviews.Scan(
			&reviewID,
			&reviewScore,
			&reviewText,
			&reviewCreatedAt,
			&userLogin,
			&userAvatar,
		); err != nil {
			errMsg := errors.Wrap(err, "error in query scan in GetMovieFromRepoByID")
			logger.Error().Err(errMsg).Msg(errMsg.Error())
			return nil, errMsg
		}

		// skip if bad user or review
		if !reviewID.Valid || !userLogin.Valid {
			continue
		}
		reviewUser := models.ReviewUserDataJSON{
			Login:  userLogin.String,
			Avatar: userAvatar.String,
		}

		review := models.ReviewJSON{
			ID:         int(reviewID.Int64),
			Score:      reviewScore.Float64,
			ReviewText: reviewText.String,
			CreatedAt:  reviewCreatedAt.String,
			User:       reviewUser,
		}

		resReviews = append(resReviews, review)
	}
	if execErr := execRowReviews.Err(); execErr != nil {
		errMsg := errors.Wrap(execErr, "error in query next in GetMovieFromRepoByID")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		return nil, errMsg
	}

	// result
	result := models.Profile{
		Username:        user.Username,
		Avatar:          user.Avatar,
		CreatedAt:       user.CreatedAt,
		MovieCollection: resMovies,
		Actors:          resStaff,
		Reviews:         resReviews,
	}
	return &result, nil
}

func (r *UserRepository) AddFavoriteMovie(ctx context.Context, login string, movieID string) error {
	logger := log.Ctx(ctx)

	execRow, err := r.pgdb.Prepare(addFavoriteMovieQuery)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "prepare statement in AddFavoriteMovie").Error())
		return errors.Wrapf(err, "prepare statement in AddFavoriteMovie")
	}

	defer func() {
		if clErr := execRow.Close(); clErr != nil {
			logger.Error().Err(clErr).Msg("failed_to_close_statement")
		}
	}()

	_, err = execRow.Exec(
		login,
		movieID,
	)
	if err != nil {
		errPg := fmt.Errorf("postgres: error while adding favorite movie - %w", err)
		logger.Error().Err(errPg).Msg(errors.Wrap(errPg, errs.ErrSomethingWentWrong).Error())

		sqlErr, ok := err.(*pq.Error)
		if ok && sqlErr.Code == uniqueViolationCode {
			return errors.New(errs.ErrAlreadyExists)
		}
		return errors.Wrap(err, errs.ErrSomethingWentWrong)
	}

	logger.Info().Msgf("successfully added movie to favorites by id: %s", movieID)
	return nil
}

func (r *UserRepository) AddFavoriteActor(ctx context.Context, login string, actorID string) error {
	logger := log.Ctx(ctx)

	execRow, err := r.pgdb.Prepare(addFavoriteActorQuery)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "prepare statement in AddFavoriteActor").Error())
		return errors.Wrapf(err, "prepare statement in AddFavoriteActor")
	}

	defer func() {
		if clErr := execRow.Close(); clErr != nil {
			logger.Error().Err(clErr).Msg("failed_to_close_statement")
		}
	}()

	_, err = execRow.Exec(
		login,
		actorID,
	)
	if err != nil {
		errPg := fmt.Errorf("postgres: error while adding favorite person - %w", err)
		logger.Error().Err(errPg).Msg(errors.Wrap(errPg, errs.ErrSomethingWentWrong).Error())

		sqlErr, ok := err.(*pq.Error)
		if ok && sqlErr.Code == uniqueViolationCode {
			return errors.New(errs.ErrAlreadyExists)
		}
		return errors.Wrap(err, errs.ErrSomethingWentWrong)
	}

	logger.Info().Msgf("successfully added person to favorites by id: %s", actorID)
	return nil
}

func (r *UserRepository) RemoveFavoriteMovie(ctx context.Context, login string, movieID string) error {
	logger := log.Ctx(ctx)

	execRow, err := r.pgdb.Prepare(removeFavoriteMovieQuery)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "prepare statement in RemoveFavoriteMovie").Error())
		return errors.Wrapf(err, "prepare statement in RemoveFavoriteMovie")
	}

	defer func() {
		if clErr := execRow.Close(); clErr != nil {
			logger.Error().Err(clErr).Msg("failed_to_close_statement")
		}
	}()

	_, err = execRow.Exec(
		login,
		movieID,
	)
	if err != nil {
		errPg := fmt.Errorf("postgres: error while removing favorite movie - %w", err)
		logger.Error().Err(errPg).Msg(errors.Wrap(errPg, errs.ErrSomethingWentWrong).Error())

		sqlErr, ok := err.(*pq.Error)
		if ok && sqlErr.Code == uniqueViolationCode {
			return errors.New(errs.ErrAlreadyExists)
		}
		return errors.Wrap(err, errs.ErrSomethingWentWrong)
	}

	logger.Info().Msgf("successfully removed movie to favorites by id: %s", movieID)
	return nil
}

func (r *UserRepository) RemoveFavoriteActor(ctx context.Context, login string, actorID string) error {
	logger := log.Ctx(ctx)

	execRow, err := r.pgdb.Prepare(removeFavoriteActorQuery)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "prepare statement in RemoveFavoriteActor").Error())
		return errors.Wrapf(err, "prepare statement in RemoveFavoriteActor")
	}

	defer func() {
		if clErr := execRow.Close(); clErr != nil {
			logger.Error().Err(clErr).Msg("failed_to_close_statement")
		}
	}()

	_, err = execRow.Exec(
		login,
		actorID,
	)
	if err != nil {
		errPg := fmt.Errorf("postgres: error while removing favorite person - %w", err)
		logger.Error().Err(errPg).Msg(errors.Wrap(errPg, errs.ErrSomethingWentWrong).Error())

		sqlErr, ok := err.(*pq.Error)
		if ok && sqlErr.Code == uniqueViolationCode {
			return errors.New(errs.ErrAlreadyExists)
		}
		return errors.Wrap(err, errs.ErrSomethingWentWrong)
	}

	logger.Info().Msgf("successfully removed person to favorites by id: %s", actorID)
	return nil
}
