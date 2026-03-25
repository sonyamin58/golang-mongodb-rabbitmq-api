.PHONY: help build run test clean migrate docker-build docker-up docker-down lint fmt deps

# Variables
APP_NAME=golang-mongodb-rabbitmq-api
DOCKER_IMAGE=golib-api
DOCKER_TAG=latest
GO=go
GOFLAGS=-v
BUILD_DIR=./bin
MAIN_FILE=./cmd/api/main.go

# Colors
RED   := \033[0;31m
GREEN := \033[0;32m
YELLOW:= \033[1;33m
NC    := \033[0m

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}'

deps: ## Download Go dependencies
	$(GO) mod download
	$(GO) mod tidy

build: deps ## Build Go binary
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/api $(MAIN_FILE)
	@echo "$(GREEN)Build successful: $(BUILD_DIR)/api$(NC)"

run: build ## Build and run the API server
	$(BUILD_DIR)/api

dev: ## Run in development mode with hot reload (requires air)
	air -c .air.toml

test: ## Run all tests
	$(GO) test -v -cover -race ./...

test-coverage: ## Run tests with coverage report
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Coverage report: coverage.html$(NC)"

lint: ## Run linters (golangci-lint)
	golangci-lint run ./...

fmt: ## Format code
	$(GO) fmt ./...
	goimports -w .

clean: ## Clean build artifacts
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Database
migrate: ## Run database migrations (requires Oracle connection)
	@echo "$(YELLOW)Running migrations...$(NC)"
	@sqlplus $(DB_USER)/$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_SERVICE) @migrations/001_init_schema.sql

migrate-check: ## Check migration status
	@sqlplus $(DB_USER)/$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_SERVICE) @migrations/check.sql

# Docker
docker-build: ## Build Docker images
	docker build --tag $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-up: ## Start all services with docker-compose
	docker-compose up -d
	@echo "$(GREEN)Services started. API: http://localhost:8080$(NC)"

docker-down: ## Stop all services
	docker-compose down

docker-logs: ## Show docker-compose logs
	docker-compose logs -f

docker-ps: ## Show running containers
	docker-compose ps

# Celery
celery-worker: ## Run Celery worker
	cd workers && celery -A celery_app worker --loglevel=info --concurrency=4

celery-beat: ## Run Celery beat scheduler
	cd workers && celery -A celery_app beat --loglevel=info

celery-flower: ## Run Celery Flower monitoring UI
	cd workers && celery -A celery_app flower --port=5555

# Oracle DB (local)
oracle-sqlplus: ## Connect to Oracle DB via SQLPlus
	docker exec -it golib-oracle sqlplus system/oracle@localhost:1521/XE

# Redis
redis-cli: ## Connect to Redis
	docker exec -it golib-redis redis-cli

# Development helpers
db-seed: ## Seed database with test data
	@echo "$(YELLOW)Seeding database...$(NC)"
	$(GO) run scripts/seed.go

githooks: ## Install git hooks
	cp scripts/pre-commit .git/hooks/
	chmod +x .git/hooks/pre-commit
