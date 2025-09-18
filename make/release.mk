.PHONY: release-check release-create build-info-print

## release-check: Check if ready for release
release-check:
	@printf "$(BLUE)🔍 Checking release readiness...$(NC)\n"
	@# Check if we're on main branch
	@current_branch=$$(git rev-parse --abbrev-ref HEAD); \
	if [ "$$current_branch" != "main" ]; then \
		printf "$(RED)❌ Must be on main branch to create release (currently on $$current_branch)$(NC)\n"; \
		exit 1; \
	fi
	@# Check if working directory is clean
	@if [ -n "$$(git status --porcelain)" ]; then \
		printf "$(RED)❌ Working directory is not clean$(NC)\n"; \
		exit 1; \
	fi
	@# Check if tests pass
	@printf "$(YELLOW)�� Running tests...$(NC)\n"
	@cd $(API_SERVICE_PATH) && go test ./...
	@printf "$(GREEN)✓ All checks passed$(NC)\n"

## release-create: Create a new release (usage: make release-create VERSION=v1.2.0)
release-create: release-check
	@test -n "$(VERSION)" || (printf "$(RED)❌ VERSION is required. Usage: make release-create VERSION=v1.2.0$(NC)\n" && exit 1)
	@printf "$(GREEN)🚀 Creating release $(VERSION)...$(NC)\n"
	@# Create annotated tag
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@# Push tag to remote
	@git push origin $(VERSION)
	@printf "$(GREEN)✓ Release $(VERSION) created and pushed$(NC)\n"
	@printf "$(BLUE)📋 Next steps:$(NC)\n"
	@printf "  1. Update release notes on GitHub/GitLab$(NC)\n"
	@printf "  2. Deploy to staging environment$(NC)\n"
	@printf "  3. Run integration tests$(NC)\n"
	@printf "  4. Deploy to production$(NC)\n"

build-info-print:
	@printf "$(BLUE)🔍 Build info:..$(NC)\n"
	@printf "$(GREEN)Version: $(VERSION)$(NC)\n"
	@printf "$(GREEN)Commit: $(COMMIT)$(NC)\n"
	@printf "$(GREEN)Build date: $(DATE)$(NC)\n"