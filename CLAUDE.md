# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go backend template built with the Gin web framework. The project uses Go 1.25.4 and follows a clean package-based architecture.

Module: `github.com/SuperIntelligence-Labs/go-backend-template`

## Project Structure

```
├── cmd/api/          # Application entry point (main.go)
├── pkg/
│   ├── response/     # Standardized HTTP response handlers
│   ├── httpclient/   # HTTP client utilities (empty)
│   └── utils/        # General utilities (empty)
└── tmp/              # Build artifacts (excluded from version control)
```

## Architecture

### Response Package (`pkg/response/`)

The response package provides standardized JSON response formats for all API endpoints:

**Success Responses** (`success_response.go`):
- All success responses follow the `Response` struct with `success`, `response`, `message`, and `timestamp` fields
- Helper functions: `OK()`, `Created()`, `Updated()`, `Deleted()`
- Each helper automatically sets appropriate HTTP status codes and messages

**Error Responses** (`error_response.go`):
- All error responses follow the `ErrorResponse` struct with `success`, `error`, and `timestamp` fields
- Error codes are constants: `ErrCodeValidation`, `ErrCodeNotFound`, `ErrCodeUnauthorized`, `ErrCodeForbidden`, `ErrCodeConflict`, `ErrCodeBadRequest`, `ErrCodeInternal`, `ErrCodeRateLimit`
- Helper functions: `BadRequest()`, `ValidationError()`, `NotFound()`, `Unauthorized()`, `Forbidden()`, `Conflict()`, `InternalError()`, `RateLimitExceeded()`
- `ValidationError()` automatically parses `validator.ValidationErrors` and formats them into user-friendly messages
- Custom validation tags are supported: `indian_phone`, `court_id` (see `formatValidationError()` for the full list)

**Important**: All response helpers call `c.AbortWithStatusJSON()` for errors to prevent further handler execution.

## Development Commands

### Build
```bash
go build -o ./tmp/main .
```

### Run with Hot Reload
```bash
air
```
Air watches for file changes and automatically rebuilds. Configuration is in `.air.toml`. Build output goes to `./tmp/main` and logs to `build-errors.log`.

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests in a specific package
go test ./pkg/response

# Run a specific test
go test -v -run TestName ./pkg/response
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
- **Sonic** (`github.com/bytedance/sonic`): High-performance JSON serialization
- **UUID** (`github.com/google/uuid`): UUID generation

## Development Notes

- The build process excludes `_test.go` files (configured in `.air.toml`)
- Air excludes `tmp/`, `vendor/`, `testdata/`, and `assets/` directories from watching
- The project uses Sonic for JSON encoding/decoding (faster than standard library)
