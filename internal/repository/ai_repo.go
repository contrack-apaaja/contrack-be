package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"contrack-be/internal/database"
	"contrack-be/internal/models"
)

// AIRepository handles AI analysis database operations
type AIRepository struct {
	db *sql.DB
}

// NewAIRepository creates a new AI repository instance
func NewAIRepository() *AIRepository {
	return &AIRepository{db: database.DB}
}

// CreateAnalysis creates a new AI analysis record
func (r *AIRepository) CreateAnalysis(analysis *models.ClauseRiskAnalysis) error {
	query := `
		INSERT INTO clause_risk_analyses (
			clause_id, risk_level, risk_score, analysis_summary, 
			identified_risks, recommendations, legal_implications, 
			compliance_notes, confidence_score, model_version, 
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id`

	// Convert slices to JSON strings
	identifiedRisksJSON, err := json.Marshal(analysis.IdentifiedRisks)
	if err != nil {
		return fmt.Errorf("failed to marshal identified risks: %w", err)
	}

	recommendationsJSON, err := json.Marshal(analysis.Recommendations)
	if err != nil {
		return fmt.Errorf("failed to marshal recommendations: %w", err)
	}

	now := time.Now()
	analysis.CreatedAt = now
	analysis.UpdatedAt = now

	err = r.db.QueryRow(
		query,
		analysis.ClauseID,
		analysis.RiskLevel,
		analysis.RiskScore,
		analysis.AnalysisSummary,
		string(identifiedRisksJSON),
		string(recommendationsJSON),
		analysis.LegalImplications,
		analysis.ComplianceNotes,
		analysis.ConfidenceScore,
		analysis.ModelVersion,
		analysis.CreatedAt,
		analysis.UpdatedAt,
	).Scan(&analysis.ID)

	if err != nil {
		return fmt.Errorf("failed to create analysis: %w", err)
	}

	return nil
}

// GetAnalysisByID retrieves an analysis by its ID
func (r *AIRepository) GetAnalysisByID(id int) (*models.ClauseRiskAnalysis, error) {
	query := `
		SELECT id, clause_id, risk_level, risk_score, analysis_summary,
			   identified_risks, recommendations, legal_implications,
			   compliance_notes, confidence_score, model_version,
			   created_at, updated_at
		FROM clause_risk_analyses
		WHERE id = $1`

	analysis := &models.ClauseRiskAnalysis{}
	var identifiedRisksJSON, recommendationsJSON string

	err := r.db.QueryRow(query, id).Scan(
		&analysis.ID,
		&analysis.ClauseID,
		&analysis.RiskLevel,
		&analysis.RiskScore,
		&analysis.AnalysisSummary,
		&identifiedRisksJSON,
		&recommendationsJSON,
		&analysis.LegalImplications,
		&analysis.ComplianceNotes,
		&analysis.ConfidenceScore,
		&analysis.ModelVersion,
		&analysis.CreatedAt,
		&analysis.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("analysis not found")
		}
		return nil, fmt.Errorf("failed to get analysis: %w", err)
	}

	// Parse JSON strings back to slices
	if err := json.Unmarshal([]byte(identifiedRisksJSON), &analysis.IdentifiedRisks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal identified risks: %w", err)
	}

	if err := json.Unmarshal([]byte(recommendationsJSON), &analysis.Recommendations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal recommendations: %w", err)
	}

	return analysis, nil
}

// GetAnalysisByClauseID retrieves the latest analysis for a specific clause
func (r *AIRepository) GetAnalysisByClauseID(clauseID int) (*models.ClauseRiskAnalysis, error) {
	query := `
		SELECT id, clause_id, risk_level, risk_score, analysis_summary,
			   identified_risks, recommendations, legal_implications,
			   compliance_notes, confidence_score, model_version,
			   created_at, updated_at
		FROM clause_risk_analyses
		WHERE clause_id = $1
		ORDER BY created_at DESC
		LIMIT 1`

	analysis := &models.ClauseRiskAnalysis{}
	var identifiedRisksJSON, recommendationsJSON string

	err := r.db.QueryRow(query, clauseID).Scan(
		&analysis.ID,
		&analysis.ClauseID,
		&analysis.RiskLevel,
		&analysis.RiskScore,
		&analysis.AnalysisSummary,
		&identifiedRisksJSON,
		&recommendationsJSON,
		&analysis.LegalImplications,
		&analysis.ComplianceNotes,
		&analysis.ConfidenceScore,
		&analysis.ModelVersion,
		&analysis.CreatedAt,
		&analysis.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("analysis not found for clause")
		}
		return nil, fmt.Errorf("failed to get analysis: %w", err)
	}

	// Parse JSON strings back to slices
	if err := json.Unmarshal([]byte(identifiedRisksJSON), &analysis.IdentifiedRisks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal identified risks: %w", err)
	}

	if err := json.Unmarshal([]byte(recommendationsJSON), &analysis.Recommendations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal recommendations: %w", err)
	}

	return analysis, nil
}

