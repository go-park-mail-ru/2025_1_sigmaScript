package router

import (
	"context"
	"testing"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/handlers"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/mocks"
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
	router := NewRouter()
	require.NotNil(t, router)

	cfg, err := config.New()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	authHandler := handlers.NewAuthHandler(config.WrapCookieContext(context.Background(), cfg.Cookie))
	require.NotEmpty(t, authHandler)

	collectionRepo := repoCollection.NewCollectionRepository(&mocks.MainPageCollections)
	collectionService := serviceCollection.NewCollectionService(collectionRepo)
	collectionHandler := deliveryCollection.NewCollectionHandler(collectionService)

	ApplyMiddlewares(router)
	SetupAuth(router, authHandler)
	SetupCollections(router, collectionHandler)
}
