package dto

// NewReviewDataJSON delivery layer review info
type NewReviewDataJSON struct {
	ReviewText string `json:"review_text,omitempty"`
	Score      int    `json:"score"`
}
