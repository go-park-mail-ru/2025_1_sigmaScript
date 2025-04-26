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

	deliveryMovie "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/delivery"
	repoMovie "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/repository"
	serviceMovie "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/movie/service"

	deliveryCSAT "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/csat/delivery"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/csat/delivery/dto"
	repoCSAT "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/csat/repository"
	serviceCSAT "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/csat/service"

	csrfDelivery "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/csrf/delivery"
	deliveryReviews "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/reviews/delivery"
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
	// // TODO fix config: it`s test database test password
	// postgres := config.Postgres{
	// 	Host:            "127.0.0.1",
	// 	Port:            5433,
	// 	User:            "filmlk_user",
	// 	Password:        "filmlk_password",
	// 	Name:            "filmlk",
	// 	MaxOpenConns:    100,
	// 	MaxIdleConns:    30,
	// 	ConnMaxLifetime: 30,
	// 	ConnMaxIdleTime: 5,
	// }

	// avatarLocalStorage := config.LocalAvatarsStorage{
	// 	UserAvatarsFullPath:     "",
	// 	UserAvatarsRelativePath: "",
	// }

	// pgDatabase := config.Databases{
	// 	Postgres:     postgres,
	// 	LocalStorage: avatarLocalStorage,
	// }

	// pgListener := config.Listener{
	// 	Port: "5433",
	// }

	// cfgDB := config.ConfigPgDB{
	// 	Listener:  pgListener,
	// 	Databases: pgDatabase,
	// }

	// ctxDb := config.WrapPgDatabaseContext(context.Background(), cfgDB)
	// ctxDb, cancel := context.WithTimeout(ctxDb, time.Second*30)
	// defer cancel()

	// pgdb, err := db.SetupDatabase(ctxDb, cancel)
	// if err != nil {
	// 	return fmt.Errorf("error couldnt connect to postgres database: %w", err)
	// }

	sessionRepo := repoAuthSessions.NewSessionRepository()
	sessionService := serviceAuth.NewSessionService(config.WrapCookieContext(context.Background(), &s.Config.Cookie), sessionRepo)

	userRepo := repoUsers.NewUserRepository(nil)
	userService := serviceUsers.NewUserService(userRepo)
	userHandler := deliveryUsers.NewUserHandler(config.WrapCookieContext(context.Background(), &s.Config.Cookie), userService, sessionService)

	authHandler := deliveryAuth.NewAuthHandler(config.WrapCookieContext(context.Background(), &s.Config.Cookie), userService, sessionService)

	csrfHandler := csrfDelivery.NewCSRFHandler(config.WrapCookieContext(context.Background(), &s.Config.Cookie), sessionService)

	staffPersonRepo := repoStaff.NewStaffPersonRepository(&mocks.ExistingActors)
	staffPersonService := serviceStaff.NewStaffPersonService(staffPersonRepo)
	staffPersonHandler := deliveryStaff.NewStaffPersonHandler(staffPersonService)

	collectionRepo := repoCollection.NewCollectionRepository(&mocks.MainPageCollections)
	collectionService := serviceCollection.NewCollectionService(collectionRepo)
	collectionHandler := deliveryCollection.NewCollectionHandler(collectionService)

	movieRepo := repoMovie.NewMovieRepository(&mocks.ExistingMovies)
	movieService := serviceMovie.NewMovieService(movieRepo)
	movieHandler := deliveryMovie.NewMovieHandler(movieService)

	movieReviewHandler := deliveryReviews.NewReviewHandler(userService, sessionService, movieService)

	csatRepo := repoCSAT.NewCSATRepository(&mocks.CSATRepo{
		Rating:  0,
		Reviews: make(map[int]*dto.CSATReviewDataJSON),
	})
	csatService := serviceCSAT.NewCSATService(csatRepo)
	csatHandler := deliveryCSAT.NewCSATHandler(userService, sessionService, csatService)

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

	router.SetupCSATReviewsHandlers(mx, csatHandler)

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
