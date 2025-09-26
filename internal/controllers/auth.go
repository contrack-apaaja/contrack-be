package controllers

import (
	"contrack-be/internal/middleware"
	"contrack-be/internal/models"
	"contrack-be/internal/services/auth"
	"contrack-be/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *auth.Service
}

// NewAuthController creates a new authentication controller
func NewAuthController(authService *auth.Service) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Register handles user registration
func (ac *AuthController) Register(c *gin.Context) {
	var req models.UserRegistrationRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	user, token, err := ac.authService.Register(&req)
	if err != nil {
		if err.Error() == "user with email "+req.Email+" already exists" {
			utils.ConflictResponse(c, "User already exists")
			return
		}
		
		utils.InternalServerErrorResponse(c, "Failed to register user")
		return
	}

	utils.CreatedResponse(c, "User registered successfully", gin.H{
		"user":  user,
		"token": token,
	})
}

// Login handles user authentication
func (ac *AuthController) Login(c *gin.Context) {
	var req models.UserLoginRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	user, token, err := ac.authService.Login(&req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			utils.UnauthorizedResponse(c, "Invalid email or password")
			return
		}
		
		utils.InternalServerErrorResponse(c, "Failed to authenticate user")
		return
	}

	utils.OKResponse(c, "Login successful", gin.H{
		"user":  user,
		"token": token,
	})
}

// Profile returns the current user's profile
func (ac *AuthController) Profile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.UnauthorizedResponse(c, "User not found in context")
		return
	}

	user, err := ac.authService.GetUserByID(userID)
	if err != nil {
		utils.NotFoundResponse(c, "User not found")
		return
	}

	utils.OKResponse(c, "User profile retrieved successfully", gin.H{
		"user": user,
	})
}

// RefreshToken generates a new token for the user
func (ac *AuthController) RefreshToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Token is required")
		return
	}

	newToken, err := ac.authService.RefreshToken(req.Token)
	if err != nil {
		utils.UnauthorizedResponse(c, "Invalid or expired token")
		return
	}

	utils.OKResponse(c, "Token refreshed successfully", gin.H{
		"token": newToken,
	})
}
