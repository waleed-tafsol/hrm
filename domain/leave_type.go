package domain

import (
	"time"
)

// LeaveType represents a leave type definition
type LeaveType struct {
	ID                 uint      `json:"id" gorm:"primaryKey"`
	Type               string    `json:"type" gorm:"uniqueIndex;not null;type:varchar(20)"`
	Name               string    `json:"name" gorm:"not null;type:varchar(100)"`
	Description        string    `json:"description" gorm:"type:text"`
	DefaultDaysPerYear int       `json:"default_days_per_year" gorm:"default:0"`
	IsActive           bool      `json:"is_active" gorm:"default:true"`
	RequiresApproval   bool      `json:"requires_approval" gorm:"default:true"`
	Color              string    `json:"color" gorm:"default:'#007bff';type:varchar(7)"`
	Icon               string    `json:"icon" gorm:"type:varchar(50)"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	// Relationships
	Leaves []Leave `gorm:"foreignKey:Type;references:Type" json:"leaves,omitempty"`
}

// LeaveTypeRepositoryInterface defines the contract for leave type data operations
type LeaveTypeRepositoryInterface interface {
	Create(leaveType *LeaveType) error
	GetByID(id uint) (*LeaveType, error)
	GetByType(leaveType string) (*LeaveType, error)
	GetAll() ([]LeaveType, error)
	GetActive() ([]LeaveType, error)
	Update(leaveType *LeaveType) error
	Delete(id uint) error
	GetWithUsageStats() ([]LeaveType, error)
}

// LeaveTypeServiceInterface defines the contract for leave type business logic
type LeaveTypeServiceInterface interface {
	CreateLeaveType(leaveType *LeaveType) error
	GetLeaveTypeByID(id uint) (*LeaveType, error)
	GetLeaveTypeByType(leaveType string) (*LeaveType, error)
	GetAllLeaveTypes() ([]LeaveType, error)
	GetActiveLeaveTypes() ([]LeaveType, error)
	UpdateLeaveType(leaveType *LeaveType) error
	DeleteLeaveType(id uint) error
	GetLeaveTypesWithUsageStats() ([]LeaveType, error)
	ValidateLeaveType(leaveType string) error
}

// ValidateLeaveType checks if a leave type is valid
func (ltd *LeaveType) ValidateLeaveType() error {
	validTypes := []string{
		"sick",
		"vacation",
		"personal",
		"maternity",
		"paternity",
		"other",
	}

	for _, validType := range validTypes {
		if ltd.Type == validType {
			return nil
		}
	}
	return ErrInvalidLeaveType
}

// IsActive returns true if the leave type is active
func (ltd *LeaveType) GetIsActive() bool {
	return ltd.IsActive
}

// RequiresApproval returns true if the leave type requires approval
func (ltd *LeaveType) GetRequiresApproval() bool {
	return ltd.RequiresApproval
}

// GetDefaultDaysPerYear returns the default days per year for this leave type
func (ltd *LeaveType) GetDefaultDaysPerYear() int {
	return ltd.DefaultDaysPerYear
}

func (LeaveType) TableName() string {
	return "leave_types"
}
