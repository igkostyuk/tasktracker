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

type taskHandler struct {
	taskUsecase domain.TaskUsecase
}

// New return routes for task resource.
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

// Fetch tasks godoc
// @Summary Get all tasks
// @Description get all tasks
// @Tags tasks
// @Produce  json
// @Success 200 {array} domain.Task
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /tasks [get]
// Fetch will fetch tasks.
func (t *taskHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	tasks, err := t.taskUsecase.Fetch(r.Context())
	if err != nil {
		web.RespondError(w, err, getStatusCode(err))

		return
	}
	web.Respond(w, tasks, http.StatusOK)
}

// FetchComments godoc
// @Summary Get comments by task id
// @Description get tasks by task id
// @Tags tasks
// @Produce  json
// @Param  id path string true "task ID" format(uuid)
// @Success 200 {array} domain.Task
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /tasks/{id}/comments [get]
// FetchComments will fetch tasks by task id.
func (t *taskHandler) FetchComments(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "taskID"))
	if err != nil {
		web.RespondError(w, domain.ErrNotFound, http.StatusNotFound)

		return
	}
	comments, err := t.taskUsecase.FetchComments(r.Context(), id)
	if err != nil {
		web.RespondError(w, err, getStatusCode(err))

		return
	}
	web.Respond(w, comments, http.StatusOK)
}

// GetByID godoc
// @Summary Show a task
// @Description get task by id
// @Tags tasks
// @Produce  json
// @Param  id path string true "task ID" format(uuid)
// @Success 200 {object} domain.Task
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /tasks/{id} [get]
// GetByID will get task by given id.
func (t *taskHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "taskID"))
	if err != nil {
		web.RespondError(w, domain.ErrNotFound, http.StatusNotFound)

		return
	}
	task, err := t.taskUsecase.GetByID(r.Context(), id)
	if err != nil {
		web.RespondError(w, err, getStatusCode(err))

		return
	}
	web.Respond(w, task, http.StatusOK)
}

func isRequestValid(m *domain.Task) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, fmt.Errorf("validation: %w", err)
	}

	return true, nil
}

// Store godoc
// @Summary Add a task
// @Description add by json task
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param project body domain.Task true "Add task"
// @Success 200 {object} domain.Task
// @Failure 400 {object} web.HTTPError
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 422 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /tasks [post]
// Store will store the task by given request body.
func (t *taskHandler) Store(w http.ResponseWriter, r *http.Request) {
	var task domain.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		web.RespondError(w, err, http.StatusUnprocessableEntity)

		return
	}
	if ok, err := isRequestValid(&task); !ok {
		web.RespondError(w, err, http.StatusBadRequest)

		return
	}
	if err := t.taskUsecase.Store(r.Context(), &task); err != nil {
		web.RespondError(w, err, getStatusCode(err))

		return
	}
	web.Respond(w, task, http.StatusOK)
}

// Delete godoc
// @Summary Delete a task
// @Description Delete by task ID
// @Tags tasks
// @Produce  json
// @Param  id path string true "task ID"
// @Success 204 "it's ok"
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /tasks/{id} [delete]
// Delete will delete task by given param.
func (t *taskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "taskID"))
	if err != nil {
		web.RespondError(w, domain.ErrNotFound, http.StatusNotFound)

		return
	}
	if err := t.taskUsecase.Delete(r.Context(), id); err != nil {
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
	default:
		return http.StatusInternalServerError
	}
}
