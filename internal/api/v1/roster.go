package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"code.forgejo.org/forgejo/classroom/internal/response"
)

// RosterHandler handles roster-related API endpoints
type RosterHandler struct {
	logger *zap.Logger
	// TODO: Add service dependencies
	// service *service.RosterService
}

// NewRosterHandler creates a new roster handler
func NewRosterHandler(logger *zap.Logger) *RosterHandler {
	return &RosterHandler{
		logger: logger,
	}
}

// RegisterRosterRoutes registers roster routes with the router group
func RegisterRosterRoutes(rg *gin.RouterGroup, logger *zap.Logger) {
	handler := NewRosterHandler(logger)

	rosters := rg.Group("/classrooms/:classroom_id/roster")
	{
		rosters.POST("/students", handler.AddStudent)
		rosters.GET("/students", handler.ListStudents)
		rosters.POST("/students/:student_id/link", handler.LinkStudent)
		rosters.POST("/import", handler.ImportRoster)
	}
}

// AddStudent handles POST /api/v1/classrooms/:classroom_id/roster/students
func (h *RosterHandler) AddStudent(c *gin.Context) {
	classroomID := c.Param("classroom_id")
	h.logger.Info("Adding student to roster", zap.String("classroom_id", classroomID))

	// TODO: Implement student addition
	response.RespondWithData(c, http.StatusCreated, gin.H{
		"message": "Student addition not yet implemented",
		"todo":    "Parse request, validate student data, call service layer",
	})
}

// ListStudents handles GET /api/v1/classrooms/:classroom_id/roster/students
func (h *RosterHandler) ListStudents(c *gin.Context) {
	classroomID := c.Param("classroom_id")
	h.logger.Info("Listing roster students", zap.String("classroom_id", classroomID))

	// TODO: Implement student listing
	response.RespondWithData(c, http.StatusOK, gin.H{
		"message": "Student listing not yet implemented",
		"todo":    "Parse filters, call service layer, return paginated results",
	})
}

// LinkStudent handles POST /api/v1/classrooms/:classroom_id/roster/students/:student_id/link
func (h *RosterHandler) LinkStudent(c *gin.Context) {
	classroomID := c.Param("classroom_id")
	studentID := c.Param("student_id")
	h.logger.Info("Linking student account", zap.String("classroom_id", classroomID), zap.String("student_id", studentID))

	// TODO: Implement account linking
	response.RespondWithData(c, http.StatusOK, gin.H{
		"message": "Account linking not yet implemented",
		"todo":    "Parse forgejo username, validate, call service layer",
	})
}

// ImportRoster handles POST /api/v1/classrooms/:classroom_id/roster/import
func (h *RosterHandler) ImportRoster(c *gin.Context) {
	classroomID := c.Param("classroom_id")
	h.logger.Info("Importing roster", zap.String("classroom_id", classroomID))

	// TODO: Implement roster import
	response.RespondWithData(c, http.StatusOK, gin.H{
		"message": "Roster import not yet implemented",
		"todo":    "Parse CSV file, validate data, bulk import students",
	})
}
