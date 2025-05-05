package delivery

import "net/http"

type SearchHandlerInterface interface {
	SearchActorsAndMovies(w http.ResponseWriter, r *http.Request)
}
