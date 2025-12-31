package example

import (
	"github.com/google/uuid"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateItemRequest represents the request payload for creating an item
type CreateItemRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

// UpdateItemRequest represents the request payload for updating an item
// Uses pointer fields to distinguish between "not provided" and "set to empty"
type UpdateItemRequest struct {
	Name        *string `json:"name" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
}

// ItemResponse represents the response payload for an item
type ItemResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
}

func (s *Service) Create(req CreateItemRequest) (*ItemResponse, error) {
	item := &Item{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repo.Create(item); err != nil {
		return nil, err
	}

	return toResponse(item), nil
}

func (s *Service) GetByID(id uuid.UUID) (*ItemResponse, error) {
	item, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return toResponse(item), nil
}

func (s *Service) GetAll(limit, offset int) ([]ItemResponse, error) {
	items, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]ItemResponse, len(items))
	for i, item := range items {
		responses[i] = *toResponse(&item)
	}

	return responses, nil
}

func (s *Service) Update(id uuid.UUID, req UpdateItemRequest) (*ItemResponse, error) {
	// Build update map for atomic update (fixes race condition)
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}

	// Perform atomic update
	if err := s.repo.UpdateFields(id, updates); err != nil {
		return nil, err
	}

	// Fetch and return the updated item
	return s.GetByID(id)
}

func (s *Service) Delete(id uuid.UUID) (int64, error) {
	return s.repo.Delete(id)
}

func toResponse(item *Item) *ItemResponse {
	if item == nil {
		return nil
	}
	return &ItemResponse{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   item.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
