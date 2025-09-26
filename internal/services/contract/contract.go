package contract

import (
	"fmt"
	"time"

	"contrack-be/internal/models"
	"contrack-be/internal/repository"
)

type Service struct {
	contractRepo       *repository.ContractRepository
	stakeholderRepo    *repository.StakeholderRepository
	contractClauseRepo *repository.ContractClauseRepository
	clauseTemplateRepo *repository.ClauseTemplateRepository
}

func NewService() *Service {
	return &Service{
		contractRepo:       repository.NewContractRepository(),
		stakeholderRepo:    repository.NewStakeholderRepository(),
		contractClauseRepo: repository.NewContractClauseRepository(),
		clauseTemplateRepo: repository.NewClauseTemplateRepository(),
	}
}

// CreateContract creates a new contract with stakeholders and clauses
func (s *Service) CreateContract(req *models.ContractCreateRequest, createdBy string) (*models.Contract, error) {
	// Validate signing date
	if req.SigningDate != nil {
		signingDate, err := time.Parse("2006-01-02", *req.SigningDate)
		if err != nil {
			return nil, fmt.Errorf("invalid signing date format: %w", err)
		}
		if signingDate.After(time.Now()) {
			return nil, fmt.Errorf("signing date cannot be in the future")
		}
	}

	// Parse signing date
	var signingDate *time.Time
	if req.SigningDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.SigningDate)
		if err != nil {
			return nil, fmt.Errorf("invalid signing date format: %w", err)
		}
		signingDate = &parsed
	}

	// Create contract
	contract := &models.Contract{
		ProjectName:       req.ProjectName,
		PackageName:       req.PackageName,
		ExternalReference: req.ExternalReference,
		ContractType:      req.ContractType,
		SigningPlace:      req.SigningPlace,
		SigningDate:       signingDate,
		TotalValue:        req.TotalValue,
		FundingSource:     req.FundingSource,
		Status:            models.StatusDraft,
		CreatedBy:         createdBy,
	}

	err := s.contractRepo.Create(contract)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract: %w", err)
	}

	// Add stakeholders if provided
	if len(req.Stakeholders) > 0 {
		for _, stakeholderReq := range req.Stakeholders {
			contractStakeholder := &models.ContractStakeholder{
				ContractID:          contract.ID,
				StakeholderID:       stakeholderReq.StakeholderID,
				RoleInContract:      stakeholderReq.RoleInContract,
				RepresentativeName:  stakeholderReq.RepresentativeName,
				RepresentativeTitle: stakeholderReq.RepresentativeTitle,
				OtherDetails:        stakeholderReq.OtherDetails,
			}
			
			err = s.stakeholderRepo.AddToContract(contractStakeholder)
			if err != nil {
				return nil, fmt.Errorf("failed to add stakeholder to contract: %w", err)
			}
		}
	}

	// Add clauses if provided
	if len(req.ClauseTemplateIDs) > 0 {
		err = s.contractClauseRepo.BulkAddClausesToContract(contract.ID, req.ClauseTemplateIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to add clauses to contract: %w", err)
		}
	}

	// Load related data
	return s.GetContractWithDetails(contract.ID)
}

// GetContract retrieves a contract by ID
func (s *Service) GetContract(id int) (*models.Contract, error) {
	return s.contractRepo.GetByID(id)
}

// GetContractWithDetails retrieves a contract with all related data
func (s *Service) GetContractWithDetails(id int) (*models.Contract, error) {
	contract, err := s.contractRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Load stakeholders
	stakeholders, err := s.stakeholderRepo.GetContractStakeholders(id)
	if err != nil {
		return nil, fmt.Errorf("failed to load contract stakeholders: %w", err)
	}
	contract.Stakeholders = stakeholders

	// Load clauses
	clauses, err := s.contractClauseRepo.GetContractClauses(id)
	if err != nil {
		return nil, fmt.Errorf("failed to load contract clauses: %w", err)
	}
	contract.Clauses = clauses

	return contract, nil
}

