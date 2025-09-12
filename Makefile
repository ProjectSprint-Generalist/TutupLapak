# TutupLapak API Makefile

.PHONY: help build run test clean deps dev

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Build the application
build: ## Build the application
	@echo "Building TutupLapak API..."
	go build -o bin/tutuplapak main.go

# Run the application
run: ## Run the application
	@echo "Starting TutupLapak API..."
	go run main.go

# Run in development mode with hot reload (requires air)
dev: ## Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
	@echo "Starting TutupLapak API in development mode..."
	air

# Install dependencies
deps: ## Install dependencies
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Run tests
test: ## Run tests
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html

# Format code
fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint: ## Lint code
	@echo "Linting code..."
	golangci-lint run

# Install development tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Setup development environment
setup: deps install-tools ## Setup development environment
	@echo "Setting up development environment..."
	@if [ ! -f .env ]; then cp .env.example .env; fi
	@echo "Development environment setup complete!"

# Docker commands
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t tutuplapak-api .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env tutuplapak-api

# Docker Compose commands
docker-up: ## Start all services with docker-compose
	@echo "Starting all services..."
	docker-compose up -d

docker-up-dev: ## Start all services in development mode
	@echo "Starting all services in development mode..."
	docker-compose -f docker-compose.dev.yml up -d

docker-down: ## Stop all services
	@echo "Stopping all services..."
	docker-compose down

docker-down-dev: ## Stop all development services
	@echo "Stopping all development services..."
	docker-compose -f docker-compose.dev.yml down

docker-logs: ## Show logs for all services
	@echo "Showing logs for all services..."
	docker-compose logs -f

docker-logs-app: ## Show logs for app service only
	@echo "Showing logs for app service..."
	docker-compose logs -f tutuplapak-app

docker-restart: ## Restart all services
	@echo "Restarting all services..."
	docker-compose restart

docker-rebuild: ## Rebuild and restart all services
	@echo "Rebuilding and restarting all services..."
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d

docker-rebuild-dev: ## Rebuild and restart all development services
	@echo "Rebuilding and restarting all development services..."
	docker-compose -f docker-compose.dev.yml down
	docker-compose -f docker-compose.dev.yml build --no-cache
	docker-compose -f docker-compose.dev.yml up -d

# Database commands
db-migrate: ## Run database migrations
	@echo "Running database migrations..."
	docker-compose exec tutuplapak-app ./tutuplapak migrate

db-seed: ## Seed database with sample data
	@echo "Seeding database..."
	docker-compose exec tutuplapak-app ./tutuplapak seed

db-reset: ## Reset database (WARNING: This will delete all data)
	@echo "Resetting database..."
	docker-compose down -v
	docker-compose up -d postgres
	@echo "Waiting for database to be ready..."
	@sleep 10
	docker-compose up -d

# Utility commands
docker-clean: ## Clean up Docker resources
	@echo "Cleaning up Docker resources..."
	docker-compose down -v
	docker system prune -f
	docker volume prune -f

docker-status: ## Show status of all services
	@echo "Service status:"
	docker-compose ps

docker-shell: ## Open shell in app container
	@echo "Opening shell in app container..."
	docker-compose exec tutuplapak-app sh

docker-db-shell: ## Open database shell
	@echo "Opening database shell..."
	docker-compose exec postgres psql -U tutuplapak -d tutuplapak_db

# Production commands
docker-up-prod: ## Start production services
	@echo "Starting production services..."
	docker-compose -f docker-compose.prod.yml up -d

docker-down-prod: ## Stop production services
	@echo "Stopping production services..."
	docker-compose -f docker-compose.prod.yml down

docker-logs-prod: ## Show production logs
	@echo "Showing production logs..."
	docker-compose -f docker-compose.prod.yml logs -f
