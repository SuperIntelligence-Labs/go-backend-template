package middleware

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"

	"github.com/SuperIntelligence-Labs/go-backend-template/internal/response"
)

type JWTClaims struct {
	UserID    uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role,omitempty"`
	TokenType string    `json:"token_type"`
	jwt.RegisteredClaims
}

// JWTMiddleware returns a configured JWT middleware with custom error handling
func JWTMiddleware(secretKey string) echo.MiddlewareFunc {
	config := echojwt.Config{
		SigningKey:  []byte(secretKey),
		TokenLookup: "header:Authorization:Bearer ,query:token",

		// Custom error handler that uses your AppError format
		ErrorHandler: func(c echo.Context, err error) error {
			return response.ErrUnauthorized("Invalid or missing authentication token")
		},

		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JWTClaims)
		},
	}

	return echojwt.WithConfig(config)
}

// GetClaims extracts JWT claims from the Echo context
func GetClaims(c echo.Context) (*JWTClaims, error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return nil, response.ErrUnauthorized("JWT token missing or invalid")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, response.ErrUnauthorized("Failed to validate token")
	}

	return claims, nil
}

// ValidateRefreshToken validates a refresh token and returns claims
func ValidateRefreshToken(tokenString, secretKey string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid || claims.TokenType != "refresh" {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// RequireRole is a middleware that checks if user has a specific role
func RequireRole(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := GetClaims(c)
			if err != nil {
				return err
			}

			for _, role := range allowedRoles {
				if claims.Role == role {
					return next(c)
				}
			}

			return response.ErrForbidden("Insufficient permissions")
		}
	}
}
