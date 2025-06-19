package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

// PostmanCollection represents the structure of a Postman collection
type PostmanCollection struct {
	Info     CollectionInfo   `json:"info"`
	Item     []CollectionItem `json:"item"`
	Variable []Variable       `json:"variable"`
}

// CollectionInfo contains metadata about the collection
type CollectionInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Schema      string `json:"schema"`
	ExportedAt  string `json:"_exporter_id"`
}

// CollectionItem represents a request or folder in the collection
type CollectionItem struct {
	Name     string           `json:"name"`
	Request  *Request         `json:"request,omitempty"`
	Response []Response       `json:"response,omitempty"`
	Item     []CollectionItem `json:"item,omitempty"`
}

// Request represents a Postman request
type Request struct {
	Method      string       `json:"method"`
	Header      []Header     `json:"header"`
	Body        *RequestBody `json:"body,omitempty"`
	URL         URL          `json:"url"`
	Description string       `json:"description"`
	Auth        *Auth        `json:"auth,omitempty"`
}

// Header represents a request header
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

// URL represents the request URL
type URL struct {
	Raw      string   `json:"raw"`
	Protocol string   `json:"protocol"`
	Host     []string `json:"host"`
	Port     string   `json:"port"`
	Path     []string `json:"path"`
	Query    []Query  `json:"query,omitempty"`
}

// Query represents URL query parameters
type Query struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// RequestBody represents the request body
type RequestBody struct {
	Mode    string  `json:"mode"`
	Raw     string  `json:"raw"`
	Options Options `json:"options"`
}

// Options for request body
type Options struct {
	Raw RawOptions `json:"raw"`
}

// RawOptions for raw body
type RawOptions struct {
	Language string `json:"language"`
}

// Auth represents authentication
type Auth struct {
	Type   string   `json:"type"`
	Bearer []Bearer `json:"bearer"`
}

// Bearer token for JWT auth
type Bearer struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

// Response represents a sample response
type Response struct {
	Name            string   `json:"name"`
	OriginalRequest Request  `json:"originalRequest"`
	Status          string   `json:"status"`
	Code            int      `json:"code"`
	Header          []Header `json:"header"`
	Body            string   `json:"body"`
}

// Variable represents collection variables
type Variable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

// APIEndpoint represents an API endpoint definition
type APIEndpoint struct {
	Name        string
	Method      string
	Path        string
	Description string
	Auth        bool
	Body        string
	Response    string
}

