package server

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/igkostyuk/tasktracker/configs"
	"github.com/igkostyuk/tasktracker/internal/middleware"
	projectDelivery "github.com/igkostyuk/tasktracker/project/delivery/http"
	projectRepository "github.com/igkostyuk/tasktracker/project/repository/postgres"
	projectUsecase "github.com/igkostyuk/tasktracker/project/usecase"
	"go.uber.org/zap"
)

func New(cfg configs.Config, logger *zap.Logger, db *sql.DB) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.NewZaplogger(logger))
	r.Use(middleware.Recoverer)

	r.Mount("/projects", projectDelivery.New(projectUsecase.New(projectRepository.New(db))))

	// nolint:exhaustivestruct
	return &http.Server{
		Addr:         cfg.APIHost,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}