// UpdateContract updates a contract
func (s *Service) UpdateContract(id int, req *models.ContractUpdateRequest, userID string) error {
	// Check if contract can be edited
	canEdit, err := s.contractRepo.CanEdit(id)
	if err != nil {
		return err
	}
	if !canEdit {
		return fmt.Errorf("contract cannot be edited in its current status")
	}

	// Prepare updates
	updates := make(map[string]interface{})
	
	if req.ProjectName != nil {
		updates["project_name"] = *req.ProjectName
	}
	if req.PackageName != nil {
		updates["package_name"] = *req.PackageName
	}
	if req.ExternalReference != nil {
		updates["external_reference"] = *req.ExternalReference
	}
	if req.ContractType != nil {
		updates["contract_type"] = *req.ContractType
	}
	if req.SigningPlace != nil {
		updates["signing_place"] = *req.SigningPlace
	}
	if req.SigningDate != nil {
		signingDate, err := time.Parse("2006-01-02", *req.SigningDate)
		if err != nil {
			return fmt.Errorf("invalid signing date format: %w", err)
		}
		if signingDate.After(time.Now()) {
			return fmt.Errorf("signing date cannot be in the future")
		}
		updates["signing_date"] = signingDate
	}
	if req.TotalValue != nil {
		updates["total_value"] = *req.TotalValue
	}
	if req.FundingSource != nil {
		updates["funding_source"] = *req.FundingSource
	}

	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	return s.contractRepo.Update(id, updates)
}

// DeleteContract soft deletes a contract
func (s *Service) DeleteContract(id int, userID string) error {
	// Check if contract can be edited (only DRAFT contracts can be deleted)
	canEdit, err := s.contractRepo.CanEdit(id)
	if err != nil {
		return err
	}
	if !canEdit {
		return fmt.Errorf("only draft contracts can be deleted")
	}

	return s.contractRepo.Delete(id)
}

// SearchContracts searches contracts with filters
func (s *Service) SearchContracts(req *models.ContractSearchRequest) (*models.ContractListResponse, error) {
	contracts, total, err := s.contractRepo.List(req)
	if err != nil {
		return nil, err
	}

	// Calculate pagination info
	pages := (total + req.Limit - 1) / req.Limit
	if pages == 0 {
		pages = 1
	}

	return &models.ContractListResponse{
		Contracts: contracts,
		Total:     total,
		Page:      req.Page,
		Limit:     req.Limit,
		Pages:     pages,
	}, nil
}

// ChangeContractStatus changes the status of a contract
func (s *Service) ChangeContractStatus(id int, req *models.ContractStatusChangeRequest, userID string) error {
	// Validate status
	if !req.Status.IsValidStatus() {
		return fmt.Errorf("invalid status: %s", req.Status)
	}

	return s.contractRepo.UpdateStatus(id, req.Status, userID, req.ChangeReason, req.Comments)
}

// GetContractStatusHistory retrieves status change history for a contract
func (s *Service) GetContractStatusHistory(id int) ([]models.ContractStatusHistory, error) {
	return s.contractRepo.GetStatusHistory(id)
}

