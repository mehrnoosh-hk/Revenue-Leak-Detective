# 📘 Revenue Leak Detective

<a href="https://codecov.io/gh/mehrnoosh-hk/Revenue-Leak-Detective" > 
 <img src="https://codecov.io/gh/mehrnoosh-hk/Revenue-Leak-Detective/graph/badge.svg?token=VG56BUUMW7"/> 
 </a>

An agent that hunts down money leaks in a SaaS: failed charges, paused subscriptions, coupon misuse, “trial forever” zombies, and quiet churn signals (no logins + downgrades). It triages issues, suggests fixes, drafts customer outreach, and files tasks automatically.

---

## 🚀 Features
- **High-Performance API**: Built with Go's standard library for optimal performance
- **Structured Logging**: Using Go's native slog package for structured, leveled logging
- **Health Monitoring**: Comprehensive health check endpoints with detailed status reporting
- **Graceful Shutdown**: Proper signal handling and graceful server shutdown
- **Middleware Stack**: Request logging, panic recovery, and CORS support
- **Docker Ready**: Multi-stage Dockerfile for production deployments
- **Test Coverage**: Comprehensive test suite with benchmarks
- **Security**: Built-in security best practices and vulnerability scanning
---

<!-- ## 🛠️ Installation
```bash
# Clone the repository
git clone 

# Navigate into the project
cd Revenue-Leak-Detective -->



## 📋 Prerequisites

### For API Service (Go)
- Go 1.23.2 or higher
- PostgreSQL (for local development)

### For Workers Service (Python)
- Python 3.12 or higher
- uv package manager (recommended) or pip
- PostgreSQL (for local development)

### Optional
- Docker & Docker Compose (for containerized deployment)
- Make (for build automation)
---

## 🛠️ Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/mehrnoosh-hk/Revenue-Leak-Detective.git
   cd Revenue-Leak-Detective
   ```

2. **Install dependencies**:
Install all dependencies
   ```bash
   make deps
   ```
   Or install each services dependencies seperately

   ```bash
   make api-deps
   make workers-deps
   ```

3. **Install development tools** (optional):
This command installs development tools such as golangci-lint, ...
   ```bash
   make install-tools
   ```

## 🚀 Quick Start

### Local Development

1. **Set up environment variables**:
   ```bash
   cp .env.example .env.dev
   # Edit .env.dev with your configuration
   ```

2. **Run the application**:
   ```bash
   make api-run
   make workers-run
   ```

   Or manually:
   ```bash
   make api-build-local
   ./bin/rld-api
   ```

### Docker Deployment

1. **Build and run with Docker**:
   ```bash
   make api-docker-build #Build Docker image for Go API service
   make api-docker-run #Run Docker container for Go API service
   make workers-docker-build #Build Docker image for Workers service
   make workers-docker-run #Run Docker container for Workers service
   ```

2. **Or use docker-compose for full stack**:
   ```bash
   make docker-compose-up
   ```

## 🏗️ Project Structure

