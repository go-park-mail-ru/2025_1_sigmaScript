package mocks

// NewRevieDataJSON delivery layer review info
type NewReviewDataJSON struct {
	ReviewText string `json:"review_text"`
	Score      int    `json:"score"`
}
