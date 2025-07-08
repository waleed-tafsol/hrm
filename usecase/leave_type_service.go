package usecase

import (
	"hrm/domain"
)

// LeaveTypeService implements LeaveTypeServiceInterface
type LeaveTypeService struct {
	leaveTypeRepo domain.LeaveTypeRepositoryInterface
}

// NewLeaveTypeService creates a new instance of LeaveTypeService
func NewLeaveTypeService(leaveTypeRepo domain.LeaveTypeRepositoryInterface) domain.LeaveTypeServiceInterface {
	return &LeaveTypeService{
		leaveTypeRepo: leaveTypeRepo,
	}
}

// CreateLeaveType creates a new leave type
func (s *LeaveTypeService) CreateLeaveType(leaveType *domain.LeaveType) error {
	// Validate the leave type
	if err := leaveType.ValidateLeaveType(); err != nil {
		return err
	}

	// Check if leave type already exists
	existing, err := s.leaveTypeRepo.GetByType(leaveType.Type)
	if err == nil && existing != nil {
		return domain.ErrInvalidLeaveType
	}

	return s.leaveTypeRepo.Create(leaveType)
}

// GetLeaveTypeByID retrieves a leave type by ID
func (s *LeaveTypeService) GetLeaveTypeByID(id uint) (*domain.LeaveType, error) {
	return s.leaveTypeRepo.GetByID(id)
}

// GetLeaveTypeByType retrieves a leave type by type string
func (s *LeaveTypeService) GetLeaveTypeByType(leaveType string) (*domain.LeaveType, error) {
	return s.leaveTypeRepo.GetByType(leaveType)
}

// GetAllLeaveTypes retrieves all leave types
func (s *LeaveTypeService) GetAllLeaveTypes() ([]domain.LeaveType, error) {
	return s.leaveTypeRepo.GetAll()
}

// GetActiveLeaveTypes retrieves all active leave types
func (s *LeaveTypeService) GetActiveLeaveTypes() ([]domain.LeaveType, error) {
	return s.leaveTypeRepo.GetActive()
}

// UpdateLeaveType updates a leave type
func (s *LeaveTypeService) UpdateLeaveType(leaveType *domain.LeaveType) error {
	// Validate the leave type
	if err := leaveType.ValidateLeaveType(); err != nil {
		return err
	}

	return s.leaveTypeRepo.Update(leaveType)
}

// DeleteLeaveType deletes a leave type by ID
func (s *LeaveTypeService) DeleteLeaveType(id uint) error {
	// Check if leave type exists
	_, err := s.leaveTypeRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Check if leave type is being used
	leaveTypesWithStats, err := s.leaveTypeRepo.GetWithUsageStats()
	if err != nil {
		return err
	}

	for _, lt := range leaveTypesWithStats {
		if lt.ID == id && len(lt.Leaves) > 0 {
			return domain.ErrInvalidLeaveType // Cannot delete leave type that is being used
		}
	}

	return s.leaveTypeRepo.Delete(id)
}

// GetLeaveTypesWithUsageStats retrieves leave types with usage statistics
func (s *LeaveTypeService) GetLeaveTypesWithUsageStats() ([]domain.LeaveType, error) {
	return s.leaveTypeRepo.GetWithUsageStats()
}

// ValidateLeaveType validates if a leave type string is valid
func (s *LeaveTypeService) ValidateLeaveType(leaveType string) error {
	validTypes := []string{
		"sick",
		"vacation",
		"personal",
		"maternity",
		"paternity",
		"other",
	}

	for _, validType := range validTypes {
		if leaveType == validType {
			return nil
		}
	}
	return domain.ErrInvalidLeaveType
}
