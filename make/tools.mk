# =============================================================================
# Tools & Development Targets
# =============================================================================

# Tools PHONY declarations
.PHONY: install-tools check-tools dev validate-env

## install-tools: Install required development tools
install-tools:
	@printf "$(MAGENTA)📦 Installing development tools...$(NC)\n"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@go install github.com/cosmtrek/air@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@printf "$(GREEN)✓ Development tools installed$(NC)\n"

## check-tools: Verify required development tools are installed
check-tools:
	@printf "$(BLUE)Checking required tools...$(NC)\n"
	@command -v go >/dev/null || (printf "$(RED)❌ Go not installed$(NC)\n" && exit 1)
	@command -v docker >/dev/null || (printf "$(RED)❌ Docker not installed$(NC)\n" && exit 1)
	@command -v sqlc >/dev/null || (printf "$(RED)❌ sqlc not installed$(NC)\n" && exit 1)
	@command -v migrate >/dev/null || (printf "$(RED)❌ golang-migrate not installed$(NC)\n" && exit 1)
	@command -v uv >/dev/null || printf "$(YELLOW)⚠️  UV not installed (required for workers development)$(NC)\n"
	@printf "$(GREEN)✓ Required tools available$(NC)\n"

## validate-env: Validate environment and required files
validate-env:
	@printf "$(BLUE)Validating environment...$(NC)\n"
	@test -f "$(ENV_FILE)" || (printf "$(RED)❌ $(ENV_FILE) not found$(NC)\n" && exit 1)
	@test -d "$(API_SERVICE_PATH)" || (printf "$(RED)❌ API service path not found$(NC)\n" && exit 1)
	@test -d "$(WORKERS_SERVICE_PATH)" || (printf "$(RED)❌ Workers service path not found$(NC)\n" && exit 1)
	@printf "$(GREEN)✓ Environment validated$(NC)\n"

## dev: Run development server with hot reload
dev: validate-env
	@printf "$(GREEN)🚀 Starting development server with hot reload...$(NC)\n"
	@command -v air >/dev/null 2>&1 || (echo "$(RED)❌ air not installed. Run: go install github.com/cosmtrek/air@latest$(NC)\n" && exit 1)
	cd $(API_SERVICE_PATH) && air
