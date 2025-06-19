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
	err = db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

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
