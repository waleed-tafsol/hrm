package handler

import (
	"hrm/domain"
	"hrm/handler/request"
	"hrm/handler/response"
	"hrm/middleware"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// LeaveHandler handles HTTP requests for leave operations
type LeaveHandler struct {
	leaveService domain.LeaveServiceInterface
	userService  domain.UserServiceInterface
}

// NewLeaveHandler creates a new instance of LeaveHandler
func NewLeaveHandler(leaveService domain.LeaveServiceInterface, userService domain.UserServiceInterface) *LeaveHandler {
	return &LeaveHandler{
		leaveService: leaveService,
		userService:  userService,
	}
}

// SetupLeaveRoutes sets up all leave-related HTTP routes
func SetupLeaveRoutes(router *gin.Engine, leaveService domain.LeaveServiceInterface, userService domain.UserServiceInterface) {
	handler := NewLeaveHandler(leaveService, userService)

	// Leave API group
	leaveGroup := router.Group("/api/leaves")
	{
		// Protected routes (require authentication)
		leaveGroup.Use(middleware.JWTAuthMiddleware())
		{
			// Leave management
			leaveGroup.POST("/", handler.CreateLeave)
			leaveGroup.GET("/", handler.GetAllLeaves)
			leaveGroup.GET("/pending", handler.GetPendingLeaves)
			leaveGroup.GET("/:id", handler.GetLeaveByID)
			leaveGroup.PUT("/:id", handler.UpdateLeave)
			leaveGroup.DELETE("/:id", handler.DeleteLeave)

			// Leave approval/rejection
			leaveGroup.POST("/:id/approve", handler.ApproveLeave)
			leaveGroup.POST("/:id/reject", handler.RejectLeave)
			leaveGroup.POST("/:id/cancel", handler.CancelLeave)

			// User-specific leaves
			leaveGroup.GET("/user/:user_id", handler.GetUserLeaves)
			leaveGroup.GET("/user/:user_id/range", handler.GetUserLeavesByDateRange)
			leaveGroup.GET("/user/:user_id/balance", handler.GetUserLeaveBalance)
		}
	}
}

// CreateLeave handles POST /api/leaves
func (h *LeaveHandler) CreateLeave(c *gin.Context) {
	// Get authenticated user ID
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	var req request.CreateLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	leave := &domain.Leave{
		Type:        domain.LeaveTypeName(req.Type.Type),
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Reason:      req.Reason,
		Description: req.Description,
	}

	if err := h.leaveService.CreateLeave(userID, leave); err != nil {
		BadRequestResponse(c, "Failed to create leave: "+err.Error())
		return
	}

	leaveResponse := response.ToLeaveResponse(leave)
	SuccessResponse(c, http.StatusCreated, "Leave created successfully", response.CreateLeaveResponse{
		Leave:   leaveResponse,
		Message: "Leave created successfully",
	})
}

// GetLeaveByID handles GET /api/leaves/:id
func (h *LeaveHandler) GetLeaveByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid leave ID: "+err.Error())
		return
	}

	leave, err := h.leaveService.GetLeaveByID(uint(id))
	if err != nil {
		NotFoundResponse(c, "Leave not found: "+err.Error())
		return
	}

	leaveResponse := response.ToLeaveResponse(leave)
	SuccessResponse(c, http.StatusOK, "Leave retrieved successfully", response.GetLeaveResponse{
		Leave: leaveResponse,
	})
}

// GetAllLeaves handles GET /api/leaves
func (h *LeaveHandler) GetAllLeaves(c *gin.Context) {
	leaves, err := h.leaveService.GetAllLeaves()
	if err != nil {
		InternalServerErrorResponse(c, "Failed to retrieve leaves: "+err.Error())
		return
	}

	leaveResponses := response.ToLeaveResponseList(leaves)
	SuccessResponse(c, http.StatusOK, "Leaves retrieved successfully", response.ListLeavesResponse{
		Leaves: leaveResponses,
		Total:  len(leaveResponses),
	})
}

// GetPendingLeaves handles GET /api/leaves/pending
func (h *LeaveHandler) GetPendingLeaves(c *gin.Context) {
	leaves, err := h.leaveService.GetPendingLeaves()
	if err != nil {
		InternalServerErrorResponse(c, "Failed to retrieve pending leaves: "+err.Error())
		return
	}

	leaveResponses := response.ToLeaveResponseList(leaves)
	SuccessResponse(c, http.StatusOK, "Pending leaves retrieved successfully", response.ListLeavesResponse{
		Leaves: leaveResponses,
		Total:  len(leaveResponses),
	})
}

// UpdateLeave handles PUT /api/leaves/:id
func (h *LeaveHandler) UpdateLeave(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid leave ID: "+err.Error())
		return
	}

	var req request.UpdateLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	leave := &domain.Leave{
		ID:          uint(id),
		Type:        domain.LeaveTypeName(req.Type.Type),
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Reason:      req.Reason,
		Description: req.Description,
	}

	if err := h.leaveService.UpdateLeave(leave); err != nil {
		BadRequestResponse(c, "Failed to update leave: "+err.Error())
		return
	}

	leaveResponse := response.ToLeaveResponse(leave)
	SuccessResponse(c, http.StatusOK, "Leave updated successfully", response.UpdateLeaveResponse{
		Leave:   leaveResponse,
		Message: "Leave updated successfully",
	})
}

