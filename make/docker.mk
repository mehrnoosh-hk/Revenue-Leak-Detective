# =============================================================================
# Docker & Deployment Targets
# =============================================================================

# Docker PHONY declarations
.PHONY: api-docker-build api-docker-run workers-docker-build workers-docker-run
.PHONY: docker-build-all docker-compose-up docker-compose-down docker-compose-logs

## api-docker-build: Build Docker image for Go API service
api-docker-build:
	@printf "$(BLUE)Building Docker image for API service...$(NC)\n"
	@test -f "$(ENV_FILE)" || (printf "$(RED)❌ $(ENV_FILE) not found$(NC)\n" && exit 1)
	@. $(ENV_FILE) && \
		RESOLVED_IMAGE="$${API_DOCKER_IMAGE:-$(API_DOCKER_IMAGE)}" && \
		RESOLVED_TAG="$${DOCKER_TAG:-$(DOCKER_TAG)}" && \
		docker build -t $$RESOLVED_IMAGE:$$RESOLVED_TAG -f deploy/docker/Dockerfile.api \
			--build-arg VERSION=$(VERSION) \
			--build-arg COMMIT=$(COMMIT) \
			--build-arg DATE=$(DATE) . && \
		printf "$(GREEN)✓ API Docker image built: $$RESOLVED_IMAGE:$$RESOLVED_TAG$(NC)\n"

## api-docker-run: Run Docker container for Go API service
api-docker-run:
	@printf "$(BLUE)Running API Docker container...$(NC)\n"
	@test -f "$(ENV_FILE)" || (printf "$(RED)❌ $(ENV_FILE) not found$(NC)\n" && exit 1)
	@. $(ENV_FILE) && \
		docker run -p 8080:8080 --env-file $(ENV_FILE) $${API_DOCKER_IMAGE:-$(API_DOCKER_IMAGE)}:$${DOCKER_TAG:-$(DOCKER_TAG)}

## workers-docker-build: Build Docker image for Python workers service
workers-docker-build:
	@printf "$(BLUE)Building Docker image for workers service...$(NC)\n"
	@test -f "$(ENV_FILE)" || (printf "$(RED)❌ $(ENV_FILE) not found$(NC)\n" && exit 1)
	@. $(ENV_FILE) && \
		RESOLVED_IMAGE="$${WORKERS_DOCKER_IMAGE:-$(WORKERS_DOCKER_IMAGE)}" && \
		RESOLVED_TAG="$${DOCKER_TAG:-$(DOCKER_TAG)}" && \
		docker build -t $$RESOLVED_IMAGE:$$RESOLVED_TAG -f deploy/docker/Dockerfile.workers . && \
		printf "$(GREEN)✓ Workers Docker image built: $$RESOLVED_IMAGE:$$RESOLVED_TAG$(NC)\n"

## workers-docker-run: Run Docker container for workers service
workers-docker-run:
	@printf "$(BLUE)Running workers Docker container...$(NC)\n"
	@test -f "$(ENV_FILE)" || (printf "$(RED)❌ $(ENV_FILE) not found$(NC)\n" && exit 1)
	@. $(ENV_FILE) && \
		docker run --env-file $(ENV_FILE) $${WORKERS_DOCKER_IMAGE:-$(WORKERS_DOCKER_IMAGE)}:$${DOCKER_TAG:-$(DOCKER_TAG)}

## docker-build-all: Build both Docker images
docker-build-all: api-docker-build workers-docker-build

## docker-compose-up: Start all services with docker-compose
docker-compose-up:
	@printf "$(BLUE)Starting all services with docker-compose...$(NC)\n"
	@test -f "$(ENV_FILE)" || (printf "$(RED)❌ $(ENV_FILE) not found$(NC)\n" && exit 1)
	@. $(ENV_FILE) && \
		export VERSION=$(VERSION) COMMIT=$(COMMIT) DATE=$(DATE) \
			POSTGRES_DB="$$POSTGRES_DB" \
			POSTGRES_USER="$$POSTGRES_USER" \
			POSTGRES_PASSWORD="$$POSTGRES_PASSWORD" && \
		docker-compose -f $(DOCKER_COMPOSE_FILE) up --build

## docker-compose-down: Stop all services
docker-compose-down:
	@printf "$(YELLOW)Stopping all services...$(NC)\n"
	docker-compose -f $(DOCKER_COMPOSE_FILE) down

## docker-compose-logs: View logs from all services
docker-compose-logs:
	@printf "$(CYAN)Viewing logs from all services...$(NC)\n"
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f
