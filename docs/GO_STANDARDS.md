# Go Coding Standards and Best Practices

## 1. Naming Conventions

### Package Names
- Use short, concise, lowercase names
- No underscores or mixedCaps
- Should be singular (e.g., `time`, not `times`)
- Avoid generic names like `util`, `common`, `misc`

### Variable Names
- Use MixedCaps or mixedCaps, not underscores
- Keep acronyms uppercase (e.g., `userID`, `HTTP`, `URL`)
- Short variable names for small scopes (e.g., `i` for loops)
- Descriptive names for larger scopes
- Avoid Hungarian notation

### Interface Names
- One-method interfaces end in -er (e.g., `Reader`, `Writer`)
- Choose good names for multiple-method interfaces
- Avoid `I` prefix (not `IReader`)

### Function and Method Names
- Use MixedCaps
- Keep names concise but descriptive
- Getters don't use `Get` prefix (e.g., `user.Name()`, not `user.GetName()`)
- Setters may use `Set` prefix (e.g., `user.SetName()`)

## 2. Code Organization

### Package Organization
```go
package main

import (
    // Standard library imports
    "fmt"
    "strings"

    // Third-party imports (separated by blank line)
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    // Local imports (separated by blank line)
    "myapp/internal/config"
    "myapp/pkg/logger"
)
```

### Directory Structure
```
project/
├── cmd/                    # Command line applications
│   └── api/               # Main application entry points
├── internal/              # Private application code
│   ├── domain/           # Business domain types
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic
│   └── handler/          # Request handlers
├── pkg/                   # Public libraries
├── test/                  # Additional test code
└── docs/                  # Documentation
```

## 3. Code Style

### General
- Use `gofmt` or `goimports` for formatting
- Maximum line length of 100-120 characters
- Group related declarations
- Order types, constants, variables
- Declare variables close to their use

### Error Handling
```go
// Preferred
if err != nil {
    return fmt.Errorf("failed to process request: %w", err)
}

// Avoid
if err != nil {
    log.Printf("failed to process request: %v", err)
    return err
}
```

### Comments
```go
// Package user provides user management functionality
package user

// User represents an authenticated user in the system
type User struct {
    // ID is the unique identifier for the user
    ID string
    
    // Name is the user's full name
    Name string
}

// NewUser creates a new user with the given name
func NewUser(name string) *User {
    // ...
}
```

## 4. Best Practices

### Interface Design
- Keep interfaces small
- Accept interfaces, return structs
- Interface satisfaction should be obvious
```go
// Good
type Reader interface {
    Read(p []byte) (n int, err error)
}

// Avoid
type DoEverything interface {
    DoA()
    DoB()
    DoC()
}
```

### Concurrency
- Don't communicate by sharing memory; share memory by communicating
- Use channels for communication between goroutines
- Use mutexes for simple state protection
- Always clean up goroutines
```go
// Good
done := make(chan bool)
go func() {
    defer close(done)
    // Do work
}()
<-done
```

### Testing
- Table-driven tests
- Use meaningful test names
- One assertion per test
- Use subtests for better organization
```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "valid input",
            input:    "hello",
            expected: "HELLO",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Something(tt.input)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## 5. Documentation

### Package Documentation
- Every package should have a package comment
- Document exported names
- Use complete sentences
- Start with the package name
```go
// Package user provides user management functionality for the application.
// It handles user creation, authentication, and authorization.
package user
```

### Function Documentation
- Begin with the function name
- Describe behavior precisely
- Document panics and thread-safety
```go
// CreateUser creates a new user with the given name and email.
// It returns an error if the email is already in use or if the
// name is empty.
func CreateUser(name, email string) (*User, error)
```

## 6. Common Mistakes to Avoid

### Don't
- Use panic for normal error handling
- Ignore errors
- Use global variables
- Return naked returns
- Use init() functions unless absolutely necessary
- Mix synchronous and asynchronous error handling

### Do
- Use context for cancellation
- Handle errors explicitly
- Use meaningful variable names
- Keep functions and methods small
- Write tests for exported functions
- Use defer for cleanup

## 7. Performance

### Tips
- Avoid premature optimization
- Profile before optimizing
- Use benchmarks to verify improvements
- Consider memory allocations
- Use buffer pools for frequent allocations

### Memory Management
```go
// Good - reuse buffer
var bufPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

// Use
buf := bufPool.Get().(*bytes.Buffer)
defer bufPool.Put(buf)
buf.Reset()
```

## 8. Tools

### Essential Tools
- `go fmt` - Format code
- `go vet` - Find subtle problems
- `golint` - Style mistakes
- `errcheck` - Find unchecked errors
- `staticcheck` - Advanced static analysis
- `goimports` - Manage imports

### IDE Integration
- Configure automatic formatting
- Enable linting on save
- Set up go tools in the development environment

## 9. Project Structure Best Practices

### Configuration
- Use environment variables for deployment-specific values
- Keep configuration separate from code
- Use structured configuration files
- Support multiple environments

### Dependency Management
- Use go modules
- Keep dependencies up to date
- Vendor dependencies when needed
- Review licenses of dependencies

### Logging
- Use structured logging
- Include relevant context
- Log at appropriate levels
- Don't log sensitive information

## 10. Security Best Practices

### General Security
- Never store secrets in code
- Use secure random numbers for security purposes
- Keep dependencies updated
- Use prepared statements for SQL
- Validate all input

### Authentication & Authorization
- Use secure password hashing
- Implement rate limiting
- Use secure session management
- Implement proper access control
- Use HTTPS in production

## 11. Testing Standards

### Unit Tests
- Test package API, not internals
- Use table-driven tests
- Mock external dependencies
- Test error conditions
- Keep tests simple

### Integration Tests
- Test real dependencies
- Use test containers
- Clean up test data
- Don't depend on test order

### Benchmarks
- Write meaningful benchmarks
- Test realistic scenarios
- Use realistic data sizes
- Compare benchmarks fairly 