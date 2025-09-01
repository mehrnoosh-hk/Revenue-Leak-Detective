# =============================================================================
# Tools & Development Targets
# =============================================================================

# Tools PHONY declarations
.PHONY: install-tools check-tools dev validate-env install-workers-tools check-workers-tools

## install-tools: Install required development tools
install-tools:
	@printf "$(MAGENTA)📦 Installing development tools...$(NC)\n"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@go install github.com/cosmtrek/air@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@printf "$(GREEN)✓ Development tools installed$(NC)\n"

## install-workers-tools: Install only workers-related development tools
install-workers-tools:
	@printf "$(MAGENTA)📦 Installing workers development tools...$(NC)\n"
	@command -v uv >/dev/null && (printf "$(GREEN)✓ UV already installed$(NC)\n") || ( \
		printf "$(BLUE)📦 Installing UV...$(NC)\n" && \
		command -v curl >/dev/null || (printf "$(RED)❌ curl not installed (required for uv installation)$(NC)\n" && exit 1) && \
		curl -LsSf https://astral.sh/uv/install.sh | sh && \
		printf "$(GREEN)✓ UV installed$(NC)\n" && \
		printf "$(YELLOW)⚠️  Please restart your shell or run: source ~/.bashrc$(NC)\n" \
	)
	@printf "$(BLUE)📦 Checking workers project dependencies...$(NC)\n"
	@export PATH="$$HOME/.local/bin:$$PATH" && cd $(WORKERS_SERVICE_PATH) && \
		(uv run ruff --version >/dev/null 2>&1 && uv run pytest --version >/dev/null 2>&1) && \
		printf "$(GREEN)✓ Workers dependencies already installed$(NC)\n" || ( \
		printf "$(BLUE)📦 Installing workers project dependencies...$(NC)\n" && \
		uv sync --group dev && \
		printf "$(GREEN)✓ Workers dependencies installed$(NC)\n" \
	)
	@printf "$(GREEN)✓ All workers development tools ready$(NC)\n"

## check-tools: Verify required development tools are installed
check-tools:
	@printf "$(BLUE)Checking required tools...$(NC)\n"
	@command -v go >/dev/null && printf "$(GREEN)✓ Go available$(NC)\n" || (printf "$(RED)❌ Go not installed$(NC)\n" && exit 1)
	@command -v docker >/dev/null && printf "$(GREEN)✓ Docker available$(NC)\n" || (printf "$(RED)❌ Docker not installed$(NC)\n" && exit 1)
	@command -v sqlc >/dev/null && printf "$(GREEN)✓ sqlc available$(NC)\n" || (printf "$(RED)❌ sqlc not installed$(NC)\n" && exit 1)
	@command -v migrate >/dev/null && printf "$(GREEN)✓ golang-migrate available$(NC)\n" || (printf "$(RED)❌ golang-migrate not installed$(NC)\n" && exit 1)
	@command -v uv >/dev/null && printf "$(GREEN)✓ UV available$(NC)\n" || printf "$(YELLOW)⚠️  UV not installed (required for workers development)$(NC)\n"
	@printf "$(GREEN)✓ Required tools available$(NC)\n"

## check-workers-tools: Verify workers-related development tools are installed
check-workers-tools:
	@printf "$(BLUE)Checking workers development tools...$(NC)\n"
	@command -v uv >/dev/null || (printf "$(RED)❌ UV not installed. Run: make install-workers-tools$(NC)\n" && exit 1)
	@printf "$(GREEN)✓ UV available$(NC)\n"
	@command -v python3 >/dev/null || (printf "$(RED)❌ Python 3 not installed$(NC)\n" && exit 1)
	@python3 --version | grep -q "Python 3\.[1-9][2-9]" || (printf "$(RED)❌ Python 3.12+ required (found: $(shell python3 --version))$(NC)\n" && exit 1)
	@printf "$(GREEN)✓ Python 3.12+ available$(NC)\n"
	@test -d "$(WORKERS_SERVICE_PATH)" || (printf "$(RED)❌ Workers service directory not found$(NC)\n" && exit 1)
	@printf "$(GREEN)✓ Workers service directory found$(NC)\n"
	@test -f "$(WORKERS_SERVICE_PATH)/pyproject.toml" || (printf "$(RED)❌ Workers pyproject.toml not found$(NC)\n" && exit 1)
	@printf "$(GREEN)✓ Workers pyproject.toml found$(NC)\n"
	@cd $(WORKERS_SERVICE_PATH) && uv run ruff --version >/dev/null 2>&1 || (printf "$(RED)❌ ruff not available in workers environment. Run: make install-workers-tools$(NC)\n" && exit 1)
	@printf "$(GREEN)✓ ruff available in workers environment$(NC)\n"
	@cd $(WORKERS_SERVICE_PATH) && uv run pytest --version >/dev/null 2>&1 || (printf "$(RED)❌ pytest not available in workers environment. Run: make install-workers-tools$(NC)\n" && exit 1)
	@printf "$(GREEN)✓ pytest available in workers environment$(NC)\n"
	@cd $(WORKERS_SERVICE_PATH) && uv run python -c "import sys; print('Python path:', sys.executable)" >/dev/null 2>&1 || (printf "$(RED)❌ UV virtual environment not properly configured. Run: make install-workers-tools$(NC)\n" && exit 1)
	@printf "$(GREEN)✓ UV virtual environment properly configured$(NC)\n"
	@printf "$(GREEN)✓ All workers development tools available$(NC)\n"

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
