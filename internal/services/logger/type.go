package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"time"
)

type Logger struct {
	log *zap.Logger
}

func (l *Logger) Info(msg string, otherArgs map[string]interface{}) {
	l.printLog(zap.InfoLevel, msg, otherArgs)
}

func (l *Logger) Error(msg string, otherArgs map[string]interface{}) {
	l.printLog(zap.ErrorLevel, msg, otherArgs)
}

func (l *Logger) Debug(msg string, otherArgs map[string]interface{}) {
	l.printLog(zap.DebugLevel, msg, otherArgs)
}

func (l *Logger) Fatal(msg string, otherArgs map[string]interface{}) {
	l.printLog(zap.FatalLevel, msg, otherArgs)
}

func (l *Logger) printLog(level zapcore.Level, msg string, otherArgs map[string]interface{}) {
	var args []zap.Field
	for name, value := range otherArgs {
		switch value.(type) {
		case string:
			args = append(args, zap.String(name, value.(string)))
		case int:
			args = append(args, zap.Int(name, value.(int)))
		case bool:
			args = append(args, zap.Bool(name, value.(bool)))
		case time.Duration:
			args = append(args, zap.Duration(name, value.(time.Duration)))
		default:
			args = append(args, zap.String(name, value.(string)))
		}

	}

	switch level {
	case zap.FatalLevel:
		l.log.Fatal(msg, args...)
	case zap.ErrorLevel:
		l.log.Error(msg, args...)
	case zap.DebugLevel:
		l.log.Debug(msg, args...)
	case zap.InfoLevel:
		l.log.Info(msg, args...)
	default:
		l.log.Info(msg, args...)
	}
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
