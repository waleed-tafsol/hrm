package usecase

import (
	"hrm/domain"
	"log"
	"time"
)

// LeaveServiceImpl implements the LeaveServiceInterface
type LeaveServiceImpl struct {
	leaveRepo domain.LeaveRepositoryInterface
	userRepo  domain.UserRepositoryInterface
}

// NewLeaveService creates and returns a new LeaveServiceImpl instance
func NewLeaveService(leaveRepo domain.LeaveRepositoryInterface, userRepo domain.UserRepositoryInterface) domain.LeaveServiceInterface {
	return &LeaveServiceImpl{
		leaveRepo: leaveRepo,
		userRepo:  userRepo,
	}
}

// CreateLeave creates a new leave request
func (s *LeaveServiceImpl) CreateLeave(userID uint, leave *domain.Leave) error {
	// Set the user ID and initial status
	leave.UserID = userID
	leave.Status = domain.LeaveStatusPending

	// Calculate the number of days
	leave.Days = s.CalculateLeaveDays(leave.StartDate, leave.EndDate)

	// Validate the leave
	if err := leave.Validate(); err != nil {
		return err
	}

	// Check for overlapping leaves
	overlappingLeaves, err := s.leaveRepo.GetByUserIDAndDateRange(userID, leave.StartDate, leave.EndDate)
	if err != nil {
		log.Printf("Error checking for overlapping leaves: %v", err)
		return err
	}

	// Filter out leaves that are cancelled or rejected
	for _, existingLeave := range overlappingLeaves {
		if existingLeave.Status != domain.LeaveStatusCancelled && existingLeave.Status != domain.LeaveStatusRejected {
			return domain.ErrLeaveOverlap
		}
	}

	// Create the leave
	return s.leaveRepo.Create(leave)
}

// GetLeaveByID retrieves a leave by its ID
func (s *LeaveServiceImpl) GetLeaveByID(id uint) (*domain.Leave, error) {
	return s.leaveRepo.GetByID(id)
}

// GetUserLeaves retrieves all leaves for a specific user
func (s *LeaveServiceImpl) GetUserLeaves(userID uint) ([]domain.Leave, error) {
	return s.leaveRepo.GetByUserID(userID)
}

// GetUserLeavesByDateRange retrieves leaves for a user within a specific date range
func (s *LeaveServiceImpl) GetUserLeavesByDateRange(userID uint, startDate, endDate time.Time) ([]domain.Leave, error) {
	return s.leaveRepo.GetByUserIDAndDateRange(userID, startDate, endDate)
}

// GetAllLeaves retrieves all leaves
func (s *LeaveServiceImpl) GetAllLeaves() ([]domain.Leave, error) {
	return s.leaveRepo.GetAll()
}

// GetPendingLeaves retrieves all pending leave requests
func (s *LeaveServiceImpl) GetPendingLeaves() ([]domain.Leave, error) {
	return s.leaveRepo.GetPendingLeaves()
}

// UpdateLeave updates an existing leave
func (s *LeaveServiceImpl) UpdateLeave(leave *domain.Leave) error {
	// Validate the leave
	if err := leave.Validate(); err != nil {
		return err
	}

	// Recalculate days if dates changed
	leave.Days = s.CalculateLeaveDays(leave.StartDate, leave.EndDate)

	return s.leaveRepo.Update(leave)
}

// DeleteLeave deletes a leave by ID
func (s *LeaveServiceImpl) DeleteLeave(id uint) error {
	return s.leaveRepo.Delete(id)
}

// ApproveLeave approves a leave request
func (s *LeaveServiceImpl) ApproveLeave(leaveID uint, approverID uint) error {
	leave, err := s.leaveRepo.GetByID(leaveID)
	if err != nil {
		return err
	}

	if !leave.CanApprove() {
		return domain.ErrLeaveAlreadyApproved
	}

	// Verify the approver exists
	_, err = s.userRepo.GetByID(approverID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	// Update leave status
	now := time.Now()
	leave.Status = domain.LeaveStatusApproved
	leave.ApprovedBy = &approverID
	leave.ApprovedAt = &now

	return s.leaveRepo.Update(leave)
}

// RejectLeave rejects a leave request
func (s *LeaveServiceImpl) RejectLeave(leaveID uint, rejecterID uint, reason string) error {
	leave, err := s.leaveRepo.GetByID(leaveID)
	if err != nil {
		return err
	}

	if !leave.CanReject() {
		return domain.ErrLeaveAlreadyRejected
	}

	// Verify the rejecter exists
	_, err = s.userRepo.GetByID(rejecterID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	// Update leave status
	now := time.Now()
	leave.Status = domain.LeaveStatusRejected
	leave.RejectedBy = &rejecterID
	leave.RejectedAt = &now
	leave.RejectReason = reason

	return s.leaveRepo.Update(leave)
}

// CancelLeave cancels a leave request
func (s *LeaveServiceImpl) CancelLeave(leaveID uint, userID uint) error {
	leave, err := s.leaveRepo.GetByID(leaveID)
	if err != nil {
		return err
	}

	// Only the leave owner can cancel their own leave
	if leave.UserID != userID {
		return domain.ErrUnauthorized
	}

	if !leave.CanCancel() {
		return domain.ErrCannotCancelApprovedLeave
	}

	leave.Status = domain.LeaveStatusCancelled
	return s.leaveRepo.Update(leave)
}

// GetUserLeaveBalance retrieves the leave balance for a user in a specific year
func (s *LeaveServiceImpl) GetUserLeaveBalance(userID uint, year int) (map[domain.LeaveTypeName]float64, error) {
	return s.leaveRepo.GetUserLeaveBalance(userID, year)
}

// CalculateLeaveDays calculates the number of days between start and end dates
// This is a simple implementation that counts calendar days
// You might want to enhance this to exclude weekends, holidays, etc.
func (s *LeaveServiceImpl) CalculateLeaveDays(startDate, endDate time.Time) float64 {
	// Normalize dates to start of day
	start := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	end := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 0, 0, 0, 0, endDate.Location())

	// Calculate the difference in days
	duration := end.Sub(start)
	days := duration.Hours() / 24

	// Add 1 to include both start and end dates
	return days + 1
}
