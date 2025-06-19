package repository

import (
	"hrm/domain"
	"log"

	"gorm.io/gorm"
)

// UserRepositoryImpl implements the UserRepositoryInterface.
// This struct handles all database operations for the User entity.
// It uses GORM as the ORM to interact with the MySQL database.
type UserRepositoryImpl struct {
	db *gorm.DB // Database connection instance
}

// NewUserRepository creates and returns a new UserRepositoryImpl instance.
// This function acts as a constructor and ensures proper dependency injection.
// It takes a GORM database connection and returns an interface, making it easy to test.
func NewUserRepository(db *gorm.DB) domain.UserRepositoryInterface {
	return &UserRepositoryImpl{db: db}
}

// Create saves a new user to the database.
// This method takes a user pointer, saves it to the database, and returns any error that occurs.
// The user object will be updated with the generated ID and timestamps.
func (r *UserRepositoryImpl) Create(user *domain.User) error {
	// Use GORM's Create method to insert the user into the database
	if err := r.db.Create(user).Error; err != nil {
		// Log the error for debugging purposes
		log.Printf("Error creating user: %v", err)
		return err
	}
	return nil
}

// GetByID retrieves a user from the database by their unique ID.
// If the user is not found, it returns a domain.ErrUserNotFound error.
// This method is commonly used for user profile operations.
func (r *UserRepositoryImpl) GetByID(id uint) (*domain.User, error) {
	var user domain.User

	// Use GORM's First method to find the user by ID
	if err := r.db.First(&user, id).Error; err != nil {
		// Check if the error is due to record not found
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		// Log other database errors
		log.Printf("Error getting user by ID: %v", err)
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user from the database by their email address.
// Email addresses are unique in our system, so this method will return at most one user.
// If no user is found with the given email, it returns domain.ErrUserNotFound.
// This method is commonly used for authentication and user lookup operations.
func (r *UserRepositoryImpl) GetByEmail(email string) (*domain.User, error) {
	var user domain.User

	// Use GORM's Where method to find the user by email
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		// Check if the error is due to record not found
		if err == gorm.ErrRecordNotFound {
			return nil, domain.ErrUserNotFound
		}
		// Log other database errors
		log.Printf("Error getting user by email: %v", err)
		return nil, err
	}

	return &user, nil
}

// Update modifies an existing user in the database.
// This method takes a user object with an existing ID and updates all its fields.
// The UpdatedAt timestamp will be automatically updated by GORM.
func (r *UserRepositoryImpl) Update(user *domain.User) error {
	// Use GORM's Save method to update the user record
	if err := r.db.Save(user).Error; err != nil {
		// Log the error for debugging purposes
		log.Printf("Error updating user: %v", err)
		return err
	}
	return nil
}

// Delete removes a user from the database by their ID.
// This is a soft delete operation that marks the user as deleted.
// The user record will still exist in the database but will be filtered out in queries.
func (r *UserRepositoryImpl) Delete(id uint) error {
	// Use GORM's Delete method to remove the user
	if err := r.db.Delete(&domain.User{}, id).Error; err != nil {
		// Log the error for debugging purposes
		log.Printf("Error deleting user: %v", err)
		return err
	}
	return nil
}

// List retrieves a paginated list of users from the database.
// This method supports pagination with limit and offset parameters.
// It returns a slice of users and any error that occurs during the operation.
// This method is commonly used for admin panels and user management interfaces.
func (r *UserRepositoryImpl) List(limit, offset int) ([]domain.User, error) {
	var users []domain.User

	// Use GORM's Limit and Offset methods for pagination
	if err := r.db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		// Log the error for debugging purposes
		log.Printf("Error listing users: %v", err)
		return nil, err
	}

	return users, nil
}
