package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIResponse represents the standard API response structure
type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

// ErrorInfo contains detailed error information
type ErrorInfo struct {
	Code    string      `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

// Response status constants
const (
	StatusSuccess = "success"
	StatusError   = "error"
	StatusFail    = "fail"
)

// SuccessResponse sends a successful response with data
func SuccessResponse(c *gin.Context, httpStatus int, message string, data interface{}) {
	response := APIResponse{
		Status:  StatusSuccess,
		Message: message,
		Data:    data,
	}
	c.JSON(httpStatus, response)
}

// ErrorResponse sends an error response
func ErrorResponse(c *gin.Context, httpStatus int, message string, errorCode string, details interface{}) {
	response := APIResponse{
		Status:  StatusError,
		Message: message,
		Error: &ErrorInfo{
			Code:    errorCode,
			Details: details,
		},
	}
	c.JSON(httpStatus, response)
}

// FailResponse sends a fail response (client error - 4xx)
func FailResponse(c *gin.Context, httpStatus int, message string, details interface{}) {
	response := APIResponse{
		Status:  StatusFail,
		Message: message,
		Data:    details,
	}
	c.JSON(httpStatus, response)
}

// ValidationErrorResponse sends a validation error response
func ValidationErrorResponse(c *gin.Context, details interface{}) {
	ErrorResponse(c, http.StatusBadRequest, "Validation failed", "VALIDATION_ERROR", details)
}

// UnauthorizedResponse sends an unauthorized response
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message, "UNAUTHORIZED", nil)
}

// ForbiddenResponse sends a forbidden response
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message, "FORBIDDEN", nil)
}

// NotFoundResponse sends a not found response
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message, "NOT_FOUND", nil)
}

// ConflictResponse sends a conflict response
func ConflictResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusConflict, message, "CONFLICT", nil)
}

// InternalServerErrorResponse sends an internal server error response
func InternalServerErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message, "INTERNAL_ERROR", nil)
}

// CreatedResponse sends a created response (201)
func CreatedResponse(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusCreated, message, data)
}

// OKResponse sends an OK response (200)
func OKResponse(c *gin.Context, message string, data interface{}) {
	SuccessResponse(c, http.StatusOK, message, data)
}

// NoContentResponse sends a no content response (204)
func NoContentResponse(c *gin.Context, message string) {
	response := APIResponse{
		Status:  StatusSuccess,
		Message: message,
	}
	c.JSON(http.StatusNoContent, response)
}
