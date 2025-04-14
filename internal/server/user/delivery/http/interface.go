package http

import "net/http"

type UserHandlerInterface interface {
	UpdateUser(w http.ResponseWriter, r *http.Request)
	UpdateUserAvatar(w http.ResponseWriter, r *http.Request)
}
