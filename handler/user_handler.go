package handler

import (
	"net/http"

	"hrm/domain"
	"hrm/handler/request"
	"hrm/handler/response"
	"hrm/middleware"

	"github.com/gin-gonic/gin"
)

// UserHandler handles all HTTP requests related to user operations.
// This struct coordinates between HTTP requests and the business logic layer.
// It's responsible for:
// - Parsing and validating HTTP requests
// - Calling appropriate business logic methods
// - Formatting HTTP responses
// - Handling HTTP-specific errors
type UserHandler struct {
	userService domain.UserServiceInterface // Dependency on user business logic
}

// NewUserHandler creates a new UserHandler instance.
// This function acts as a constructor and ensures proper dependency injection.
// It takes a user service interface, making it easy to test with mock services.
func NewUserHandler(userService domain.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// SetupUserRoutes sets up all user-related HTTP routes.
// This function configures the routing for all user operations including:
// - User registration (POST /api/users/signup)
// - User authentication (POST /api/users/signin)
// - Current user profile (GET /api/users/me) - requires JWT
// - User retrieval (GET /api/users/:id)
// - User updates (PUT /api/users/:id)
// - User deletion (DELETE /api/users/:id)
// - User listing (GET /api/users)
func SetupUserRoutes(router *gin.Engine, userService domain.UserServiceInterface) {
	handler := NewUserHandler(userService)

	// Group all user routes under /api/users
	users := router.Group("/api/users")
	{
		users.POST("/signup", handler.SignUp)                                    // Register new user
		users.POST("/signin", handler.SignIn)                                    // Authenticate user
		users.GET("/me", middleware.JWTAuthMiddleware(), handler.GetCurrentUser) // Get current user (requires JWT)
		users.GET("/:id", handler.GetUserByID)                                   // Get user by ID
		users.PUT("/:id", handler.UpdateUser)                                    // Update user
		users.DELETE("/:id", handler.DeleteUser)                                 // Delete user
		users.GET("/", handler.ListUsers)                                        // List users with pagination
	}
}

// SignUp handles user registration requests.
// This method:
// 1. Parses and validates the JSON request body
// 2. Converts the request to a domain User object
// 3. Calls the business logic to create the user
// 4. Generates a JWT token for the newly created user
// 5. Returns appropriate HTTP response with user data and token
func (h *UserHandler) SignUp(c *gin.Context) {
	// Step 1: Parse and validate the JSON request body
	var req request.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request data: "+err.Error())
		return
	}

	// Step 2: Convert request to domain User object
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	// Step 3: Call business logic to create user
	if err := h.userService.SignUp(user); err != nil {
		// Handle different types of business errors
		switch err {
		case domain.ErrUserAlreadyExists:
			BadRequestResponse(c, "User already exists")
		case domain.ErrInvalidEmail, domain.ErrInvalidPassword, domain.ErrInvalidName:
			BadRequestResponse(c, err.Error())
		default:
			InternalServerErrorResponse(c, "Failed to create user")
		}
		return
	}

	// Step 4: Generate JWT token for the newly created user
	token, err := h.userService.GenerateJWTToken(user)
	if err != nil {
		InternalServerErrorResponse(c, "Failed to generate authentication token")
		return
	}

	// Step 5: Return success response with user data and token
	userResponse := response.ToUserResponse(user)
	signUpResponse := response.SignUpResponse{
		User:  userResponse,
		Token: token,
	}
	SuccessResponse(c, http.StatusCreated, "User created successfully", signUpResponse)
}

// SignIn handles user authentication requests.
// This method:
// 1. Parses and validates the JSON request body
// 2. Calls the business logic to authenticate the user
// 3. Generates a JWT token for the authenticated user
// 4. Returns appropriate HTTP response with user data and token
func (h *UserHandler) SignIn(c *gin.Context) {
	// Step 1: Parse and validate the JSON request body
	var req request.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request data: "+err.Error())
		return
	}

	// Step 2: Call business logic to authenticate user
	user, err := h.userService.SignIn(req.Email, req.Password)
	if err != nil {
		// Handle authentication errors
		switch err {
		case domain.ErrInvalidCredentials:
			UnauthorizedResponse(c, "Invalid email or password")
		default:
			InternalServerErrorResponse(c, "Failed to sign in")
		}
		return
	}

	// Step 3: Generate JWT token for the authenticated user
	token, err := h.userService.GenerateJWTToken(user)
	if err != nil {
		InternalServerErrorResponse(c, "Failed to generate authentication token")
		return
	}

	// Step 4: Return success response with user data and token
	userResponse := response.ToUserResponse(user)
	signInResponse := response.SignInResponse{
		User:  userResponse,
		Token: token,
	}
	SuccessResponse(c, http.StatusOK, "Sign in successful", signInResponse)
}

