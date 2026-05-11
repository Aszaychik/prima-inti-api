package category

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrCategoryNotFound      = errors.New("category not found")
	ErrCategoryAlreadyExists = errors.New("category already exists")
)

type Service interface {
	Create(ctx context.Context, req *CreateCategoryRequest) (*Category, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Category, error)
	List(ctx context.Context) ([]Category, error)
	Update(ctx context.Context, id uuid.UUID, req *UpdateCategoryRequest) (*Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req *CreateCategoryRequest) (*Category, error) {
	cat := &Category{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := s.repo.Create(ctx, cat); err != nil {
		// check for unique violation (depends on driver, simplify for now)
		return nil, ErrCategoryAlreadyExists
	}
	return cat, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*Category, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) List(ctx context.Context) ([]Category, error) {
	return s.repo.List(ctx)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, req *UpdateCategoryRequest) (*Category, error) {
	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrCategoryNotFound
	}
	if req.Name != nil {
		cat.Name = *req.Name
	}
	if req.Description != nil {
		cat.Description = *req.Description
	}
	if err := s.repo.Update(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return ErrCategoryNotFound
	}
	return nil
}