// GetAnalyses retrieves analyses with pagination and filtering
func (r *AIRepository) GetAnalyses(params models.ClauseRiskAnalysisSearchRequest) (*models.ClauseRiskAnalysisListResponse, error) {
	// Build the WHERE clause
	whereConditions := []string{}
	args := []interface{}{}
	argIndex := 1

	if params.ClauseID > 0 {
		whereConditions = append(whereConditions, fmt.Sprintf("clause_id = $%d", argIndex))
		args = append(args, params.ClauseID)
		argIndex++
	}

	if params.RiskLevel != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("risk_level = $%d", argIndex))
		args = append(args, params.RiskLevel)
		argIndex++
	}

	if params.MinRiskScore >= 0 {
		whereConditions = append(whereConditions, fmt.Sprintf("risk_score >= $%d", argIndex))
		args = append(args, params.MinRiskScore)
		argIndex++
	}

	if params.MaxRiskScore > 0 {
		whereConditions = append(whereConditions, fmt.Sprintf("risk_score <= $%d", argIndex))
		args = append(args, params.MaxRiskScore)
		argIndex++
	}

	if params.MinConfidence >= 0 {
		whereConditions = append(whereConditions, fmt.Sprintf("confidence_score >= $%d", argIndex))
		args = append(args, params.MinConfidence)
		argIndex++
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Build the ORDER BY clause
	orderBy := "created_at DESC"
	if params.SortBy != "" {
		orderBy = params.SortBy
		if params.SortDir != "" {
			orderBy += " " + strings.ToUpper(params.SortDir)
		}
	}

	// Count total records
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM clause_risk_analyses %s", whereClause)
	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count analyses: %w", err)
	}

	// Calculate pagination
	offset := (params.Page - 1) * params.Limit
	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	// Build the main query
	query := fmt.Sprintf(`
		SELECT id, clause_id, risk_level, risk_score, analysis_summary,
			   identified_risks, recommendations, legal_implications,
			   compliance_notes, confidence_score, model_version,
			   created_at, updated_at
		FROM clause_risk_analyses
		%s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		whereClause, orderBy, argIndex, argIndex+1)

	args = append(args, params.Limit, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query analyses: %w", err)
	}
	defer rows.Close()

	var analyses []models.ClauseRiskAnalysis
	for rows.Next() {
		analysis := models.ClauseRiskAnalysis{}
		var identifiedRisksJSON, recommendationsJSON string

		err := rows.Scan(
			&analysis.ID,
			&analysis.ClauseID,
			&analysis.RiskLevel,
			&analysis.RiskScore,
			&analysis.AnalysisSummary,
			&identifiedRisksJSON,
			&recommendationsJSON,
			&analysis.LegalImplications,
			&analysis.ComplianceNotes,
			&analysis.ConfidenceScore,
			&analysis.ModelVersion,
			&analysis.CreatedAt,
			&analysis.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan analysis: %w", err)
		}

		// Parse JSON strings back to slices
		if err := json.Unmarshal([]byte(identifiedRisksJSON), &analysis.IdentifiedRisks); err != nil {
			return nil, fmt.Errorf("failed to unmarshal identified risks: %w", err)
		}

		if err := json.Unmarshal([]byte(recommendationsJSON), &analysis.Recommendations); err != nil {
			return nil, fmt.Errorf("failed to unmarshal recommendations: %w", err)
		}

		analyses = append(analyses, analysis)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating analyses: %w", err)
	}

	response := &models.ClauseRiskAnalysisListResponse{
		Analyses: analyses,
		Pagination: models.PaginationInfo{
			Page:       params.Page,
			Limit:      params.Limit,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    params.Page < totalPages,
			HasPrev:    params.Page > 1,
		},
	}

	return response, nil
}

// UpdateAnalysis updates an existing analysis
func (r *AIRepository) UpdateAnalysis(analysis *models.ClauseRiskAnalysis) error {
	query := `
		UPDATE clause_risk_analyses SET
			risk_level = $2, risk_score = $3, analysis_summary = $4,
			identified_risks = $5, recommendations = $6, legal_implications = $7,
			compliance_notes = $8, confidence_score = $9, model_version = $10,
			updated_at = $11
		WHERE id = $1`

	// Convert slices to JSON strings
	identifiedRisksJSON, err := json.Marshal(analysis.IdentifiedRisks)
	if err != nil {
		return fmt.Errorf("failed to marshal identified risks: %w", err)
	}

	recommendationsJSON, err := json.Marshal(analysis.Recommendations)
	if err != nil {
		return fmt.Errorf("failed to marshal recommendations: %w", err)
	}

	analysis.UpdatedAt = time.Now()

	result, err := r.db.Exec(
		query,
		analysis.ID,
		analysis.RiskLevel,
		analysis.RiskScore,
		analysis.AnalysisSummary,
		string(identifiedRisksJSON),
		string(recommendationsJSON),
		analysis.LegalImplications,
		analysis.ComplianceNotes,
		analysis.ConfidenceScore,
		analysis.ModelVersion,
		analysis.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to update analysis: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("analysis not found")
	}

	return nil
}

// DeleteAnalysis deletes an analysis by ID
func (r *AIRepository) DeleteAnalysis(id int) error {
	query := "DELETE FROM clause_risk_analyses WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete analysis: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("analysis not found")
	}

	return nil
}

// GetAnalysisWithClause retrieves an analysis with its associated clause
func (r *AIRepository) GetAnalysisWithClause(analysisID int) (*models.ClauseRiskAnalysisResponse, error) {
	query := `
		SELECT 
			a.id, a.clause_id, a.risk_level, a.risk_score, a.analysis_summary,
			a.identified_risks, a.recommendations, a.legal_implications,
			a.compliance_notes, a.confidence_score, a.model_version,
			a.created_at, a.updated_at,
			c.id, c.clause_code, c.title, c.type, c.content, c.is_active,
			c.created_at, c.updated_at
		FROM clause_risk_analyses a
		JOIN clause_templates c ON a.clause_id = c.id
		WHERE a.id = $1`

	analysis := &models.ClauseRiskAnalysis{}
	clause := &models.ClauseTemplate{}
	var identifiedRisksJSON, recommendationsJSON string

	err := r.db.QueryRow(query, analysisID).Scan(
		&analysis.ID,
		&analysis.ClauseID,
		&analysis.RiskLevel,
		&analysis.RiskScore,
		&analysis.AnalysisSummary,
		&identifiedRisksJSON,
		&recommendationsJSON,
		&analysis.LegalImplications,
		&analysis.ComplianceNotes,
		&analysis.ConfidenceScore,
		&analysis.ModelVersion,
		&analysis.CreatedAt,
		&analysis.UpdatedAt,
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
			return nil, fmt.Errorf("analysis not found")
		}
		return nil, fmt.Errorf("failed to get analysis with clause: %w", err)
	}

	// Parse JSON strings back to slices
	if err := json.Unmarshal([]byte(identifiedRisksJSON), &analysis.IdentifiedRisks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal identified risks: %w", err)
	}

	if err := json.Unmarshal([]byte(recommendationsJSON), &analysis.Recommendations); err != nil {
		return nil, fmt.Errorf("failed to unmarshal recommendations: %w", err)
	}

	response := &models.ClauseRiskAnalysisResponse{
		Analysis: *analysis,
		Clause:   *clause,
	}

	return response, nil
}

// GetContractRecommendations retrieves all recommendations for a specific contract
func (r *AIRepository) GetContractRecommendations(contractID int) (*models.ContractRecommendations, error) {
	// For now, return all available AI analyses as recommendations
	// TODO: Implement proper contract-clause relationship when contract_clauses table is available
	
	query := `
		SELECT 
			id, clause_id, risk_level, risk_score, analysis_summary,
			identified_risks, recommendations, legal_implications, 
			compliance_notes, confidence_score, model_version, created_at
		FROM clause_risk_analyses 
		ORDER BY created_at DESC
		LIMIT 10
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query contract recommendations: %w", err)
	}
	defer rows.Close()
	
	var clauseRecommendations []models.ClauseRecommendation
	var totalRiskScore float64
	var maxRiskLevel models.RiskLevel = models.RiskLevelLow
	var totalClauses int
	
	for rows.Next() {
		var analysis models.ClauseRiskAnalysis
		var identifiedRisksJSON, recommendationsJSON string
		
		err := rows.Scan(
			&analysis.ID,
			&analysis.ClauseID,
			&analysis.RiskLevel,
			&analysis.RiskScore,
			&analysis.AnalysisSummary,
			&identifiedRisksJSON,
			&recommendationsJSON,
			&analysis.LegalImplications,
			&analysis.ComplianceNotes,
			&analysis.ConfidenceScore,
			&analysis.ModelVersion,
			&analysis.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan analysis: %w", err)
		}
		
		// Parse JSON arrays
		if err := json.Unmarshal([]byte(identifiedRisksJSON), &analysis.IdentifiedRisks); err != nil {
			analysis.IdentifiedRisks = []string{}
		}
		if err := json.Unmarshal([]byte(recommendationsJSON), &analysis.Recommendations); err != nil {
			analysis.Recommendations = []string{}
		}
		
		// Convert to ClauseRecommendation
		clauseRec := models.ClauseRecommendation{
			ClauseID:          analysis.ClauseID,
			RiskLevel:         analysis.RiskLevel,
			RiskScore:         analysis.RiskScore,
			AnalysisSummary:   analysis.AnalysisSummary,
			Recommendations:   analysis.Recommendations,
			IdentifiedRisks:   analysis.IdentifiedRisks,
			LegalImplications: analysis.LegalImplications,
			ComplianceNotes:   analysis.ComplianceNotes,
			ConfidenceScore:   analysis.ConfidenceScore,
			CreatedAt:         analysis.CreatedAt,
		}
		
		clauseRecommendations = append(clauseRecommendations, clauseRec)
		totalRiskScore += analysis.RiskScore
		totalClauses++
		
		// Update max risk level
		if analysis.RiskLevel == models.RiskLevelCritical || 
		   (analysis.RiskLevel == models.RiskLevelHigh && maxRiskLevel != models.RiskLevelCritical) ||
		   (analysis.RiskLevel == models.RiskLevelMedium && maxRiskLevel == models.RiskLevelLow) {
			maxRiskLevel = analysis.RiskLevel
		}
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating analyses: %w", err)
	}
	
	// Calculate overall metrics
	var overallRiskScore float64
	if totalClauses > 0 {
		overallRiskScore = totalRiskScore / float64(totalClauses)
	}
	
	// Generate overall recommendations and key risks
	var overallRecommendations []string
	var keyRisks []string
	
	for _, clause := range clauseRecommendations {
		overallRecommendations = append(overallRecommendations, clause.Recommendations...)
		keyRisks = append(keyRisks, clause.IdentifiedRisks...)
	}
	
	// Remove duplicates
	overallRecommendations = removeDuplicates(overallRecommendations)
	keyRisks = removeDuplicates(keyRisks)
	
	// Get the latest analysis date
	var latestDate time.Time
	if len(clauseRecommendations) > 0 {
		latestDate = clauseRecommendations[0].CreatedAt
	}
	
	result := &models.ContractRecommendations{
		ContractID:            contractID,
		OverallRiskLevel:      maxRiskLevel,
		OverallRiskScore:      overallRiskScore,
		TotalClauses:          totalClauses,
		ClauseRecommendations: clauseRecommendations,
		OverallRecommendations: overallRecommendations,
		KeyRisks:              keyRisks,
		CreatedAt:             latestDate,
	}
	
	return result, nil
}

