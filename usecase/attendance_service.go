package usecase

import (
	"errors"
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
func (attendanceService *AttendanceService) CreateAttendance(userID uint, date time.Time) (*domain.Attendance, error) {
	// Check if user exists
	_, err := attendanceService.userRepo.GetByID(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Check if attendance already exists for this user and date
	existingAttendance, err := attendanceService.attendanceRepo.GetByUserID(userID, date)
	if err == nil && existingAttendance != nil {
		return existingAttendance, nil // Return existing attendance
	}

	// Create new attendance record
	attendance := &domain.Attendance{
		UserID: userID,
		Date:   date,
		Status: "absent",
	}

	if err := attendanceService.attendanceRepo.Create(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

// CheckIn records the check-in time for a user on a specific date
func (attendanceService *AttendanceService) CheckIn(userID uint, date time.Time) (*domain.Attendance, error) {
	// Check if user exists
	_, err := attendanceService.userRepo.GetByID(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Get or create attendance record
	attendance, err := attendanceService.attendanceRepo.GetByUserID(userID, date)
	if err != nil {
		if errors.Is(err, domain.ErrAttendanceNotFound) {
			// Create new attendance record
			attendance = &domain.Attendance{
				UserID: userID,
				Date:   date,
				Status: "present",
			}
			if err := attendanceService.attendanceRepo.Create(attendance); err != nil {
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
	if err := attendanceService.attendanceRepo.Update(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

// CheckOut records the check-out time for a user on a specific date
func (attendanceService *AttendanceService) CheckOut(userID uint, date time.Time) (*domain.Attendance, error) {
	// Check if user exists
	_, err := attendanceService.userRepo.GetByID(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Get attendance record
	attendance, err := attendanceService.attendanceRepo.GetByUserID(userID, date)
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
	attendance.Status = attendance.GetStatus()

	// Calculate work hours
	attendance.CalculateWorkHours()

	// Update attendance record
	if err := attendanceService.attendanceRepo.Update(attendance); err != nil {
		return nil, err
	}

	return attendance, nil
}

// GetAttendanceByID retrieves an attendance record by its ID
func (attendanceService *AttendanceService) GetAttendanceByID(id uint) (*domain.Attendance, error) {
	return attendanceService.attendanceRepo.GetByID(id)
}

// GetUserAttendance retrieves attendance record for a user on a specific date
func (attendanceService *AttendanceService) GetUserAttendance(userID uint, date time.Time) (*domain.Attendance, error) {
	// Check if user exists
	_, err := attendanceService.userRepo.GetByID(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return attendanceService.attendanceRepo.GetByUserID(userID, date)
}

// GetUserAttendanceRange retrieves attendance records for a user within a date range
func (attendanceService *AttendanceService) GetUserAttendanceRange(userID uint, startDate, endDate time.Time) ([]domain.Attendance, error) {
	// Check if user exists
	_, err := attendanceService.userRepo.GetByID(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return attendanceService.attendanceRepo.GetByUserIDAndDateRange(userID, startDate, endDate)
}

// GetAllAttendance retrieves all attendance records
func (attendanceService *AttendanceService) GetAllAttendance() ([]domain.Attendance, error) {
	return attendanceService.attendanceRepo.GetAll()
}

// UpdateAttendance modifies an existing attendance record
func (attendanceService *AttendanceService) UpdateAttendance(attendance *domain.Attendance) error {
	// Check if attendance exists
	existingAttendance, err := attendanceService.attendanceRepo.GetByID(attendance.ID)
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

	return attendanceService.attendanceRepo.Update(attendance)
}

// DeleteAttendance removes an attendance record
func (attendanceService *AttendanceService) DeleteAttendance(id uint) error {
	return attendanceService.attendanceRepo.Delete(id)
}

// CalculateWorkHours calculates and updates the work hours for an attendance record
func (attendanceService *AttendanceService) CalculateWorkHours(attendance *domain.Attendance) error {
	attendance.CalculateWorkHours()
	return attendanceService.attendanceRepo.Update(attendance)
}

// GetLastNAttendanceByUserID retrieves the last N attendance records for a user
func (attendanceService *AttendanceService) GetLastNAttendanceByUserID(userID uint, limit int) ([]domain.Attendance, error) {
	// Check if user exists
	_, err := attendanceService.userRepo.GetByID(userID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return attendanceService.attendanceRepo.GetLastNByUserID(userID, limit)
}
