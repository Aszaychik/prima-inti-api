package product

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, p *Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	List(ctx context.Context) ([]Product, error)
	ListByBrand(ctx context.Context, brandID uuid.UUID) ([]Product, error)
	ListByCategory(ctx context.Context, categoryID uuid.UUID) ([]Product, error)
	ListBySeries(ctx context.Context, seriesID uuid.UUID) ([]Product, error)
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, p *Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	var p Product
	err := r.db.WithContext(ctx).First(&p, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) List(ctx context.Context) ([]Product, error) {
	var list []Product
	err := r.db.WithContext(ctx).Order("model ASC").Find(&list).Error
	return list, err
}

func (r *repository) ListByBrand(ctx context.Context, brandID uuid.UUID) ([]Product, error) {
	var list []Product
	err := r.db.WithContext(ctx).Where("brand_id = ?", brandID).Order("model ASC").Find(&list).Error
	return list, err
}

func (r *repository) ListByCategory(ctx context.Context, categoryID uuid.UUID) ([]Product, error) {
	var list []Product
	err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Order("model ASC").Find(&list).Error
	return list, err
}

func (r *repository) ListBySeries(ctx context.Context, seriesID uuid.UUID) ([]Product, error) {
	var list []Product
	err := r.db.WithContext(ctx).Where("series_id = ?", seriesID).Order("model ASC").Find(&list).Error
	return list, err
}

func (r *repository) Update(ctx context.Context, p *Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&Product{}, "id = ?", id).Error
}
