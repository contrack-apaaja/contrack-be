package stakeholder

import (
	"fmt"

	"contrack-be/internal/models"
	"contrack-be/internal/repository"
)

type Service struct {
	stakeholderRepo *repository.StakeholderRepository
}

func NewService() *Service {
	return &Service{
		stakeholderRepo: repository.NewStakeholderRepository(),
	}
}

// CreateStakeholder creates a new stakeholder
func (s *Service) CreateStakeholder(req *models.StakeholderCreateRequest) (*models.Stakeholder, error) {
	stakeholder := &models.Stakeholder{
		LegalName: req.LegalName,
		Address:   req.Address,
		Type:      req.Type,
	}

	err := s.stakeholderRepo.Create(stakeholder)
	if err != nil {
		return nil, fmt.Errorf("failed to create stakeholder: %w", err)
	}

	return stakeholder, nil
}

// GetStakeholder retrieves a stakeholder by ID
func (s *Service) GetStakeholder(id int) (*models.Stakeholder, error) {
	return s.stakeholderRepo.GetByID(id)
}

// UpdateStakeholder updates a stakeholder
func (s *Service) UpdateStakeholder(id int, req *models.StakeholderUpdateRequest) error {
	updates := make(map[string]interface{})
	
	if req.LegalName != nil {
		updates["legal_name"] = *req.LegalName
	}
	if req.Address != nil {
		updates["address"] = *req.Address
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}

	if len(updates) == 0 {
		return fmt.Errorf("no updates provided")
	}

	return s.stakeholderRepo.Update(id, updates)
}

// DeleteStakeholder soft deletes a stakeholder
func (s *Service) DeleteStakeholder(id int) error {
	return s.stakeholderRepo.Delete(id)
}

// ListStakeholders retrieves stakeholders with pagination and filtering
func (s *Service) ListStakeholders(search string, stakeholderType *models.StakeholderType, page, limit int) ([]models.Stakeholder, int, error) {
	return s.stakeholderRepo.List(search, stakeholderType, page, limit)
}

// GetStakeholderTypes returns all available stakeholder types
func (s *Service) GetStakeholderTypes() []models.StakeholderType {
	return []models.StakeholderType{
		models.StakeholderIndividual,
		models.StakeholderCompany,
		models.StakeholderGovernment,
		models.StakeholderNGO,
		models.StakeholderOther,
	}
}