```
.
├── services/                      # Go API service
│   └── api/
│       ├── cmd/main.go            # Application entrypoint
│       ├── config/                # Configuration management
│       │   ├── types.go           # Configuration structs & constants
│       │   ├── loader.go          # Configuration loading & validation
│       │   ├── validator.go       # Configuration validation logic
│       │   ├── database.go        # Database config utilities
│       │   ├── http.go           # HTTP config utilities
│       │   ├── environment.go     # Environment detection utilities
│       │   ├── build_info.go      # Build info management
│       │   ├── errors.go          # Configuration error definitions
│       │   └── README.md          # Configuration documentation
│       ├── handlers/              # HTTP request handlers
│       │   ├── health-check.go    # Basic health check endpoints
│       │   ├── live_ready.go      # Kubernetes-style probes
│       │   └── health-check_test.go
│       ├── internal/              # Private application packages
│       │   ├── app/               # Application layer (Clean Architecture)
│       │   │   ├── app.go         # Main application setup & DI
│       │   │   └── server.go      # HTTP server implementation
│       │   ├── domain/            # Domain layer (business logic)
│       │   │   ├── health/        # Health domain services
│       │   │   │   ├── health.go  # Health service interface & impl
│       │   │   │   ├── errors.go  # Domain-specific errors
│       │   │   │   └── health_test.go
│       │   │   ├── models/        # Auto-generated domain models
│       │   │   │   └── generated_models.go # Domain models & params
│       │   │   └── README.md      # Domain layer documentation
│       │   ├── db/                # Database layer
│       │   │   ├── queries/       # SQL queries for sqlc
│       │   │   │   ├── users.sql  # User-related queries
│       │   │   │   └── README.md  # Query documentation
│       │   │   ├── repository/    # Repository pattern implementation
│       │   │   │   ├── health_repository.go     # Health repository
│       │   │   │   ├── pgx_adapter.go          # PostgreSQL adapter
│       │   │   │   ├── errors.go               # Repository errors
│       │   │   │   ├── health_repository_test.go
│       │   │   │   └── README.md               # Repository docs
│       │   │   ├── sqlc/          # Generated type-safe DB code
│       │   │   │   ├── db.go      # Database connection
│       │   │   │   ├── models.go  # Generated models
│       │   │   │   ├── querier.go # Query interfaces
│       │   │   │   └── users.sql.go # Generated queries
│       │   │   └── migrations/    # Database migrations
│       │   │       ├── 001_initial_schema.up.sql
│       │   │       ├── 001_initial_schema.down.sql
│       │   │       ├── 002_email_case_insensitive.up.sql
│       │   │       ├── 002_email_case_insensitive.down.sql
│       │   │       └── README.md  # Migration docs
│       │   └── middleware/        # HTTP middleware stack
│       │       ├── middleware.go # Core middleware functions
│       │       ├── request_id.go # Request ID middleware
│       │       └── middleware_test.go
│       ├── go.mod & go.sum        # Go module dependencies
│       ├── sqlc.yml              # sqlc configuration
│       └── ARCHITECTURE_EVOLUTION.md # Architecture documentation
├── workers/                       # Python workers service
│   ├── pyproject.toml            # Python project configuration
│   ├── uv.lock                  # uv dependency lock file
│   ├── pytest.ini               # pytest configuration
│   ├── src/agent/               # Main agent code
│   │   ├── __init__.py
│   │   ├── __main__.py          # Entry point
│   │   └── run.py               # Agent runner
│   ├── tests/                   # Test suite
│   │   ├── __init__.py
│   │   ├── conftest.py          # pytest fixtures
│   │   └── test_run.py          # Agent tests
│   └── README.md                # Workers documentation
├── deploy/                       # Deployment configurations
│   ├── docker/
│   │   ├── Dockerfile.api        # Go API Dockerfile
│   │   └── Dockerfile.workers    # Python workers Dockerfile
│   └── dev/
│       └── docker-compose.yml    # Development docker-compose
├── make/                         # Makefile modules
│   ├── api.mk                   # API service make targets
│   ├── workers.mk               # Workers service make targets
│   ├── db.mk                    # Database make targets
│   ├── docker.mk                # Docker make targets
│   ├── tools.mk                 # Development tools
│   └── variables.mk             # Shared variables
├── scripts/                      # Utility scripts
│   └── get-git-info.sh          # Git info extraction
├── docs/                        # Documentation
├── examples/                    # Code examples
├── .github/workflows/           # GitHub Actions CI/CD
│   ├── go-ci.yml               # Go service CI pipeline
│   └── python-ci.yml           # Python service CI pipeline
├── Makefile                     # Main makefile
├── .golangci.yml               # Go linting configuration
├── .gitignore                  # Git ignore rules
├── .gitattributes              # Git attributes
└── README.md                    # This file
```

## 🔧 Configuration

The application uses environment variables for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `API_PORT` | Server port | `3030` |
| `API_HOST` | Server host | `0.0.0.0` |
| `LOG_LEVEL` | Logging level (DEBUG, INFO, WARN, ERROR) | `INFO` |
| `ENVIRONMENT` | Environment (development, staging, production) | `development` |
| `POSTGRES_URL` | Database connection URL | - |
| `POSTGRES_HOST` | Database host | `localhost` |
| `POSTGRES_PORT` | Database port | `5432` |
| `POSTGRES_USER` | Database user | `postgres` |
| `POSTGRES_PASSWORD` | Database password | `password` |
| `POSTGRES_DB` | Database name | `revenue_leak_detective_dev` |

### Architecture

The application follows **Clean Architecture** principles with clear separation of concerns:

- **HTTP Handlers**: Handle HTTP requests and responses
- **Domain Services**: Contain business logic and rules
- **Repository Layer**: Abstract data access through interfaces
- **Database Layer**: SQLC-generated type-safe database operations

This architecture ensures:
- Easy testing with mocks
- Loose coupling between layers
- Clear separation of concerns
- Maintainable and scalable code

## 🏛️ Architecture Evaluation

### Overall Assessment: 8.5/10 ⭐

