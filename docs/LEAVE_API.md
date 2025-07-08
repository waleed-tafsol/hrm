# Leave Management API Documentation

This document describes the Leave Management API endpoints for the HRM system. The leave management system allows employees to request leaves, and managers to approve or reject them.

## Authentication

All leave endpoints require JWT authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Leave Types

The system supports the following leave types:

- `sick` - Sick leave
- `vacation` - Vacation leave
- `personal` - Personal leave
- `maternity` - Maternity leave
- `paternity` - Paternity leave
- `other` - Other types of leave

## Leave Status

Leaves can have the following statuses:

- `pending` - Waiting for approval
- `approved` - Approved by manager
- `rejected` - Rejected by manager
- `cancelled` - Cancelled by employee

## API Endpoints

### 1. Create Leave Request

**POST** `/api/leaves`

Creates a new leave request for the authenticated user.

**Request Body:**
```json
{
  "type": "vacation",
  "start_date": "2024-01-15T00:00:00Z",
  "end_date": "2024-01-17T00:00:00Z",
  "reason": "Family vacation",
  "description": "Taking time off to spend with family"
}
```

**Response (201 Created):**
```json
{
  "success": true,
  "message": "Leave created successfully",
  "data": {
    "leave": {
      "id": 1,
      "user_id": 1,
      "type": "vacation",
      "status": "pending",
      "start_date": "2024-01-15T00:00:00Z",
      "end_date": "2024-01-17T00:00:00Z",
      "days": 3,
      "reason": "Family vacation",
      "description": "Taking time off to spend with family",
      "approved_by": null,
      "approved_at": null,
      "rejected_by": null,
      "rejected_at": null,
      "reject_reason": "",
      "created_at": "2024-01-10T10:30:00Z",
      "updated_at": "2024-01-10T10:30:00Z"
    },
    "message": "Leave created successfully"
  }
}
```

### 2. Get Leave by ID

**GET** `/api/leaves/:id`

Retrieves a specific leave request by its ID.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Leave retrieved successfully",
  "data": {
    "leave": {
      "id": 1,
      "user_id": 1,
      "type": "vacation",
      "status": "pending",
      "start_date": "2024-01-15T00:00:00Z",
      "end_date": "2024-01-17T00:00:00Z",
      "days": 3,
      "reason": "Family vacation",
      "description": "Taking time off to spend with family",
      "approved_by": null,
      "approved_at": null,
      "rejected_by": null,
      "rejected_at": null,
      "reject_reason": "",
      "created_at": "2024-01-10T10:30:00Z",
      "updated_at": "2024-01-10T10:30:00Z",
      "user": {
        "id": 1,
        "name": "John Doe",
        "email": "john@example.com",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    }
  }
}
```

### 3. Update Leave Request

**PUT** `/api/leaves/:id`

Updates an existing leave request. Only the leave owner can update their own leaves.

**Request Body:**
```json
{
  "type": "vacation",
  "start_date": "2024-01-16T00:00:00Z",
  "end_date": "2024-01-18T00:00:00Z",
  "reason": "Family vacation",
  "description": "Taking time off to spend with family"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Leave updated successfully",
  "data": {
    "leave": {
      "id": 1,
      "user_id": 1,
      "type": "vacation",
      "status": "pending",
      "start_date": "2024-01-16T00:00:00Z",
      "end_date": "2024-01-18T00:00:00Z",
      "days": 3,
      "reason": "Family vacation",
      "description": "Taking time off to spend with family",
      "approved_by": null,
      "approved_at": null,
      "rejected_by": null,
      "rejected_at": null,
      "reject_reason": "",
      "created_at": "2024-01-10T10:30:00Z",
      "updated_at": "2024-01-10T11:00:00Z"
    },
    "message": "Leave updated successfully"
  }
}
```

### 4. Delete Leave Request

**DELETE** `/api/leaves/:id`

Deletes a leave request. Only the leave owner can delete their own leaves.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Leave deleted successfully",
  "data": {
    "message": "Leave deleted successfully"
  }
}
```

### 5. Approve Leave Request

**POST** `/api/leaves/:id/approve`

