package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-park-mail-ru/2025_1_sigmaScript/config"
	errs "github.com/go-park-mail-ru/2025_1_sigmaScript/internal/errors"
	"github.com/go-park-mail-ru/2025_1_sigmaScript/internal/server"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal().Err(errors.Wrap(err, errs.ErrLoadConfig)).Msg(errors.Wrap(err, errs.ErrLoadConfig).Error())
	}

	srv := server.New(cfg)
	log.Info().Msg("Starting server")

	go func() {
		if err = srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(errors.Wrap(err, errs.ErrStartServer)).Msg(errors.Wrap(err, errs.ErrStartServer).Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Info().Msg("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(errors.Wrap(err, errs.ErrShutdown)).Msg(errors.Wrap(err, errs.ErrShutdown).Error())
	}
	log.Info().Msg("Server is shut down gracefully")
}
