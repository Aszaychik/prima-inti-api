package series

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrSeriesNotFound      = errors.New("series not found")
	ErrSeriesAlreadyExists = errors.New("series already exists for this brand")
)

type Service interface {
	Create(ctx context.Context, req *CreateSeriesRequest) (*Series, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Series, error)
	List(ctx context.Context) ([]Series, error)
	ListByBrand(ctx context.Context, brandID uuid.UUID) ([]Series, error)
	Update(ctx context.Context, id uuid.UUID, req *UpdateSeriesRequest) (*Series, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req *CreateSeriesRequest) (*Series, error) {
	ser := &Series{
		BrandID: req.BrandID,
		Name:    req.Name,
	}
	if err := s.repo.Create(ctx, ser); err != nil {
		// Unique constraint violation on (brand_id, name) would be caught here
		return nil, ErrSeriesAlreadyExists
	}
	return ser, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Series, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) List(ctx context.Context) ([]Series, error) {
	return s.repo.List(ctx)
}

func (s *service) ListByBrand(ctx context.Context, brandID uuid.UUID) ([]Series, error) {
	return s.repo.ListByBrand(ctx, brandID)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, req *UpdateSeriesRequest) (*Series, error) {
	ser, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrSeriesNotFound
	}
	if req.Name != nil {
		ser.Name = *req.Name
	}
	if err := s.repo.Update(ctx, ser); err != nil {
		return nil, err
	}
	return ser, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return ErrSeriesNotFound
	}
	return nil
}
