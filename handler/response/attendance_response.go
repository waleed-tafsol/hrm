package response

import (
	"hrm/domain"
	"time"
)

// AttendanceResponse represents the response structure for attendance data
type AttendanceResponse struct {
	ID             uint       `json:"id"`
	UserID         uint       `json:"user_id"`
	Date           time.Time  `json:"date"`
	CheckInTime    *time.Time `json:"check_in_time"`
	CheckOutTime   *time.Time `json:"check_out_time"`
	TotalWorkHours float64    `json:"total_work_hours"`
	Status         string     `json:"status"` // "present", "absent", "late", "early_leave", "completed"
	//CreatedAt      time.Time       `json:"created_at"`
	//UpdatedAt      time.Time       `json:"updated_at"`
	Breaks []BreakResponse `json:"breaks"`
}

// AttendanceListResponse represents the response structure for a list of attendances
type AttendanceListResponse struct {
	Attendances []AttendanceResponse `json:"attendances"`
	Total       int                  `json:"total"`
}

// CheckInResponse represents the response structure for check-in operation
type CheckInResponse struct {
	Attendance AttendanceResponse `json:"attendance"`
}

// CheckOutResponse represents the response structure for check-out operation
type CheckOutResponse struct {
	Attendance AttendanceResponse `json:"attendance"`
}

// ToAttendanceResponse converts a domain Attendance to AttendanceResponse
func ToAttendanceResponse(attendance *domain.Attendance) AttendanceResponse {
	breaks := make([]BreakResponse, len(attendance.Breaks))
	for i, breakItem := range attendance.Breaks {
		breaks[i] = ToBreakResponse(&breakItem)
	}

	return AttendanceResponse{
		ID:             attendance.ID,
		UserID:         attendance.UserID,
		Date:           attendance.Date,
		CheckInTime:    attendance.CheckInTime,
		CheckOutTime:   attendance.CheckOutTime,
		TotalWorkHours: attendance.TotalWorkHours,
		Status:         attendance.GetStatus(),
		Breaks:         breaks,
	}
}

// ToAttendanceResponseList converts a slice of domain Attendances to AttendanceResponse slice
func ToAttendanceResponseList(attendances []domain.Attendance) []AttendanceResponse {
	var responses []AttendanceResponse
	for _, attendance := range attendances {
		responses = append(responses, ToAttendanceResponse(&attendance))
	}
	return responses
}
