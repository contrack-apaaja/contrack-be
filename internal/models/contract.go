package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// ContractStatus represents the possible contract statuses
type ContractStatus string

const (
	StatusDraft               ContractStatus = "DRAFT"
	StatusPendingLegalReview  ContractStatus = "PENDING_LEGAL_REVIEW"
	StatusPendingSignature    ContractStatus = "PENDING_SIGNATURE"
	StatusActive              ContractStatus = "ACTIVE"
	StatusExpired             ContractStatus = "EXPIRED"
	StatusTerminated          ContractStatus = "TERMINATED"
)

// IsValidStatus checks if the status is valid
func (s ContractStatus) IsValidStatus() bool {
	switch s {
	case StatusDraft, StatusPendingLegalReview, StatusPendingSignature, StatusActive, StatusExpired, StatusTerminated:
		return true
	}
	return false
}

// CanTransitionTo checks if status transition is allowed
func (s ContractStatus) CanTransitionTo(newStatus ContractStatus) bool {
	transitions := map[ContractStatus][]ContractStatus{
		StatusDraft:               {StatusPendingLegalReview},
		StatusPendingLegalReview:  {StatusDraft, StatusPendingSignature},
		StatusPendingSignature:    {StatusPendingLegalReview, StatusActive},
		StatusActive:              {StatusExpired, StatusTerminated},
		StatusExpired:             {StatusTerminated},
		StatusTerminated:          {}, // Terminal status
	}
	
	allowed, exists := transitions[s]
	if !exists {
		return false
	}
	
	for _, allowedStatus := range allowed {
		if allowedStatus == newStatus {
			return true
		}
	}
	return false
}

// StakeholderType represents the type of stakeholder
type StakeholderType string

const (
	StakeholderIndividual StakeholderType = "INDIVIDUAL"
	StakeholderCompany    StakeholderType = "COMPANY"
	StakeholderGovernment StakeholderType = "GOVERNMENT"
	StakeholderNGO        StakeholderType = "NGO"
	StakeholderOther      StakeholderType = "OTHER"
)

// JSONB represents a PostgreSQL JSONB field
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j)
	case string:
		return json.Unmarshal([]byte(v), j)
	default:
		return errors.New("cannot scan into JSONB")
	}
}

// Contract represents a contract
type Contract struct {
	ID                int            `json:"id" db:"id"`
	BaseID            string         `json:"base_id" db:"base_id"`
	VersionNumber     int            `json:"version_number" db:"version_number"`
	ProjectName       string         `json:"project_name" db:"project_name"`
	PackageName       *string        `json:"package_name" db:"package_name"`
	ContractNumber    string         `json:"contract_number" db:"contract_number"`
	ExternalReference *string        `json:"external_reference" db:"external_reference"`
	ContractType      string         `json:"contract_type" db:"contract_type"`
	SigningPlace      *string        `json:"signing_place" db:"signing_place"`
	SigningDate       *time.Time     `json:"signing_date" db:"signing_date"`
	TotalValue        float64        `json:"total_value" db:"total_value"`
	FundingSource     *string        `json:"funding_source" db:"funding_source"`
	Status            ContractStatus `json:"status" db:"status"`
	CreatedBy         string         `json:"created_by" db:"created_by"`
	CreatedAt         time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt         *time.Time     `json:"deleted_at,omitempty" db:"deleted_at"`
	IsDeleted         bool           `json:"is_deleted" db:"is_deleted"`

	// Related entities (loaded separately)
	Stakeholders []ContractStakeholder `json:"stakeholders,omitempty"`
	Clauses      []ContractClause      `json:"clauses,omitempty"`
}

