package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"contrack-be/internal/models"
)

// SaveContractAnalysis saves AI analysis results for a contract
func (r *ContractRepository) SaveContractAnalysis(contractID int, analysisResult *models.ContractAnalysisResult, reviewedBy string) error {
	// Update contract status to PENDING_LEGAL_REVIEW
	_, err := r.db.Exec(`
		UPDATE contracts 
		SET status = $1, updated_at = NOW() 
		WHERE id = $2`,
		models.StatusPendingLegalReview, contractID)
	if err != nil {
		return fmt.Errorf("failed to update contract status: %w", err)
	}

	// Save each clause analysis
	for _, clauseAnalysis := range analysisResult.ClauseAnalyses {
		// Convert arrays to JSON
		identifiedRisksJSON, _ := json.Marshal(clauseAnalysis.IdentifiedRisks)
		recommendationsJSON, _ := json.Marshal(clauseAnalysis.Recommendations)

		_, err := r.db.Exec(`
			INSERT INTO clause_risk_analyses (
				clause_id, risk_level, risk_score, analysis_summary,
				identified_risks, recommendations, legal_implications,
				compliance_notes, confidence_score, model_version, created_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())`,
			clauseAnalysis.ClauseID,
			clauseAnalysis.RiskLevel,
			clauseAnalysis.RiskScore,
			clauseAnalysis.AnalysisSummary,
			string(identifiedRisksJSON),
			string(recommendationsJSON),
			clauseAnalysis.LegalImplications,
			clauseAnalysis.ComplianceNotes,
			clauseAnalysis.ConfidenceScore,
			clauseAnalysis.ModelVersion,
		)
		if err != nil {
			return fmt.Errorf("failed to save clause analysis: %w", err)
		}
	}

	// Record status change in history
	_, err = r.db.Exec(`
		INSERT INTO contract_status_history (contract_id, from_status, to_status, changed_by, change_reason, comments)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		contractID, "DRAFT", models.StatusPendingLegalReview, reviewedBy, "AI Analysis Completed", "Contract analyzed by AI and ready for legal review")
	if err != nil {
		return fmt.Errorf("failed to record status change: %w", err)
	}

	return nil
}

// ProcessLegalReview processes legal review decision
func (r *ContractRepository) ProcessLegalReview(contractID int, decision string, notes string, rejectedReason string, reviewedBy string) (*models.LegalReviewResponse, error) {
	var newStatus models.ContractStatus
	var message string

	if decision == "approve" {
		newStatus = models.StatusPendingSignature
		message = "Contract approved by legal team and ready for signature"
	} else {
		newStatus = models.StatusRejected
		message = "Contract rejected by legal team"
	}

	// Update contract status
	_, err := r.db.Exec(`
		UPDATE contracts 
		SET status = $1, updated_at = NOW() 
		WHERE id = $2`,
		newStatus, contractID)
	if err != nil {
		return nil, fmt.Errorf("failed to update contract status: %w", err)
	}

	// Record status change in history
	changeReason := "Legal Review - Approved"
	if decision == "reject" {
		changeReason = "Legal Review - Rejected"
	}

	_, err = r.db.Exec(`
		INSERT INTO contract_status_history (contract_id, from_status, to_status, changed_by, change_reason, comments)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		contractID, models.StatusPendingLegalReview, newStatus, reviewedBy, changeReason, notes)
	if err != nil {
		return nil, fmt.Errorf("failed to record status change: %w", err)
	}

	response := &models.LegalReviewResponse{
		ContractID:     contractID,
		Decision:       decision,
		Notes:          notes,
		RejectedReason: rejectedReason,
		ReviewedBy:     reviewedBy,
		ReviewedAt:     time.Now(),
		NewStatus:      string(newStatus),
		Message:        message,
	}

	return response, nil
}
