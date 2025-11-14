# Go Backend Template

A production-ready Go backend template with clean architecture, built using the Gin web framework and following industry best practices.

## Features

- **Layer-based Architecture** - Clean separation of concerns with handler, service, repository, and model layers
- **Standardized API Responses** - Consistent success and error response formats with user-friendly messages
- **Hot Reload** - Development with Air for automatic rebuilds on file changes
- **Type-safe Validation** - Request validation using go-playground/validator with custom error messages
- **High Performance** - Sonic JSON serialization for faster response times
- **PostgreSQL Ready** - Structure prepared for database integration
- **Production Ready** - Proper error handling, logging, and configuration management

## Tech Stack

- **Go 1.25.4**
- **Gin** - HTTP web framework
- **Validator** - Request validation
- **Sonic** - High-performance JSON serialization
- **Air** - Hot reload for development
- **Viper** - Configuration management (ready to integrate)

## Project Structure

```
.
├── cmd/api/              # Application entry point
├── internal/
│   ├── handler/          # HTTP request handlers
│   ├── service/          # Business logic layer
│   ├── repository/       # Data access layer
│   ├── model/            # Domain models
│   ├── dto/              # Data transfer objects
│   │   ├── request/      # API request structures
│   │   └── response/     # API response structures
│   ├── middleware/       # Gin middleware
│   ├── response/         # Standard response helpers
│   ├── router/           # Route definitions
│   ├── database/         # Database connection
│   └── validator/        # Custom validators
├── config/               # Configuration files
├── migrations/           # Database migrations
├── scripts/              # Build and deployment scripts
└── pkg/                  # Reusable packages
```

## Architecture Guide

### Understanding the Layers

This template follows a **layer-based architecture** where each layer has a specific responsibility:

#### 1. `cmd/api/` - Application Entry Point
**Purpose:** Contains the main.go file that bootstraps the application.

**What goes here:**
- Application initialization
- Dependency injection setup
- Server configuration and startup

**Example:**
```go
package main

func main() {
    // Load config
    // Initialize database
    // Setup router
    // Start server
}
```

#### 2. `internal/handler/` - HTTP Handlers (Controllers)
**Purpose:** Handles HTTP requests and responses. This is the entry point for all API endpoints.

**Responsibilities:**
- Receive HTTP requests
- Validate request data using DTOs
- Call service layer methods
- Return standardized responses

**What goes here:**
- `user_handler.go` - User-related endpoints
- `auth_handler.go` - Authentication endpoints
- `product_handler.go` - Product endpoints

**Example:**
```go
func (h *UserHandler) GetUser(c *gin.Context) {
    id := c.Param("id")
    user, err := h.userService.GetByID(id)
    if err != nil {
        response.NotFound(c, "User")
        return
    }
    response.OK(c, user)
}
```

#### 3. `internal/service/` - Business Logic Layer
**Purpose:** Contains all business logic and rules. This is where your application's core functionality lives.

**Responsibilities:**
- Implement business rules
- Coordinate between multiple repositories
- Transform data between models and DTOs
- Handle complex operations

**What goes here:**
- `user_service.go` - User business logic
- `auth_service.go` - Authentication logic
- `email_service.go` - Email sending logic

**Example:**
```go
func (s *UserService) CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error) {
    // Business logic: hash password, validate uniqueness, etc.
    hashedPassword := hashPassword(req.Password)

    user := &model.User{
        Email:    req.Email,
        Password: hashedPassword,
    }

    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }

    return toUserResponse(user), nil
}
```

#### 4. `internal/repository/` - Data Access Layer
**Purpose:** Handles all database operations. This is the only layer that talks to the database.

**Responsibilities:**
- CRUD operations
- Database queries
- Transaction management
- Return domain models

**What goes here:**
- `user_repository.go` - User database operations
- `product_repository.go` - Product database operations

**Example:**
```go
func (r *UserRepository) Create(user *model.User) error {
    return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
    var user model.User
    err := r.db.Where("email = ?", email).First(&user).Error
    return &user, err
}
```

#### 5. `internal/model/` - Domain Models
**Purpose:** Defines your database entities and domain objects.

**What goes here:**
- Database table structures (with GORM tags)
- Domain entities

**Example:**
```go
type User struct {
    ID        uuid.UUID  `gorm:"type:uuid;primary_key"`
    Email     string     `gorm:"uniqueIndex;not null"`
    Password  string     `gorm:"not null"`
    Name      string     `gorm:"not null"`
    CreatedAt time.Time  `gorm:"autoCreateTime"`
    UpdatedAt time.Time  `gorm:"autoUpdateTime"`
}
```

#### 6. `internal/dto/` - Data Transfer Objects
**Purpose:** Defines the structure of data coming in (requests) and going out (responses) of your API.

**Why separate from models?** DTOs decouple your API contract from your database schema. You can change your database without breaking your API.

**What goes here:**

`dto/request/` - API Request structures:
```go
type CreateUserRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Name     string `json:"name" binding:"required"`
}
```

`dto/response/` - API Response structures:
```go
type UserResponse struct {
    ID        string `json:"id"`
    Email     string `json:"email"`
    Name      string `json:"name"`
    CreatedAt int64  `json:"created_at"`
}
```

#### 7. `internal/middleware/` - Gin Middleware
**Purpose:** Contains middleware functions that run before/after request handlers.

**What goes here:**
- `auth.go` - JWT authentication middleware
- `logger.go` - Request logging
- `cors.go` - CORS configuration
- `rate_limit.go` - Rate limiting

**Example:**
```go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if !validateToken(token) {
            response.Unauthorized(c, "")
            return
        }
        c.Next()
    }
}
```

