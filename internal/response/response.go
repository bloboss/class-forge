package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains detailed error information
type ErrorDetail struct {
	Code             string                 `json:"code"`
	Message          string                 `json:"message"`
	Details          map[string]interface{} `json:"details,omitempty"`
	RequestID        string                 `json:"request_id"`
	Timestamp        string                 `json:"timestamp"`
	DocumentationURL string                 `json:"documentation_url,omitempty"`
}

// SuccessResponse represents a standardized success response
type SuccessResponse struct {
	Data interface{} `json:"data"`
	Meta *MetaInfo   `json:"meta,omitempty"`
}

// MetaInfo contains metadata about the response
type MetaInfo struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"per_page,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
	TotalCount int `json:"total_count,omitempty"`
}

// RespondWithError sends a standardized error response
func RespondWithError(c *gin.Context, statusCode int, errorCode, message string, details map[string]interface{}) {
	requestID := getRequestID(c)

	response := ErrorResponse{
		Error: ErrorDetail{
			Code:      errorCode,
			Message:   message,
			Details:   details,
			RequestID: requestID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Add documentation URL for known error codes
	if docURL := getDocumentationURL(errorCode); docURL != "" {
		response.Error.DocumentationURL = docURL
	}

	c.JSON(statusCode, response)
}

// RespondWithSuccess sends a standardized success response
func RespondWithSuccess(c *gin.Context, statusCode int, data interface{}, meta *MetaInfo) {
	response := SuccessResponse{
		Data: data,
		Meta: meta,
	}

	c.JSON(statusCode, response)
}

// RespondWithData is a shorthand for success responses without metadata
func RespondWithData(c *gin.Context, statusCode int, data interface{}) {
	RespondWithSuccess(c, statusCode, data, nil)
}

// Common error response helpers
func BadRequest(c *gin.Context, errorCode, message string, details map[string]interface{}) {
	RespondWithError(c, http.StatusBadRequest, errorCode, message, details)
}

func Unauthorized(c *gin.Context, errorCode, message string) {
	RespondWithError(c, http.StatusUnauthorized, errorCode, message, nil)
}

func Forbidden(c *gin.Context, errorCode, message string) {
	RespondWithError(c, http.StatusForbidden, errorCode, message, nil)
}

func NotFound(c *gin.Context, errorCode, message string) {
	RespondWithError(c, http.StatusNotFound, errorCode, message, nil)
}

func Conflict(c *gin.Context, errorCode, message string, details map[string]interface{}) {
	RespondWithError(c, http.StatusConflict, errorCode, message, details)
}

func InternalServerError(c *gin.Context, errorCode, message string) {
	RespondWithError(c, http.StatusInternalServerError, errorCode, message, nil)
}

// Helper functions
func getRequestID(c *gin.Context) string {
	// TODO: Implement request ID generation/extraction from middleware
	return "req_" + time.Now().Format("20060102150405")
}

func getDocumentationURL(errorCode string) string {
	// TODO: Return documentation URLs for error codes
	baseURL := "https://docs.forgejo-classroom.org/errors/"
	return baseURL + errorCode
}
