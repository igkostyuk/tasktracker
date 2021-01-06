package server

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/igkostyuk/tasktracker/configs"
	"github.com/igkostyuk/tasktracker/internal/middleware"
	"go.uber.org/zap"

	columnDelivery "github.com/igkostyuk/tasktracker/column/delivery/http"
	columnRepository "github.com/igkostyuk/tasktracker/column/repository/postgres"
	columnUsecase "github.com/igkostyuk/tasktracker/column/usecase"

	commentDelivery "github.com/igkostyuk/tasktracker/comment/delivery/http"
	commentRepository "github.com/igkostyuk/tasktracker/comment/repository/postgres"
	commentUsecase "github.com/igkostyuk/tasktracker/comment/usecase"

	projectDelivery "github.com/igkostyuk/tasktracker/project/delivery/http"
	projectRepository "github.com/igkostyuk/tasktracker/project/repository/postgres"
	projectUsecase "github.com/igkostyuk/tasktracker/project/usecase"

	taskDelivery "github.com/igkostyuk/tasktracker/task/delivery/http"
	taskRepository "github.com/igkostyuk/tasktracker/task/repository/postgres"
	taskUsecase "github.com/igkostyuk/tasktracker/task/usecase"
)

func New(cfg configs.Config, logger *zap.Logger, db *sql.DB) *http.Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.NewZaplogger(logger))
	r.Use(middleware.Recoverer)

	r.Mount("/projects", projectDelivery.New(projectUsecase.New(projectRepository.New(db), columnRepository.New(db))))
	r.Mount("/columns", columnDelivery.New(columnUsecase.New(columnRepository.New(db), taskRepository.New(db))))
	r.Mount("/tasks", taskDelivery.New(taskUsecase.New(taskRepository.New(db), commentRepository.New(db))))
	r.Mount("/comments", commentDelivery.New(commentUsecase.New(commentRepository.New(db))))

	// nolint:exhaustivestruct
	return &http.Server{
		Addr:         cfg.APIHost,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}
