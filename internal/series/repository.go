package series

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, s *Series) error
	GetByID(ctx context.Context, id uuid.UUID) (*Series, error)
	List(ctx context.Context) ([]Series, error)
	ListByBrand(ctx context.Context, brandID uuid.UUID) ([]Series, error)
	Update(ctx context.Context, s *Series) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, s *Series) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Series, error) {
	var s Series
	err := r.db.WithContext(ctx).First(&s, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) List(ctx context.Context) ([]Series, error) {
	var list []Series
	err := r.db.WithContext(ctx).Order("name ASC").Find(&list).Error
	return list, err
}

func (r *repository) ListByBrand(ctx context.Context, brandID uuid.UUID) ([]Series, error) {
	var list []Series
	err := r.db.WithContext(ctx).Where("brand_id = ?", brandID).Order("name ASC").Find(&list).Error
	return list, err
}

func (r *repository) Update(ctx context.Context, s *Series) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&Series{}, "id = ?", id).Error
}
