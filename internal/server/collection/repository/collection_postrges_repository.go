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
	cm.collection_id as collection_id,
	c.name AS collection_name,
	m.id as id,
	m.name as title,
	m.poster as preview_url,
	m.duration as duration
FROM collection_movie cm
JOIN movie m ON cm.movie_id = m.id
JOIN collection c ON cm.collection_id = c.id
where c.is_main_collection = TRUE;
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

		var movieID int
		var movieTitle sql.NullString
		var moviePreviewURL sql.NullString
		var movieDuration sql.NullString

		if err := execRow.Scan(
			&collectionId,
			&collectionName,
			&movieID,
			&movieTitle,
			&moviePreviewURL,
			&movieDuration,
		); err != nil {
			errMsg := errors.Wrap(err, "error in query scan in GetMainPageCollectionsFromRepo")
			logger.Error().Err(errMsg).Msg(errMsg.Error())
			return nil, errMsg
		}

		if _, ok := resCollections[collectionName]; !ok {
			resCollections[collectionName] = mocks.Collection{}
		}

		movie := mocks.Movie{
			ID:         movieID,
			Title:      movieTitle.String,
			PreviewURL: moviePreviewURL.String,
			Duration:   movieDuration.String,
		}
		resCollections[collectionName][movie.ID] = movie

	}

	if execErr := execRow.Err(); execErr != nil {
		errMsg := errors.Wrap(execErr, "error in query next in GetMainPageCollectionsFromRepo")
		logger.Error().Err(errMsg).Msg(errMsg.Error())
		return nil, errMsg
	}

	return resCollections, nil
}