func main() {
	// Define API endpoints based on your HRM application
	endpoints := []APIEndpoint{
		// Health Check
		{
			Name:        "Health Check",
			Method:      "GET",
			Path:        "/health",
			Description: "Check if the HRM API is running",
			Auth:        false,
			Response:    `{"status":"ok","message":"HRM API is running"}`,
		},
		// Authentication
		{
			Name:        "Register User",
			Method:      "POST",
			Path:        "/api/v1/auth/register",
			Description: "Register a new user account",
			Auth:        false,
			Body: `{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "password": "securepassword123",
  "phone": "+1234567890",
  "department": "Engineering",
  "position": "Software Engineer",
  "hire_date": "2024-01-15"
}`,
			Response: `{
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
}`,
		},
		{
			Name:        "Sign In",
			Method:      "POST",
			Path:        "/api/v1/auth/signin",
			Description: "Sign in with email and password",
			Auth:        false,
			Body: `{
  "email": "john.doe@example.com",
  "password": "securepassword123"
}`,
			Response: `{
  "success": true,
  "message": "User signed in successfully",
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
}`,
		},
		// User Management
		{
			Name:        "Get Current User",
			Method:      "GET",
			Path:        "/api/v1/users/me",
			Description: "Get current user details",
			Auth:        true,
			Response: `{
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
}`,
		},
		{
			Name:        "Get All Users",
			Method:      "GET",
			Path:        "/api/v1/users",
			Description: "Get all users (paginated)",
			Auth:        true,
			Response: `{
  "success": true,
  "message": "Users retrieved successfully",
  "data": [
    {
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
  ]
}`,
		},
		{
			Name:        "Get User by ID",
			Method:      "GET",
			Path:        "/api/v1/users/{{user_id}}",
			Description: "Get user details by ID",
			Auth:        true,
			Response: `{
  "success": true,
  "message": "User retrieved successfully",
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
}`,
		},
		{
			Name:        "Update User",
			Method:      "PUT",
			Path:        "/api/v1/users/{{user_id}}",
			Description: "Update user details",
			Auth:        true,
			Body: `{
  "first_name": "John",
  "last_name": "Smith",
  "phone": "+1234567890",
  "department": "Product",
  "position": "Product Manager"
}`,
			Response: `{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "first_name": "John",
    "last_name": "Smith",
    "email": "john.doe@example.com",
    "phone": "+1234567890",
    "department": "Product",
    "position": "Product Manager",
    "hire_date": "2024-01-15T00:00:00Z",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}`,
		},
		{
			Name:        "Delete User",
			Method:      "DELETE",
			Path:        "/api/v1/users/{{user_id}}",
			Description: "Delete a user",
			Auth:        true,
			Response: `{
  "success": true,
  "message": "User deleted successfully"
}`,
		},
		// Attendance Management
		{
			Name:        "Check In",
			Method:      "POST",
			Path:        "/api/v1/attendance/checkin",
			Description: "Check in for the day",
			Auth:        false,
			Body: `{
  "user_id": 1,
  "date": "2024-01-15"
}`,
			Response: `{
  "success": true,
  "message": "Check-in successful",
  "data": {
    "attendance": {
      "id": 1,
      "user_id": 1,
      "date": "2024-01-15T00:00:00Z",
      "check_in_time": "2024-01-15T09:00:00Z",
      "check_out_time": null,
      "total_work_hours": 0,
      "status": "present",
      "created_at": "2024-01-15T09:00:00Z",
      "updated_at": "2024-01-15T09:00:00Z",
      "breaks": []
    },
    "message": "Check-in successful",
    "check_in_time": "2024-01-15T09:00:00Z"
  }
}`,
		},
		{
			Name:        "Check Out",
			Method:      "POST",
			Path:        "/api/v1/attendance/checkout",
			Description: "Check out for the day",
			Auth:        false,
			Body: `{
  "user_id": 1,
  "date": "2024-01-15"
}`,
			Response: `{
  "success": true,
  "message": "Check-out successful",
  "data": {
    "attendance": {
      "id": 1,
      "user_id": 1,
      "date": "2024-01-15T00:00:00Z",
      "check_in_time": "2024-01-15T09:00:00Z",
      "check_out_time": "2024-01-15T17:00:00Z",
      "total_work_hours": 8,
      "status": "completed",
      "created_at": "2024-01-15T09:00:00Z",
      "updated_at": "2024-01-15T17:00:00Z",
      "breaks": []
    },
    "message": "Check-out successful",
    "check_out_time": "2024-01-15T17:00:00Z",
    "total_work_hours": 8
  }
}`,
		},
		{
			Name:        "Create Attendance",
			Method:      "POST",
			Path:        "/api/v1/attendance",
			Description: "Create a new attendance record",
			Auth:        true,
			Body: `{
  "user_id": 1,
  "date": "2024-01-15"
}`,
			Response: `{
  "success": true,
  "message": "Attendance created successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "date": "2024-01-15T00:00:00Z",
    "check_in_time": null,
    "check_out_time": null,
    "total_work_hours": 0,
    "status": "absent",
    "created_at": "2024-01-15T00:00:00Z",
    "updated_at": "2024-01-15T00:00:00Z",
    "breaks": []
  }
}`,
		},
		{
			Name:        "Get All Attendance",
			Method:      "GET",
			Path:        "/api/v1/attendance",
			Description: "Get all attendance records",
			Auth:        true,
			Response: `{
  "success": true,
  "message": "All attendance records retrieved successfully",
  "data": {
    "attendances": [
      {
        "id": 1,
        "user_id": 1,
        "date": "2024-01-15T00:00:00Z",
        "check_in_time": "2024-01-15T09:00:00Z",
        "check_out_time": "2024-01-15T17:00:00Z",
        "total_work_hours": 8,
        "status": "completed",
        "created_at": "2024-01-15T09:00:00Z",
        "updated_at": "2024-01-15T17:00:00Z",
        "breaks": []
      }
    ],
    "total": 1
  }
}`,
		},
		{
			Name:        "Get Attendance by ID",
			Method:      "GET",
			Path:        "/api/v1/attendance/{{attendance_id}}",
			Description: "Get attendance record by ID",
			Auth:        true,
			Response: `{
  "success": true,
  "message": "Attendance retrieved successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "date": "2024-01-15T00:00:00Z",
    "check_in_time": "2024-01-15T09:00:00Z",
    "check_out_time": "2024-01-15T17:00:00Z",
    "total_work_hours": 8,
    "status": "completed",
    "created_at": "2024-01-15T09:00:00Z",
    "updated_at": "2024-01-15T17:00:00Z",
    "breaks": []
  }
}`,
		},
		{
			Name:        "Get User Attendance",
			Method:      "GET",
			Path:        "/api/v1/attendance/user/{{user_id}}",
			Description: "Get user attendance for a specific date",
			Auth:        true,
			Response: `{
  "success": true,
  "message": "User attendance retrieved successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "date": "2024-01-15T00:00:00Z",
    "check_in_time": "2024-01-15T09:00:00Z",
    "check_out_time": "2024-01-15T17:00:00Z",
    "total_work_hours": 8,
    "status": "completed",
    "created_at": "2024-01-15T09:00:00Z",
    "updated_at": "2024-01-15T17:00:00Z",
    "breaks": []
  }
}`,
		},
		{
			Name:        "Get User Attendance Range",
			Method:      "POST",
			Path:        "/api/v1/attendance/user/range",
			Description: "Get user attendance within a date range",
			Auth:        true,
			Body: `{
  "user_id": 1,
  "start_date": "2024-01-01",
  "end_date": "2024-01-31"
}`,
			Response: `{
  "success": true,
  "message": "Attendance range retrieved successfully",
  "data": {
    "attendances": [
      {
        "id": 1,
        "user_id": 1,
        "date": "2024-01-15T00:00:00Z",
        "check_in_time": "2024-01-15T09:00:00Z",
        "check_out_time": "2024-01-15T17:00:00Z",
        "total_work_hours": 8,
        "status": "completed",
        "created_at": "2024-01-15T09:00:00Z",
        "updated_at": "2024-01-15T17:00:00Z",
        "breaks": []
      }
    ],
    "total": 1
  }
}`,
		},
		{
			Name:        "Add Break",
			Method:      "POST",
			Path:        "/api/v1/attendance/breaks",
			Description: "Add a break to an attendance record",
			Auth:        true,
			Body: `{
  "attendance_id": 1,
  "start_time": "2024-01-15T12:00:00Z",
  "reason": "Lunch break"
}`,
			Response: `{
  "success": true,
  "message": "Break added successfully",
  "data": {
    "break": {
      "id": 1,
      "attendance_id": 1,
      "start_time": "2024-01-15T12:00:00Z",
      "end_time": null,
      "duration": 0,
      "reason": "Lunch break",
      "created_at": "2024-01-15T12:00:00Z",
      "updated_at": "2024-01-15T12:00:00Z"
    },
    "message": "Break added successfully"
  }
}`,
		},
		{
			Name:        "End Break",
			Method:      "PUT",
			Path:        "/api/v1/attendance/breaks/end",
			Description: "End an existing break",
			Auth:        true,
			Body: `{
  "break_id": 1,
  "end_time": "2024-01-15T13:00:00Z"
}`,
			Response: `{
  "success": true,
  "message": "Break ended successfully",
  "data": {
    "break": {
      "id": 1,
      "attendance_id": 1,
      "start_time": "2024-01-15T12:00:00Z",
      "end_time": "2024-01-15T13:00:00Z",
      "duration": 60,
      "reason": "Lunch break",
      "created_at": "2024-01-15T12:00:00Z",
      "updated_at": "2024-01-15T13:00:00Z"
    },
    "message": "Break ended successfully",
    "duration": 60
  }
}`,
		},
		{
			Name:        "Delete Attendance",
			Method:      "DELETE",
			Path:        "/api/v1/attendance/{{attendance_id}}",
			Description: "Delete an attendance record",
			Auth:        true,
			Response: `{
  "success": true,
  "message": "Attendance deleted successfully"
}`,
		},
	}

	// Create Postman collection
	collection := createPostmanCollection(endpoints)

	// Convert to JSON
	jsonData, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal collection: %v", err)
	}

	// Write to file
	filename := "HRM_API_Collection.json"
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		log.Fatalf("Failed to write collection file: %v", err)
	}

	fmt.Printf("✅ Postman collection generated successfully: %s\n", filename)
	fmt.Printf("📁 File location: %s\n", filename)
	fmt.Printf("📊 Total endpoints: %d\n", len(endpoints))
	fmt.Printf("🔗 Import this file into Postman to test your HRM API\n")
}