**Architectural Strengths:**
- ✅ **Excellent Layered Architecture**: Clean separation between HTTP, domain, repository, and database layers following Clean Architecture principles
- ✅ **SOLID Principles**: Strong adherence to SOLID principles with proper dependency injection and interface segregation
- ✅ **Modern Go Practices**: Uses contemporary Go idioms (slog, context, proper HTTP server setup)
- ✅ **Repository Pattern**: Well-implemented repository pattern with clean abstraction over data access
- ✅ **Testing Approach**: Excellent use of interfaces and mocks for testable design
- ✅ **Error Handling**: Custom error types and proper error wrapping throughout
- ✅ **Configuration Management**: Comprehensive configuration system with validation and environment-specific loading

**Areas of Excellence:**
- **Clean Code Structure**: Clear package organization with proper internal/external package separation
- **Dependency Injection**: Well-implemented DI pattern that makes testing and extension easy
- **Type Safety**: sqlc integration provides compile-time type safety for database operations
- **Middleware Stack**: Comprehensive HTTP middleware with proper chaining
- **Database Design**: Proper migration system with up/down scripts and sqlc integration

**Areas for Improvement:**

#### 🔴 High Priority (Immediate Action Recommended)
- **Missing Health Check Implementation**: The `--health` flag in main.go has TODO comments and doesn't actually perform health checks
- **Error Consistency**: Some error wrapping could be more consistent throughout the codebase

#### 🟡 Medium Priority (Consider for Next Iteration)
- **Observability**: Add metrics collection, distributed tracing, and better monitoring capabilities
- **Configuration Validation**: Enhance environment variable validation with better error messages
- **Resource Management**: Add connection pooling limits and resource cleanup improvements

#### 🟢 Low Priority (Nice to Have)
- **Code Organization**: Minor refactoring for better file organization and naming consistency
- **Documentation**: Add more inline documentation for complex business logic
- **Performance**: Add connection pooling optimizations and caching layers

### Technical Implementation Quality

**Go Idioms & Best Practices:** ⭐⭐⭐⭐⭐
- Excellent use of modern Go features (generics where appropriate, proper context usage)
- Proper error handling with custom error types
- Good use of interfaces for dependency injection
- Proper package organization and visibility rules

**Software Design Principles:** ⭐⭐⭐⭐⭐
- Strong SOLID principles implementation
- Clean Architecture with proper layer separation
- Repository pattern correctly implemented
- Good separation of concerns

**Code Quality:** ⭐⭐⭐⭐⭐
- Comprehensive test coverage with proper mocking
- Clean, readable code with good naming conventions
- Proper documentation and comments
- Good use of Go's standard library

**Maintainability:** ⭐⭐⭐⭐⭐
- Modular design makes changes isolated and predictable
- Interface-based design allows easy testing and extension
- Clear boundaries between layers prevent tight coupling
- Good documentation and architectural documentation

### Recommendations

1. **Complete Missing Features** (High Priority)
   - Implement the `--health` flag functionality in main.go
   - Add comprehensive error wrapping consistency

2. **Add Observability** (Medium Priority)
   - Integrate Prometheus metrics
   - Add distributed tracing (OpenTelemetry)
   - Enhance logging with structured fields

3. **Performance Optimizations** (Medium Priority)
   - Add database connection pooling limits
   - Implement caching layers where appropriate
   - Add performance monitoring

4. **Documentation Enhancement** (Low Priority)
   - Add API documentation (Swagger/OpenAPI)
   - Enhance inline code documentation
   - Add deployment guides and troubleshooting docs

### Conclusion

This is a **well-architected, production-ready codebase** that demonstrates excellent understanding of software design principles and Go best practices. The layered architecture, proper dependency injection, and clean code structure make it highly maintainable and scalable. The main areas for improvement are relatively minor and the foundation is solid for enterprise-level applications.

The architecture successfully balances:
- **Clean Architecture principles** with practical implementation
- **Modern Go practices** with maintainable code structure
- **Testability** with performance considerations
- **Scalability** with simplicity

**Recommendation**: This codebase is ready for production with the noted improvements. The architectural foundation is excellent and will support future growth and feature additions effectively.

---

## 📡 API Endpoints

### Health Check Endpoints

The application provides multiple health check endpoints following Kubernetes probe patterns:

- **GET** `/healthz` - Basic health check endpoint
- **GET** `/health` - Alternative health check endpoint  
- **GET** `/live` - Liveness probe (checks if application is alive)
- **GET** `/ready` - Readiness probe (checks if application is ready to serve requests)

**Basic Health Check Response**:
```json
{
  "status": "OK",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0"
}
```

**Probe Response**:
```json
{
  "status": "OK",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Health Check Behavior**:
- `/live` - Always returns 200 if the application is running (no external dependencies)
- `/ready` - Returns 200 if database is accessible, 503 if not ready
- `/healthz` & `/health` - Basic status endpoints for load balancers

## 🧪 Testing

### API Service Tests

```bash
# Run API tests
make api-test

