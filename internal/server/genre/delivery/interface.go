package delivery

import "net/http"

type GenreHandlerInterface interface {
	GetGenreByID(w http.ResponseWriter, r *http.Request)
	GetGenres(w http.ResponseWriter, r *http.Request)
}
