package middleware

import (
	"database/sql"
	"strings"

	"contrack-be/internal/database"
	"contrack-be/internal/models"
	"contrack-be/internal/services/jwt"
	"contrack-be/internal/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a middleware that validates JWT tokens
func AuthMiddleware(jwtService *jwt.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header is required")
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.UnauthorizedResponse(c, "Invalid authorization header format. Use: Bearer <token>")
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// Validate token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Store user information in context for use in handlers
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		
		// Get user role from database and store in context
		userRole, err := getUserRoleFromDB(claims.UserID)
		if err != nil {
			utils.UnauthorizedResponse(c, "Failed to retrieve user role")
			c.Abort()
			return
		}
		c.Set("user_role", userRole)
		
		c.Next()
	}
}

// GetUserID extracts the user ID from the Gin context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	
	userIDStr, ok := userID.(string)
	return userIDStr, ok
}

// GetUserEmail extracts the user email from the Gin context
func GetUserEmail(c *gin.Context) (string, bool) {
	userEmail, exists := c.Get("user_email")
	if !exists {
		return "", false
	}
	
	userEmailStr, ok := userEmail.(string)
	return userEmailStr, ok
}

// getUserRoleFromDB retrieves the user role from the database
func getUserRoleFromDB(userID string) (models.UserRole, error) {
	query := `SELECT role FROM users WHERE id = $1`
	
	var role string
	err := database.DB.QueryRow(query, userID).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", err
		}
		return "", err
	}
	
	return models.UserRole(role), nil
}