// DeleteLeave handles DELETE /api/leaves/:id
func (h *LeaveHandler) DeleteLeave(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid leave ID: "+err.Error())
		return
	}

	if err := h.leaveService.DeleteLeave(uint(id)); err != nil {
		BadRequestResponse(c, "Failed to delete leave: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Leave deleted successfully", response.DeleteLeaveResponse{
		Message: "Leave deleted successfully",
	})
}

// ApproveLeave handles POST /api/leaves/:id/approve
func (h *LeaveHandler) ApproveLeave(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid leave ID: "+err.Error())
		return
	}

	// Get authenticated user ID (approver)
	approverID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	if err := h.leaveService.ApproveLeave(uint(id), approverID); err != nil {
		BadRequestResponse(c, "Failed to approve leave: "+err.Error())
		return
	}

	// Get updated leave
	leave, err := h.leaveService.GetLeaveByID(uint(id))
	if err != nil {
		InternalServerErrorResponse(c, "Failed to retrieve updated leave: "+err.Error())
		return
	}

	leaveResponse := response.ToLeaveResponse(leave)
	SuccessResponse(c, http.StatusOK, "Leave approved successfully", response.ApproveLeaveResponse{
		Leave:   leaveResponse,
		Message: "Leave approved successfully",
	})
}

// RejectLeave handles POST /api/leaves/:id/reject
func (h *LeaveHandler) RejectLeave(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid leave ID: "+err.Error())
		return
	}

	var req request.RejectLeaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	// Get authenticated user ID (rejecter)
	rejecterID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	if err := h.leaveService.RejectLeave(uint(id), rejecterID, req.RejectReason); err != nil {
		BadRequestResponse(c, "Failed to reject leave: "+err.Error())
		return
	}

	// Get updated leave
	leave, err := h.leaveService.GetLeaveByID(uint(id))
	if err != nil {
		InternalServerErrorResponse(c, "Failed to retrieve updated leave: "+err.Error())
		return
	}

	leaveResponse := response.ToLeaveResponse(leave)
	SuccessResponse(c, http.StatusOK, "Leave rejected successfully", response.RejectLeaveResponse{
		Leave:   leaveResponse,
		Message: "Leave rejected successfully",
	})
}

// CancelLeave handles POST /api/leaves/:id/cancel
func (h *LeaveHandler) CancelLeave(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid leave ID: "+err.Error())
		return
	}

	// Get authenticated user ID
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	if err := h.leaveService.CancelLeave(uint(id), userID); err != nil {
		BadRequestResponse(c, "Failed to cancel leave: "+err.Error())
		return
	}

	// Get updated leave
	leave, err := h.leaveService.GetLeaveByID(uint(id))
	if err != nil {
		InternalServerErrorResponse(c, "Failed to retrieve updated leave: "+err.Error())
		return
	}

	leaveResponse := response.ToLeaveResponse(leave)
	SuccessResponse(c, http.StatusOK, "Leave cancelled successfully", response.CancelLeaveResponse{
		Leave:   leaveResponse,
		Message: "Leave cancelled successfully",
	})
}

// GetUserLeaves handles GET /api/leaves/user/:user_id
func (h *LeaveHandler) GetUserLeaves(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid user ID: "+err.Error())
		return
	}

	leaves, err := h.leaveService.GetUserLeaves(uint(userID))
	if err != nil {
		InternalServerErrorResponse(c, "Failed to retrieve user leaves: "+err.Error())
		return
	}

	leaveResponses := response.ToLeaveResponseList(leaves)
	SuccessResponse(c, http.StatusOK, "User leaves retrieved successfully", response.GetUserLeavesResponse{
		Leaves: leaveResponses,
		UserID: uint(userID),
	})
}

// GetUserLeavesByDateRange handles GET /api/leaves/user/:user_id/range
func (h *LeaveHandler) GetUserLeavesByDateRange(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid user ID: "+err.Error())
		return
	}

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		BadRequestResponse(c, "Invalid start date format. Use YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		BadRequestResponse(c, "Invalid end date format. Use YYYY-MM-DD")
		return
	}

	leaves, err := h.leaveService.GetUserLeavesByDateRange(uint(userID), startDate, endDate)
	if err != nil {
		InternalServerErrorResponse(c, "Failed to retrieve user leaves by date range: "+err.Error())
		return
	}

	leaveResponses := response.ToLeaveResponseList(leaves)
	SuccessResponse(c, http.StatusOK, "User leaves by date range retrieved successfully", response.GetUserLeavesResponse{
		Leaves: leaveResponses,
		UserID: uint(userID),
	})
}

// GetUserLeaveBalance handles GET /api/leaves/user/:user_id/balance
func (h *LeaveHandler) GetUserLeaveBalance(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid user ID: "+err.Error())
		return
	}

	yearStr := c.Query("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		BadRequestResponse(c, "Invalid year parameter")
		return
	}

	balance, err := h.leaveService.GetUserLeaveBalance(uint(userID), year)
	if err != nil {
		InternalServerErrorResponse(c, "Failed to retrieve user leave balance: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "User leave balance retrieved successfully", response.GetUserLeaveBalanceResponse{
		UserID:  uint(userID),
		Year:    year,
		Balance: balance,
	})
}
