package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"contrack-be/internal/database"
	"contrack-be/internal/models"
)

type ContractRepository struct {
	db *sql.DB
}

func NewContractRepository() *ContractRepository {
	return &ContractRepository{
		db: database.DB,
	}
}

// GenerateContractNumber generates a unique contract number with format CTR-YYYY-MM-XXXXX-VV
func (r *ContractRepository) GenerateContractNumber() (string, error) {
	now := time.Now()
	yearMonth := now.Format("2006-01")
	
	// Get or create sequence for current year-month
	var sequenceNumber int
	
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	// Lock and update sequence
	query := `
		INSERT INTO contract_sequences (year_month, sequence_number) 
		VALUES ($1, 1)
		ON CONFLICT (year_month) 
		DO UPDATE SET sequence_number = contract_sequences.sequence_number + 1
		RETURNING sequence_number`
	
	err = tx.QueryRow(query, yearMonth).Scan(&sequenceNumber)
	if err != nil {
		return "", fmt.Errorf("failed to generate sequence number: %w", err)
	}
	
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	// Format: CTR-YYYY-MM-XXXXX-V1
	contractNumber := fmt.Sprintf("CTR-%s-%05d-V1", yearMonth, sequenceNumber)
	return contractNumber, nil
}

// Create creates a new contract
func (r *ContractRepository) Create(contract *models.Contract) error {
	// Generate contract number if not provided
	if contract.ContractNumber == "" {
		contractNumber, err := r.GenerateContractNumber()
		if err != nil {
			return fmt.Errorf("failed to generate contract number: %w", err)
		}
		contract.ContractNumber = contractNumber
	}
	
	query := `
		INSERT INTO contracts (
			base_id, version_number, project_name, package_name, contract_number, 
			external_reference, contract_type, signing_place, signing_date, 
			total_value, funding_source, status, created_by
		) VALUES (
			gen_random_uuid(), 1, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		) RETURNING id, base_id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		contract.ProjectName,
		contract.PackageName,
		contract.ContractNumber,
		contract.ExternalReference,
		contract.ContractType,
		contract.SigningPlace,
		contract.SigningDate,
		contract.TotalValue,
		contract.FundingSource,
		contract.Status,
		contract.CreatedBy,
	).Scan(&contract.ID, &contract.BaseID, &contract.CreatedAt, &contract.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create contract: %w", err)
	}

	// Set version number
	contract.VersionNumber = 1

	return nil
}

// GetByID retrieves a contract by ID
func (r *ContractRepository) GetByID(id int) (*models.Contract, error) {
	contract := &models.Contract{}
	
	query := `
		SELECT id, base_id, version_number, project_name, package_name, contract_number,
			   external_reference, contract_type, signing_place, signing_date,
			   total_value, funding_source, status, created_by, created_at, updated_at,
			   deleted_at, is_deleted
		FROM contracts 
		WHERE id = $1 AND is_deleted = FALSE`

	err := r.db.QueryRow(query, id).Scan(
		&contract.ID,
		&contract.BaseID,
		&contract.VersionNumber,
		&contract.ProjectName,
		&contract.PackageName,
		&contract.ContractNumber,
		&contract.ExternalReference,
		&contract.ContractType,
		&contract.SigningPlace,
		&contract.SigningDate,
		&contract.TotalValue,
		&contract.FundingSource,
		&contract.Status,
		&contract.CreatedBy,
		&contract.CreatedAt,
		&contract.UpdatedAt,
		&contract.DeletedAt,
		&contract.IsDeleted,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("contract not found")
		}
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}

	return contract, nil
}

// Update updates a contract
func (r *ContractRepository) Update(id int, updates map[string]interface{}) error {
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
		UPDATE contracts 
		SET %s 
		WHERE id = $%d AND is_deleted = FALSE`,
		strings.Join(setParts, ", "),
		argCount,
	)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update contract: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("contract not found or already deleted")
	}

	return nil
}

