# =============================================================================
# Revenue Leak Detective - Modular Makefile
# =============================================================================

# Include shared variables and modules
include make/variables.mk
include make/tools.mk
include make/api.mk
include make/db.mk
include make/workers.mk
include make/docker.mk

# =============================================================================
# PHONY Target Declarations (Main Targets Only)
# =============================================================================

.PHONY: help all clean deps

# =============================================================================
# Main Targets
# =============================================================================

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
help:
	@printf "$(BOLD)$(BLUE)$(PROJECT_NAME) - Available Commands:$(NC)\n\n"
	@printf "$(BOLD)🚀 Quick Start:$(NC)\n"
	@printf "  make install-tools       Install development tools\n"
	@printf "  make deps               Install dependencies for all services\n"
	@printf "  make dev                Start development server with hot reload\n"
	@printf "  make all                Run complete CI pipeline\n"
	@printf "\n$(BOLD)🏗️  Build & Test:$(NC)\n"
	@grep -E '^## (api-build|api-test|workers-test).*:' make/*.mk | sed 's/^[^:]*://; s/^##/  /' | column -t -s ':'
	@printf "\n$(BOLD)📊 Database:$(NC)\n"
	@grep -E '^## (migrate|sqlc|db-reset).*:' make/*.mk | sed 's/^[^:]*://; s/^##/  /' | column -t -s ':'
	@printf "\n$(BOLD)🐳 Docker:$(NC)\n"
	@grep -E '^## (docker-).*:' make/*.mk | sed 's/^[^:]*://; s/^##/  /' | column -t -s ':'
	@printf "\n$(BOLD)⚡ Quality:$(NC)\n"
	@grep -E '^## (api-lint|api-fmt|workers-lint).*:' make/*.mk | sed 's/^[^:]*://; s/^##/  /' | column -t -s ':'
	@printf "\n$(BOLD)🔧 Tools:$(NC)\n"
	@grep -E '^## (install-tools|check-tools|validate-env).*:' make/*.mk | sed 's/^[^:]*://; s/^##/  /' | column -t -s ':'

## all: Run complete CI pipeline (format, lint, test, build)
all: validate-env api-fmt api-vet api-lint api-test api-build workers-all

## clean: Clean all build artifacts and caches
clean: api-clean
	@printf "$(GREEN)✓ All artifacts cleaned$(NC)\n"

## deps: Install dependencies for both services
deps: api-deps workers-deps

