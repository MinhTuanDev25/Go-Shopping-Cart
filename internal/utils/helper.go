package utils

import (
	"go-shopping-cart/pkg/logger"
	"os"

	"github.com/rs/zerolog"
)

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func NewLoggerWithPath(path string, level string) *zerolog.Logger {
	config := logger.LoggerConfig{
		Level:     level,
		Filename:  path,
		MaxSize:   1, // megabytes
		MaxBackup: 5,
		MaxAge:    5, //days
		Compress:  true,
		IsDev:     GetEnv("APP_ENV", "development"),
	}
	return logger.NewLogger(config)
}
