package controllers

import (
	"net/http"
	"strconv"

	"contrack-be/internal/models"
	contractService "contrack-be/internal/services/contract"
	"contrack-be/internal/utils"

	"github.com/gin-gonic/gin"
)

type ContractController struct {
	contractService *contractService.Service
}

func NewContractController() *ContractController {
	return &ContractController{
		contractService: contractService.NewService(),
	}
}

// CreateContract creates a new contract
func (cc *ContractController) CreateContract(c *gin.Context) {
	var req models.ContractCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "VALIDATION_ERROR", err.Error())
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", "AUTH_ERROR", nil)
		return
	}

	contract, err := cc.contractService.CreateContract(&req, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create contract", "CREATE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Contract created successfully", contract)
}

// GetContract retrieves a contract by ID
func (cc *ContractController) GetContract(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid contract ID", "VALIDATION_ERROR", err.Error())
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", "AUTH_ERROR", nil)
		return
	}

	// Validate access
	err = cc.contractService.ValidateContractAccess(id, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "ACCESS_ERROR", err.Error())
		return
	}

	contract, err := cc.contractService.GetContractWithDetails(id)
	if err != nil {
		if err.Error() == "contract not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Contract not found", "NOT_FOUND", nil)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get contract", "FETCH_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Contract retrieved successfully", contract)
}

// UpdateContract updates a contract
func (cc *ContractController) UpdateContract(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid contract ID", "VALIDATION_ERROR", err.Error())
		return
	}

	var req models.ContractUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "VALIDATION_ERROR", err.Error())
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", "AUTH_ERROR", nil)
		return
	}

	// Validate access
	err = cc.contractService.ValidateContractAccess(id, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "ACCESS_ERROR", err.Error())
		return
	}

	err = cc.contractService.UpdateContract(id, &req, userID.(string))
	if err != nil {
		if err.Error() == "contract not found or already deleted" {
			utils.ErrorResponse(c, http.StatusNotFound, "Contract not found", "NOT_FOUND", nil)
			return
		}
		if err.Error() == "contract cannot be edited in its current status" {
			utils.ErrorResponse(c, http.StatusBadRequest, "Cannot edit contract", "STATUS_ERROR", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update contract", "UPDATE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Contract updated successfully", nil)
}

// DeleteContract soft deletes a contract
func (cc *ContractController) DeleteContract(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid contract ID", "VALIDATION_ERROR", err.Error())
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", "AUTH_ERROR", nil)
		return
	}

	// Validate access
	err = cc.contractService.ValidateContractAccess(id, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "ACCESS_ERROR", err.Error())
		return
	}

	err = cc.contractService.DeleteContract(id, userID.(string))
	if err != nil {
		if err.Error() == "contract not found or already deleted" {
			utils.ErrorResponse(c, http.StatusNotFound, "Contract not found", "NOT_FOUND", nil)
			return
		}
		if err.Error() == "only draft contracts can be deleted" {
			utils.ErrorResponse(c, http.StatusBadRequest, "Cannot delete contract", "STATUS_ERROR", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete contract", "DELETE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Contract deleted successfully", nil)
}

// ListContracts searches and lists contracts
func (cc *ContractController) ListContracts(c *gin.Context) {
	var req models.ContractSearchRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid query parameters", "VALIDATION_ERROR", err.Error())
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", "AUTH_ERROR", nil)
		return
	}

	// Filter by current user (Phase 1 restriction)
	req.CreatedBy = userID.(string)

	response, err := cc.contractService.SearchContracts(&req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to search contracts", "SEARCH_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Contracts retrieved successfully", response)
}

// ChangeContractStatus changes the status of a contract
func (cc *ContractController) ChangeContractStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid contract ID", "VALIDATION_ERROR", err.Error())
		return
	}

	var req models.ContractStatusChangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "VALIDATION_ERROR", err.Error())
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", "AUTH_ERROR", nil)
		return
	}

	// Validate access
	err = cc.contractService.ValidateContractAccess(id, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "ACCESS_ERROR", err.Error())
		return
	}

	err = cc.contractService.ChangeContractStatus(id, &req, userID.(string))
	if err != nil {
		if err.Error() == "contract not found" {
			utils.ErrorResponse(c, http.StatusNotFound, "Contract not found", "NOT_FOUND", nil)
			return
		}
		if len(err.Error()) > 15 && (err.Error()[:15] == "invalid status " || err.Error()[:26] == "invalid status transition ") {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid status change", "STATUS_ERROR", err.Error())
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to change contract status", "STATUS_CHANGE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Contract status changed successfully", nil)
}

// GetContractStatusHistory retrieves status change history for a contract
func (cc *ContractController) GetContractStatusHistory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid contract ID", "VALIDATION_ERROR", err.Error())
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", "AUTH_ERROR", nil)
		return
	}

	// Validate access
	err = cc.contractService.ValidateContractAccess(id, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied", "ACCESS_ERROR", err.Error())
		return
	}

	history, err := cc.contractService.GetContractStatusHistory(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get status history", "FETCH_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Status history retrieved successfully", history)
}

// GetContractStats returns statistics about contracts for the current user
func (cc *ContractController) GetContractStats(c *gin.Context) {
	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", "AUTH_ERROR", nil)
		return
	}

	stats, err := cc.contractService.GetContractStats(userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get contract statistics", "STATS_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Contract statistics retrieved successfully", stats)
}

// CreateContractVersion creates a new version of an existing contract
func (cc *ContractController) CreateContractVersion(c *gin.Context) {
	baseID := c.Param("baseId")

	var req models.ContractCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", "VALIDATION_ERROR", err.Error())
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "User not authenticated", "AUTH_ERROR", nil)
		return
	}

	contract, err := cc.contractService.CreateContractVersion(baseID, &req, userID.(string))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create contract version", "VERSION_CREATE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Contract version created successfully", contract)
}

// GetContractVersions retrieves all versions of a contract
func (cc *ContractController) GetContractVersions(c *gin.Context) {
	baseID := c.Param("baseId")

	contracts, err := cc.contractService.GetContractVersions(baseID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get contract versions", "VERSION_FETCH_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Contract versions retrieved successfully", contracts)
}