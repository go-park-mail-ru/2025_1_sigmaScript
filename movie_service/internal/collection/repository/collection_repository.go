package repository

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/rs/zerolog/log"
)

type CollectionRepository struct {
	db *mocks.Collections
}

func NewCollectionRepository(repo *mocks.Collections) *CollectionRepository {
	return &CollectionRepository{
		db: repo,
	}
}

func (r *CollectionRepository) GetMainPageCollectionsFromRepo(ctx context.Context) (mocks.Collections, error) {
	logger := log.Ctx(ctx)

	logger.Info().Msg("Get Collections from repo")

	return *r.db, nil
}
