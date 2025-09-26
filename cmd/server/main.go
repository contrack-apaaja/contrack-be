package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"contrack-be/internal/config"
	"contrack-be/internal/database"
	"contrack-be/internal/middleware"
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
	
	// Add CORS middleware - using permissive middleware for development
	// For production, use: middleware.CORSMiddlewareWithOrigins([]string{"https://yourdomain.com"})
	r.Use(middleware.CorsMiddleware())

	// Setup routes
	router.Setup(r, jwtSvc, authSvc, cfg)

	fmt.Printf("🚀 Server starting on port %s\n", cfg.Port)
	fmt.Printf("📚 API Documentation:\n")
	fmt.Printf("   POST /api/auth/register - Register new user\n")
	fmt.Printf("   POST /api/auth/login    - Login user\n")
	fmt.Printf("   POST /api/auth/refresh  - Refresh token\n")
	fmt.Printf("   GET  /api/profile       - Get user profile (protected)\n")
	fmt.Printf("   GET  /api/hello         - Health check\n")
	fmt.Printf("\n🤖 AI Analysis Endpoints:\n")
	fmt.Printf("   POST /api/ai/analyze           - Analyze clause risk using AI\n")
	fmt.Printf("   GET  /api/ai/analysis/:id      - Get analysis by ID\n")
	fmt.Printf("   GET  /api/ai/analysis/clause/:id - Get analysis by clause ID\n")
	fmt.Printf("   GET  /api/ai/analyses          - Get analyses with pagination\n")
	fmt.Printf("   DELETE /api/ai/analysis/:id    - Delete analysis\n")
	fmt.Printf("   GET  /api/ai/stats             - Get analysis statistics\n")
	
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
