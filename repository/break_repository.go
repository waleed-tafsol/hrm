package repository

import (
	"hrm/domain"
	"time"

	"gorm.io/gorm"
)

// BreakRepository implements the BreakRepositoryInterface
// This struct handles all database operations related to breaks
type BreakRepository struct {
	db *gorm.DB
}

// NewBreakRepository creates a new instance of BreakRepository
func NewBreakRepository(db *gorm.DB) domain.BreakRepositoryInterface {
	return &BreakRepository{db: db}
}

// Create saves a new break record to the database
func (r *BreakRepository) Create(breakItem *domain.Break) error {
	// Validate break data before saving
	if err := breakItem.Validate(); err != nil {
		return err
	}

	// Set timestamps
	now := time.Now()
	breakItem.CreatedAt = now
	breakItem.UpdatedAt = now

	// Save to database
	return r.db.Create(breakItem).Error
}

// GetByID retrieves a break record by its ID
func (r *BreakRepository) GetByID(id uint) (*domain.Break, error) {
	var breakItem domain.Break

	err := r.db.First(&breakItem, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrBreakNotFound
		}
		return nil, err
	}

	return &breakItem, nil
}

// GetByAttendanceID retrieves all breaks for a specific attendance
func (r *BreakRepository) GetByAttendanceID(attendanceID uint) ([]domain.Break, error) {
	var breaks []domain.Break

	err := r.db.Where("attendance_id = ?", attendanceID).
		Order("start_time ASC").
		Find(&breaks).Error

	if err != nil {
		return nil, err
	}

	return breaks, nil
}

// Update modifies an existing break record
func (r *BreakRepository) Update(breakItem *domain.Break) error {
	// Validate break data before updating
	if err := breakItem.Validate(); err != nil {
		return err
	}

	// Update timestamp
	breakItem.UpdatedAt = time.Now()

	// Update in database
	result := r.db.Save(breakItem)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrBreakNotFound
	}

	return nil
}

// Delete removes a break record from the database
func (r *BreakRepository) Delete(id uint) error {
	result := r.db.Delete(&domain.Break{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domain.ErrBreakNotFound
	}

	return nil
}

// GetAll retrieves all break records
func (r *BreakRepository) GetAll() ([]domain.Break, error) {
	var breaks []domain.Break

	err := r.db.Order("start_time DESC").Find(&breaks).Error
	if err != nil {
		return nil, err
	}

	return breaks, nil
}

// GetActiveBreakByAttendanceID retrieves the active (not ended) break for an attendance
func (r *BreakRepository) GetActiveBreakByAttendanceID(attendanceID uint) (*domain.Break, error) {
	var breakItem domain.Break

	err := r.db.Where("attendance_id = ? AND end_time IS NULL", attendanceID).
		First(&breakItem).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrBreakNotFound
		}
		return nil, err
	}

	return &breakItem, nil
}

// GetBreaksByDateRange retrieves breaks within a date range
func (r *BreakRepository) GetBreaksByDateRange(startDate, endDate time.Time) ([]domain.Break, error) {
	var breaks []domain.Break

	err := r.db.Where("start_time >= ? AND start_time <= ?", startDate, endDate).
		Order("start_time ASC").
		Find(&breaks).Error

	if err != nil {
		return nil, err
	}

	return breaks, nil
}
