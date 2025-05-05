package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/db"
	deliveryAuth "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/delivery"
	repoAuthSessions "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/repository"
	serviceAuth "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/service"
	repoUsers "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/user/repository"

	deliveryUsers "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/user/delivery/http"
	serviceUsers "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/user/service"

	deliveryCollection "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/collection/delivery"
	repoCollection "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/collection/repository"
	serviceCollection "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/collection/service"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/router"

	deliveryStaff "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/staff_person/delivery"
	repoStaff "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/staff_person/repository"
	serviceStaff "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/staff_person/service"

	deliveryMovie "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/delivery"
	repoMovie "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/repository"
	serviceMovie "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/service"

	csrfDelivery "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/csrf/delivery"
	deliveryReviews "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/reviews/delivery"

	deliveryGenre "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/genre/delivery"
	repoGenre "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/genre/repository"
	serviceGenre "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/genre/service"

	deliverySearch "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/search/delivery"
	repoSearch "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/search/repository"
	serviceSearch "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/search/service"

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
	ctxPgDb := config.WrapPgDatabaseContext(context.Background(), s.Config.PostgresConfig)
	ctxPgDb, cancelPgDb := context.WithTimeout(ctxPgDb, time.Second*30)
	defer cancelPgDb()

	pgdb, err := db.SetupDatabase(ctxPgDb, cancelPgDb)
	if err != nil {
		return fmt.Errorf("error couldnt connect to postgres database: %w", err)
	}

	sessionRepo := repoAuthSessions.NewSessionRepository()
	sessionService := serviceAuth.NewSessionService(config.WrapCookieContext(context.Background(), &s.Config.Cookie), sessionRepo)

	userRepo := repoUsers.NewUserRepository(pgdb)
	userService := serviceUsers.NewUserService(userRepo)
	userHandler := deliveryUsers.NewUserHandler(config.WrapCookieContext(context.Background(), &s.Config.Cookie), userService, sessionService)

	authHandler := deliveryAuth.NewAuthHandler(config.WrapCookieContext(context.Background(), &s.Config.Cookie), userService, sessionService)

	csrfHandler := csrfDelivery.NewCSRFHandler(config.WrapCookieContext(context.Background(), &s.Config.Cookie), sessionService)

	staffPersonRepo := repoStaff.NewStaffPersonPostgresRepository(pgdb)
	staffPersonService := serviceStaff.NewStaffPersonService(staffPersonRepo)
	staffPersonHandler := deliveryStaff.NewStaffPersonHandler(staffPersonService)

	collectionRepo := repoCollection.NewCollectionPostgresRepository(pgdb)
	collectionService := serviceCollection.NewCollectionService(collectionRepo)
	collectionHandler := deliveryCollection.NewCollectionHandler(collectionService)

	movieRepo := repoMovie.NewMoviePostgresRepository(pgdb)
	movieService := serviceMovie.NewMovieService(movieRepo)
	movieHandler := deliveryMovie.NewMovieHandler(movieService)

	movieReviewHandler := deliveryReviews.NewReviewHandler(userService, sessionService, movieService)

	genreRepo := repoGenre.NewGenreRepository(pgdb)
	genreService := serviceGenre.NewGenreService(genreRepo)
	genreHandler := deliveryGenre.NewGenreHandler(genreService)

	searchRepo := repoSearch.NewSearchRepository(pgdb)
	searchService := serviceSearch.NewSearchService(searchRepo)
	searchHandler := deliverySearch.NewSearchHandler(searchService)

	mx := router.NewRouter()

	log.Info().Msg("Configuring routes")

	router.ApplyMiddlewares(mx)
	router.SetupAuth(mx, authHandler)

	router.SetupCsrf(mx, csrfHandler)

	router.SetupCollections(mx, collectionHandler)
	router.SetupStaffPersonHandlers(mx, staffPersonHandler)
	router.SetupUserHandlers(mx, userHandler)
	router.SetupMovieHandlers(mx, movieHandler)
	router.SetupReviewsHandlers(mx, movieReviewHandler)
	router.SetupGenresHandlers(mx, genreHandler)
	router.SetupSearchHandlers(mx, searchHandler)

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
