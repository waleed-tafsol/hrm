package response

import (
	"time"
)

// BreakResponse represents the response structure for break data
type BreakResponse struct {
	ID           uint       `json:"id"`
	AttendanceID uint       `json:"attendance_id"`
	StartTime    time.Time  `json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	Duration     float64    `json:"duration"`
	Reason       string     `json:"reason"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// BreakAddResponse represents the response structure for break operations
type BreakAddResponse struct {
	Break   BreakResponse `json:"break"`
	Message string        `json:"message"`
}

// BreakEndResponse represents the response structure for ending a break
type BreakEndResponse struct {
	Break    BreakResponse `json:"break"`
	Message  string        `json:"message"`
	Duration float64       `json:"duration"`
}

// BreakListResponse represents the response structure for a list of breaks
type BreakListResponse struct {
	Breaks []BreakResponse `json:"breaks"`
	Total  int             `json:"total"`
}

// BreakUpdateResponse represents the response structure for updating break details
type BreakUpdateResponse struct {
	Break   BreakResponse `json:"break"`
	Message string        `json:"message"`
}

// BreakDetailResponse represents the response structure for a single break with additional details
type BreakDetailResponse struct {
	Break         BreakResponse `json:"break"`
	AttendanceID  uint          `json:"attendance_id"`
	UserID        uint          `json:"user_id"`
	Date          time.Time     `json:"date"`
	TotalDuration float64       `json:"total_duration"`
}