#### 8. `internal/response/` - Standard Response Helpers
**Purpose:** Provides consistent response formats across all endpoints.

**Already implemented:**
- Success responses: `OK()`, `Created()`, `Updated()`, `Deleted()`
- Error responses: `BadRequest()`, `NotFound()`, `Unauthorized()`, etc.

**Usage:**
```go
response.OK(c, userData)
response.Created(c, newUser)
response.NotFound(c, "User")
response.ValidationError(c, err)
```

#### 9. `internal/router/` - Route Definitions
**Purpose:** Central place to define all your API routes.

**What goes here:**
- Route registration
- Middleware application
- Route grouping

**Example:**
```go
func SetupRouter(userHandler *handler.UserHandler) *gin.Engine {
    r := gin.Default()

    api := r.Group("/api/v1")
    {
        users := api.Group("/users")
        {
            users.GET("/:id", userHandler.GetUser)
            users.POST("", userHandler.CreateUser)
            users.PUT("/:id", userHandler.UpdateUser)
            users.DELETE("/:id", userHandler.DeleteUser)
        }
    }

    return r
}
```

#### 10. `internal/database/` - Database Connection
**Purpose:** Database initialization and connection management.

**What goes here:**
- `postgres.go` - PostgreSQL connection setup
- Connection pooling configuration
- Database health checks

**Example:**
```go
func NewPostgresDB(cfg *config.Config) (*gorm.DB, error) {
    dsn := cfg.GetDatabaseDSN()
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }
    return db, nil
}
```

#### 11. `internal/validator/` - Custom Validators
**Purpose:** Custom validation rules beyond the standard validator tags.

**What goes here:**
- Custom validation functions
- Business-specific validators

**Example:**
```go
func ValidateIndianPhone(fl validator.FieldLevel) bool {
    phone := fl.Field().String()
    matched, _ := regexp.MatchString(`^[6-9]\d{9}$`, phone)
    return matched
}
```

#### 12. `config/` - Configuration Files
**Purpose:** Application configuration files.

**What goes here:**
- `config.yaml` - Non-sensitive settings
- `development.yaml`, `production.yaml` - Environment-specific configs
- `.env.example` - Template for environment variables

#### 13. `migrations/` - Database Migrations
**Purpose:** Version-controlled database schema changes.

**What goes here:**
- SQL migration files
- Naming: `000001_create_users_table.up.sql`, `000001_create_users_table.down.sql`

**Example:**
```sql
-- 000001_create_users_table.up.sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### 14. `scripts/` - Build and Deployment Scripts
**Purpose:** Automation scripts for development and deployment.

**What goes here:**
- `migrate.sh` - Run database migrations
- `seed.sh` - Seed database with test data
- `deploy.sh` - Deployment scripts

#### 15. `pkg/` - Reusable Packages
**Purpose:** Code that could be extracted into a separate library.

**What goes here:**
- `httpclient/` - Reusable HTTP client
- `utils/` - General utility functions
- Only truly reusable, project-agnostic code

**Note:** Most application-specific code should go in `internal/`, not `pkg/`.

## Request Flow

Here's how a typical request flows through the architecture:

```
1. HTTP Request
   ↓
2. Router (matches route)
   ↓
3. Middleware (auth, logging, etc.)
   ↓
4. Handler (validates request DTO)
   ↓
5. Service (business logic)
   ↓
6. Repository (database operation)
   ↓
7. Service (transforms model to response DTO)
   ↓
8. Handler (uses response helper to send JSON)
   ↓
9. HTTP Response
```

## Adding a New Feature

**Example: Adding a "Post" feature**

1. **Model** - Create `internal/model/post.go`:
```go
type Post struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key"`
    Title     string    `gorm:"not null"`
    Content   string
    AuthorID  uuid.UUID `gorm:"type:uuid;not null"`
    CreatedAt time.Time
}
```

2. **DTOs** - Create request/response structs:
- `internal/dto/request/post_request.go`
- `internal/dto/response/post_response.go`

3. **Repository** - Create `internal/repository/post_repository.go`

4. **Service** - Create `internal/service/post_service.go`

5. **Handler** - Create `internal/handler/post_handler.go`

6. **Routes** - Add routes in `internal/router/router.go`

7. **Migration** - Create migration in `migrations/`

## Getting Started

### Prerequisites

- Go 1.25.4 or higher
- Air (for hot reload): `go install github.com/air-verse/air@latest`

### Installation

1. Clone the repository:
```bash
git clone <your-repo-url>
cd go-template
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application:
```bash
# Development (with hot reload)
air

# Production build
go build -o bin/api ./cmd/api
./bin/api
```

## API Response Format

### Success Response
```json
{
  "success": true,
  "message": "Request successful",
  "response": { ... },
  "timestamp": 1234567890
}
```

### Error Response
```json
{
  "success": false,
  "message": "Unable to process your request",
  "error": {
    "code": "BAD_REQUEST",
    "message": "Detailed error message",
    "details": { ... }
  },
  "timestamp": 1234567890
}
```

## Development

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/handler
```

### Building
```bash
# Build binary
go build -o bin/api ./cmd/api

# Build with optimizations
go build -ldflags="-s -w" -o bin/api ./cmd/api
```

## Configuration

Configuration can be managed using environment variables or config files (Viper setup ready).

### Environment Variables
- `APP_ENV` - Application environment (development, staging, production)
- `DATABASE_PASSWORD` - Database password
- `JWT_SECRET` - JWT signing secret

## Contributing

1. Create a feature branch
2. Make your changes
3. Add tests
4. Submit a pull request

## License

MIT License - feel free to use this template for your projects.
