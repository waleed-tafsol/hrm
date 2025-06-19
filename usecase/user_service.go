package usecase

import (
	"errors"
	"os"
	"strconv"
	"time"

	"hrm/domain"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// UserService implements the UserServiceInterface and contains all business logic
// for user operations. This layer orchestrates between the repository layer and
// domain entities, applying business rules and validation.
type UserService struct {
	userRepository domain.UserRepositoryInterface // Dependency on user repository
}

// NewUserService creates and returns a new UserService instance.
// This function acts as a constructor and ensures proper dependency injection.
// It takes a user repository interface, making it easy to test with mock repositories.
func NewUserService(userRepository domain.UserRepositoryInterface) domain.UserServiceInterface {
	return &UserService{userRepository: userRepository}
}

// SignUp registers a new user with the system.
// This method performs the following business operations:
// 1. Validates the user input data
// 2. Checks if a user with the same email already exists
// 3. Hashes the password for security
// 4. Creates the user in the database
// 5. Sanitizes the user data before returning
func (s *UserService) SignUp(user *domain.User) error {
	// Step 1: Validate user input data
	if err := user.Validate(); err != nil {
		return err
	}

	// Step 2: Check if user already exists with the same email
	existingUser, err := s.userRepository.GetByEmail(user.Email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		// If it's not a "not found" error, something went wrong with the database
		return err
	}
	if existingUser != nil {
		// User with this email already exists
		return domain.ErrUserAlreadyExists
	}

	// Step 3: Hash the password for security
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		return err
	}
	user.Password = string(hashedPassword)

	// Step 4: Create the user in the database
	if err := s.userRepository.Create(user); err != nil {
		return err
	}

	// Step 5: Sanitize user data before returning (remove password)
	user.Sanitize()
	return nil
}

// SignIn authenticates a user with their email and password.
// This method performs the following business operations:
// 1. Retrieves the user by email
// 2. Verifies the provided password against the stored hash
// 3. Returns the authenticated user (with password removed)
func (s *UserService) SignIn(email, password string) (*domain.User, error) {
	// Step 1: Get user by email
	user, err := s.userRepository.GetByEmail(email)
	if err != nil {
		// Don't reveal whether the email exists or not for security
		return nil, domain.ErrInvalidCredentials
	}

	// Step 2: Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// Password doesn't match
		return nil, domain.ErrInvalidCredentials
	}

	// Step 3: Return authenticated user (with password removed)
	user.Sanitize()
	return user, nil
}

// GetUserByID retrieves a user by their ID.
// This method performs the following business operations:
// 1. Retrieves the user from the database
// 2. Sanitizes the user data before returning
func (s *UserService) GetUserByID(id uint) (*domain.User, error) {
	// Step 1: Get user from database
	user, err := s.userRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Step 2: Sanitize user data before returning
	user.Sanitize()
	return user, nil
}

// GetCurrentUser retrieves the current authenticated user.
// This method is used when a user wants to get their own profile information.
// It performs the following business operations:
// 1. Retrieves the user from the database using the authenticated user ID
// 2. Sanitizes the user data before returning
func (s *UserService) GetCurrentUser(userID uint) (*domain.User, error) {
	// Step 1: Get user from database using authenticated user ID
	user, err := s.userRepository.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Step 2: Sanitize user data before returning
	user.Sanitize()
	return user, nil
}

// UpdateUser modifies an existing user's information.
// This method performs the following business operations:
// 1. Validates the updated user data
// 2. Checks if the user exists
// 3. Hashes the password if it has changed
// 4. Updates the user in the database
func (s *UserService) UpdateUser(user *domain.User) error {
	// Step 1: Validate user input data
	if err := user.Validate(); err != nil {
		return err
	}

	// Step 2: Check if user exists
	existingUser, err := s.userRepository.GetByID(user.ID)
	if err != nil {
		return err
	}

	// Step 3: Hash password if it has changed
	if user.Password != existingUser.Password {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %v", err)
			return err
		}
		user.Password = string(hashedPassword)
	}

	// Step 4: Update user in database
	return s.userRepository.Update(user)
}

// DeleteUser removes a user from the system.
// This method performs the following business operations:
// 1. Checks if the user exists
// 2. Deletes the user from the database
func (s *UserService) DeleteUser(id uint) error {
	// Step 1: Check if user exists
	_, err := s.userRepository.GetByID(id)
	if err != nil {
		return err
	}

	// Step 2: Delete user from database
	return s.userRepository.Delete(id)
}

// ListUsers retrieves a paginated list of users.
// This method performs the following business operations:
// 1. Retrieves users from the database with pagination
// 2. Sanitizes all user data before returning
func (s *UserService) ListUsers(limit, offset int) ([]domain.User, error) {
	// Step 1: Get users from database
	users, err := s.userRepository.List(limit, offset)
	if err != nil {
		return nil, err
	}

	// Step 2: Sanitize all users (remove passwords)
	for i := range users {
		users[i].Sanitize()
	}

	return users, nil
}

// GenerateJWTToken generates a JWT token for the user.
// This method creates a secure token containing user information that can be used
// for authentication in subsequent requests. The token includes:
// - User ID
// - User email
// - Token expiration time
// - Token issuance time
func (s *UserService) GenerateJWTToken(user *domain.User) (string, error) {
	// Step 1: Get JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return "", errors.New("JWT_SECRET environment variable is required")
	}

	// Step 2: Get token expiry from environment (default 24 hours)
	expiryHours := 24
	if expiryStr := os.Getenv("JWT_EXPIRY_HOURS"); expiryStr != "" {
		if hours, err := strconv.Atoi(expiryStr); err == nil {
			expiryHours = hours
		}
	}

	// Step 3: Create JWT claims
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Duration(expiryHours) * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	// Step 4: Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Step 5: Sign the token with the secret
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		log.Printf("Error generating JWT token: %v", err)
		return "", err
	}

	return tokenString, nil
}
