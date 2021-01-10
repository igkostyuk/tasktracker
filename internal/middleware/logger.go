package middleware

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type logContextKey string

const (
	// LogEntryCtxKey is the context.Context key to store the request log entry.
	LogEntryCtxKey logContextKey = "LogEntry"
)

// LogFormatter initiates the beginning of a new LogEntry per request.
type LogFormatter interface {
	NewLogEntry(r *http.Request) LogEntry
}

// LogEntry records the final log when a request completes.
type LogEntry interface {
	Write(status int, header http.Header, elapsed time.Duration, extra interface{})
	WriteError(err error)
	Panic(v interface{}, stack []byte)
}

func NewZaplogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return RequestLogger(&zapStructuredLogger{Logger: logger})
}

// RequestLogger returns a logger handler using a custom LogFormatter.
func RequestLogger(f LogFormatter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := f.NewLogEntry(r)
			ww := wrapResponseWriter(w)

			t1 := time.Now()
			defer func() {
				entry.Write(ww.Status(), ww.Header(), time.Since(t1), nil)
			}()

			next.ServeHTTP(ww, WithLogEntry(r, entry))
		}

		return http.HandlerFunc(fn)
	}
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true
}

// GetLogEntry returns the in-context LogEntry for a request.
func GetLogEntry(r *http.Request) LogEntry {
	entry, _ := r.Context().Value(LogEntryCtxKey).(LogEntry)

	return entry
}

// WithLogEntry sets the in-context LogEntry for a request.
func WithLogEntry(r *http.Request, entry LogEntry) *http.Request {
	r = r.WithContext(context.WithValue(r.Context(), LogEntryCtxKey, entry))

	return r
}
