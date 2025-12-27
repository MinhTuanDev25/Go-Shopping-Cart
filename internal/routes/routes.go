package routes

import (
	"go-shopping-cart/internal/middleware"
	"go-shopping-cart/internal/utils"
	"go-shopping-cart/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {
	logPath := "../../internal/logs/http.log"
	recoveryPath := "../../internal/logs/recovery.log"
	rateLimiterPath := "../../internal/logs/rate_limiter.log"

	httpLogger := newLoggerWithPath(logPath, "info")

	recoveryLogger := newLoggerWithPath(recoveryPath, "warning")

	rateLimiterLogger := newLoggerWithPath(rateLimiterPath, "warning")

	r.Use(
		middleware.RateLimiterMiddleware(rateLimiterLogger),
		middleware.LoggerMiddleware(httpLogger),
		middleware.RecoveryMiddleware(recoveryLogger),
		middleware.ApiKeyMiddleware(),
		middleware.AuthMiddleware(),
	)

	v1api := r.Group("/api/v1")

	for _, route := range routes {
		route.Register(v1api)
	}
}

func newLoggerWithPath(path string, level string) *zerolog.Logger {
	config := logger.LoggerConfig{
		Level:     level,
		Filename:  path,
		MaxSize:   1, // megabytes
		MaxBackup: 5,
		MaxAge:    5, //days
		Compress:  true,
		IsDev:     utils.GetEnv("APP_ENV", "development"),
	}
	return logger.NewLogger(config)
}
