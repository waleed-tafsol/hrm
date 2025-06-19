# Testing Structure

This directory should contain all tests for the HRM application.

## Recommended Test Structure

```
tests/
├── unit/                   # Unit tests for individual components
│   ├── domain/            # Domain entity tests
│   ├── repository/        # Repository layer tests
│   ├── usecase/           # Business logic tests
│   └── handler/           # Handler tests
├── integration/           # Integration tests
│   ├── api/              # API endpoint tests
│   └── database/         # Database integration tests
├── e2e/                  # End-to-end tests
└── mocks/                # Mock implementations for testing
```

## Test Files to Create

### Unit Tests
- `domain/user_test.go` - Test user validation and business rules
- `repository/user_repository_test.go` - Test data access operations
- `usecase/user_service_test.go` - Test business logic
- `handler/user_handler_test.go` - Test HTTP handlers

### Integration Tests
- `integration/api/user_api_test.go` - Test API endpoints
- `integration/database/db_test.go` - Test database operations

### Mock Files
- `mocks/user_repository_mock.go` - Mock repository for testing
- `mocks/user_service_mock.go` - Mock service for testing

## Testing Best Practices

1. **Use table-driven tests** for multiple scenarios
2. **Mock external dependencies** (database, external APIs)
3. **Test both success and failure cases**
4. **Use descriptive test names**
5. **Keep tests fast and isolated**
6. **Use test fixtures for consistent test data**

## Example Test Command

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test file
go test ./usecase -v

# Run tests with race detection
go test -race ./...
``` 