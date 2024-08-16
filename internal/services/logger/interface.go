package logger

import "go.uber.org/zap"

type AppLogger interface {
	Info(msg string, args ...zap.Field)
	Error(msg string, args ...zap.Field)
	Fatal(msg string, args ...zap.Field)
}
