package router

import (
  "context"
  "testing"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/config"
  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server/handlers"
  "github.com/stretchr/testify/require"
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

  ApplyMiddlewares(router)
  SetupAuth(router, authHandler)
  SetupCollections(router)
}
