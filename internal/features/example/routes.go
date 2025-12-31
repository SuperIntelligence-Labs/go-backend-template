package example

import "github.com/labstack/echo/v4"

// RegisterRoutes registers all example feature routes
func RegisterRoutes(g *echo.Group, h *Handler) {
	g.POST("", h.Create)
	g.GET("", h.GetAll)
	g.GET("/:id", h.GetByID)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}
