package repository

import (
	"hrm/domain"
	"log"
	"time"

	"gorm.io/gorm"
)

// LeaveRepositoryImpl implements the LeaveRepositoryInterface
type LeaveRepositoryImpl struct {
	db *gorm.DB
}

// NewLeaveRepository creates and returns a new LeaveRepositoryImpl instance
func NewLeaveRepository(db *gorm.DB) domain.LeaveRepositoryInterface {
	return &LeaveRepositoryImpl{db: db}
}

// Create saves a new leave to the database
func (r *LeaveRepositoryImpl) Create(leave *domain.Leave) error {
	if err := r.db.Create(leave).Error; err != nil {
		log.Printf("Error creating leave: %v", err)
		return err
	}
	return nil
}

// GetByID retrieves a leave from the database by its unique ID
func (r *LeaveRepositoryImpl) GetByID(id uint) (*domain.Leave, error) {
	var leave domain.Leave
	if err := r.db.First(&leave, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrLeaveNotFound
		}
		log.Printf("Error getting leave by ID: %v", err)
		return nil, err
	}
	return &leave, nil
}

// GetByUserID retrieves all leaves for a specific user
func (r *LeaveRepositoryImpl) GetByUserID(userID uint) ([]domain.Leave, error) {
	var leaves []domain.Leave
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&leaves).Error; err != nil {
		log.Printf("Error getting leaves by user ID: %v", err)
		return nil, err
	}
	return leaves, nil
}

// GetByUserIDAndDateRange retrieves leaves for a user within a specific date range
func (r *LeaveRepositoryImpl) GetByUserIDAndDateRange(userID uint, startDate, endDate time.Time) ([]domain.Leave, error) {
	var leaves []domain.Leave
	if err := r.db.Where("user_id = ? AND ((start_date BETWEEN ? AND ?) OR (end_date BETWEEN ? AND ?) OR (start_date <= ? AND end_date >= ?))",
		userID, startDate, endDate, startDate, endDate, startDate, endDate).
		Order("created_at DESC").Find(&leaves).Error; err != nil {
		log.Printf("Error getting leaves by user ID and date range: %v", err)
		return nil, err
	}
	return leaves, nil
}

// GetByStatus retrieves all leaves with a specific status
func (r *LeaveRepositoryImpl) GetByStatus(status domain.LeaveStatus) ([]domain.Leave, error) {
	var leaves []domain.Leave
	if err := r.db.Where("status = ?", status).Order("created_at DESC").Find(&leaves).Error; err != nil {
		log.Printf("Error getting leaves by status: %v", err)
		return nil, err
	}
	return leaves, nil
}

// GetByType retrieves all leaves of a specific type
func (r *LeaveRepositoryImpl) GetByType(leaveType domain.LeaveTypeName) ([]domain.Leave, error) {
	var leaves []domain.Leave
	if err := r.db.Where("type = ?", leaveType).Order("created_at DESC").Find(&leaves).Error; err != nil {
		log.Printf("Error getting leaves by type: %v", err)
		return nil, err
	}
	return leaves, nil
}

// GetPendingLeaves retrieves all pending leave requests
func (r *LeaveRepositoryImpl) GetPendingLeaves() ([]domain.Leave, error) {
	return r.GetByStatus(domain.LeaveStatusPending)
}

// Update modifies an existing leave in the database
func (r *LeaveRepositoryImpl) Update(leave *domain.Leave) error {
	if err := r.db.Save(leave).Error; err != nil {
		log.Printf("Error updating leave: %v", err)
		return err
	}
	return nil
}

// Delete removes a leave from the database by ID
func (r *LeaveRepositoryImpl) Delete(id uint) error {
	if err := r.db.Delete(&domain.Leave{}, id).Error; err != nil {
		log.Printf("Error deleting leave: %v", err)
		return err
	}
	return nil
}

// GetAll retrieves all leaves from the database
func (r *LeaveRepositoryImpl) GetAll() ([]domain.Leave, error) {
	var leaves []domain.Leave
	if err := r.db.Order("created_at DESC").Find(&leaves).Error; err != nil {
		log.Printf("Error getting all leaves: %v", err)
		return nil, err
	}
	return leaves, nil
}

// GetWithUser retrieves a leave with user information
func (r *LeaveRepositoryImpl) GetWithUser(id uint) (*domain.Leave, error) {
	var leave domain.Leave
	if err := r.db.Preload("User").First(&leave, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrLeaveNotFound
		}
		log.Printf("Error getting leave with user: %v", err)
		return nil, err
	}
	return &leave, nil
}

// GetUserLeaveBalance calculates the leave balance for a user in a specific year
func (r *LeaveRepositoryImpl) GetUserLeaveBalance(userID uint, year int) (map[domain.LeaveTypeName]float64, error) {
	startOfYear := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endOfYear := time.Date(year, 12, 31, 23, 59, 59, 999999999, time.UTC)

	var leaves []domain.Leave
	if err := r.db.Where("user_id = ? AND status = ? AND start_date BETWEEN ? AND ?",
		userID, domain.LeaveStatusApproved, startOfYear, endOfYear).
		Find(&leaves).Error; err != nil {
		log.Printf("Error getting user leave balance: %v", err)
		return nil, err
	}

	balance := make(map[domain.LeaveTypeName]float64)

	// Initialize balance for all leave types
	balance[domain.LeaveTypeSick] = 0
	balance[domain.LeaveTypeVacation] = 0
	balance[domain.LeaveTypePersonal] = 0
	balance[domain.LeaveTypeMaternity] = 0
	balance[domain.LeaveTypePaternity] = 0
	balance[domain.LeaveTypeOther] = 0

	// Calculate used leaves
	for _, leave := range leaves {
		balance[leave.Type] += leave.Days
	}

	return balance, nil
}
