package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/igkostyuk/tasktracker/app/config"
	"github.com/igkostyuk/tasktracker/app/handlers"
	zapLogger "github.com/igkostyuk/tasktracker/app/logger"
	"github.com/igkostyuk/tasktracker/store/postgres"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "development"

func main() {
	zl, err := zapLogger.New(build)
	if err != nil {
		log.Fatal("building logger", err)
	}
	c, err := config.FromFile("")
	if err != nil {
		log.Fatal("parsing config", err)
	}

	if err := run(c, zl); err != nil {
		zl.Error("main: error:", zap.Error(err))
		os.Exit(1)
	}
}

func run(cfg config.Config, logger *zap.Logger) error {
	logger.Info("main: Started : Application initializing.")
	defer logger.Info("main: Completed")
	// =========================================================================
	logger.Info("main: Initializing database support")
	db, err := postgres.Open(cfg.Postgres)
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}
	defer func() {
		log.Printf("main: Database Stopping")
		db.Close()
	}()
	// =========================================================================
	// Start Debug Service
	// /debug/pprof - Added to the default mux by importing the net/http/pprof package.
	// /debug/vars - Added to the default mux by importing the expvar package.
	logger.Info("main: Initializing debugging support")
	go func() { // Not concerned with shutting this down when the application is shutdown.
		logger.Sugar().Infof("main: Debug linstening: %s", cfg.DebugHost)
		if err := http.ListenAndServe(cfg.DebugHost, http.DefaultServeMux); err != nil {
			logger.Error("main: Debug Listener closed :", zap.Error(err))
		}
	}()
	logger.Info("main: Initializing API support")
	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// nolint:exhaustivestruct
	api := http.Server{
		Addr:         cfg.APIHost,
		Handler:      handlers.API(logger),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
	serverErrors := make(chan error, 1)
	// Start the service listening for requests.
	go func() {
		logger.Sugar().Infof("main: API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()
	// =========================================================================
	select { // Blocking main and waiting for shutdown.
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")
	case sig := <-shutdown:
		logger.Sugar().Infof("main: %v : Start shutdown", sig)
		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		// Asking listener to shutdown and shed load.
		if err := api.Shutdown(ctx); err != nil {
			api.Close()

			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
