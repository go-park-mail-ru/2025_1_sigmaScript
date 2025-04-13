package delivery

import "net/http"

type MovieHandlerInterface interface {
	GetMovie(w http.ResponseWriter, r *http.Request)
}
