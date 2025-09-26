package controllers

import (
	"strconv"

	"contrack-be/internal/models"
	"contrack-be/internal/services/clause"
	"contrack-be/internal/utils"

	"github.com/gin-gonic/gin"
)

type ClauseController struct {
	clauseService *clause.Service
}

// NewClauseController creates a new clause controller
func NewClauseController() *ClauseController {
	return &ClauseController{
		clauseService: clause.NewClauseService(),
	}
}

// CreateClauseTemplate handles clause template creation
// POST /api/clauses
func (cc *ClauseController) CreateClauseTemplate(c *gin.Context) {
	var req models.ClauseTemplateCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	clauseTemplate, err := cc.clauseService.CreateClauseTemplate(&req)
	if err != nil {
		if err.Error() == "clause template with code '"+req.ClauseCode+"' already exists" {
			utils.ConflictResponse(c, "Clause template with this code already exists")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to create clause template")
		return
	}

	utils.CreatedResponse(c, "Clause template created successfully", gin.H{
		"clause_template": clauseTemplate,
	})
}

// GetClauseTemplate handles retrieving a single clause template
// GET /api/clauses/:id
func (cc *ClauseController) GetClauseTemplate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid clause template ID")
		return
	}

	clauseTemplate, err := cc.clauseService.GetClauseTemplateByID(id)
	if err != nil {
		if err.Error() == "clause template not found" {
			utils.NotFoundResponse(c, "Clause template not found")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to retrieve clause template")
		return
	}

	utils.OKResponse(c, "Clause template retrieved successfully", gin.H{
		"clause_template": clauseTemplate,
	})
}

// GetClauseTemplateByCode handles retrieving a clause template by code
// GET /api/clauses/by-code/:code
func (cc *ClauseController) GetClauseTemplateByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		utils.ValidationErrorResponse(c, "Clause code is required")
		return
	}

	clauseTemplate, err := cc.clauseService.GetClauseTemplateByCode(code)
	if err != nil {
		if err.Error() == "clause template not found" {
			utils.NotFoundResponse(c, "Clause template not found")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to retrieve clause template")
		return
	}

	utils.OKResponse(c, "Clause template retrieved successfully", gin.H{
		"clause_template": clauseTemplate,
	})
}

// UpdateClauseTemplate handles clause template updates
// PUT /api/clauses/:id
func (cc *ClauseController) UpdateClauseTemplate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid clause template ID")
		return
	}

	var req models.ClauseTemplateUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	clauseTemplate, err := cc.clauseService.UpdateClauseTemplate(id, &req)
	if err != nil {
		if err.Error() == "clause template not found" {
			utils.NotFoundResponse(c, "Clause template not found")
			return
		}

		if err.Error() == "clause template with code '"+*req.ClauseCode+"' already exists" {
			utils.ConflictResponse(c, "Clause template with this code already exists")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to update clause template")
		return
	}

	utils.OKResponse(c, "Clause template updated successfully", gin.H{
		"clause_template": clauseTemplate,
	})
}

// DeleteClauseTemplate handles clause template deletion
// DELETE /api/clauses/:id
func (cc *ClauseController) DeleteClauseTemplate(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid clause template ID")
		return
	}

	err = cc.clauseService.DeleteClauseTemplate(id)
	if err != nil {
		if err.Error() == "clause template not found" {
			utils.NotFoundResponse(c, "Clause template not found")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to delete clause template")
		return
	}

	utils.OKResponse(c, "Clause template deleted successfully", nil)
}

// ListClauseTemplates handles listing clause templates with pagination and search
// GET /api/clauses
func (cc *ClauseController) ListClauseTemplates(c *gin.Context) {
	var req models.ClauseTemplateSearchRequest

	// Bind query parameters
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	result, err := cc.clauseService.ListClauseTemplates(&req)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve clause templates")
		return
	}

	utils.OKResponse(c, "Clause templates retrieved successfully", gin.H{
		"clause_templates": result.ClauseTemplates,
		"pagination":       result.Pagination,
	})
}

// SearchClauseTemplates handles searching clause templates
// GET /api/clauses/search
func (cc *ClauseController) SearchClauseTemplates(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.ValidationErrorResponse(c, "Search query is required")
		return
	}

	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	result, err := cc.clauseService.SearchClauseTemplates(query, limit)
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to search clause templates")
		return
	}

	utils.OKResponse(c, "Search completed successfully", gin.H{
		"query":            query,
		"clause_templates": result.ClauseTemplates,
		"pagination":       result.Pagination,
	})
}

// GetClauseTypes handles retrieving all unique clause types
// GET /api/clauses/types
func (cc *ClauseController) GetClauseTypes(c *gin.Context) {
	types, err := cc.clauseService.GetClauseTypes()
	if err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve clause types")
		return
	}

	utils.OKResponse(c, "Clause types retrieved successfully", gin.H{
		"types": types,
		"count": len(types),
	})
}

// ToggleClauseTemplateStatus handles toggling the active status of a clause template
// PATCH /api/clauses/:id/toggle-status
func (cc *ClauseController) ToggleClauseTemplateStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid clause template ID")
		return
	}

	clauseTemplate, err := cc.clauseService.ToggleClauseTemplateStatus(id)
	if err != nil {
		if err.Error() == "clause template not found" {
			utils.NotFoundResponse(c, "Clause template not found")
			return
		}

		utils.InternalServerErrorResponse(c, "Failed to toggle clause template status")
		return
	}

	status := "deactivated"
	if clauseTemplate.IsActive {
		status = "activated"
	}

	utils.OKResponse(c, "Clause template "+status+" successfully", gin.H{
		"clause_template": clauseTemplate,
	})
}
