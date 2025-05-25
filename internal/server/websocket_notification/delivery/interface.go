package delivery

import "net/http"

type NotificationHandlerInterface interface {
	WSHandler(w http.ResponseWriter, r *http.Request)
	Stop()
}
