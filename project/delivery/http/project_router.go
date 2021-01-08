package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
	"github.com/igkostyuk/tasktracker/domain"
	"github.com/igkostyuk/tasktracker/internal/web"
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

// Fetch godoc
// @Summary List projects
// @Description fetch projects
// @Tags projects
// @Produce  json
// @Success 200 {array} domain.Project
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /projects [get]
func (p *projectHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	projects, err := p.projectUsecase.Fetch(r.Context())
	if err != nil {
		web.RespondError(w, err, getStatusCode(err))

		return
	}
	web.Respond(w, projects, http.StatusOK)
}

// FetchColumns godoc
// @Summary List columns
// @Description fetch columns by project id
// @Tags projects
// @Produce  json
// @Param  id path string true "project ID"
// @Success 200 {array} domain.Column
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /projects/{id}/columns [get]
func (p *projectHandler) FetchColumns(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "projectID")
	columns, err := p.projectUsecase.FetchColumns(r.Context(), id)
	if err != nil {
		web.RespondError(w, err, getStatusCode(err))

		return
	}
	web.Respond(w, columns, http.StatusOK)
}

// FetchTasks godoc
// @Summary List tasks
// @Description fetch tasks by project id
// @Tags projects
// @Produce  json
// @Param  id path string true "project ID"
// @Success 200 {array} domain.Task
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /projects/{id}/tasks [get]
func (p *projectHandler) FetchTasks(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "projectID")
	tasks, err := p.projectUsecase.FetchTasks(r.Context(), id)
	if err != nil {
		web.RespondError(w, err, getStatusCode(err))

		return
	}
	web.Respond(w, tasks, http.StatusOK)
}

// GetByID will get project by given id.
// GetByID godoc
// @Summary Show a account
// @Description get project by id
// @Tags projects
// @Produce  json
// @Param  id path string true "project ID"
// @Success 200 {object} domain.Project
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /projects/{id} [get]
func (p *projectHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "projectID")
	project, err := p.projectUsecase.GetByID(r.Context(), id)
	if err != nil {
		web.RespondError(w, err, getStatusCode(err))

		return
	}
	web.Respond(w, project, http.StatusOK)
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
// Store godoc
// @Summary Add an project
// @Description add by json project
// @Tags projects
// @Accept  json
// @Produce  json
// @Param account body domain.Project true "Add project"
// @Success 200 {object} domain.Project
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /projects [post]
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
		web.RespondError(w, err, getStatusCode(err))

		return
	}
	web.Respond(w, project, http.StatusOK)
}

// Delete will delete project by given param.
// Delete godoc
// @Summary Delete a project
// @Description Delete by project ID
// @Tags projects
// @Produce  json
// @Param  id path string true "project ID"
// @Success 204 "it's ok"
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /projects/{id} [delete]
func (p *projectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "projectID")
	err := p.projectUsecase.Delete(r.Context(), id)
	if err != nil {
		web.RespondError(w, err, getStatusCode(err))

		return
	}
	web.Respond(w, map[string]interface{}{}, http.StatusNoContent)
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
