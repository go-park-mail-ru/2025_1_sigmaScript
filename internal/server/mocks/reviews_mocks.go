package mocks

// ReviewJSON delivery layer review info
type ReviewJSON struct {
	ID         int    `json:"id"`
	UserID     int    `json:"user_id"`
	MovieID    int    `json:"movie_id"`
	ReviewText string `json:"review_text"`
	Score      string `json:"score"`
}

type Reviews map[int]ReviewJSON

var ExistingReviews = Reviews{
	1: {ID: 1, UserID: 1, MovieID: 1},
	2: {ID: 2, UserID: 2, MovieID: 2},
	3: {ID: 3, UserID: 3, MovieID: 1},
}
