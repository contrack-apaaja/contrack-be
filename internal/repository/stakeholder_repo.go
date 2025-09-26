package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"contrack-be/internal/database"
	"contrack-be/internal/models"
)

type StakeholderRepository struct {
	db *sql.DB
}

func NewStakeholderRepository() *StakeholderRepository {
	return &StakeholderRepository{
		db: database.DB,
	}
}

// Create creates a new stakeholder
func (r *StakeholderRepository) Create(stakeholder *models.Stakeholder) error {
	query := `
		INSERT INTO stakeholders (legal_name, address, type)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		stakeholder.LegalName,
		stakeholder.Address,
		stakeholder.Type,
	).Scan(&stakeholder.ID, &stakeholder.CreatedAt, &stakeholder.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create stakeholder: %w", err)
	}

	return nil
}

// GetByID retrieves a stakeholder by ID
func (r *StakeholderRepository) GetByID(id int) (*models.Stakeholder, error) {
	stakeholder := &models.Stakeholder{}
	
	query := `
		SELECT id, legal_name, address, type, created_at, updated_at, deleted_at, is_deleted
		FROM stakeholders 
		WHERE id = $1 AND is_deleted = FALSE`

	err := r.db.QueryRow(query, id).Scan(
		&stakeholder.ID,
		&stakeholder.LegalName,
		&stakeholder.Address,
		&stakeholder.Type,
		&stakeholder.CreatedAt,
		&stakeholder.UpdatedAt,
		&stakeholder.DeletedAt,
		&stakeholder.IsDeleted,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("stakeholder not found")
		}
		return nil, fmt.Errorf("failed to get stakeholder: %w", err)
	}

	return stakeholder, nil
}

// Update updates a stakeholder
func (r *StakeholderRepository) Update(id int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	// Build dynamic query
	setParts := []string{}
	args := []interface{}{}
	argCount := 1

	for column, value := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", column, argCount))
		args = append(args, value)
		argCount++
	}

	// Always update the updated_at timestamp
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argCount))
	args = append(args, time.Now())
	argCount++

	// Add WHERE clause
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE stakeholders 
		SET %s 
		WHERE id = $%d AND is_deleted = FALSE`,
		strings.Join(setParts, ", "),
		argCount,
	)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update stakeholder: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("stakeholder not found or already deleted")
	}

	return nil
}

// Delete soft deletes a stakeholder
func (r *StakeholderRepository) Delete(id int) error {
	query := `
		UPDATE stakeholders 
		SET is_deleted = TRUE, deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND is_deleted = FALSE`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete stakeholder: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("stakeholder not found or already deleted")
	}

	return nil
}

// List retrieves stakeholders with pagination and filtering
func (r *StakeholderRepository) List(search string, stakeholderType *models.StakeholderType, page, limit int) ([]models.Stakeholder, int, error) {
	// Set defaults
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	// Build WHERE conditions
	whereConditions := []string{"is_deleted = FALSE"}
	args := []interface{}{}
	argCount := 1

	// Search filter
	if search != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("legal_name ILIKE $%d", argCount))
		args = append(args, "%"+search+"%")
		argCount++
	}

	// Type filter
	if stakeholderType != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("type = $%d", argCount))
		args = append(args, *stakeholderType)
		argCount++
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM stakeholders WHERE %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count stakeholders: %w", err)
	}

	// Fetch records with pagination
	offset := (page - 1) * limit
	args = append(args, limit, offset)

	query := fmt.Sprintf(`
		SELECT id, legal_name, address, type, created_at, updated_at, deleted_at, is_deleted
		FROM stakeholders 
		WHERE %s 
		ORDER BY legal_name ASC
		LIMIT $%d OFFSET $%d`,
		whereClause,
		argCount, argCount+1,
	)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query stakeholders: %w", err)
	}
	defer rows.Close()

	stakeholders := []models.Stakeholder{}
	for rows.Next() {
		stakeholder := models.Stakeholder{}
		err := rows.Scan(
			&stakeholder.ID,
			&stakeholder.LegalName,
			&stakeholder.Address,
			&stakeholder.Type,
			&stakeholder.CreatedAt,
			&stakeholder.UpdatedAt,
			&stakeholder.DeletedAt,
			&stakeholder.IsDeleted,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan stakeholder: %w", err)
		}
		stakeholders = append(stakeholders, stakeholder)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate stakeholders: %w", err)
	}

	return stakeholders, total, nil
}

