package handler

import (
	"hrm/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LeaveTypeHandler handles HTTP requests for leave type operations
type LeaveTypeHandler struct {
	leaveTypeService domain.LeaveTypeServiceInterface
}

// NewLeaveTypeHandler creates a new instance of LeaveTypeHandler
func NewLeaveTypeHandler(leaveTypeService domain.LeaveTypeServiceInterface) *LeaveTypeHandler {
	return &LeaveTypeHandler{
		leaveTypeService: leaveTypeService,
	}
}

// CreateLeaveType handles POST /api/leave-types
func (h *LeaveTypeHandler) CreateLeaveType(c *gin.Context) {
	var leaveType domain.LeaveType
	if err := c.ShouldBindJSON(&leaveType); err != nil {
		BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	if err := h.leaveTypeService.CreateLeaveType(&leaveType); err != nil {
		BadRequestResponse(c, "Failed to create leave type: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusCreated, "Leave type created successfully", leaveType)
}

// GetLeaveTypeByID handles GET /api/leave-types/:id
func (h *LeaveTypeHandler) GetLeaveTypeByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid leave type ID: "+err.Error())
		return
	}

	leaveType, err := h.leaveTypeService.GetLeaveTypeByID(uint(id))
	if err != nil {
		NotFoundResponse(c, "Leave type not found: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Leave type retrieved successfully", leaveType)
}

// GetLeaveTypeByType handles GET /api/leave-types/type/:type
func (h *LeaveTypeHandler) GetLeaveTypeByType(c *gin.Context) {
	leaveTypeStr := c.Param("type")

	leaveType, err := h.leaveTypeService.GetLeaveTypeByType(leaveTypeStr)
	if err != nil {
		NotFoundResponse(c, "Leave type not found: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Leave type retrieved successfully", leaveType)
}

// GetAllLeaveTypes handles GET /api/leave-types
func (h *LeaveTypeHandler) GetAllLeaveTypes(c *gin.Context) {
	leaveTypes, err := h.leaveTypeService.GetAllLeaveTypes()
	if err != nil {
		InternalServerErrorResponse(c, "Failed to retrieve leave types: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Leave types retrieved successfully", leaveTypes)
}

// GetActiveLeaveTypes handles GET /api/leave-types/active
func (h *LeaveTypeHandler) GetActiveLeaveTypes(c *gin.Context) {
	leaveTypes, err := h.leaveTypeService.GetActiveLeaveTypes()
	if err != nil {
		InternalServerErrorResponse(c, "Failed to retrieve active leave types: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Active leave types retrieved successfully", leaveTypes)
}

// UpdateLeaveType handles PUT /api/leave-types/:id
func (h *LeaveTypeHandler) UpdateLeaveType(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid leave type ID: "+err.Error())
		return
	}

	var leaveType domain.LeaveType
	if err := c.ShouldBindJSON(&leaveType); err != nil {
		BadRequestResponse(c, "Invalid request body: "+err.Error())
		return
	}

	leaveType.ID = uint(id)

	if err := h.leaveTypeService.UpdateLeaveType(&leaveType); err != nil {
		BadRequestResponse(c, "Failed to update leave type: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Leave type updated successfully", leaveType)
}

// DeleteLeaveType handles DELETE /api/leave-types/:id
func (h *LeaveTypeHandler) DeleteLeaveType(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid leave type ID: "+err.Error())
		return
	}

	if err := h.leaveTypeService.DeleteLeaveType(uint(id)); err != nil {
		BadRequestResponse(c, "Failed to delete leave type: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Leave type deleted successfully", nil)
}

// GetLeaveTypesWithUsageStats handles GET /api/leave-types/stats
func (h *LeaveTypeHandler) GetLeaveTypesWithUsageStats(c *gin.Context) {
	leaveTypes, err := h.leaveTypeService.GetLeaveTypesWithUsageStats()
	if err != nil {
		InternalServerErrorResponse(c, "Failed to retrieve leave types with usage stats: "+err.Error())
		return
	}

	SuccessResponse(c, http.StatusOK, "Leave types with usage stats retrieved successfully", leaveTypes)
}
