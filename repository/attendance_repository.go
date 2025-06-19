package repository

import (
	"hrm/domain"
	"time"

	"gorm.io/gorm"
)

// AttendanceRepository implements the AttendanceRepositoryInterface
// This struct handles all database operations related to attendance
type AttendanceRepository struct {
	db *gorm.DB
}

// NewAttendanceRepository creates a new instance of AttendanceRepository
func NewAttendanceRepository(db *gorm.DB) domain.AttendanceRepositoryInterface {
	return &AttendanceRepository{db: db}
}

// Create saves a new attendance record to the database
func (r *AttendanceRepository) Create(attendance *domain.Attendance) error {
	// Validate attendance data before saving
	if err := attendance.Validate(); err != nil {
		return err
	}

	// Set timestamps
	now := time.Now()
	attendance.CreatedAt = now
	attendance.UpdatedAt = now

	// Save to database
	return r.db.Create(attendance).Error
}

// GetByID retrieves an attendance record by its ID
func (r *AttendanceRepository) GetByID(id uint) (*domain.Attendance, error) {
	var attendance domain.Attendance

	err := r.db.First(&attendance, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrAttendanceNotFound
		}
		return nil, err
	}

	return &attendance, nil
}

// GetByUserID retrieves an attendance record for a specific user and date
func (r *AttendanceRepository) GetByUserID(userID uint, date time.Time) (*domain.Attendance, error) {
	var attendance domain.Attendance

	// Convert date to start and end of day for accurate comparison
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.db.Where("user_id = ? AND date >= ? AND date < ?", userID, startOfDay, endOfDay).First(&attendance).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrAttendanceNotFound
		}
		return nil, err
	}

	return &attendance, nil
}

// GetByUserIDAndDateRange retrieves attendance records for a user within a date range
func (r *AttendanceRepository) GetByUserIDAndDateRange(userID uint, startDate, endDate time.Time) ([]domain.Attendance, error) {
	var attendances []domain.Attendance

	err := r.db.Where("user_id = ? AND date >= ? AND date <= ?", userID, startDate, endDate).
		Order("date ASC").
		Find(&attendances).Error

	if err != nil {
		return nil, err
	}

	return attendances, nil
}

// GetByDate retrieves all attendance records for a specific date
func (r *AttendanceRepository) GetByDate(date time.Time) ([]domain.Attendance, error) {
	var attendances []domain.Attendance

	// Convert date to start and end of day for accurate comparison
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := r.db.Where("date >= ? AND date < ?", startOfDay, endOfDay).
		Order("user_id ASC").
		Find(&attendances).Error

	if err != nil {
		return nil, err
	}

	return attendances, nil
}

// Update modifies an existing attendance record
func (r *AttendanceRepository) Update(attendance *domain.Attendance) error {
	// Validate attendance data before updating
	if err := attendance.Validate(); err != nil {
		return err
	}

	// Update timestamp
	attendance.UpdatedAt = time.Now()

	// Update in database
	result := r.db.Save(attendance)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrAttendanceNotFound
	}

	return nil
}

// Delete removes an attendance record from the database
func (r *AttendanceRepository) Delete(id uint) error {
	result := r.db.Delete(&domain.Attendance{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrAttendanceNotFound
	}

	return nil
}

// GetAll retrieves all attendance records
func (r *AttendanceRepository) GetAll() ([]domain.Attendance, error) {
	var attendances []domain.Attendance

	err := r.db.Order("date DESC, user_id ASC").Find(&attendances).Error
	if err != nil {
		return nil, err
	}

	return attendances, nil
}

// GetWithBreaks retrieves an attendance record with its associated breaks
func (r *AttendanceRepository) GetWithBreaks(id uint) (*domain.Attendance, error) {
	var attendance domain.Attendance

	err := r.db.Preload("Breaks").First(&attendance, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrAttendanceNotFound
		}
		return nil, err
	}

	return &attendance, nil
}
