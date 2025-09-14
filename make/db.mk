# =============================================================================
# Database Targets
# =============================================================================

# Database PHONY declarations
.PHONY: migrate-up migrate-down migrate-create migrate-version migrate-force
.PHONY: migrate-up-step migrate-down-step db-reset sqlc sqlc-check domain-models-generate
.PHONY: migrate-check _validate-env _validate-postgres-url _validate-dev-env

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
		printf "$(GREEN)✓ All migrations rolled back$(NC)\n"

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

## migrate-check: Validate migrations by applying and rolling back (matches CI workflow)
migrate-check:
	@printf "$(BLUE)🔍 Validating database migrations...$(NC)\n"
	@printf "$(YELLOW)⚠️  This will create a temporary test database and apply/rollback all migrations$(NC)\n"
	@printf "$(BLUE)📋 Checking if migrate CLI is available...$(NC)\n"
	@which migrate > /dev/null || (printf "$(RED)❌ migrate CLI not found. Please install it first.$(NC)\n" && exit 1)
	@printf "$(GREEN)✓ migrate CLI found$(NC)\n"
	@printf "$(BLUE)📋 Creating temporary test database...$(NC)\n"
	@cd $(API_SERVICE_PATH) && \
		TEMP_DB_NAME="rld_migrate_check_$$(date +%s)" && \
		TEMP_DB_URL="postgres://postgres:postgres@localhost:5432/$$TEMP_DB_NAME?sslmode=disable" && \
		printf "$(BLUE)📊 Test database: $$TEMP_DB_NAME$(NC)\n" && \
		createdb "$$TEMP_DB_NAME" 2>/dev/null || (printf "$(RED)❌ Failed to create test database. Make sure PostgreSQL is running and accessible.$(NC)\n" && exit 1) && \
		printf "$(GREEN)✓ Test database created$(NC)\n" && \
		printf "$(BLUE)⬆️  Applying all migrations...$(NC)\n" && \
		migrate -path ./migrations -database "$$TEMP_DB_URL" up && \
		printf "$(GREEN)✓ All migrations applied successfully$(NC)\n" && \
		printf "$(BLUE)⬇️  Rolling back all migrations...$(NC)\n" && \
		migrate -path ./migrations -database "$$TEMP_DB_URL" down -all && \
		printf "$(GREEN)✓ All migrations rolled back successfully$(NC)\n" && \
		printf "$(BLUE)🧹 Cleaning up test database...$(NC)\n" && \
		dropdb "$$TEMP_DB_NAME" 2>/dev/null && \
		printf "$(GREEN)✓ Test database cleaned up$(NC)\n" && \
		printf "$(GREEN)✅ Migration validation completed successfully$(NC)\n"

# =============================================================================
# Migration File Management
# =============================================================================

## migrate-create: Create a new migration file with timestamp (usage: make migrate-create name=migration_name)
## Note: Uses timestamp format by default (YYYYMMDDHHMMSS). For custom format, modify the command below.
migrate-create:
	@test -n "$(name)" || (printf "$(RED)❌ name is required. Usage: make migrate-create name=migration_name$(NC)\n" && exit 1)
	@printf "$(BLUE)📝 Creating migration: $(name)...$(NC)\n"
	@cd $(API_SERVICE_PATH) && \
		migrate create -ext sql -seq -digits 3 -dir ./migrations $(name) && \
		printf "$(GREEN)✓ Migration files created in ./migrations/$(NC)\n"

# =============================================================================
# SQLC Generation and Management
# =============================================================================

## sqlc-generate: Generate SQLC code from SQL queries
sqlc-generate:
	@sqlc version | grep -q "1.30.0" || (printf "$(RED)❌ SQLC version is not 1.30.0. Please update SQLC to 1.30.0.$(NC)\n" && exit 1)
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

## domain-models-generate: Generate domain models from SQLC models
domain-models-generate:
	@printf "$(BLUE)📝 Generating domain models...$(NC)\n"
	@cd $(API_SERVICE_PATH) && go run scripts/generate_domain_models.go
	@printf "$(GREEN)✓ Domain models generated$(NC)\n"

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
