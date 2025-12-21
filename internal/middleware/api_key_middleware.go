package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func ApiKeyMiddleware() gin.HandlerFunc {
	expectedKey := os.Getenv("API_KEY")
	if expectedKey == "" {
		expectedKey = "secret-key"
	}

	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("API-KEY")
		if apiKey == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Missing API-KEY"})
			return
		}

		if apiKey != expectedKey {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
			return
		}

		ctx.Next()
	}
}
