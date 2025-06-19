# Attendance API Documentation

This document describes the Attendance Management API endpoints for the HRM system.

## Overview

The Attendance API provides functionality to manage employee attendance records, including check-in/check-out operations, break management, and attendance reporting.

## Base URL

```
http://localhost:8080/api/v1/attendance
```

## Authentication

Most endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Data Models

### Attendance

```json
{
  "id": 1,
  "user_id": 1,
  "date": "2024-01-15T00:00:00Z",
  "check_in_time": "2024-01-15T09:00:00Z",
  "check_out_time": "2024-01-15T17:00:00Z",
  "total_work_hours": 8.0,
  "status": "completed",
  "created_at": "2024-01-15T09:00:00Z",
  "updated_at": "2024-01-15T17:00:00Z",
  "breaks": []
}
```

### Break

```json
{
  "id": 1,
  "attendance_id": 1,
  "start_time": "2024-01-15T12:00:00Z",
  "end_time": "2024-01-15T13:00:00Z",
  "duration": 60.0,
  "reason": "Lunch break",
  "created_at": "2024-01-15T12:00:00Z",
  "updated_at": "2024-01-15T13:00:00Z"
}
```

## Endpoints

### 1. Check In

**POST** `/api/v1/attendance/checkin`

Records the check-in time for a user on a specific date.

**Request Body:**
```json
{
  "user_id": 1,
  "date": "2024-01-15"
}
```

**Response:**
```json
{
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
}
```

**Error Responses:**
- `409 Conflict`: Already checked in for this date
- `404 Not Found`: User not found

### 2. Check Out

**POST** `/api/v1/attendance/checkout`

Records the check-out time for a user on a specific date.

**Request Body:**
```json
{
  "user_id": 1,
  "date": "2024-01-15"
}
```

**Response:**
```json
{
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
}
```

**Error Responses:**
- `409 Conflict`: Already checked out for this date
- `400 Bad Request`: Not checked in yet
- `404 Not Found`: User or attendance not found

### 3. Create Attendance

**POST** `/api/v1/attendance`

Creates a new attendance record for a user on a specific date.

**Authentication:** Required

**Request Body:**
```json
{
  "user_id": 1,
  "date": "2024-01-15"
}
```

**Response:**
```json
{
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
}
```

### 4. Get All Attendance

**GET** `/api/v1/attendance`

Retrieves all attendance records.

**Authentication:** Required

**Response:**
```json
{
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
}
```

### 5. Get Attendance by ID

**GET** `/api/v1/attendance/{id}`

Retrieves a specific attendance record by ID.

**Authentication:** Required

**Response:**
```json
{
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
}
```

**Error Responses:**
- `404 Not Found`: Attendance not found

### 6. Get User Attendance

**GET** `/api/v1/attendance/user/{user_id}`

Retrieves attendance for a specific user on a specific date.

**Authentication:** Required

**Query Parameters:**
- `date` (optional): Date in YYYY-MM-DD format (defaults to today)

**Response:**
```json
{
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
}
```

**Error Responses:**
- `404 Not Found`: User or attendance not found

### 7. Get User Attendance Range

**POST** `/api/v1/attendance/user/range`

Retrieves attendance records for a user within a date range.

**Authentication:** Required

**Request Body:**
```json
{
  "user_id": 1,
  "start_date": "2024-01-01",
  "end_date": "2024-01-31"
}
```

**Response:**
```json
{
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
}
```

**Error Responses:**
- `404 Not Found`: User not found

### 8. Add Break

**POST** `/api/v1/attendance/breaks`

Adds a new break to an attendance record.

**Authentication:** Required

**Request Body:**
```json
{
  "attendance_id": 1,
  "start_time": "2024-01-15T12:00:00Z",
  "reason": "Lunch break"
}
```

**Response:**
```json
{
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
}
```

**Error Responses:**
- `404 Not Found`: Attendance not found

### 9. End Break

**PUT** `/api/v1/attendance/breaks/end`

Ends an existing break and calculates its duration.

**Authentication:** Required

**Request Body:**
```json
{
  "break_id": 1,
  "end_time": "2024-01-15T13:00:00Z"
}
```

**Response:**
```json
{
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
}
```

**Error Responses:**
- `404 Not Found`: Break not found
- `409 Conflict`: Break already ended

### 10. Delete Attendance

**DELETE** `/api/v1/attendance/{id}`

Deletes an attendance record.

**Authentication:** Required

**Response:**
```json
{
  "success": true,
  "message": "Attendance deleted successfully"
}
```

**Error Responses:**
- `404 Not Found`: Attendance not found

## Status Values

The attendance status can be one of the following:

- `absent`: No check-in recorded
- `present`: Checked in but not checked out
- `completed`: Checked in and checked out

## Work Hours Calculation

Total work hours are calculated as:
```
Total Work Hours = (Check-out Time - Check-in Time) - Sum of Break Durations
```

## Error Handling

All endpoints return consistent error responses:

```json
{
  "success": false,
  "message": "Error description"
}
```

Common HTTP status codes:
- `200 OK`: Success
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict (e.g., already checked in)
- `500 Internal Server Error`: Server error

## Examples

### Complete Work Day Flow

1. **Check In:**
```bash
curl -X POST http://localhost:8080/api/v1/attendance/checkin \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "date": "2024-01-15"}'
```

2. **Add Break:**
```bash
curl -X POST http://localhost:8080/api/v1/attendance/breaks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"attendance_id": 1, "start_time": "2024-01-15T12:00:00Z", "reason": "Lunch"}'
```

3. **End Break:**
```bash
curl -X PUT http://localhost:8080/api/v1/attendance/breaks/end \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"break_id": 1, "end_time": "2024-01-15T13:00:00Z"}'
```

4. **Check Out:**
```bash
curl -X POST http://localhost:8080/api/v1/attendance/checkout \
  -H "Content-Type: application/json" \
  -d '{"user_id": 1, "date": "2024-01-15"}'
```

### Get Monthly Report

```bash
curl -X POST http://localhost:8080/api/v1/attendance/user/range \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"user_id": 1, "start_date": "2024-01-01", "end_date": "2024-01-31"}'
``` 