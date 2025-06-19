package domain

import (
	"errors"
	"time"
)

// User represents a user entity in the HRM system.
// This is the core business object that contains all user-related data.
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`                       // Unique identifier for the user
	Name      string    `json:"name" gorm:"not null"`                       // Full name of the user
	Email     string    `json:"email" gorm:"uniqueIndex;not null;size:255"` // Email address (unique, max 255 chars)
	Password  string    `json:"-" gorm:"not null"`                          // Hashed password (hidden from JSON)
	CreatedAt time.Time `json:"created_at"`                                 // When the user was created
	UpdatedAt time.Time `json:"updated_at"`                                 // When the user was last updated
}

// UserRepositoryInterface defines the contract for user data access operations.
// This interface allows us to easily swap database implementations (MySQL, PostgreSQL, etc.)
// and makes testing easier by allowing us to create mock repositories.
type UserRepositoryInterface interface {
	// Create saves a new user to the database
	Create(user *User) error

	// GetByID retrieves a user by their unique ID
	GetByID(id uint) (*User, error)

	// GetByEmail retrieves a user by their email address
	GetByEmail(email string) (*User, error)

	// Update modifies an existing user in the database
	Update(user *User) error

	// Delete removes a user from the database by ID
	Delete(id uint) error

	// List retrieves a paginated list of users
	List(limit, offset int) ([]User, error)
}

// UserServiceInterface defines the contract for user business logic operations.
// This interface contains all the business rules and operations that can be performed on users.
type UserServiceInterface interface {
	// SignUp registers a new user with the system
	SignUp(user *User) error

	// SignIn authenticates a user with email and password
	SignIn(email, password string) (*User, error)

	// GetUserByID retrieves a user by ID (with business logic)
	GetUserByID(id uint) (*User, error)

	// GetCurrentUser retrieves the current authenticated user
	GetCurrentUser(userID uint) (*User, error)

	// UpdateUser modifies an existing user (with validation)
	UpdateUser(user *User) error

	// DeleteUser removes a user from the system
	DeleteUser(id uint) error

	// ListUsers retrieves a paginated list of users
	ListUsers(limit, offset int) ([]User, error)

	// GenerateJWTToken generates a JWT token for the user
	GenerateJWTToken(user *User) (string, error)
}

// Domain-specific errors that can occur during business operations.
// These errors are defined here so they can be used consistently across all layers.
var (
	ErrInvalidEmail       = errors.New("invalid email format")                   // Email format is not valid
	ErrInvalidPassword    = errors.New("password must be at least 6 characters") // Password is too short
	ErrInvalidName        = errors.New("name cannot be empty")                   // Name field is required
	ErrUserNotFound       = errors.New("user not found")                         // User doesn't exist
	ErrUserAlreadyExists  = errors.New("user already exists")                    // User with this email already exists
	ErrInvalidCredentials = errors.New("invalid credentials")                    // Wrong email or password
	ErrUnauthorized       = errors.New("unauthorized access")                    // User not authorized
)

// Validate performs business rule validation on the User entity.
// This method ensures that the user data meets all business requirements
// before it can be saved to the database.
func (u *User) Validate() error {
	// Check if name is provided
	if u.Name == "" {
		return ErrInvalidName
	}

	// Check if email is provided
	if u.Email == "" {
		return ErrInvalidEmail
	}

	// Check if password meets minimum length requirement
	if len(u.Password) < 6 {
		return ErrInvalidPassword
	}

	return nil
}

// Sanitize removes sensitive information from the user object
// before sending it to the client. This ensures that passwords
// and other sensitive data are never exposed in API responses.
func (u *User) Sanitize() {
	u.Password = "" // Remove password from response
}
