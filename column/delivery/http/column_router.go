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
	r.Route("/{columnID}", func(r chi.Router) {
		r.Get("/", handler.GetByID)
		r.Put("/", handler.Update)
		r.Delete("/", handler.Delete)
		r.Get("/tasks", handler.FetchTasks)
	})

	return r
}

// Fetch column godoc
// @Summary Get all columns
// @Description get all columns
// @Tags columns
// @Produce  json
// @Success 200 {array} domain.Column
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /columns [get]
// Fetch will fetch columns.
func (c *columnHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	columns, err := c.columnUsecase.Fetch(r.Context())
	if err != nil {
		web.RespondError(w, r, err, getStatusCode(err))

		return
	}
	web.Respond(w, r, columns, http.StatusOK)
}

// FetchTasks godoc
// @Summary Get tasks by column id
// @Description get tasks by column id
// @Tags tasks
// @Produce  json
// @Param  id path string true "column ID" format(uuid)
// @Success 200 {array} domain.Task
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /columns/{id}/tasks [get]
// FetchTasks will fetch tasks by column id.
func (c *columnHandler) FetchTasks(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "columnID"))
	if err != nil {
		web.RespondError(w, r, domain.ErrNotFound, http.StatusNotFound)

		return
	}
	tasks, err := c.columnUsecase.FetchTasks(r.Context(), id)
	if err != nil {
		web.RespondError(w, r, err, getStatusCode(err))

		return
	}
	web.Respond(w, r, tasks, http.StatusOK)
}

// GetByID godoc
// @Summary Show a column
// @Description get column by id
// @Tags columns
// @Produce  json
// @Param  id path string true "column ID" format(uuid)
// @Success 200 {object} domain.Column
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /columns/{id} [get]
// GetByID will get column by given id.
func (c *columnHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "columnID"))
	if err != nil {
		web.RespondError(w, r, domain.ErrNotFound, http.StatusNotFound)

		return
	}
	column, err := c.columnUsecase.GetByID(r.Context(), id)
	if err != nil {
		web.RespondError(w, r, err, getStatusCode(err))

		return
	}
	web.Respond(w, r, column, http.StatusOK)
}

func isRequestValid(m *domain.Column) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, fmt.Errorf("validation: %w", err)
	}

	return true, nil
}

// Store godoc
// @Summary Update a column
// @Description update by json column
// @Tags columns
// @Accept  json
// @Produce  json
// @Param  id path string true "column ID"
// @Param column body domain.Column true "Update column"
// @Success 200 {object} domain.Column
// @Failure 400 {object} web.HTTPError
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 422 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /columns/{id} [put]
// Update will store the column by given request body.
func (c *columnHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "columnID"))
	if err != nil {
		web.RespondError(w, r, domain.ErrNotFound, http.StatusNotFound)

		return
	}
	var column domain.Column
	if err := json.NewDecoder(r.Body).Decode(&column); err != nil {
		web.RespondError(w, r, err, http.StatusUnprocessableEntity)

		return
	}
	column.ID = id
	if ok, err := isRequestValid(&column); !ok {
		web.RespondError(w, r, err, http.StatusBadRequest)

		return
	}
	if err := c.columnUsecase.Update(r.Context(), &column); err != nil {
		web.RespondError(w, r, err, getStatusCode(err))

		return
	}
	web.Respond(w, r, column, http.StatusOK)
}

// Delete godoc
// @Summary Delete a column
// @Description Delete by column ID
// @Tags columns
// @Produce  json
// @Param  id path string true "column ID"
// @Success 204 "it's ok"
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /columns/{id} [delete]
// Delete will delete column by given param.
func (c *columnHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "columnID"))
	if err != nil {
		web.RespondError(w, r, domain.ErrNotFound, http.StatusNotFound)

		return
	}
	if err := c.columnUsecase.Delete(r.Context(), id); err != nil {
		web.RespondError(w, r, err, getStatusCode(err))

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
	case errors.Is(err, domain.ErrUnique):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