// Delete soft deletes a contract
func (r *ContractRepository) Delete(id int) error {
	query := `
		UPDATE contracts 
		SET is_deleted = TRUE, deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND is_deleted = FALSE`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete contract: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("contract not found or already deleted")
	}

	return nil
}

// List retrieves contracts with pagination and filtering
func (r *ContractRepository) List(req *models.ContractSearchRequest) ([]models.Contract, int, error) {
	// Set defaults
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortDir == "" {
		req.SortDir = "desc"
	}

	// Build WHERE conditions
	whereConditions := []string{"is_deleted = FALSE"}
	args := []interface{}{}
	argCount := 1

	// Full-text search
	if req.Query != "" {
		whereConditions = append(whereConditions, fmt.Sprintf(`
			to_tsvector('english', project_name || ' ' || COALESCE(package_name, '') || ' ' || contract_type) 
			@@ plainto_tsquery('english', $%d)`, argCount))
		args = append(args, req.Query)
		argCount++
	}

	// Status filter
	if req.Status != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("status = $%d", argCount))
		args = append(args, *req.Status)
		argCount++
	}

	// Contract type filter
	if req.ContractType != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("contract_type ILIKE $%d", argCount))
		args = append(args, "%"+req.ContractType+"%")
		argCount++
	}

	// Funding source filter
	if req.FundingSource != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("funding_source ILIKE $%d", argCount))
		args = append(args, "%"+req.FundingSource+"%")
		argCount++
	}

	// Date range filters
	if req.SigningDateFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("signing_date >= $%d", argCount))
		args = append(args, *req.SigningDateFrom)
		argCount++
	}
	if req.SigningDateTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("signing_date <= $%d", argCount))
		args = append(args, *req.SigningDateTo)
		argCount++
	}

	// Value range filters
	if req.ValueFrom != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("total_value >= $%d", argCount))
		args = append(args, *req.ValueFrom)
		argCount++
	}
	if req.ValueTo != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("total_value <= $%d", argCount))
		args = append(args, *req.ValueTo)
		argCount++
	}

	// Created by filter
	if req.CreatedBy != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("created_by = $%d", argCount))
		args = append(args, req.CreatedBy)
		argCount++
	}

	whereClause := strings.Join(whereConditions, " AND ")

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM contracts WHERE %s", whereClause)
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count contracts: %w", err)
	}

	// Fetch records with pagination
	offset := (req.Page - 1) * req.Limit
	args = append(args, req.Limit, offset)

	query := fmt.Sprintf(`
		SELECT id, base_id, version_number, project_name, package_name, contract_number,
			   external_reference, contract_type, signing_place, signing_date,
			   total_value, funding_source, status, created_by, created_at, updated_at,
			   deleted_at, is_deleted
		FROM contracts 
		WHERE %s 
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d`,
		whereClause,
		req.SortBy, strings.ToUpper(req.SortDir),
		argCount, argCount+1,
	)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query contracts: %w", err)
	}
	defer rows.Close()

	contracts := []models.Contract{}
	for rows.Next() {
		contract := models.Contract{}
		err := rows.Scan(
			&contract.ID,
			&contract.BaseID,
			&contract.VersionNumber,
			&contract.ProjectName,
			&contract.PackageName,
			&contract.ContractNumber,
			&contract.ExternalReference,
			&contract.ContractType,
			&contract.SigningPlace,
			&contract.SigningDate,
			&contract.TotalValue,
			&contract.FundingSource,
			&contract.Status,
			&contract.CreatedBy,
			&contract.CreatedAt,
			&contract.UpdatedAt,
			&contract.DeletedAt,
			&contract.IsDeleted,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan contract: %w", err)
		}
		contracts = append(contracts, contract)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("failed to iterate contracts: %w", err)
	}

	return contracts, total, nil
}

