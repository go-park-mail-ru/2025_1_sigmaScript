package mocks

import "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/csat/delivery/dto"

type CSATRepo struct {
	Rating  float64
	Reviews map[int]*dto.CSATReviewDataJSON
}
