package server

import (
	"github.com/SuperIntelligence-Labs/go-backend-template/internal/features/example"
)

type RoutesConfig struct {
	ExampleHandler *example.Handler
}

func (s *Server) RegisterRoutes(cfg RoutesConfig) {
	api := s.Echo.Group("/api/v1")

	// Example feature routes
	itemsGroup := api.Group("/items")
	example.RegisterRoutes(itemsGroup, cfg.ExampleHandler)
}
