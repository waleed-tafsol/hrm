package repository

import (
	"hrm/domain"

	"gorm.io/gorm"
)

// LeaveTypeRepository implements LeaveTypeRepositoryInterface
type LeaveTypeRepository struct {
	db *gorm.DB
}

// NewLeaveTypeRepository creates a new instance of LeaveTypeRepository
func NewLeaveTypeRepository(db *gorm.DB) domain.LeaveTypeRepositoryInterface {
	return &LeaveTypeRepository{db: db}
}

// Create creates a new leave type
func (r *LeaveTypeRepository) Create(leaveType *domain.LeaveType) error {
	return r.db.Create(leaveType).Error
}

// GetByID retrieves a leave type by ID
func (r *LeaveTypeRepository) GetByID(id uint) (*domain.LeaveType, error) {
	var leaveType domain.LeaveType
	err := r.db.First(&leaveType, id).Error
	if err != nil {
		return nil, err
	}
	return &leaveType, nil
}

// GetByType retrieves a leave type by type string
func (r *LeaveTypeRepository) GetByType(leaveType string) (*domain.LeaveType, error) {
	var lt domain.LeaveType
	err := r.db.Where("type = ?", leaveType).First(&lt).Error
	if err != nil {
		return nil, err
	}
	return &lt, nil
}

// GetAll retrieves all leave types
func (r *LeaveTypeRepository) GetAll() ([]domain.LeaveType, error) {
	var leaveTypes []domain.LeaveType
	err := r.db.Find(&leaveTypes).Error
	return leaveTypes, err
}

// GetActive retrieves all active leave types
func (r *LeaveTypeRepository) GetActive() ([]domain.LeaveType, error) {
	var leaveTypes []domain.LeaveType
	err := r.db.Where("is_active = ?", true).Find(&leaveTypes).Error
	return leaveTypes, err
}

// Update updates a leave type
func (r *LeaveTypeRepository) Update(leaveType *domain.LeaveType) error {
	return r.db.Save(leaveType).Error
}

// Delete deletes a leave type by ID
func (r *LeaveTypeRepository) Delete(id uint) error {
	return r.db.Delete(&domain.LeaveType{}, id).Error
}

// GetWithUsageStats retrieves leave types with usage statistics
func (r *LeaveTypeRepository) GetWithUsageStats() ([]domain.LeaveType, error) {
	var leaveTypes []domain.LeaveType
	err := r.db.Preload("Leaves").Find(&leaveTypes).Error
	return leaveTypes, err
}
