package router

import (
	"contrack-be/internal/controllers"
	"contrack-be/internal/middleware"
	authService "contrack-be/internal/services/auth"
	jwtService "contrack-be/internal/services/jwt"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, jwtSvc *jwtService.Service, authSvc *authService.Service) {
	// Create auth controller
	authController := controllers.NewAuthController(authSvc)
	
	api := r.Group("/api")
	{
		// Public routes (no authentication required)
		api.GET("/hello", controllers.Hello)
		
		// Authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/refresh", authController.RefreshToken)
		}
		
		// Protected routes (authentication required)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(jwtSvc))
		{
			// User profile routes
			protected.GET("/profile", authController.Profile)
			
			// Other protected routes
			protected.GET("/users", controllers.ListUsers)
		}
	}
}
