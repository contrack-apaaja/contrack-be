package config

import (
	"os"
)

type Config struct {
	Port        string
	SupabaseURL string
	SupabaseKey string
}

// Load reads config from environment variables with sensible defaults
func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port:        port,
		SupabaseURL: os.Getenv("SUPABASE_URL"),
		SupabaseKey: os.Getenv("SUPABASE_KEY"),
	}
}
