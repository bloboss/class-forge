package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"code.forgejo.org/forgejo/classroom/internal/api/v1"
	"code.forgejo.org/forgejo/classroom/internal/config"
)

// NewRouter creates and configures the main API router
func NewRouter(cfg *config.Config, logger *zap.Logger) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Health check
	router.GET("/health", healthCheck)

	// API version info
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "forgejo-classroom",
			"version": "dev", // TODO: Get from build info
			"api_versions": []string{
				"v1",
			},
		})
	})

	// API v1 routes
	v1Group := router.Group("/api/v1")
	{
		// TODO: Add authentication middleware
		// v1Group.Use(authMiddleware())

		// Register v1 handlers
		v1.RegisterClassroomRoutes(v1Group, logger)
		v1.RegisterAssignmentRoutes(v1Group, logger)
		v1.RegisterRosterRoutes(v1Group, logger)
		v1.RegisterSubmissionRoutes(v1Group, logger)
		v1.RegisterTeamRoutes(v1Group, logger)
	}

	return router
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"service": "forgejo-classroom",
	})
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}