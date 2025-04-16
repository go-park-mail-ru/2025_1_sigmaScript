package delivery

import "net/http"

type ReviewHandlerInterface interface {
	GetAllReviewsOfMovie(w http.ResponseWriter, r *http.Request)
	CreateReview(w http.ResponseWriter, r *http.Request)
}