# Run API tests with coverage
make api-test-coverage

# Run API benchmarks
make api-benchmark

# Run all API quality checks
make api-all
```

### Workers Service Tests

```bash
# Run workers tests
make workers-test

# Run workers formatting
make workers-format

# Run workers linting
make workers-lint

# Run all workers quality checks
make workers-all
```

### Combined Testing

```bash
# Run all services tests
make test

# Run complete CI pipeline
make all
```

### Code Quality

```bash
# API service quality checks
make api-fmt api-lint api-vet

# Workers service quality checks
make workers-format workers-lint

# Security scan (API service)
make api-security
```

## 🏗️ Build & Deployment

### API Service Build

```bash
# Build API service for current platform
make api-build-local

# Build API service for Linux (production)
make api-build

# Build and run API service
make api-run
```

### Workers Service Build

```bash
# Install workers dependencies
make workers-deps

# Run workers service
make workers-run

# Run workers in dry-run mode
make workers-run-dry
```

### Combined Build

```bash
# Install all dependencies
make deps

# Build all services
make all

# Quick start development
make dev
```

### Docker Build & Deployment

```bash
# Build API Docker image
make api-docker-build

# Build workers Docker image
make workers-docker-build

# Build all Docker images
make docker-build-all

# Run API Docker container
make api-docker-run

# Run workers Docker container
make workers-docker-run

# Start full development stack
make docker-compose-up

# Stop development stack
make docker-compose-down
```

## 🗄️ Database Management

### Database Setup

```bash
# Run database migrations
make migrate-up

# Rollback last migration
make migrate-down

# Check migration status
make migrate-version

# Generate SQLC code after schema changes
make sqlc-generate

# Validate migrations and SQLC sync
make sqlc-check
```

### Development Database

```bash
# Reset database (development only)
make db-reset

# Create new migration
make migrate-create NAME=migration_name
```

## 🔍 Monitoring & Observability

### Health Checks

The service provides comprehensive health check endpoints that can be used by:

- Load balancers
- Container orchestrators (Kubernetes, Docker Swarm)
- Monitoring systems

### Logging

Structured logging with different formats:
- **Development**: Human-readable text format
- **Production**: JSON format for log aggregation

Log levels: DEBUG, INFO, WARN, ERROR

### CI/CD Pipelines

The project includes comprehensive CI/CD pipelines:

#### Go Service CI (`.github/workflows/go-ci.yml`)
- **Linting**: golangci-lint with custom configuration
- **Testing**: Unit tests with race detection and coverage
- **Database**: Migration validation and SQLC code synchronization
- **Security**: gosec security scanning
- **Coverage**: Codecov integration

#### Python Service CI (`.github/workflows/python-ci.yml`)
- **Linting**: ruff linting and formatting
- **Testing**: pytest with coverage reporting
- **Dependencies**: uv dependency management

### Development Tools

```bash
# Install all development tools
make install-tools

# Check tool versions
make check-tools

# Generate git environment info
make git-env ENV_FILE=.env.dev

# Validate environment
make validate-env
```

### Metrics

The service is designed to be easily extended with metrics collection using:
- Prometheus metrics
- Custom business metrics
- Performance monitoring

## 🔒 Security

### Built-in Security Features

- Input validation
- Structured error handling
- Panic recovery middleware
- Security headers via CORS middleware
- Non-root container execution
- Static binary with minimal attack surface

<!-- ### Security Scanning

```bash
# Run security vulnerability scan
make security
``` -->

<!-- ## 🚀 Production Deployment -->

<!-- ### Docker

```bash
# Build production image
make docker-build

# Deploy with docker-compose
docker-compose -f deploy/prod/docker-compose.yml up -d
```

### Kubernetes

```yaml
# Example Kubernetes deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: rld-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: rld-api
  template:
    metadata:
      labels:
        app: rld-api
    spec:
      containers:
      - name: rld-api
        image: rld-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: ENVIRONMENT
          value: "production"
        - name: LOG_LEVEL
          value: "INFO"
        livenessProbe:
          httpGet:
            path: /healthz
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
``` -->

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests and linting (`make all`)
4. Commit your changes (`git commit -m 'Add amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

### Development Guidelines

- Follow Go and Python best practices and idioms
- Write tests for new functionality
- Update documentation as needed
- Run `make all` before committing
- Use conventional commit messages

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 📞 Support

For support and questions:

- Create an issue in the GitHub repository
- Check the documentation
- Review the test files for usage examples

---

**Built with ❤️ using Go and Python**
