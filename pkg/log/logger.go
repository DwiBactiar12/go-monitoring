package logger

import (
	"go.uber.org/zap"
)

var log *zap.Logger

// InitLogger initializes the zap logger
func InitLogger(isProduction bool) {
	var err error
	if isProduction {
		log, err = zap.NewProduction()
	} else {
		log, err = zap.NewDevelopment()
	}

	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
}

// GetLogger returns the logger instance
func GetLogger() *zap.Logger {
	if log == nil {
		InitLogger(false) // default to development mode
	}
	return log
}
