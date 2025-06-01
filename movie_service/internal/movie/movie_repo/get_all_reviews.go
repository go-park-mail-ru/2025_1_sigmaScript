package movie_repo

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func (r *MoviePostgresRepository) GetAllReviewsOfMovieFromRepoByID(ctx context.Context, movieID int) (*[]mocks.ReviewJSON, error) {
	logger := log.Ctx(ctx)

	resReviews := []mocks.ReviewJSON{}
	execRowReviews, err := r.pgdb.Query(getMovieReviewsQuery, movieID)
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
		var movieName string

		var reviewID int
		var reviewScore sql.NullString
		var reviewText sql.NullString
		var reviewCreatedAt sql.NullString

		var userLogin string
		var userAvatar sql.NullString

		if err := execRowReviews.Scan(
			&movieName,
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

		reviewUser := mocks.ReviewUserDataJSON{
			Login:  userLogin,
			Avatar: userAvatar.String,
		}

		reviewScoreFloat, err := strconv.ParseFloat(reviewScore.String, 64)
		if err != nil {
			errMsg := errors.Wrap(err, "error in query statement in GetMovieFromRepoByID")
			logger.Error().Err(errMsg).Msg(errMsg.Error())
			return nil, errMsg
		}

		review := mocks.ReviewJSON{
			ID:         reviewID,
			Score:      reviewScoreFloat,
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

	return &resReviews, nil
}
