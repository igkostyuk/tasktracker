package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
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
		r.Put("/", handler.Update)
		r.Delete("/", handler.Delete)
		r.Get("/columns", handler.FetchColumns)
		r.Post("/columns", handler.StoreColumn)
		r.Get("/tasks", handler.FetchTasks)
	})

	return r
}

// Fetch project godoc
// @Summary Get all projects
// @Description get all projects
// @Tags projects
// @Produce  json
// @Success 200 {array} domain.Project
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /projects [get]
// Fetch will fetch projects.
func (p *projectHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	projects, err := p.projectUsecase.Fetch(r.Context())
	if err != nil {
		web.RespondError(w, r, err, getStatusCode(err))

		return
	}
	web.Respond(w, r, projects, http.StatusOK)
}

// FetchColumns godoc
// @Summary Get columns by project id
// @Description get columns by project id
// @Tags columns
// @Produce  json
// @Param  id path string true "project ID" format(uuid)
// @Success 200 {array} domain.Column
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /projects/{id}/columns [get]
// FetchColumns will fetch columns by project id.
func (p *projectHandler) FetchColumns(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		web.RespondError(w, r, domain.ErrNotFound, http.StatusNotFound)

		return
	}
	columns, err := p.projectUsecase.FetchColumns(r.Context(), id)
	if err != nil {
		web.RespondError(w, r, err, getStatusCode(err))

		return
	}
	web.Respond(w, r, columns, http.StatusOK)
}

func isColumnRequestValid(m *domain.Column) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, fmt.Errorf("validation: %w", err)
	}

	return true, nil
}

// Store godoc
// @Summary Add a column
// @Description add by json column
// @Tags columns
// @Accept  json
// @Produce  json
// @Param id path string true "project ID" format(uuid)
// @Param column body domain.Column true "Add column"
// @Success 200 {object} domain.Column
// @Failure 400 {object} web.HTTPError
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 422 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /projects/{id}/columns [post]
// Store will store the column by given request body and Project ID.
func (p *projectHandler) StoreColumn(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		web.RespondError(w, r, domain.ErrNotFound, http.StatusNotFound)

		return
	}
	var column domain.Column
	if err := json.NewDecoder(r.Body).Decode(&column); err != nil {
		web.RespondError(w, r, err, http.StatusUnprocessableEntity)

		return
	}
	column.ProjectID = id
	if ok, err := isColumnRequestValid(&column); !ok {
		web.RespondError(w, r, err, http.StatusBadRequest)

		return
	}
	if err := p.projectUsecase.StoreColumn(r.Context(), &column); err != nil {
		web.RespondError(w, r, err, getStatusCode(err))

		return
	}
	web.Respond(w, r, column, http.StatusOK)
}

// FetchTasks godoc
// @Summary Get tasks by project id
// @Description get tasks by project id
// @Tags tasks
// @Produce  json
// @Param  id path string true "project ID" format(uuid)
// @Success 200 {array} domain.Task
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /projects/{id}/tasks [get]
// FetchTasks will fetch tasks by project id.
func (p *projectHandler) FetchTasks(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		web.RespondError(w, r, domain.ErrNotFound, http.StatusNotFound)

		return
	}
	tasks, err := p.projectUsecase.FetchTasks(r.Context(), id)
	if err != nil {
		web.RespondError(w, r, err, getStatusCode(err))

		return
	}
	web.Respond(w, r, tasks, http.StatusOK)
}

// GetByID godoc
// @Summary Show a project
// @Description get project by id
// @Tags projects
// @Produce  json
// @Param  id path string true "project ID" format(uuid)
// @Success 200 {object} domain.Project
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /projects/{id} [get]
// GetByID will get project by given id.
func (p *projectHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		web.RespondError(w, r, domain.ErrNotFound, http.StatusNotFound)

		return
	}
	project, err := p.projectUsecase.GetByID(r.Context(), id)
	if err != nil {
		web.RespondError(w, r, err, getStatusCode(err))

		return
	}
	web.Respond(w, r, project, http.StatusOK)
}

func isRequestValid(m *domain.Project) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, fmt.Errorf("validation: %w", err)
	}

	return true, nil
}

// Store godoc
// @Summary Add a project
// @Description add by json project
// @Tags projects
// @Accept  json
// @Produce  json
// @Param project body domain.Project true "Add project"
// @Success 200 {object} domain.Project
// @Failure 400 {object} web.HTTPError
// @Failure 422 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /projects [post]
// Store will store the project by given request body.
func (p *projectHandler) Store(w http.ResponseWriter, r *http.Request) {
	var project domain.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		web.RespondError(w, r, err, http.StatusUnprocessableEntity)

		return
	}
	if ok, err := isRequestValid(&project); !ok {
		web.RespondError(w, r, err, http.StatusBadRequest)

		return
	}
	if err := p.projectUsecase.Store(r.Context(), &project); err != nil {
		web.RespondError(w, r, err, getStatusCode(err))

		return
	}

	web.Respond(w, r, project, http.StatusOK)
}

// Store godoc
// @Summary Update a project
// @Description update by json project
// @Tags projects
// @Accept  json
// @Produce  json
// @Param  id path string true "project ID"
// @Param project body domain.Project true "Update project"
// @Success 200 {object} domain.Project
// @Failure 400 {object} web.HTTPError
// @Failure 422 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /projects/{id} [put]
// Update will update the project by given request body.
func (p *projectHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		web.RespondError(w, r, domain.ErrNotFound, http.StatusNotFound)

		return
	}
	var project domain.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		web.RespondError(w, r, err, http.StatusUnprocessableEntity)

		return
	}
	project.ID = id
	if ok, err := isRequestValid(&project); !ok {
		web.RespondError(w, r, err, http.StatusBadRequest)

		return
	}
	if err := p.projectUsecase.Update(r.Context(), &project); err != nil {
		web.RespondError(w, r, err, getStatusCode(err))

		return
	}

	web.Respond(w, r, project, http.StatusOK)
}

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
// Delete will delete project by given param.
func (p *projectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "projectID"))
	if err != nil {
		web.RespondError(w, r, domain.ErrNotFound, http.StatusNotFound)

		return
	}
	if err := p.projectUsecase.Delete(r.Context(), id); err != nil {
		web.RespondError(w, r, err, getStatusCode(err))

		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func getStatusCode(err error) int {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, domain.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, domain.ErrUnique):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
