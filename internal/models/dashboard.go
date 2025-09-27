package models

import "time"

// SimpleContractList represents a simple contract for dashboard table
type SimpleContractList struct {
	ID             int            `json:"id" db:"id"`
	ProjectName    string         `json:"project_name" db:"project_name"`
	ContractType   string         `json:"contract_type" db:"contract_type"`
	Status         ContractStatus `json:"status" db:"status"`
	StatusDisplay  string         `json:"status_display"`
	TotalValue     float64        `json:"total_value" db:"total_value"`
	SigningDate    *time.Time     `json:"signing_date" db:"signing_date"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
}

// StatusCount represents count for each contract status
type StatusCount struct {
	Status        ContractStatus `json:"status" db:"status"`
	StatusDisplay string         `json:"status_display"`
	Count         int            `json:"count" db:"count"`
}

// ProjectValueDistribution represents total value per project for bar chart
type ProjectValueDistribution struct {
	ProjectName string  `json:"project_name" db:"project_name"`
	TotalValue  float64 `json:"total_value" db:"total_value"`
}

// ContractTypeDistribution represents contract count by type for donut/pie chart
type ContractTypeDistribution struct {
	ContractType string `json:"contract_type" db:"contract_type"`
	Count        int    `json:"count" db:"count"`
}

// DashboardVisualizationData represents comprehensive data for dashboard visualizations
type DashboardVisualizationData struct {
	StatusCounts           []StatusCount                `json:"status_counts"`
	ProjectValueDist      []ProjectValueDistribution   `json:"project_value_distribution"`
	ContractTypeDist      []ContractTypeDistribution   `json:"contract_type_distribution"`
	TotalContracts        int                          `json:"total_contracts"`
	TotalValue            float64                      `json:"total_value"`
}

// GetStatusDisplay returns a human-readable string for the contract status
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