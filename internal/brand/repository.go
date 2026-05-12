package brand

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, b *Brand) error
	GetByID(ctx context.Context, id uuid.UUID) (*Brand, error)
	List(ctx context.Context) ([]Brand, error)
	Update(ctx context.Context, b *Brand) error
	Delete(ctx context.Context, id uuid.UUID) error
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, b *Brand) error {
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Brand, error) {
	var b Brand
	err := r.db.WithContext(ctx).First(&b, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *repository) List(ctx context.Context) ([]Brand, error) {
	var brands []Brand
	err := r.db.WithContext(ctx).Order("name ASC").Find(&brands).Error
	return brands, err
}

func (r *repository) Update(ctx context.Context, b *Brand) error {
	return r.db.WithContext(ctx).Save(b).Error
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&Brand{}, "id = ?", id).Error
}

func (r *repository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Brand{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
