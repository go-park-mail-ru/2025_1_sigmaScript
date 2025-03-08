package server

import (
  "context"
  "fmt"
  "testing"
  "time"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/config"
  "github.com/stretchr/testify/require"
)

func TestAppIntegration(t *testing.T) {
  cfg, err := config.New()
  require.NoError(t, err)
  require.NotNil(t, cfg)

  srv := New(cfg)
  require.NotNil(t, srv)
  fmt.Println(cfg)
  done := make(chan struct{})

  go func() {
    defer close(done)

    err := srv.Run()
    require.Equal(t, err.Error(), "http: Server closed")
  }()

  time.Sleep(1 * time.Second)
  ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
  defer cancel()
  require.NoError(t, srv.Shutdown(ctx), "failed to shut down server")

  <-done
}