Approves a leave request. Only managers can approve leaves.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Leave approved successfully",
  "data": {
    "leave": {
      "id": 1,
      "user_id": 1,
      "type": "vacation",
      "status": "approved",
      "start_date": "2024-01-15T00:00:00Z",
      "end_date": "2024-01-17T00:00:00Z",
      "days": 3,
      "reason": "Family vacation",
      "description": "Taking time off to spend with family",
      "approved_by": 2,
      "approved_at": "2024-01-10T12:00:00Z",
      "rejected_by": null,
      "rejected_at": null,
      "reject_reason": "",
      "created_at": "2024-01-10T10:30:00Z",
      "updated_at": "2024-01-10T12:00:00Z",
      "approver": {
        "id": 2,
        "name": "Jane Manager",
        "email": "jane@example.com",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    },
    "message": "Leave approved successfully"
  }
}
```

### 6. Reject Leave Request

**POST** `/api/leaves/:id/reject`

Rejects a leave request. Only managers can reject leaves.

**Request Body:**
```json
{
  "reject_reason": "Insufficient notice period"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Leave rejected successfully",
  "data": {
    "leave": {
      "id": 1,
      "user_id": 1,
      "type": "vacation",
      "status": "rejected",
      "start_date": "2024-01-15T00:00:00Z",
      "end_date": "2024-01-17T00:00:00Z",
      "days": 3,
      "reason": "Family vacation",
      "description": "Taking time off to spend with family",
      "approved_by": null,
      "approved_at": null,
      "rejected_by": 2,
      "rejected_at": "2024-01-10T12:00:00Z",
      "reject_reason": "Insufficient notice period",
      "created_at": "2024-01-10T10:30:00Z",
      "updated_at": "2024-01-10T12:00:00Z",
      "rejecter": {
        "id": 2,
        "name": "Jane Manager",
        "email": "jane@example.com",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    },
    "message": "Leave rejected successfully"
  }
}
```

### 7. Cancel Leave Request

**POST** `/api/leaves/:id/cancel`

Cancels a leave request. Only the leave owner can cancel their own leaves.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Leave cancelled successfully",
  "data": {
    "leave": {
      "id": 1,
      "user_id": 1,
      "type": "vacation",
      "status": "cancelled",
      "start_date": "2024-01-15T00:00:00Z",
      "end_date": "2024-01-17T00:00:00Z",
      "days": 3,
      "reason": "Family vacation",
      "description": "Taking time off to spend with family",
      "approved_by": null,
      "approved_at": null,
      "rejected_by": null,
      "rejected_at": null,
      "reject_reason": "",
      "created_at": "2024-01-10T10:30:00Z",
      "updated_at": "2024-01-10T13:00:00Z"
    },
    "message": "Leave cancelled successfully"
  }
}
```

### 8. List All Leaves

**GET** `/api/leaves`

Retrieves all leaves with pagination support.

**Query Parameters:**
- `limit` (optional): Number of leaves per page (default: 10, max: 100)
- `offset` (optional): Number of leaves to skip (default: 0)

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Leaves retrieved successfully",
  "data": {
    "leaves": [
      {
        "id": 1,
        "user_id": 1,
        "type": "vacation",
        "status": "approved",
        "start_date": "2024-01-15T00:00:00Z",
        "end_date": "2024-01-17T00:00:00Z",
        "days": 3,
        "reason": "Family vacation",
        "description": "Taking time off to spend with family",
        "approved_by": 2,
        "approved_at": "2024-01-10T12:00:00Z",
        "rejected_by": null,
        "rejected_at": null,
        "reject_reason": "",
        "created_at": "2024-01-10T10:30:00Z",
        "updated_at": "2024-01-10T12:00:00Z"
      }
    ],
    "total": 1,
    "limit": 10,
    "offset": 0
  }
}
```

### 9. Get Pending Leaves

**GET** `/api/leaves/pending`

Retrieves all pending leave requests.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Pending leaves retrieved successfully",
  "data": {
    "leaves": [
      {
        "id": 2,
        "user_id": 3,
        "type": "sick",
        "status": "pending",
        "start_date": "2024-01-20T00:00:00Z",
        "end_date": "2024-01-20T00:00:00Z",
        "days": 1,
        "reason": "Not feeling well",
        "description": "",
        "approved_by": null,
        "approved_at": null,
        "rejected_by": null,
        "rejected_at": null,
        "reject_reason": "",
        "created_at": "2024-01-19T09:00:00Z",
        "updated_at": "2024-01-19T09:00:00Z"
      }
    ],
    "total": 1,
    "limit": 1,
    "offset": 0
  }
}
```

