package repository

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/rs/zerolog/log"
)

const (
	getMainPageCollectionsQuery = `
SELECT
	c.id as collection_id,
	c.name AS collection_name,
	m.id as id,
	m.name as title,
	CASE
        WHEN c.name = 'promo' THEN m.promo_poster
        ELSE m.poster
    END AS preview_url,
	m.duration as duration,
	m.rating as rating
FROM collection c
LEFT JOIN collection_movie cm ON cm.collection_id = c.id
LEFT JOIN movie m ON cm.movie_id = m.id
where c.is_main_collection = TRUE;
`

	getMainPageReleasesCalendarQuery = `
 SELECT
    m.id AS id,
    m.name AS title,
    m.poster AS preview_url,
    m.duration AS duration,
    m.release_year AS release_date,
		m.rating as rating
FROM movie m
WHERE
    m.release_year >= CURRENT_TIMESTAMP AND
    m.release_year <= CURRENT_TIMESTAMP + INTERVAL '4 weeks';
 `
)

type CollectionPostgresRepository struct {
	pgdb *sql.DB
}

func NewCollectionPostgresRepository(db *sql.DB) *CollectionPostgresRepository {
	return &CollectionPostgresRepository{
		pgdb: db,
	}
}

func (r *CollectionPostgresRepository) GetMainPageCollectionsFromRepo(ctx context.Context) (mocks.Collections, error) {
	logger := log.Ctx(ctx)
	logger.Info().Msg("Get Collections from postgres repo")

	resCollections := mocks.Collections{}
	// get main collections
	execRow, err := r.pgdb.Query(getMainPageCollectionsQuery)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "error in query statement in GetMainPageCollectionsFromRepo").Error())
		return nil, errors.Wrap(err, "error in prepare query statement in GetMainPageCollectionsFromRepo")
	}
	defer func() {
		if closeErr := execRow.Close(); closeErr != nil {
			logger.Error().Err(closeErr).Msg("failed_to_close_statement")
			return
		}
	}()

	for execRow.Next() {
		var collectionId int
		var collectionName string

		var movieID sql.NullInt64
		var movieTitle sql.NullString
		var moviePreviewURL sql.NullString
		var movieDuration sql.NullString
		var movieRating sql.NullFloat64

		if err := execRow.Scan(
			&collectionId,
			&collectionName,
			&movieID,
			&movieTitle,
			&moviePreviewURL,
			&movieDuration,
			&movieRating,
		); err != nil {
			errMsg := errors.Wrap(err, "error in query scan in GetMainPageCollectionsFromRepo")
			logger.Error().Err(errMsg).Msg(errMsg.Error())
			return nil, errMsg
		}

		// skip if not found
		if !movieID.Valid {
			continue
		}

		if _, ok := resCollections[collectionName]; !ok {
			resCollections[collectionName] = mocks.Collection{}
		}

		movie := mocks.Movie{
			ID:         int(movieID.Int64),
			Title:      movieTitle.String,
			PreviewURL: moviePreviewURL.String,
			Duration:   movieDuration.String,
			Rating:     movieRating.Float64,
		}
		resCollections[collectionName][movie.ID] = movie

	}

	if execErr := execRow.Err(); execErr != nil {
		errMsg := errors.Wrap(execErr, "error in query next in GetMainPageCollectionsFromRepo")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		return nil, errMsg
	}

	// get new releases collections
	execReleasesRow, err := r.pgdb.Query(getMainPageReleasesCalendarQuery)
	if err != nil {
		logger.Error().Err(err).Msg(errors.Wrapf(err, "error in query statement in GetMainPageCollectionsFromRepo").Error())
		return nil, errors.Wrap(err, "error in prepare query statement in GetMainPageCollectionsFromRepo")
	}
	defer func() {
		if closeErr := execReleasesRow.Close(); closeErr != nil {
			logger.Error().Err(closeErr).Msg("failed_to_close_statement")
			return
		}
	}()

	for execReleasesRow.Next() {
		var movieID sql.NullInt64
		var movieTitle sql.NullString
		var moviePreviewURL sql.NullString
		var movieDuration sql.NullString
		var movieReleaseDate sql.NullString
		var movieRating sql.NullFloat64

		if err := execReleasesRow.Scan(
			&movieID,
			&movieTitle,
			&moviePreviewURL,
			&movieDuration,
			&movieReleaseDate,
			&movieRating,
		); err != nil {
			errMsg := errors.Wrap(err, "error in query scan in GetMainPageReleaseCalendarCollectionsFromRepo")
			logger.Error().Err(errMsg).Msg(errMsg.Error())
			return nil, errMsg
		}

		// skip if not found
		if !movieID.Valid {
			continue
		}

		collectionName := "calendar"

		if _, ok := resCollections[collectionName]; !ok {
			resCollections[collectionName] = mocks.Collection{}
		}

		movie := mocks.Movie{
			ID:         int(movieID.Int64),
			Title:      movieTitle.String,
			PreviewURL: moviePreviewURL.String,
			Duration:   movieDuration.String,
			Rating:     movieRating.Float64,
		}
		resCollections[collectionName][movie.ID] = movie

	}

	if execErr := execReleasesRow.Err(); execErr != nil {
		errMsg := errors.Wrap(execErr, "error in query next in GetMainPageCollectionsFromRepo")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		return nil, errMsg
	}

	return resCollections, nil
}
