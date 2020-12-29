package middleware

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapStructuredLogger struct {
	Logger *zap.Logger
}

func (l *zapStructuredLogger) NewLogEntry(r *http.Request) LogEntry {
	fields := []zapcore.Field{}
	if reqID := GetRequestID(r.Context()); reqID != "" {
		fields = append(fields, zap.String("request-id", reqID))
	}
	fields = append(fields,
		zap.String("method", r.Method),
		zap.String("request", r.RequestURI),
		zap.String("remote", r.RemoteAddr),
	)
	l.Logger.Info("request started", fields...)

	return &zapStructuredLoggerEntry{Logger: l.Logger, fields: fields}
}

type zapStructuredLoggerEntry struct {
	Logger *zap.Logger
	fields []zapcore.Field
}

func (l *zapStructuredLoggerEntry) Write(status int, header http.Header, elapsed time.Duration, extra interface{}) {
	fields := append(l.fields,
		zap.Int("status", status),
		zap.Float64("elapsed ms", float64(elapsed.Nanoseconds())/1000000.0))

	l.Logger.Info("request complete", fields...)
}

func (l *zapStructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	fields := append(
		l.fields,
		zap.String("stack", string(stack)),
		zap.String("panic", fmt.Sprintf("%+v", v)))

	l.Logger.Error("panic", fields...)
}
