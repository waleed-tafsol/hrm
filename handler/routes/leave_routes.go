package routes

import (
	"hrm/domain"
	"hrm/handler"

	"github.com/gin-gonic/gin"
)

// SetupLeaveRoutes configures all leave-related routes
func SetupLeaveRoutes(router *gin.Engine, leaveService domain.LeaveServiceInterface, userService domain.UserServiceInterface) {
	handler.SetupLeaveRoutes(router, leaveService, userService)
}
