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
WITH user_id_res AS (
	SELECT id
	FROM public."user"
	WHERE login = $1
)
INSERT INTO review (user_id, movie_id, review_text, score)
SELECT (SELECT id FROM user_id_res), $2, $3, $4;
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
	movieID int,
	newReview mocks.ReviewJSON) error {
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
		newReview.User.Login,
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

	logger.Info().Msgf("successfully updated movie rating by movie id: %d", movieID)
	return nil
}
