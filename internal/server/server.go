package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	auth "github.com/go-park-mail-ru/2025_1_sigmaScript/auth_service/api/auth_api_v1/proto"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	deliveryAuth "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/delivery"
	movie "github.com/go-park-mail-ru/2025_1_sigmaScript/movie_service/api/movie_api_v1/proto"
	user "github.com/go-park-mail-ru/2025_1_sigmaScript/user_service/api/user_api_v1/proto"

	// repoAuthSessions "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/movie_repo"
	// serviceAuth "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/auth/movie_service"
	//repoUsers "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/user/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	deliveryCollection "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/collection/delivery"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/router"

	csrfDelivery "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/csrf/delivery"
	deliveryMovie "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/delivery"

	//repoMovie "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/repository"
	//serviceMovie "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/service"

	deliveryReviews "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/reviews/delivery"
	deliveryStaff "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/staff_person/delivery"

	deliveryGenre "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/genre/delivery"
	client "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/grpc_client"
	deliveryUsers "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/user/delivery/http"

	deliverySearch "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/search/delivery"

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
	log.Info().Msg("Trying to connect to auth movie_service")
	aGrpcConn, err := grpc.NewClient(
		":8081",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("error couldnt connect to grpc: %w", err)
	}
	log.Info().Msg("Auth movie_service connection opened successfully")

	defer func() {
		if clErr := aGrpcConn.Close(); clErr != nil {
			log.Error().Msg("couldn't close auth microservice grpc connection")
		}
	}()

	log.Info().Msg("Trying to connect to auth movie_service")
	mGrpcConn, err := grpc.NewClient(
		":8083",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("error couldnt connect to grpc: %w", err)
	}
	log.Info().Msg("Auth movie_service connection opened successfully")

	defer func() {
		if clErr := mGrpcConn.Close(); clErr != nil {
			log.Error().Msg("couldn't close auth microservice grpc connection")
		}
	}()

	uGrpcConn, err := grpc.NewClient(
		":8082",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("error couldnt connect to grpc: %w", err)
	}
	log.Info().Msg("Auth service connection opened successfully")

	defer func() {
		if clErr := uGrpcConn.Close(); clErr != nil {
			log.Error().Msg("couldn't close auth microservice grpc connection")
		}
	}()

	ctxPgDb := config.WrapPgDatabaseContext(context.Background(), s.Config.PostgresConfig)
	ctxPgDb, cancelPgDb := context.WithTimeout(ctxPgDb, time.Second*30)
	defer cancelPgDb()

	// pgdb, err := db.SetupDatabase(ctxPgDb, cancelPgDb)
	// if err != nil {
	// 	return fmt.Errorf("error couldnt connect to postgres database: %w", err)
	// }

	// sessionRepo := repoAuthSessions.NewSessionRepository()
	// sessionService := serviceAuth.NewSessionService(config.WrapCookieContext(context.Background(), &s.Config.Cookie), sessionRepo)

	sessionService := client.NewAuthClient(auth.NewSessionRPCClient(aGrpcConn))

	//userRepo := repoUsers.NewUserRepository(pgdb)
	//userService := serviceUsers.NewUserService(userRepo)
	userService := client.NewUserClient(user.NewUserServiceClient(uGrpcConn))
	userHandler := deliveryUsers.NewUserHandler(config.WrapCookieContext(context.Background(), &s.Config.Cookie), userService, sessionService)

	authHandler := deliveryAuth.NewAuthHandler(config.WrapCookieContext(context.Background(), &s.Config.Cookie), userService, sessionService)

	csrfHandler := csrfDelivery.NewCSRFHandler(config.WrapCookieContext(context.Background(), &s.Config.Cookie), sessionService)

	movieService := client.NewMovieClient(movie.NewMovieRPCClient(mGrpcConn))
	movieHandler := deliveryMovie.NewMovieHandler(movieService)

	//staffPersonRepo := repoStaff.NewStaffPersonPostgresRepository(pgdb)
	//staffPersonService := serviceStaff.NewStaffPersonService(staffPersonRepo)
	staffPersonHandler := deliveryStaff.NewStaffPersonHandler(movieService)

	//collectionRepo := repoCollection.NewCollectionPostgresRepository(pgdb)
	//collectionService := serviceCollection.NewCollectionService(collectionRepo)
	collectionHandler := deliveryCollection.NewCollectionHandler(movieService)

	movieReviewHandler := deliveryReviews.NewReviewHandler(userService, sessionService, movieService)

	//genreRepo := repoGenre.NewGenreRepository(pgdb)
	//genreService := serviceGenre.NewGenreService(genreRepo)
	genreHandler := deliveryGenre.NewGenreHandler(movieService)

	//searchRepo := repoSearch.NewSearchRepository(pgdb)
	//searchService := serviceSearch.NewSearchService(searchRepo)
	searchHandler := deliverySearch.NewSearchHandler(movieService)

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
