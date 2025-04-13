package http

import "net/http"

type UserHandlerInterface interface {
	UpdateUser(w http.ResponseWriter, r *http.Request)
}
