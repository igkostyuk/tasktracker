package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
	"github.com/igkostyuk/tasktracker/domain"
)

type taskHandler struct {
	taskUsecase domain.TaskUsecase
}

func New(us domain.TaskUsecase) chi.Router {
	handler := &taskHandler{
		taskUsecase: us,
	}
	r := chi.NewRouter()
	r.Get("/", handler.Fetch)
	r.Post("/", handler.Store)
	r.Route("/{taskID}", func(r chi.Router) {
		r.Get("/", handler.GetByID)
		r.Delete("/", handler.Delete)
		r.Get("/comments", handler.FetchComments)
	})

	return r
}

func (t *taskHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	columns, err := t.taskUsecase.Fetch(r.Context())
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))

		return
	}
	jsonData, err := json.Marshal(columns)
	if err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (t *taskHandler) FetchComments(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "taskID")
	columns, err := t.taskUsecase.FetchComments(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))

		return
	}
	jsonData, err := json.Marshal(columns)
	if err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// GetByID will get task by given id.
func (t *taskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "taskID")

	ctx := r.Context()

	task, err := t.taskUsecase.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))

		return
	}

	jsonData, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func isRequestValid(m *domain.Task) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, fmt.Errorf("validation: %w", err)
	}

	return true, nil
}

// Store will store the task by given request body.
func (t *taskHandler) Store(w http.ResponseWriter, r *http.Request) {
	var task domain.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}
	var ok bool
	if ok, err = isRequestValid(&task); !ok {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	if err = t.taskUsecase.Store(r.Context(), &task); err != nil {
		http.Error(w, err.Error(), getStatusCode(err))

		return
	}
	jsonData, err := json.Marshal(task)
	if err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

// Delete will delete project by given param.
func (t *taskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "taskID")

	ctx := r.Context()

	err := t.taskUsecase.Delete(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, domain.ErrLastColumn):
		return http.StatusConflict
	case errors.Is(err, domain.ErrColumnName):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
