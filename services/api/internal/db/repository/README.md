# Database Repository Layer

This directory contains the database repository layer that implements the repository pattern and provides a clean abstraction over database operations.

## Architecture Overview

### Repository Pattern Implementation

The repository pattern provides a clean abstraction over data access logic, separating business logic from data access concerns.

```
Domain Services → Repository Interfaces → Repository Implementations → Database (sqlc)
```

### Components

#### 1. Repository Interfaces (`health_repository.go`)
- Define contracts for data access operations
- Domain-specific repository interfaces
- Independent of database implementation

#### 2. Database Adapters (`pgx_adapter.go`)
- Bridge between repository interfaces and actual database
- Implement database-specific logic
- Handle database connection details

#### 3. Error Definitions (`errors.go`)
- Repository-specific error types
- Consistent error handling across repositories

## Current Repositories

### Health Repository

The health repository demonstrates the repository pattern for health-related database operations:

```go
// Interface definition
type HealthRepository interface {
    Ping(ctx context.Context) error
}

// Implementation using sqlc
type healthRepository struct {
    db Database
}

// Database abstraction
type Database interface {
    Ping(ctx context.Context) error
}
```

#### Usage:
```go
// Create repository with pgx adapter
pgxAdapter := repository.NewPgxAdapter(db)
healthRepo := repository.NewHealthRepository(pgxAdapter)

// Use in domain services
healthService := health.NewHealthService(healthRepo)
```

#### Testing:
```go
// Easy to test with mocks
mockDB := &MockDatabase{pingError: nil}
repo := repository.NewHealthRepository(mockDB)
err := repo.Ping(ctx)
// Assert expectations
```

## Benefits of Repository Pattern

### 1. **Separation of Concerns**
- Business logic separated from data access
- Clear boundaries between layers

### 2. **Testability**
- Easy to mock repositories for unit testing
- No dependency on real database in tests

### 3. **Flexibility**
- Can swap database implementations
- Easy to add caching, logging, etc.

### 4. **Maintainability**
- Changes to database don't affect business logic
- Consistent data access patterns

## Adding New Repositories

When adding new repositories, follow this pattern:

1. **Define the interface** in `{domain}_repository.go`
2. **Implement the repository** with business logic
3. **Create database adapters** if needed
4. **Define repository errors** in `errors.go`
5. **Write comprehensive tests**

Example structure:
```
repository/
├── health_repository.go      # Health repository interface
├── user_repository.go        # User repository interface
├── pgx_adapter.go           # Database adapters
├── errors.go                # Repository errors
├── health_repository_test.go
└── user_repository_test.go
```

## Integration with sqlc

The repository layer works seamlessly with sqlc:

1. **sqlc generates** type-safe database operations
2. **Repository interfaces** define the contract
3. **Repository implementations** use sqlc-generated code
4. **Database adapters** bridge the gap

Example with sqlc:
```go
// sqlc generates this
type Queries struct {
    db DBTX
}

// Repository uses it
type userRepository struct {
    queries *sqlc.Queries
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*User, error) {
    return r.queries.GetUserById(ctx, id)
}
```

**Current Implementation:**
The health repository demonstrates a simple, clean implementation:
- Uses a `Database` interface for abstraction
- Implements the repository pattern with minimal complexity
- Provides easy testing through interface mocking
- Follows Go best practices for dependency injection

## Best Practices

### 1. **Interface Segregation**
- Keep repository interfaces focused
- Don't expose unnecessary database details

### 2. **Error Handling**
- Use repository-specific errors
- Don't leak database errors to domain layer

### 3. **Testing**
- Mock repositories for unit tests
- Use integration tests for repository implementations

### 4. **Naming Conventions**
- `{Domain}Repository` for interfaces
- `{domain}Repository` for implementations
- `{Database}Adapter` for database adapters

This repository layer ensures clean separation between business logic and data access, making the application more maintainable, testable, and scalable.
