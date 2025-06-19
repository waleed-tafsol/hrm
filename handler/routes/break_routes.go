package routes

import (
	"hrm/domain"
	"hrm/handler"
	"hrm/middleware"

	"github.com/gin-gonic/gin"
)

// SetupBreakRoutes configures all break-related routes
func SetupBreakRoutes(router *gin.Engine, breakService domain.BreakServiceInterface) {
	// Create break handler
	breakHandler := handler.NewBreakHandler(breakService)

	// Break API group
	breakGroup := router.Group("/api/v1/breaks")
	{
		// Protected routes (require authentication)
		breakGroup.Use(middleware.JWTAuthMiddleware())
		{
			// Break management
			breakGroup.POST("/", breakHandler.AddBreak)
			breakGroup.GET("/", breakHandler.GetAllBreaks)
			breakGroup.GET("/:id", breakHandler.GetBreakByID)
			breakGroup.PUT("/end", breakHandler.EndBreak)
			breakGroup.DELETE("/:id", breakHandler.DeleteBreak)

			// Attendance-specific breaks
			breakGroup.GET("/attendance/:attendance_id", breakHandler.GetBreaksByAttendanceID)
		}
	}
}
