package controllers

import (
	"net/http"
	"strconv"

	"contrack-be/internal/models"
	"contrack-be/internal/repository"
	"contrack-be/internal/utils"

	"github.com/gin-gonic/gin"
)

// DashboardController handles dashboard-related endpoints
type DashboardController struct {
	contractRepo *repository.ContractRepository
}

// NewDashboardController creates a new dashboard controller instance
func NewDashboardController(contractRepo *repository.ContractRepository) *DashboardController {
	return &DashboardController{
		contractRepo: contractRepo,
	}
}

// GetDashboardStats retrieves comprehensive dashboard statistics
// @Summary Get dashboard statistics
// @Description Retrieves comprehensive dashboard statistics including contract counts, status stats, and recent contracts
// @Tags Dashboard
// @Accept json
// @Produce json
// @Success 200 {object} models.DashboardStats
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/dashboard/stats [get]
func (c *DashboardController) GetDashboardStats(ctx *gin.Context) {
	stats, err := c.contractRepo.GetDashboardStats()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve dashboard stats", "DATABASE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Dashboard statistics retrieved successfully", stats)
}

// GetDashboardContracts retrieves contracts for dashboard with filtering and pagination
// @Summary Get dashboard contracts
// @Description Retrieves contracts for dashboard with summary information, filtering, and pagination
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Param status query string false "Filter by contract status"
// @Param contract_type query string false "Filter by contract type"
// @Param sort_by query string false "Sort field (project_name, contract_number, total_value, signing_date, status, created_at)"
// @Param sort_dir query string false "Sort direction (asc, desc)"
// @Success 200 {array} models.DashboardContractSummary
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/dashboard/contracts [get]
func (c *DashboardController) GetDashboardContracts(ctx *gin.Context) {
	// Parse query parameters
	params := models.GetDefaultDashboardParams()

	if pageStr := ctx.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
	}

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			params.Limit = limit
		}
	}

	if status := ctx.Query("status"); status != "" {
		params.Status = status
	}

	if contractType := ctx.Query("contract_type"); contractType != "" {
		params.ContractType = contractType
	}

	if sortBy := ctx.Query("sort_by"); sortBy != "" {
		params.SortBy = sortBy
	}

	if sortDir := ctx.Query("sort_dir"); sortDir != "" {
		params.SortDir = sortDir
	}

	// Get contracts from repository
	contracts, err := c.contractRepo.GetDashboardContracts(params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve dashboard contracts", "DATABASE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Dashboard contracts retrieved successfully", contracts)
}

// GetContractStatusStats retrieves statistics for each contract status
// @Summary Get contract status statistics
// @Description Retrieves statistics for each contract status including count, percentage, and total value
// @Tags Dashboard
// @Accept json
// @Produce json
// @Success 200 {array} models.ContractStatusStats
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/dashboard/status-stats [get]
func (c *DashboardController) GetContractStatusStats(ctx *gin.Context) {
	stats, err := c.contractRepo.GetContractStatusStats()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve status statistics", "DATABASE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Status statistics retrieved successfully", stats)
}

// GetRecentContracts retrieves recent contracts for dashboard
// @Summary Get recent contracts
// @Description Retrieves recent contracts for dashboard display
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param limit query int false "Number of recent contracts to retrieve (default: 5, max: 20)"
// @Success 200 {array} models.DashboardContractSummary
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/dashboard/recent [get]
func (c *DashboardController) GetRecentContracts(ctx *gin.Context) {
	limit := 5
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 20 {
			limit = l
		}
	}

	params := models.DashboardRequest{
		Page:    1,
		Limit:   limit,
		SortBy:  "created_at",
		SortDir: "desc",
	}

	contracts, err := c.contractRepo.GetDashboardContracts(params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve recent contracts", "DATABASE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Recent contracts retrieved successfully", contracts)
}

// GetExpiringContracts retrieves contracts that are expiring soon
// @Summary Get expiring contracts
// @Description Retrieves contracts that are expiring soon (within 30 days)
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param limit query int false "Number of expiring contracts to retrieve (default: 5, max: 20)"
// @Success 200 {array} models.DashboardContractSummary
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/dashboard/expiring [get]
func (c *DashboardController) GetExpiringContracts(ctx *gin.Context) {
	limit := 5
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 20 {
			limit = l
		}
	}

	params := models.DashboardRequest{
		Page:    1,
		Limit:   limit,
		SortBy:  "signing_date",
		SortDir: "asc",
	}

	contracts, err := c.contractRepo.GetExpiringContracts(params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve expiring contracts", "DATABASE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Expiring contracts retrieved successfully", contracts)
}

// GetHighValueContracts retrieves high value contracts for dashboard
// @Summary Get high value contracts
// @Description Retrieves high value contracts sorted by total value
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param limit query int false "Number of high value contracts to retrieve (default: 5, max: 20)"
// @Success 200 {array} models.DashboardContractSummary
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/dashboard/high-value [get]
func (c *DashboardController) GetHighValueContracts(ctx *gin.Context) {
	limit := 5
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 20 {
			limit = l
		}
	}

	params := models.DashboardRequest{
		Page:    1,
		Limit:   limit,
		SortBy:  "total_value",
		SortDir: "desc",
	}

	contracts, err := c.contractRepo.GetDashboardContracts(params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve high value contracts", "DATABASE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "High value contracts retrieved successfully", contracts)
}
