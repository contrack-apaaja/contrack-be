package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"contrack-be/internal/database"
	"contrack-be/internal/models"
)

type ContractClauseRepository struct {
	db *sql.DB
}

func NewContractClauseRepository() *ContractClauseRepository {
	return &ContractClauseRepository{
		db: database.DB,
	}
}

// AddClauseToContract adds a clause to a contract
func (r *ContractClauseRepository) AddClauseToContract(contractClause *models.ContractClause) error {
	query := `
		INSERT INTO contract_clauses (
			contract_id, clause_template_id, display_order, custom_content
		) VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		contractClause.ContractID,
		contractClause.ClauseTemplateID,
		contractClause.DisplayOrder,
		contractClause.CustomContent,
	).Scan(&contractClause.ID, &contractClause.CreatedAt, &contractClause.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to add clause to contract: %w", err)
	}

	return nil
}

// RemoveClauseFromContract removes a clause from a contract
func (r *ContractClauseRepository) RemoveClauseFromContract(contractID, clauseTemplateID int) error {
	query := `
		DELETE FROM contract_clauses 
		WHERE contract_id = $1 AND clause_template_id = $2`

	result, err := r.db.Exec(query, contractID, clauseTemplateID)
	if err != nil {
		return fmt.Errorf("failed to remove clause from contract: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("contract clause not found")
	}

	return nil
}

// GetContractClauses retrieves all clauses for a contract
func (r *ContractClauseRepository) GetContractClauses(contractID int) ([]models.ContractClause, error) {
	query := `
		SELECT cc.id, cc.contract_id, cc.clause_template_id, cc.display_order,
			   cc.custom_content, cc.created_at, cc.updated_at,
			   ct.id, ct.clause_code, ct.title, ct.type, ct.content, ct.is_active,
			   ct.created_at, ct.updated_at
		FROM contract_clauses cc
		JOIN clause_templates ct ON cc.clause_template_id = ct.id
		WHERE cc.contract_id = $1 AND ct.is_active = TRUE
		ORDER BY cc.display_order ASC`

	rows, err := r.db.Query(query, contractID)
	if err != nil {
		return nil, fmt.Errorf("failed to query contract clauses: %w", err)
	}
	defer rows.Close()

	contractClauses := []models.ContractClause{}
	for rows.Next() {
		cc := models.ContractClause{}
		clauseTemplate := models.ClauseTemplate{}
		
		err := rows.Scan(
			&cc.ID,
			&cc.ContractID,
			&cc.ClauseTemplateID,
			&cc.DisplayOrder,
			&cc.CustomContent,
			&cc.CreatedAt,
			&cc.UpdatedAt,
			&clauseTemplate.ID,
			&clauseTemplate.ClauseCode,
			&clauseTemplate.Title,
			&clauseTemplate.Type,
			&clauseTemplate.Content,
			&clauseTemplate.IsActive,
			&clauseTemplate.CreatedAt,
			&clauseTemplate.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contract clause: %w", err)
		}
		
		cc.ClauseTemplate = &clauseTemplate
		contractClauses = append(contractClauses, cc)
	}

	return contractClauses, rows.Err()
}

// UpdateContractClause updates a contract clause
func (r *ContractClauseRepository) UpdateContractClause(id int, displayOrder *int, customContent *string) error {
	updates := make(map[string]interface{})
	
	if displayOrder != nil {
		updates["display_order"] = *displayOrder
	}
	if customContent != nil {
		updates["custom_content"] = *customContent
	}

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
		UPDATE contract_clauses 
		SET %s 
		WHERE id = $%d`,
		strings.Join(setParts, ", "),
		argCount,
	)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update contract clause: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("contract clause not found")
	}

	return nil
}

// ReorderClauses updates the display order of all clauses for a contract
func (r *ContractClauseRepository) ReorderClauses(contractID int, clauseOrders []struct {
	ClauseID     int `json:"clause_id"`
	DisplayOrder int `json:"display_order"`
}) error {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Update each clause order
	for _, order := range clauseOrders {
		_, err = tx.Exec(`
			UPDATE contract_clauses 
			SET display_order = $1, updated_at = NOW()
			WHERE id = $2 AND contract_id = $3`,
			order.DisplayOrder, order.ClauseID, contractID)
		if err != nil {
			return fmt.Errorf("failed to update clause order: %w", err)
		}
	}

	return tx.Commit()
}

// GetNextDisplayOrder gets the next available display order for a contract
func (r *ContractClauseRepository) GetNextDisplayOrder(contractID int) (int, error) {
	var nextOrder int
	err := r.db.QueryRow(`
		SELECT COALESCE(MAX(display_order), 0) + 1 
		FROM contract_clauses 
		WHERE contract_id = $1`,
		contractID).Scan(&nextOrder)
	
	if err != nil {
		return 0, fmt.Errorf("failed to get next display order: %w", err)
	}
	
	return nextOrder, nil
}

// BulkAddClausesToContract adds multiple clauses to a contract in one transaction
func (r *ContractClauseRepository) BulkAddClausesToContract(contractID int, clauseTemplateIDs []int) error {
	if len(clauseTemplateIDs) == 0 {
		return nil
	}

	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get starting display order
	var startOrder int
	err = tx.QueryRow(`
		SELECT COALESCE(MAX(display_order), 0) + 1 
		FROM contract_clauses 
		WHERE contract_id = $1`,
		contractID).Scan(&startOrder)
	if err != nil {
		return fmt.Errorf("failed to get starting display order: %w", err)
	}

	// Add each clause
	for i, templateID := range clauseTemplateIDs {
		_, err = tx.Exec(`
			INSERT INTO contract_clauses (contract_id, clause_template_id, display_order)
			VALUES ($1, $2, $3)`,
			contractID, templateID, startOrder+i)
		if err != nil {
			return fmt.Errorf("failed to add clause template %d to contract: %w", templateID, err)
		}
	}

	return tx.Commit()
}

// GetContractClauseByID retrieves a specific contract clause by ID
func (r *ContractClauseRepository) GetContractClauseByID(id int) (*models.ContractClause, error) {
	cc := &models.ContractClause{}
	clauseTemplate := &models.ClauseTemplate{}
	
	query := `
		SELECT cc.id, cc.contract_id, cc.clause_template_id, cc.display_order,
			   cc.custom_content, cc.created_at, cc.updated_at,
			   ct.id, ct.clause_code, ct.title, ct.type, ct.content, ct.is_active,
			   ct.created_at, ct.updated_at
		FROM contract_clauses cc
		JOIN clause_templates ct ON cc.clause_template_id = ct.id
		WHERE cc.id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&cc.ID,
		&cc.ContractID,
		&cc.ClauseTemplateID,
		&cc.DisplayOrder,
		&cc.CustomContent,
		&cc.CreatedAt,
		&cc.UpdatedAt,
		&clauseTemplate.ID,
		&clauseTemplate.ClauseCode,
		&clauseTemplate.Title,
		&clauseTemplate.Type,
		&clauseTemplate.Content,
		&clauseTemplate.IsActive,
		&clauseTemplate.CreatedAt,
		&clauseTemplate.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("contract clause not found")
		}
		return nil, fmt.Errorf("failed to get contract clause: %w", err)
	}

	cc.ClauseTemplate = clauseTemplate
	return cc, nil
}