// CreateContractVersion creates a new version of an existing contract
func (s *Service) CreateContractVersion(baseID string, req *models.ContractCreateRequest, createdBy string) (*models.Contract, error) {
	// Validate signing date
	if req.SigningDate != nil {
		signingDate, err := time.Parse("2006-01-02", *req.SigningDate)
		if err != nil {
			return nil, fmt.Errorf("invalid signing date format: %w", err)
		}
		if signingDate.After(time.Now()) {
			return nil, fmt.Errorf("signing date cannot be in the future")
		}
	}

	// Parse signing date
	var signingDate *time.Time
	if req.SigningDate != nil {
		parsed, err := time.Parse("2006-01-02", *req.SigningDate)
		if err != nil {
			return nil, fmt.Errorf("invalid signing date format: %w", err)
		}
		signingDate = &parsed
	}

	// Create new version
	// TODO: Memperbaiki Atribut Contract
	fmt.Println(">>>>>>>>>> req (old contract): ", req)
	contract := &models.Contract{
		ProjectName:       req.ProjectName,
		PackageName:       req.PackageName,
		ExternalReference: req.ExternalReference,
		ContractType:      req.ContractType,
		SigningPlace:      req.SigningPlace,
		SigningDate:       signingDate,
		TotalValue:        req.TotalValue,
		FundingSource:     req.FundingSource,
		CreatedBy:         createdBy,
	}
	fmt.Println("<<<<<<<<<< contract: ", contract)

	err := s.contractRepo.CreateVersion(baseID, contract)
	if err != nil {
		return nil, fmt.Errorf("failed to create contract version: %w", err)
	}

	// Add stakeholders if provided
	if len(req.Stakeholders) > 0 {
		for _, stakeholderReq := range req.Stakeholders {
			contractStakeholder := &models.ContractStakeholder{
				ContractID:          contract.ID,
				StakeholderID:       stakeholderReq.StakeholderID,
				RoleInContract:      stakeholderReq.RoleInContract,
				RepresentativeName:  stakeholderReq.RepresentativeName,
				RepresentativeTitle: stakeholderReq.RepresentativeTitle,
				OtherDetails:        stakeholderReq.OtherDetails,
			}
			
			err = s.stakeholderRepo.AddToContract(contractStakeholder)
			if err != nil {
				return nil, fmt.Errorf("failed to add stakeholder to contract: %w", err)
			}
		}
	}

	// Add clauses if provided
	if len(req.ClauseTemplateIDs) > 0 {
		err = s.contractClauseRepo.BulkAddClausesToContract(contract.ID, req.ClauseTemplateIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to add clauses to contract: %w", err)
		}
	}

	// Load related data
	return s.GetContractWithDetails(contract.ID)
}

// GetContractVersions retrieves all versions of a contract
func (s *Service) GetContractVersions(baseID string) ([]models.Contract, error) {
	return s.contractRepo.GetVersions(baseID)
}

// AddStakeholderToContract adds a stakeholder to a contract
func (s *Service) AddStakeholderToContract(contractID int, req *models.ContractStakeholderCreateRequest) error {
	// Check if contract can be edited
	canEdit, err := s.contractRepo.CanEdit(contractID)
	if err != nil {
		return err
	}
	if !canEdit {
		return fmt.Errorf("contract cannot be edited in its current status")
	}

	contractStakeholder := &models.ContractStakeholder{
		ContractID:          contractID,
		StakeholderID:       req.StakeholderID,
		RoleInContract:      req.RoleInContract,
		RepresentativeName:  req.RepresentativeName,
		RepresentativeTitle: req.RepresentativeTitle,
		OtherDetails:        req.OtherDetails,
	}

	return s.stakeholderRepo.AddToContract(contractStakeholder)
}

// RemoveStakeholderFromContract removes a stakeholder from a contract
func (s *Service) RemoveStakeholderFromContract(contractID, stakeholderID int, role string) error {
	// Check if contract can be edited
	canEdit, err := s.contractRepo.CanEdit(contractID)
	if err != nil {
		return err
	}
	if !canEdit {
		return fmt.Errorf("contract cannot be edited in its current status")
	}

	return s.stakeholderRepo.RemoveFromContract(contractID, stakeholderID, role)
}

// AddClauseToContract adds a clause to a contract
func (s *Service) AddClauseToContract(contractID int, req *models.ContractClauseCreateRequest) error {
	// Check if contract can be edited
	canEdit, err := s.contractRepo.CanEdit(contractID)
	if err != nil {
		return err
	}
	if !canEdit {
		return fmt.Errorf("contract cannot be edited in its current status")
	}

	contractClause := &models.ContractClause{
		ContractID:       contractID,
		ClauseTemplateID: req.ClauseTemplateID,
		DisplayOrder:     req.DisplayOrder,
		CustomContent:    req.CustomContent,
	}

	return s.contractClauseRepo.AddClauseToContract(contractClause)
}

