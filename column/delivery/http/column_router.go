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
		web.RespondError(w, err, getStatusCode(err))

		return
	}
	web.Respond(w, columns, http.StatusOK)
}

// FetchTasks godoc
// @Summary Get tasks by column id
// @Description get tasks by column id
// @Tags columns
// @Produce  json
// @Param  id path string true "column ID" format(uuid)
// @Success 200 {array} domain.Task
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /columns/{id}/tasks [get]
// FetchTasks will fetch tasks by column id.
func (c *columnHandler) FetchTasks(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "columnID")
	tasks, err := c.columnUsecase.FetchTasks(r.Context(), id)
	if err != nil {
		web.RespondError(w, err, getStatusCode(err))

		return
	}
	web.Respond(w, tasks, http.StatusOK)
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
	id := chi.URLParam(r, "columnID")
	column, err := c.columnUsecase.GetByID(r.Context(), id)
	if err != nil {
		web.RespondError(w, err, getStatusCode(err))

		return
	}
	web.Respond(w, column, http.StatusOK)
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
// @Summary Add a column
// @Description add by json column
// @Tags columns
// @Accept  json
// @Produce  json
// @Param project body domain.Column true "Add column"
// @Success 200 {object} domain.Column
// @Failure 400 {object} web.HTTPError
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 422 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /columns [post]
// Store will store the column by given request body.
func (c *columnHandler) Store(w http.ResponseWriter, r *http.Request) {
	var column domain.Column
	if err := json.NewDecoder(r.Body).Decode(&column); err != nil {
		web.RespondError(w, err, http.StatusUnprocessableEntity)

		return
	}
	if ok, err := isRequestValid(&column); !ok {
		web.RespondError(w, err, http.StatusBadRequest)

		return
	}
	if err := c.columnUsecase.Store(r.Context(), &column); err != nil {
		web.RespondError(w, err, getStatusCode(err))

		return
	}
	web.Respond(w, column, http.StatusOK)
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
	id := chi.URLParam(r, "columnID")
	if err := c.columnUsecase.Delete(r.Context(), id); err != nil {
		web.RespondError(w, err, getStatusCode(err))

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
