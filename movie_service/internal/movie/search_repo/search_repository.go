package search_repo

import (
	"context"
	"database/sql"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	searchActorsBySearchStringQuery = `
SELECT
    p.id as id,
    p.full_name as full_name,
    p.photo as photo
FROM person p
WHERE
    p.full_name ILIKE '%' || $1 || '%'
    OR p.en_full_name ILIKE '%' || $1 || '%'
    OR p.about ILIKE '%' || $1 || '%'
ORDER BY p.id;
`

	searchMoviesBySearchStringQuery = `
SELECT
    m.id as id,
    m.name as title,
    m.poster as preview_url
FROM movie m
WHERE
    m.name ILIKE '%' || $1 || '%'
    OR m.original_name ILIKE '%' || $1 || '%'
    OR m.about ILIKE '%' || $1 || '%'
ORDER BY m.id;
`
)

type SearchRepository struct {
	pgdb *sql.DB
}

func NewSearchRepository(db *sql.DB) *SearchRepository {
	return &SearchRepository{
		pgdb: db,
	}
}

func (r *SearchRepository) SearchActorsAndMovies(ctx context.Context, searchStr string) (*models.SearchResponseJSON, error) {
	logger := log.Ctx(ctx)

	resStaff := []mocks.PersonJSON{}
	execRowStaff, err := r.pgdb.Query(searchActorsBySearchStringQuery, searchStr)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "error in query statement in SearchActorsAndMovies").Error())
		return nil, errors.Wrap(err, "error in prepare query statement in SearchActorsAndMovies")
	}
	defer func() {
		if closeErr := execRowStaff.Close(); closeErr != nil {
			logger.Error().Err(closeErr).Msg("failed_to_close_statement")
			return
		}
	}()

	for execRowStaff.Next() {
		var personID sql.NullInt64
		var personFullName sql.NullString
		var personPhoto sql.NullString

		if err := execRowStaff.Scan(
			&personID,
			&personFullName,
			&personPhoto,
		); err != nil {
			errMsg := errors.Wrap(err, "error in query scan in SearchActorsAndMovies")
			logger.Error().Err(errMsg).Msg(errMsg.Error())
			return nil, errMsg
		}

		if !personID.Valid {
			continue
		}

		person := mocks.PersonJSON{
			ID:       int(personID.Int64),
			FullName: personFullName.String,
			Photo:    personPhoto.String,
		}

		resStaff = append(resStaff, person)
	}
	if execErr := execRowStaff.Err(); execErr != nil {
		errMsg := errors.Wrap(execErr, "error in query next in SearchActorsAndMovies")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		return nil, errMsg
	}

	resMovies := []mocks.Movie{}
	execRowMovie, err := r.pgdb.Query(searchMoviesBySearchStringQuery, searchStr)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "error in query statement in SearchActorsAndMovies").Error())
		return nil, errors.Wrap(err, "error in prepare query statement in SearchActorsAndMovies")
	}
	defer func() {
		if closeErr := execRowMovie.Close(); closeErr != nil {
			logger.Error().Err(closeErr).Msg("failed_to_close_statement")
			return
		}
	}()

	for execRowMovie.Next() {
		var movieID sql.NullInt64
		var movieTitle sql.NullString
		var moViePreviewURL sql.NullString

		if err := execRowMovie.Scan(
			&movieID,
			&movieTitle,
			&moViePreviewURL,
		); err != nil {
			errMsg := errors.Wrap(err, "error in query scan in SearchActorsAndMovies")
			logger.Error().Err(errMsg).Msg(errMsg.Error())
			return nil, errMsg
		}

		if !movieID.Valid {
			continue
		}

		movie := mocks.Movie{
			ID:         int(movieID.Int64),
			Title:      movieTitle.String,
			PreviewURL: moViePreviewURL.String,
		}

		resMovies = append(resMovies, movie)
	}
	if execErr := execRowMovie.Err(); execErr != nil {
		errMsg := errors.Wrap(execErr, "error in query next in SearchActorsAndMovies")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		return nil, errMsg
	}

	result := models.SearchResponseJSON{
		MovieCollection: resMovies,
		Actors:          resStaff,
	}
	return &result, nil

}
