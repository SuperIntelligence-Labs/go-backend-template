package server

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	appMiddleware "github.com/SuperIntelligence-Labs/go-backend-template/internal/middleware"
	"github.com/SuperIntelligence-Labs/go-backend-template/internal/response"
)

// Server wraps the Echo HTTP server.
type Server struct {
	Echo *echo.Echo
}

func New() *Server {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = response.ErrorHandler
	e.Validator = response.NewValidator()

	// Middleware
	e.Use(middleware.RequestID())
	e.Use(appMiddleware.Zerolog())
	e.Use(middleware.Recover())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 30 * time.Second,
	}))

	// Health endpoint
	e.GET("/health", func(c echo.Context) error {
		return response.OK(c, "Server is healthy and running", map[string]string{
			"status": "healthy",
		})
	})

	// Invalid route handler
	e.Any("/*", func(c echo.Context) error {
		return response.ErrNotFound("Route not found")
	})

	return &Server{Echo: e}
}

// Start begins listening and handles graceful shutdown on SIGINT/SIGTERM.
func (s *Server) Start(addr string) error {
	go func() {
		if err := s.Echo.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.Echo.Logger.Fatalf("listen err: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.Echo.Shutdown(ctx)
}
