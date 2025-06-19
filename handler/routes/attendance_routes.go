package routes

import (
	"hrm/domain"
	"hrm/handler"
	"hrm/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAttendanceRoutes configures all attendance-related routes
func SetupAttendanceRoutes(router *gin.Engine, attendanceService domain.AttendanceServiceInterface) {
	// Create attendance handler
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)

	// Attendance API group
	attendanceGroup := router.Group("/api/v1/attendance")
	{
		// Public routes (no authentication required)
		attendanceGroup.POST("/checkin", attendanceHandler.CheckIn)
		attendanceGroup.POST("/checkout", attendanceHandler.CheckOut)

		// Protected routes (require authentication)
		protected := attendanceGroup.Group("")
		protected.Use(middleware.JWTAuthMiddleware())
		{
			// Attendance management
			protected.POST("/", attendanceHandler.CreateAttendance)
			protected.GET("/", attendanceHandler.GetAllAttendance)
			protected.GET("/:id", attendanceHandler.GetAttendanceByID)
			protected.DELETE("/:id", attendanceHandler.DeleteAttendance)

			// User-specific attendance
			protected.GET("/user/:user_id", attendanceHandler.GetUserAttendance)
			protected.POST("/user/range", attendanceHandler.GetUserAttendanceRange)
		}
	}
}
