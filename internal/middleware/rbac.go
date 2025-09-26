package middleware

import (
	"net/http"

	"contrack-be/internal/models"
	"contrack-be/internal/utils"

	"github.com/gin-gonic/gin"
)

// RequireRole creates a middleware that requires specific roles
func RequireRole(allowedRoles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by AuthMiddleware)
		userRole, exists := c.Get("user_role")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "User role not found", "ROLE_NOT_FOUND", "Authentication required")
			c.Abort()
			return
		}

		role, ok := userRole.(models.UserRole)
		if !ok {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Invalid role format", "INVALID_ROLE", "Role data corrupted")
			c.Abort()
			return
		}

		// Check if user role is in allowed roles
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		// Access denied
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "INSUFFICIENT_PERMISSIONS", 
			"Your role '"+string(role)+"' is not authorized for this action")
		c.Abort()
	}
}

// RequireLegalRole middleware for legal-only endpoints
func RequireLegalRole() gin.HandlerFunc {
	return RequireRole(models.RoleLegal, models.RoleManagement)
}

// RequireManagementRole middleware for management-only endpoints
func RequireManagementRole() gin.HandlerFunc {
	return RequireRole(models.RoleManagement)
}

// RequireAnyRole middleware that allows any authenticated user
func RequireAnyRole() gin.HandlerFunc {
	return RequireRole(models.RoleUser, models.RoleLegal, models.RoleManagement)
}

// GetUserRole extracts the user role from the Gin context
func GetUserRole(c *gin.Context) (models.UserRole, bool) {
	userRole, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	
	role, ok := userRole.(models.UserRole)
	return role, ok
}