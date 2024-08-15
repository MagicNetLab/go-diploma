package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log = Logger{log: zap.NewNop()}

func Initialize() error {
	cnf := zap.NewProductionEncoderConfig()
	cnf.EncodeTime = zapcore.ISO8601TimeEncoder
	cnf.TimeKey = "timestamp"

	zl, err := zap.NewProduction()
	if err != nil {
		return err
	}

	log = Logger{log: zl}

	return nil
}

func Info(msg string, args map[string]interface{}) {
	log.Info(msg, args)
}
func Error(msg string, args map[string]interface{}) {
	log.Error(msg, args)
}
func Debug(msg string, args map[string]interface{}) {
	log.Debug(msg, args)
}
func Fatal(msg string, args map[string]interface{}) {
	log.Fatal(msg, args)
}

func Sync() {
	log.log.Sync()
}
