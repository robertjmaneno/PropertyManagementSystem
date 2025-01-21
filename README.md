# Go Clean Architecture Template

A production-ready template for building scalable web applications in Go, following clean architecture principles.

## 🌟 Features

### Core Features
- Clean Architecture design
- Dependency Injection
- Configuration management with Viper
- Structured logging with Zap
- OpenAPI/Swagger documentation
- Database-agnostic migrations (SQLite, PostgreSQL, MySQL)
- User management with search functionality
- Comprehensive testing suite

### Security Features
- Third-party authentication integration
- Role-based access control (RBAC)
- CSRF protection (configurable per route)
- Secure session management
- Comprehensive error handling
- Service-to-service authentication
- Security headers and CORS configuration

## 🚀 Getting Started

### Prerequisites
- Go 1.21 or higher
- Make
- (Optional) Docker for containerization

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/projectname.git
cd projectname
```

2. Install dependencies and tools:
```bash
make setup  # Installs everything (dependencies + tools)
# Or individually:
make deps        # Install project dependencies only
make deps-tools  # Install development tools only
```

3. Copy example configuration:
```bash
cp config/config.example.yaml config/config.yaml
```

4. Configure your environment:
   - Update database settings in config.yaml
   - Set authentication service URL and credentials
   - Configure CORS settings if needed

5. Run the application:
```bash
make dev  # Runs swagger + build + migrate + run
```

## 📡 API Endpoints

### User Management
- `GET /api/v1/users` - List all users
- `GET /api/v1/users/search` - Search users by name or email
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users/{id}` - Get user by ID
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user

### Pagination and Sorting
The list and search endpoints support pagination and sorting:

#### Pagination Parameters
- `page` - Page number (default: 1)
- `page_size` - Number of items per page (default: 10, max: 100)

#### Sorting Parameters
- `sort` - Comma-separated list of fields to sort by
- Format: `field1:asc,field2:desc`
- Available fields: `first_name`, `last_name`, `email`, `created_at`
- Example: `/api/v1/users?sort=first_name:asc,created_at:desc`

Examples:
```bash
# Get the first page with 10 items
GET /api/v1/users?page=1&page_size=10

# Get the second page with 20 items, sorted by email
GET /api/v1/users?page=2&page_size=20&sort=email:asc

# Sort by multiple fields
GET /api/v1/users?sort=last_name:asc,first_name:asc
```

## 🗄️ Database Management

The project uses GORM for database operations and supports multiple database types:

### Supported Databases
- PostgreSQL (default)
- MySQL
- SQLite

### Running Migrations

```bash
# Using default PostgreSQL database
make migrate

# Using different PostgreSQL credentials
DB_DSN="postgres://different_user:pass@localhost:5432/dbname?sslmode=disable" make migrate

# Using MySQL
DB_DSN="mysql://user:pass@tcp(localhost:3306)/dbname" make migrate

# Using SQLite
DB_DSN="sqlite3://data/app.db" make migrate
```

### Creating New Migrations

1. Add your domain model in `internal/domain/`
2. Add migration functions in `internal/migrations/migrations.go`:
```go
// Add to migrations slice
migrations := []struct {
    name string
    fn   func(*gorm.DB) error
}{
    {"existing_migration", existingMigration},
    {"your_new_migration", yourNewMigration}, // Add here
}

// Add to rollbacks map
rollbacks := map[string]func(*gorm.DB) error{
    "existing_migration": existingRollback,
    "your_new_migration": yourNewRollback, // Add here
}

// Add migration functions
func yourNewMigration(db *gorm.DB) error {
    return db.AutoMigrate(&YourModel{})
}

func yourNewRollback(db *gorm.DB) error {
    return db.Migrator().DropTable(&YourModel{})
}
```

## 🏗️ Project Structure

```
.
├── cmd/                    # Application entry points
│   ├── api/               # API server
│   └── migrate/           # Migration tool
├── config/                # Configuration files
├── internal/              # Private application code
│   ├── client/           # External service clients
│   ├── config/           # Configuration structures
│   ├── domain/           # Domain models and interfaces
│   ├── dto/              # Data Transfer Objects
│   ├── middleware/       # HTTP middleware
│   ├── migrations/       # Database migrations
│   ├── router/           # HTTP routing and handlers
│   ├── repository/       # Data access layer
│   │   ├── base_repository.go     # Base repository utilities
│   │   ├── generic_repository.go  # Generic CRUD operations
│   │   ├── interfaces.go          # Repository interfaces
│   │   └── user_repository.go     # User-specific repository
│   └── service/          # Business logic layer
├── pkg/                  # Public libraries
│   ├── logger/          # Logging utilities
│   └── errors/          # Error handling
└── test/                # Test suites
    ├── unit/           # Unit tests
    ├── integration/    # Integration tests
    └── testutil/       # Test utilities
```