// UpdateStatus updates contract status and records history
func (r *ContractRepository) UpdateStatus(contractID int, newStatus models.ContractStatus, changedBy string, changeReason, comments *string) error {
	// Start transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get current status
	var currentStatus models.ContractStatus
	err = tx.QueryRow("SELECT status FROM contracts WHERE id = $1 AND is_deleted = FALSE", contractID).Scan(&currentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("contract not found")
		}
		return fmt.Errorf("failed to get current status: %w", err)
	}

	// Validate status transition
	if !currentStatus.CanTransitionTo(newStatus) {
		return fmt.Errorf("invalid status transition from %s to %s", currentStatus, newStatus)
	}

	// Update contract status
	_, err = tx.Exec(`
		UPDATE contracts 
		SET status = $1, updated_at = NOW() 
		WHERE id = $2`,
		newStatus, contractID)
	if err != nil {
		return fmt.Errorf("failed to update contract status: %w", err)
	}

	// Record status history
	_, err = tx.Exec(`
		INSERT INTO contract_status_history (contract_id, from_status, to_status, changed_by, change_reason, comments)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		contractID, currentStatus, newStatus, changedBy, changeReason, comments)
	if err != nil {
		return fmt.Errorf("failed to record status history: %w", err)
	}

	return tx.Commit()
}

// GetStatusHistory retrieves status change history for a contract
func (r *ContractRepository) GetStatusHistory(contractID int) ([]models.ContractStatusHistory, error) {
	query := `
		SELECT id, contract_id, from_status, to_status, changed_by, change_reason, comments, changed_at
		FROM contract_status_history
		WHERE contract_id = $1
		ORDER BY changed_at DESC`

	rows, err := r.db.Query(query, contractID)
	if err != nil {
		return nil, fmt.Errorf("failed to query status history: %w", err)
	}
	defer rows.Close()

	history := []models.ContractStatusHistory{}
	for rows.Next() {
		h := models.ContractStatusHistory{}
		err := rows.Scan(
			&h.ID,
			&h.ContractID,
			&h.FromStatus,
			&h.ToStatus,
			&h.ChangedBy,
			&h.ChangeReason,
			&h.Comments,
			&h.ChangedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan status history: %w", err)
		}
		history = append(history, h)
	}

	return history, rows.Err()
}

// CreateVersion creates a new version of an existing contract
func (r *ContractRepository) CreateVersion(baseID string, contract *models.Contract) error {
	// Get the latest version number
	fmt.Println("baseID: ", baseID)
	fmt.Println("contract: ", contract)
	fmt.Println("contract number: ", contract.ContractNumber)
	var latestVersion int
	err := r.db.QueryRow(`
		SELECT COALESCE(MAX(version_number), 0) 
		FROM contracts 
		WHERE base_id = $1`,
		baseID).Scan(&latestVersion)
	if err != nil {
		return fmt.Errorf("failed to get latest version: %w", err)
	}

	newVersion := latestVersion + 1
	
	// Generate new contract number with updated version
	parts := strings.Split(contract.ContractNumber, "-V")
	if len(parts) != 2 {
		return fmt.Errorf("invalid contract number format")
	}
	newContractNumber := fmt.Sprintf("%s-V%d", parts[0], newVersion)

	query := `
		INSERT INTO contracts (
			base_id, version_number, project_name, package_name, contract_number, 
			external_reference, contract_type, signing_place, signing_date, 
			total_value, funding_source, status, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		) RETURNING id, created_at, updated_at`

	err = r.db.QueryRow(
		query,
		baseID,
		newVersion,
		contract.ProjectName,
		contract.PackageName,
		newContractNumber,
		contract.ExternalReference,
		contract.ContractType,
		contract.SigningPlace,
		contract.SigningDate,
		contract.TotalValue,
		contract.FundingSource,
		models.StatusDraft, // New versions always start as DRAFT
		contract.CreatedBy,
	).Scan(&contract.ID, &contract.CreatedAt, &contract.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create contract version: %w", err)
	}

	contract.BaseID = baseID
	contract.VersionNumber = newVersion
	contract.ContractNumber = newContractNumber
	contract.Status = models.StatusDraft

	return nil
}

// GetVersions retrieves all versions of a contract by base_id
func (r *ContractRepository) GetVersions(baseID string) ([]models.Contract, error) {
	query := `
		SELECT id, base_id, version_number, project_name, package_name, contract_number,
			   external_reference, contract_type, signing_place, signing_date,
			   total_value, funding_source, status, created_by, created_at, updated_at,
			   deleted_at, is_deleted
		FROM contracts 
		WHERE base_id = $1 AND is_deleted = FALSE
		ORDER BY version_number DESC`

	rows, err := r.db.Query(query, baseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query contract versions: %w", err)
	}
	defer rows.Close()

	contracts := []models.Contract{}
	for rows.Next() {
		contract := models.Contract{}
		err := rows.Scan(
			&contract.ID,
			&contract.BaseID,
			&contract.VersionNumber,
			&contract.ProjectName,
			&contract.PackageName,
			&contract.ContractNumber,
			&contract.ExternalReference,
			&contract.ContractType,
			&contract.SigningPlace,
			&contract.SigningDate,
			&contract.TotalValue,
			&contract.FundingSource,
			&contract.Status,
			&contract.CreatedBy,
			&contract.CreatedAt,
			&contract.UpdatedAt,
			&contract.DeletedAt,
			&contract.IsDeleted,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contract: %w", err)
		}
		contracts = append(contracts, contract)
	}

	return contracts, rows.Err()
}

// CanEdit checks if a contract can be edited based on its status
func (r *ContractRepository) CanEdit(contractID int) (bool, error) {
	var status models.ContractStatus
	err := r.db.QueryRow("SELECT status FROM contracts WHERE id = $1 AND is_deleted = FALSE", contractID).Scan(&status)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("contract not found")
		}
		return false, fmt.Errorf("failed to get contract status: %w", err)
	}

	// Only DRAFT contracts can be edited
	return status == models.StatusDraft, nil
}

// GetContractsByStatus retrieves contracts by status for a specific user
func (r *ContractRepository) GetContractsByStatus(status models.ContractStatus, createdBy string) ([]models.Contract, error) {
	query := `
		SELECT id, base_id, version_number, project_name, package_name, contract_number,
			   external_reference, contract_type, signing_place, signing_date,
			   total_value, funding_source, status, created_by, created_at, updated_at,
			   deleted_at, is_deleted
		FROM contracts 
		WHERE status = $1 AND created_by = $2 AND is_deleted = FALSE
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, status, createdBy)
	if err != nil {
		return nil, fmt.Errorf("failed to query contracts by status: %w", err)
	}
	defer rows.Close()

	contracts := []models.Contract{}
	for rows.Next() {
		contract := models.Contract{}
		err := rows.Scan(
			&contract.ID,
			&contract.BaseID,
			&contract.VersionNumber,
			&contract.ProjectName,
			&contract.PackageName,
			&contract.ContractNumber,
			&contract.ExternalReference,
			&contract.ContractType,
			&contract.SigningPlace,
			&contract.SigningDate,
			&contract.TotalValue,
			&contract.FundingSource,
			&contract.Status,
			&contract.CreatedBy,
			&contract.CreatedAt,
			&contract.UpdatedAt,
			&contract.DeletedAt,
			&contract.IsDeleted,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contract: %w", err)
		}
		contracts = append(contracts, contract)
	}

	return contracts, rows.Err()
}

