package category

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, cat *Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*Category, error)
	List(ctx context.Context) ([]Category, error)
	Update(ctx context.Context, cat *Category) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, cat *Category) error {
	return r.db.WithContext(ctx).Create(cat).Error
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Category, error) {
	var cat Category
	err := r.db.WithContext(ctx).First(&cat, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *repository) List(ctx context.Context) ([]Category, error) {
	var cats []Category
	err := r.db.WithContext(ctx).Order("name ASC").Find(&cats).Error
	return cats, err
}

func (r *repository) Update(ctx context.Context, cat *Category) error {
	return r.db.WithContext(ctx).Save(cat).Error
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&Category{}, "id = ?", id).Error
}