// AddToContract adds a stakeholder to a contract with specific role
func (r *StakeholderRepository) AddToContract(contractStakeholder *models.ContractStakeholder) error {
	query := `
		INSERT INTO contract_stakeholders (
			contract_id, stakeholder_id, role_in_contract, 
			representative_name, representative_title, other_details
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		contractStakeholder.ContractID,
		contractStakeholder.StakeholderID,
		contractStakeholder.RoleInContract,
		contractStakeholder.RepresentativeName,
		contractStakeholder.RepresentativeTitle,
		contractStakeholder.OtherDetails,
	).Scan(&contractStakeholder.ID, &contractStakeholder.CreatedAt, &contractStakeholder.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to add stakeholder to contract: %w", err)
	}

	return nil
}

// RemoveFromContract removes a stakeholder from a contract
func (r *StakeholderRepository) RemoveFromContract(contractID, stakeholderID int, role string) error {
	query := `
		DELETE FROM contract_stakeholders 
		WHERE contract_id = $1 AND stakeholder_id = $2 AND role_in_contract = $3`

	result, err := r.db.Exec(query, contractID, stakeholderID, role)
	if err != nil {
		return fmt.Errorf("failed to remove stakeholder from contract: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("stakeholder relationship not found")
	}

	return nil
}

// GetContractStakeholders retrieves all stakeholders for a contract
func (r *StakeholderRepository) GetContractStakeholders(contractID int) ([]models.ContractStakeholder, error) {
	query := `
		SELECT cs.id, cs.contract_id, cs.stakeholder_id, cs.role_in_contract,
			   cs.representative_name, cs.representative_title, cs.other_details,
			   cs.created_at, cs.updated_at,
			   s.id, s.legal_name, s.address, s.type, s.created_at, s.updated_at,
			   s.deleted_at, s.is_deleted
		FROM contract_stakeholders cs
		JOIN stakeholders s ON cs.stakeholder_id = s.id
		WHERE cs.contract_id = $1 AND s.is_deleted = FALSE
		ORDER BY cs.role_in_contract, s.legal_name`

	rows, err := r.db.Query(query, contractID)
	if err != nil {
		return nil, fmt.Errorf("failed to query contract stakeholders: %w", err)
	}
	defer rows.Close()

	contractStakeholders := []models.ContractStakeholder{}
	for rows.Next() {
		cs := models.ContractStakeholder{}
		stakeholder := models.Stakeholder{}
		
		err := rows.Scan(
			&cs.ID,
			&cs.ContractID,
			&cs.StakeholderID,
			&cs.RoleInContract,
			&cs.RepresentativeName,
			&cs.RepresentativeTitle,
			&cs.OtherDetails,
			&cs.CreatedAt,
			&cs.UpdatedAt,
			&stakeholder.ID,
			&stakeholder.LegalName,
			&stakeholder.Address,
			&stakeholder.Type,
			&stakeholder.CreatedAt,
			&stakeholder.UpdatedAt,
			&stakeholder.DeletedAt,
			&stakeholder.IsDeleted,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contract stakeholder: %w", err)
		}
		
		cs.Stakeholder = &stakeholder
		contractStakeholders = append(contractStakeholders, cs)
	}

	return contractStakeholders, rows.Err()
}

// UpdateContractStakeholder updates a contract stakeholder relationship
func (r *StakeholderRepository) UpdateContractStakeholder(id int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	// Build dynamic query
	setParts := []string{}
	args := []interface{}{}
	argCount := 1

	for column, value := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", column, argCount))
		args = append(args, value)
		argCount++
	}

	// Always update the updated_at timestamp
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argCount))
	args = append(args, time.Now())
	argCount++

	// Add WHERE clause
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE contract_stakeholders 
		SET %s 
		WHERE id = $%d`,
		strings.Join(setParts, ", "),
		argCount,
	)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update contract stakeholder: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("contract stakeholder not found")
	}

	return nil
}