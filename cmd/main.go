package main

import (
	"fmt"
	"log"

	"hrm/handler"

	"github.com/gin-gonic/gin"
)

// main is the entry point of the HRM application.
// This function orchestrates the entire application startup process including:
// 1. Loading configuration and establishing database connection
// 2. Setting up dependency injection container
// 3. Configuring HTTP server with middleware
// 4. Setting up API routes
// 5. Starting the HTTP server
func main() {
	// Step 1: Initialize dependency injection container
	// This creates all the necessary dependencies (repositories, services, handlers)
	// and establishes the database connection
	container := NewContainer()

	// Step 2: Configure Gin framework mode
	// Set to release mode when running on production host for better performance
	if container.Config.Server.Host == "0.0.0.0" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Step 3: Initialize HTTP router with middleware
	// Create a new Gin router and apply all necessary middleware for:
	// - Request logging
	// - CORS handling
	// - Error recovery
	// - Custom error handling
	router := gin.New()
	router.Use(handler.Recovery())       // Recover from panics
	router.Use(handler.Logger())         // Log all requests
	router.Use(handler.CORSMiddleware()) // Handle CORS
	router.Use(handler.ErrorHandler())   // Handle errors globally

	// Step 4: Setup all API routes
	// Configure all the HTTP endpoints for the application
	container.SetupRoutes(router)

	// Step 5: Start the HTTP server
	// Build the server address and start listening for requests
	serverAddr := fmt.Sprintf("%s:%s", container.Config.Server.Host, container.Config.Server.Port)
	log.Printf("Server starting on %s", serverAddr)

	// Start the server and handle any startup errors
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
