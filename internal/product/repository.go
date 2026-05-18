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
	ListWithFilters(ctx context.Context, filters ProductFilters) ([]Product, error)
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
	err := r.db.WithContext(ctx).
		Preload("Brand").
		Preload("Category").
		Preload("Series").
		First(&p, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) List(ctx context.Context) ([]Product, error) {
	var list []Product
	err := r.db.WithContext(ctx).
		Preload("Brand").
		Preload("Category").
		Preload("Series").
		Order("model ASC").
		Find(&list).Error
	return list, err
}

func (r *repository) ListWithFilters(ctx context.Context, filters ProductFilters) ([]Product, error) {
	query := r.db.WithContext(ctx).
		Preload("Brand").
		Preload("Category").
		Preload("Series")

	// Filter by model (LIKE)
	if filters.Model != nil && *filters.Model != "" {
		query = query.Where("model ILIKE ?", "%"+*filters.Model+"%")
	}

	// Filter by brand (either by ID or name)
	if filters.BrandID != nil {
		query = query.Where("brand_id = ?", *filters.BrandID)
	} else if filters.BrandName != nil && *filters.BrandName != "" {
		// Need to join brands table to search by name
		query = query.Joins("JOIN brands ON brands.id = products.brand_id").
			Where("brands.name ILIKE ?", "%"+*filters.BrandName+"%")
	}

	// Filter by category
	if filters.CategoryID != nil {
		query = query.Where("category_id = ?", *filters.CategoryID)
	} else if filters.CategoryName != nil && *filters.CategoryName != "" {
		query = query.Joins("JOIN categories ON categories.id = products.category_id").
			Where("categories.name ILIKE ?", "%"+*filters.CategoryName+"%")
	}

	// Filter by series (optional, could be null)
	if filters.SeriesID != nil {
		query = query.Where("series_id = ?", *filters.SeriesID)
	} else if filters.SeriesName != nil && *filters.SeriesName != "" {
		// Join series table (left join to allow null series)
		query = query.Joins("LEFT JOIN series ON series.id = products.series_id").
			Where("series.name ILIKE ?", "%"+*filters.SeriesName+"%")
	}

	var products []Product
	err := query.Order("model ASC").Find(&products).Error
	return products, err
}

func (r *repository) ListByBrand(ctx context.Context, brandID uuid.UUID) ([]Product, error) {
	var list []Product
	err := r.db.WithContext(ctx).
		Preload("Brand").
		Preload("Category").
		Preload("Series").
		Where("brand_id = ?", brandID).
		Order("model ASC").
		Find(&list).Error
	return list, err
}

func (r *repository) ListByCategory(ctx context.Context, categoryID uuid.UUID) ([]Product, error) {
	var list []Product
	err := r.db.WithContext(ctx).
		Preload("Brand").
		Preload("Category").
		Preload("Series").
		Where("category_id = ?", categoryID).
		Order("model ASC").
		Find(&list).Error
	return list, err
}

func (r *repository) ListBySeries(ctx context.Context, seriesID uuid.UUID) ([]Product, error) {
	var list []Product
	err := r.db.WithContext(ctx).
		Preload("Brand").
		Preload("Category").
		Preload("Series").
		Where("series_id = ?", seriesID).
		Order("model ASC").
		Find(&list).Error
	return list, err
}

func (r *repository) Update(ctx context.Context, p *Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&Product{}, "id = ?", id).Error
}
