package models

import (
	"time"
)

// RiskLevel represents the risk level of a clause
type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "low"
	RiskLevelMedium RiskLevel = "medium"
	RiskLevelHigh   RiskLevel = "high"
	RiskLevelCritical RiskLevel = "critical"
)

// ClauseRiskAnalysis represents the AI analysis result for a clause
type ClauseRiskAnalysis struct {
	ID                int       `json:"id" db:"id"`
	ClauseID          int       `json:"clause_id" db:"clause_id"`
	RiskLevel         RiskLevel `json:"risk_level" db:"risk_level"`
	RiskScore         float64   `json:"risk_score" db:"risk_score"` // 0-100 scale
	AnalysisSummary   string    `json:"analysis_summary" db:"analysis_summary"`
	IdentifiedRisks   []string  `json:"identified_risks" db:"identified_risks"` // JSON array
	Recommendations   []string  `json:"recommendations" db:"recommendations"`   // JSON array
	LegalImplications string    `json:"legal_implications" db:"legal_implications"`
	ComplianceNotes   string    `json:"compliance_notes" db:"compliance_notes"`
	ConfidenceScore   float64   `json:"confidence_score" db:"confidence_score"` // 0-100 scale
	ModelVersion      string    `json:"model_version" db:"model_version"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// ContractRecommendations represents all recommendations for a contract
type ContractRecommendations struct {
	ContractID        int                    `json:"contract_id"`
	OverallRiskLevel  RiskLevel              `json:"overall_risk_level"`
	OverallRiskScore  float64                `json:"overall_risk_score"`
	TotalClauses      int                    `json:"total_clauses"`
	ClauseRecommendations []ClauseRecommendation `json:"clause_recommendations"`
	OverallRecommendations []string           `json:"overall_recommendations"`
	KeyRisks          []string               `json:"key_risks"`
	CreatedAt         time.Time              `json:"created_at"`
}

// ClauseRecommendation represents recommendations for a specific clause
type ClauseRecommendation struct {
	ClauseID          int       `json:"clause_id"`
	RiskLevel         RiskLevel `json:"risk_level"`
	RiskScore         float64   `json:"risk_score"`
	AnalysisSummary   string    `json:"analysis_summary"`
	Recommendations   []string  `json:"recommendations"`
	IdentifiedRisks   []string  `json:"identified_risks"`
	LegalImplications string    `json:"legal_implications"`
	ComplianceNotes   string    `json:"compliance_notes"`
	ConfidenceScore   float64   `json:"confidence_score"`
	CreatedAt         time.Time `json:"created_at"`
}

// ClauseRiskAnalysisRequest represents the request for AI analysis
type ClauseRiskAnalysisRequest struct {
	ClauseID int `json:"clause_id" binding:"required"`
}

// ClauseRiskAnalysisResponse represents the response for AI analysis
type ClauseRiskAnalysisResponse struct {
	Analysis ClauseRiskAnalysis `json:"analysis"`
	Clause   ClauseTemplate     `json:"clause"`
}

// ClauseRiskAnalysisListResponse represents paginated response for analysis list
type ClauseRiskAnalysisListResponse struct {
	Analyses   []ClauseRiskAnalysis `json:"analyses"`
	Pagination PaginationInfo       `json:"pagination"`
}

// ClauseRiskAnalysisSearchRequest represents search parameters for analysis
type ClauseRiskAnalysisSearchRequest struct {
	ClauseID      int        `form:"clause_id" binding:"omitempty"`
	RiskLevel     RiskLevel  `form:"risk_level" binding:"omitempty,oneof=low medium high critical"`
	MinRiskScore  float64    `form:"min_risk_score" binding:"omitempty,min=0,max=100"`
	MaxRiskScore  float64    `form:"max_risk_score" binding:"omitempty,min=0,max=100"`
	MinConfidence float64    `form:"min_confidence" binding:"omitempty,min=0,max=100"`
	Page          int        `form:"page" binding:"omitempty,min=1"`
	Limit         int        `form:"limit" binding:"omitempty,min=1,max=100"`
	SortBy        string     `form:"sort_by" binding:"omitempty,oneof=id clause_id risk_level risk_score created_at updated_at"`
	SortDir       string     `form:"sort_dir" binding:"omitempty,oneof=asc desc"`
}

// GetDefaultAnalysisSearchParams returns default search parameters for analysis
func GetDefaultAnalysisSearchParams() ClauseRiskAnalysisSearchRequest {
	return ClauseRiskAnalysisSearchRequest{
		Page:    1,
		Limit:   10,
		SortBy:  "created_at",
		SortDir: "desc",
	}
}

// IsValidRiskLevel checks if the risk level is valid
func IsValidRiskLevel(level RiskLevel) bool {
	switch level {
	case RiskLevelLow, RiskLevelMedium, RiskLevelHigh, RiskLevelCritical:
		return true
	default:
		return false
	}
}

// GetRiskLevelFromScore converts risk score to risk level
func GetRiskLevelFromScore(score float64) RiskLevel {
	switch {
	case score >= 0 && score < 25:
		return RiskLevelLow
	case score >= 25 && score < 50:
		return RiskLevelMedium
	case score >= 50 && score < 75:
		return RiskLevelHigh
	case score >= 75 && score <= 100:
		return RiskLevelCritical
	default:
		return RiskLevelLow
	}
}

// ContractAnalysisRequest represents the request for analyzing a contract
type ContractAnalysisRequest struct {
	ContractID        int   `json:"contract_id" binding:"required"`
	ClauseTemplateIDs []int `json:"clause_template_ids" binding:"required,min=1"`
}

// ContractAnalysisResult represents the AI analysis result for a contract
type ContractAnalysisResult struct {
	ContractID        int                    `json:"contract_id"`
	ClauseAnalyses    []ClauseRiskAnalysis   `json:"clause_analyses"`
	OverallRiskLevel  RiskLevel              `json:"overall_risk_level"`
	OverallRiskScore  float64                `json:"overall_risk_score"`
	ContractSummary   string                 `json:"contract_summary"`
	KeyRisks          []string               `json:"key_risks"`
	Recommendations   []string               `json:"recommendations"`
	CreatedAt         time.Time              `json:"created_at"`
}