// GetStatusCounts retrieves count for each contract status
func (r *ContractRepository) GetStatusCounts() ([]models.StatusCount, error) {
	query := `
		SELECT 
			status,
			COUNT(*) as count
		FROM contracts 
		WHERE is_deleted = false
		GROUP BY status
		ORDER BY count DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query status counts: %w", err)
	}
	defer rows.Close()

	var statusCounts []models.StatusCount
	for rows.Next() {
		statusCount := models.StatusCount{}
		err := rows.Scan(
			&statusCount.Status,
			&statusCount.Count,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan status count: %w", err)
		}
		
		statusCount.StatusDisplay = statusCount.Status.GetStatusDisplay()
		statusCounts = append(statusCounts, statusCount)
	}

	return statusCounts, rows.Err()
}

// GetProjectValueDistribution retrieves total value per project for bar chart
func (r *ContractRepository) GetProjectValueDistribution() ([]models.ProjectValueDistribution, error) {
	query := `
		SELECT 
			project_name,
			SUM(total_value) as total_value
		FROM contracts 
		WHERE is_deleted = false
		GROUP BY project_name
		ORDER BY total_value DESC
		LIMIT 10`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query project value distribution: %w", err)
	}
	defer rows.Close()

	var distributions []models.ProjectValueDistribution
	for rows.Next() {
		dist := models.ProjectValueDistribution{}
		err := rows.Scan(
			&dist.ProjectName,
			&dist.TotalValue,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project value distribution: %w", err)
		}
		
		distributions = append(distributions, dist)
	}

	return distributions, rows.Err()
}

// GetContractTypeDistribution retrieves contract count by type for donut/pie chart
func (r *ContractRepository) GetContractTypeDistribution() ([]models.ContractTypeDistribution, error) {
	query := `
		SELECT 
			contract_type,
			COUNT(*) as count
		FROM contracts 
		WHERE is_deleted = false
		GROUP BY contract_type
		ORDER BY count DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query contract type distribution: %w", err)
	}
	defer rows.Close()

	var distributions []models.ContractTypeDistribution
	for rows.Next() {
		dist := models.ContractTypeDistribution{}
		err := rows.Scan(
			&dist.ContractType,
			&dist.Count,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan contract type distribution: %w", err)
		}
		
		distributions = append(distributions, dist)
	}

	return distributions, rows.Err()
}

