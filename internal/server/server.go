package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	deliveryAuth "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/delivery"
	repoAuthSessions "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/repository"
	serviceAuth "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/service"
	repoUsers "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/user/repository"

	deliveryUsers "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/user/delivery/http"
	serviceUsers "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/user/service"

	deliveryCollection "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/collection/delivery"
	repoCollection "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/collection/repository"
	serviceCollection "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/collection/service"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/router"

	deliveryStaff "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/staff_person/delivery"
	repoStaff "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/staff_person/repository"
	serviceStaff "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/staff_person/service"

	"github.com/rs/zerolog/log"
)

type Server struct {
	Config     *config.Config
	httpServer *http.Server
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down server")
	return s.httpServer.Shutdown(ctx)
}

func New(cfg *config.Config) *Server {
	log.Info().Msg("Initializing server")

	s := &Server{
		Config: cfg,
	}

	log.Info().Msg("Server initialized successfully")
	return s
}

func (s *Server) Run() error {
	sessionRepo := repoAuthSessions.NewSessionRepository()
	sessionService := serviceAuth.NewSessionService(config.WrapCookieContext(context.Background(), &s.Config.Cookie), sessionRepo)

	userRepo := repoUsers.NewUserRepository()
	userService := serviceUsers.NewUserService(userRepo)
	userHandler := deliveryUsers.NewUserHandler(config.WrapCookieContext(context.Background(), &s.Config.Cookie), userService, sessionService)

	authHandler := deliveryAuth.NewAuthHandler(config.WrapCookieContext(context.Background(), &s.Config.Cookie), userService, sessionService)

	staffPersonRepo := repoStaff.NewStaffPersonRepository(&mocks.ExistingActors)
	staffPersonService := serviceStaff.NewStaffPersonService(staffPersonRepo)
	staffPersonHandler := deliveryStaff.NewStaffPersonHandler(staffPersonService)

	collectionRepo := repoCollection.NewCollectionRepository(&mocks.MainPageCollections)
	collectionService := serviceCollection.NewCollectionService(collectionRepo)
	collectionHandler := deliveryCollection.NewCollectionHandler(collectionService)

	mx := router.NewRouter()

	log.Info().Msg("Configuring routes")

	router.ApplyMiddlewares(mx)
	router.SetupAuth(mx, authHandler)
	router.SetupCollections(mx, collectionHandler)
	router.SetupStaffPersonHandlers(mx, staffPersonHandler)
	router.SetupUserHandlers(mx, userHandler)

	log.Info().Msg("Routes configured successfully")

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.Config.Server.Address, s.Config.Server.Port),
		ReadTimeout:  s.Config.Server.ReadTimeout,
		WriteTimeout: s.Config.Server.WriteTimeout,
		IdleTimeout:  s.Config.Server.IdleTimeout,
		Handler:      mx,
	}

	s.httpServer = srv

	log.Info().Msg("Running server")
	return s.httpServer.ListenAndServe()
}
