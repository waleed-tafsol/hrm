package domain

import (
	"errors"
	"time"
)

// Break represents a break period during attendance
type Break struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	AttendanceID uint       `gorm:"not null" json:"attendance_id"`
	StartTime    time.Time  `gorm:"not null" json:"start_time"`
	EndTime      *time.Time `json:"end_time"`
	Duration     float64    `json:"duration"` // in minutes
	Reason       string     `json:"reason"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Attendance Attendance `gorm:"foreignKey:AttendanceID" json:"-"`
}

// BreakRepositoryInterface defines the contract for break data operations
type BreakRepositoryInterface interface {
	Create(breakItem *Break) error
	GetByID(id uint) (*Break, error)
	GetByAttendanceID(attendanceID uint) ([]Break, error)
	GetActiveBreakByAttendanceID(attendanceID uint) (*Break, error)
	Update(breakItem *Break) error
	Delete(id uint) error
	GetAll() ([]Break, error)
	GetBreaksByDateRange(startDate, endDate time.Time) ([]Break, error)
}

// BreakServiceInterface defines the contract for break business logic
type BreakServiceInterface interface {
	CreateBreak(attendanceID uint, startTime time.Time, reason string) (*Break, error)
	GetBreakByID(id uint) (*Break, error)
	GetBreaksByAttendanceID(attendanceID uint) ([]Break, error)
	GetAllBreaks() ([]Break, error)
	UpdateBreak(breakItem *Break) error
	DeleteBreak(id uint) error
	EndBreak(breakID uint, endTime time.Time) error
	CalculateBreakDuration(breakItem *Break) error
}

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

// Domain-specific errors for break operations
var (
	ErrInvalidAttendanceID = errors.New("invalid attendance ID")
	ErrInvalidBreakTime    = errors.New("invalid break time")
	ErrBreakInProgress     = errors.New("break already in progress")
	ErrBreakNotFound       = errors.New("break not found")
	ErrBreakAlreadyEnded   = errors.New("break already ended")
)

// Validate checks if the break data is valid
func (b *Break) Validate() error {
	if b.AttendanceID == 0 {
		return ErrInvalidAttendanceID
	}
	if b.StartTime.IsZero() {
		return ErrInvalidBreakTime
	}
	return nil
}

// CalculateDuration calculates the duration of the break in minutes
func (b *Break) CalculateDuration() {
	if b.EndTime == nil {
		b.Duration = 0
		return
	}

	duration := b.EndTime.Sub(b.StartTime)
	b.Duration = duration.Minutes()
}

// IsEnded returns true if the break has ended
func (b *Break) IsEnded() bool {
	return b.EndTime != nil
}

// IsInProgress returns true if the break is currently in progress
func (b *Break) IsInProgress() bool {
	return !b.IsEnded()
}

// CanEnd returns true if the break can be ended (not already ended)
func (b *Break) CanEnd() bool {
	return b.IsInProgress()
}

// GetDurationInHours returns the break duration in hours
func (b *Break) GetDurationInHours() float64 {
	return b.Duration / 60.0
}

// GetDurationInMinutes returns the break duration in minutes
func (b *Break) GetDurationInMinutes() float64 {
	return b.Duration
}

// GetDurationInSeconds returns the break duration in seconds
func (b *Break) GetDurationInSeconds() float64 {
	return b.Duration * 60.0
}
