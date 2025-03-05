package router

import (
	"net/http"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/handlers"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/middleware"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	return router
}

func SetupAuth(router *mux.Router, authHandler handlers.AuthHandlerInterdace) {
	authSubRouter := router.PathPrefix("/auth").Subrouter()

	authSubRouter.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost, http.MethodOptions).Name("LoginRoute")
	authSubRouter.HandleFunc("/logout", authHandler.Logout).Methods(http.MethodPost, http.MethodOptions).Name("LogoutRoute")
	authSubRouter.HandleFunc("/register", authHandler.Register).Methods(http.MethodPost, http.MethodOptions).Name("RegisterRoute")
	authSubRouter.HandleFunc("/session", authHandler.Session).Methods(http.MethodGet, http.MethodOptions).Name("SessionRoute")
}

func SetupCollections(router *mux.Router) {
	router.HandleFunc("/collections/", handlers.GetCollections).Methods(http.MethodGet, http.MethodOptions).Name("CollectionsRoute")
}

func ApplyMiddlewares(router *mux.Router) {
	// router.Use(middleware.PreventPanicMiddleware)
	router.Use(middleware.MiddlewareCors)
}
