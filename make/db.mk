# =============================================================================
# Database Targets
# =============================================================================

# Database PHONY declarations
.PHONY: migrate-up migrate-down migrate-create migrate-version migrate-force
.PHONY: migrate-up-step migrate-down-step db-reset sqlc sqlc-check
.PHONY: _validate-env _validate-postgres-url _validate-dev-env

# =============================================================================
# Helper Functions
# =============================================================================

# Validate environment file exists
_validate-env:
	@test -f "$(ENV_FILE)" || (printf "$(RED)❌ $(ENV_FILE) not found. Create it from .env.example$(NC)\n" && exit 1)

# Validate POSTGRES_URL is set
_validate-postgres-url: _validate-env
	@. $(ENV_FILE) && \
		test -n "$$POSTGRES_URL" || (printf "$(RED)❌ POSTGRES_URL not set in $(ENV_FILE)$(NC)\n" && exit 1)

# Validate environment is development/dev for destructive operations
_validate-dev-env: _validate-postgres-url
	@. $(ENV_FILE) && \
		test -n "$$ENVIRONMENT" || (printf "$(RED)❌ ENVIRONMENT not set in $(ENV_FILE)$(NC)\n" && exit 1) && \
		(echo "$$ENVIRONMENT" | grep -E "^(development|dev)$$" > /dev/null) || (printf "$(RED)❌ This command is only allowed in development or dev environment. Current: $$ENVIRONMENT$(NC)\n" && exit 1)

# =============================================================================
# Migration Commands
# =============================================================================

## migrate-up: Apply all database migrations
migrate-up: _validate-postgres-url
	@printf "$(BLUE)⬆️  Applying database migrations...$(NC)\n"
	@. $(ENV_FILE) && \
		cd $(API_SERVICE_PATH) && \
		migrate -path ./migrations -database "$$POSTGRES_URL" up && \
		printf "$(GREEN)✓ Database migrations applied$(NC)\n"

## migrate-down: Rollback last database migration
migrate-down: _validate-postgres-url
	@printf "$(YELLOW)⬇️  Rolling back last database migration...$(NC)\n"
	@. $(ENV_FILE) && \
		cd $(API_SERVICE_PATH) && \
		migrate -path ./migrations -database "$$POSTGRES_URL" down -all && \
		printf "$(GREEN)✓ Last migration rolled back$(NC)\n"

## migrate-up-step: Apply specific number of migration steps (usage: make migrate-up-step STEPS=1)
migrate-up-step: _validate-postgres-url
	@test -n "$(STEPS)" || (printf "$(RED)❌ STEPS is required. Usage: make migrate-up-step STEPS=1$(NC)\n" && exit 1)
	@printf "$(BLUE)⬆️  Applying $(STEPS) migration step(s)...$(NC)\n"
	@. $(ENV_FILE) && \
		cd $(API_SERVICE_PATH) && \
		migrate -path ./migrations -database "$$POSTGRES_URL" up $(STEPS) && \
		printf "$(GREEN)✓ $(STEPS) migration step(s) applied$(NC)\n"

## migrate-down-step: Rollback specific number of migration steps (usage: make migrate-down-step STEPS=1)
migrate-down-step: _validate-postgres-url
	@test -n "$(STEPS)" || (printf "$(RED)❌ STEPS is required. Usage: make migrate-down-step STEPS=1$(NC)\n" && exit 1)
	@printf "$(YELLOW)⬇️  Rolling back $(STEPS) migration step(s)...$(NC)\n"
	@. $(ENV_FILE) && \
		cd $(API_SERVICE_PATH) && \
		migrate -path ./migrations -database "$$POSTGRES_URL" down $(STEPS) && \
		printf "$(GREEN)✓ $(STEPS) migration step(s) rolled back$(NC)\n"

## migrate-version: Show current migration version
migrate-version: _validate-postgres-url
	@printf "$(BLUE)📊 Checking migration version...$(NC)\n"
	@. $(ENV_FILE) && \
		cd $(API_SERVICE_PATH) && \
		migrate -path ./migrations -database "$$POSTGRES_URL" version

## migrate-force: Force migration to specific version (usage: make migrate-force VERSION=1)
migrate-force: _validate-postgres-url
	@test -n "$(VERSION)" || (printf "$(RED)❌ VERSION is required. Usage: make migrate-force VERSION=1$(NC)\n" && exit 1)
	@printf "$(RED)⚠️  Force setting migration version to $(VERSION)...$(NC)\n"
	@. $(ENV_FILE) && \
		cd $(API_SERVICE_PATH) && \
		migrate -path ./migrations -database "$$POSTGRES_URL" force $(VERSION) && \
		printf "$(GREEN)✓ Migration version forced to $(VERSION)$(NC)\n"

# =============================================================================
# Migration File Management
# =============================================================================

## migrate-create: Create a new migration file (usage: make migrate-create NAME=migration_name)
migrate-create:
	@test -n "$(NAME)" || (printf "$(RED)❌ NAME is required. Usage: make migrate-create NAME=migration_name$(NC)\n" && exit 1)
	@printf "$(BLUE)📝 Creating migration: $(NAME)...$(NC)\n"
	@cd $(API_SERVICE_PATH) && \
		migrate create -ext sql -dir ./migrations -seq $(NAME) && \
		printf "$(GREEN)✓ Migration files created in ./migrations/$(NC)\n"

# =============================================================================
# SQLC Generation and Management
# =============================================================================

## sqlc-generate: Generate SQLC code from SQL queries
sqlc-generate:
	@printf "$(BLUE)📝 Generating SQLC code...$(NC)\n"
	@cd $(API_SERVICE_PATH) && sqlc generate
	@printf "$(GREEN)✓ SQLC code generated$(NC)\n"

## sqlc-check: Check if generated SQLC code matches committed code (matches CI workflow)
sqlc-check:
	@printf "$(BLUE)🔍 Checking SQLC code synchronization...$(NC)\n"
	@cd $(API_SERVICE_PATH) && \
		if ! git diff --quiet -- ./internal/db/sqlc/; then \
			printf "$(RED)❌ Generated sqlc code is out of sync with committed code.$(NC)\n"; \
			printf "$(YELLOW)Please run 'make sqlc-generate' and commit the changes.$(NC)\n\n"; \
			printf "$(BLUE)Differences found in:$(NC)\n"; \
			git diff --name-only -- ./internal/db/sqlc/; \
			printf "\n$(BLUE)Full diff:$(NC)\n"; \
			git diff -- ./internal/db/sqlc/; \
			exit 1; \
		else \
			printf "$(GREEN)✅ sqlc generated code is up to date$(NC)\n"; \
		fi

# =============================================================================
# Database Management
# =============================================================================

## db-reset: Drop and recreate database (development/dev environments only)
db-reset: _validate-dev-env
	@printf "$(RED)⚠️  Dropping and recreating database...$(NC)\n"
	@. $(ENV_FILE) && \
		cd $(API_SERVICE_PATH) && \
		migrate -path ./migrations -database "$$POSTGRES_URL" drop -f && \
		migrate -path ./migrations -database "$$POSTGRES_URL" up && \
		printf "$(GREEN)✓ Database dropped and recreated$(NC)\n"
