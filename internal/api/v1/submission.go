package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"code.forgejo.org/forgejo/classroom/internal/response"
)

// SubmissionHandler handles submission-related API endpoints
type SubmissionHandler struct {
	logger *zap.Logger
	// TODO: Add service dependencies
	// service *service.SubmissionService
}

// NewSubmissionHandler creates a new submission handler
func NewSubmissionHandler(logger *zap.Logger) *SubmissionHandler {
	return &SubmissionHandler{
		logger: logger,
	}
}

// RegisterSubmissionRoutes registers submission routes with the router group
func RegisterSubmissionRoutes(rg *gin.RouterGroup, logger *zap.Logger) {
	handler := NewSubmissionHandler(logger)

	submissions := rg.Group("/submissions")
	{
		submissions.GET("", handler.ListSubmissions)
		submissions.GET("/:id", handler.GetSubmission)
		submissions.GET("/:id/download", handler.DownloadSubmission)
	}

	// Assignment-specific submissions
	assignmentSubmissions := rg.Group("/assignments/:assignment_id/submissions")
	{
		assignmentSubmissions.GET("", handler.ListAssignmentSubmissions)
		assignmentSubmissions.GET("/download", handler.DownloadAllSubmissions)
	}
}

// ListSubmissions handles GET /api/v1/submissions
func (h *SubmissionHandler) ListSubmissions(c *gin.Context) {
	h.logger.Info("Listing submissions")

	// TODO: Implement submission listing
	response.RespondWithData(c, http.StatusOK, gin.H{
		"message": "Submission listing not yet implemented",
		"todo":    "Parse filters, call service layer, return paginated results",
	})
}

// GetSubmission handles GET /api/v1/submissions/:id
func (h *SubmissionHandler) GetSubmission(c *gin.Context) {
	id := c.Param("id")
	h.logger.Info("Getting submission", zap.String("id", id))

	// TODO: Implement submission retrieval
	response.RespondWithData(c, http.StatusOK, gin.H{
		"message": "Submission retrieval not yet implemented",
		"id":      id,
		"todo":    "Validate ID, call service layer, return submission details",
	})
}

// DownloadSubmission handles GET /api/v1/submissions/:id/download
func (h *SubmissionHandler) DownloadSubmission(c *gin.Context) {
	id := c.Param("id")
	h.logger.Info("Downloading submission", zap.String("id", id))

	// TODO: Implement submission download
	response.RespondWithData(c, http.StatusOK, gin.H{
		"message": "Submission download not yet implemented",
		"id":      id,
		"todo":    "Generate archive, stream file response",
	})
}

// ListAssignmentSubmissions handles GET /api/v1/assignments/:assignment_id/submissions
func (h *SubmissionHandler) ListAssignmentSubmissions(c *gin.Context) {
	assignmentID := c.Param("assignment_id")
	h.logger.Info("Listing assignment submissions", zap.String("assignment_id", assignmentID))

	// TODO: Implement assignment submissions listing
	response.RespondWithData(c, http.StatusOK, gin.H{
		"message": "Assignment submissions listing not yet implemented",
		"todo":    "Parse filters, call service layer, return submissions",
	})
}

// DownloadAllSubmissions handles GET /api/v1/assignments/:assignment_id/submissions/download
func (h *SubmissionHandler) DownloadAllSubmissions(c *gin.Context) {
	assignmentID := c.Param("assignment_id")
	h.logger.Info("Downloading all assignment submissions", zap.String("assignment_id", assignmentID))

	// TODO: Implement bulk submission download
	response.RespondWithData(c, http.StatusOK, gin.H{
		"message": "Bulk submission download not yet implemented",
		"todo":    "Generate bulk archive, stream response",
	})
}
