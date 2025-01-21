.PHONY: build test test-unit test-integration test-short run clean swagger lint mock migrate migrate-up migrate-down deps deps-tools setup dev ci docs

# Build settings
BINARY_NAME=api
BUILD_DIR=bin
MAIN_FILE=cmd/api/main.go

# Test settings
TEST_FLAGS=-v -race
COVERAGE_FILE=coverage.out

# Environment settings
CONFIG_FILE?=config/config.yaml
export CONFIG_FILE

# Database settings
DB_USER?=postgres
DB_PASS?=postgres
DB_NAME?=template
DB_HOST?=localhost
DB_PORT?=5432
DB_SCHEMA?=app
DB_DSN?=postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable&search_path=$(DB_SCHEMA)
export DB_DSN

# Docker settings
DOCKER_IMAGE=yourusername/projectname
DOCKER_TAG?=latest

# Install project dependencies
deps:
	@echo "Installing project dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies installed successfully."

# Install development tools
deps-tools:
	@echo "Installing development tools..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/golang/mock/mockgen@latest
	@echo "Development tools installed successfully."
	
# Install everything
setup: deps deps-tools
	@echo "Setup completed successfully."

# Build the application
build:
	@echo "Building application..."
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build completed successfully."

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f $(COVERAGE_FILE)
	@echo "Clean completed successfully."

# Run the application
run:
	@echo "Starting application..."
	go run $(MAIN_FILE)

# Testing targets
test: test-unit test-integration

test-unit:
	@echo "Running unit tests..."
	go test $(TEST_FLAGS) ./test/unit/...

test-integration:
	@echo "Running integration tests..."
	go test $(TEST_FLAGS) ./test/integration/...

test-short:
	@echo "Running short tests..."
	go test $(TEST_FLAGS) -short ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test $(TEST_FLAGS) -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -html=$(COVERAGE_FILE)
	@echo "Coverage report generated at $(COVERAGE_FILE)"

# Development tools
lint:
	@echo "Running linter..."
	golangci-lint run
	@echo "Linting completed successfully."

swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/api/main.go -o docs
	@echo "Swagger documentation generated successfully."

# Generate mocks
mock: mock-repository mock-service

mock-repository:
	@echo "Generating repository mocks..."
	mockgen -source=internal/repository/interfaces.go -destination=test/mocks/repository/mock_repository.go
	mockgen -source=internal/repository/base_repository.go -destination=test/mocks/repository/mock_base_repository.go
	@echo "Repository mocks generated successfully."

mock-service:
	@echo "Generating service mocks..."
	mockgen -source=internal/service/user_service.go -destination=test/mocks/service/mock_user_service.go
	@echo "Service mocks generated successfully."

# Database migrations
migrate-tool:
	@echo "Building migration tool..."
	go build -o $(BUILD_DIR)/migrate cmd/migrate/main.go

migrate-up: migrate-tool
	@echo "Running migrations..."
	./$(BUILD_DIR)/migrate up

migrate-down: migrate-tool
	@echo "Rolling back migration..."
	./$(BUILD_DIR)/migrate down

migrate: migrate-up

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "Docker image built successfully."

docker-push:
	@echo "Pushing Docker image..."
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo "Docker image pushed successfully."

# Documentation
docs: swagger
	@echo "Documentation updated successfully."

# Development workflow targets
dev: swagger build migrate run

# CI targets
ci: lint test-coverage

# Help target
help:
	@echo "Available targets:"
	@echo "  setup          - Install all dependencies and tools"
	@echo "  deps           - Install project dependencies"
	@echo "  deps-tools     - Install development tools"
	@echo "  build          - Build the application"
	@echo "  clean          - Clean build artifacts"
	@echo "  run            - Run the application"
	@echo "  test           - Run all tests"
	@echo "  test-unit      - Run unit tests"
	@echo "  test-integration - Run integration tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  lint           - Run linter"
	@echo "  swagger        - Generate Swagger documentation"
	@echo "  mock           - Generate all mocks"
	@echo "  migrate        - Run database migrations"
	@echo "  migrate-down   - Rollback last migration"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-push    - Push Docker image"
	@echo "  dev            - Full development setup"
	@echo "  ci             - Run CI checks"
	@echo "  docs           - Update documentation"
	@echo "  help           - Show this help message"

# Database commands
db-create:
	@echo "Creating database and schema..."
	@PGPASSWORD=$(DB_PASS) createdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME) || true
	@PGPASSWORD=$(DB_PASS) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -c "CREATE SCHEMA IF NOT EXISTS $(DB_SCHEMA);" || true
	@echo "Database and schema created successfully."

db-drop:
	@echo "Dropping database..."
	@PGPASSWORD=$(DB_PASS) dropdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME) || true
	@echo "Database dropped successfully."

db-reset: db-drop db-create migrate