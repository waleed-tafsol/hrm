package request

import (
	"time"
)

// AttendanceRequest represents the request structure for attendance operations
type AttendanceRequest struct {
	Date time.Time `json:"date" binding:"required"`
}

// CheckInRequest represents the request structure for check-in
type CheckInRequest struct {
	Date time.Time `json:"date" binding:"required"`
}

// CheckOutRequest represents the request structure for check-out
type CheckOutRequest struct {
	Date time.Time `json:"date" binding:"required"`
}

// AttendanceRangeRequest represents the request structure for getting attendance range
type AttendanceRangeRequest struct {
	StartDate time.Time `json:"start_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}

// AttendanceUpdateRequest represents the request structure for updating attendance
type AttendanceUpdateRequest struct {
	CheckInTime  *time.Time `json:"check_in_time"`
	CheckOutTime *time.Time `json:"check_out_time"`
}
