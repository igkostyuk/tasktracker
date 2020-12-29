package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

var RequestIDHeader = "X-Request-Id"

type contextKey string

const (
	requestIDKey contextKey = "request_id"
)

func GetRequestID(ctx context.Context) (value string) {
	value, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		return ""
	}

	return value
}

func WithRequestID(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, requestIDKey, value)
}

func RequestID(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestID := r.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx = WithRequestID(ctx, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
