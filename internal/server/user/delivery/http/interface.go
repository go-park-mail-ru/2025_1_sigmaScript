package http

import "net/http"

//go:generate mockgen -source=interface.go -destination=mocks/mock.go
type UserHandlerInterface interface {
	UpdateUser(w http.ResponseWriter, r *http.Request)
	UpdateUserAvatar(w http.ResponseWriter, r *http.Request)
	GetProfile(w http.ResponseWriter, r *http.Request)
	AddFavoriteMovie(w http.ResponseWriter, r *http.Request)
	AddFavoriteActor(w http.ResponseWriter, r *http.Request)
	RemoveFavoriteMovie(w http.ResponseWriter, r *http.Request)
	RemoveFavoriteActor(w http.ResponseWriter, r *http.Request)
}
