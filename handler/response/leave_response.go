package response

import (
	"hrm/domain"
	"time"
)

// LeaveResponse represents the response model for leave data
type LeaveResponse struct {
	ID           uint                 `json:"id"`
	UserID       uint                 `json:"user_id"`
	Type         domain.LeaveTypeName `json:"type"`
	Status       domain.LeaveStatus   `json:"status"`
	StartDate    time.Time            `json:"start_date"`
	EndDate      time.Time            `json:"end_date"`
	Days         float64              `json:"days"`
	Reason       string               `json:"reason"`
	Description  string               `json:"description"`
	ApprovedBy   *uint                `json:"approved_by"`
	ApprovedAt   *time.Time           `json:"approved_at"`
	RejectedBy   *uint                `json:"rejected_by"`
	RejectedAt   *time.Time           `json:"rejected_at"`
	RejectReason string               `json:"reject_reason"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
	User         *UserResponse        `json:"user,omitempty"`
	Approver     *UserResponse        `json:"approver,omitempty"`
	Rejecter     *UserResponse        `json:"rejecter,omitempty"`
}

// CreateLeaveResponse represents the response model for creating a leave
type CreateLeaveResponse struct {
	Leave   LeaveResponse `json:"leave"`
	Message string        `json:"message"`
}

// GetLeaveResponse represents the response model for getting a leave
type GetLeaveResponse struct {
	Leave LeaveResponse `json:"leave"`
}

// UpdateLeaveResponse represents the response model for updating a leave
type UpdateLeaveResponse struct {
	Leave   LeaveResponse `json:"leave"`
	Message string        `json:"message"`
}

// DeleteLeaveResponse represents the response model for deleting a leave
type DeleteLeaveResponse struct {
	Message string `json:"message"`
}

// ApproveLeaveResponse represents the response model for approving a leave
type ApproveLeaveResponse struct {
	Leave   LeaveResponse `json:"leave"`
	Message string        `json:"message"`
}

// RejectLeaveResponse represents the response model for rejecting a leave
type RejectLeaveResponse struct {
	Leave   LeaveResponse `json:"leave"`
	Message string        `json:"message"`
}

// CancelLeaveResponse represents the response model for cancelling a leave
type CancelLeaveResponse struct {
	Leave   LeaveResponse `json:"leave"`
	Message string        `json:"message"`
}

// ListLeavesResponse represents the response model for listing leaves
type ListLeavesResponse struct {
	Leaves []LeaveResponse `json:"leaves"`
	Total  int             `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}

// GetUserLeavesResponse represents the response model for getting user leaves
type GetUserLeavesResponse struct {
	Leaves []LeaveResponse `json:"leaves"`
	UserID uint            `json:"user_id"`
}

// GetUserLeaveBalanceResponse represents the response model for getting user leave balance
type GetUserLeaveBalanceResponse struct {
	UserID  uint                             `json:"user_id"`
	Year    int                              `json:"year"`
	Balance map[domain.LeaveTypeName]float64 `json:"balance"`
}

// ToLeaveResponse converts a domain Leave to LeaveResponse
func ToLeaveResponse(leave *domain.Leave) LeaveResponse {
	response := LeaveResponse{
		ID:           leave.ID,
		UserID:       leave.UserID,
		Type:         leave.Type,
		Status:       leave.Status,
		StartDate:    leave.StartDate,
		EndDate:      leave.EndDate,
		Days:         leave.Days,
		Reason:       leave.Reason,
		Description:  leave.Description,
		ApprovedBy:   leave.ApprovedBy,
		ApprovedAt:   leave.ApprovedAt,
		RejectedBy:   leave.RejectedBy,
		RejectedAt:   leave.RejectedAt,
		RejectReason: leave.RejectReason,
		CreatedAt:    leave.CreatedAt,
		UpdatedAt:    leave.UpdatedAt,
	}

	// Include user information if available
	if leave.User.ID != 0 {
		userResp := ToUserResponse(&leave.User)
		response.User = &userResp
	}

	// Include approver information if available
	if leave.Approver != nil && leave.Approver.ID != 0 {
		approverResp := ToUserResponse(leave.Approver)
		response.Approver = &approverResp
	}

	// Include rejecter information if available
	if leave.Rejecter != nil && leave.Rejecter.ID != 0 {
		rejecterResp := ToUserResponse(leave.Rejecter)
		response.Rejecter = &rejecterResp
	}

	return response
}

// ToLeaveResponseList converts a slice of domain Leaves to LeaveResponse slice
func ToLeaveResponseList(leaves []domain.Leave) []LeaveResponse {
	var responses []LeaveResponse
	for _, leave := range leaves {
		responses = append(responses, ToLeaveResponse(&leave))
	}
	return responses
}
