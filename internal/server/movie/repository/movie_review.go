package repository

import (
	"context"
	"fmt"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	uniqueViolationCode = "23505"

	insertNewReviewQuery = `
INSERT INTO review (user_id, movie_id, review_text, score)
VALUES ($1, $2, $3, $4);
`

	updateReviewQuery = `
UPDATE review SET (review_text, score) = 
($3, $4) WHERE user_id = $1 and movie_id = $2;
`

	updateMovieRatingQuery = `
UPDATE movie
SET rating = (
    SELECT AVG(score)
    FROM review
    WHERE movie_id = $1
)
WHERE id = $1;
`
)

func (r *MoviePostgresRepository) CreateNewMovieReviewInRepo(
	ctx context.Context,
	userID string,
	movieID string,
	newReview mocks.NewReviewDataJSON) error {
	logger := log.Ctx(ctx)

	execRow, err := r.pgdb.Prepare(insertNewReviewQuery)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "prepare statement in CreateNewMovieReviewInRepo").Error())
		return errors.Wrapf(err, "prepare statement in CreateNewMovieReviewInRepo")
	}
	defer func() {
		if clErr := execRow.Close(); clErr != nil {
			logger.Error().Err(clErr).Msg("failed_to_close_statement")
		}
	}()

	_, err = execRow.Exec(
		userID,
		movieID,
		newReview.ReviewText,
		newReview.Score,
	)
	if err != nil {
		errPg := fmt.Errorf("postgres: error while creating review - %w", err)
		logger.Error().Err(errPg).Msg(errors.Wrap(errPg, errs.ErrSomethingWentWrong).Error())
		sqlErr, ok := err.(*pq.Error)
		if ok && sqlErr.Code == uniqueViolationCode {
			return errors.New(errs.ErrAlreadyExists)
		}
		return errors.New(errs.ErrSomethingWentWrong)
	}

	execRowRating, err := r.pgdb.Prepare(updateMovieRatingQuery)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "prepare statement in CreateNewMovieReviewInRepo").Error())
		return errors.Wrapf(err, "prepare statement in CreateNewMovieReviewInRepo")
	}
	defer func() {
		if clErr := execRowRating.Close(); clErr != nil {
			logger.Error().Err(clErr).Msg("failed_to_close_statement")
		}
	}()

	_, err = execRowRating.Exec(movieID)
	if err != nil {
		errPg := fmt.Errorf("postgres: error while updating movie rating - %w", err)
		logger.Error().Err(errPg).Msg(errors.Wrap(errPg, errs.ErrSomethingWentWrong).Error())
		sqlErr, ok := err.(*pq.Error)
		if ok && sqlErr.Code == uniqueViolationCode {
			return errors.New(errs.ErrAlreadyExists)
		}
		return errors.New(errs.ErrSomethingWentWrong)
	}

	logger.Info().Msgf("successfully updated movie rating by movie id: %s", movieID)
	return nil
}

func (r *MoviePostgresRepository) UpdateMovieReviewInRepo(
	ctx context.Context,
	userID string,
	movieID string,
	newReview mocks.NewReviewDataJSON) error {
	logger := log.Ctx(ctx)

	execRow, err := r.pgdb.Prepare(updateReviewQuery)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "prepare statement in UpdateMovieReviewInRepo").Error())
		return errors.Wrapf(err, "prepare statement in UpdateMovieReviewInRepo")
	}
	defer func() {
		if clErr := execRow.Close(); clErr != nil {
			logger.Error().Err(clErr).Msg("failed_to_close_statement")
		}
	}()

	result, errExec := execRow.Exec(
		userID,
		movieID,
		newReview.ReviewText,
		newReview.Score,
	)
	if errExec != nil {
		errPg := fmt.Errorf("postgres: error while updating review - %w", errExec)
		logger.Error().Err(errPg).Msg(errors.Wrap(errPg, errs.ErrSomethingWentWrong).Error())
		sqlErr, ok := errExec.(*pq.Error)
		if ok && sqlErr.Code == uniqueViolationCode {
			return errors.New(errs.ErrAlreadyExists)
		}
		return errors.New(errs.ErrSomethingWentWrong)
	}
	// check rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		errPg := fmt.Errorf("postgres: error while updating review - %w", err)
		logger.Error().Err(errPg).Msg(errors.Wrap(errPg, errs.ErrSomethingWentWrong).Error())
		return errors.New(errs.ErrSomethingWentWrong)
	}

	if rowsAffected == 0 {
		return r.CreateNewMovieReviewInRepo(ctx, userID, movieID, newReview)
	}

	// update movie rating
	execRowRating, err := r.pgdb.Prepare(updateMovieRatingQuery)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "prepare statement in UpdateMovieReviewInRepo").Error())
		return errors.Wrapf(err, "prepare statement in UpdateMovieReviewInRepo")
	}
	defer func() {
		if clErr := execRowRating.Close(); clErr != nil {
			logger.Error().Err(clErr).Msg("failed_to_close_statement")
		}
	}()

	_, err = execRowRating.Exec(movieID)
	if err != nil {
		errPg := fmt.Errorf("postgres: error while updating movie rating - %w", err)
		logger.Error().Err(errPg).Msg(errors.Wrap(errPg, errs.ErrSomethingWentWrong).Error())
		sqlErr, ok := err.(*pq.Error)
		if ok && sqlErr.Code == uniqueViolationCode {
			return errors.New(errs.ErrAlreadyExists)
		}
		return errors.New(errs.ErrSomethingWentWrong)
	}

	logger.Info().Msgf("successfully updated movie rating by movie id: %s", movieID)
	return nil
}
