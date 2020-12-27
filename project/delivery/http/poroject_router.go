package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
	"github.com/igkostyuk/tasktracker/domain"
)

type projectHandler struct {
	projectUsecase domain.ProjectUsecase
}

func New(us domain.ProjectUsecase) chi.Router {
	handler := &projectHandler{
		projectUsecase: us,
	}
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Route("/{projectID}", func(r chi.Router) {
			r.Get("/", handler.GetByID)
			r.Post("/", handler.GetByID)
			r.Delete("/", handler.Delete)
		})
	})

	return r
}

// GetByID will get project by given id.
func (p *projectHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idP, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		http.Error(w, domain.ErrNotFound.Error(), http.StatusNotFound)

		return
	}

	id := int64(idP)
	ctx := r.Context()

	comm, err := p.projectUsecase.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&comm)
	if err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
	}
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(&project)
	if err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
	}
}

// Delete will delete project by given param.
func (p *projectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idP, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, domain.ErrNotFound.Error(), http.StatusNotFound)

		return
	}

	id := int64(idP)
	ctx := r.Context()

	err = p.projectUsecase.Delete(ctx, id)
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
