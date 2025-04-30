package service

import (
	"context"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
)

type SessionRepositoryInterface interface {
	GetSession(ctx context.Context, sessionID string) (string, error)
}

type ReviewService struct {
	sessionRepo SessionRepositoryInterface
}

func NewReviewService(ctx context.Context, sessionRepo SessionRepositoryInterface) *ReviewService {
	return &ReviewService{
		sessionRepo: sessionRepo,
	}
}

func (s *ReviewService) GetReviewsOfMovie(ctx context.Context, movieID int, paginatorPageNumber ...int) []mocks.ReviewJSON {
	return []mocks.ReviewJSON{}
}
