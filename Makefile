# HRM Application Makefile

# Variables
APP_NAME=hrm
BUILD_DIR=build
MAIN_FILE=cmd/main.go

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=$(APP_NAME)
BINARY_UNIX=$(BINARY_NAME)_unix

# Default target
all: clean build

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(MAIN_FILE)

# Build for multiple platforms
build-all: clean
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_UNIX) -v $(MAIN_FILE)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME).exe -v $(MAIN_FILE)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)_darwin -v $(MAIN_FILE)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe
	rm -rf $(BUILD_DIR)

# Run the application
run:
	$(GOCMD) run $(MAIN_FILE)

# Run with hot reload (requires air)
dev:
	air

# Install dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run tests with race detection
test-race:
	$(GOTEST) -race ./...

# Run benchmarks
bench:
	$(GOTEST) -bench=. ./...

# Format code
fmt:
	$(GOCMD) fmt ./...

# Run linter
lint:
	golangci-lint run

# Install linter
install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install air for hot reload
install-air:
	go install github.com/cosmtrek/air@latest

# Generate API documentation
docs:
	swag init -g $(MAIN_FILE)

# Install swagger
install-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest

# Generate Postman collection
postman:
	$(GOCMD) run cmd/postman-generator/main.go

# Test API endpoints
test-api:
	./scripts/test-api.sh

# Database migrations
migrate-up:
	# Add your migration command here
	echo "Run database migrations"

migrate-down:
	# Add your rollback command here
	echo "Rollback database migrations"

# Docker commands
docker-build:
	docker build -t $(APP_NAME) .

docker-run:
	docker run -p 8080:8080 $(APP_NAME)

# Help
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  build-all    - Build for multiple platforms"
	@echo "  clean        - Clean build artifacts"
	@echo "  run          - Run the application"
	@echo "  dev          - Run with hot reload (requires air)"
	@echo "  deps         - Install dependencies"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage"
	@echo "  test-race    - Run tests with race detection"
	@echo "  bench        - Run benchmarks"
	@echo "  fmt          - Format code"
	@echo "  lint         - Run linter"
	@echo "  docs         - Generate API documentation"
	@echo "  postman      - Generate Postman collection"
	@echo "  test-api     - Test API endpoints"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  help         - Show this help"

.PHONY: all build build-all clean run dev deps test test-coverage test-race bench fmt lint docs postman test-api docker-build docker-run help 