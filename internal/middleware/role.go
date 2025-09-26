package middleware

import (
	"contrack-be/internal/models"
	"contrack-be/internal/utils"

	"github.com/gin-gonic/gin"
)

// RequireRole creates a middleware that requires specific roles
func RequireRole(allowedRoles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by AuthMiddleware)
		userRole, exists := GetUserRole(c)
		if !exists {
			utils.UnauthorizedResponse(c, "User role not found in context")
			return
		}

		// Check if user role is in allowed roles
		allowed := false
		for _, role := range allowedRoles {
			if userRole == role {
				allowed = true
				break
			}
		}

		if !allowed {
			utils.ForbiddenResponse(c, "Insufficient permissions")
			return
		}

		c.Next()
	}
}

// RequireContractUpdatePermission creates a middleware that requires LEGAL or MANAGEMENT role
func RequireContractUpdatePermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context
		userRole, exists := GetUserRole(c)
		if !exists {
			utils.UnauthorizedResponse(c, "User role not found in context")
			return
		}

		// Check if user can update contracts
		if !userRole.CanUpdateContracts() {
			utils.ForbiddenResponse(c, "Only LEGAL and MANAGEMENT users can update contracts and clauses")
			return
		}

		c.Next()
	}
}

// GetUserRole retrieves the user role from the context
func GetUserRole(c *gin.Context) (models.UserRole, bool) {
	role, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	
	userRole, ok := role.(models.UserRole)
	return userRole, ok
}

// SetUserRole sets the user role in the context
func SetUserRole(c *gin.Context, role models.UserRole) {
	c.Set("user_role", role)
}
