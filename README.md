# Go Backend Template

A production-ready Go backend template with clean architecture, built using the Echo web framework and following industry best practices.

## Features

- **Feature-based Architecture** - Clean separation with handler, service, repository layers per feature
- **Standardized API Responses** - Consistent success and error response formats
- **Hot Reload** - Development with Air for automatic rebuilds
- **JWT Authentication** - Ready-to-use JWT middleware
- **PostgreSQL + GORM** - Database integration with migrations support
- **Structured Logging** - Zerolog for production-ready logging
- **Configuration Management** - Viper with environment variable support
- **Request Validation** - go-playground/validator with custom error messages

## Tech Stack

- **Go 1.25.4**
- **Echo** - HTTP web framework
- **GORM** - ORM for database operations
- **Zerolog** - Structured logging
- **Viper** - Configuration management
- **Validator** - Request validation
- **Air** - Hot reload for development

## Project Structure

```
.
├── cmd/api/              # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── database/         # Database connection
│   ├── errors/           # Common error definitions
│   ├── features/         # Feature modules
│   │   └── example/      # Example CRUD feature
│   │       ├── model.go
│   │       ├── repository.go
│   │       ├── service.go
│   │       ├── handler.go
│   │       └── routes.go
│   ├── logger/           # Logging setup
│   ├── middleware/       # JWT, logging middleware
│   ├── response/         # Response helpers
│   └── server/           # Server and router
├── migrations/           # SQL migrations
├── scripts/              # Build scripts
└── pkg/                  # Reusable packages
```

## Getting Started

### Prerequisites

- Go 1.25.4 or higher
- PostgreSQL
- Air (optional, for hot reload): `go install github.com/air-verse/air@latest`

### Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd go-backend-template
```

2. Copy environment file:
```bash
cp .env.example .env
```

3. Configure `.env` with your settings

4. Install dependencies:
```bash
go mod download
```

5. Run the application:
```bash
# Development (with hot reload)
air

# Or directly
go run ./cmd/api
```

## Configuration

Environment variables (`.env`):

```env
SERVER_HOST=localhost
SERVER_PORT=8080
SERVER_ENV=development

LOG_LEVEL=debug

JWT_AT_SECRET=your-access-token-secret
JWT_AT_EXPIRES_IN=15
JWT_RT_SECRET=your-refresh-token-secret
JWT_RT_EXPIRES_IN=10080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=myapp
DB_SSL_MODE=disable
```

## API Endpoints

### Health Check
```
GET /health
```

### Example Feature (Items)
```
GET    /api/v1/items      # List all items
POST   /api/v1/items      # Create item
GET    /api/v1/items/:id  # Get item by ID
PUT    /api/v1/items/:id  # Update item
DELETE /api/v1/items/:id  # Delete item
```

## API Response Format

### Success Response
```json
{
  "success": true,
  "message": "Request successful",
  "data": { ... },
  "request_id": "uuid",
  "timestamp": "2025-01-01T00:00:00Z"
}
```

### Error Response
```json
{
  "success": false,
  "message": "Error description",
  "error_code": "ERR_BAD_REQUEST",
  "request_id": "uuid",
  "timestamp": "2025-01-01T00:00:00Z"
}
```

---

## Architecture Guide

### Layer Flow

```
HTTP Request → Router → Middleware → Handler → Service → Repository → Database
                                         ↓
HTTP Response ← Response Helper ← Handler ← Service ← Repository
```

### Layer Responsibilities

| Layer | Responsibility |
|-------|----------------|
| **Handler** | Parse requests, validate input, call service, return responses |
| **Service** | Business logic, orchestrate repositories, transform data |
| **Repository** | Database operations only (CRUD) |
| **Model** | Database entity with GORM tags |

---

## Adding a New Feature

Follow this step-by-step guide to add a new feature (e.g., "Users"):

### Step 1: Create Feature Directory

```bash
mkdir internal/features/users
```

### Step 2: Create Model (`model.go`)

```go
package users

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Email     string    `gorm:"type:varchar(255);uniqueIndex;not null"`
    Name      string    `gorm:"type:varchar(255);not null"`
    Password  string    `gorm:"type:varchar(255);not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
    return "users"
}
```

### Step 3: Create Repository (`repository.go`)

```go
package users

import (
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) Create(user *User) error {
    return r.db.Create(user).Error
}

func (r *Repository) FindByID(id uuid.UUID) (*User, error) {
    var user User
    err := r.db.Where("id = ?", id).First(&user).Error
    return &user, err
}

func (r *Repository) FindByEmail(email string) (*User, error) {
    var user User
    err := r.db.Where("email = ?", email).First(&user).Error
    return &user, err
}

func (r *Repository) FindAll(limit, offset int) ([]User, error) {
    var users []User
    err := r.db.Limit(limit).Offset(offset).Find(&users).Error
    return users, err
}

func (r *Repository) Update(user *User) error {
    return r.db.Save(user).Error
}

func (r *Repository) Delete(id uuid.UUID) error {
    return r.db.Delete(&User{}, "id = ?", id).Error
}
```

### Step 4: Create Service (`service.go`)

