package controllers

import (
	"net/http"
	"strconv"

	"contrack-be/internal/models"
	"contrack-be/internal/repository"
	"contrack-be/internal/services/ai"
	"contrack-be/internal/utils"

	"github.com/gin-gonic/gin"
)

// AIController handles AI-related endpoints
type AIController struct {
	aiService      *ai.OpenAIService
	aiRepo         *repository.AIRepository
	clauseRepo     *repository.ClauseTemplateRepository
}

// NewAIController creates a new AI controller instance
func NewAIController(aiService *ai.OpenAIService, aiRepo *repository.AIRepository, clauseRepo *repository.ClauseTemplateRepository) *AIController {
	return &AIController{
		aiService:  aiService,
		aiRepo:     aiRepo,
		clauseRepo: clauseRepo,
	}
}

// AnalyzeClauseRisk analyzes a clause for potential risks using AI
// @Summary Analyze clause risk using AI
// @Description Analyzes a clause for potential legal risks and provides recommendations
// @Tags AI Analysis
// @Accept json
// @Produce json
// @Param request body models.ClauseRiskAnalysisRequest true "Analysis request"
// @Success 200 {object} models.ClauseRiskAnalysisResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/ai/analyze [post]
func (c *AIController) AnalyzeClauseRisk(ctx *gin.Context) {
	var req models.ClauseRiskAnalysisRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request format", "VALIDATION_ERROR", err.Error())
		return
	}

	// Get the clause from database
	clause, err := c.clauseRepo.GetByID(req.ClauseID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Clause not found", "CLAUSE_NOT_FOUND", err.Error())
		return
	}

	// Check if clause is active
	if !clause.IsActive {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Cannot analyze inactive clause", "INACTIVE_CLAUSE", "Clause is not active")
		return
	}

	// Perform AI analysis
	aiResult, err := c.aiService.AnalyzeClauseRisk(clause)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "AI analysis failed", "AI_ANALYSIS_ERROR", err.Error())
		return
	}

	// Convert AI result to database model
	analysis := &models.ClauseRiskAnalysis{
		ClauseID:          clause.ID,
		RiskLevel:         aiResult.RiskLevel,
		RiskScore:         aiResult.RiskScore,
		AnalysisSummary:   aiResult.AnalysisSummary,
		IdentifiedRisks:   aiResult.IdentifiedRisks,
		Recommendations:   aiResult.Recommendations,
		LegalImplications: aiResult.LegalImplications,
		ComplianceNotes:   aiResult.ComplianceNotes,
		ConfidenceScore:   aiResult.ConfidenceScore,
		ModelVersion:      c.aiService.GetModelVersion(),
	}

	// Save analysis to database
	if err := c.aiRepo.CreateAnalysis(analysis); err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to save analysis", "DATABASE_ERROR", err.Error())
		return
	}

	// Prepare response
	response := models.ClauseRiskAnalysisResponse{
		Analysis: *analysis,
		Clause:   *clause,
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Analysis completed successfully", response)
}

// GetAnalysisByID retrieves an analysis by its ID
// @Summary Get analysis by ID
// @Description Retrieves a specific AI analysis by its ID
// @Tags AI Analysis
// @Accept json
// @Produce json
// @Param id path int true "Analysis ID"
// @Success 200 {object} models.ClauseRiskAnalysisResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/ai/analysis/{id} [get]
func (c *AIController) GetAnalysisByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid analysis ID", "INVALID_ID", err.Error())
		return
	}

	analysis, err := c.aiRepo.GetAnalysisWithClause(id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Analysis not found", "ANALYSIS_NOT_FOUND", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Analysis retrieved successfully", analysis)
}

