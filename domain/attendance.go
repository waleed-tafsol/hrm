package domain

import (
	"errors"
	"time"
)

// Attendance represents an employee's attendance record for a specific date
type Attendance struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	UserID         uint       `gorm:"not null" json:"user_id"`
	Date           time.Time  `gorm:"not null;type:date" json:"date"`
	CheckInTime    *time.Time `json:"check_in_time"`
	CheckOutTime   *time.Time `json:"check_out_time"`
	TotalWorkHours float64    `json:"total_work_hours"` // in hours

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    string    `json:"-"` // "present", "absent", "late", "early_leave", "completed"

	User   User    `gorm:"foreignKey:UserID" json:"-"`
	Breaks []Break `gorm:"foreignKey:AttendanceID" json:"breaks,omitempty"`
}

// AttendanceRepositoryInterface defines the contract for attendance data operations
type AttendanceRepositoryInterface interface {
	Create(attendance *Attendance) error
	GetByID(id uint) (*Attendance, error)
	GetByUserID(userID uint, date time.Time) (*Attendance, error)
	GetByUserIDAndDateRange(userID uint, startDate, endDate time.Time) ([]Attendance, error)
	GetByDate(date time.Time) ([]Attendance, error)
	Update(attendance *Attendance) error
	Delete(id uint) error
	GetAll() ([]Attendance, error)
	GetWithBreaks(id uint) (*Attendance, error)
	GetLastNByUserID(userID uint, limit int) ([]Attendance, error)
}

// AttendanceServiceInterface defines the contract for attendance business logic
type AttendanceServiceInterface interface {
	CreateAttendance(userID uint, date time.Time) (*Attendance, error)
	CheckIn(userID uint, date time.Time) (*Attendance, error)
	CheckOut(userID uint, date time.Time) (*Attendance, error)
	GetAttendanceByID(id uint) (*Attendance, error)
	GetUserAttendance(userID uint, date time.Time) (*Attendance, error)
	GetUserAttendanceRange(userID uint, startDate, endDate time.Time) ([]Attendance, error)
	GetAllAttendance() ([]Attendance, error)
	UpdateAttendance(attendance *Attendance) error
	DeleteAttendance(id uint) error
	CalculateWorkHours(attendance *Attendance) error
	GetLastNAttendanceByUserID(userID uint, limit int) ([]Attendance, error)
}

// Domain-specific errors for attendance operations
var (
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrInvalidDate        = errors.New("invalid date")
	ErrAttendanceNotFound = errors.New("attendance not found")
	ErrAlreadyCheckedIn   = errors.New("already checked in for this date")
	ErrAlreadyCheckedOut  = errors.New("already checked out for this date")
	ErrNotCheckedIn       = errors.New("not checked in yet")
)

// Validate checks if the attendance data is valid
func (a *Attendance) Validate() error {
	if a.UserID == 0 {
		return ErrInvalidUserID
	}
	if a.Date.IsZero() {
		return ErrInvalidDate
	}
	return nil
}

// CalculateWorkHours calculates the total work hours for the attendance
func (a *Attendance) CalculateWorkHours() {
	if a.CheckInTime == nil || a.CheckOutTime == nil {
		a.TotalWorkHours = 0
		return
	}

	duration := a.CheckOutTime.Sub(*a.CheckInTime)
	a.TotalWorkHours = duration.Hours()

	// Subtract break durations
	for _, breakItem := range a.Breaks {
		if breakItem.EndTime != nil {
			breakDuration := breakItem.EndTime.Sub(breakItem.StartTime)
			a.TotalWorkHours -= breakDuration.Hours()
		}
	}

	// Ensure work hours is not negative
	if a.TotalWorkHours < 0 {
		a.TotalWorkHours = 0
	}
}

// GetStatus returns the attendance status based on check-in/out times
func (a *Attendance) GetStatus() string {
	if a.CheckInTime == nil {
		return "absent"
	}
	if a.CheckOutTime == nil {
		return "present"
	}

	// You can add more logic here for "late" or "early_leave" based on your business rules
	return "completed"
}

// IsCheckedIn returns true if the user has checked in
func (a *Attendance) IsCheckedIn() bool {
	return a.CheckInTime != nil
}

// IsCheckedOut returns true if the user has checked out
func (a *Attendance) IsCheckedOut() bool {
	return a.CheckOutTime != nil
}

// CanCheckIn returns true if the user can check in (not already checked in)
func (a *Attendance) CanCheckIn() bool {
	return !a.IsCheckedIn()
}

// CanCheckOut returns true if the user can check out (checked in but not checked out)
func (a *Attendance) CanCheckOut() bool {
	return a.IsCheckedIn() && !a.IsCheckedOut()
}
