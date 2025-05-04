package router

import (
	"net/http"

	authDelivery "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/delivery"
	collectionDelivery "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/collection/delivery"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/middleware"
	movieDelivery "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/delivery"
	reviewsDelivery "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/reviews/delivery"
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

func SetupReviewsHandlers(router *mux.Router, reviewsHandler reviewsDelivery.ReviewHandlerInterface) {
	router.HandleFunc("/movie/{movie_id}/reviews", reviewsHandler.GetAllReviewsOfMovie).Methods(http.MethodGet, http.MethodOptions).Name("GetReviewsOfMovieRoute")
	router.HandleFunc("/movie/{movie_id}/reviews", reviewsHandler.CreateReview).Methods(http.MethodPost, http.MethodOptions).Name("CreateReviewOfMovieRoute")
}

func SetupMovieHandlers(router *mux.Router, movieHandler movieDelivery.MovieHandlerInterface) {
	router.HandleFunc("/movie/{movie_id}", movieHandler.GetMovie).Methods(http.MethodGet, http.MethodOptions).Name("MovieRoute")
}

func SetupUserHandlers(router *mux.Router, userHandler userDelivery.UserHandlerInterface) {
	router.HandleFunc("/users", userHandler.UpdateUser).Methods(http.MethodPost, http.MethodOptions).Name("UpdateUserRoute")
	router.HandleFunc("/users/avatar", userHandler.UpdateUserAvatar).Methods(http.MethodPost, http.MethodOptions).Name("UpdateUserAvatarRoute")
	router.HandleFunc("/profile", userHandler.GetProfile).Methods(http.MethodGet, http.MethodOptions).Name("GetProfileRoute")

	router.HandleFunc("/movie/{movie_id}/favorite", userHandler.AddFavoriteMovie).Methods(http.MethodPost, http.MethodOptions).Name("AddFavoriteMovieRoute")
	router.HandleFunc("/name/{person_id}/favorite", userHandler.AddFavoriteActor).Methods(http.MethodPost, http.MethodOptions).Name("AddFavoriteActorRoute")
}

func ApplyMiddlewares(router *mux.Router) {
	router.Use(middleware.RequestWithLoggerMiddleware)
	router.Use(middleware.PreventPanicMiddleware)
	router.Use(middleware.MiddlewareCors)

	// router.Use(middleware.CsrfTokenMiddleware)
}
