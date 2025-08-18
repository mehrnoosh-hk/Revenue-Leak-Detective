# Revenue Leak Detective

<a href="https://codecov.io/gh/mehrnoosh-hk/Revenue-Leak-Detective" > 
 <img src="https://codecov.io/gh/mehrnoosh-hk/Revenue-Leak-Detective/graph/badge.svg?token=VG56BUUMW7"/> 
 </a>

A production-grade Go microservice for detecting revenue leaks in business operations.

## 🚀 Features

- **High-Performance API**: Built with Go's standard library for optimal performance
- **Structured Logging**: Using Go's native slog package for structured, leveled logging
- **Health Monitoring**: Comprehensive health check endpoints with detailed status reporting
- **Graceful Shutdown**: Proper signal handling and graceful server shutdown
- **Middleware Stack**: Request logging, panic recovery, and CORS support
- **Docker Ready**: Multi-stage Dockerfile for production deployments
- **Test Coverage**: Comprehensive test suite with benchmarks
- **Security**: Built-in security best practices and vulnerability scanning

## 📋 Prerequisites

- Go 1.23.2 or higher
- Docker (optional, for containerized deployment)
- Make (for using the Makefile)

## 🛠️ Installation

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd Revenue-Leak-Detective
   ```

2. **Install dependencies**:
   ```bash
   make deps
   ```

3. **Install development tools** (optional):
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
   make dev
   ```

   Or manually:
   ```bash
   make build-local
   ./bin/rld-api
   ```

### Docker Deployment

1. **Build and run with Docker**:
   ```bash
   make docker-build
   make docker-run
   ```

2. **Or use docker-compose for full stack**:
   ```bash
   docker-compose -f deploy/dev/docker-compose.yml up
   ```

## 🏗️ Project Structure

```
.
├── services/
│   └── api/
│       ├── cmd/                    # Application entrypoints
│       │   └── main.go
│       ├── config/                 # Configuration management
│       │   └── config.go
│       ├── handlers/               # HTTP handlers
│       │   ├── app.go
│       │   ├── health-check.go
│       │   └── health-check_test.go
│       └── internal/               # Internal packages
│           ├── middleware/         # HTTP middleware
│           │   └── middleware.go
│           └── server/             # Server implementation
│               └── server.go
├── workers/                        # Background workers (Python)
├── deploy/
│   ├── docker/
│   │   └── Dockerfile
│   └── dev/
│       └── docker-compose.yml
├── .golangci.yml                   # Linting configuration
├── Makefile                        # Build automation
└── README.md
```

## 🔧 Configuration

The application uses environment variables for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `API_PORT` | Server port | `8080` |
| `LOG_LEVEL` | Logging level (DEBUG, INFO, WARN, ERROR) | `INFO` |
| `ENVIRONMENT` | Environment (development, staging, production) | `development` |

## 📡 API Endpoints

### Health Check

- **GET** `/healthz` - Primary health check endpoint
- **GET** `/health` - Alternative health check endpoint

**Response**:
```json
{
  "status": "OK",
  "timestamp": "2024-01-15T10:30:00Z",
  "version": "1.0.0"
}
```

## 🧪 Testing

### Run Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
make benchmark
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run static analysis
make vet

# Security scan
make security
```

## 🏗️ Build & Deployment

### Local Build

```bash
# Build for current platform
make build-local

# Build for Linux (production)
make build
```

### Docker Build

```bash
# Build Docker image
make docker-build

# Run Docker container
make docker-run
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

### Security Scanning

```bash
# Run security vulnerability scan
make security
```

## 🚀 Production Deployment

### Docker

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
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests and linting (`make all`)
4. Commit your changes (`git commit -m 'Add amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

### Development Guidelines

- Follow Go best practices and idioms
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

**Built with ❤️ using Go**
