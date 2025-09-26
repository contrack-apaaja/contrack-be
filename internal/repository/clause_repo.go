package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"contrack-be/internal/database"
	"contrack-be/internal/models"
)

type ClauseTemplateRepository struct {
	db *sql.DB
}

// NewClauseTemplateRepository creates a new clause template repository
func NewClauseTemplateRepository() *ClauseTemplateRepository {
	return &ClauseTemplateRepository{
		db: database.DB,
	}
}

// Create inserts a new clause template
func (r *ClauseTemplateRepository) Create(clause *models.ClauseTemplate) (*models.ClauseTemplate, error) {
	query := `
		INSERT INTO clause_templates (clause_code, title, type, content, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		query,
		clause.ClauseCode,
		clause.Title,
		clause.Type,
		clause.Content,
		clause.IsActive,
	).Scan(&clause.ID, &clause.CreatedAt, &clause.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create clause template: %w", err)
	}

	return clause, nil
}

// GetByID retrieves a clause template by ID
func (r *ClauseTemplateRepository) GetByID(id int) (*models.ClauseTemplate, error) {
	query := `
		SELECT id, clause_code, title, type, content, is_active, created_at, updated_at
		FROM clause_templates
		WHERE id = $1
	`

	var clause models.ClauseTemplate
	err := r.db.QueryRow(query, id).Scan(
		&clause.ID,
		&clause.ClauseCode,
		&clause.Title,
		&clause.Type,
		&clause.Content,
		&clause.IsActive,
		&clause.CreatedAt,
		&clause.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("clause template not found")
		}
		return nil, fmt.Errorf("failed to get clause template: %w", err)
	}

	return &clause, nil
}

// GetByClauseCode retrieves a clause template by clause code
func (r *ClauseTemplateRepository) GetByClauseCode(clauseCode string) (*models.ClauseTemplate, error) {
	query := `
		SELECT id, clause_code, title, type, content, is_active, created_at, updated_at
		FROM clause_templates
		WHERE clause_code = $1
	`

	var clause models.ClauseTemplate
	err := r.db.QueryRow(query, clauseCode).Scan(
		&clause.ID,
		&clause.ClauseCode,
		&clause.Title,
		&clause.Type,
		&clause.Content,
		&clause.IsActive,
		&clause.CreatedAt,
		&clause.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("clause template not found")
		}
		return nil, fmt.Errorf("failed to get clause template: %w", err)
	}

	return &clause, nil
}

// Update updates an existing clause template
func (r *ClauseTemplateRepository) Update(clause *models.ClauseTemplate) (*models.ClauseTemplate, error) {
	query := `
		UPDATE clause_templates
		SET clause_code = $1, title = $2, type = $3, content = $4, is_active = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		clause.ClauseCode,
		clause.Title,
		clause.Type,
		clause.Content,
		clause.IsActive,
		clause.ID,
	).Scan(&clause.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("clause template not found")
		}
		return nil, fmt.Errorf("failed to update clause template: %w", err)
	}

	return clause, nil
}

// Delete deletes a clause template by ID
func (r *ClauseTemplateRepository) Delete(id int) error {
	query := `DELETE FROM clause_templates WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete clause template: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("clause template not found")
	}

	return nil
}

// List retrieves clause templates with pagination and search
func (r *ClauseTemplateRepository) List(params models.ClauseTemplateSearchRequest) (*models.ClauseTemplateListResponse, error) {
	// Build WHERE conditions
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Search query (title and content)
	if params.Query != "" {
		conditions = append(conditions, fmt.Sprintf("to_tsvector('english', title || ' ' || content) @@ plainto_tsquery('english', $%d)", argIndex))
		args = append(args, params.Query)
		argIndex++
	}

	// Filter by type
	if params.Type != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argIndex))
		args = append(args, params.Type)
		argIndex++
	}

	// Filter by active status
	if params.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *params.IsActive)
		argIndex++
	}

	// Build WHERE clause
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM clause_templates %s", whereClause)
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count clause templates: %w", err)
	}

	// Build ORDER BY clause
	orderBy := fmt.Sprintf("ORDER BY %s %s", params.SortBy, strings.ToUpper(params.SortDir))

	// Build LIMIT and OFFSET
	offset := (params.Page - 1) * params.Limit
	limitOffset := fmt.Sprintf("LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, params.Limit, offset)

	// Build final query
	query := fmt.Sprintf(`
		SELECT id, clause_code, title, type, content, is_active, created_at, updated_at
		FROM clause_templates
		%s
		%s
		%s
	`, whereClause, orderBy, limitOffset)

	// Execute query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query clause templates: %w", err)
	}
	defer rows.Close()

	// Scan results
	var clauses []models.ClauseTemplate
	for rows.Next() {
		var clause models.ClauseTemplate
		err := rows.Scan(
			&clause.ID,
			&clause.ClauseCode,
			&clause.Title,
			&clause.Type,
			&clause.Content,
			&clause.IsActive,
			&clause.CreatedAt,
			&clause.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan clause template: %w", err)
		}
		clauses = append(clauses, clause)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating clause templates: %w", err)
	}

	// Calculate pagination info
	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))
	hasNext := params.Page < totalPages
	hasPrev := params.Page > 1

	pagination := models.PaginationInfo{
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}

	return &models.ClauseTemplateListResponse{
		ClauseTemplates: clauses,
		Pagination:      pagination,
	}, nil
}

// GetActiveTypes retrieves all unique types from active clause templates
func (r *ClauseTemplateRepository) GetActiveTypes() ([]string, error) {
	query := `
		SELECT DISTINCT type
		FROM clause_templates
		WHERE is_active = TRUE
		ORDER BY type
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get clause types: %w", err)
	}
	defer rows.Close()

	var types []string
	for rows.Next() {
		var clauseType string
		err := rows.Scan(&clauseType)
		if err != nil {
			return nil, fmt.Errorf("failed to scan clause type: %w", err)
		}
		types = append(types, clauseType)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating clause types: %w", err)
	}

	return types, nil
}
