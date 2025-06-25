package routes

import (
	"hrm/domain"
	"hrm/handler"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configures all user-related routes
func SetupUserRoutes(router *gin.Engine, userService domain.UserServiceInterface, attendanceService domain.AttendanceServiceInterface) {
	handler.SetupUserRoutes(router, userService, attendanceService)
}

// SetupHealthRoutes configures health check routes
func SetupHealthRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "HRM API is running",
		})
	})
}