// RemoveClauseFromContract removes a clause from a contract
func (s *Service) RemoveClauseFromContract(contractID, clauseTemplateID int) error {
	// Check if contract can be edited
	canEdit, err := s.contractRepo.CanEdit(contractID)
	if err != nil {
		return err
	}
	if !canEdit {
		return fmt.Errorf("contract cannot be edited in its current status")
	}

	return s.contractClauseRepo.RemoveClauseFromContract(contractID, clauseTemplateID)
}

// UpdateContractClause updates a contract clause
func (s *Service) UpdateContractClause(contractID, clauseID int, req *models.ContractClauseUpdateRequest) error {
	// Check if contract can be edited
	canEdit, err := s.contractRepo.CanEdit(contractID)
	if err != nil {
		return err
	}
	if !canEdit {
		return fmt.Errorf("contract cannot be edited in its current status")
	}

	return s.contractClauseRepo.UpdateContractClause(clauseID, req.DisplayOrder, req.CustomContent)
}

// ReorderContractClauses reorders clauses in a contract
func (s *Service) ReorderContractClauses(contractID int, clauseOrders []struct {
	ClauseID     int `json:"clause_id"`
	DisplayOrder int `json:"display_order"`
}) error {
	// Check if contract can be edited
	canEdit, err := s.contractRepo.CanEdit(contractID)
	if err != nil {
		return err
	}
	if !canEdit {
		return fmt.Errorf("contract cannot be edited in its current status")
	}

	return s.contractClauseRepo.ReorderClauses(contractID, clauseOrders)
}

// GetContractsByStatus retrieves contracts by status for a user
func (s *Service) GetContractsByStatus(status models.ContractStatus, createdBy string) ([]models.Contract, error) {
	return s.contractRepo.GetContractsByStatus(status, createdBy)
}

// ValidateContractAccess checks if a user can access a contract
func (s *Service) ValidateContractAccess(contractID int, userID string) error {
	contract, err := s.contractRepo.GetByID(contractID)
	if err != nil {
		return err
	}

	// For now, users can only access contracts they created
	// In future phases, this will be extended for different roles
	if contract.CreatedBy != userID {
		return fmt.Errorf("access denied: you can only access contracts you created")
	}

	return nil
}

// GetContractForPDF checks if a contract can be downloaded as PDF
func (s *Service) GetContractForPDF(contractID int, userID string) (*models.Contract, error) {
	// Validate access
	err := s.ValidateContractAccess(contractID, userID)
	if err != nil {
		return nil, err
	}

	contract, err := s.GetContractWithDetails(contractID)
	if err != nil {
		return nil, err
	}

	// Check if contract can be downloaded
	if contract.Status != models.StatusPendingSignature && contract.Status != models.StatusActive {
		return nil, fmt.Errorf("contract can only be downloaded when status is PENDING_SIGNATURE or ACTIVE")
	}

	return contract, nil
}

// ValidateStatusTransition validates if a status transition is allowed
func (s *Service) ValidateStatusTransition(contractID int, newStatus models.ContractStatus) error {
	contract, err := s.contractRepo.GetByID(contractID)
	if err != nil {
		return err
	}

	if !contract.Status.CanTransitionTo(newStatus) {
		return fmt.Errorf("invalid status transition from %s to %s", contract.Status, newStatus)
	}

	return nil
}

// GetContractStats returns statistics about contracts for a user
func (s *Service) GetContractStats(userID string) (map[string]int, error) {
	stats := make(map[string]int)
	
	statuses := []models.ContractStatus{
		models.StatusDraft,
		models.StatusPendingLegalReview,
		models.StatusPendingSignature,
		models.StatusActive,
		models.StatusExpired,
		models.StatusTerminated,
	}

	for _, status := range statuses {
		contracts, err := s.contractRepo.GetContractsByStatus(status, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get contracts for status %s: %w", status, err)
		}
		stats[string(status)] = len(contracts)
	}

	return stats, nil
}