func createPostmanCollection(endpoints []APIEndpoint) PostmanCollection {
	// Create collection info
	info := CollectionInfo{
		Name:        "HRM API Collection",
		Description: "Complete API collection for HRM (Human Resource Management) system",
		Schema:      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		ExportedAt:  time.Now().Format(time.RFC3339),
	}

	// Create variables
	variables := []Variable{
		{Key: "base_url", Value: "http://localhost:8080", Type: "string"},
		{Key: "jwt_token", Value: "{{auth_token}}", Type: "string"},
		{Key: "user_id", Value: "1", Type: "string"},
	}

	// Group endpoints by category
	authEndpoints := []CollectionItem{}
	userEndpoints := []CollectionItem{}
	healthEndpoints := []CollectionItem{}

	for _, endpoint := range endpoints {
		item := createCollectionItem(endpoint)

		// Group by category
		switch {
		case endpoint.Path == "/health":
			healthEndpoints = append(healthEndpoints, item)
		case endpoint.Path == "/api/v1/auth/register" || endpoint.Path == "/api/v1/auth/signin":
			authEndpoints = append(authEndpoints, item)
		default:
			userEndpoints = append(userEndpoints, item)
		}
	}

	// Create main items
	items := []CollectionItem{}

	// Health Check folder
	if len(healthEndpoints) > 0 {
		items = append(items, CollectionItem{
			Name: "Health Check",
			Item: healthEndpoints,
		})
	}

	// Authentication folder
	if len(authEndpoints) > 0 {
		items = append(items, CollectionItem{
			Name: "Authentication",
			Item: authEndpoints,
		})
	}

	// User Management folder
	if len(userEndpoints) > 0 {
		items = append(items, CollectionItem{
			Name: "User Management",
			Item: userEndpoints,
		})
	}

	return PostmanCollection{
		Info:     info,
		Item:     items,
		Variable: variables,
	}
}

