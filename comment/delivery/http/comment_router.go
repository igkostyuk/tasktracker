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

type commentHandler struct {
	commentUsecase domain.CommentUsecase
}

func New(us domain.CommentUsecase) chi.Router {
	handler := &commentHandler{
		commentUsecase: us,
	}
	r := chi.NewRouter()
	r.Get("/", handler.Fetch)
	r.Post("/", handler.Store)
	r.Route("/{commmentID}", func(r chi.Router) {
		r.Get("/", handler.GetByID)
		r.Delete("/", handler.Delete)
	})

	return r
}

func (c *commentHandler) Fetch(w http.ResponseWriter, r *http.Request) {
	columns, err := c.commentUsecase.Fetch(r.Context())
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

// GetByID will get comment by given id.
func (c *commentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "commentID")

	ctx := r.Context()

	comment, err := c.commentUsecase.GetByID(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))

		return
	}

	jsonData, err := json.Marshal(comment)
	if err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func isRequestValid(m *domain.Comment) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, fmt.Errorf("validation: %w", err)
	}

	return true, nil
}

// Store will store the comment by given request body.
func (c *commentHandler) Store(w http.ResponseWriter, r *http.Request) {
	var comment domain.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}
	var ok bool
	if ok, err = isRequestValid(&comment); !ok {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	if err := c.commentUsecase.Store(r.Context(), &comment); err != nil {
		http.Error(w, err.Error(), getStatusCode(err))

		return
	}
	jsonData, err := json.Marshal(comment)
	if err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)

		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

// Delete will delete comment by given param.
func (c *commentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "commentID")

	ctx := r.Context()

	err := c.commentUsecase.Delete(ctx, id)
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
