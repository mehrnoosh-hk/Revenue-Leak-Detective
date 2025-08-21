# =============================================================================
# Revenue Leak Detective - Multi-Service Makefile
# =============================================================================

# Shell configuration
SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c

# Colors for output
RED := \e[31m
GREEN := \e[32m
YELLOW := \e[33m
BLUE := \e[34m
MAGENTA := \e[35m
CYAN := \e[36m
BOLD := \e[1m
NC := \e[0m

# =============================================================================
# Configuration Variables
# =============================================================================

# Project metadata
PROJECT_NAME := revenue-leak-detective
ENV_FILE := .env.dev

# API Service (Go) Configuration
API_SERVICE_PATH := ./services/api
API_BINARY_NAME := rld-api
API_BINARY_PATH := ./bin/$(API_BINARY_NAME)
API_MAIN_PATH := ./cmd/main.go
API_DOCKER_IMAGE := rld-api

# Workers Service (Python) Configuration
WORKERS_SERVICE_PATH := ./workers
WORKERS_DOCKER_IMAGE := rld-workers

# Build metadata
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE) -w -s"

# Docker configuration
DOCKER_TAG ?= latest
DOCKER_COMPOSE_FILE := deploy/dev/docker-compose.yml

# Tool configuration
GOLINT_CONFIG := $(API_SERVICE_PATH)/.golangci.yml

# =============================================================================
# PHONY Target Declarations
# =============================================================================

.PHONY: help all clean validate-env check-tools
.PHONY: api-build api-build-local api-clean api-test api-test-coverage 
.PHONY: api-benchmark api-lint api-fmt api-fmt-check api-vet api-deps api-security api-all
.PHONY: workers-test workers-lint workers-format workers-install workers-all
.PHONY: api-docker-build api-docker-run workers-docker-build workers-docker-run
.PHONY: docker-build-all docker-compose-up docker-compose-down docker-compose-logs
.PHONY: dev install-tools
# Backward compatibility aliases
.PHONY: build build-local test test-coverage benchmark lint fmt fmt-check vet deps security
.PHONY: docker-build docker-run

# =============================================================================
# Default & Meta Targets
# =============================================================================

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
help:
	@printf "$(BOLD)$(BLUE)$(PROJECT_NAME) - Available Commands:$(NC)\n\n"
	@printf "$(BOLD)🏗️  Build Targets:$(NC)\n"
	@grep -E '^## (api-build|workers-|docker-).*:' $(MAKEFILE_LIST) | sed 's/^##/  /' | column -t -s ':'
	@printf "\n$(BOLD)🧪 Test & Quality:$(NC)\n"
	@grep -E '^## (api-test|api-lint|api-fmt|api-security|workers-test|workers-lint).*:' $(MAKEFILE_LIST) | sed 's/^##/  /' | column -t -s ':'
	@printf "\n$(BOLD)🚀 Development:$(NC)\n"
	@grep -E '^## (dev|install-tools|all|clean).*:' $(MAKEFILE_LIST) | sed 's/^##/  /' | column -t -s ':'

## all: Run complete CI pipeline (format, lint, test, build)
all: validate-env api-fmt api-vet api-lint api-test api-build

## clean: Clean all build artifacts and caches
clean: api-clean
	@printf "$(GREEN)✓ All artifacts cleaned$(NC)\n"

## validate-env: Validate environment and required files
validate-env:
	@printf "$(BLUE)Validating environment...$(NC)\n"
	@test -f "$(ENV_FILE)" || (printf "$(RED)❌ $(ENV_FILE) not found$(NC)\n" && exit 1)
	@test -d "$(API_SERVICE_PATH)" || (printf "$(RED)❌ API service path not found$(NC)\n" && exit 1)
	@test -d "$(WORKERS_SERVICE_PATH)" || (printf "$(RED)❌ Workers service path not found$(NC)\n" && exit 1)
	@printf "$(GREEN)✓ Environment validated$(NC)\n"

## check-tools: Verify required development tools are installed
check-tools:
	@printf "$(BLUE)Checking required tools...$(NC)\n"
	@command -v go >/dev/null || (printf "$(RED)❌ Go not installed$(NC)\n" && exit 1)
	@command -v docker >/dev/null || (printf "$(RED)❌ Docker not installed$(NC)\n" && exit 1)
	@command -v uv >/dev/null || printf "$(YELLOW)⚠️  UV not installed (required for workers development)$(NC)\n"
	@printf "$(GREEN)✓ Required tools available$(NC)\n"

# =============================================================================
# API Service Targets (Go)
# =============================================================================

## api-build: Build the API service for production (Linux)
api-build: validate-env
	@printf "$(BLUE)Building $(API_BINARY_NAME)...$(NC)\n"
	@mkdir -p bin
	cd $(API_SERVICE_PATH) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o ../../$(API_BINARY_PATH) $(API_MAIN_PATH)
	@printf "$(GREEN)✓ Build complete: $(API_BINARY_PATH)$(NC)\n"

## api-build-local: Build API service for local development
api-build-local: validate-env
	@printf "$(BLUE)Building $(API_BINARY_NAME) for local development...$(NC)\n"
	@mkdir -p bin
	cd $(API_SERVICE_PATH) && go build $(LDFLAGS) -o ../../$(API_BINARY_PATH) $(API_MAIN_PATH)
	@printf "$(GREEN)✓ Local build complete: $(API_BINARY_PATH)$(NC)\n"

## api-clean: Clean API service build artifacts
api-clean:
	@printf "$(YELLOW)🧹 Cleaning API build artifacts...$(NC)\n"
	cd $(API_SERVICE_PATH) && go clean
	@rm -rf bin/ dist/
	@printf "$(GREEN)✓ API clean complete$(NC)\n"

## api-test: Run tests for the API service
api-test:
	@printf "$(CYAN)🧪 Running API tests...$(NC)\n"
	@cd $(API_SERVICE_PATH) && go test -v -race -timeout 30s ./... && \
        printf "$(GREEN)✓ All tests passed$(NC)\n"

## api-test-coverage: Run API tests with coverage
api-test-coverage:
	@printf "$(CYAN)📊 Running API tests with coverage...$(NC)\n"
	cd $(API_SERVICE_PATH) && go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	cd $(API_SERVICE_PATH) && go tool cover -html=coverage.out -o coverage.html
	@printf "$(GREEN)✓ Coverage report: $(API_SERVICE_PATH)/coverage.html$(NC)\n"

## api-benchmark: Run benchmarks for the API service
api-benchmark:
	@printf "$(MAGENTA)⚡ Running API benchmarks...$(NC)\n"
	cd $(API_SERVICE_PATH) && go test -bench=. -benchmem ./...
	@printf "$(GREEN)✓ API benchmarks completed$(NC)\n"

## api-lint: Run linter for the API service
api-lint:
	@printf "$(CYAN)🔍 Running API linter...$(NC)\n"
	@cd $(API_SERVICE_PATH) && golangci-lint run --config .golangci.yml && \
		printf "$(GREEN)✓ API linting passed$(NC)\n"

## api-fmt: Format Go code in API service
api-fmt:
	@printf "$(YELLOW)🎨 Formatting API code...$(NC)\n"
	cd $(API_SERVICE_PATH) && gofmt -s -w .
	@printf "$(GREEN)✓ API code formatted$(NC)\n"

## api-fmt-check: Check if API code is formatted
api-fmt-check:
	@printf "$(YELLOW)📋 Checking API code formatting...$(NC)\n"
	@cd $(API_SERVICE_PATH) && test -z "$$(gofmt -l .)" || (echo "$(RED)❌ Code not formatted, run 'make api-fmt'$(NC)\n" && exit 1)
	@printf "$(GREEN)✓ API code is properly formatted$(NC)\n"

## api-vet: Run go vet for the API service
api-vet:
	@printf "$(YELLOW)🔎 Running API vet analysis...$(NC)\n"
	@cd $(API_SERVICE_PATH) && go vet ./... && \
		printf "$(GREEN)✓ API vet analysis passed$(NC)\n"

## api-deps: Download and verify API dependencies
api-deps:
	@printf "$(BLUE)📦 Managing API dependencies...$(NC)\n"
	cd $(API_SERVICE_PATH) && go mod download
	cd $(API_SERVICE_PATH) && go mod verify
	cd $(API_SERVICE_PATH) && go mod tidy
	@printf "$(GREEN)✓ API dependencies updated$(NC)\n"

## api-security: Run security scan for the API service
api-security:
	@printf "$(RED)🔒 Running API security scan...$(NC)\n"
	@cd $(API_SERVICE_PATH) && go run github.com/securego/gosec/v2/cmd/gosec@latest ./... && \
		printf "$(GREEN)✓ API security scan passed$(NC)\n"

## api-all: Run all API quality checks
api-all: api-fmt api-fmt-check api-vet api-lint api-test api-test-coverage api-deps api-security

# =============================================================================
# Workers Service Targets (Python)
# =============================================================================

## workers-install: Install workers dependencies
workers-install:
	@printf "$(BLUE)📦 Installing workers dependencies...$(NC)\n"
	cd $(WORKERS_SERVICE_PATH) && uv sync
	@printf "$(GREEN)✓ Workers dependencies installed$(NC)\n"

## workers-test: Run tests for the workers service
workers-test:
	@printf "$(CYAN)🧪 Running workers tests...$(NC)\n"
	@cd $(WORKERS_SERVICE_PATH) && uv run pytest tests/ -v && \
		printf "$(GREEN)✓ All workers tests passed$(NC)\n"

## workers-lint: Run linting for the workers service
workers-lint:
	@printf "$(CYAN)🔍 Running workers linter...$(NC)\n"
	@cd $(WORKERS_SERVICE_PATH) && uv run ruff check . && \
		printf "$(GREEN)✓ Workers linting passed$(NC)\n"

## workers-format: Format workers code
workers-format:
	@printf "$(YELLOW)🎨 Formatting workers code...$(NC)\n"
	cd $(WORKERS_SERVICE_PATH) && uv run ruff format .
	@printf "$(GREEN)✓ Workers code formatted$(NC)\n"

## workers-format-check: Check workers code formatting
workers-format-check:
	@printf "$(YELLOW)📋 Checking workers code formatting...$(NC)\n"
	@cd $(WORKERS_SERVICE_PATH) && uv run ruff format --check . && \
		printf "$(GREEN)✓ Workers code formatting is correct$(NC)\n"

## workers-run: Run the workers agent
workers-run:
	@printf "$(GREEN)🤖 Starting workers agent...$(NC)\n"
	cd $(WORKERS_SERVICE_PATH) && uv run python -m agent

## workers-run-dry: Run workers in dry-run mode
workers-run-dry:
	@printf "$(GREEN)🤖 Starting workers agent (dry-run)...$(NC)\n"
	cd $(WORKERS_SERVICE_PATH) && uv run python -m agent --dry-run

## workers-all: Run all workers quality checks
workers-all: workers-format workers-format-check workers-lint workers-test

# =============================================================================
# Docker Targets
# =============================================================================

## api-docker-build: Build Docker image for Go API service
api-docker-build:
	@printf "$(BLUE)Building Docker image for API service...$(NC)\n"
	docker build -t $(API_DOCKER_IMAGE):$(DOCKER_TAG) -f deploy/docker/Dockerfile.api \
		--build-arg VERSION=$(VERSION) \
		--build-arg COMMIT=$(COMMIT) \
		--build-arg DATE=$(DATE) .
	@printf "$(GREEN)✓ API Docker image built: $(API_DOCKER_IMAGE):$(DOCKER_TAG)$(NC)\n"

## api-docker-run: Run Docker container for Go API service
api-docker-run:
	@printf "$(BLUE)Running API Docker container...$(NC)\n"
	docker run -p 8080:8080 --env-file .env.dev $(API_DOCKER_IMAGE):$(DOCKER_TAG)

## workers-docker-build: Build Docker image for Python workers service
workers-docker-build:
	@printf "$(BLUE)Building Docker image for workers service...$(NC)\n"
	docker build -t rld-workers:$(DOCKER_TAG) -f deploy/docker/Dockerfile.workers .
	@printf "$(GREEN)✓ Workers Docker image built: rld-workers:$(DOCKER_TAG)$(NC)\n"

## workers-docker-run: Run Docker container for workers service
workers-docker-run:
	@printf "$(BLUE)Running workers Docker container...$(NC)\n"
	docker run --env-file .env.dev rld-workers:$(DOCKER_TAG)

## docker-build-all: Build both Docker images
docker-build-all: api-docker-build workers-docker-build

## docker-compose-up: Start all services with docker-compose
docker-compose-up:
	@printf "$(BLUE)Starting all services with docker-compose...$(NC)\n"
	export VERSION=$(VERSION) COMMIT=$(COMMIT) DATE=$(DATE) && \
	docker-compose -f deploy/dev/docker-compose.yml up --build

## docker-compose-down: Stop all services
docker-compose-down:
	@printf "$(YELLOW)Stopping all services...$(NC)\n"
	docker-compose -f deploy/dev/docker-compose.yml down

## docker-compose-logs: View logs from all services
docker-compose-logs:
	@printf "$(CYAN)Viewing logs from all services...$(NC)\n"
	docker-compose -f deploy/dev/docker-compose.yml logs -f

# =============================================================================
# Development & Utility Targets
# =============================================================================

## dev: Run development server with hot reload
dev: validate-env
	@printf "$(GREEN)🚀 Starting development server with hot reload...$(NC)\n"
	@command -v air >/dev/null 2>&1 || (echo "$(RED)❌ air not installed. Run: go install github.com/cosmtrek/air@latest$(NC)\n" && exit 1)
	cd $(API_SERVICE_PATH) && air

## install-tools: Install required development tools
install-tools:
	@printf "$(MAGENTA)📦 Installing development tools...$(NC)\n"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@go install github.com/cosmtrek/air@latest
	@printf "$(GREEN)✓ Development tools installed$(NC)\n"

# =============================================================================
# Backward Compatibility Aliases
# =============================================================================

# Legacy build targets
build: api-build
build-local: api-build-local

# Legacy test & quality targets
test: api-test
test-coverage: api-test-coverage
benchmark: api-benchmark
lint: api-lint
fmt: api-fmt
fmt-check: api-fmt-check
vet: api-vet
deps: api-deps
security: api-security

# Legacy docker targets
docker-build: api-docker-build
docker-run: api-docker-run
