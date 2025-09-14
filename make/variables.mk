# =============================================================================
# Shared Variables for Revenue Leak Detective
# =============================================================================

# Shell configuration
SHELL := /bin/sh
.SHELLFLAGS := -ec

# Colors for output
RED := \033[31m
GREEN := \033[32m
YELLOW := \033[33m
BLUE := \033[34m
MAGENTA := \033[35m
CYAN := \033[36m
BOLD := \033[1m
NC := \033[0m

# Project metadata
PROJECT_NAME := revenue-leak-detective
ENV_FILE := .env.dev

# API Service (Go) Configuration
API_SERVICE_PATH := ./services/api
API_BINARY_NAME := rld-api
API_BINARY_PATH := ./bin/$(API_BINARY_NAME)
API_MAIN_PATH := $(API_SERVICE_PATH)/cmd
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
DOCKER_COMPOSE_CMD ?= docker compose

# Tool configuration
GOLINT_CONFIG := $(API_SERVICE_PATH)/.golangci.yml
