package logger

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Logger struct {
	log *zap.Logger
}

func (l *Logger) Info(msg string, args map[string]interface{}) {
	l.log.Info(msg, l.prepareArgs(args)...)
}

func (l *Logger) Error(msg string, args map[string]interface{}) {
	l.log.Error(msg, l.prepareArgs(args)...)
}

func (l *Logger) Fatal(msg string, args map[string]interface{}) {
	l.log.Fatal(msg, l.prepareArgs(args)...)
}

func (l *Logger) prepareArgs(args map[string]interface{}) []zap.Field {
	var r []zap.Field
	for k, v := range args {
		switch t := v.(type) {
		case string:
			r = append(r, zap.String(k, t))
		case int:
			r = append(r, zap.Int(k, t))
		case time.Duration:
			r = append(r, zap.Duration(k, t))
		default:
			fmt.Printf("undefuned logger params type %s", t)
		}
	}

	return r
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
