package product

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrBrandNotFound    = errors.New("brand not found")
	ErrCategoryNotFound = errors.New("category not found")
	ErrSeriesNotFound   = errors.New("series not found")
)

type brandChecker interface {
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

type categoryChecker interface {
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

type seriesChecker interface {
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
}

type Service interface {
	Create(ctx context.Context, req *CreateProductRequest) (*Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	List(ctx context.Context) ([]Product, error)
	ListWithFilters(ctx context.Context, filters ProductFilters) ([]Product, error)
	ListByBrand(ctx context.Context, brandID uuid.UUID) ([]Product, error)
	ListByCategory(ctx context.Context, categoryID uuid.UUID) ([]Product, error)
	ListBySeries(ctx context.Context, seriesID uuid.UUID) ([]Product, error)
	Update(ctx context.Context, id uuid.UUID, req *UpdateProductRequest) (*Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo         Repository
	brandRepo    brandChecker
	categoryRepo categoryChecker
	seriesRepo   seriesChecker
}

func NewService(repo Repository, brandRepo brandChecker, categoryRepo categoryChecker, seriesRepo seriesChecker) Service {
	return &service{
		repo:         repo,
		brandRepo:    brandRepo,
		categoryRepo: categoryRepo,
		seriesRepo:   seriesRepo,
	}
}

func (s *service) Create(ctx context.Context, req *CreateProductRequest) (*Product, error) {
	// Validate brand exists
	exists, err := s.brandRepo.Exists(ctx, req.BrandID)
	if err != nil || !exists {
		return nil, ErrBrandNotFound
	}
	// Validate category exists
	exists, err = s.categoryRepo.Exists(ctx, req.CategoryID)
	if err != nil || !exists {
		return nil, ErrCategoryNotFound
	}
	// Validate series if provided
	if req.SeriesID != nil {
		exists, err = s.seriesRepo.Exists(ctx, *req.SeriesID)
		if err != nil || !exists {
			return nil, ErrSeriesNotFound
		}
	}
	prod := &Product{
		BrandID:     req.BrandID,
		CategoryID:  req.CategoryID,
		SeriesID:    req.SeriesID,
		Model:       req.Model,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		ImageURL:    req.ImageURL,
	}
	if err := s.repo.Create(ctx, prod); err != nil {
		return nil, err
	}
	return prod, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) List(ctx context.Context) ([]Product, error) {
	return s.repo.List(ctx)
}

func (s *service) ListWithFilters(ctx context.Context, filters ProductFilters) ([]Product, error) {
	return s.repo.ListWithFilters(ctx, filters)
}

func (s *service) ListByBrand(ctx context.Context, brandID uuid.UUID) ([]Product, error) {
	return s.repo.ListByBrand(ctx, brandID)
}

func (s *service) ListByCategory(ctx context.Context, categoryID uuid.UUID) ([]Product, error) {
	return s.repo.ListByCategory(ctx, categoryID)
}

func (s *service) ListBySeries(ctx context.Context, seriesID uuid.UUID) ([]Product, error) {
	return s.repo.ListBySeries(ctx, seriesID)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, req *UpdateProductRequest) (*Product, error) {
	prod, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrProductNotFound
	}
	if req.BrandID != nil {
		exists, err := s.brandRepo.Exists(ctx, *req.BrandID)
		if err != nil || !exists {
			return nil, ErrBrandNotFound
		}
		prod.BrandID = *req.BrandID
	}
	if req.CategoryID != nil {
		exists, err := s.categoryRepo.Exists(ctx, *req.CategoryID)
		if err != nil || !exists {
			return nil, ErrCategoryNotFound
		}
		prod.CategoryID = *req.CategoryID
	}
	if req.SeriesID != nil {
		if *req.SeriesID == uuid.Nil {
			prod.SeriesID = nil
		} else {
			exists, err := s.seriesRepo.Exists(ctx, *req.SeriesID)
			if err != nil || !exists {
				return nil, ErrSeriesNotFound
			}
			prod.SeriesID = req.SeriesID
		}
	}
	if req.Model != nil {
		prod.Model = *req.Model
	}
	if req.Description != nil {
		prod.Description = *req.Description
	}
	if req.Price != nil {
		prod.Price = req.Price
	}
	if req.Stock != nil {
		prod.Stock = *req.Stock
	}
	if req.ImageURL != nil {
		prod.ImageURL = *req.ImageURL
	}
	if err := s.repo.Update(ctx, prod); err != nil {
		return nil, err
	}
	return prod, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return ErrProductNotFound
	}
	return nil
}
