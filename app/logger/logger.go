package logger

import (
	"go.uber.org/zap"
)

type log struct {
	logger *zap.Logger
}

var l *log

func Init() error {
	l = new(log)

	var err error
	l.logger, err = zap.NewProduction()
	if err != nil {
		return err
	}
	return nil
}

func Info(message string, fields ...zap.Field) {
	l.logger.Info(message, fields...)
}

func Error(message string, fields ...zap.Field) {
	l.logger.Error(message, fields...)
}

func Warn(message string, fields ...zap.Field) {
	l.logger.Warn(message, fields...)
}

func Fatal(message string, fields ...zap.Field) {
	l.logger.Fatal(message, fields...)
}