// Stakeholder represents a stakeholder entity
type Stakeholder struct {
	ID        int             `json:"id" db:"id"`
	LegalName string          `json:"legal_name" db:"legal_name"`
	Address   *string         `json:"address" db:"address"`
	Type      StakeholderType `json:"type" db:"type"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
	DeletedAt *time.Time      `json:"deleted_at,omitempty" db:"deleted_at"`
	IsDeleted bool            `json:"is_deleted" db:"is_deleted"`
}

// ContractStakeholder represents the relationship between contracts and stakeholders
type ContractStakeholder struct {
	ID                    int       `json:"id" db:"id"`
	ContractID            int       `json:"contract_id" db:"contract_id"`
	StakeholderID         int       `json:"stakeholder_id" db:"stakeholder_id"`
	RoleInContract        string    `json:"role_in_contract" db:"role_in_contract"`
	RepresentativeName    *string   `json:"representative_name" db:"representative_name"`
	RepresentativeTitle   *string   `json:"representative_title" db:"representative_title"`
	OtherDetails          JSONB     `json:"other_details" db:"other_details"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`

	// Related entity
	Stakeholder *Stakeholder `json:"stakeholder,omitempty"`
}

// ContractClause represents a clause instance within a contract
type ContractClause struct {
	ID                int            `json:"id" db:"id"`
	ContractID        int            `json:"contract_id" db:"contract_id"`
	ClauseTemplateID  int            `json:"clause_template_id" db:"clause_template_id"`
	DisplayOrder      int            `json:"display_order" db:"display_order"`
	CustomContent     *string        `json:"custom_content" db:"custom_content"`
	CreatedAt         time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at" db:"updated_at"`

	// Related entity
	ClauseTemplate *ClauseTemplate `json:"clause_template,omitempty"`
}

// ContractStatusHistory tracks status changes
type ContractStatusHistory struct {
	ID           int            `json:"id" db:"id"`
	ContractID   int            `json:"contract_id" db:"contract_id"`
	FromStatus   *ContractStatus `json:"from_status" db:"from_status"`
	ToStatus     ContractStatus `json:"to_status" db:"to_status"`
	ChangedBy    string         `json:"changed_by" db:"changed_by"`
	ChangeReason *string        `json:"change_reason" db:"change_reason"`
	Comments     *string        `json:"comments" db:"comments"`
	ChangedAt    time.Time      `json:"changed_at" db:"changed_at"`
}

// Request/Response Models

// ContractCreateRequest represents the request body for creating a contract
type ContractCreateRequest struct {
	ProjectName       string    `json:"project_name" binding:"required,min=3,max=255"`
	PackageName       *string   `json:"package_name" binding:"omitempty,max=255"`
	ExternalReference *string   `json:"external_reference" binding:"omitempty,max=100"`
	ContractType      string    `json:"contract_type" binding:"required,min=3,max=100"`
	SigningPlace      *string   `json:"signing_place" binding:"omitempty,max=255"`
	SigningDate       *string   `json:"signing_date" binding:"omitempty"` // ISO date format
	TotalValue        float64   `json:"total_value" binding:"required,gt=0"`
	FundingSource     *string   `json:"funding_source" binding:"omitempty,max=255"`
	Stakeholders      []ContractStakeholderCreateRequest `json:"stakeholders"`
	ClauseTemplateIDs []int     `json:"clause_template_ids"`
}

// ContractUpdateRequest represents the request body for updating a contract
type ContractUpdateRequest struct {
	ProjectName       *string   `json:"project_name,omitempty" binding:"omitempty,min=3,max=255"`
	PackageName       *string   `json:"package_name,omitempty" binding:"omitempty,max=255"`
	ExternalReference *string   `json:"external_reference,omitempty" binding:"omitempty,max=100"`
	ContractType      *string   `json:"contract_type,omitempty" binding:"omitempty,min=3,max=100"`
	SigningPlace      *string   `json:"signing_place,omitempty" binding:"omitempty,max=255"`
	SigningDate       *string   `json:"signing_date,omitempty"` // ISO date format
	TotalValue        *float64  `json:"total_value,omitempty" binding:"omitempty,gt=0"`
	FundingSource     *string   `json:"funding_source,omitempty" binding:"omitempty,max=255"`
}

// ContractStatusChangeRequest represents the request body for changing contract status
type ContractStatusChangeRequest struct {
	Status       ContractStatus `json:"status" binding:"required"`
	ChangeReason *string        `json:"change_reason" binding:"omitempty,max=500"`
	Comments     *string        `json:"comments" binding:"omitempty,max=1000"`
}

