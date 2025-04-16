package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mocks.go -package=service_mocks CollectionRepositoryInterface
type CollectionRepositoryInterface interface {
	GetMainPageCollectionsFromRepo(ctx context.Context) (mocks.Collections, error)
}

type CollectionService struct {
	collectionRepo CollectionRepositoryInterface
}

func NewCollectionService(collectionRepo CollectionRepositoryInterface) *CollectionService {
	return &CollectionService{
		collectionRepo: collectionRepo,
	}
}

func (s *CollectionService) GetMainPageCollections(ctx context.Context) (mocks.Collections, error) {
	logger := log.Ctx(ctx)

	logger.Info().Msg("Get Collections from service")

	collections, err := s.collectionRepo.GetMainPageCollectionsFromRepo(ctx)
	if err != nil {
		logger.Err(err).Msg(err.Error())
		return nil, err
	}

	return collections, nil
}
