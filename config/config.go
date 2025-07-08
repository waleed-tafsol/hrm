package config

import (
	"hrm/domain"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config holds all application configuration settings.
// This struct centralizes all configuration data including database connection,
// server settings, and other environment-specific configurations.
type Config struct {
	DB     *gorm.DB     // Database connection instance
	Server ServerConfig // Server configuration settings
}

// ServerConfig holds server-specific configuration settings.
// This struct contains all the settings needed to configure the HTTP server.
type ServerConfig struct {
	Port string // Server port (e.g., "8080")
	Host string // Server host (e.g., "0.0.0.0" for all interfaces)
}

// LoadConfig loads and initializes all application configuration.
// This function:
// 1. Loads environment variables from .env file
// 2. Establishes database connection
// 3. Creates and returns a Config struct with all settings
//
// Returns:
//   - *Config: Complete configuration object
//   - Panics if critical configuration fails (database connection, .env loading)
func LoadConfig() *Config {
	// Step 1: Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Step 2: Establish database connection
	db := connectDB()

	// Step 3: Create and return configuration object
	return &Config{
		DB: db,
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
	}
}

// connectDB establishes a connection to the MySQL database.
// This function:
// 1. Reads the database connection string from environment variables
// 2. Connects to the MySQL database using GORM
// 3. Runs database migrations for the User entity
// 4. Returns the database connection instance
//
// Returns:
//   - *gorm.DB: Database connection instance
//   - Panics if connection fails or migrations fail
func connectDB() *gorm.DB {
	// Step 1: Get database connection string from environment
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	// Step 2: Connect to MySQL database using GORM
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Step 3: Run database migrations for User entity
	// This will create the users table if it doesn't exist
	err = db.AutoMigrate(
		&domain.User{},
		&domain.Attendance{},
		&domain.Break{},
		&domain.LeaveType{}, // Create leave_types table first
		&domain.Leave{},     // Then create leaves table
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	// Step 4: Seed leave_types table with default data
	seedLeaveTypes(db)

	// Step 5: Create foreign key constraint after both tables exist
	createForeignKeyConstraint(db)

	log.Println("Database connected successfully")
	return db
}

// getEnv retrieves an environment variable with a fallback default value.
// This function provides a safe way to get environment variables with sensible defaults.
//
// Parameters:
//   - key: The environment variable name
//   - defaultValue: The default value to return if the environment variable is not set
//
// Returns:
//   - string: The environment variable value or the default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// seedLeaveTypes seeds the leave_types table with default data.
// This function is called after the leave_types table is created in the database.
//
// Parameters:
//   - db: The database connection instance
func seedLeaveTypes(db *gorm.DB) {
	// Check if leave_types table has data
	var count int64
	db.Model(&domain.LeaveType{}).Count(&count)

	// If table is empty, seed with default data
	if count == 0 {
		defaultLeaveTypes := []domain.LeaveType{
			{
				Type:               "sick",
				Name:               "Sick Leave",
				Description:        "Medical leave for illness, injury, or health-related issues.",
				DefaultDaysPerYear: 10,
				IsActive:           true,
				RequiresApproval:   true,
				Color:              "#dc3545",
				Icon:               "medical",
			},
			{
				Type:               "vacation",
				Name:               "Vacation Leave",
				Description:        "Annual vacation time for rest, relaxation, and personal activities.",
				DefaultDaysPerYear: 20,
				IsActive:           true,
				RequiresApproval:   true,
				Color:              "#28a745",
				Icon:               "vacation",
			},
			{
				Type:               "personal",
				Name:               "Personal Leave",
				Description:        "Personal time off for various personal matters.",
				DefaultDaysPerYear: 5,
				IsActive:           true,
				RequiresApproval:   true,
				Color:              "#ffc107",
				Icon:               "personal",
			},
			{
				Type:               "maternity",
				Name:               "Maternity Leave",
				Description:        "Leave for expecting mothers before and after childbirth.",
				DefaultDaysPerYear: 90,
				IsActive:           true,
				RequiresApproval:   true,
				Color:              "#e83e8c",
				Icon:               "maternity",
			},
			{
				Type:               "paternity",
				Name:               "Paternity Leave",
				Description:        "Leave for new fathers to bond with their newborn child.",
				DefaultDaysPerYear: 14,
				IsActive:           true,
				RequiresApproval:   true,
				Color:              "#17a2b8",
				Icon:               "paternity",
			},
			{
				Type:               "other",
				Name:               "Other Leave",
				Description:        "Other types of leave for special circumstances.",
				DefaultDaysPerYear: 0,
				IsActive:           true,
				RequiresApproval:   true,
				Color:              "#6c757d",
				Icon:               "other",
			},
		}

		// Insert default leave types
		for _, leaveType := range defaultLeaveTypes {
			if err := db.Create(&leaveType).Error; err != nil {
				log.Printf("Error seeding leave type %s: %v", leaveType.Type, err)
			}
		}

		log.Println("Leave types seeded successfully")
	}
}

// createForeignKeyConstraint creates the foreign key constraint after both tables exist.
// This function is called after the leave_types table is seeded.
//
// Parameters:
//   - db: The database connection instance
func createForeignKeyConstraint(db *gorm.DB) {
	// Check if foreign key constraint already exists
	var count int64
	db.Raw("SELECT COUNT(*) FROM information_schema.KEY_COLUMN_USAGE WHERE TABLE_NAME = 'leaves' AND REFERENCED_TABLE_NAME = 'leave_types' AND CONSTRAINT_NAME = 'fk_leaves_leave_type'").Count(&count)

	if count == 0 {
		// Add foreign key constraint between leaves.type and leave_types.type
		err := db.Exec("ALTER TABLE leaves ADD CONSTRAINT fk_leaves_leave_type FOREIGN KEY (type) REFERENCES leave_types(type) ON DELETE RESTRICT ON UPDATE CASCADE").Error
		if err != nil {
			log.Printf("Warning: Could not create foreign key constraint: %v", err)
			log.Println("This is normal if the constraint already exists or if you're using a database that doesn't support this constraint")
		} else {
			log.Println("Foreign key constraint created successfully")
		}
	} else {
		log.Println("Foreign key constraint already exists")
	}
}
