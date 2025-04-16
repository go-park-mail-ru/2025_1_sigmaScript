package repository

import (
	synccredmap "github.com/go-park-mail-ru/2025_1_sigmaScript/pkg/sync_cred_map"
)

type ReviewsRepository struct {
	// sessionID --> username
	rdb synccredmap.SyncCredentialsMap
}

func NewReviewsRepository() *ReviewsRepository {
	res := &ReviewsRepository{
		rdb: *synccredmap.NewSyncCredentialsMap(),
	}

	return res
}
