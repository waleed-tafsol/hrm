package handler

import (
	"net/http"
	"strconv"
	"time"

	"hrm/domain"
	"hrm/handler/request"
	"hrm/handler/response"
	"hrm/middleware"

	"github.com/gin-gonic/gin"
)

// AttendanceHandler handles HTTP requests related to attendance operations
type AttendanceHandler struct {
	attendanceService domain.AttendanceServiceInterface
}

// NewAttendanceHandler creates a new instance of AttendanceHandler
func NewAttendanceHandler(attendanceService domain.AttendanceServiceInterface) *AttendanceHandler {
	return &AttendanceHandler{
		attendanceService: attendanceService,
	}
}

// CreateAttendance creates a new attendance record
func (h *AttendanceHandler) CreateAttendance(c *gin.Context) {
	var req request.AttendanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request data: "+err.Error())
		return
	}

	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		UnauthorizedResponse(c, "User ID not found in token")
		return
	}

	attendance, err := h.attendanceService.CreateAttendance(userID, req.Date)
	if err != nil {
		InternalServerErrorResponse(c, "Failed to create attendance: "+err.Error())
		return
	}

	attendanceResp := h.convertToAttendanceResponse(*attendance)
	SuccessResponse(c, http.StatusCreated, "Attendance created successfully", attendanceResp)
}

// CheckIn handles user check-in
func (h *AttendanceHandler) CheckIn(c *gin.Context) {
	var req request.CheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request data: "+err.Error())
		return
	}

	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		UnauthorizedResponse(c, "User ID not found in token")
		return
	}

	attendance, err := h.attendanceService.CheckIn(userID, req.Date)
	if err != nil {
		if err == domain.ErrAlreadyCheckedIn {
			c.JSON(http.StatusConflict, Response{
				Success: false,
				Message: "Already checked in for this date",
			})
		} else if err == domain.ErrUserNotFound {
			NotFoundResponse(c, "User not found")
		} else {
			InternalServerErrorResponse(c, "Failed to check in: "+err.Error())
		}
		return
	}

	attendanceResp := h.convertToAttendanceResponse(*attendance)
	checkInResp := response.CheckInResponse{
		Attendance:  attendanceResp,
		Message:     "Check-in successful",
		CheckInTime: *attendance.CheckInTime,
	}

	SuccessResponse(c, http.StatusOK, "Check-in successful", checkInResp)
}

// CheckOut handles user check-out
func (h *AttendanceHandler) CheckOut(c *gin.Context) {
	var req request.CheckOutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request data: "+err.Error())
		return
	}

	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		UnauthorizedResponse(c, "User ID not found in token")
		return
	}

	attendance, err := h.attendanceService.CheckOut(userID, req.Date)
	if err != nil {
		if err == domain.ErrAlreadyCheckedOut {
			c.JSON(http.StatusConflict, Response{
				Success: false,
				Message: "Already checked out for this date",
			})
		} else if err == domain.ErrNotCheckedIn {
			BadRequestResponse(c, "Not checked in yet")
		} else if err == domain.ErrUserNotFound || err == domain.ErrAttendanceNotFound {
			NotFoundResponse(c, "User or attendance not found")
		} else {
			InternalServerErrorResponse(c, "Failed to check out: "+err.Error())
		}
		return
	}

	attendanceResp := h.convertToAttendanceResponse(*attendance)
	checkOutResp := response.CheckOutResponse{
		Attendance:     attendanceResp,
		Message:        "Check-out successful",
		CheckOutTime:   *attendance.CheckOutTime,
		TotalWorkHours: attendance.TotalWorkHours,
	}

	SuccessResponse(c, http.StatusOK, "Check-out successful", checkOutResp)
}

// GetAttendanceByID retrieves an attendance record by ID
func (h *AttendanceHandler) GetAttendanceByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid attendance ID")
		return
	}

	attendance, err := h.attendanceService.GetAttendanceByID(uint(id))
	if err != nil {
		if err == domain.ErrAttendanceNotFound {
			NotFoundResponse(c, "Attendance not found")
		} else {
			InternalServerErrorResponse(c, "Failed to get attendance: "+err.Error())
		}
		return
	}

	attendanceResp := h.convertToAttendanceResponse(*attendance)
	SuccessResponse(c, http.StatusOK, "Attendance retrieved successfully", attendanceResp)
}

