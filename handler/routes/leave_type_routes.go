package routes

import (
	"hrm/handler"

	"github.com/gin-gonic/gin"
)

// SetupLeaveTypeRoutes configures the leave type routes
func SetupLeaveTypeRoutes(router *gin.Engine, leaveTypeHandler *handler.LeaveTypeHandler) {
	// Leave Type routes
	leaveTypeRoutes := router.Group("/api/leave-types")
	{
		leaveTypeRoutes.POST("", leaveTypeHandler.CreateLeaveType)
		leaveTypeRoutes.GET("", leaveTypeHandler.GetAllLeaveTypes)
		leaveTypeRoutes.GET("/active", leaveTypeHandler.GetActiveLeaveTypes)
		leaveTypeRoutes.GET("/stats", leaveTypeHandler.GetLeaveTypesWithUsageStats)
		leaveTypeRoutes.GET("/type/:type", leaveTypeHandler.GetLeaveTypeByType)
		leaveTypeRoutes.GET("/:id", leaveTypeHandler.GetLeaveTypeByID)
		leaveTypeRoutes.PUT("/:id", leaveTypeHandler.UpdateLeaveType)
		leaveTypeRoutes.DELETE("/:id", leaveTypeHandler.DeleteLeaveType)
	}
}
