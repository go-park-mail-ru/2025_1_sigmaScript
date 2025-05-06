package router

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	router := NewRouter()
	require.NotNil(t, router)
}

// func TestSetup(t *testing.T) {
// 	cfg, err := config.New()
// 	require.NoError(t, err)
// 	require.NotNil(t, cfg)

// 	sessionRepo := repoAuthSessions.NewSessionRepository()
// 	sessionService := serviceAuth.NewSessionService(config.WrapCookieContext(context.Background(), &cfg.Cookie), sessionRepo)

// 	// userRepo := repoUsers.NewUserRepository(nil)
// 	userService := serviceUsers.NewUserService(nil)
// 	userHandler := deliveryUsers.NewUserHandler(config.WrapCookieContext(context.Background(), &cfg.Cookie), userService, sessionService)

// 	authHandler := deliveryAuth.NewAuthHandler(config.WrapCookieContext(context.Background(), &cfg.Cookie), userService, sessionService)

// 	staffPersonRepo := repoStaff.NewStaffPersonRepository(&mocks.ExistingActors)
// 	staffPersonService := serviceStaff.NewStaffPersonService(staffPersonRepo)
// 	staffPersonHandler := deliveryStaff.NewStaffPersonHandler(staffPersonService)

// 	collectionRepo := repoCollection.NewCollectionRepository(&mocks.MainPageCollections)
// 	collectionService := serviceCollection.NewCollectionService(collectionRepo)
// 	collectionHandler := deliveryCollection.NewCollectionHandler(collectionService)

// 	movieRepo := repoMovie.NewMovieRepository(&mocks.ExistingMovies)
// 	movieService := serviceMovie.NewMovieService(movieRepo)
// 	movieHandler := deliveryMovie.NewMovieHandler(movieService)

// 	mx := NewRouter()

// 	log.Info().Msg("Configuring routes")

// 	ApplyMiddlewares(mx)
// 	SetupAuth(mx, authHandler)
// 	SetupCollections(mx, collectionHandler)
// 	SetupStaffPersonHandlers(mx, staffPersonHandler)
// 	SetupUserHandlers(mx, userHandler)
// 	SetupMovieHandlers(mx, movieHandler)
// }
