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

type projectHandler struct {
	projectUsecase domain.ProjectUsecase
}

// New return routes for project resource.
func New(us domain.ProjectUsecase) chi.Router {
	handler := &projectHandler{
		projectUsecase: us,
	}
	r := chi.NewRouter()
	r.Get("/", handler.Fetch)
	r.Post("/", handler.Store)
	r.Route("/{projectID}", func(r chi.Router) {
		r.Get("/", handler.GetByID)
		r.Delete("/", handler.Delete)
		r.Get("/columns", handler.FetchColumns)
		r.Get("/tasks", handler.FetchTasks)
	})

	return r
}

func (p *projectHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	projects, err := p.projectUsecase.Fetch(r.Context())
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))

		return
	}
	jsonData, err := json.Marshal(projects)
	if err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (p *projectHandler) FetchColumns(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "projectID")
	columns, err := p.projectUsecase.FetchColumns(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))

		return
	}
	jsonData, err := json.Marshal(columns)
	if err != nil {
		http.Error(w, "columns encoding error", http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (p *projectHandler) FetchTasks(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "projectID")
	tasks, err := p.projectUsecase.FetchTasks(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))

		return
	}
	jsonData, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, "columns encoding error", http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// GetByID will get project by given id.
func (p *projectHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "projectID")
	project, err := p.projectUsecase.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))

		return
	}

	jsonData, err := json.Marshal(project)
	if err != nil {
		http.Error(w, "project encoding error", http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func isRequestValid(m *domain.Project) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, fmt.Errorf("validation: %w", err)
	}

	return true, nil
}

// Store will store the project by given request body.
func (p *projectHandler) Store(w http.ResponseWriter, r *http.Request) {
	var project domain.Project
	err := json.NewDecoder(r.Body).Decode(&project)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	var ok bool
	if ok, err = isRequestValid(&project); !ok {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	ctx := r.Context()
	err = p.projectUsecase.Store(ctx, &project)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))

		return
	}

	jsonData, err := json.Marshal(project)
	if err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

// Delete will delete project by given param.
func (p *projectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "projectID")
	err := p.projectUsecase.Delete(r.Context(), id)
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
	default:
		return http.StatusInternalServerError
	}
}
