package repository

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	getGenreMoviesQuery = `
SELECT
    g.name as genre_name,
    m.id as id,
    m.name as title,
    m.poster as preview_url,
    m.rating as rating
FROM genre g
LEFT JOIN movie_genre mg ON g.id = mg.genre_id
LEFT JOIN movie m ON m.id = mg.movie_id
WHERE g.id = $1;
`

	getAllGenresQuery = `
SELECT
	g.id as genre_id,
	g.name AS genre_name,
	m.id as id,
	m.name as title,
	m.poster as preview_url,
	m.rating as rating
FROM genre g
LEFT JOIN movie_genre mg ON mg.genre_id = g.id
LEFT JOIN movie m ON mg.movie_id = m.id
ORDER BY g.id, m.id;
`
)

type GenreRepository struct {
	pgdb *sql.DB
}

func NewGenreRepository(db *sql.DB) *GenreRepository {
	return &GenreRepository{
		pgdb: db,
	}
}

func (r *GenreRepository) GetGenreFromRepoByID(ctx context.Context, genreID string) (*mocks.Genre, error) {
	logger := log.Ctx(ctx)

	// get movies
	var genreName string
	resMovies := []mocks.Movie{}
	execRowMovie, err := r.pgdb.Query(getGenreMoviesQuery, genreID)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "error in query statement in GetGenreFromRepoByID").Error())
		return nil, errors.Wrap(err, "error in prepare query statement in GetGenreFromRepoByID")
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
		var moviePreviewURL sql.NullString
		var movieRating sql.NullFloat64

		if err := execRowMovie.Scan(
			&genreName,
			&movieID,
			&movieTitle,
			&moviePreviewURL,
			&movieRating,
		); err != nil {
			errMsg := errors.Wrap(err, "error in query scan in GetGenreFromRepoByID")
			logger.Error().Err(errMsg).Msg(errMsg.Error())
			return nil, errMsg
		}

		// skip if not found
		if !movieID.Valid {
			continue
		}

		movie := mocks.Movie{
			ID:         int(movieID.Int64),
			Title:      movieTitle.String,
			PreviewURL: moviePreviewURL.String,
			Rating:     movieRating.Float64,
		}

		resMovies = append(resMovies, movie)
	}
	if execErr := execRowMovie.Err(); execErr != nil {
		errMsg := errors.Wrap(execErr, "error in query next in GetGenreFromRepoByID")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		return nil, errMsg
	}

	result := mocks.Genre{
		ID:     genreID,
		Name:   genreName,
		Movies: resMovies,
	}
	return &result, nil
}

func (r *GenreRepository) GetAllGenresFromRepo(ctx context.Context) (*[]mocks.Genre, error) {
	logger := log.Ctx(ctx)
	logger.Info().Msg("Get Collections from postgres repo")

	var resultGenres []mocks.Genre

	execRow, err := r.pgdb.Query(getAllGenresQuery)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "error in query statement in GetAllGenresFromRepo").Error())
		return nil, errors.Wrap(err, "error in prepare query statement in GetAllGenresFromRepo")
	}
	defer func() {
		if closeErr := execRow.Close(); closeErr != nil {
			logger.Error().Err(closeErr).Msg("failed_to_close_statement")
			return
		}
	}()

	for execRow.Next() {
		var genreID int
		var genreName string

		var movieID sql.NullInt64
		var movieTitle sql.NullString
		var moviePreviewURL sql.NullString
		var movieRating sql.NullFloat64

		if err := execRow.Scan(
			&genreID,
			&genreName,
			&movieID,
			&movieTitle,
			&moviePreviewURL,
			&movieRating,
		); err != nil {
			errMsg := errors.Wrap(err, "error in query scan in GetAllGenresFromRepo")
			logger.Error().Err(errMsg).Msg(errMsg.Error())
			return nil, errMsg
		}

		// check if we have enough genres in resultGenres
		if len(resultGenres) < genreID {
			resultGenres = append(resultGenres, mocks.Genre{})
			resultGenres[genreID-1].ID = strconv.Itoa(genreID)
			resultGenres[genreID-1].Name = genreName
		}

		// skip if not found movie
		if !movieID.Valid {
			continue
		}

		movie := mocks.Movie{
			ID:         int(movieID.Int64),
			Title:      movieTitle.String,
			PreviewURL: moviePreviewURL.String,
			Rating:     movieRating.Float64,
		}

		resultGenres[genreID-1].Movies = append(resultGenres[genreID-1].Movies, movie)
	}

	if execErr := execRow.Err(); execErr != nil {
		errMsg := errors.Wrap(execErr, "error in query next in GetAllGenresFromRepo")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		return nil, errMsg
	}

	return &resultGenres, nil
}
