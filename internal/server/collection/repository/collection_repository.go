package repository

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/rs/zerolog/log"
)

type CollectionRepository struct {
	repo mocks.Collections
}

func NewCollectionRepository(repo *mocks.Collections) *CollectionRepository {
	return &CollectionRepository{
		repo: *repo,
	}
}

func (r *CollectionRepository) GetMainPageCollectionsFromRepo(ctx context.Context) (mocks.Collections, error) {
	logger := log.Ctx(ctx)

	logger.Info().Msg("Get Collections from repo")

	return r.repo, nil
}
