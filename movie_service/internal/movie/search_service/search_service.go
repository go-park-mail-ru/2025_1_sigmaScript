package search_service

import (
	"context"
	"errors"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/models"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mocks.go -package=service_mocks SearchRepositoryInterface
type SearchRepositoryInterface interface {
	SearchActorsAndMovies(ctx context.Context, searchStr string) (*models.SearchResponseJSON, error)
}

type SearchService struct {
	searchRepo SearchRepositoryInterface
}

func NewSearchService(searchRepo SearchRepositoryInterface) *SearchService {
	return &SearchService{
		searchRepo: searchRepo,
	}
}

func (s *SearchService) SearchActorsAndMovies(ctx context.Context, searchStr string) (*models.SearchResponseJSON, error) {
	logger := log.Ctx(ctx)

	if len([]rune(searchStr)) < 3 {
		errMsg := errors.New(errs.ErrMsgLengthTooShort)
		logger.Error().Err(errMsg).Msgf("error happened while getting search results by search request '%s': %v", searchStr, errMsg.Error())
		return nil, errMsg
	}

	genre, err := s.searchRepo.SearchActorsAndMovies(ctx, searchStr)
	if err != nil {
		logger.Error().Err(err).Msgf("error happened while getting search results by search request '%s': %v", searchStr, err.Error())
		return nil, err
	}

	return genre, nil
}