// GetAnalysisByClauseID retrieves the latest analysis for a specific clause
// @Summary Get analysis by clause ID
// @Description Retrieves the latest AI analysis for a specific clause
// @Tags AI Analysis
// @Accept json
// @Produce json
// @Param clause_id path int true "Clause ID"
// @Success 200 {object} models.ClauseRiskAnalysisResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/ai/analysis/clause/{clause_id} [get]
func (c *AIController) GetAnalysisByClauseID(ctx *gin.Context) {
	clauseIDStr := ctx.Param("clause_id")
	clauseID, err := strconv.Atoi(clauseIDStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid clause ID", "INVALID_CLAUSE_ID", err.Error())
		return
	}

	// Get the latest analysis for the clause
	analysis, err := c.aiRepo.GetAnalysisByClauseID(clauseID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Analysis not found for clause", "ANALYSIS_NOT_FOUND", err.Error())
		return
	}

	// Get the clause details
	clause, err := c.clauseRepo.GetByID(clauseID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Clause not found", "CLAUSE_NOT_FOUND", err.Error())
		return
	}

	response := models.ClauseRiskAnalysisResponse{
		Analysis: *analysis,
		Clause:   *clause,
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Analysis retrieved successfully", response)
}

// GetAnalyses retrieves analyses with pagination and filtering
// @Summary Get analyses with pagination
// @Description Retrieves AI analyses with pagination and filtering options
// @Tags AI Analysis
// @Accept json
// @Produce json
// @Param clause_id query int false "Filter by clause ID"
// @Param risk_level query string false "Filter by risk level (low, medium, high, critical)"
// @Param min_risk_score query number false "Minimum risk score"
// @Param max_risk_score query number false "Maximum risk score"
// @Param min_confidence query number false "Minimum confidence score"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Items per page (default: 10)"
// @Param sort_by query string false "Sort field (id, clause_id, risk_level, risk_score, created_at, updated_at)"
// @Param sort_dir query string false "Sort direction (asc, desc)"
// @Success 200 {object} models.ClauseRiskAnalysisListResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/ai/analyses [get]
func (c *AIController) GetAnalyses(ctx *gin.Context) {
	// Parse query parameters
	params := models.GetDefaultAnalysisSearchParams()

	if clauseIDStr := ctx.Query("clause_id"); clauseIDStr != "" {
		if clauseID, err := strconv.Atoi(clauseIDStr); err == nil {
			params.ClauseID = clauseID
		}
	}

	if riskLevel := ctx.Query("risk_level"); riskLevel != "" {
		params.RiskLevel = models.RiskLevel(riskLevel)
	}

	if minRiskScoreStr := ctx.Query("min_risk_score"); minRiskScoreStr != "" {
		if minRiskScore, err := strconv.ParseFloat(minRiskScoreStr, 64); err == nil {
			params.MinRiskScore = minRiskScore
		}
	}

	if maxRiskScoreStr := ctx.Query("max_risk_score"); maxRiskScoreStr != "" {
		if maxRiskScore, err := strconv.ParseFloat(maxRiskScoreStr, 64); err == nil {
			params.MaxRiskScore = maxRiskScore
		}
	}

	if minConfidenceStr := ctx.Query("min_confidence"); minConfidenceStr != "" {
		if minConfidence, err := strconv.ParseFloat(minConfidenceStr, 64); err == nil {
			params.MinConfidence = minConfidence
		}
	}

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

	if sortBy := ctx.Query("sort_by"); sortBy != "" {
		params.SortBy = sortBy
	}

	if sortDir := ctx.Query("sort_dir"); sortDir != "" {
		params.SortDir = sortDir
	}

	// Get analyses from repository
	response, err := c.aiRepo.GetAnalyses(params)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to retrieve analyses", "DATABASE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Analyses retrieved successfully", response)
}

// DeleteAnalysis deletes an analysis by ID
// @Summary Delete analysis
// @Description Deletes a specific AI analysis by its ID
// @Tags AI Analysis
// @Accept json
// @Produce json
// @Param id path int true "Analysis ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/ai/analysis/{id} [delete]
func (c *AIController) DeleteAnalysis(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid analysis ID", "INVALID_ID", err.Error())
		return
	}

	err = c.aiRepo.DeleteAnalysis(id)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Failed to delete analysis", "DELETE_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Analysis deleted successfully", nil)
}

// GetAnalysisStats retrieves statistics about analyses
// @Summary Get analysis statistics
// @Description Retrieves statistics about AI analyses
// @Tags AI Analysis
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/ai/stats [get]
func (c *AIController) GetAnalysisStats(ctx *gin.Context) {
	// This would require additional repository methods to get statistics
	// For now, return a placeholder response
	stats := map[string]interface{}{
		"total_analyses": 0,
		"risk_distribution": map[string]int{
			"low":      0,
			"medium":   0,
			"high":     0,
			"critical": 0,
		},
		"average_risk_score": 0.0,
		"average_confidence": 0.0,
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Statistics retrieved successfully", stats)
}

// AnalyzeContractRisk analyzes a contract with multiple clauses for potential risks using AI
// @Summary Analyze contract risk using AI
// @Description Analyzes a contract with multiple clauses for potential legal risks and provides recommendations
// @Tags AI Analysis
// @Accept json
// @Produce json
// @Param request body models.ContractAnalysisRequest true "Contract analysis request"
// @Success 200 {object} models.ContractAnalysisResult
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/ai/analyze-contract [post]
func (c *AIController) AnalyzeContractRisk(ctx *gin.Context) {
	var req models.ContractAnalysisRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid request format", "VALIDATION_ERROR", err.Error())
		return
	}

	// Get user ID from context
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.ErrorResponse(ctx, http.StatusUnauthorized, "User not authenticated", "AUTH_ERROR", nil)
		return
	}

	// Analyze contract with multiple clauses
	result, err := c.aiService.AnalyzeContractRisk(req.ContractID, req.ClauseTemplateIDs, userID.(string))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to analyze contract", "AI_ERROR", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Contract analysis completed successfully", result)
}
