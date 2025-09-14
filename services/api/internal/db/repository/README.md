# Database Repository Layer

This directory contains the database repository layer that implements the repository pattern and provides a clean abstraction over database operations.

## Architecture Overview

### Repository Pattern Implementation

The repository pattern provides a clean abstraction over data access logic, separating business logic from data access concerns.

```
Domain Services → Repository Interfaces → Repository Implementations → Database (sqlc)
```

### Components

#### 1. Repository Interfaces (`interfaces.go`)
- Define contracts for data access operations
- Domain-specific repository interfaces
- Independent of database implementation

#### 2. Database Adapters
- Bridge between repository interfaces and actual database
- Implement database-specific logic
- Handle database connection details

#### 3. Error Definitions (`errors.go`)
- Repository-specific error types
- Consistent error handling across repositories

### 4.Tenant Isolation
- Guard the database operations with tenant context
- Set the tenant ID in the database session
- Set the service account flag in the database session
- Use Postgres RLS to enforce tenant isolation

### 5. Service Account Isolation
- Guard the database operations with service account context
- Set the service account flag in the database session
- Use Postgres RLS to enforce service account isolation

#### Features:
- **Domain Decoupling**: Uses domain-specific parameter structs instead of SQLC types
- **Type Safety**: Uses sqlc-generated types for compile-time safety
- **Error Handling**: Proper error wrapping and custom repository errors
- **Input Validation**: Validates input parameters before database operations
- **UUID Support**: Handles PostgreSQL UUID types properly with domain types
- **Null Handling**: Properly handles nullable database fields
- **Auto-generation**: Parameter structs are automatically generated from SQLC models

## Domain Model Generation

The repository layer automatically generates domain models and parameter structs from SQLC models:

### Auto-generated Components:
- **Domain Models**: Clean domain entities with proper Go types (e.g., `uuid.UUID` instead of `pgtype.UUID`)
- **Parameter Structs**: Type-safe parameters for CRUD operations
- **Conversion Functions**: Bidirectional conversion between domain and SQLC types (see `adapters.go`)

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

1. **Define the interface** in `interfaces.go`
2. **Create a file for the repository implementation** in `{domain}.go`
3. **Implement the repository** with business logic
4. **Create database adapters** if needed in `adapters.go`
5. **Define repository errors** in `errors.go`
6. **Write comprehensive tests** (see `{domain}_test.go`)

Example structure:
```
repository/
├── interfaces.go        # Repository interface
├── adapters.go        # Database adapters
├── health.go      # Health repository implementation
├── events.go        # Events repository implementation
├── errors.go                # Repository errors
├── health_test.go
└── events_test.go
```

## Integration with sqlc

The repository layer works seamlessly with sqlc:

1. **sqlc generates** type-safe database operations
2. **Repository interfaces** define the contract
3. **Repository implementations** use sqlc-generated code
4. **Database adapters** bridge the gap

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