```go
package users

import (
    "github.com/google/uuid"
)

type Service struct {
    repo *Repository
}

func NewService(repo *Repository) *Service {
    return &Service{repo: repo}
}

// Request DTOs
type CreateUserRequest struct {
    Email    string `json:"email" validate:"required,email"`
    Name     string `json:"name" validate:"required,min=2,max=100"`
    Password string `json:"password" validate:"required,min=8"`
}

type UpdateUserRequest struct {
    Name string `json:"name" validate:"omitempty,min=2,max=100"`
}

// Response DTO
type UserResponse struct {
    ID        uuid.UUID `json:"id"`
    Email     string    `json:"email"`
    Name      string    `json:"name"`
    CreatedAt string    `json:"created_at"`
}

func (s *Service) Create(req CreateUserRequest) (*UserResponse, error) {
    user := &User{
        Email:    req.Email,
        Name:     req.Name,
        Password: hashPassword(req.Password), // implement this
    }

    if err := s.repo.Create(user); err != nil {
        return nil, err
    }

    return toResponse(user), nil
}

func (s *Service) GetByID(id uuid.UUID) (*UserResponse, error) {
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, err
    }
    return toResponse(user), nil
}

func (s *Service) GetAll(limit, offset int) ([]UserResponse, error) {
    users, err := s.repo.FindAll(limit, offset)
    if err != nil {
        return nil, err
    }

    responses := make([]UserResponse, len(users))
    for i, user := range users {
        responses[i] = *toResponse(&user)
    }
    return responses, nil
}

func toResponse(user *User) *UserResponse {
    return &UserResponse{
        ID:        user.ID,
        Email:     user.Email,
        Name:      user.Name,
        CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
    }
}

func hashPassword(password string) string {
    // TODO: implement bcrypt hashing
    return password
}
```

### Step 5: Create Handler (`handler.go`)

```go
package users

import (
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "gorm.io/gorm"

    "github.com/SuperIntelligence-Labs/go-backend-template/internal/response"
)

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) Create(c echo.Context) error {
    var req CreateUserRequest
    if err := c.Bind(&req); err != nil {
        return response.ErrBadRequest("Invalid request body", nil)
    }

    if err := c.Validate(&req); err != nil {
        return err
    }

    user, err := h.service.Create(req)
    if err != nil {
        return response.ErrInternalError(err)
    }

    return response.Created(c, "User created successfully", user)
}

func (h *Handler) GetByID(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return response.ErrBadRequest("Invalid user ID", nil)
    }

    user, err := h.service.GetByID(id)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return response.ErrNotFound("User not found")
        }
        return response.ErrInternalError(err)
    }

    return response.OK(c, "User retrieved successfully", user)
}

func (h *Handler) GetAll(c echo.Context) error {
    users, err := h.service.GetAll(20, 0)
    if err != nil {
        return response.ErrInternalError(err)
    }

    return response.OK(c, "Users retrieved successfully", users)
}
```

### Step 6: Create Routes (`routes.go`)

```go
package users

import "github.com/labstack/echo/v4"

func RegisterRoutes(g *echo.Group, h *Handler) {
    g.POST("", h.Create)
    g.GET("", h.GetAll)
    g.GET("/:id", h.GetByID)
}
```

### Step 7: Wire Up in `cmd/api/main.go`

Add to imports:
```go
"github.com/SuperIntelligence-Labs/go-backend-template/internal/features/users"
```

Add DI setup after database connection:
```go
// Users Feature
userRepo := users.NewRepository(db)
userService := users.NewService(userRepo)
userHandler := users.NewHandler(userService)
```

Add to AutoMigrate:
```go
err = db.AutoMigrate(&example.Item{}, &users.User{})
```

Update RoutesConfig:
```go
srv.RegisterRoutes(server.RoutesConfig{
    ExampleHandler: exampleHandler,
    UserHandler:    userHandler,
})
```

### Step 8: Update Router (`internal/server/router.go`)

```go
type RoutesConfig struct {
    ExampleHandler *example.Handler
    UserHandler    *users.Handler
}

func (s *Server) RegisterRoutes(cfg RoutesConfig) {
    api := s.Echo.Group("/api/v1")

    // Example feature
    itemsGroup := api.Group("/items")
    example.RegisterRoutes(itemsGroup, cfg.ExampleHandler)

    // Users feature
    usersGroup := api.Group("/users")
    users.RegisterRoutes(usersGroup, cfg.UserHandler)
}
```

---

## Using JWT Middleware

### Protect Routes

```go
import "github.com/SuperIntelligence-Labs/go-backend-template/internal/middleware"

// In router.go
usersGroup := api.Group("/users")
usersGroup.Use(middleware.JWTMiddleware(cfg.JWTSecret))
users.RegisterRoutes(usersGroup, cfg.UserHandler)
```

### Access Claims in Handler

```go
func (h *Handler) GetProfile(c echo.Context) error {
    claims, err := middleware.GetClaims(c)
    if err != nil {
        return err
    }

    userID := claims.UserID
    // ... use userID
}
```

---

## Response Helpers

### Success Responses

```go
response.OK(c, "Message", data)           // 200
response.Created(c, "Message", data)      // 201
response.Accepted(c, "Message", data)     // 202
response.NoContent(c)                     // 204
```

### Error Responses

```go
response.ErrBadRequest("Message", details)    // 400
response.ErrUnauthorized("Message")           // 401
response.ErrForbidden("Message")              // 403
response.ErrNotFound("Message")               // 404
response.ErrConflict("Message")               // 409
response.ErrValidationFailed(validationErrs)  // 422
response.ErrInternalError(err)                // 500
```

---

## Development

### Running Tests
```bash
go test ./...
go test -cover ./...
go test -v ./internal/features/example
```

### Building
```bash
go build -o bin/api ./cmd/api
go build -ldflags="-s -w" -o bin/api ./cmd/api  # optimized
```

### Database Migrations

Create migration files in `migrations/`:
```sql
-- migrations/000001_create_users_table.up.sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- migrations/000001_create_users_table.down.sql
DROP TABLE IF EXISTS users;
```

---

## License

MIT License - feel free to use this template for your projects.