// ContractSearchRequest represents the request parameters for searching contracts
type ContractSearchRequest struct {
	Query         string          `form:"q" binding:"omitempty,min=2"`
	Status        *ContractStatus `form:"status" binding:"omitempty"`
	ContractType  string          `form:"contract_type" binding:"omitempty"`
	FundingSource string          `form:"funding_source" binding:"omitempty"`
	SigningDateFrom *string       `form:"signing_date_from" binding:"omitempty"`
	SigningDateTo   *string       `form:"signing_date_to" binding:"omitempty"`
	ValueFrom     *float64        `form:"value_from" binding:"omitempty,gte=0"`
	ValueTo       *float64        `form:"value_to" binding:"omitempty,gte=0"`
	CreatedBy     string          `form:"created_by" binding:"omitempty"`
	Page          int             `form:"page" binding:"omitempty,min=1"`
	Limit         int             `form:"limit" binding:"omitempty,min=1,max=100"`
	SortBy        string          `form:"sort_by" binding:"omitempty,oneof=id project_name contract_number contract_type total_value signing_date status created_at updated_at"`
	SortDir       string          `form:"sort_dir" binding:"omitempty,oneof=asc desc"`
}

// ContractListResponse represents the paginated response for contract list
type ContractListResponse struct {
	Contracts []Contract `json:"contracts"`
	Total     int        `json:"total"`
	Page      int        `json:"page"`
	Limit     int        `json:"limit"`
	Pages     int        `json:"pages"`
}

// StakeholderCreateRequest represents the request body for creating a stakeholder
type StakeholderCreateRequest struct {
	LegalName string          `json:"legal_name" binding:"required,min=2,max=255"`
	Address   *string         `json:"address" binding:"omitempty"`
	Type      StakeholderType `json:"type" binding:"required"`
}

// StakeholderUpdateRequest represents the request body for updating a stakeholder
type StakeholderUpdateRequest struct {
	LegalName *string          `json:"legal_name,omitempty" binding:"omitempty,min=2,max=255"`
	Address   *string          `json:"address,omitempty"`
	Type      *StakeholderType `json:"type,omitempty"`
}

// ContractStakeholderCreateRequest represents the request for adding stakeholders to contracts
type ContractStakeholderCreateRequest struct {
	StakeholderID       int     `json:"stakeholder_id" binding:"required"`
	RoleInContract      string  `json:"role_in_contract" binding:"required,min=2,max=100"`
	RepresentativeName  *string `json:"representative_name" binding:"omitempty,max=255"`
	RepresentativeTitle *string `json:"representative_title" binding:"omitempty,max=255"`
	OtherDetails        JSONB   `json:"other_details"`
}

// ContractClauseCreateRequest represents the request for adding clauses to contracts
type ContractClauseCreateRequest struct {
	ClauseTemplateID int     `json:"clause_template_id" binding:"required"`
	DisplayOrder     int     `json:"display_order" binding:"required,min=1"`
	CustomContent    *string `json:"custom_content"`
}

// ContractClauseUpdateRequest represents the request for updating contract clauses
type ContractClauseUpdateRequest struct {
	DisplayOrder  *int    `json:"display_order,omitempty" binding:"omitempty,min=1"`
	CustomContent *string `json:"custom_content,omitempty"`
}

// ContractApprovalRequest represents the request for one-click contract approval
type ContractApprovalRequest struct {
	ContractID int `json:"contract_id" binding:"required"`
}

// ContractApprovalResponse represents the response for contract approval
type ContractApprovalResponse struct {
	ContractID        int            `json:"contract_id"`
	ContractName      string         `json:"contract_name"`
	TotalValue        float64        `json:"total_value"`
	RiskLevel         string         `json:"risk_level"`
	RiskScore         float64        `json:"risk_score"`
	ApprovalStatus    string         `json:"approval_status"` // "ACTIVE", "REVIEW_REQUIRED"
	ApprovalMessage   string         `json:"approval_message"`
	RequiresReview    bool           `json:"requires_review"`
	ReviewReasons     []string       `json:"review_reasons,omitempty"`
	KeyRisks          []string       `json:"key_risks,omitempty"`
	Recommendations   []string       `json:"recommendations,omitempty"`
	ApprovedAt        *time.Time     `json:"approved_at,omitempty"`
	ApprovedBy        string         `json:"approved_by,omitempty"`
	NextSteps         []string       `json:"next_steps,omitempty"`
}

// ContractApprovalCriteria represents the criteria for automatic approval
type ContractApprovalCriteria struct {
	MaxRiskScore      float64 `json:"max_risk_score"`      // Maximum risk score for auto-approval
	MaxValue          float64 `json:"max_value"`            // Maximum contract value for auto-approval
	MinConfidenceScore float64 `json:"min_confidence_score"` // Minimum AI confidence score
}