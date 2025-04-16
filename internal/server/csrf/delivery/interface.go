package delivery

import "net/http"

//go:generate mockgen -source=interface.go -destination=mocks/mock.go
type CSRFHandlerInterface interface {
	CreateCSRFTokenHandler(w http.ResponseWriter, r *http.Request)
}
