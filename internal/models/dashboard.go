package models

import "time"

// DashboardContractSummary represents a contract summary for dashboard
type DashboardContractSummary struct {
	ID                int            `json:"id"`
	BaseID            string         `json:"base_id"`
	ProjectName       string         `json:"project_name"`
	ContractNumber    string         `json:"contract_number"`
	ContractType      string         `json:"contract_type"`
	Status            ContractStatus `json:"status"`
	TotalValue        float64        `json:"total_value"`
	SigningDate       *time.Time     `json:"signing_date"`
	CreatedBy         string         `json:"created_by"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	
	// Additional fields for dashboard
	StakeholderCount  int    `json:"stakeholder_count"`
	ClauseCount       int    `json:"clause_count"`
	DaysSinceCreated  int    `json:"days_since_created"`
	StatusDisplay     string `json:"status_display"`
}

// ContractStatusStats represents statistics for each contract status
type ContractStatusStats struct {
	Status        ContractStatus `json:"status"`
	StatusDisplay string         `json:"status_display"`
	Count         int            `json:"count"`
	Percentage    float64        `json:"percentage"`
	TotalValue    float64        `json:"total_value"`
}

// DashboardStats represents overall dashboard statistics
type DashboardStats struct {
	TotalContracts     int                  `json:"total_contracts"`
	TotalValue         float64              `json:"total_value"`
	StatusStats        []ContractStatusStats `json:"status_stats"`
	RecentContracts    []DashboardContractSummary `json:"recent_contracts"`
	ExpiringSoon       []DashboardContractSummary `json:"expiring_soon"`
	HighValueContracts []DashboardContractSummary `json:"high_value_contracts"`
}

// DashboardRequest represents request parameters for dashboard
type DashboardRequest struct {
	Page          int    `form:"page" binding:"omitempty,min=1"`
	Limit         int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Status        string `form:"status" binding:"omitempty"`
	ContractType  string `form:"contract_type" binding:"omitempty"`
	SortBy        string `form:"sort_by" binding:"omitempty,oneof=project_name contract_number total_value signing_date status created_at"`
	SortDir       string `form:"sort_dir" binding:"omitempty,oneof=asc desc"`
	IncludeStats  bool   `form:"include_stats" binding:"omitempty"`
}

// GetDefaultDashboardParams returns default dashboard parameters
func GetDefaultDashboardParams() DashboardRequest {
	return DashboardRequest{
		Page:         1,
		Limit:        10,
		SortBy:       "created_at",
		SortDir:      "desc",
		IncludeStats: true,
	}
}

// GetStatusDisplay returns a human-readable status display
func (s ContractStatus) GetStatusDisplay() string {
	switch s {
	case StatusDraft:
		return "Draft"
	case StatusPendingLegalReview:
		return "Pending Legal Review"
	case StatusPendingSignature:
		return "Pending Signature"
	case StatusActive:
		return "Active"
	case StatusExpired:
		return "Expired"
	case StatusTerminated:
		return "Terminated"
	default:
		return string(s)
	}
}