// GetUserAttendance retrieves attendance for a user on a specific date
func (h *AttendanceHandler) GetUserAttendance(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid user ID")
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		BadRequestResponse(c, "Invalid date format. Use YYYY-MM-DD")
		return
	}

	attendance, err := h.attendanceService.GetUserAttendance(uint(userID), date)
	if err != nil {
		if err == domain.ErrAttendanceNotFound || err == domain.ErrUserNotFound {
			NotFoundResponse(c, "User or attendance not found")
		} else {
			InternalServerErrorResponse(c, "Failed to get user attendance: "+err.Error())
		}
		return
	}

	attendanceResp := h.convertToAttendanceResponse(*attendance)
	SuccessResponse(c, http.StatusOK, "User attendance retrieved successfully", attendanceResp)
}

// GetUserAttendanceRange retrieves attendance records for a user within a date range
func (h *AttendanceHandler) GetUserAttendanceRange(c *gin.Context) {
	var req request.AttendanceRangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request data: "+err.Error())
		return
	}

	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		UnauthorizedResponse(c, "User ID not found in token")
		return
	}

	attendances, err := h.attendanceService.GetUserAttendanceRange(userID, req.StartDate, req.EndDate)
	if err != nil {
		if err == domain.ErrUserNotFound {
			NotFoundResponse(c, "User not found")
		} else {
			InternalServerErrorResponse(c, "Failed to get attendance range: "+err.Error())
		}
		return
	}

	attendanceResponses := make([]response.AttendanceResponse, len(attendances))
	for i, attendance := range attendances {
		attendanceResponses[i] = h.convertToAttendanceResponse(attendance)
	}

	listResp := response.AttendanceListResponse{
		Attendances: attendanceResponses,
		Total:       len(attendanceResponses),
	}

	SuccessResponse(c, http.StatusOK, "Attendance range retrieved successfully", listResp)
}

// GetAllAttendance retrieves all attendance records
func (h *AttendanceHandler) GetAllAttendance(c *gin.Context) {
	attendances, err := h.attendanceService.GetAllAttendance()
	if err != nil {
		InternalServerErrorResponse(c, "Failed to get all attendance records: "+err.Error())
		return
	}

	attendanceResponses := make([]response.AttendanceResponse, len(attendances))
	for i, attendance := range attendances {
		attendanceResponses[i] = h.convertToAttendanceResponse(attendance)
	}

	listResp := response.AttendanceListResponse{
		Attendances: attendanceResponses,
		Total:       len(attendanceResponses),
	}

	SuccessResponse(c, http.StatusOK, "All attendance records retrieved successfully", listResp)
}

// DeleteAttendance deletes an attendance record
func (h *AttendanceHandler) DeleteAttendance(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid attendance ID")
		return
	}

	err = h.attendanceService.DeleteAttendance(uint(id))
	if err != nil {
		if err == domain.ErrAttendanceNotFound {
			NotFoundResponse(c, "Attendance not found")
		} else {
			InternalServerErrorResponse(c, "Failed to delete attendance: "+err.Error())
		}
		return
	}

	SuccessResponse(c, http.StatusOK, "Attendance deleted successfully", nil)
}

// convertToAttendanceResponse converts domain Attendance to response AttendanceResponse
func (h *AttendanceHandler) convertToAttendanceResponse(attendance domain.Attendance) response.AttendanceResponse {
	breakResponses := make([]response.BreakResponse, len(attendance.Breaks))
	for i, breakItem := range attendance.Breaks {
		breakResponses[i] = h.convertToBreakResponse(&breakItem)
	}

	return response.AttendanceResponse{
		ID:             attendance.ID,
		UserID:         attendance.UserID,
		Date:           attendance.Date,
		CheckInTime:    attendance.CheckInTime,
		CheckOutTime:   attendance.CheckOutTime,
		TotalWorkHours: attendance.TotalWorkHours,
		Status:         attendance.GetStatus(),
		CreatedAt:      attendance.CreatedAt,
		UpdatedAt:      attendance.UpdatedAt,
		Breaks:         breakResponses,
	}
}

// convertToBreakResponse converts domain Break to response BreakResponse
func (h *AttendanceHandler) convertToBreakResponse(breakItem *domain.Break) response.BreakResponse {
	return response.BreakResponse{
		ID:           breakItem.ID,
		AttendanceID: breakItem.AttendanceID,
		StartTime:    breakItem.StartTime,
		EndTime:      breakItem.EndTime,
		Duration:     breakItem.Duration,
		Reason:       breakItem.Reason,
		CreatedAt:    breakItem.CreatedAt,
		UpdatedAt:    breakItem.UpdatedAt,
	}
}
