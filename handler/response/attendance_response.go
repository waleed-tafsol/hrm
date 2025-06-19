package response

import (
	"time"
)

// AttendanceResponse represents the response structure for attendance data
type AttendanceResponse struct {
	ID             uint            `json:"id"`
	UserID         uint            `json:"user_id"`
	Date           time.Time       `json:"date"`
	CheckInTime    *time.Time      `json:"check_in_time"`
	CheckOutTime   *time.Time      `json:"check_out_time"`
	TotalWorkHours float64         `json:"total_work_hours"`
	Status         string          `json:"status"` // "present", "absent", "late", "early_leave", "completed"
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Breaks         []BreakResponse `json:"breaks,omitempty"`
}

// AttendanceListResponse represents the response structure for a list of attendances
type AttendanceListResponse struct {
	Attendances []AttendanceResponse `json:"attendances"`
	Total       int                  `json:"total"`
}

// CheckInResponse represents the response structure for check-in operation
type CheckInResponse struct {
	Attendance  AttendanceResponse `json:"attendance"`
	Message     string             `json:"message"`
	CheckInTime time.Time          `json:"check_in_time"`
}

// CheckOutResponse represents the response structure for check-out operation
type CheckOutResponse struct {
	Attendance     AttendanceResponse `json:"attendance"`
	Message        string             `json:"message"`
	CheckOutTime   time.Time          `json:"check_out_time"`
	TotalWorkHours float64            `json:"total_work_hours"`
}
