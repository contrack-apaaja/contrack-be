package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CorsMiddleware returns a Gin middleware function for CORS
func CorsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Set CORS headers
        c.Header("Access-Control-Allow-Origin", "*")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Handle preflight request (OPTIONS)
        if c.Request.Method == http.MethodOptions {
            c.AbortWithStatus(http.StatusOK)
            return
        }

        c.Next()
    }
}

// CORSMiddlewareWithOrigins creates a CORS middleware with specific allowed origins

