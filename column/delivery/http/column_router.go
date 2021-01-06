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

type columnHandler struct {
	columnUsecase domain.ColumnUsecase
}

// New return routes for column resource.
func New(us domain.ColumnUsecase) chi.Router {
	handler := &columnHandler{
		columnUsecase: us,
	}
	r := chi.NewRouter()
	r.Get("/", handler.Fetch)
	r.Post("/", handler.Store)
	r.Route("/{columnID}", func(r chi.Router) {
		r.Get("/", handler.GetByID)
		r.Delete("/", handler.Delete)
		r.Get("/tasks", handler.FetchTasks)
	})

	return r
}

func (c *columnHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	columns, err := c.columnUsecase.Fetch(r.Context())
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

func (c *columnHandler) FetchTasks(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "columnID")
	tasks, err := c.columnUsecase.FetchTasks(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))

		return
	}
	jsonData, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// GetByID will get project by given id.
func (c *columnHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "columnID")

	ctx := r.Context()

	column, err := c.columnUsecase.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))

		return
	}

	jsonData, err := json.Marshal(column)
	if err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func isRequestValid(m *domain.Column) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, fmt.Errorf("validation: %w", err)
	}

	return true, nil
}

// Store will store the project by given request body.
func (c *columnHandler) Store(w http.ResponseWriter, r *http.Request) {
	var column domain.Column
	err := json.NewDecoder(r.Body).Decode(&column)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}
	var ok bool
	if ok, err = isRequestValid(&column); !ok {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	if err := c.columnUsecase.Store(r.Context(), &column); err != nil {
		http.Error(w, err.Error(), getStatusCode(err))

		return
	}
	jsonData, err := json.Marshal(column)
	if err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

// Delete will delete project by given param.
func (c *columnHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "columnID")

	ctx := r.Context()

	err := c.columnUsecase.Delete(ctx, id)
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
