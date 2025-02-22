package main

import (
  "context"
  "errors"
  "log"
  "net/http"
  "os"
  "os/signal"
  "syscall"
  "time"

  "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server"
)

func main() {
  srv := server.New()
  log.Println("Server is starting on :8080...")

  go func() {
    if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
      log.Fatalf("Server failed to start: %v", err)
    }
  }()

  stop := make(chan os.Signal, 1)
  signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

  <-stop
  log.Println("Server is shutting down...")

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
  defer cancel()

  if err := srv.Shutdown(ctx); err != nil {
    log.Fatalf("Server shutdown failed: %v", err)
  }
  log.Println("Server is shut down gracefully")
}
