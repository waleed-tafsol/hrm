#!/bin/bash

# HRM API Testing Script
# This script tests all the endpoints of the HRM API

BASE_URL="http://localhost:8080"
TOKEN=""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    local status=$1
    local message=$2
    case $status in
        "SUCCESS")
            echo -e "${GREEN}✓ $message${NC}"
            ;;
        "ERROR")
            echo -e "${RED}✗ $message${NC}"
            ;;
        "INFO")
            echo -e "${BLUE}ℹ $message${NC}"
            ;;
        "WARNING")
            echo -e "${YELLOW}⚠ $message${NC}"
            ;;
    esac
}

# Function to test an endpoint
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4
    
    echo -e "\n${BLUE}Testing: $description${NC}"
    echo "Endpoint: $method $BASE_URL$endpoint"
    
    if [ "$method" = "GET" ]; then
        if [ -n "$TOKEN" ]; then
            response=$(curl -s -w "\n%{http_code}" -H "Authorization: Bearer $TOKEN" "$BASE_URL$endpoint")
        else
            response=$(curl -s -w "\n%{http_code}" "$BASE_URL$endpoint")
        fi
    else
        if [ -n "$TOKEN" ]; then
            response=$(curl -s -w "\n%{http_code}" -X "$method" -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d "$data" "$BASE_URL$endpoint")
        else
            response=$(curl -s -w "\n%{http_code}" -X "$method" -H "Content-Type: application/json" -d "$data" "$BASE_URL$endpoint")
        fi
    fi
    
    # Extract status code and response body
    status_code=$(echo "$response" | tail -n1)
    response_body=$(echo "$response" | head -n -1)
    
    if [ "$status_code" -ge 200 ] && [ "$status_code" -lt 300 ]; then
        print_status "SUCCESS" "Status: $status_code"
        echo "Response: $response_body" | jq '.' 2>/dev/null || echo "Response: $response_body"
    else
        print_status "ERROR" "Status: $status_code"
        echo "Response: $response_body" | jq '.' 2>/dev/null || echo "Response: $response_body"
    fi
}

# Check if server is running
print_status "INFO" "Checking if server is running..."
if curl -s "$BASE_URL/health" > /dev/null; then
    print_status "SUCCESS" "Server is running"
else
    print_status "ERROR" "Server is not running. Please start the server first."
    exit 1
fi

print_status "INFO" "Starting API tests..."

# Test 1: Health Check
test_endpoint "GET" "/health" "" "Health Check"

# Test 2: Register User
register_data='{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "password": "password123"
}'
test_endpoint "POST" "/api/v1/auth/register" "$register_data" "Register User"

# Test 3: Sign In
signin_data='{
  "email": "john.doe@example.com",
  "password": "password123"
}'
test_endpoint "POST" "/api/v1/auth/signin" "$signin_data" "Sign In"

# Extract token from signin response for subsequent requests
TOKEN=$(curl -s -X POST -H "Content-Type: application/json" -d "$signin_data" "$BASE_URL/api/v1/auth/signin" | jq -r '.data.token' 2>/dev/null)

if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
    print_status "SUCCESS" "Token extracted successfully"
else
    print_status "WARNING" "Could not extract token, some tests may fail"
fi

# Test 4: Get Current User
test_endpoint "GET" "/api/v1/users/me" "" "Get Current User"

# Test 5: Get All Users
test_endpoint "GET" "/api/v1/users" "" "Get All Users"

# Test 6: Get User by ID
test_endpoint "GET" "/api/v1/users/1" "" "Get User by ID"

# Test 7: Update User
update_data='{
  "name": "John Smith",
  "email": "john.smith@example.com"
}'
test_endpoint "PUT" "/api/v1/users/1" "$update_data" "Update User"

# Attendance Tests
print_status "INFO" "Starting Attendance Tests..."

# Test 8: Check In
checkin_data='{
  "user_id": 1,
  "date": "'$(date +%Y-%m-%d)'"
}'
test_endpoint "POST" "/api/v1/attendance/checkin" "$checkin_data" "Check In"

# Test 9: Create Attendance
attendance_data='{
  "user_id": 1,
  "date": "'$(date +%Y-%m-%d)'"
}'
test_endpoint "POST" "/api/v1/attendance" "$attendance_data" "Create Attendance"

# Test 10: Get All Attendance
test_endpoint "GET" "/api/v1/attendance" "" "Get All Attendance"

# Test 11: Get Attendance by ID
test_endpoint "GET" "/api/v1/attendance/1" "" "Get Attendance by ID"

# Test 12: Get User Attendance
test_endpoint "GET" "/api/v1/attendance/user/1?date=$(date +%Y-%m-%d)" "" "Get User Attendance"

# Test 13: Get User Attendance Range
range_data='{
  "user_id": 1,
  "start_date": "'$(date -d '7 days ago' +%Y-%m-%d)'",
  "end_date": "'$(date +%Y-%m-%d)'"
}'
test_endpoint "POST" "/api/v1/attendance/user/range" "$range_data" "Get User Attendance Range"

# Test 14: Add Break
break_data='{
  "attendance_id": 1,
  "start_time": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
  "reason": "Lunch break"
}'
test_endpoint "POST" "/api/v1/attendance/breaks" "$break_data" "Add Break"

# Test 15: End Break
end_break_data='{
  "break_id": 1,
  "end_time": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'"
}'
test_endpoint "PUT" "/api/v1/attendance/breaks/end" "$end_break_data" "End Break"

# Test 16: Check Out
checkout_data='{
  "user_id": 1,
  "date": "'$(date +%Y-%m-%d)'"
}'
test_endpoint "POST" "/api/v1/attendance/checkout" "$checkout_data" "Check Out"

# Test 17: Delete Attendance
test_endpoint "DELETE" "/api/v1/attendance/1" "" "Delete Attendance"

# Test 18: Delete User
test_endpoint "DELETE" "/api/v1/users/1" "" "Delete User"

print_status "INFO" "API testing completed!" 