## 📦 Repository Pattern

This project implements a robust repository pattern with generics for data access. The architecture is designed to be:
- Type-safe with Go generics
- Easy to extend
- Consistent across different entities
- Testable with clear separation of concerns

### Repository Structure

1. **Base Repository** (`base_repository.go`)
   - Provides common utilities and helper methods
   - Handles database connection management
   - Implements common query modifiers (pagination, sorting, includes)

2. **Generic Repository** (`generic_repository.go`)
   - Implements common CRUD operations
   - Type-safe with Go generics
   - Can be embedded in specific repositories

3. **Entity Repositories** (e.g., `user_repository.go`)
   - Extend generic repository
   - Implement entity-specific operations
   - Override generic methods if needed

### Adding a New Repository

1. Define your entity in `internal/domain`:
```go
type YourEntity struct {
    ID        string    `json:"id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    // Add your fields
}
```

2. Add repository interface in `internal/repository/interfaces.go`:
```go
type YourEntityRepository interface {
    Repository[domain.YourEntity]  // Embed generic repository
    // Add custom methods
    SearchByField(ctx context.Context, field string) ([]*domain.YourEntity, error)
}
```

3. Implement the repository in `internal/repository/your_entity_repository.go`:
```go
type yourEntityRepository struct {
    *GenericRepository[domain.YourEntity]
}

func NewYourEntityRepository(db *gorm.DB) YourEntityRepository {
    return &yourEntityRepository{
        GenericRepository: NewGenericRepository(db, domain.YourEntity{}),
    }
}

// Implement custom methods
func (r *yourEntityRepository) SearchByField(ctx context.Context, field string) ([]*domain.YourEntity, error) {
    // Implementation
}
```

### Using Repositories

1. **Basic CRUD Operations**
```go
// Create
entity := &domain.YourEntity{...}
err := repo.Create(ctx, entity)

// Read
entity, err := repo.GetByID(ctx, "id")

// Update
entity.Field = "new value"
err := repo.Update(ctx, entity)

// Delete
err := repo.Delete(ctx, "id")
```

2. **Pagination and Sorting**
```go
opts := repository.ListOptions{
    Pagination: &pagination.Params{
        Page:     1,
        PageSize: 10,
    },
    Sort: sorting.NewOptions("created_at", sorting.DESC),
}

// Get paginated results
results, err := repo.GetPaged(ctx, opts)

// Get total count
total, err := repo.Count(ctx, opts)
```

## 🛠️ Development

### Available Make Commands

```bash
# Development
make dev           # Full development setup
make run           # Run the application
make build         # Build the application
make clean         # Clean build artifacts

# Dependencies
make setup         # Install all dependencies and tools
make deps          # Install project dependencies
make deps-tools    # Install development tools

# Database
make migrate       # Run all migrations
make migrate-up    # Same as migrate
make migrate-down  # Rollback last migration

# Testing
make test              # Run all tests
make test-unit         # Run unit tests
make test-integration  # Run integration tests
make test-coverage     # Run tests with coverage

# Tools
make swagger      # Generate Swagger docs
make lint         # Run linter
make mock         # Generate mocks

# CI/CD
make ci           # Run CI checks (lint + tests)
```

## 🔐 Security Features

1. **Authentication**:
   - Third-party authentication integration
   - Token validation and verification
   - Role-based access control

2. **CSRF Protection**:
   - Configurable per route group
   - Web routes protected by default
   - API routes configurable

3. **Error Handling**:
   - Secure error messages
   - Detailed internal logging
   - Client-safe error responses

4. **Configuration**:
   - Environment-based configuration
   - Secure defaults
   - Flexible CORS settings

## 📚 API Documentation

Swagger documentation is automatically generated and available at `/swagger/index.html` when running the application.

To regenerate the documentation:
```bash
make swagger
```

### Adding Swagger Annotations
Swagger annotations should be added in the router file (`internal/router/router.go`), not in the handlers. This is because the router is where the routes are defined and where Swagger looks for documentation.

Example of documenting an endpoint:
```go
// ListUsers handles listing all users
// @Summary List users
// @Description Get a list of all users with pagination and sorting
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number (default: 1)" minimum(1)
// @Param page_size query int false "Number of items per page (default: 10)" minimum(1) maximum(100)
// @Param sort query string false "Sort fields (format: field1:asc,field2:desc)"
// @Success 200 {object} dto.ListUsersResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /users [get]
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details. 