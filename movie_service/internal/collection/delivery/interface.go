package delivery

import "net/http"

type CollectionHandlerInterface interface {
	GetMainPageCollections(w http.ResponseWriter, r *http.Request)
}
