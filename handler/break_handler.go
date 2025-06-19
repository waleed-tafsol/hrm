package handler

import (
	"net/http"
	"strconv"

	"hrm/domain"
	"hrm/handler/request"
	"hrm/handler/response"

	"github.com/gin-gonic/gin"
)

// BreakHandler handles HTTP requests related to break operations
type BreakHandler struct {
	breakService domain.BreakServiceInterface
}

// NewBreakHandler creates a new instance of BreakHandler
func NewBreakHandler(breakService domain.BreakServiceInterface) *BreakHandler {
	return &BreakHandler{
		breakService: breakService,
	}
}

// AddBreak adds a new break to an attendance record
func (h *BreakHandler) AddBreak(c *gin.Context) {
	var req request.BreakRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request data: "+err.Error())
		return
	}

	breakItem, err := h.breakService.CreateBreak(req.AttendanceID, req.StartTime, req.Reason)
	if err != nil {
		if err == domain.ErrAttendanceNotFound {
			NotFoundResponse(c, "Attendance not found")
		} else if err == domain.ErrBreakInProgress {
			c.JSON(http.StatusConflict, Response{
				Success: false,
				Message: "Break already in progress",
			})
		} else {
			InternalServerErrorResponse(c, "Failed to add break: "+err.Error())
		}
		return
	}

	breakResp := h.convertToBreakResponse(*breakItem)
	addBreakResp := response.BreakAddResponse{
		Break:   breakResp,
		Message: "Break added successfully",
	}

	SuccessResponse(c, http.StatusCreated, "Break added successfully", addBreakResp)
}

// EndBreak ends an existing break
func (h *BreakHandler) EndBreak(c *gin.Context) {
	var req request.EndBreakRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request data: "+err.Error())
		return
	}

	err := h.breakService.EndBreak(req.BreakID, req.EndTime)
	if err != nil {
		if err == domain.ErrBreakNotFound {
			NotFoundResponse(c, "Break not found")
		} else if err == domain.ErrBreakAlreadyEnded {
			c.JSON(http.StatusConflict, Response{
				Success: false,
				Message: "Break already ended",
			})
		} else if err == domain.ErrInvalidBreakTime {
			BadRequestResponse(c, "Invalid break time")
		} else {
			InternalServerErrorResponse(c, "Failed to end break: "+err.Error())
		}
		return
	}

	// Get the updated break record
	breakItem, err := h.breakService.GetBreakByID(req.BreakID)
	if err != nil {
		InternalServerErrorResponse(c, "Failed to get updated break: "+err.Error())
		return
	}

	breakResp := h.convertToBreakResponse(*breakItem)
	endBreakResp := response.BreakEndResponse{
		Break:    breakResp,
		Message:  "Break ended successfully",
		Duration: breakItem.Duration,
	}

	SuccessResponse(c, http.StatusOK, "Break ended successfully", endBreakResp)
}

// GetBreakByID retrieves a break record by ID
func (h *BreakHandler) GetBreakByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid break ID")
		return
	}

	breakItem, err := h.breakService.GetBreakByID(uint(id))
	if err != nil {
		if err == domain.ErrBreakNotFound {
			NotFoundResponse(c, "Break not found")
		} else {
			InternalServerErrorResponse(c, "Failed to get break: "+err.Error())
		}
		return
	}

	breakResp := h.convertToBreakResponse(*breakItem)
	SuccessResponse(c, http.StatusOK, "Break retrieved successfully", breakResp)
}

// GetBreaksByAttendanceID retrieves all breaks for an attendance
func (h *BreakHandler) GetBreaksByAttendanceID(c *gin.Context) {
	attendanceIDStr := c.Param("attendance_id")
	attendanceID, err := strconv.ParseUint(attendanceIDStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid attendance ID")
		return
	}

	breaks, err := h.breakService.GetBreaksByAttendanceID(uint(attendanceID))
	if err != nil {
		if err == domain.ErrAttendanceNotFound {
			NotFoundResponse(c, "Attendance not found")
		} else {
			InternalServerErrorResponse(c, "Failed to get breaks: "+err.Error())
		}
		return
	}

	breakResponses := make([]response.BreakResponse, len(breaks))
	for i, breakItem := range breaks {
		breakResponses[i] = h.convertToBreakResponse(breakItem)
	}

	listResp := response.BreakListResponse{
		Breaks: breakResponses,
		Total:  len(breakResponses),
	}

	SuccessResponse(c, http.StatusOK, "Breaks retrieved successfully", listResp)
}

// GetAllBreaks retrieves all break records
func (h *BreakHandler) GetAllBreaks(c *gin.Context) {
	breaks, err := h.breakService.GetAllBreaks()
	if err != nil {
		InternalServerErrorResponse(c, "Failed to get all breaks: "+err.Error())
		return
	}

	breakResponses := make([]response.BreakResponse, len(breaks))
	for i, breakItem := range breaks {
		breakResponses[i] = h.convertToBreakResponse(breakItem)
	}

	listResp := response.BreakListResponse{
		Breaks: breakResponses,
		Total:  len(breakResponses),
	}

	SuccessResponse(c, http.StatusOK, "All breaks retrieved successfully", listResp)
}

// DeleteBreak deletes a break record
func (h *BreakHandler) DeleteBreak(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		BadRequestResponse(c, "Invalid break ID")
		return
	}

	err = h.breakService.DeleteBreak(uint(id))
	if err != nil {
		if err == domain.ErrBreakNotFound {
			NotFoundResponse(c, "Break not found")
		} else {
			InternalServerErrorResponse(c, "Failed to delete break: "+err.Error())
		}
		return
	}

	SuccessResponse(c, http.StatusOK, "Break deleted successfully", nil)
}

// convertToBreakResponse converts domain Break to response BreakResponse
func (h *BreakHandler) convertToBreakResponse(breakItem domain.Break) response.BreakResponse {
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
