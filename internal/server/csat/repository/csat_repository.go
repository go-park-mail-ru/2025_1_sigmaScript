package repository

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/csat/delivery/dto"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/rs/zerolog/log"

	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
)

const (
	DEFAULT_CSAT_SCORE = 5
)

type CSATRepository struct {
	mu sync.RWMutex
	db *mocks.CSATRepo
	// pgdb *sql.DB
}

func NewCSATRepository(CSATDB *mocks.CSATRepo) *CSATRepository {
	return &CSATRepository{db: CSATDB}
}

func (r *CSATRepository) GetAllCSATReviews(ctx context.Context) (*[]dto.CSATReviewDataJSON, error) {
	logger := log.Ctx(ctx)

	r.mu.RLock()
	defer r.mu.Unlock()

	reviews := []dto.CSATReviewDataJSON{}
	for _, val := range r.db.Reviews {
		reviews = append(reviews, *val)
	}

	if len(reviews) == 0 {
		logger.Err(errs.ErrCSATReviewsNotFound).Msg(errs.ErrCSATReviewsNotFound.Error())
		return nil, errs.ErrCSATReviewsNotFound
	}

	return &reviews, nil
}

func (r *CSATRepository) CreateNewCSAT(ctx context.Context, newReview dto.CSATReviewDataJSON) error {
	// logger := log.Ctx(ctx)
	r.mu.Lock()
	defer r.mu.Unlock()

	reviewsCount := len(r.db.Reviews)

	// update average rating
	r.db.Rating = float64(math.Trunc((r.db.Rating+(float64(newReview.Score)-r.db.Rating)/float64(reviewsCount+1))*100)) / 100

	newReviewID := reviewsCount + 1
	csatReview := &dto.CSATReviewDataJSON{
		ID:        newReviewID,
		Score:     newReview.Score,
		CSATText:  newReview.CSATText,
		CreatedAt: time.Now().String(),
		User:      newReview.User,
	}

	(r.db.Reviews)[newReviewID] = csatReview
	return nil
}

func (r *CSATRepository) GetCSATStatistic(ctx context.Context) (*dto.CSATStatisticDataJSON, error) {
	logger := log.Ctx(ctx)

	r.mu.Lock()
	defer r.mu.Unlock()

	reviews := []dto.CSATReviewDataJSON{}
	for _, val := range r.db.Reviews {
		reviews = append(reviews, *val)
	}

	if len(reviews) == 0 {
		logger.Err(errs.ErrCSATReviewsNotFound).Msg(errs.ErrCSATReviewsNotFound.Error())
	}

	statistic := dto.CSATStatisticDataJSON{
		Statistic: dto.AverageCSATStatisticData{
			AverageRating: r.db.Rating,
			ReviewsCount:  len(reviews),
		},
		Reviews: reviews,
	}

	return &statistic, nil
}
