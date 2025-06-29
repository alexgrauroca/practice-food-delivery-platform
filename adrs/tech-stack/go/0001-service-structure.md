# 1. Service Structure

## Status

Accepted

Date: 2025-06-24

## Context

As we build domain services using Go, we need a consistent project structure that:

- Supports the separation of concerns
- Enables efficient testing at various levels
- Follows Go best practices and conventions
- Provides a clear organization for new team members to understand
- Facilitates dependency management and code reuse
- Enforces consistent patterns across services

Without a standardized structure, we risk inconsistent implementations, increased maintenance overhead, and slower 
onboarding for new team members.

## Decision

We will adopt a standardized project structure for Go domain services following these principles:

### Top-Level Organization

```
/services/{service-name}/
├── cmd/                      # Command applications
├── docs/                     # Documentation
├── internal/                 # Private application code
├── Makefile                  # Build automation
├── Dockerfile                # Container definition
├── go.mod                    # Dependency management
├── go.sum                    # Dependency checksums
└── README.md                 # Service documentation
```

### Internal Package Organization

The `internal` directory contains all service-specific code that should not be imported by other services:

```
/internal/
├── {domain-features}/        # Business feature packages (e.g., customers)
│   ├── handler.go            # HTTP handlers
│   ├── service.go            # Business logic
│   ├── repository.go         # Data access
│   ├── errors.go             # Domain-specific errors
│   └── *_test.go             # Tests for each component
├── config/                   # Configuration management
├── clock/                    # Time abstractions
├── middleware/               # HTTP middleware
├── logctx/                   # Logging context utilities
└── infrastructure/           # External dependencies implementation
```

### Package Design Principles

1. **Layered Architecture**: Each domain feature follows a layered architecture:
   - Handler layer (HTTP/API interface)
   - Service layer (business logic)
   - Repository layer (data access)

2. **Clear Separation of Concerns**:
   - Handlers handle HTTP concerns only
   - Services contain business logic
   - Repositories handle data persistence

3. **Shared Utilities**: Cross-cutting concerns (logging, time, etc.) are in dedicated packages

## Consequences

### Positive

- Consistent structure improves developer experience
- Clear separation of concerns reduces complexity
- Layered architecture improves testability
- Reduced cognitive load when switching between services

### Negative

- More initial structure may seem complex for very simple services
- Requires discipline to maintain the structure

### Neutral

- May need to evolve as the organization grows
- Different services may need slight variations based on specific requirements

## Implementation Notes

### File Naming Conventions

- `handler.go`: HTTP handlers and routing logic
- `service.go`: Business logic implementation
- `repository.go`: Data access logic
- `errors.go`: Domain-specific error types
- `models.go`: Data structures for the domain

### Entry Point (cmd/main.go)

The main application entry point coordinates all components:

```go
func main() {
    // Initialize logger, database, etc.

    // Initialize and wire components together
    repo := feature.NewRepository(logger, db, ...)
    service := feature.NewService(logger, repo, ...)
    handler := feature.NewHandler(logger, service)

    // Register routes
    router := gin.Default()
    handler.RegisterRoutes(router)

    // Start the server
    router.Run(":8080")
}
```

## Related Documents

- [The Go Project Layout](https://github.com/golang-standards/project-layout)

## Contributors

- Àlex Grau Roca