func createCollectionItem(endpoint APIEndpoint) CollectionItem {
	// Create URL
	url := URL{
		Raw:      fmt.Sprintf("{{base_url}}%s", endpoint.Path),
		Protocol: "http",
		Host:     []string{"{{base_url}}"},
		Path:     []string{endpoint.Path},
	}

	// Create headers
	headers := []Header{
		{Key: "Content-Type", Value: "application/json", Type: "text"},
	}

	// Add auth header if required
	var auth *Auth
	if endpoint.Auth {
		auth = &Auth{
			Type: "bearer",
			Bearer: []Bearer{
				{Key: "token", Value: "{{jwt_token}}", Type: "string"},
			},
		}
	}

	// Create request body if provided
	var body *RequestBody
	if endpoint.Body != "" {
		body = &RequestBody{
			Mode: "raw",
			Raw:  endpoint.Body,
			Options: Options{
				Raw: RawOptions{
					Language: "json",
				},
			},
		}
	}

	// Create request
	request := &Request{
		Method:      endpoint.Method,
		Header:      headers,
		Body:        body,
		URL:         url,
		Description: endpoint.Description,
		Auth:        auth,
	}

	// Create sample response
	var responses []Response
	if endpoint.Response != "" {
		responses = append(responses, Response{
			Name:            "Sample Response",
			OriginalRequest: *request,
			Status:          "OK",
			Code:            200,
			Header: []Header{
				{Key: "Content-Type", Value: "application/json", Type: "text"},
			},
			Body: endpoint.Response,
		})
	}

	return CollectionItem{
		Name:     endpoint.Name,
		Request:  request,
		Response: responses,
	}
}
