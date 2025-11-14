# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go backend template built with the Gin web framework. The project uses Go 1.25.4 and follows a clean layer-based architecture.

Module: `github.com/SuperIntelligence-Labs/go-backend-template`

## Project Structure

```
.
├── cmd/api/              # Application entry point (main.go - package main)
├── internal/             # Private application code
│   ├── handler/          # HTTP request handlers (controllers)
│   ├── service/          # Business logic layer
│   ├── repository/       # Data access layer
│   ├── model/            # Domain models / DB entities
│   ├── dto/              # Data Transfer Objects
│   │   ├── request/      # API request structures
│   │   └── response/     # API response structures
│   ├── middleware/       # Gin middleware (auth, logging, etc.)
│   ├── response/         # Standard response helpers (IMPLEMENTED)
│   ├── router/           # Route definitions
│   ├── database/         # Database connection setup
│   └── validator/        # Custom validation functions
├── config/               # Configuration files (yaml, env)
├── migrations/           # Database migrations (SQL files)
├── scripts/              # Build and deployment scripts
├── pkg/                  # Reusable packages (project-agnostic)
│   ├── httpclient/       # HTTP client utilities
│   └── utils/            # General utilities
└── tmp/                  # Build artifacts (gitignored)
```

## Architecture

This project follows a **layer-based architecture** where each layer has specific responsibilities:

### Layer Flow
```
HTTP Request → Router → Middleware → Handler → Service → Repository → Database
                                         ↓
HTTP Response ← Response Helper ← Handler ← Service ← Repository
```

### Key Layers

1. **Handler** - Validates requests, calls services, returns responses
2. **Service** - Contains business logic, coordinates repositories
3. **Repository** - Database operations only (CRUD)
4. **Model** - Database entities with GORM tags
5. **DTO** - API request/response structures (decoupled from models)

### Response Package (`internal/response/`)

The response package provides standardized JSON response formats for all API endpoints:

**Success Responses** (`success_response.go`):
- Structure: `Response` struct with `success`, `response`, `message`, `timestamp` fields
- Helper functions: `OK()`, `Created()`, `Updated()`, `Deleted()`
- Each helper automatically sets appropriate HTTP status codes and messages

**Error Responses** (`error_response.go`):
- Structure: `ErrorResponse` struct with `success`, `message`, `error`, `timestamp` fields
- Both top-level `message` (user-friendly) and `error.message` (detailed) are set
- Error codes: `ErrCodeValidation`, `ErrCodeNotFound`, `ErrCodeUnauthorized`, `ErrCodeForbidden`, `ErrCodeConflict`, `ErrCodeBadRequest`, `ErrCodeInternal`, `ErrCodeRateLimit`
- Helper functions: `BadRequest()`, `ValidationError()`, `NotFound()`, `Unauthorized()`, `Forbidden()`, `Conflict()`, `InternalError()`, `RateLimitExceeded()`
- `ValidationError()` automatically parses `validator.ValidationErrors` and formats them into user-friendly messages
- Custom validation tags supported: `indian_phone`, `court_id` (see `formatValidationError()` in error_response.go:110-130)

**Response Format Examples:**

Success:
```json
{
  "success": true,
  "message": "Request successful",
  "response": { ... },
  "timestamp": 1234567890
}
```

Error:
```json
{
  "success": false,
  "message": "Unable to process your request",
  "error": {
    "code": "BAD_REQUEST",
    "message": "Detailed error",
    "details": { ... }
  },
  "timestamp": 1234567890
}
```

**Important Notes:**
- All error response helpers call `c.AbortWithStatusJSON()` to prevent further handler execution
- Always use the response helpers instead of raw `c.JSON()` for consistency
- Top-level `message` field is now mandatory in both success and error responses

## Development Commands

### Build
```bash
# Build from cmd/api
go build -o ./tmp/main ./cmd/api

# Build with optimizations
go build -ldflags="-s -w" -o ./tmp/main ./cmd/api
```

