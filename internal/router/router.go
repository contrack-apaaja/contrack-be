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

	// Create repositories
	contractRepo := repository.NewContractRepository()

	// Create controllers
	authController := controllers.NewAuthController(authSvc)
	clauseController := controllers.NewClauseController()
	aiController := controllers.NewAIController(aiService, aiRepo, clauseRepo)
	contractController := controllers.NewContractController()
	stakeholderController := controllers.NewStakeholderController()
	dashboardController := controllers.NewDashboardController(contractRepo)

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
		protected.Use(middleware.CorsMiddleware())
		protected.Use(middleware.AuthMiddleware(jwtSvc))
		{
			// User profile routes
			protected.GET("/profile", authController.Profile)

			// User management routes
			protected.GET("/users", controllers.ListUsers)

			// Clause template routes
			clauses := protected.Group("/clauses")
			{
				clauses.POST("/", clauseController.CreateClauseTemplate)
				clauses.GET("/", clauseController.ListClauseTemplates)
				clauses.GET("/:id", clauseController.GetClauseTemplate)
				clauses.PUT("/:id", clauseController.UpdateClauseTemplate)
				clauses.DELETE("/:id", clauseController.DeleteClauseTemplate)

				clauses.GET("/by-code/:code", clauseController.GetClauseTemplateByCode)
				clauses.GET("/search", clauseController.SearchClauseTemplates)
				clauses.GET("/types", clauseController.GetClauseTypes)
				clauses.PATCH("/:id/toggle-status", clauseController.ToggleClauseTemplateStatus)
			}

			// AI Analysis routes
			ai := protected.Group("/ai")
			{
				ai.POST("/analyze", aiController.AnalyzeClauseRisk)
				ai.POST("/analyze-contract", aiController.AnalyzeContractRisk)
				ai.GET("/analysis/:id", aiController.GetAnalysisByID)
				ai.GET("/analysis/clause/:clause_id", aiController.GetAnalysisByClauseID)
				ai.GET("/analyses", aiController.GetAnalyses)
				ai.DELETE("/analysis/:id", aiController.DeleteAnalysis)
				ai.GET("/stats", aiController.GetAnalysisStats)
			} // <<– tutup group AI di sini

			// Contract routes
			contracts := protected.Group("/contracts")
			{
				contracts.GET("/stats", contractController.GetContractStats)

				contracts.POST("/", contractController.CreateContract)
				contracts.GET("/", contractController.ListContracts)
				contracts.GET("/:id", contractController.GetContract)
				contracts.PUT("/:id", contractController.UpdateContract)
				contracts.DELETE("/:id", contractController.DeleteContract)

				contracts.POST("/:id/status", contractController.ChangeContractStatus)
				contracts.GET("/:id/status-history", contractController.GetContractStatusHistory)
			}

			// Contract versioning routes
			contractVersions := protected.Group("/contract-versions")
			{
				contractVersions.POST("/:baseId", contractController.CreateContractVersion)
				contractVersions.GET("/:baseId", contractController.GetContractVersions)
			}

			// Stakeholder routes
			stakeholders := protected.Group("/stakeholders")
			{
				stakeholders.POST("/", stakeholderController.CreateStakeholder)
				stakeholders.GET("/", stakeholderController.ListStakeholders)
				stakeholders.GET("/:id", stakeholderController.GetStakeholder)
				stakeholders.PUT("/:id", stakeholderController.UpdateStakeholder)
				stakeholders.DELETE("/:id", stakeholderController.DeleteStakeholder)

				stakeholders.GET("/types", stakeholderController.GetStakeholderTypes)
			}

			// Dashboard routes
			dashboard := protected.Group("/dashboard")
			{
				dashboard.GET("/status-counts", dashboardController.GetStatusCounts)
				dashboard.GET("/contracts", dashboardController.GetContractList)
			}
		}
	}
}
