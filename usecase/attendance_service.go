package usecase

import (
	"hrm/domain"
	"time"
)

// AttendanceService implements the AttendanceServiceInterface
// This struct contains all the business logic for attendance operations
type AttendanceService struct {
	attendanceRepo domain.AttendanceRepositoryInterface
	userRepo       domain.UserRepositoryInterface
}

// NewAttendanceService creates a new instance of AttendanceService
func NewAttendanceService(
	attendanceRepo domain.AttendanceRepositoryInterface,
	userRepo domain.UserRepositoryInterface,
) domain.AttendanceServiceInterface {
	return &AttendanceService{
		attendanceRepo: attendanceRepo,
		userRepo:       userRepo,
	}
}

// CreateAttendance creates a new attendance record for a user on a specific date
func (s *AttendanceService) CreateAttendance(userID uint, date time.Time) (*domain.Attendance, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Check if attendance already exists for this user and date
	existingAttendance, err := s.attendanceRepo.GetByUserID(userID, date)
	if err == nil && existingAttendance != nil {
		return existingAttendance, nil // Return existing attendance
	}

	// Create new attendance record
	attendance := &domain.Attendance{
		UserID: userID,
		Date:   date,
	}

	if err := s.attendanceRepo.Create(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

// CheckIn records the check-in time for a user on a specific date
func (s *AttendanceService) CheckIn(userID uint, date time.Time) (*domain.Attendance, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Get or create attendance record
	attendance, err := s.attendanceRepo.GetByUserID(userID, date)
	if err != nil {
		if err == domain.ErrAttendanceNotFound {
			// Create new attendance record
			attendance = &domain.Attendance{
				UserID: userID,
				Date:   date,
			}
			if err := s.attendanceRepo.Create(attendance); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	// Check if already checked in
	if attendance.IsCheckedIn() {
		return nil, domain.ErrAlreadyCheckedIn
	}

	// Set check-in time
	now := time.Now()
	attendance.CheckInTime = &now

	// Update attendance record
	if err := s.attendanceRepo.Update(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

// CheckOut records the check-out time for a user on a specific date
func (s *AttendanceService) CheckOut(userID uint, date time.Time) (*domain.Attendance, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Get attendance record
	attendance, err := s.attendanceRepo.GetByUserID(userID, date)
	if err != nil {
		return nil, err
	}

	// Check if already checked out
	if attendance.IsCheckedOut() {
		return nil, domain.ErrAlreadyCheckedOut
	}

	// Check if checked in
	if !attendance.IsCheckedIn() {
		return nil, domain.ErrNotCheckedIn
	}

	// Set check-out time
	now := time.Now()
	attendance.CheckOutTime = &now

	// Calculate work hours
	attendance.CalculateWorkHours()

	// Update attendance record
	if err := s.attendanceRepo.Update(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

// GetAttendanceByID retrieves an attendance record by its ID
func (s *AttendanceService) GetAttendanceByID(id uint) (*domain.Attendance, error) {
	return s.attendanceRepo.GetByID(id)
}

// GetUserAttendance retrieves attendance record for a user on a specific date
func (s *AttendanceService) GetUserAttendance(userID uint, date time.Time) (*domain.Attendance, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return s.attendanceRepo.GetByUserID(userID, date)
}

// GetUserAttendanceRange retrieves attendance records for a user within a date range
func (s *AttendanceService) GetUserAttendanceRange(userID uint, startDate, endDate time.Time) ([]domain.Attendance, error) {
	// Check if user exists
	_, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return s.attendanceRepo.GetByUserIDAndDateRange(userID, startDate, endDate)
}

// GetAllAttendance retrieves all attendance records
func (s *AttendanceService) GetAllAttendance() ([]domain.Attendance, error) {
	return s.attendanceRepo.GetAll()
}

// UpdateAttendance modifies an existing attendance record
func (s *AttendanceService) UpdateAttendance(attendance *domain.Attendance) error {
	// Check if attendance exists
	existingAttendance, err := s.attendanceRepo.GetByID(attendance.ID)
	if err != nil {
		return err
	}

	// Preserve existing check-in/out times if not provided
	if attendance.CheckInTime == nil {
		attendance.CheckInTime = existingAttendance.CheckInTime
	}
	if attendance.CheckOutTime == nil {
		attendance.CheckOutTime = existingAttendance.CheckOutTime
	}

	// Recalculate work hours
	attendance.CalculateWorkHours()

	return s.attendanceRepo.Update(attendance)
}

// DeleteAttendance removes an attendance record
func (s *AttendanceService) DeleteAttendance(id uint) error {
	return s.attendanceRepo.Delete(id)
}

// CalculateWorkHours calculates and updates the work hours for an attendance record
func (s *AttendanceService) CalculateWorkHours(attendance *domain.Attendance) error {
	attendance.CalculateWorkHours()
	return s.attendanceRepo.Update(attendance)
}
