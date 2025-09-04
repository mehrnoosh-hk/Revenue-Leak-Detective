# Database Repository Layer

This directory contains the database repository layer that implements the repository pattern and provides a clean abstraction over database operations.

## Architecture Overview

### Repository Pattern Implementation

The repository pattern provides a clean abstraction over data access logic, separating business logic from data access concerns.

```
Domain Services → Repository Interfaces → Repository Implementations → Database (sqlc)
```

### Components

#### 1. Repository Interfaces (`repo_interface.go`)
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

### User Repository

The user repository provides complete CRUD operations for user management using domain-specific types and sqlc-generated type-safe database operations:

```go
// Interface definition
type UserRepository interface {
    CreateUser(ctx context.Context, arg models.CreateUserParams) (models.User, error)
    DeleteUser(ctx context.Context, id uuid.UUID) (int64, error)
    GetAllUsers(ctx context.Context) ([]models.User, error)
    GetUserByEmail(ctx context.Context, email string) (models.User, error)
    GetUserById(ctx context.Context, id uuid.UUID) (models.User, error)
    UpdateUser(ctx context.Context, arg models.UpdateUserParams) (models.User, error)
}

// Domain parameter structs (auto-generated)
type CreateUserParams struct {
    Email string `json:"email"`
    Name  string `json:"name"`
}

type UpdateUserParams struct {
    ID        uuid.UUID  `json:"id"`         // Primary key
    Email     *string    `json:"email"`      // Optional for updates
    Name      *string    `json:"name"`       // Optional for updates
    UpdatedAt *time.Time `json:"updated_at"` // Optional for updates
}
```

#### Features:
- **Domain Decoupling**: Uses domain-specific parameter structs instead of SQLC types
- **Type Safety**: Uses sqlc-generated types for compile-time safety
- **Error Handling**: Proper error wrapping and custom repository errors
- **Input Validation**: Validates input parameters before database operations
- **UUID Support**: Handles PostgreSQL UUID types properly with domain types
- **Null Handling**: Properly handles nullable database fields
- **Auto-generation**: Parameter structs are automatically generated from SQLC models

#### Testing:
```go
// Comprehensive test coverage with mocks
mockQueries := &MockQueries{}
mockQueries.On("CreateUser", mock.Anything, mock.AnythingOfType("models.CreateUserParams")).Return(mockUser, nil)

repo := repository.NewUserRepository(mockQueries)
user, err := repo.CreateUser(ctx, models.CreateUserParams{
    Email: "test@example.com",
    Name:  "Test User",
})
// Assert expectations
```

## Domain Model Generation

The repository layer automatically generates domain models and parameter structs from SQLC models:

### Auto-generated Components:
- **Domain Models**: Clean domain entities with proper Go types (e.g., `uuid.UUID` instead of `pgtype.UUID`)
- **Parameter Structs**: Type-safe parameters for CRUD operations
- **Conversion Functions**: Bidirectional conversion between domain and SQLC types

### Generation Command:
```bash
make domain-models-generate
```

This command:
1. Parses SQLC-generated models
2. Generates domain models with proper types
3. Creates parameter structs for repository operations
4. Provides conversion functions between layers

### Benefits:
- **Decoupling**: Repository interface independent of SQLC implementation
- **Type Safety**: Domain types provide better compile-time safety
- **Maintainability**: Changes to SQLC don't break repository contracts
- **Consistency**: Uniform parameter patterns across all repositories

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

1. **Define the interface** in `repo_interface.go`
2. **Create a file for the repository implementation** in `{domain}_repository.go`
3. **Implement the repository** with business logic
4. **Create database adapters** if needed
5. **Define repository errors** in `errors.go`
6. **Write comprehensive tests**

Example structure:
```
repository/
├── repo_interface.go        # Repository interface
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
