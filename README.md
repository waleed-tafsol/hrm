# HRM (Human Resource Management) API

A clean, well-structured HRM API built with Go using Clean Architecture principles.

## 🏗️ Architecture Overview

This project follows Clean Architecture principles with clear separation of concerns:

```
┌─────────────────────────────────────────────────────────────┐
│                    Presentation Layer                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Routes    │  │  Handlers   │  │ Middleware  │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                    Business Logic Layer                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Use Cases │  │   Services  │  │ Validation  │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                      Domain Layer                           │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Entities  │  │ Interfaces  │  │   Errors    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                    Data Access Layer                        │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │ Repository  │  │   Database  │  │   Models    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
└─────────────────────────────────────────────────────────────┘
```

## 📁 Project Structure

```
hrm/
├── cmd/                          # Application entry points
│   ├── main.go                   # Main application file
│   └── container.go              # Dependency injection container
├── config/                       # Configuration management
│   └── config.go                 # Database and server configuration
├── domain/                       # Core business entities and interfaces
│   └── user.go                   # User entity and business interfaces
├── repository/                   # Data access layer
│   └── user_repository.go        # User data operations
├── usecase/                      # Business logic layer
│   └── user_service.go           # User business operations
├── handler/                      # HTTP interface layer
│   ├── request/                  # Request models
│   │   └── user_request.go       # User request structures
│   ├── response/                 # Response models
│   │   └── user_response.go      # User response structures
│   ├── routes/                   # Route configurations
│   │   └── user_routes.go        # User route definitions
│   ├── user_handler.go           # User HTTP handlers
│   ├── response.go               # Generic response helpers
│   └── middleware.go             # HTTP middleware
├── .env                          # Environment variables
├── .env.example                  # Environment template
├── go.mod                        # Go module dependencies
└── README.md                     # This file
```

## 🚀 Quick Start

### Prerequisites
- Go 1.24 or higher
- MySQL 8.0 or higher
- Git

### 1. Clone the Repository
```bash
git clone <repository-url>
cd hrm
```

### 2. Set Up Environment
```bash
# Copy environment template
cp .env.example .env

# Edit .env file with your database credentials
DB_DSN=root:your_password@tcp(localhost:3306)/hrm?charset=utf8mb4&parseTime=True&loc=Local
```

### 3. Install Dependencies
```bash
go mod tidy
```

### 4. Build and Run
```bash
# Build the application
go build -o hrm ./cmd/...

# Run the application
./hrm
```

### 5. Test the API
```bash
# Health check
curl http://localhost:8080/health

# Create a user
curl -X POST http://localhost:8080/api/users/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

## 📚 API Documentation

### Base URL
```
http://localhost:8080
```

### Endpoints

#### Health Check
```http
GET /health
```
**Response:**
```json
{
  "status": "ok",
  "message": "HRM API is running"
}
```

#### User Registration
```http
POST /api/users/signup
```
**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```
**Response:**
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### User Authentication
```http
POST /api/users/signin
```
**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

#### Get User by ID
```http
GET /api/users/{id}
```

#### Update User
```http
PUT /api/users/{id}
```

#### Delete User
```http
DELETE /api/users/{id}
```

#### List Users (with pagination)
```http
GET /api/users?limit=10&offset=0
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_DSN` | MySQL connection string | Required |
| `SERVER_PORT` | Server port | 8080 |
| `SERVER_HOST` | Server host | 0.0.0.0 |
| `ENVIRONMENT` | Environment mode | development |
| `JWT_SECRET` | JWT signing secret | your_super_secret_jwt_key_here |
| `JWT_EXPIRY_HOURS` | JWT token expiry | 24 |

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./usecase -v
```

## 📝 Code Style

This project follows Go best practices:

- **Package Naming**: Use descriptive, lowercase names
- **Function Naming**: Use camelCase for private, PascalCase for public
- **Error Handling**: Always check and handle errors explicitly
- **Documentation**: Comment all exported functions and types
- **Testing**: Write tests for all business logic

## 🔍 Understanding the Code

### 1. Domain Layer (`domain/`)
Contains the core business entities and interfaces. This is the heart of your application.

### 2. Repository Layer (`repository/`)
Handles data persistence. Implements the interfaces defined in the domain layer.

### 3. Use Case Layer (`usecase/`)
Contains business logic. Orchestrates operations between repositories and domain entities.

### 4. Handler Layer (`handler/`)
Manages HTTP requests and responses. Separated into request models, response models, and route configurations.

### 5. Configuration (`config/`)
Manages application configuration and database connections.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License. 