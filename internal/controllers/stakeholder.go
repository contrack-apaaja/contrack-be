package controllers

import (
	"net/http"
	"strconv"

	"contrack-be/internal/models"
	stakeholderService "contrack-be/internal/services/stakeholder"
	"contrack-be/internal/utils"

	"github.com/gin-gonic/gin"
)

type StakeholderController struct {
	stakeholderService *stakeholderService.Service
}

func NewStakeholderController() *StakeholderController {
	return &StakeholderController{
		stakeholderService: stakeholderService.NewService(),
	}
}

// CreateStakeholder creates a new stakeholder
func (sc *StakeholderController) CreateStakeholder(c *gin.Context) {
	var req models.StakeholderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "VALIDATION_ERROR", err.Error())
		return
	}

	stakeholder, err := sc.stakeholderService.CreateStakeholder(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create stakeholder", "CREATE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Stakeholder created successfully", stakeholder)
}

// GetStakeholder retrieves a stakeholder by ID
func (sc *StakeholderController) GetStakeholder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid stakeholder ID", "VALIDATION_ERROR", err.Error())
		return
	}

	stakeholder, err := sc.stakeholderService.GetStakeholder(id)
	if err != nil {
		if err.Error() == "stakeholder not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Stakeholder not found", "NOT_FOUND", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get stakeholder", "FETCH_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Stakeholder retrieved successfully", stakeholder)
}

// UpdateStakeholder updates a stakeholder
func (sc *StakeholderController) UpdateStakeholder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid stakeholder ID", "VALIDATION_ERROR", err.Error())
		return
	}

	var req models.StakeholderUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "VALIDATION_ERROR", err.Error())
		return
	}

	err = sc.stakeholderService.UpdateStakeholder(id, &req)
	if err != nil {
		if err.Error() == "stakeholder not found or already deleted" {
			utils.ErrorResponse(c, http.StatusNotFound, "Stakeholder not found", "NOT_FOUND", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update stakeholder", "UPDATE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Stakeholder updated successfully", nil)
}

// DeleteStakeholder soft deletes a stakeholder
func (sc *StakeholderController) DeleteStakeholder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid stakeholder ID", "VALIDATION_ERROR", err.Error())
		return
	}

	err = sc.stakeholderService.DeleteStakeholder(id)
	if err != nil {
		if err.Error() == "stakeholder not found or already deleted" {
			utils.ErrorResponse(c, http.StatusNotFound, "Stakeholder not found", "NOT_FOUND", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete stakeholder", "DELETE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Stakeholder deleted successfully", nil)
}

// ListStakeholders retrieves stakeholders with pagination and filtering
func (sc *StakeholderController) ListStakeholders(c *gin.Context) {
	search := c.Query("search")
	typeStr := c.Query("type")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	var stakeholderType *models.StakeholderType
	if typeStr != "" {
		sType := models.StakeholderType(typeStr)
		stakeholderType = &sType
	}

	stakeholders, total, err := sc.stakeholderService.ListStakeholders(search, stakeholderType, page, limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to list stakeholders", "LIST_ERROR", err.Error())
		return
	}

	// Calculate pagination info
	pages := (total + limit - 1) / limit
	if pages == 0 {
		pages = 1
	}

	response := map[string]interface{}{
		"stakeholders": stakeholders,
		"total":        total,
		"page":         page,
		"limit":        limit,
		"pages":        pages,
	}

	utils.SuccessResponse(c, http.StatusOK, "Stakeholders retrieved successfully", response)
}

// GetStakeholderTypes returns all available stakeholder types
func (sc *StakeholderController) GetStakeholderTypes(c *gin.Context) {
	types := sc.stakeholderService.GetStakeholderTypes()
	utils.SuccessResponse(c, http.StatusOK, "Stakeholder types retrieved successfully", types)
}