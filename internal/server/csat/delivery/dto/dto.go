package dto

// NewCSATReviewDataJSON delivery layer review info
type NewCSATReviewDataJSON struct {
	ReviewText string `json:"csat_text,omitempty"`
	Score      int    `json:"score"`
}

type ReviewUserDataJSON struct {
	Login  string `json:"login"`
	Avatar string `json:"avatar,omitempty"`
}

type CSATReviewDataJSON struct {
	ID        int                `json:"id"`
	Score     int                `json:"score"`
	CSATText  string             `json:"csat_text"`
	CreatedAt string             `json:"created_at"`
	User      ReviewUserDataJSON `json:"user"`
}

type AverageCSATStatisticData struct {
	AverageRating float64 `json:"average_rating"`
	ReviewsCount  int     `json:"reviews_count"`
}

type CSATStatisticDataJSON struct {
	Statistic AverageCSATStatisticData `json:"statistic"`
	Reviews   []CSATReviewDataJSON     `json:"reviews,omitempty"`
}
