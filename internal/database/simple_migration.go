package database

import (
	"fmt"
	"log"
)

// SimpleRoleMigration performs a simple role column migration
func SimpleRoleMigration() error {
	// Step 1: Try to add the role column (ignore error if it already exists)
	addColumnQuery := `ALTER TABLE users ADD COLUMN role VARCHAR(20) DEFAULT 'REGULAR'`
	_, err := DB.Exec(addColumnQuery)
	if err != nil {
		// Check if error is because column already exists
		if !isColumnExistsError(err) {
			return fmt.Errorf("failed to add role column: %w", err)
		}
		log.Println("Role column already exists, skipping...")
	} else {
		log.Println("Role column added successfully")
	}

	// Step 2: Update existing users
	updateQuery := `UPDATE users SET role = 'REGULAR' WHERE role IS NULL OR role = ''`
	_, err = DB.Exec(updateQuery)
	if err != nil {
		return fmt.Errorf("failed to update existing users: %w", err)
	}

	// Step 3: Try to add constraint (ignore error if it already exists)
	constraintQuery := `
		ALTER TABLE users ADD CONSTRAINT users_role_check 
		CHECK (role IN ('REGULAR', 'LEGAL', 'MANAGEMENT'))
	`
	_, err = DB.Exec(constraintQuery)
	if err != nil {
		if !isConstraintExistsError(err) {
			return fmt.Errorf("failed to add role constraint: %w", err)
		}
		log.Println("Role constraint already exists, skipping...")
	} else {
		log.Println("Role constraint added successfully")
	}

	// Step 4: Create index
	indexQuery := `CREATE INDEX IF NOT EXISTS idx_users_role ON users(role)`
	_, err = DB.Exec(indexQuery)
	if err != nil {
		return fmt.Errorf("failed to create role index: %w", err)
	}

	log.Println("Simple role migration completed successfully")
	return nil
}

// isColumnExistsError checks if the error is due to column already existing
func isColumnExistsError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "already exists") || 
		   contains(errStr, "duplicate column") ||
		   contains(errStr, "column \"role\" of relation \"users\" already exists")
}

// isConstraintExistsError checks if the error is due to constraint already existing
func isConstraintExistsError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "already exists") || 
		   contains(errStr, "duplicate key") ||
		   contains(errStr, "constraint \"users_role_check\" already exists")
}

// contains checks if a string contains a substring (case insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    (len(s) > len(substr) && 
		     (s[:len(substr)] == substr || 
		      s[len(s)-len(substr):] == substr || 
		      containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
