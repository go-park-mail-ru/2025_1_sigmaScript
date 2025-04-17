package router

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	deliveryAuth "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/delivery"
	repoAuthSessions "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/repository"
	serviceAuth "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/service"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
	deliveryMovie "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/delivery"
	repoMovie "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/repository"
	serviceMovie "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/service"
	deliveryStaff "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/staff_person/delivery"
	repoStaff "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/staff_person/repository"
	serviceStaff "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/staff_person/service"
	deliveryUsers "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/user/delivery/http"
	repoUsers "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/user/repository"
	serviceUsers "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/user/service"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"

	deliveryCollection "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/collection/delivery"
	repoCollection "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/collection/repository"
	serviceCollection "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/collection/service"
)

func TestNew(t *testing.T) {
	router := NewRouter()
	require.NotNil(t, router)
}

func TestSetup(t *testing.T) {
	cfg, err := config.New()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	sessionRepo := repoAuthSessions.NewSessionRepository()
	sessionService := serviceAuth.NewSessionService(config.WrapCookieContext(context.Background(), &cfg.Cookie), sessionRepo)

	userRepo := repoUsers.NewUserRepository(nil)
	userService := serviceUsers.NewUserService(userRepo)
	userHandler := deliveryUsers.NewUserHandler(config.WrapCookieContext(context.Background(), &cfg.Cookie), userService, sessionService)

	authHandler := deliveryAuth.NewAuthHandler(config.WrapCookieContext(context.Background(), &cfg.Cookie), userService, sessionService)

	staffPersonRepo := repoStaff.NewStaffPersonRepository(&mocks.ExistingActors)
	staffPersonService := serviceStaff.NewStaffPersonService(staffPersonRepo)
	staffPersonHandler := deliveryStaff.NewStaffPersonHandler(staffPersonService)

	collectionRepo := repoCollection.NewCollectionRepository(&mocks.MainPageCollections)
	collectionService := serviceCollection.NewCollectionService(collectionRepo)
	collectionHandler := deliveryCollection.NewCollectionHandler(collectionService)

	movieRepo := repoMovie.NewMovieRepository(&mocks.ExistingMovies)
	movieService := serviceMovie.NewMovieService(movieRepo)
	movieHandler := deliveryMovie.NewMovieHandler(movieService)

	mx := NewRouter()

	log.Info().Msg("Configuring routes")

	ApplyMiddlewares(mx)
	SetupAuth(mx, authHandler)
	SetupCollections(mx, collectionHandler)
	SetupStaffPersonHandlers(mx, staffPersonHandler)
	SetupUserHandlers(mx, userHandler)
	SetupMovieHandlers(mx, movieHandler)
}
