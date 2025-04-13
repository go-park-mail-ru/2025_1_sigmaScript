package delivery

import "net/http"

type StaffPersonHandlerInterface interface {
	GetPerson(w http.ResponseWriter, r *http.Request)
}
