package http

import "net/http"

//go:generate mockgen -source=interface.go -destination=mocks/mock.go
type UserHandlerInterface interface {
	UpdateUser(w http.ResponseWriter, r *http.Request)
}
