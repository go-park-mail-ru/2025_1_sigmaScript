package router

import (
	"net/http"

	"github.com/gorilla/mux"

	csrfDelivery "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/csrf/delivery"
)

func SetupCsrf(router *mux.Router, csrfTokenHandler csrfDelivery.CSRFHandlerInterface) {
	router.HandleFunc("/auth/csrf-token", csrfTokenHandler.CreateCSRFTokenHandler).Methods(http.MethodGet, http.MethodOptions).Name("CsrfTokenRoute")
}
