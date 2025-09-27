package clause

import (
	"fmt"
	"strings"

	"contrack-be/internal/models"
	"contrack-be/internal/repository"
)

type Service struct {
	repo *repository.ClauseTemplateRepository
}

// NewClauseService creates a new clause service
func NewClauseService() *Service {
	return &Service{
		repo: repository.NewClauseTemplateRepository(),
	}
}

// CreateClauseTemplate creates a new clause template
func (s *Service) CreateClauseTemplate(req *models.ClauseTemplateCreateRequest) (*models.ClauseTemplate, error) {
	// Validate clause code uniqueness
	existing, err := s.repo.GetByClauseCode(req.ClauseCode)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("clause template with code '%s' already exists", req.ClauseCode)
	}

	// Convert request to clause template
	clause := req.ToClauseTemplate()

	// Create the clause template
	createdClause, err := s.repo.Create(clause)
	if err != nil {
		return nil, fmt.Errorf("failed to create clause template: %w", err)
	}

	return createdClause, nil
}

// GetClauseTemplateByID retrieves a clause template by ID
func (s *Service) GetClauseTemplateByID(id int) (*models.ClauseTemplate, error) {
	clause, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return clause, nil
}

// GetClauseTemplateByCode retrieves a clause template by clause code
func (s *Service) GetClauseTemplateByCode(clauseCode string) (*models.ClauseTemplate, error) {
	clause, err := s.repo.GetByClauseCode(clauseCode)
	if err != nil {
		return nil, err
	}

	return clause, nil
}

// UpdateClauseTemplate updates an existing clause template
func (s *Service) UpdateClauseTemplate(id int, req *models.ClauseTemplateUpdateRequest) (*models.ClauseTemplate, error) {
	// Get existing clause template
	clause, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check clause code uniqueness if it's being updated
	if req.ClauseCode != nil && *req.ClauseCode != clause.ClauseCode {
		existing, err := s.repo.GetByClauseCode(*req.ClauseCode)
		if err == nil && existing != nil {
			return nil, fmt.Errorf("clause template with code '%s' already exists", *req.ClauseCode)
		}
	}

	// Apply updates
	clause.ApplyUpdates(req)

	// Update in database
	updatedClause, err := s.repo.Update(clause)
	if err != nil {
		return nil, fmt.Errorf("failed to update clause template: %w", err)
	}

	return updatedClause, nil
}

// DeleteClauseTemplate deletes a clause template by ID
func (s *Service) DeleteClauseTemplate(id int) error {
	// Check if clause template exists
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Delete the clause template
	err = s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete clause template: %w", err)
	}

	return nil
}

// ListClauseTemplates retrieves clause templates with pagination and search
func (s *Service) ListClauseTemplates(req *models.ClauseTemplateSearchRequest) (*models.ClauseTemplateListResponse, error) {
	// Set default values if not provided
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 50
	}
	if req.SortBy == "" {
		req.SortBy = "created_at"
	}
	if req.SortDir == "" {
		req.SortDir = "desc"
	}

	// Validate sort parameters
	validSortFields := []string{"id", "title", "type", "created_at", "updated_at"}
	if !contains(validSortFields, req.SortBy) {
		req.SortBy = "created_at"
	}

	req.SortDir = strings.ToLower(req.SortDir)
	if req.SortDir != "asc" && req.SortDir != "desc" {
		req.SortDir = "desc"
	}

	// Clean up search query
	if req.Query != "" {
		req.Query = strings.TrimSpace(req.Query)
	}

	// Get results from repository
	result, err := s.repo.List(*req)
	if err != nil {
		return nil, fmt.Errorf("failed to list clause templates: %w", err)
	}

	return result, nil
}

// GetClauseTypes retrieves all unique clause types from active templates
func (s *Service) GetClauseTypes() ([]string, error) {
	types, err := s.repo.GetActiveTypes()
	if err != nil {
		return nil, fmt.Errorf("failed to get clause types: %w", err)
	}

	return types, nil
}

// SearchClauseTemplates performs a search on clause templates
func (s *Service) SearchClauseTemplates(query string, limit int) (*models.ClauseTemplateListResponse, error) {
	searchReq := models.ClauseTemplateSearchRequest{
		Query:   query,
		Page:    1,
		Limit:   limit,
		SortBy:  "created_at",
		SortDir: "desc",
	}

	if limit == 0 {
		searchReq.Limit = 50
	}

	return s.ListClauseTemplates(&searchReq)
}

// ToggleClauseTemplateStatus toggles the active status of a clause template
func (s *Service) ToggleClauseTemplateStatus(id int) (*models.ClauseTemplate, error) {
	// Get existing clause template
	clause, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Toggle status
	clause.IsActive = !clause.IsActive

	// Update in database
	updatedClause, err := s.repo.Update(clause)
	if err != nil {
		return nil, fmt.Errorf("failed to toggle clause template status: %w", err)
	}

	return updatedClause, nil
}

// helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
