package example

import (
	"errors"
	"strconv"

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

// Create handles POST /items
func (h *Handler) Create(c echo.Context) error {
	var req CreateItemRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrBadRequest("Invalid request body", nil)
	}

	if err := c.Validate(&req); err != nil {
		details := response.ToValidationErrors(err)
		return response.ErrValidationFailed(details)
	}

	item, err := h.service.Create(req)
	if err != nil {
		return response.ErrInternalError(err)
	}

	return response.Created(c, "Item created successfully", item)
}

// GetByID handles GET /items/:id
func (h *Handler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ErrBadRequest("Invalid item ID", nil)
	}

	item, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ErrNotFound("Item not found")
		}
		return response.ErrInternalError(err)
	}

	return response.OK(c, "Item retrieved successfully", item)
}

// GetAll handles GET /items
func (h *Handler) GetAll(c echo.Context) error {
	// Parse pagination parameters with defaults and limits
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	// Apply defaults
	if limit <= 0 {
		limit = 20
	}
	// Enforce maximum limit to prevent DoS
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	items, err := h.service.GetAll(limit, offset)
	if err != nil {
		return response.ErrInternalError(err)
	}

	return response.OK(c, "Items retrieved successfully", items)
}

// Update handles PUT /items/:id
func (h *Handler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ErrBadRequest("Invalid item ID", nil)
	}

	var req UpdateItemRequest
	if err := c.Bind(&req); err != nil {
		return response.ErrBadRequest("Invalid request body", nil)
	}

	if err := c.Validate(&req); err != nil {
		details := response.ToValidationErrors(err)
		return response.ErrValidationFailed(details)
	}

	item, err := h.service.Update(id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return response.ErrNotFound("Item not found")
		}
		return response.ErrInternalError(err)
	}

	return response.OK(c, "Item updated successfully", item)
}

// Delete handles DELETE /items/:id
func (h *Handler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.ErrBadRequest("Invalid item ID", nil)
	}

	rowsAffected, err := h.service.Delete(id)
	if err != nil {
		return response.ErrInternalError(err)
	}

	if rowsAffected == 0 {
		return response.ErrNotFound("Item not found")
	}

	return response.NoContent(c)
}
