package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/delivery"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/repository"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/service"
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
	// authHandler := handlers.NewAuthHandler(config.WrapCookieContext(context.Background(), &s.Config.Cookie))

	authRepository := repository.NewAuthRepository()
	authService := service.NewAuthService(config.WrapCookieContext(context.Background(), &s.Config.Cookie), authRepository)
	authHandler := delivery.NewAuthHandler(config.WrapCookieContext(context.Background(), &s.Config.Cookie), authService)

	staffPersonRepo := repoStaff.NewStaffPersonRepository(&mocks.ExistingActors)
	staffPersonService := serviceStaff.NewStaffPersonService(staffPersonRepo)
	staffPersonHandler := deliveryStaff.NewStaffPersonHandler(staffPersonService)

	mx := router.NewRouter()

	log.Info().Msg("Configuring routes")

	router.ApplyMiddlewares(mx)
	router.SetupAuth(mx, authHandler)
	router.SetupCollections(mx)
	router.SetupStaffPersonHandlers(mx, staffPersonHandler)

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
