package example

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

func (r *Repository) Create(item *Item) error {
	return r.db.Create(item).Error
}

func (r *Repository) FindByID(id uuid.UUID) (*Item, error) {
	var item Item
	err := r.db.Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *Repository) FindAll(limit, offset int) ([]Item, error) {
	var items []Item
	err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&items).Error
	return items, err
}

func (r *Repository) Update(item *Item) error {
	return r.db.Save(item).Error
}

// UpdateFields performs an atomic update of specific fields
func (r *Repository) UpdateFields(id uuid.UUID, fields map[string]interface{}) error {
	if len(fields) == 0 {
		return nil // No fields to update
	}
	result := r.db.Model(&Item{}).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *Repository) Delete(id uuid.UUID) (int64, error) {
	result := r.db.Delete(&Item{}, "id = ?", id)
	return result.RowsAffected, result.Error
}
