package main

import (
  "context"
  "errors"
  "net/http"
  "os"
  "os/signal"
  "syscall"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/config"
  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server"
  "github.com/rs/zerolog/log"
)

func main() {
  cfg, err := config.New()
  if err != nil {
    log.Fatal().Err(err).
      Str("package", "main").
      Str("func", "main").
      Msg("Error loading config")
  }

  srv := server.New(&cfg.Server)
  log.Info().
    Str("package", "main").
    Str("func", "main").
    Str("address", cfg.Server.Address).
    Int("port", cfg.Server.Port).
    Msg("Starting server")

  go func() {
    if err = srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
      log.Fatal().
        Err(err).
        Str("package", "main").
        Str("func", "main").
        Msg("Error starting server")
    }
  }()

  stop := make(chan os.Signal, 1)
  signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

  <-stop
  log.Info().
    Str("package", "main").
    Str("func", "main").
    Msg("Server is shutting down...")

  ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
  defer cancel()

  if err = srv.Shutdown(ctx); err != nil {
    log.Fatal().
      Err(err).
      Str("package", "main").
      Str("func", "main").
      Msg("Error shutting down")
  }
  log.Info().
    Str("package", "main").
    Str("func", "main").
    Msg("Server is shut down gracefully")
}
