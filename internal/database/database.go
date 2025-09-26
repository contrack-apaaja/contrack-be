package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Connect establishes a connection to the PostgreSQL database
func Connect(databaseURL string) error {
	var err error
	
	if databaseURL == "" {
		return fmt.Errorf("database URL is required")
	}

	DB, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// CreateUsersTable creates the users table if it doesn't exist
func CreateUsersTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	log.Println("Users table created successfully")
	return nil
}

// CreateClauseTemplatesTable creates the clause_templates table if it doesn't exist
func CreateClauseTemplatesTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS clause_templates (
		id SERIAL PRIMARY KEY,
		clause_code VARCHAR(50) UNIQUE NOT NULL,
		title VARCHAR(255) NOT NULL,
		type VARCHAR(100) NOT NULL,
		content TEXT NOT NULL,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	-- Indexes for better search performance
	CREATE INDEX IF NOT EXISTS idx_clause_templates_clause_code ON clause_templates(clause_code);
	CREATE INDEX IF NOT EXISTS idx_clause_templates_title ON clause_templates(title);
	CREATE INDEX IF NOT EXISTS idx_clause_templates_type ON clause_templates(type);
	CREATE INDEX IF NOT EXISTS idx_clause_templates_is_active ON clause_templates(is_active);
	CREATE INDEX IF NOT EXISTS idx_clause_templates_created_at ON clause_templates(created_at);
	
	-- Full text search index for title and content
	CREATE INDEX IF NOT EXISTS idx_clause_templates_search ON clause_templates USING gin(to_tsvector('english', title || ' ' || content));
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create clause_templates table: %w", err)
	}

	log.Println("Clause templates table created successfully")
	return nil
}

// Migrate runs all necessary database migrations
func Migrate() error {
	if err := CreateUsersTable(); err != nil {
		return err
	}
	
	if err := CreateClauseTemplatesTable(); err != nil {
		return err
	}
	
	log.Println("Database migrations completed successfully")
	return nil
}