// GetCurrentUser handles requests to get the current authenticated user's profile.
// This method:
// 1. Extracts the authenticated user ID from the JWT token (via middleware)
// 2. Calls the business logic to retrieve the current user
// 3. Returns appropriate HTTP response with user data
//
// This endpoint requires JWT authentication via middleware.
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	// Step 1: Get authenticated user ID from JWT middleware
	userID, exists := middleware.GetUserIDFromContext(c)
	if !exists {
		UnauthorizedResponse(c, "User not authenticated")
		return
	}

	// Step 2: Call business logic to get current user
	user, err := h.userService.GetCurrentUser(userID)
	if err != nil {
		// Handle different types of errors
		switch err {
		case domain.ErrUserNotFound:
			NotFoundResponse(c, "User not found")
		default:
			InternalServerErrorResponse(c, "Failed to get user")
		}
		return
	}

	// Step 3: Return success response with user data
	userResponse := response.ToUserResponse(user)
	SuccessResponse(c, http.StatusOK, "Current user retrieved successfully", userResponse)
}

// GetUserByID handles requests to retrieve a user by their ID.
// This method:
// 1. Parses and validates the URL parameter (user ID)
// 2. Calls the business logic to retrieve the user
// 3. Returns appropriate HTTP response with user data
func (h *UserHandler) GetUserByID(c *gin.Context) {
	// Step 1: Parse and validate the URL parameter
	var req request.GetUserByIDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		BadRequestResponse(c, "Invalid user ID")
		return
	}

	// Step 2: Call business logic to get user
	user, err := h.userService.GetUserByID(req.ID)
	if err != nil {
		// Handle different types of errors
		switch err {
		case domain.ErrUserNotFound:
			NotFoundResponse(c, "User not found")
		default:
			InternalServerErrorResponse(c, "Failed to get user")
		}
		return
	}

	// Step 3: Return success response with user data
	userResponse := response.ToUserResponse(user)
	SuccessResponse(c, http.StatusOK, "User retrieved successfully", userResponse)
}

// UpdateUser handles requests to update user information.
// This method:
// 1. Parses and validates both URL parameters and JSON request body
// 2. Converts the request to a domain User object
// 3. Calls the business logic to update the user
// 4. Returns appropriate HTTP response with updated user data
func (h *UserHandler) UpdateUser(c *gin.Context) {
	// Step 1: Parse and validate the URL parameter (user ID)
	var uriReq request.GetUserByIDRequest
	if err := c.ShouldBindUri(&uriReq); err != nil {
		BadRequestResponse(c, "Invalid user ID")
		return
	}

	// Step 2: Parse and validate the JSON request body
	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestResponse(c, "Invalid request data: "+err.Error())
		return
	}

	// Step 3: Convert request to domain User object
	user := &domain.User{
		ID:       uriReq.ID,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	// Step 4: Call business logic to update user
	if err := h.userService.UpdateUser(user); err != nil {
		// Handle different types of errors
		switch err {
		case domain.ErrUserNotFound:
			NotFoundResponse(c, "User not found")
		case domain.ErrInvalidEmail, domain.ErrInvalidPassword, domain.ErrInvalidName:
			BadRequestResponse(c, err.Error())
		default:
			InternalServerErrorResponse(c, "Failed to update user")
		}
		return
	}

	// Step 5: Return success response with updated user data
	userResponse := response.ToUserResponse(user)
	SuccessResponse(c, http.StatusOK, "User updated successfully", userResponse)
}

// DeleteUser handles requests to delete a user.
// This method:
// 1. Parses and validates the URL parameter (user ID)
// 2. Calls the business logic to delete the user
// 3. Returns appropriate HTTP response
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// Step 1: Parse and validate the URL parameter
	var req request.DeleteUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		BadRequestResponse(c, "Invalid user ID")
		return
	}

	// Step 2: Call business logic to delete user
	if err := h.userService.DeleteUser(req.ID); err != nil {
		// Handle different types of errors
		switch err {
		case domain.ErrUserNotFound:
			NotFoundResponse(c, "User not found")
		default:
			InternalServerErrorResponse(c, "Failed to delete user")
		}
		return
	}

	// Step 3: Return success response
	SuccessResponse(c, http.StatusOK, "User deleted successfully", nil)
}

// ListUsers handles requests to retrieve a paginated list of users.
// This method:
// 1. Parses and validates query parameters (limit, offset)
// 2. Calls the business logic to retrieve users
// 3. Returns appropriate HTTP response with paginated user data
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Step 1: Parse and validate query parameters
	var req request.ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		// Set default values if binding fails
		req.Limit = 10
		req.Offset = 0
	}

	// Step 2: Validate and set defaults for pagination
	if req.Limit <= 0 || req.Limit > 100 {
		req.Limit = 10
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	// Step 3: Call business logic to get users
	users, err := h.userService.ListUsers(req.Limit, req.Offset)
	if err != nil {
		InternalServerErrorResponse(c, "Failed to list users")
		return
	}

	// Step 4: Convert to response format and return
	userResponses := response.ToUserResponseList(users)
	listResponse := response.ListUsersResponse{
		Users:  userResponses,
		Total:  len(userResponses),
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	SuccessResponse(c, http.StatusOK, "Users retrieved successfully", listResponse)
}
