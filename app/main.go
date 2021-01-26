package main

import (

	// Register the expvar handler.
	"context"
	"expvar"
	"fmt"
	"log"
	"net/http"

	// nolint:gosec //Register the pprof handlers on other port than api.
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/igkostyuk/tasktracker/app/server"
	"github.com/igkostyuk/tasktracker/configs"
	"github.com/igkostyuk/tasktracker/store/postgres"
	"go.uber.org/zap"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "development"

func main() {
	configPath := ""

	lgr, err := configs.GetLoggerConfig(build).Build()
	if err != nil {
		log.Fatal("building logger", err)
	}

	if err := run(configPath, lgr); err != nil {
		lgr.Error("main: error:", zap.Error(err))
		os.Exit(1)
	}
}

//nolint:funlen
func run(cfgPath string, logger *zap.Logger) error {
	// =========================================================================
	// Configuration
	cfg, err := configs.FromFile(cfgPath)
	if err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}
	// =========================================================================
	// App Starting
	// Print the build version for our logs. Also expose it under /debug/vars.
	expvar.NewString("build").Set(build)
	log.Printf("main : Started : Application initializing : version %q", build)
	defer logger.Info("main: Completed")
	// =========================================================================
	// Start Database
	logger.Info("main: Initializing database support")
	db, err := postgres.Open(cfg.Postgres)
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
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
	// =========================================================================
	// Start API Service
	logger.Info("main: Initializing API support")
	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)
	api := server.New(cfg, logger, db)
	serverErrors := make(chan error, 1)
	// Start the service listening for requests.
	go func() {
		logger.Sugar().Infof("main: API listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()
	// =========================================================================
	// Shutdown
	select { // Blocking main and waiting for shutdown.
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		logger.Sugar().Infof("main: %v : Start shutdown", sig)
		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		// Asking listener to shutdown and shed load.
		if err := api.Shutdown(ctx); err != nil {
			api.Close()

			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
