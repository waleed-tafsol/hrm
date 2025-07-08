package main

import (
	"hrm/config"
	"hrm/domain"
	"hrm/handler"
	"hrm/handler/routes"
	"hrm/repository"
	"hrm/usecase"

	"github.com/gin-gonic/gin"
)

// Container holds all application dependencies and provides dependency injection.
// This struct follows the Dependency Injection pattern to manage all the components
// needed by the application. It centralizes the creation and wiring of:
// - Configuration settings
// - Database repositories
// - Business logic services
// - HTTP handlers
type Container struct {
	Config            *config.Config                       // Application configuration
	UserRepo          domain.UserRepositoryInterface       // User data access layer
	UserService       domain.UserServiceInterface          // User business logic layer
	AttendanceRepo    domain.AttendanceRepositoryInterface // Attendance data access layer
	BreakRepo         domain.BreakRepositoryInterface      // Break data access layer
	LeaveRepo         domain.LeaveRepositoryInterface      // Leave data access layer
	LeaveTypeRepo     domain.LeaveTypeRepositoryInterface  // Leave type data access layer
	AttendanceService domain.AttendanceServiceInterface    // Attendance business logic layer
	BreakService      domain.BreakServiceInterface         // Break business logic layer
	LeaveService      domain.LeaveServiceInterface         // Leave business logic layer
	LeaveTypeService  domain.LeaveTypeServiceInterface     // Leave type business logic layer
}

// NewContainer creates and initializes all application dependencies.
// This function acts as a factory that creates all the necessary components
// and wires them together properly. The dependency injection follows this flow:
// Config -> Repository -> Service -> Handler
//
// Returns:
//   - *Container: Fully initialized container with all dependencies
func NewContainer() *Container {
	// Step 1: Load application configuration
	// This establishes database connection and loads environment settings
	cfg := config.LoadConfig()

	// Step 2: Initialize repositories (Data Access Layer)
	// Repositories handle all database operations and implement domain interfaces
	userRepo := repository.NewUserRepository(cfg.DB)
	attendanceRepo := repository.NewAttendanceRepository(cfg.DB)
	breakRepo := repository.NewBreakRepository(cfg.DB)
	leaveRepo := repository.NewLeaveRepository(cfg.DB)
	leaveTypeRepo := repository.NewLeaveTypeRepository(cfg.DB)

	// Step 3: Initialize services (Business Logic Layer)
	// Services contain business logic and orchestrate operations between repositories
	userService := usecase.NewUserService(userRepo)
	attendanceService := usecase.NewAttendanceService(attendanceRepo, userRepo)
	breakService := usecase.NewBreakService(breakRepo, attendanceRepo)
	leaveService := usecase.NewLeaveService(leaveRepo, userRepo)
	leaveTypeService := usecase.NewLeaveTypeService(leaveTypeRepo)

	// Step 4: Create and return the container with all dependencies
	return &Container{
		Config:            cfg,
		UserRepo:          userRepo,
		UserService:       userService,
		AttendanceRepo:    attendanceRepo,
		BreakRepo:         breakRepo,
		LeaveRepo:         leaveRepo,
		LeaveTypeRepo:     leaveTypeRepo,
		AttendanceService: attendanceService,
		BreakService:      breakService,
		LeaveService:      leaveService,
		LeaveTypeService:  leaveTypeService,
	}
}

// SetupRoutes configures all HTTP routes for the application.
// This function sets up the routing structure and connects HTTP endpoints
// to their corresponding handlers. It organizes routes into logical groups:
// - Health check routes
// - User management routes
// - Attendance management routes
// - Break management routes
// - Leave management routes
// - Leave type management routes
//
// Parameters:
//   - router: The Gin router instance to configure
func (c *Container) SetupRoutes(router *gin.Engine) {
	// Step 1: Setup health check routes
	// These routes provide basic health monitoring for the application
	routes.SetupHealthRoutes(router)

	// Step 2: Setup user management routes
	// These routes handle all user-related operations (CRUD, authentication)
	routes.SetupUserRoutes(router, c.UserService, c.AttendanceService)

	// Step 3: Setup attendance management routes
	// These routes handle all attendance-related operations (check-in, check-out)
	routes.SetupAttendanceRoutes(router, c.AttendanceService)

	// Step 4: Setup break management routes
	// These routes handle all break-related operations (add, end, manage breaks)
	routes.SetupBreakRoutes(router, c.BreakService)

	// Step 5: Setup leave management routes
	// These routes handle all leave-related operations (create, approve, reject, etc.)
	routes.SetupLeaveRoutes(router, c.LeaveService, c.UserService)

	// Step 6: Setup leave type management routes
	// These routes handle all leave type-related operations (CRUD for leave types)
	leaveTypeHandler := handler.NewLeaveTypeHandler(c.LeaveTypeService)
	routes.SetupLeaveTypeRoutes(router, leaveTypeHandler)
}
