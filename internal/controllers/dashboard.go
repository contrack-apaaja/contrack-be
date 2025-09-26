package controllers

import (
	"net/http"

	"contrack-be/internal/repository"
	"contrack-be/internal/utils"

	"github.com/gin-gonic/gin"
)

// DashboardController handles dashboard-related endpoints
type DashboardController struct {
	contractRepo *repository.ContractRepository
}

// NewDashboardController creates a new Dashboard controller instance
func NewDashboardController(contractRepo *repository.ContractRepository) *DashboardController {
	return &DashboardController{
		contractRepo: contractRepo,
	}
}

// GetStatusCounts retrieves count for each contract status with visualization data
// @Summary Get contract status counts with visualization data
// @Description Retrieves the count of contracts for each status plus monthly trends and totals
// @Tags Dashboard
// @Accept json
// @Produce json
// @Success 200 {object} models.DashboardVisualizationData
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/dashboard/status-counts [get]
func (c *DashboardController) GetStatusCounts(ctx *gin.Context) {
	visualizationData, err := c.contractRepo.GetDashboardVisualizationData()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve dashboard visualization data", "DATABASE_ERROR", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Dashboard visualization data retrieved successfully", visualizationData)
}

// GetContractList retrieves simple contract list for dashboard table
// @Summary Get contract list for dashboard
// @Description Retrieves a simple list of contracts for the current user
// @Tags Dashboard
// @Accept json
// @Produce json
// @Success 200 {array} models.SimpleContractList
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/dashboard/contracts [get]
func (c *DashboardController) GetContractList(ctx *gin.Context) {
	// Get user ID from JWT token (assuming it's stored in context by auth middleware)
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User ID not found in context", "AUTH_ERROR", "User not authenticated")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Invalid user ID format", "AUTH_ERROR", "User ID is not a string")
		return
	}

	contracts, err := c.contractRepo.GetSimpleContractList(userIDStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve contract list", "DATABASE_ERROR", err.Error())
		return
	}
	utils.SuccessResponse(ctx, http.StatusOK, "Contract list retrieved successfully", contracts)
}