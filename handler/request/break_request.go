package request

import (
	"time"
)

// BreakRequest represents the request structure for break operations
type BreakRequest struct {
	AttendanceID uint      `json:"attendance_id" binding:"required"`
	StartTime    time.Time `json:"start_time" binding:"required"`
	Reason       string    `json:"reason"`
}

// EndBreakRequest represents the request structure for ending a break
type EndBreakRequest struct {
	BreakID uint      `json:"break_id" binding:"required"`
	EndTime time.Time `json:"end_time" binding:"required"`
}

// BreakUpdateRequest represents the request structure for updating break details
type BreakUpdateRequest struct {
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`
	Reason    *string    `json:"reason"`
}

// BreakRangeRequest represents the request structure for getting breaks within a date range
type BreakRangeRequest struct {
	UserID    uint      `json:"user_id" binding:"required"`
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}
