package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/csat/delivery/dto"
	"github.com/rs/zerolog/log"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/mocks.go -package=service_mocks MovieRepositoryInterface
type CSATRepositoryInterface interface {
	GetAllCSATReviews(ctx context.Context) (*[]dto.CSATReviewDataJSON, error)
	CreateNewCSAT(ctx context.Context, newReview dto.CSATReviewDataJSON) error
	GetCSATStatistic(ctx context.Context) (*dto.CSATStatisticDataJSON, error)
}

type CSATService struct {
	CSATRepository CSATRepositoryInterface
}

func NewCSATService(CSATRepo CSATRepositoryInterface) *CSATService {
	return &CSATService{
		CSATRepository: CSATRepo,
	}
}

func (s *CSATService) GetAllCSATReviews(ctx context.Context) (*[]dto.CSATReviewDataJSON, error) {
	logger := log.Ctx(ctx)

	reviews, err := s.CSATRepository.GetAllCSATReviews(ctx)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		return nil, err
	}

	return reviews, nil
}

func (s *CSATService) GetCSATStatistic(ctx context.Context) (*dto.CSATStatisticDataJSON, error) {
	logger := log.Ctx(ctx)

	statistic, err := s.CSATRepository.GetCSATStatistic(ctx)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		return nil, err
	}

	return statistic, nil
}

func (s *CSATService) CreateNewCSATReview(ctx context.Context, newReview dto.CSATReviewDataJSON) error {
	logger := log.Ctx(ctx)

	err := s.CSATRepository.CreateNewCSAT(ctx, newReview)
	if err != nil {
		logger.Error().Err(err).Msg(err.Error())
		return err
	}

	return nil
}
