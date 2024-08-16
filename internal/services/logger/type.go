package logger

import (
	"net/http"

	"go.uber.org/zap"
)

type Logger struct {
	log *zap.Logger
}

func (l *Logger) Info(msg string, args ...zap.Field) {
	l.log.Info(msg, args...)
}

func (l *Logger) Error(msg string, args ...zap.Field) {
	l.log.Error(msg, args...)
}

func (l *Logger) Fatal(msg string, args ...zap.Field) {
	l.log.Fatal(msg, args...)
}

func (l *Logger) Sync() {
	l.log.Sync()
}

type ResponseData struct {
	Status int
	Size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *ResponseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.Size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.Status = statusCode
}
