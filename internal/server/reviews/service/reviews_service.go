package service

import (
	"context"
	"fmt"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/common"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
)

const (
	noData = ""
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
	pageNumber := 1
	if len(paginatorPageNumber) > 0 && paginatorPageNumber[0] > 0 {
		pageNumber = paginatorPageNumber[0]
	}

	fmt.Println(pageNumber, common.REVIEWS_PER_PAGE)
	return []mocks.ReviewJSON{}
}

func (s *ReviewService) GetReview(ctx context.Context, movieID, userID int) (mocks.ReviewJSON, error)

func (s *ReviewService) CreateReview(ctx context.Context, newReview mocks.ReviewJSON) error
func (s *ReviewService) UpdateReview(ctx context.Context, updatedReview mocks.ReviewJSON) error
func (s *ReviewService) DeleteReview(ctx context.Context, reviewID int) error
