package request

import (
	"hrm/domain"
	"time"
)

// CreateLeaveRequest represents the request model for creating a leave
type CreateLeaveRequest struct {
	Type        domain.LeaveType `json:"type" binding:"required"`
	StartDate   time.Time        `json:"start_date" binding:"required"`
	EndDate     time.Time        `json:"end_date" binding:"required"`
	Reason      string           `json:"reason" binding:"required"`
	Description string           `json:"description"`
}

// UpdateLeaveRequest represents the request model for updating a leave
type UpdateLeaveRequest struct {
	Type        domain.LeaveType `json:"type" binding:"required"`
	StartDate   time.Time        `json:"start_date" binding:"required"`
	EndDate     time.Time        `json:"end_date" binding:"required"`
	Reason      string           `json:"reason" binding:"required"`
	Description string           `json:"description"`
}

// ApproveLeaveRequest represents the request model for approving a leave
type ApproveLeaveRequest struct {
	LeaveID uint `uri:"id" binding:"required"`
}

// RejectLeaveRequest represents the request model for rejecting a leave
type RejectLeaveRequest struct {
	LeaveID      uint   `uri:"id" binding:"required"`
	RejectReason string `json:"reject_reason" binding:"required"`
}

// CancelLeaveRequest represents the request model for cancelling a leave
type CancelLeaveRequest struct {
	LeaveID uint `uri:"id" binding:"required"`
}

// GetLeaveByIDRequest represents the request model for getting a leave by ID
type GetLeaveByIDRequest struct {
	LeaveID uint `uri:"id" binding:"required"`
}

// DeleteLeaveRequest represents the request model for deleting a leave
type DeleteLeaveRequest struct {
	LeaveID uint `uri:"id" binding:"required"`
}

// GetUserLeavesRequest represents the request model for getting user leaves
type GetUserLeavesRequest struct {
	UserID uint `uri:"user_id" binding:"required"`
}

// GetUserLeavesByDateRangeRequest represents the request model for getting user leaves by date range
type GetUserLeavesByDateRangeRequest struct {
	UserID    uint      `uri:"user_id" binding:"required"`
	StartDate time.Time `form:"start_date" binding:"required"`
	EndDate   time.Time `form:"end_date" binding:"required"`
}

// GetUserLeaveBalanceRequest represents the request model for getting user leave balance
type GetUserLeaveBalanceRequest struct {
	UserID uint `uri:"user_id" binding:"required"`
	Year   int  `form:"year" binding:"required,min=2000,max=2100"`
}

// GetLeavesByStatusRequest represents the request model for getting leaves by status
type GetLeavesByStatusRequest struct {
	Status domain.LeaveStatus `form:"status" binding:"required"`
}

// GetLeavesByTypeRequest represents the request model for getting leaves by type
type GetLeavesByTypeRequest struct {
	Type domain.LeaveType `form:"type" binding:"required"`
}

// ListLeavesRequest represents the request model for listing leaves with pagination
type ListLeavesRequest struct {
	Limit  int `form:"limit" binding:"min=1,max=100"`
	Offset int `form:"offset" binding:"min=0"`
}