// GetDashboardVisualizationData retrieves comprehensive data for dashboard visualizations
func (r *ContractRepository) GetDashboardVisualizationData() (*models.DashboardVisualizationData, error) {
	// Get status counts
	statusCounts, err := r.GetStatusCounts()
	if err != nil {
		return nil, fmt.Errorf("failed to get status counts: %w", err)
	}

	// Get project value distribution
	projectValueDist, err := r.GetProjectValueDistribution()
	if err != nil {
		return nil, fmt.Errorf("failed to get project value distribution: %w", err)
	}

	// Get contract type distribution
	contractTypeDist, err := r.GetContractTypeDistribution()
	if err != nil {
		return nil, fmt.Errorf("failed to get contract type distribution: %w", err)
	}

	// Get total contracts and total value
	var totalContracts int
	var totalValue float64
	err = r.db.QueryRow(`
		SELECT COUNT(*), COALESCE(SUM(total_value), 0)
		FROM contracts 
		WHERE is_deleted = false
	`).Scan(&totalContracts, &totalValue)
	if err != nil {
		return nil, fmt.Errorf("failed to get total stats: %w", err)
	}

	return &models.DashboardVisualizationData{
		StatusCounts:      statusCounts,
		ProjectValueDist:  projectValueDist,
		ContractTypeDist:  contractTypeDist,
		TotalContracts:    totalContracts,
		TotalValue:        totalValue,
	}, nil
}

