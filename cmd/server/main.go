package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"contrack-be/internal/config"
	"contrack-be/internal/database"
	"contrack-be/internal/router"
	authService "contrack-be/internal/services/auth"
	jwtService "contrack-be/internal/services/jwt"
	sup "contrack-be/internal/services/supabase"
)

func main() {
	// Load environment variables from .env file (optional)
	_ = godotenv.Load()

	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	if cfg.DatabaseURL != "" {
		if err := database.Connect(cfg.DatabaseURL); err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer database.Close()

		// Run database migrations
		if err := database.Migrate(); err != nil {
			log.Fatalf("Failed to run database migrations: %v", err)
		}
	} else {
		log.Println("Warning: DATABASE_URL not set, database features will not be available")
	}

	// Initialize supabase (internal service reads env too; keep for clarity)
	if err := sup.Init(); err != nil {
		log.Printf("Warning: Failed to initialize Supabase: %v", err)
	}

	// Initialize JWT service
	jwtSvc := jwtService.NewJWTService(cfg.JWTSecret, cfg.JWTExpiresIn)

	// Initialize auth service
	authSvc := authService.NewAuthService(jwtSvc)

	// Setup Gin router
	r := gin.Default()
	
	// Add CORS middleware for development
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Setup routes
	router.Setup(r, jwtSvc, authSvc)

	fmt.Printf("🚀 Server starting on port %s\n", cfg.Port)
	fmt.Printf("📚 API Documentation:\n")
	fmt.Printf("   POST /api/auth/register - Register new user\n")
	fmt.Printf("   POST /api/auth/login    - Login user\n")
	fmt.Printf("   POST /api/auth/refresh  - Refresh token\n")
	fmt.Printf("   GET  /api/profile       - Get user profile (protected)\n")
	fmt.Printf("   GET  /api/hello         - Health check\n")
	
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