### 10. Get User Leaves

**GET** `/api/users/:user_id/leaves`

Retrieves all leaves for a specific user.

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User leaves retrieved successfully",
  "data": {
    "leaves": [
      {
        "id": 1,
        "user_id": 1,
        "type": "vacation",
        "status": "approved",
        "start_date": "2024-01-15T00:00:00Z",
        "end_date": "2024-01-17T00:00:00Z",
        "days": 3,
        "reason": "Family vacation",
        "description": "Taking time off to spend with family",
        "approved_by": 2,
        "approved_at": "2024-01-10T12:00:00Z",
        "rejected_by": null,
        "rejected_at": null,
        "reject_reason": "",
        "created_at": "2024-01-10T10:30:00Z",
        "updated_at": "2024-01-10T12:00:00Z"
      }
    ],
    "user_id": 1
  }
}
```

### 11. Get User Leaves by Date Range

**GET** `/api/users/:user_id/leaves/range`

Retrieves leaves for a specific user within a date range.

**Query Parameters:**
- `start_date` (required): Start date in ISO format
- `end_date` (required): End date in ISO format

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User leaves by date range retrieved successfully",
  "data": {
    "leaves": [
      {
        "id": 1,
        "user_id": 1,
        "type": "vacation",
        "status": "approved",
        "start_date": "2024-01-15T00:00:00Z",
        "end_date": "2024-01-17T00:00:00Z",
        "days": 3,
        "reason": "Family vacation",
        "description": "Taking time off to spend with family",
        "approved_by": 2,
        "approved_at": "2024-01-10T12:00:00Z",
        "rejected_by": null,
        "rejected_at": null,
        "reject_reason": "",
        "created_at": "2024-01-10T10:30:00Z",
        "updated_at": "2024-01-10T12:00:00Z"
      }
    ],
    "user_id": 1
  }
}
```

### 12. Get User Leave Balance

**GET** `/api/users/:user_id/leaves/balance`

Retrieves the leave balance for a specific user in a given year.

**Query Parameters:**
- `year` (optional): Year to get balance for (default: current year)

**Response (200 OK):**
```json
{
  "success": true,
  "message": "User leave balance retrieved successfully",
  "data": {
    "user_id": 1,
    "year": 2024,
    "balance": {
      "sick": 5,
      "vacation": 10,
      "personal": 3,
      "maternity": 0,
      "paternity": 0,
      "other": 2
    }
  }
}
```

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "message": "Invalid request data: start_date is required"
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "message": "User not authenticated"
}
```

### 404 Not Found
```json
{
  "success": false,
  "message": "Leave not found"
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "message": "Failed to create leave"
}
```

## Business Rules

1. **Leave Date Validation**: Leave start date cannot be in the past
2. **Date Range Validation**: End date must be after or equal to start date
3. **Overlap Prevention**: Users cannot have overlapping leave requests
4. **Ownership**: Users can only update, delete, or cancel their own leaves
5. **Status Transitions**: 
   - Pending leaves can be approved, rejected, or cancelled
   - Approved leaves can only be cancelled
   - Rejected and cancelled leaves cannot be modified
6. **Leave Calculation**: Days are calculated inclusively (start and end dates count as full days)

## Example Usage

### Creating a Leave Request
```bash
curl -X POST http://localhost:8080/api/leaves \
  -H "Authorization: Bearer <your-jwt-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "vacation",
    "start_date": "2024-02-15T00:00:00Z",
    "end_date": "2024-02-17T00:00:00Z",
    "reason": "Family vacation",
    "description": "Taking time off to spend with family"
  }'
```

### Approving a Leave Request
```bash
curl -X POST http://localhost:8080/api/leaves/1/approve \
  -H "Authorization: Bearer <your-jwt-token>"
```

### Getting User Leave Balance
```bash
curl -X GET "http://localhost:8080/api/users/1/leaves/balance?year=2024" \
  -H "Authorization: Bearer <your-jwt-token>"
``` 