// GetAllRecommendations retrieves all AI recommendations
func (r *AIRepository) GetAllRecommendations() ([]models.ClauseRiskAnalysis, error) {
	query := `
		SELECT 
			id, clause_id, risk_level, risk_score, analysis_summary,
			identified_risks, recommendations, legal_implications, 
			compliance_notes, confidence_score, model_version, created_at
		FROM clause_risk_analyses 
		ORDER BY created_at DESC
		LIMIT 20
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query recommendations: %w", err)
	}
	defer rows.Close()

	var recommendations []models.ClauseRiskAnalysis

	for rows.Next() {
		var analysis models.ClauseRiskAnalysis
		var identifiedRisksJSON, recommendationsJSON string

		err := rows.Scan(
			&analysis.ID,
			&analysis.ClauseID,
			&analysis.RiskLevel,
			&analysis.RiskScore,
			&analysis.AnalysisSummary,
			&identifiedRisksJSON,
			&recommendationsJSON,
			&analysis.LegalImplications,
			&analysis.ComplianceNotes,
			&analysis.ConfidenceScore,
			&analysis.ModelVersion,
			&analysis.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan analysis: %w", err)
		}

		// Parse JSON arrays
		if err := json.Unmarshal([]byte(identifiedRisksJSON), &analysis.IdentifiedRisks); err != nil {
			analysis.IdentifiedRisks = []string{}
		}
		if err := json.Unmarshal([]byte(recommendationsJSON), &analysis.Recommendations); err != nil {
			analysis.Recommendations = []string{}
		}

		recommendations = append(recommendations, analysis)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating analyses: %w", err)
	}

	return recommendations, nil
}

// removeDuplicates removes duplicate strings from a slice
func removeDuplicates(slice []string) []string {
	keys := make(map[string]bool)
	var result []string
	
	for _, item := range slice {
		if !keys[item] {
			keys[item] = true
			result = append(result, item)
		}
	}
	
	return result
}