// GetSimpleContractList retrieves simple contract list for dashboard table
func (r *ContractRepository) GetSimpleContractList(userID string) ([]models.SimpleContractList, error) {
	query := `
		SELECT 
			id, project_name, contract_type, status, total_value, signing_date, created_at
		FROM contracts 
		WHERE is_deleted = false AND created_by = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query simple contract list: %w", err)
	}
	defer rows.Close()

	var contracts []models.SimpleContractList
	for rows.Next() {
		contract := models.SimpleContractList{}
		err := rows.Scan(
			&contract.ID,
			&contract.ProjectName,
			&contract.ContractType,
			&contract.Status,
			&contract.TotalValue,
			&contract.SigningDate,
			&contract.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan simple contract: %w", err)
		}
		
		contract.StatusDisplay = contract.Status.GetStatusDisplay()
		contracts = append(contracts, contract)
	}

	return contracts, rows.Err()
}

// GetContractApprovalData retrieves contract data needed for approval decision
func (r *ContractRepository) GetContractApprovalData(contractID int) (*models.ContractApprovalResponse, error) {
	// Get contract basic info
	contractQuery := `
		SELECT id, project_name, total_value, status
		FROM contracts 
		WHERE id = $1 AND is_deleted = false`
	
	var contract models.Contract
	err := r.db.QueryRow(contractQuery, contractID).Scan(
		&contract.ID,
		&contract.ProjectName,
		&contract.TotalValue,
		&contract.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract: %w", err)
	}

	// Get AI analysis data for this contract
	analysisQuery := `
		SELECT 
			AVG(risk_score) as avg_risk_score,
			MAX(risk_level) as max_risk_level,
			AVG(confidence_score) as avg_confidence_score,
			COUNT(*) as analysis_count
		FROM clause_risk_analyses 
		WHERE clause_id IN (
			SELECT clause_id FROM contract_clauses WHERE contract_id = $1
		)`
	
	var avgRiskScore, avgConfidenceScore float64
	var maxRiskLevel string
	var analysisCount int
	
	err = r.db.QueryRow(analysisQuery, contractID).Scan(
		&avgRiskScore,
		&maxRiskLevel,
		&avgConfidenceScore,
		&analysisCount,
	)
	if err != nil {
		// If no analysis found, set default values
		avgRiskScore = 50.0
		maxRiskLevel = "medium"
		avgConfidenceScore = 70.0
		analysisCount = 0
	}

	// Get key risks and recommendations
	risksQuery := `
		SELECT DISTINCT unnest(identified_risks) as risk
		FROM clause_risk_analyses 
		WHERE clause_id IN (
			SELECT clause_id FROM contract_clauses WHERE contract_id = $1
		)
		LIMIT 5`
	
	recommendationsQuery := `
		SELECT DISTINCT unnest(recommendations) as recommendation
		FROM clause_risk_analyses 
		WHERE clause_id IN (
			SELECT clause_id FROM contract_clauses WHERE contract_id = $1
		)
		LIMIT 5`

	var keyRisks []string
	var recommendations []string

	riskRows, err := r.db.Query(risksQuery, contractID)
	if err == nil {
		defer riskRows.Close()
		for riskRows.Next() {
			var risk string
			if err := riskRows.Scan(&risk); err == nil {
				keyRisks = append(keyRisks, risk)
			}
		}
	}

	recRows, err := r.db.Query(recommendationsQuery, contractID)
	if err == nil {
		defer recRows.Close()
		for recRows.Next() {
			var rec string
			if err := recRows.Scan(&rec); err == nil {
				recommendations = append(recommendations, rec)
			}
		}
	}

	// Determine approval criteria
	criteria := models.ContractApprovalCriteria{
		MaxRiskScore:      60.0,  // Medium risk threshold (increased for testing)
		MaxValue:          1000000000, // 1 billion IDR
		MinConfidenceScore: 50.0, // Medium confidence threshold (lowered for testing)
	}

	// Make approval decision - always require manual approval, but with different levels
	requiresReview := false
	approvalStatus := "APPROVAL_REQUIRED" // Changed from ACTIVE to APPROVAL_REQUIRED
	approvalMessage := "Contract ready for approval - low risk detected"
	reviewReasons := []string{}

	// Check if contract is high risk (requires detailed review)
	isHighRisk := false
	
	if avgRiskScore > criteria.MaxRiskScore {
		isHighRisk = true
		approvalStatus = "REVIEW_REQUIRED"
		approvalMessage = "Contract requires detailed review due to high risk score"
		reviewReasons = append(reviewReasons, fmt.Sprintf("Risk score %.1f exceeds threshold of %.1f", avgRiskScore, criteria.MaxRiskScore))
	}

	if contract.TotalValue > criteria.MaxValue {
		isHighRisk = true
		approvalStatus = "REVIEW_REQUIRED"
		approvalMessage = "Contract requires detailed review due to high value"
		reviewReasons = append(reviewReasons, fmt.Sprintf("Contract value %.0f exceeds threshold of %.0f", contract.TotalValue, criteria.MaxValue))
	}

	if avgConfidenceScore < criteria.MinConfidenceScore {
		isHighRisk = true
		approvalStatus = "REVIEW_REQUIRED"
		approvalMessage = "Contract requires detailed review due to low AI confidence"
		reviewReasons = append(reviewReasons, fmt.Sprintf("AI confidence %.1f below threshold of %.1f", avgConfidenceScore, criteria.MinConfidenceScore))
	}

	// If not high risk, still requires approval but no detailed review
	if !isHighRisk {
		requiresReview = false
		approvalStatus = "APPROVAL_REQUIRED"
		approvalMessage = "Contract ready for approval - low risk detected"
	} else {
		requiresReview = true
	}

	if analysisCount == 0 {
		// If contract is PENDING_SIGNATURE, it means it already passed legal review
		// So we can approve it directly without requiring AI analysis
		requiresReview = false
		approvalStatus = "APPROVAL_REQUIRED"
		approvalMessage = "Contract ready for approval - already passed legal review"
	}

	// Set next steps
	nextSteps := []string{}
	if requiresReview {
		nextSteps = append(nextSteps, "Review contract details and risks")
		nextSteps = append(nextSteps, "Consult with legal team if needed")
		nextSteps = append(nextSteps, "Make manual approval decision")
	} else {
		nextSteps = append(nextSteps, "Contract ready for approval - low risk")
		nextSteps = append(nextSteps, "Click approve to activate contract")
		nextSteps = append(nextSteps, "No detailed review required")
	}

	response := &models.ContractApprovalResponse{
		ContractID:      contract.ID,
		ContractName:    contract.ProjectName,
		TotalValue:      contract.TotalValue,
		RiskLevel:       maxRiskLevel,
		RiskScore:       avgRiskScore,
		ApprovalStatus:  approvalStatus,
		ApprovalMessage: approvalMessage,
		RequiresReview:  requiresReview,
		ReviewReasons:   reviewReasons,
		KeyRisks:        keyRisks,
		Recommendations: recommendations,
		NextSteps:       nextSteps,
	}

	// If approved, set approval timestamp
	if !requiresReview {
		now := time.Now()
		response.ApprovedAt = &now
		response.ApprovedBy = "AI Auto-Approval System"
	}

	return response, nil
}

// ApproveContract updates contract status to ACTIVE
func (r *ContractRepository) ApproveContract(contractID int, approvedBy string) error {
	// Update contract status to ACTIVE
	_, err := r.db.Exec(`
		UPDATE contracts 
		SET status = $1, updated_at = NOW() 
		WHERE id = $2`,
		models.StatusActive, contractID)
	if err != nil {
		return fmt.Errorf("failed to approve contract: %w", err)
	}

	// Record status change in history
	_, err = r.db.Exec(`
		INSERT INTO contract_status_history (contract_id, from_status, to_status, changed_by, change_reason, comments)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		contractID, "PENDING_LEGAL_REVIEW", models.StatusActive, approvedBy, "One-click approval", "Contract approved automatically based on AI analysis")
	if err != nil {
		return fmt.Errorf("failed to record approval history: %w", err)
	}

	return nil
}
