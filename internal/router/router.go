package router

import (
	"contrack-be/internal/config"
	"contrack-be/internal/controllers"
	"contrack-be/internal/middleware"
	"contrack-be/internal/repository"
	"contrack-be/internal/services/ai"
	authService "contrack-be/internal/services/auth"
	jwtService "contrack-be/internal/services/jwt"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, jwtSvc *jwtService.Service, authSvc *authService.Service, cfg *config.Config) {
	// Create repositories
	clauseRepo := repository.NewClauseTemplateRepository()
	aiRepo := repository.NewAIRepository()
	
	// Create AI service
	aiService := ai.NewOpenAIService(cfg)
	
	// Create controllers
	authController := controllers.NewAuthController(authSvc)
	clauseController := controllers.NewClauseController()
	aiController := controllers.NewAIController(aiService, aiRepo, clauseRepo)
	
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
			
			// User management routes
			protected.GET("/users", controllers.ListUsers)
			
			// Clause template routes
			clauses := protected.Group("/clauses")
			{
				// CRUD operations
				clauses.POST("/", clauseController.CreateClauseTemplate)
				clauses.GET("/", clauseController.ListClauseTemplates)
				clauses.GET("/:id", clauseController.GetClauseTemplate)
				clauses.PUT("/:id", clauseController.UpdateClauseTemplate)
				clauses.DELETE("/:id", clauseController.DeleteClauseTemplate)
				
				// Additional endpoints
				clauses.GET("/by-code/:code", clauseController.GetClauseTemplateByCode)
				clauses.GET("/search", clauseController.SearchClauseTemplates)
				clauses.GET("/types", clauseController.GetClauseTypes)
				clauses.PATCH("/:id/toggle-status", clauseController.ToggleClauseTemplateStatus)
			}
			
			// AI Analysis routes
			ai := protected.Group("/ai")
			{
				// AI analysis endpoints
				ai.POST("/analyze", aiController.AnalyzeClauseRisk)
				ai.GET("/analysis/:id", aiController.GetAnalysisByID)
				ai.GET("/analysis/clause/:clause_id", aiController.GetAnalysisByClauseID)
				ai.GET("/analyses", aiController.GetAnalyses)
				ai.DELETE("/analysis/:id", aiController.DeleteAnalysis)
				ai.GET("/stats", aiController.GetAnalysisStats)
			}
		}
	}
}
