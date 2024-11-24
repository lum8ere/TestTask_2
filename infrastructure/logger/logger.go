package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

// NewLogger создает новый инстанс логгера
func NewLogger() *Logger {
	rawLogger, _ := zap.NewProduction()
	defer rawLogger.Sync()

	return &Logger{rawLogger.Sugar()}
}
