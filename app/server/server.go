package server

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	columnDelivery "github.com/igkostyuk/tasktracker/column/delivery/http"
	columnRepository "github.com/igkostyuk/tasktracker/column/repository/postgres"
	columnUsecase "github.com/igkostyuk/tasktracker/column/usecase"
	commentDelivery "github.com/igkostyuk/tasktracker/comment/delivery/http"
	commentRepository "github.com/igkostyuk/tasktracker/comment/repository/postgres"
	commentUsecase "github.com/igkostyuk/tasktracker/comment/usecase"
	"github.com/igkostyuk/tasktracker/configs"
	"github.com/igkostyuk/tasktracker/docs"
	"github.com/igkostyuk/tasktracker/internal/middleware"
	projectDelivery "github.com/igkostyuk/tasktracker/project/delivery/http"
	projectRepository "github.com/igkostyuk/tasktracker/project/repository/postgres"
	projectUsecase "github.com/igkostyuk/tasktracker/project/usecase"
	taskDelivery "github.com/igkostyuk/tasktracker/task/delivery/http"
	taskRepository "github.com/igkostyuk/tasktracker/task/repository/postgres"
	taskUsecase "github.com/igkostyuk/tasktracker/task/usecase"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// @title Task Tracker API
// @version 1.0
// @BasePath /v1

// New constructs an http.Server with main routes.
func New(cfg configs.Config, logger *zap.Logger, db *sql.DB) *http.Server {
	r := chi.NewRouter()

	columnRepo := columnRepository.New(db)
	commentRepo := commentRepository.New(db)
	projectRepo := projectRepository.New(db)
	taskRepo := taskRepository.New(db)

	r.Route("/v1", func(r chi.Router) {
		r.Use(
			middleware.RequestID,
			middleware.NewZaplogger(logger),
			middleware.Recoverer,
		)
		r.Mount("/projects", projectDelivery.New(projectUsecase.New(projectRepo, columnRepo, taskRepo)))
		r.Mount("/columns", columnDelivery.New(columnUsecase.New(columnRepo, taskRepo)))
		r.Mount("/tasks", taskDelivery.New(taskUsecase.New(columnRepo, taskRepo, commentRepo)))
		r.Mount("/comments", commentDelivery.New(commentUsecase.New(commentRepo)))
	})

	docs.SwaggerInfo.Host = cfg.APIHost
	r.Get("/swagger/*", httpSwagger.Handler())

	// nolint:exhaustivestruct
	return &http.Server{
		Addr:         cfg.APIHost,
		Handler:      cors.AllowAll().Handler(r),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}
