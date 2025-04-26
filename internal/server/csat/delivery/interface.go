package delivery

import "net/http"

type CSATHandlerInterface interface {
	GetAllCSATReviews(w http.ResponseWriter, r *http.Request)
	CreateCSATReview(w http.ResponseWriter, r *http.Request)
	GetCSATStatistic(w http.ResponseWriter, r *http.Request)
}
