package routes

import (
	"go-shopping-cart/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {
	logPath := "../../internal/logs/http.log"
	recoveryPath := "../../internal/logs/recovery.log"

	httpLogger := zerolog.New(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    1, // megabytes
		MaxBackups: 5,
		MaxAge:     5, //days
		Compress:   true,
	}).With().Timestamp().Logger()

	recoveryLogger := zerolog.New(&lumberjack.Logger{
		Filename:   recoveryPath,
		MaxSize:    1, // megabytes
		MaxBackups: 5,
		MaxAge:     5, //days
		Compress:   true,
	}).With().Timestamp().Logger()

	r.Use(
		middleware.LoggerMiddleware(httpLogger),
		middleware.RecoveryMiddleware(recoveryLogger),
		middleware.ApiKeyMiddleware(),
		middleware.AuthMiddleware(),
		middleware.RateLimiterMiddleware(),
	)

	v1api := r.Group("/api/v1")

	for _, route := range routes {
		route.Register(v1api)
	}
}
