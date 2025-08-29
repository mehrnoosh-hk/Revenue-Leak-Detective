# =============================================================================
# API Service Targets (Go)
# =============================================================================

# API PHONY declarations
.PHONY: api-build api-build-local api-clean api-test api-test-coverage 
.PHONY: api-benchmark api-lint api-fmt api-fmt-check api-vet api-deps api-security api-all
.PHONY: api-run

## api-build: Build the API service for production (Linux)
api-build: validate-env
	@printf "$(BLUE)Building $(API_BINARY_NAME)...$(NC)\n"
	@mkdir -p $(dir $(API_BINARY_PATH))
	cd $(API_SERVICE_PATH) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o ../../$(API_BINARY_PATH) $(API_MAIN_PATH)
	@printf "$(GREEN)✓ Build complete: $(API_BINARY_PATH)$(NC)\n"

## api-build-local: Build API service for local development
api-build-local: validate-env
	@printf "$(BLUE)Building $(API_BINARY_NAME) for local development...$(NC)\n"
	@mkdir -p $(dir $(API_BINARY_PATH))
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



## api-run: Run the API service
api-run:
	@printf "$(GREEN) Running the API server...$(NC)\n"
	@cd $(API_SERVICE_PATH) && go run ./cmd/main.go

## api-all: Run all API quality checks
api-all: api-fmt api-fmt-check api-vet api-lint api-test api-test-coverage api-deps api-security
