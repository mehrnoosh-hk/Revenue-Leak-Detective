# Domain Layer Architecture

This directory contains the domain layer of the application, which implements clean architecture principles and separation of concerns.

## Architecture Overview

### Before (Problematic Architecture)
```
HTTP Handlers → Direct Database Access
```

**Problems:**
- Handlers directly coupled to database implementation
- Business logic mixed with HTTP concerns
- Difficult to test without real database
- Violation of separation of concerns
- Tight coupling between layers

### After (Clean Architecture)
```
HTTP Handlers → Domain Services → Repository Interfaces → Repository Implementations → Database (sqlc)
```

**Benefits:**
- Clear separation of concerns
- Easy to test with mocks
- Business logic isolated in domain layer
- Loose coupling between layers
- Dependency inversion through interfaces

## Current Domain Services

### Health Service (`health/`)

The health service demonstrates proper layering:

#### Components:
1. **HealthService Interface** (`health.go`): Defines business logic contracts
2. **Repository Integration**: Uses `repository.HealthRepository` for data access
3. **Error Definitions** (`errors.go`): Domain-specific errors

#### Health Service Methods:
- **CheckReadiness**: Verifies if the application is ready to serve requests (checks database connectivity)
- **CheckLiveness**: Verifies if the application is alive (no external dependencies)

#### Usage:
```go
// In your application setup
pgxAdapter := repository.NewPgxAdapter(db)
healthRepo := repository.NewHealthRepository(pgxAdapter)
healthService := health.NewHealthService(healthRepo)

// In your handlers
func ReadyHandler(healthService health.HealthService, logger *slog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := healthService.CheckReadiness(r.Context()); err != nil {
            // Handle error
        }
        // Return success response
    }
}

func LiveHandler(healthService health.HealthService, logger *slog.Logger) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := healthService.CheckLiveness(r.Context()); err != nil {
            // Handle error
        }
        // Return success response
    }
}
```

#### Testing:
```go
// Easy to test with mocks
mockRepo := &MockHealthRepository{pingError: nil}
service := health.NewHealthService(mockRepo)
err := service.CheckReadiness(ctx)
// Assert expectations
```

## Best Practices

### 1. Interface Segregation
- Define small, focused interfaces
- Keep business logic separate from infrastructure concerns

### 2. Dependency Injection
- Inject dependencies through constructors
- Use interfaces for loose coupling

### 3. Error Handling
- Define domain-specific errors
- Don't leak infrastructure errors to domain layer

### 4. Testing
- Use mocks for external dependencies
- Test business logic in isolation
- Keep tests fast and reliable

## Adding New Domain Services

When adding new domain services, follow this pattern:

1. **Define the interface** in a `service.go` file
2. **Implement the service** with business logic
3. **Create adapters** for external dependencies
4. **Define domain errors** in `errors.go`
5. **Write comprehensive tests**

Example structure:
```
domain/
├── users/
│   ├── service.go          # UserService interface
│   ├── user_service.go     # Implementation
│   ├── errors.go           # Domain errors
│   └── user_service_test.go
└── health/
    └── ... (current implementation)

repository/
├── health_repository.go    # Health repository interface
├── user_repository.go      # User repository interface
├── pgx_adapter.go         # Database adapters
├── errors.go              # Repository errors
└── *_repository_test.go   # Repository tests
```

## Migration Guide

To migrate existing handlers to use domain services:

1. **Extract business logic** from handlers into domain services
2. **Create interfaces** for external dependencies
3. **Implement adapters** for database/API calls
4. **Update handlers** to use domain services
5. **Write tests** for the new domain layer

This approach ensures your application follows clean architecture principles and is maintainable, testable, and scalable.

## Database Adapter

The database adapter is used to abstract the database connection for health operations.
It is used to test the health service without depending on the actual database connection.
It is also used to test the health service with a mock database connection.


HTTP Handlers → Domain Services → Repository Interfaces → Repository Implementations → sqlc Generated Code → Database