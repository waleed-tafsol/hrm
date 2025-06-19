package request

// SignUpRequest represents the request model for user registration
type SignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// SignInRequest represents the request model for user authentication
type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateUserRequest represents the request model for user updates
type UpdateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// ListUsersRequest represents the request model for listing users with pagination
type ListUsersRequest struct {
	Limit  int `form:"limit" binding:"min=1,max=100"`
	Offset int `form:"offset" binding:"min=0"`
}

// GetUserByIDRequest represents the request model for getting a user by ID
type GetUserByIDRequest struct {
	ID uint `uri:"id" binding:"required"`
}

// DeleteUserRequest represents the request model for deleting a user
type DeleteUserRequest struct {
	ID uint `uri:"id" binding:"required"`
}
