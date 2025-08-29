# =============================================================================
# Workers Service Targets (Python)
# =============================================================================

# Workers PHONY declarations
.PHONY: workers-deps workers-test workers-lint workers-format workers-format-check
.PHONY: workers-run workers-run-dry workers-all

## workers-deps: Install workers dependencies
workers-deps: check-workers-tools
	@printf "$(BLUE)📦 Installing workers dependencies...$(NC)\n"
	cd $(WORKERS_SERVICE_PATH) && uv sync
	@printf "$(GREEN)✓ Workers dependencies installed$(NC)\n"

## workers-test: Run tests for the workers service
workers-test: check-workers-tools
	@printf "$(CYAN)🧪 Running workers tests...$(NC)\n"
	@cd $(WORKERS_SERVICE_PATH) && uv run pytest tests/ -v && \
		printf "$(GREEN)✓ All workers tests passed$(NC)\n"

## workers-lint: Run linting for the workers service
workers-lint: check-workers-tools
	@printf "$(CYAN)🔍 Running workers linter...$(NC)\n"
	@cd $(WORKERS_SERVICE_PATH) && uv run ruff check . && \
		printf "$(GREEN)✓ Workers linting passed$(NC)\n"

## workers-format: Format workers code
workers-format: check-workers-tools
	@printf "$(YELLOW)🎨 Formatting workers code...$(NC)\n"
	cd $(WORKERS_SERVICE_PATH) && uv run ruff format .
	@printf "$(GREEN)✓ Workers code formatted$(NC)\n"

## workers-format-check: Check workers code formatting
workers-format-check: check-workers-tools
	@printf "$(YELLOW)📋 Checking workers code formatting...$(NC)\n"
	@cd $(WORKERS_SERVICE_PATH) && uv run ruff format --check . && \
		printf "$(GREEN)✓ Workers code formatting is correct$(NC)\n"

## workers-run: Run the workers agent
workers-run: check-workers-tools
	@printf "$(GREEN)🤖 Starting workers agent...$(NC)\n"
	cd $(WORKERS_SERVICE_PATH) && uv run python -m agent

## workers-run-dry: Run workers in dry-run mode
workers-run-dry: check-workers-tools
	@printf "$(GREEN)🤖 Starting workers agent (dry-run)...$(NC)\n"
	cd $(WORKERS_SERVICE_PATH) && uv run python -m agent --dry-run

## workers-all: Run all workers quality checks
workers-all: workers-format workers-format-check workers-lint workers-test
