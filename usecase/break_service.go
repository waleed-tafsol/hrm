package usecase

import (
	"hrm/domain"
	"time"
)

// BreakService implements the BreakServiceInterface
// This struct contains all the business logic for break operations
type BreakService struct {
	breakRepo      domain.BreakRepositoryInterface
	attendanceRepo domain.AttendanceRepositoryInterface
}

// NewBreakService creates a new instance of BreakService
func NewBreakService(
	breakRepo domain.BreakRepositoryInterface,
	attendanceRepo domain.AttendanceRepositoryInterface,
) domain.BreakServiceInterface {
	return &BreakService{
		breakRepo:      breakRepo,
		attendanceRepo: attendanceRepo,
	}
}

// CreateBreak creates a new break record for an attendance
func (s *BreakService) CreateBreak(attendanceID uint, startTime time.Time, reason string) (*domain.Break, error) {
	// Check if attendance exists
	attendance, err := s.attendanceRepo.GetByID(attendanceID)
	if err != nil {
		return nil, err
	}

	// Check if there's already an active break for this attendance
	activeBreak, err := s.breakRepo.GetActiveBreakByAttendanceID(attendanceID)
	if err == nil && activeBreak != nil {
		return nil, domain.ErrBreakInProgress
	}

	// Create break record
	breakItem := &domain.Break{
		AttendanceID: attendanceID,
		StartTime:    startTime,
		Reason:       reason,
	}

	if err := s.breakRepo.Create(breakItem); err != nil {
		return nil, err
	}

	// Recalculate work hours for the attendance
	attendance.CalculateWorkHours()
	if err := s.attendanceRepo.Update(attendance); err != nil {
		return nil, err
	}

	return breakItem, nil
}

// GetBreakByID retrieves a break record by its ID
func (s *BreakService) GetBreakByID(id uint) (*domain.Break, error) {
	return s.breakRepo.GetByID(id)
}

// GetBreaksByAttendanceID retrieves all breaks for a specific attendance
func (s *BreakService) GetBreaksByAttendanceID(attendanceID uint) ([]domain.Break, error) {
	// Check if attendance exists
	_, err := s.attendanceRepo.GetByID(attendanceID)
	if err != nil {
		return nil, err
	}

	return s.breakRepo.GetByAttendanceID(attendanceID)
}

// GetAllBreaks retrieves all break records
func (s *BreakService) GetAllBreaks() ([]domain.Break, error) {
	return s.breakRepo.GetAll()
}

// UpdateBreak modifies an existing break record
func (s *BreakService) UpdateBreak(breakItem *domain.Break) error {
	// Check if break exists
	existingBreak, err := s.breakRepo.GetByID(breakItem.ID)
	if err != nil {
		return err
	}

	// Preserve existing end time if not provided
	if breakItem.EndTime == nil {
		breakItem.EndTime = existingBreak.EndTime
	}

	// Recalculate duration
	breakItem.CalculateDuration()

	// Update break record
	if err := s.breakRepo.Update(breakItem); err != nil {
		return err
	}

	// Recalculate work hours for the attendance
	attendance, err := s.attendanceRepo.GetByID(breakItem.AttendanceID)
	if err != nil {
		return err
	}

	attendance.CalculateWorkHours()
	return s.attendanceRepo.Update(attendance)
}

// DeleteBreak removes a break record
func (s *BreakService) DeleteBreak(id uint) error {
	// Get break to find attendance ID
	breakItem, err := s.breakRepo.GetByID(id)
	if err != nil {
		return err
	}

	attendanceID := breakItem.AttendanceID

	// Delete break record
	if err := s.breakRepo.Delete(id); err != nil {
		return err
	}

	// Recalculate work hours for the attendance
	attendance, err := s.attendanceRepo.GetByID(attendanceID)
	if err != nil {
		return err
	}

	attendance.CalculateWorkHours()
	return s.attendanceRepo.Update(attendance)
}

// EndBreak ends an existing break and calculates its duration
func (s *BreakService) EndBreak(breakID uint, endTime time.Time) error {
	// Get break record
	breakItem, err := s.breakRepo.GetByID(breakID)
	if err != nil {
		return err
	}

	// Check if break is already ended
	if breakItem.IsEnded() {
		return domain.ErrBreakAlreadyEnded
	}

	// Validate end time is after start time
	if endTime.Before(breakItem.StartTime) {
		return domain.ErrInvalidBreakTime
	}

	// Set end time
	breakItem.EndTime = &endTime

	// Calculate duration
	breakItem.CalculateDuration()

	// Update break record
	if err := s.breakRepo.Update(breakItem); err != nil {
		return err
	}

	// Recalculate work hours for the attendance
	attendance, err := s.attendanceRepo.GetByID(breakItem.AttendanceID)
	if err != nil {
		return err
	}

	attendance.CalculateWorkHours()
	return s.attendanceRepo.Update(attendance)
}

// CalculateBreakDuration calculates and updates the duration for a break record
func (s *BreakService) CalculateBreakDuration(breakItem *domain.Break) error {
	breakItem.CalculateDuration()
	return s.breakRepo.Update(breakItem)
}
