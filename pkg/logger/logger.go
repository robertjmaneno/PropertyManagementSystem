package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger(level string) *Logger {
	var zapLogger *zap.Logger
	var err error

	if level == "development" {
		zapLogger, err = zap.NewDevelopment()
	} else {
		zapLogger, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}

	sugar := zapLogger.Sugar()
	return &Logger{sugar}
} 