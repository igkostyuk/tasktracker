package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapMiddleware returns a new Zap Middleware handler.
func NewZapMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{Logger: logger})
}

type StructuredLogger struct {
	Logger *zap.Logger
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	fields := []zapcore.Field{
		zap.String("remote", r.RemoteAddr),
		zap.String("request", r.RequestURI),
		zap.String("method", r.Method),
	}
	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		fields = append(fields, zap.String("request-id", reqID))
	}
	l.Logger.Info("request started", fields...)

	return &StructuredLoggerEntry{Logger: l.Logger, fields: fields}
}

type StructuredLoggerEntry struct {
	Logger *zap.Logger
	fields []zapcore.Field
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	fields := append(l.fields,
		zap.Int("status", status),
		zap.Int("bytes length", bytes),
		zap.Float64("elapsed ms", float64(elapsed.Nanoseconds())/1000000.0))

	l.Logger.Info("request complete", fields...)
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	fields := append(
		l.fields,
		zap.String("stack", string(stack)),
		zap.String("panic", fmt.Sprintf("%+v", v)))

	l.Logger.Error("panic", fields...)
}
