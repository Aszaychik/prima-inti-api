package brand

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrBrandNotFound      = errors.New("brand not found")
	ErrBrandAlreadyExists = errors.New("brand already exists")
)

type Service interface {
	Create(ctx context.Context, req *CreateBrandRequest) (*Brand, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Brand, error)
	List(ctx context.Context) ([]Brand, error)
	Update(ctx context.Context, id uuid.UUID, req *UpdateBrandRequest) (*Brand, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req *CreateBrandRequest) (*Brand, error) {
	b := &Brand{
		Name:    req.Name,
		LogoURL: req.LogoURL,
	}
	if err := s.repo.Create(ctx, b); err != nil {
		// For simplicity, treat any error as already exists (you can improve later)
		return nil, ErrBrandAlreadyExists
	}
	return b, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Brand, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) List(ctx context.Context) ([]Brand, error) {
	return s.repo.List(ctx)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, req *UpdateBrandRequest) (*Brand, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrBrandNotFound
	}
	if req.Name != nil {
		b.Name = *req.Name
	}
	if req.LogoURL != nil {
		b.LogoURL = *req.LogoURL
	}
	if err := s.repo.Update(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return ErrBrandNotFound
	}
	return nil
}
