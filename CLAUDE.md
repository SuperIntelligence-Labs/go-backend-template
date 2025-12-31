# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go backend template built with the Echo web framework. The project uses Go 1.25.4 and follows a clean layer-based, feature-based architecture.

Module: `github.com/SuperIntelligence-Labs/go-backend-template`

## Project Structure

```
.
├── cmd/api/              # Application entry point (main.go - package main)
├── internal/             # Private application code
│   ├── config/           # Configuration management (Viper)
│   ├── database/         # Database connection setup (PostgreSQL/GORM)
│   ├── errors/           # Common error definitions
│   ├── features/         # Feature modules (each with handler/service/repository/model)
│   │   └── example/      # Example feature demonstrating the pattern
│   ├── logger/           # Zerolog logger setup
│   ├── middleware/       # Echo middleware (JWT, logging)
│   ├── response/         # Standard response helpers
│   └── server/           # Server and router setup
├── migrations/           # Database migrations (SQL files)
├── scripts/              # Build and deployment scripts
├── pkg/                  # Reusable packages (project-agnostic)
└── tmp/                  # Build artifacts (gitignored)
```

## Architecture

This project follows a **feature-based, layer-based architecture**:

### Layer Flow
```
HTTP Request → Router → Middleware → Handler → Service → Repository → Database
                                         ↓
HTTP Response ← Response Helper ← Handler ← Service ← Repository
```

### Feature Structure

Each feature in `internal/features/` contains:
- **model.go** - Database entity with GORM tags
- **repository.go** - Database operations (CRUD)
- **service.go** - Business logic + request/response DTOs
- **handler.go** - HTTP request handlers
- **routes.go** - Route registration

See `internal/features/example/` for a complete reference implementation.

## Development Commands

### Build
```bash
go build -o ./tmp/main ./cmd/api
```

### Run with Hot Reload
```bash
air
```

### Testing
```bash
go test ./...
go test -v ./internal/features/example
```

### Dependencies
```bash
go mod download
go mod tidy
```

## Key Dependencies

- **Echo** (`github.com/labstack/echo/v4`): HTTP web framework
- **GORM** (`gorm.io/gorm`): ORM for database operations
- **Validator** (`github.com/go-playground/validator/v10`): Struct validation
- **Zerolog** (`github.com/rs/zerolog`): Structured logging
- **Viper** (`github.com/spf13/viper`): Configuration management
- **JWT** (`github.com/golang-jwt/jwt/v5`): JWT authentication

## Adding a New Feature

1. Create directory: `internal/features/yourfeature/`
2. Create files following the example pattern:
   - `model.go` - Database entity
   - `repository.go` - CRUD operations
   - `service.go` - Business logic + DTOs
   - `handler.go` - HTTP handlers
   - `routes.go` - Route registration
3. Wire up in `cmd/api/main.go` (DI setup)
4. Register routes in `internal/server/router.go`

## Configuration

Copy `.env.example` to `.env` and configure:
- Server settings (host, port, env)
- Database credentials
- JWT secrets
- Logging level

## Common Patterns

### Handler Pattern
```go
func (h *Handler) Create(c echo.Context) error {
    var req CreateRequest
    if err := c.Bind(&req); err != nil {
        return response.ErrBadRequest("Invalid request", nil)
    }
    if err := c.Validate(&req); err != nil {
        return err
    }
    result, err := h.service.Create(req)
    if err != nil {
        return response.ErrInternalError(err)
    }
    return response.Created(c, "Created successfully", result)
}
```

### Service Pattern
```go
func (s *Service) Create(req CreateRequest) (*Response, error) {
    model := &Model{Name: req.Name}
    if err := s.repo.Create(model); err != nil {
        return nil, err
    }
    return toResponse(model), nil
}
```

### Repository Pattern
```go
func (r *Repository) Create(model *Model) error {
    return r.db.Create(model).Error
}
```
