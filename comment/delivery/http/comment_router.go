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

type commentHandler struct {
	commentUsecase domain.CommentUsecase
}

// New return routes for comment resource.
func New(us domain.CommentUsecase) chi.Router {
	handler := &commentHandler{
		commentUsecase: us,
	}
	r := chi.NewRouter()
	r.Get("/", handler.Fetch)
	r.Post("/", handler.Store)
	r.Route("/{commentID}", func(r chi.Router) {
		r.Get("/", handler.GetByID)
		r.Delete("/", handler.Delete)
	})

	return r
}

// Fetch comments godoc
// @Summary Get all comments
// @Description get all comments
// @Tags comments
// @Produce  json
// @Success 200 {array} domain.Comment
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /comments [get]
// Fetch will fetch comments.
func (c *commentHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	comments, err := c.commentUsecase.Fetch(r.Context())
	if err != nil {
		web.RespondError(w, err, getStatusCode(err))

		return
	}
	web.Respond(w, comments, http.StatusOK)
}

// GetByID godoc
// @Summary Show a comment
// @Description get comment by id
// @Tags comments
// @Produce  json
// @Param  id path string true "comment ID" format(uuid)
// @Success 200 {object} domain.Comment
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /comments/{id} [get]
// GetByID will get comment by given id.
func (c *commentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "commentID"))
	if err != nil {
		web.RespondError(w, domain.ErrNotFound, http.StatusNotFound)

		return
	}
	comment, err := c.commentUsecase.GetByID(r.Context(), id)
	if err != nil {
		web.RespondError(w, err, getStatusCode(err))

		return
	}
	web.Respond(w, comment, http.StatusOK)
}

func isRequestValid(m *domain.Comment) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, fmt.Errorf("validation: %w", err)
	}

	return true, nil
}

// Store godoc
// @Summary Add a comment
// @Description add by json comment
// @Tags comments
// @Accept  json
// @Produce  json
// @Param project body domain.Comment true "Add comment"
// @Success 200 {object} domain.Comment
// @Failure 400 {object} web.HTTPError
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 422 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /comments [post]
// Store will store the comment by given request body.
func (c *commentHandler) Store(w http.ResponseWriter, r *http.Request) {
	var comment domain.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		web.RespondError(w, err, http.StatusUnprocessableEntity)

		return
	}
	if ok, err := isRequestValid(&comment); !ok {
		web.RespondError(w, err, http.StatusBadRequest)

		return
	}
	if err := c.commentUsecase.Store(r.Context(), &comment); err != nil {
		web.RespondError(w, err, getStatusCode(err))

		return
	}
	web.Respond(w, comment, http.StatusOK)
}

// Delete godoc
// @Summary Delete a comment
// @Description Delete by comment ID
// @Tags comments
// @Produce  json
// @Param  id path string true "comment ID"
// @Success 204 "it's ok"
// @Failure 404 {object} web.HTTPError
// @Failure 409 {object} web.HTTPError
// @Failure 500 {object} web.HTTPError
// @Router /comments/{id} [delete]
// Delete will delete comment by given param.
func (c *commentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "commentID"))
	if err != nil {
		web.RespondError(w, domain.ErrNotFound, http.StatusNotFound)

		return
	}
	if err := c.commentUsecase.Delete(r.Context(), id); err != nil {
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
