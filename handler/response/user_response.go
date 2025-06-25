package response

import (
	"hrm/domain"
	"time"
)

// UserResponse represents the response model for user data
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// SignUpResponse represents the response model for user registration
type SignUpResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

// SignInResponse represents the response model for user authentication
type SignInResponse struct {
	User            UserResponse         `json:"user"`
	Token           string               `json:"token"`
	LastAttendances []AttendanceResponse `json:"last_attendances"`
}

// GetUserResponse represents the response model for getting a user
type GetUserResponse struct {
	User UserResponse `json:"user"`
}

// GetCurrentUserResponse represents the response model for getting current user
type GetCurrentUserResponse struct {
	User UserResponse `json:"user"`
}

// UpdateUserResponse represents the response model for updating a user
type UpdateUserResponse struct {
	User    UserResponse `json:"user"`
	Message string       `json:"message"`
}

// DeleteUserResponse represents the response model for deleting a user
type DeleteUserResponse struct {
	Message string `json:"message"`
}

// ListUsersResponse represents the response model for listing users
type ListUsersResponse struct {
	Users  []UserResponse `json:"users"`
	Total  int            `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

// ToUserResponse converts a domain User to UserResponse
func ToUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToUserResponseList converts a slice of domain Users to UserResponse slice
func ToUserResponseList(users []domain.User) []UserResponse {
	var responses []UserResponse
	for _, user := range users {
		responses = append(responses, ToUserResponse(&user))
	}
	return responses
}
