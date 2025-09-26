package models

import (
	"time"
)

// ClauseTemplate represents a reusable clause template
type ClauseTemplate struct {
	ID        int       `json:"id" db:"id"`
	ClauseCode string   `json:"clause_code" db:"clause_code" binding:"required,min=3,max=50"`
	Title     string    `json:"title" db:"title" binding:"required,min=5,max=255"`
	Type      string    `json:"type" db:"type" binding:"required,min=3,max=100"`
	Content   string    `json:"content" db:"content" binding:"required,min=10"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ClauseTemplateCreateRequest represents the request body for creating a clause template
type ClauseTemplateCreateRequest struct {
	ClauseCode string `json:"clause_code" binding:"required,min=3,max=50"`
	Title      string `json:"title" binding:"required,min=5,max=255"`
	Type       string `json:"type" binding:"required,min=3,max=100"`
	Content    string `json:"content" binding:"required,min=10"`
	IsActive   *bool  `json:"is_active"` // Pointer to allow explicit false
}

// ClauseTemplateUpdateRequest represents the request body for updating a clause template
type ClauseTemplateUpdateRequest struct {
	ClauseCode *string `json:"clause_code,omitempty" binding:"omitempty,min=3,max=50"`
	Title      *string `json:"title,omitempty" binding:"omitempty,min=5,max=255"`
	Type       *string `json:"type,omitempty" binding:"omitempty,min=3,max=100"`
	Content    *string `json:"content,omitempty" binding:"omitempty,min=10"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

// ClauseTemplateSearchRequest represents the request parameters for searching clause templates
type ClauseTemplateSearchRequest struct {
	Query    string `form:"q" binding:"omitempty,min=2"`        // Search query
	Type     string `form:"type" binding:"omitempty"`           // Filter by type
	IsActive *bool  `form:"is_active" binding:"omitempty"`      // Filter by active status
	Page     int    `form:"page" binding:"omitempty,min=1"`     // Page number (default: 1)
	Limit    int    `form:"limit" binding:"omitempty,min=1,max=100"` // Items per page (default: 10)
	SortBy   string `form:"sort_by" binding:"omitempty,oneof=id title type created_at updated_at"` // Sort field
	SortDir  string `form:"sort_dir" binding:"omitempty,oneof=asc desc"` // Sort direction
}

// ClauseTemplateListResponse represents the paginated response for clause template list
type ClauseTemplateListResponse struct {
	ClauseTemplates []ClauseTemplate `json:"clause_templates"`
	Pagination      PaginationInfo   `json:"pagination"`
}

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

// ToCreateRequest converts ClauseTemplateCreateRequest to ClauseTemplate
func (req *ClauseTemplateCreateRequest) ToClauseTemplate() *ClauseTemplate {
	isActive := true // Default to active
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	return &ClauseTemplate{
		ClauseCode: req.ClauseCode,
		Title:      req.Title,
		Type:       req.Type,
		Content:    req.Content,
		IsActive:   isActive,
	}
}

// ApplyUpdates applies update request to existing clause template
func (clause *ClauseTemplate) ApplyUpdates(req *ClauseTemplateUpdateRequest) {
	if req.ClauseCode != nil {
		clause.ClauseCode = *req.ClauseCode
	}
	if req.Title != nil {
		clause.Title = *req.Title
	}
	if req.Type != nil {
		clause.Type = *req.Type
	}
	if req.Content != nil {
		clause.Content = *req.Content
	}
	if req.IsActive != nil {
		clause.IsActive = *req.IsActive
	}
	clause.UpdatedAt = time.Now()
}

// GetDefaultSearchParams returns default search parameters
func GetDefaultSearchParams() ClauseTemplateSearchRequest {
	return ClauseTemplateSearchRequest{
		Page:    1,
		Limit:   10,
		SortBy:  "created_at",
		SortDir: "desc",
	}
}
