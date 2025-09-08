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

# Docker build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t tutuplapak-api .

# Docker run
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 --env-file .env tutuplapak-api
