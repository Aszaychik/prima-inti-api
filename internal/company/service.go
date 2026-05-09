package company

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrCompanyNotFound      = errors.New("company profile not found")
	ErrCompanyAlreadyExists = errors.New("company profile already exists")
	ErrLinkNotFound         = errors.New("external link not found")
)

type Service interface {
	GetCompanyProfile(ctx context.Context) (*CompanyProfile, error)
	CreateCompanyProfile(ctx context.Context, req *CreateCompanyRequest) (*CompanyProfile, error)
	UpdateCompanyProfile(ctx context.Context, req *UpdateCompanyRequest) (*CompanyProfile, error)

	AddExternalLink(ctx context.Context, companyID uuid.UUID, req *CreateLinkRequest) (*ExternalLink, error)
	UpdateExternalLink(ctx context.Context, linkID uuid.UUID, req *UpdateLinkRequest) (*ExternalLink, error)
	DeleteExternalLink(ctx context.Context, linkID uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetCompanyProfile(ctx context.Context) (*CompanyProfile, error) {
	return s.repo.GetCompany(ctx)
}

func (s *service) CreateCompanyProfile(ctx context.Context, req *CreateCompanyRequest) (*CompanyProfile, error) {
	existing, _ := s.repo.GetCompany(ctx)
	if existing != nil {
		return nil, ErrCompanyAlreadyExists
	}

	comp := &CompanyProfile{
		Name:    req.Name,
		Phone:   req.Phone,
		Email:   req.Email,
		Address: req.Address,
	}
	if err := s.repo.CreateCompany(ctx, comp); err != nil {
		return nil, err
	}
	return comp, nil
}

func (s *service) UpdateCompanyProfile(ctx context.Context, req *UpdateCompanyRequest) (*CompanyProfile, error) {
	comp, err := s.repo.GetCompany(ctx)
	if err != nil {
		return nil, err
	}
	if comp == nil {
		return nil, ErrCompanyNotFound
	}

	if req.Name != nil {
		comp.Name = *req.Name
	}
	if req.Phone != nil {
		comp.Phone = *req.Phone
	}
	if req.Email != nil {
		comp.Email = *req.Email
	}
	if req.Address != nil {
		comp.Address = *req.Address
	}

	if err := s.repo.UpdateCompany(ctx, comp); err != nil {
		return nil, err
	}
	return comp, nil
}

func (s *service) AddExternalLink(ctx context.Context, companyID uuid.UUID, req *CreateLinkRequest) (*ExternalLink, error) {
	link := &ExternalLink{
		CompanyID: companyID,
		Platform:  req.Platform,
		URL:       req.URL,
	}
	if err := s.repo.CreateLink(ctx, link); err != nil {
		return nil, err
	}
	return link, nil
}

func (s *service) UpdateExternalLink(ctx context.Context, linkID uuid.UUID, req *UpdateLinkRequest) (*ExternalLink, error) {
	link, err := s.repo.GetLinkByID(ctx, linkID)
	if err != nil {
		return nil, ErrLinkNotFound
	}

	if req.Platform != nil {
		link.Platform = *req.Platform
	}
	if req.URL != nil {
		link.URL = *req.URL
	}

	if err := s.repo.UpdateLink(ctx, link); err != nil {
		return nil, err
	}
	return link, nil
}

func (s *service) DeleteExternalLink(ctx context.Context, linkID uuid.UUID) error {
	if err := s.repo.DeleteLink(ctx, linkID); err != nil {
		return ErrLinkNotFound
	}
	return nil
}
