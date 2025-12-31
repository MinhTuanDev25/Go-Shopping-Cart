package routes

import (
	"go-shopping-cart/internal/middleware"
	"go-shopping-cart/internal/utils"

	"github.com/gin-gonic/gin"
)

type Route interface {
	Register(r *gin.RouterGroup)
}

func RegisterRoutes(r *gin.Engine, routes ...Route) {
	logPath := "../../internal/logs/http.log"
	recoveryPath := "../../internal/logs/recovery.log"
	rateLimiterPath := "../../internal/logs/rate_limiter.log"

	httpLogger := utils.NewLoggerWithPath(logPath, "info")

	recoveryLogger := utils.NewLoggerWithPath(recoveryPath, "warning")

	rateLimiterLogger := utils.NewLoggerWithPath(rateLimiterPath, "warning")
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
