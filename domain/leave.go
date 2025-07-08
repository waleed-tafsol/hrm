package domain

import (
	"errors"
	"time"
)

// LeaveTypeName represents the type of leave
type LeaveTypeName string

const (
	LeaveTypeSick      LeaveTypeName = "sick"
	LeaveTypeVacation  LeaveTypeName = "vacation"
	LeaveTypePersonal  LeaveTypeName = "personal"
	LeaveTypeMaternity LeaveTypeName = "maternity"
	LeaveTypePaternity LeaveTypeName = "paternity"
	LeaveTypeOther     LeaveTypeName = "other"
)

// LeaveStatus represents the status of a leave request
type LeaveStatus string

const (
	LeaveStatusPending   LeaveStatus = "pending"
	LeaveStatusApproved  LeaveStatus = "approved"
	LeaveStatusRejected  LeaveStatus = "rejected"
	LeaveStatusCancelled LeaveStatus = "cancelled"
)

// Leave represents an employee's leave request
type Leave struct {
	ID           uint          `json:"id" gorm:"primaryKey"`
	UserID       uint          `json:"user_id" gorm:"not null"`
	Type         LeaveTypeName `json:"type" gorm:"not null;type:varchar(20)"`
	Status       LeaveStatus   `json:"status" gorm:"not null;type:varchar(20);default:'pending'"`
	StartDate    time.Time     `json:"start_date" gorm:"not null;type:date"`
	EndDate      time.Time     `json:"end_date" gorm:"not null;type:date"`
	Days         float64       `json:"days" gorm:"not null"` // Number of days (can be fractional)
	Reason       string        `json:"reason" gorm:"not null;type:text"`
	Description  string        `json:"description" gorm:"type:text"`
	ApprovedBy   *uint         `json:"approved_by" gorm:"index"`
	ApprovedAt   *time.Time    `json:"approved_at"`
	RejectedBy   *uint         `json:"rejected_by" gorm:"index"`
	RejectedAt   *time.Time    `json:"rejected_at"`
	RejectReason string        `json:"reject_reason" gorm:"type:text"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User     User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Approver *User `gorm:"foreignKey:ApprovedBy" json:"approver,omitempty"`
	Rejecter *User `gorm:"foreignKey:RejectedBy" json:"rejecter,omitempty"`
}

// LeaveRepositoryInterface defines the contract for leave data operations
type LeaveRepositoryInterface interface {
	Create(leave *Leave) error
	GetByID(id uint) (*Leave, error)
	GetByUserID(userID uint) ([]Leave, error)
	GetByUserIDAndDateRange(userID uint, startDate, endDate time.Time) ([]Leave, error)
	GetByStatus(status LeaveStatus) ([]Leave, error)
	GetByType(leaveType LeaveTypeName) ([]Leave, error)
	GetPendingLeaves() ([]Leave, error)
	Update(leave *Leave) error
	Delete(id uint) error
	GetAll() ([]Leave, error)
	GetWithUser(id uint) (*Leave, error)
	GetUserLeaveBalance(userID uint, year int) (map[LeaveTypeName]float64, error)
}

// LeaveServiceInterface defines the contract for leave business logic
type LeaveServiceInterface interface {
	CreateLeave(userID uint, leave *Leave) error
	GetLeaveByID(id uint) (*Leave, error)
	GetUserLeaves(userID uint) ([]Leave, error)
	GetUserLeavesByDateRange(userID uint, startDate, endDate time.Time) ([]Leave, error)
	GetAllLeaves() ([]Leave, error)
	GetPendingLeaves() ([]Leave, error)
	UpdateLeave(leave *Leave) error
	DeleteLeave(id uint) error
	ApproveLeave(leaveID uint, approverID uint) error
	RejectLeave(leaveID uint, rejecterID uint, reason string) error
	CancelLeave(leaveID uint, userID uint) error
	GetUserLeaveBalance(userID uint, year int) (map[LeaveTypeName]float64, error)
	CalculateLeaveDays(startDate, endDate time.Time) float64
}

// Domain-specific errors for leave operations
var (
	ErrInvalidLeaveType          = errors.New("invalid leave type")
	ErrInvalidLeaveStatus        = errors.New("invalid leave status")
	ErrInvalidDateRange          = errors.New("start date must be before or equal to end date")
	ErrLeaveNotFound             = errors.New("leave not found")
	ErrLeaveAlreadyApproved      = errors.New("leave is already approved")
	ErrLeaveAlreadyRejected      = errors.New("leave is already rejected")
	ErrLeaveAlreadyCancelled     = errors.New("leave is already cancelled")
	ErrCannotCancelApprovedLeave = errors.New("cannot cancel an approved leave")
	ErrInsufficientLeaveBalance  = errors.New("insufficient leave balance")
	ErrLeaveDateInPast           = errors.New("leave date cannot be in the past")
	ErrLeaveOverlap              = errors.New("leave dates overlap with existing leave")
)

// Validate checks if the leave data is valid
func (l *Leave) Validate() error {
	if l.UserID == 0 {
		return ErrInvalidUserID
	}

	if !l.isValidLeaveType() {
		return ErrInvalidLeaveType
	}

	if !l.isValidLeaveStatus() {
		return ErrInvalidLeaveStatus
	}

	if l.StartDate.After(l.EndDate) {
		return ErrInvalidDateRange
	}

	if l.StartDate.Before(time.Now().Truncate(24 * time.Hour)) {
		return ErrLeaveDateInPast
	}

	if l.Reason == "" {
		return errors.New("reason is required")
	}

	return nil
}

// isValidLeaveType checks if the leave type is valid
func (l *Leave) isValidLeaveType() bool {
	validTypes := []LeaveTypeName{
		LeaveTypeSick,
		LeaveTypeVacation,
		LeaveTypePersonal,
		LeaveTypeMaternity,
		LeaveTypePaternity,
		LeaveTypeOther,
	}

	for _, validType := range validTypes {
		if l.Type == validType {
			return true
		}
	}
	return false
}

// isValidLeaveStatus checks if the leave status is valid
func (l *Leave) isValidLeaveStatus() bool {
	validStatuses := []LeaveStatus{
		LeaveStatusPending,
		LeaveStatusApproved,
		LeaveStatusRejected,
		LeaveStatusCancelled,
	}

	for _, validStatus := range validStatuses {
		if l.Status == validStatus {
			return true
		}
	}
	return false
}

// CanApprove returns true if the leave can be approved
func (l *Leave) CanApprove() bool {
	return l.Status == LeaveStatusPending
}

// CanReject returns true if the leave can be rejected
func (l *Leave) CanReject() bool {
	return l.Status == LeaveStatusPending
}

// CanCancel returns true if the leave can be cancelled
func (l *Leave) CanCancel() bool {
	return l.Status == LeaveStatusPending || l.Status == LeaveStatusApproved
}

// IsPending returns true if the leave is pending approval
func (l *Leave) IsPending() bool {
	return l.Status == LeaveStatusPending
}

// IsApproved returns true if the leave is approved
func (l *Leave) IsApproved() bool {
	return l.Status == LeaveStatusApproved
}

// IsRejected returns true if the leave is rejected
func (l *Leave) IsRejected() bool {
	return l.Status == LeaveStatusRejected
}

// IsCancelled returns true if the leave is cancelled
func (l *Leave) IsCancelled() bool {
	return l.Status == LeaveStatusCancelled
}
