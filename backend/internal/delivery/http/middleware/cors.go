package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	allowedOrigins := map[string]bool{
		"http://localhost":      true,
		"http://localhost:3000": true,
		"http://127.0.0.1":      true,
		"http://127.0.0.1:3000": true,
	}
	for _, origin := range strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",") {
		if trimmed := strings.TrimSpace(origin); trimmed != "" {
			allowedOrigins[trimmed] = true
		}
	}
	if frontendURL := strings.TrimSpace(os.Getenv("FRONTEND_URL")); frontendURL != "" {
		allowedOrigins[frontendURL] = true
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if allowedOrigins[origin] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
