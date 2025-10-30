package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"code.forgejo.org/forgejo/classroom/internal/response"
)

// ClassroomHandler handles classroom-related API endpoints
type ClassroomHandler struct {
	logger *zap.Logger
	// TODO: Add service dependencies
	// service *service.ClassroomService
}

// NewClassroomHandler creates a new classroom handler
func NewClassroomHandler(logger *zap.Logger) *ClassroomHandler {
	return &ClassroomHandler{
		logger: logger,
	}
}

// RegisterClassroomRoutes registers classroom routes with the router group
func RegisterClassroomRoutes(rg *gin.RouterGroup, logger *zap.Logger) {
	handler := NewClassroomHandler(logger)

	classrooms := rg.Group("/classrooms")
	{
		classrooms.POST("", handler.CreateClassroom)
		classrooms.GET("", handler.ListClassrooms)
		classrooms.GET("/:id", handler.GetClassroom)
		classrooms.PUT("/:id", handler.UpdateClassroom)
		classrooms.DELETE("/:id", handler.DeleteClassroom)
		classrooms.POST("/:id/archive", handler.ArchiveClassroom)
	}
}

// CreateClassroom handles POST /api/v1/classrooms
func (h *ClassroomHandler) CreateClassroom(c *gin.Context) {
	h.logger.Info("Creating classroom", zap.String("request_id", c.GetString("request_id")))

	// TODO: Implement classroom creation
	// 1. Parse and validate request body
	// 2. Call service layer to create classroom
	// 3. Return created classroom

	response.RespondWithData(c, http.StatusCreated, gin.H{
		"message": "Classroom creation not yet implemented",
		"todo":    "Parse request, validate input, call service layer",
	})
}

// ListClassrooms handles GET /api/v1/classrooms
func (h *ClassroomHandler) ListClassrooms(c *gin.Context) {
	h.logger.Info("Listing classrooms", zap.String("request_id", c.GetString("request_id")))

	// TODO: Implement classroom listing
	// 1. Parse query parameters (filters, pagination)
	// 2. Call service layer to get classrooms
	// 3. Return paginated list

	response.RespondWithData(c, http.StatusOK, gin.H{
		"message": "Classroom listing not yet implemented",
		"todo":    "Parse query params, call service layer, return paginated results",
	})
}

// GetClassroom handles GET /api/v1/classrooms/:id
func (h *ClassroomHandler) GetClassroom(c *gin.Context) {
	id := c.Param("id")
	h.logger.Info("Getting classroom", zap.String("id", id), zap.String("request_id", c.GetString("request_id")))

	// TODO: Implement classroom retrieval
	// 1. Validate ID parameter
	// 2. Call service layer to get classroom
	// 3. Return classroom details

	response.RespondWithData(c, http.StatusOK, gin.H{
		"message": "Classroom retrieval not yet implemented",
		"id":      id,
		"todo":    "Validate ID, call service layer, return classroom",
	})
}

// UpdateClassroom handles PUT /api/v1/classrooms/:id
func (h *ClassroomHandler) UpdateClassroom(c *gin.Context) {
	id := c.Param("id")
	h.logger.Info("Updating classroom", zap.String("id", id), zap.String("request_id", c.GetString("request_id")))

	// TODO: Implement classroom update
	// 1. Validate ID parameter
	// 2. Parse and validate request body
	// 3. Call service layer to update classroom
	// 4. Return updated classroom

	response.RespondWithData(c, http.StatusOK, gin.H{
		"message": "Classroom update not yet implemented",
		"id":      id,
		"todo":    "Parse request, validate input, call service layer",
	})
}

// DeleteClassroom handles DELETE /api/v1/classrooms/:id
func (h *ClassroomHandler) DeleteClassroom(c *gin.Context) {
	id := c.Param("id")
	h.logger.Info("Deleting classroom", zap.String("id", id), zap.String("request_id", c.GetString("request_id")))

	// TODO: Implement classroom deletion
	// 1. Validate ID parameter
	// 2. Check permissions
	// 3. Call service layer to delete classroom
	// 4. Return success response

	response.RespondWithData(c, http.StatusNoContent, gin.H{
		"message": "Classroom deletion not yet implemented",
		"id":      id,
		"todo":    "Validate ID, check permissions, call service layer",
	})
}

// ArchiveClassroom handles POST /api/v1/classrooms/:id/archive
func (h *ClassroomHandler) ArchiveClassroom(c *gin.Context) {
	id := c.Param("id")
	h.logger.Info("Archiving classroom", zap.String("id", id), zap.String("request_id", c.GetString("request_id")))

	// TODO: Implement classroom archiving
	// 1. Validate ID parameter
	// 2. Check permissions
	// 3. Call service layer to archive classroom
	// 4. Return archived classroom

	response.RespondWithData(c, http.StatusOK, gin.H{
		"message": "Classroom archiving not yet implemented",
		"id":      id,
		"todo":    "Validate ID, check permissions, call service layer",
	})
}