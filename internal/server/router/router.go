package router

import (
	"context"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func New(ctx context.Context) *mux.Router {
	log.Info().Msg("Configuring routes")

	authHandler := handlers.NewAuthHandler(ctx)

	router := mux.NewRouter()
	router.HandleFunc("/collections/", handlers.GetCollections).Methods(http.MethodGet)
	router.HandleFunc("/auth/login/", authHandler.LoginHandler).Methods(http.MethodPost)
	router.HandleFunc("/auth/logout/", authHandler.LogoutHandler).Methods(http.MethodPost)
	router.HandleFunc("/auth/register/", authHandler.RegisterHandler).Methods(http.MethodPost)

	log.Info().Msg("Routes configured successfully")
	return router
}
