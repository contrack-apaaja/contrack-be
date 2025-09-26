package config

import (
	"os"
	"time"
)

type Config struct {
	Port         string
	SupabaseURL  string
	SupabaseKey  string
	DatabaseURL  string
	JWTSecret    string
	JWTExpiresIn time.Duration
	OpenAIAPIKey string
}

// Load reads config from environment variables with sensible defaults
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key-change-this-in-production"
	}

	// Default JWT expiration to 24 hours
	jwtExpiresIn := 24 * time.Hour

	return &Config{
		Port:         port,
		SupabaseURL:  os.Getenv("SUPABASE_URL"),
		SupabaseKey:  os.Getenv("SUPABASE_KEY"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		JWTSecret:    jwtSecret,
		JWTExpiresIn: jwtExpiresIn,
		OpenAIAPIKey: os.Getenv("OPENAI_API_KEY"),
	}
}
