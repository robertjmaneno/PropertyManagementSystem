# Development Guide

This guide provides detailed information for developers working with the template.

## Development Environment Setup

### Required Tools

1. **Go 1.21+**
   ```bash
   # macOS
   brew install go
   
   # Linux
   wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
   ```

2. **PostgreSQL 13+**
   ```bash
   # macOS
   brew install postgresql@13
   
   # Linux
   sudo apt-get install postgresql-13
   ```

3. **Docker & Docker Compose**
   ```bash
   # macOS
   brew install docker docker-compose
   
   # Linux
   sudo apt-get install docker.io docker-compose
   ```

4. **Development Tools**
   ```bash
   # Install required Go tools
   make setup-dev
   ```

### IDE Setup

#### VSCode
1. Install Go extension
2. Recommended settings:
   ```json
   {
     "go.useLanguageServer": true,
     "go.lintTool": "golangci-lint",
     "go.formatTool": "gofmt",
     "editor.formatOnSave": true
   }
   ```

#### GoLand
1. Enable Go modules integration
2. Configure golangci-lint
3. Enable format on save

## Development Workflow

### 1. Branch Management

```bash
# Create feature branch
git checkout -b feature/your-feature

# Create bugfix branch
git checkout -b fix/bug-description

# Create release branch
git checkout -b release/v1.0.0
```

### 2. Code Generation

```bash
# Generate mocks
make mock

# Generate Swagger docs
make swagger

# Generate all
make generate
```

### 3. Testing

```bash
# Run tests during development
make test-watch

# Run specific test
go test ./internal/service -run TestUserService_Create

# Run tests with race detection
make test-race

# Generate test coverage report
make test-coverage
```

### 4. Database Management

```bash
# Create new migration
make migrate-create name=add_user_roles

# Apply migrations
make migrate-up

# Rollback last migration
make migrate-down

# Reset database
make migrate-reset
```

### 5. Local Development

```bash
# Start development server with hot reload
make run-dev

# Start with specific config
CONFIG_FILE=config/local.yaml make run

# Start with debugger
make run-debug
```

## Code Organization

### Package Structure

1. **Domain Layer** (`internal/domain`)
   - Business entities
   - Repository interfaces
   - Service interfaces
   - Domain errors

2. **Repository Layer** (`internal/repository`)
   - Database implementations
   - Caching logic
   - External service clients

3. **Service Layer** (`internal/service`)
   - Business logic
   - Transaction management
   - Service composition

4. **Handler Layer** (`internal/handler`)
   - HTTP handlers
   - Request/response mapping
   - Input validation

### Dependency Management

1. **Adding Dependencies**
   ```bash
   go get -u github.com/example/package
   ```

2. **Updating Dependencies**
   ```bash
   go get -u ./...
   go mod tidy
   ```

3. **Vendoring (if needed)**
   ```bash
   go mod vendor
   ```

## Testing Strategy

### 1. Unit Tests
- Test business logic in isolation
- Use mocks for dependencies
- Focus on edge cases
- Example:
  ```go
  func TestUserService_Create(t *testing.T) {
      mockRepo := new(MockUserRepository)
      service := NewUserService(mockRepo)
      // Test cases...
  }
  ```

### 2. Integration Tests
- Test component interactions
- Use test containers
- Focus on happy path
- Example:
  ```go
  func TestUserRepository_Create(t *testing.T) {
      db := testutil.NewTestDB(t)
      repo := NewUserRepository(db)
      // Test cases...
  }
  ```

### 3. E2E Tests
- Test complete flows
- Use real dependencies
- Focus on user scenarios
- Example:
  ```go
  func TestUserFlow(t *testing.T) {
      app := SetupTestApp(t)
      client := NewTestClient(app.Port)
      // Test scenarios...
  }
  ```

## Error Handling

### 1. Domain Errors
```go
var (
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidInput = errors.New("invalid input")
)
```

### 2. HTTP Errors
```go
func handleError(c *gin.Context, err error) {
    switch {
    case errors.Is(err, domain.ErrUserNotFound):
        c.JSON(http.StatusNotFound, ErrorResponse{...})
    default:
        c.JSON(http.StatusInternalServerError, ErrorResponse{...})
    }
}
```

## Logging

### 1. Configuration
```go
logger := zap.NewProduction()
defer logger.Sync()
```

### 2. Usage
```go
logger.Info("user created",
    zap.String("user_id", user.ID),
    zap.String("email", user.Email),
)
```

## Monitoring

### 1. Metrics
```go
// Register metrics
userCreationCounter := prometheus.NewCounter(...)

// Use in code
userCreationCounter.Inc()
```

### 2. Tracing
```go
span, ctx := tracer.Start(ctx, "CreateUser")
defer span.End()
```

## Security

### 1. Input Validation
```go
type CreateUserRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Username string `json:"username" binding:"required,min=3"`
    Password string `json:"password" binding:"required,min=8"`
}
```

### 2. Authentication
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // JWT validation logic
    }
}
```

## Performance

### 1. Database
- Use indexes appropriately
- Implement caching where needed
- Use database connection pooling

### 2. API
- Implement rate limiting
- Use pagination for list endpoints
- Cache responses where appropriate

## Deployment

### 1. Building
```bash
# Build for production
make build-prod

# Build for specific platform
GOOS=linux GOARCH=amd64 make build
```

### 2. Docker
```bash
# Build image
docker build -t myapp:latest .

# Run container
docker run -p 8080:8080 myapp:latest
```

## Troubleshooting

### Common Issues

1. **Database Connection**
   ```bash
   # Check database status
   make db-status
   
   # Reset database
   make db-reset
   ```

2. **Application Logs**
   ```bash
   # View logs
   make logs
   
   # View specific service logs
   make logs-api
   ```

3. **Performance Issues**
   ```bash
   # Profile application
   make profile
   
   # View metrics
   make metrics
   