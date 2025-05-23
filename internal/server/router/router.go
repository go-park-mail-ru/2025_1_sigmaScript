package router

import (
	"net/http"

	authDelivery "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/delivery"
	collectionDelivery "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/collection/delivery"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/middleware"
	movieDelivery "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/delivery"
	staffDelivery "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/staff_person/delivery"
	userDelivery "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/user/delivery/http"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	return router
}

func SetupAuth(router *mux.Router, authHandler authDelivery.AuthHandlerInterface) {
	authSubRouter := router.PathPrefix("/auth").Subrouter()

	authSubRouter.HandleFunc("/login", authHandler.Login).Methods(http.MethodPost, http.MethodOptions).Name("LoginRoute")
	authSubRouter.HandleFunc("/logout", authHandler.Logout).Methods(http.MethodPost, http.MethodOptions).Name("LogoutRoute")
	authSubRouter.HandleFunc("/register", authHandler.Register).Methods(http.MethodPost, http.MethodOptions).Name("RegisterRoute")
	authSubRouter.HandleFunc("/session", authHandler.Session).Methods(http.MethodGet, http.MethodOptions).Name("SessionRoute")
}

func SetupCollections(router *mux.Router, collectionHandler collectionDelivery.CollectionHandlerInterface) {
	router.HandleFunc("/collections/", collectionHandler.GetMainPageCollections).Methods(http.MethodGet, http.MethodOptions).Name("CollectionsRoute")
}

func SetupStaffPersonHandlers(router *mux.Router, staffPersonHandler staffDelivery.StaffPersonHandlerInterface) {
	router.HandleFunc("/name/{person_id}", staffPersonHandler.GetPerson).Methods(http.MethodGet, http.MethodOptions).Name("StaffPersonRoute")
}

func SetupMovieHandlers(router *mux.Router, movieHandler movieDelivery.MovieHandlerInterface) {
	router.HandleFunc("/movie/{movie_id}", movieHandler.GetMovie).Methods(http.MethodGet, http.MethodOptions).Name("MovieRoute")
}

func SetupUserHandlers(router *mux.Router, userHandler userDelivery.UserHandlerInterface) {
	router.HandleFunc("/users", userHandler.UpdateUser).Methods(http.MethodPost, http.MethodOptions).Name("UpdateUserRoute")
}

func ApplyMiddlewares(router *mux.Router) {
	router.Use(middleware.RequestWithLoggerMiddleware)
	router.Use(middleware.PreventPanicMiddleware)
	router.Use(middleware.MiddlewareCors)
}
