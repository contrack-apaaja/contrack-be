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

// CreateContractsTable creates the contracts table with versioning support
func CreateContractsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS contracts (
		id SERIAL PRIMARY KEY,
		base_id UUID NOT NULL DEFAULT gen_random_uuid(),
		version_number INTEGER NOT NULL DEFAULT 1,
		project_name VARCHAR(255) NOT NULL,
		package_name VARCHAR(255),
		contract_number VARCHAR(50) UNIQUE NOT NULL,
		external_reference VARCHAR(100),
		contract_type VARCHAR(100) NOT NULL,
		signing_place VARCHAR(255),
		signing_date DATE,
		total_value DECIMAL(15,2) NOT NULL CHECK (total_value > 0),
		funding_source VARCHAR(255),
		status VARCHAR(50) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'PENDING_LEGAL_REVIEW', 'PENDING_SIGNATURE', 'ACTIVE', 'EXPIRED', 'TERMINATED')),
		created_by UUID NOT NULL REFERENCES users(id),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		deleted_at TIMESTAMP WITH TIME ZONE,
		is_deleted BOOLEAN DEFAULT FALSE,
		UNIQUE(base_id, version_number)
	);

	-- Indexes for better performance
	CREATE INDEX IF NOT EXISTS idx_contracts_base_id ON contracts(base_id);
	CREATE INDEX IF NOT EXISTS idx_contracts_contract_number ON contracts(contract_number);
	CREATE INDEX IF NOT EXISTS idx_contracts_status ON contracts(status);
	CREATE INDEX IF NOT EXISTS idx_contracts_created_by ON contracts(created_by);
	CREATE INDEX IF NOT EXISTS idx_contracts_project_name ON contracts(project_name);
	CREATE INDEX IF NOT EXISTS idx_contracts_contract_type ON contracts(contract_type);
	CREATE INDEX IF NOT EXISTS idx_contracts_signing_date ON contracts(signing_date);
	CREATE INDEX IF NOT EXISTS idx_contracts_is_deleted ON contracts(is_deleted);
	
	-- Full text search index
	CREATE INDEX IF NOT EXISTS idx_contracts_search ON contracts USING gin(to_tsvector('english', project_name || ' ' || COALESCE(package_name, '') || ' ' || contract_type));
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create contracts table: %w", err)
	}

	log.Println("Contracts table created successfully")
	return nil
}

// CreateStakeholdersTable creates the stakeholders table
func CreateStakeholdersTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS stakeholders (
		id SERIAL PRIMARY KEY,
		legal_name VARCHAR(255) NOT NULL,
		address TEXT,
		type VARCHAR(100) NOT NULL CHECK (type IN ('INDIVIDUAL', 'COMPANY', 'GOVERNMENT', 'NGO', 'OTHER')),
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		deleted_at TIMESTAMP WITH TIME ZONE,
		is_deleted BOOLEAN DEFAULT FALSE
	);

	-- Indexes
	CREATE INDEX IF NOT EXISTS idx_stakeholders_legal_name ON stakeholders(legal_name);
	CREATE INDEX IF NOT EXISTS idx_stakeholders_type ON stakeholders(type);
	CREATE INDEX IF NOT EXISTS idx_stakeholders_is_deleted ON stakeholders(is_deleted);
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create stakeholders table: %w", err)
	}

	log.Println("Stakeholders table created successfully")
	return nil
}

// CreateContractStakeholdersTable creates the junction table for contract-stakeholder relationships
func CreateContractStakeholdersTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS contract_stakeholders (
		id SERIAL PRIMARY KEY,
		contract_id INTEGER NOT NULL REFERENCES contracts(id) ON DELETE CASCADE,
		stakeholder_id INTEGER NOT NULL REFERENCES stakeholders(id) ON DELETE CASCADE,
		role_in_contract VARCHAR(100) NOT NULL,
		representative_name VARCHAR(255),
		representative_title VARCHAR(255),
		other_details JSONB,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		UNIQUE(contract_id, stakeholder_id, role_in_contract)
	);

	-- Indexes
	CREATE INDEX IF NOT EXISTS idx_contract_stakeholders_contract_id ON contract_stakeholders(contract_id);
	CREATE INDEX IF NOT EXISTS idx_contract_stakeholders_stakeholder_id ON contract_stakeholders(stakeholder_id);
	CREATE INDEX IF NOT EXISTS idx_contract_stakeholders_role ON contract_stakeholders(role_in_contract);
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create contract_stakeholders table: %w", err)
	}

	log.Println("Contract stakeholders table created successfully")
	return nil
}

// CreateContractClausesTable creates the table for contract clauses
func CreateContractClausesTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS contract_clauses (
		id SERIAL PRIMARY KEY,
		contract_id INTEGER NOT NULL REFERENCES contracts(id) ON DELETE CASCADE,
		clause_template_id INTEGER NOT NULL REFERENCES clause_templates(id),
		display_order INTEGER NOT NULL,
		custom_content TEXT,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		UNIQUE(contract_id, clause_template_id),
		UNIQUE(contract_id, display_order)
	);

	-- Indexes
	CREATE INDEX IF NOT EXISTS idx_contract_clauses_contract_id ON contract_clauses(contract_id);
	CREATE INDEX IF NOT EXISTS idx_contract_clauses_template_id ON contract_clauses(clause_template_id);
	CREATE INDEX IF NOT EXISTS idx_contract_clauses_display_order ON contract_clauses(contract_id, display_order);
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create contract_clauses table: %w", err)
	}

	log.Println("Contract clauses table created successfully")
	return nil
}

// CreateContractStatusHistoryTable creates the table to track status changes
func CreateContractStatusHistoryTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS contract_status_history (
		id SERIAL PRIMARY KEY,
		contract_id INTEGER NOT NULL REFERENCES contracts(id) ON DELETE CASCADE,
		from_status VARCHAR(50),
		to_status VARCHAR(50) NOT NULL,
		changed_by UUID NOT NULL REFERENCES users(id),
		change_reason TEXT,
		comments TEXT,
		changed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);

	-- Indexes
	CREATE INDEX IF NOT EXISTS idx_contract_status_history_contract_id ON contract_status_history(contract_id);
	CREATE INDEX IF NOT EXISTS idx_contract_status_history_changed_by ON contract_status_history(changed_by);
	CREATE INDEX IF NOT EXISTS idx_contract_status_history_changed_at ON contract_status_history(changed_at);
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create contract_status_history table: %w", err)
	}

	log.Println("Contract status history table created successfully")
	return nil
}

// CreateContractSequenceTable creates the sequence table for contract numbering
func CreateContractSequenceTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS contract_sequences (
		year_month VARCHAR(7) PRIMARY KEY, -- Format: YYYY-MM
		sequence_number INTEGER NOT NULL DEFAULT 0
	);
	`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create contract_sequences table: %w", err)
	}

	log.Println("Contract sequences table created successfully")
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

	if err := CreateContractsTable(); err != nil {
		return err
	}

	if err := CreateStakeholdersTable(); err != nil {
		return err
	}

	if err := CreateContractStakeholdersTable(); err != nil {
		return err
	}

	if err := CreateContractClausesTable(); err != nil {
		return err
	}

	if err := CreateContractStatusHistoryTable(); err != nil {
		return err
	}

	if err := CreateContractSequenceTable(); err != nil {
		return err
	}
	
	log.Println("Database migrations completed successfully")
	return nil
}
