package middleware_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/igkostyuk/tasktracker/internal/middleware"
)

func TestRequestID(t *testing.T) {
	h := middleware.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, middleware.GetRequestID(r.Context()))
	}))

	t.Run("it return id from Header", func(t *testing.T) {
		want := "id123"

		w := httptest.NewRecorder()
		ctx := context.Background()
		req, err := http.NewRequestWithContext(ctx, "GET", "http://example.com/foo", nil)
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set(middleware.RequestIDHeader, want)
		h.ServeHTTP(w, req)
		got := w.Body.String()
		if want != got {
			t.Errorf("want %s got %s", want, got)
		}
	})

	t.Run("it return generated uuid id", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx := context.Background()
		req, err := http.NewRequestWithContext(ctx, "GET", "http://example.com/foo", nil)
		if err != nil {
			t.Fatal(err)
		}
		h.ServeHTTP(w, req)
		got := w.Body.String()
		_, err = uuid.Parse(got)
		if err != nil {
			t.Errorf("got %s invvalid uuid %w", got, err)
		}
	})
}

func TestGetRequestID(t *testing.T) {
	t.Run("it returns empty string if no context key", func(t *testing.T) {
		ctx := context.Background()
		got := middleware.GetRequestID(ctx)
		if got != "" {
			t.Errorf("got %s want empty string", got)
		}
	})
	t.Run("it returns request_id from context", func(t *testing.T) {
		want := "idfoo"
		ctx := context.WithValue(context.Background(), middleware.RequestIDKey, want)
		got := middleware.GetRequestID(ctx)
		if got != want {
			t.Errorf("got %s want %s", got, want)
		}
	})
}