### Run with Hot Reload
```bash
air
```

Air configuration (`.air.toml`):
- Watches all `.go` files except `_test.go`
- Builds from `./cmd/api` (not root)
- Output: `./tmp/main`
- Logs: `./tmp/build-errors.log`
- Post-build: Automatically runs `chmod +x ./tmp/main` to ensure binary is executable

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests in a specific package
go test ./internal/handler

# Run a specific test
go test -v -run TestName ./internal/handler

# Run with coverage
go test -cover ./...
```

### Dependencies
```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Verify dependencies
go mod verify
```

## Key Dependencies

- **Gin** (`github.com/gin-gonic/gin`): Web framework for routing and HTTP handling
- **Validator** (`github.com/go-playground/validator/v10`): Struct validation with custom tags
- **Sonic** (`github.com/bytedance/sonic`): High-performance JSON serialization (used by Gin)
- **UUID** (`github.com/google/uuid`): UUID generation
- **Viper** (optional): Configuration management (structure ready, not yet implemented)
- **GORM** (optional): ORM for database operations (structure ready)

## Adding a New Feature

Follow this order when adding a new feature (e.g., "Posts"):

1. **Model** - `internal/model/post.go` - Database entity with GORM tags
2. **DTOs** - `internal/dto/request/post_request.go` and `internal/dto/response/post_response.go`
3. **Repository** - `internal/repository/post_repository.go` - Database operations
4. **Service** - `internal/service/post_service.go` - Business logic
5. **Handler** - `internal/handler/post_handler.go` - HTTP endpoints
6. **Routes** - `internal/router/router.go` - Register routes
7. **Migration** - `migrations/000XXX_create_posts_table.up.sql` - Database schema

## Development Notes

- **Main package**: `cmd/api/main.go` must use `package main` (not `package api`)
- **Air configuration**: Automatically rebuilds when Go files change, excluding tests
- **Empty directories**: Use `.gitkeep` files to track empty directory structure in git
- **Response consistency**: Always use `internal/response` helpers for all API responses
- **Error handling**: All error responses abort further execution automatically
- **JSON performance**: Project uses Sonic for faster JSON encoding/decoding
- **Internal package**: Code in `internal/` cannot be imported by external projects (Go convention)
- **Pkg vs Internal**: Only put truly reusable, project-agnostic code in `pkg/`

## Configuration Management

The project is structured to support Viper configuration:
- Config files go in root or `config/` directory
- Config loading code goes in `internal/config/config.go`
- Use `.env` for secrets (gitignored), `config.yaml` for non-sensitive settings
- See README.md "Configuration" section for detailed Viper setup guide

## Common Patterns

### Handler Pattern
```go
func (h *Handler) Method(c *gin.Context) {
    var req dto.Request
    if err := c.ShouldBindJSON(&req); err != nil {
        response.ValidationError(c, err)
        return
    }

    result, err := h.service.Method(req)
    if err != nil {
        response.InternalError(c)
        return
    }

    response.OK(c, result)
}
```

### Service Pattern
```go
func (s *Service) Method(req dto.Request) (*dto.Response, error) {
    // Business logic here
    model := toModel(req)

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

func (r *Repository) FindByID(id string) (*Model, error) {
    var model Model
    err := r.db.Where("id = ?", id).First(&model).Error
    return &model, err
}
```

## Testing Guidelines

- Place test files alongside source files: `handler_test.go` next to `handler.go`
- Use table-driven tests for multiple test cases
- Mock dependencies (repositories, external services)
- Test helpers should return errors, not call `t.Fatal()` directly
- Aim for high coverage in service layer (business logic)

## Additional Resources

- See `README.md` for comprehensive architecture guide
- See `README.md` for detailed explanation of each directory's purpose
- See `README.md` for step-by-step feature addition guide
