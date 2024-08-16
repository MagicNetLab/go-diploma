package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var log = Logger{log: zap.NewNop()}

func Initialize() error {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		zap.InfoLevel,
	)

	logger := zap.New(core)
	log = Logger{log: logger}

	return nil
}

func Info(msg string, args ...zap.Field) {
	log.Info(msg, args...)
}

func Error(msg string, args ...zap.Field) {
	log.Error(msg, args...)
}

func Fatal(msg string, args ...zap.Field) {
	log.Fatal(msg, args...)
}

func Sync() {
	log.log.Sync()
}
