# 🚀 HRM API Postman Collection Guide

## 📋 Overview

This guide shows you how to automatically generate and use Postman collections for your HRM API. The system includes:

- **🔧 Auto-generated Postman collection** from Go code
- **🧪 Comprehensive API testing script**
- **📁 Organized endpoint structure**
- **🔐 JWT authentication support**
- **📊 Sample requests and responses**

## 🎯 Quick Start

### 1. Generate Postman Collection
```bash
# Generate the collection
make postman

# This creates: HRM_API_Collection.json
```

### 2. Test Health Endpoint
```bash
# Start the server
make run

# In another terminal, test health
curl http://localhost:8080/health
```

### 3. Test All API Endpoints
```bash
# Run comprehensive API tests
make test-api
```

## 📁 Generated Collection Structure

The Postman collection includes:

### 🏥 Health Check
- **GET** `/health` - API health status

### 🔐 Authentication
- **POST** `/api/v1/auth/register` - User registration
- **POST** `/api/v1/auth/signin` - User sign in

### 👥 User Management
- **GET** `/api/v1/users/me` - Get current user
- **GET** `/api/v1/users` - Get all users
- **GET** `/api/v1/users/{id}` - Get user by ID
- **PUT** `/api/v1/users/{id}` - Update user
- **DELETE** `/api/v1/users/{id}` - Delete user

## 🔧 How to Use

### Step 1: Import into Postman

1. **Open Postman**
2. **Click "Import"** button
3. **Select** `HRM_API_Collection.json`
4. **Import** the collection

### Step 2: Set Up Environment

1. **Create Environment** in Postman
2. **Add Variables**:
   - `base_url`: `http://localhost:8080`
   - `jwt_token`: (leave empty initially)
   - `user_id`: `1`

### Step 3: Test Authentication Flow

1. **Register User**:
   ```json
   POST {{base_url}}/api/v1/auth/register
   {
     "first_name": "John",
     "last_name": "Doe",
     "email": "john.doe@example.com",
     "password": "securepassword123",
     "phone": "+1234567890",
     "department": "Engineering",
     "position": "Software Engineer",
     "hire_date": "2024-01-15"
   }
   ```

2. **Copy JWT Token** from response
3. **Set Environment Variable** `jwt_token` to the token value
4. **Test Protected Endpoints**

## 🧪 Testing Commands

### Manual Testing
```bash
# Health check
curl http://localhost:8080/health

# Register user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Test",
    "last_name": "User",
    "email": "test@example.com",
    "password": "password123",
    "phone": "+1234567890",
    "department": "Engineering",
    "position": "Developer",
    "hire_date": "2024-01-15"
  }'

# Sign in
curl -X POST http://localhost:8080/api/v1/auth/signin \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'

# Get current user (with JWT token)
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Automated Testing
```bash
# Run all API tests
make test-api

# This will:
# 1. Check if server is running
# 2. Test health endpoint
# 3. Register a test user
# 4. Extract JWT token
# 5. Test all authenticated endpoints
```

## 📊 Sample Responses

### Health Check
```json
{
  "status": "ok",
  "message": "HRM API is running"
}
```

### User Registration
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": 1,
      "first_name": "John",
      "last_name": "Doe",
      "email": "john.doe@example.com",
      "phone": "+1234567890",
      "department": "Engineering",
      "position": "Software Engineer",
      "hire_date": "2024-01-15T00:00:00Z",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T10:30:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

### Get Current User
```json
{
  "success": true,
  "message": "User details retrieved successfully",
  "data": {
    "id": 1,
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "phone": "+1234567890",
    "department": "Engineering",
    "position": "Software Engineer",
    "hire_date": "2024-01-15T00:00:00Z",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

## 🔄 Complete Workflow

### Development Workflow
```bash
# 1. Start development
make run

# 2. Generate Postman collection
make postman

# 3. Test API endpoints
make test-api

# 4. Import collection into Postman
# 5. Test manually in Postman
```

### Testing Workflow
```bash
# 1. Start server
make run

# 2. Run automated tests
make test-api

# 3. Check results
# 4. Fix any issues
# 5. Re-run tests
```

## 🛠️ Customization

### Adding New Endpoints

Edit `cmd/postman-generator/main.go`:

```go
{
    Name:        "New Endpoint",
    Method:      "POST",
    Path:        "/api/v1/new-endpoint",
    Description: "Description of the endpoint",
    Auth:        true, // or false
    Body:        `{"key": "value"}`, // optional
    Response:    `{"success": true}`, // optional
},
```

### Modifying Variables

Update the variables in the generator:

```go
variables := []Variable{
    {Key: "base_url", Value: "http://localhost:8080", Type: "string"},
    {Key: "jwt_token", Value: "{{auth_token}}", Type: "string"},
    {Key: "user_id", Value: "1", Type: "string"},
    {Key: "new_var", Value: "new_value", Type: "string"},
}
```

## 🔧 Troubleshooting

### Common Issues

1. **Server Not Running**
   ```bash
   # Start the server
   make run
   ```

2. **Port Already in Use**
   ```bash
   # Check what's using port 8080
   lsof -i :8080
   
   # Kill the process
   kill -9 <PID>
   ```

3. **JWT Token Issues**
   - Ensure token is copied correctly
   - Check token expiration
   - Verify token format

4. **Database Issues**
   - Check database connection
   - Verify environment variables
   - Run migrations if needed

### Debug Commands

```bash
# Check server status
curl -v http://localhost:8080/health

# Check server logs
# Look at terminal where server is running

# Test specific endpoint
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}' \
  -v
```

## 📚 Related Files

- `HRM_API_Collection.json` - Generated Postman collection
- `cmd/postman-generator/main.go` - Collection generator
- `scripts/test-api.sh` - API testing script
- `docs/API.md` - API documentation
- `Makefile` - Build and test commands

## 🎯 Benefits

- **✅ Automated**: No manual collection creation
- **✅ Consistent**: Standardized request format
- **✅ Documented**: Includes descriptions and responses
- **✅ Organized**: Logical folder structure
- **✅ Reusable**: Variables for easy testing
- **✅ Updatable**: Regenerate when API changes
- **✅ Testable**: Automated testing included

## 🚀 Next Steps

1. **Import the collection** into Postman
2. **Set up environment variables**
3. **Test authentication flow**
4. **Explore all endpoints**
5. **Customize for your needs**
6. **Add more endpoints** as needed

Your HRM API is now fully equipped with professional Postman collection generation and testing capabilities! 🎉 