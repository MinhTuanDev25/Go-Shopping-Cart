package middleware

import (
	"context"
	"go-shopping-cart/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func TraceMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceID := ctx.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		contextValue := context.WithValue(ctx.Request.Context(), logger.TraceIDKey, traceID)
		ctx.Request = ctx.Request.WithContext(contextValue)

		// Set X-Trace-ID header response
		ctx.Writer.Header().Set("X-Trace-ID", traceID)

		ctx.Set(string(logger.TraceIDKey), traceID)

		ctx.Next()
	}
}
