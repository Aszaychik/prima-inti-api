package company

import (
	"context"

	"github.com/google/uuid"
)

// MockRepository is a mock implementation of Repository interface
type MockRepository struct {
	GetCompanyFunc        func(ctx context.Context) (*CompanyProfile, error)
	CreateCompanyFunc     func(ctx context.Context, comp *CompanyProfile) error
	UpdateCompanyFunc     func(ctx context.Context, comp *CompanyProfile) error
	GetLinksByCompanyFunc func(ctx context.Context, companyID uuid.UUID) ([]ExternalLink, error)
	CreateLinkFunc        func(ctx context.Context, link *ExternalLink) error
	UpdateLinkFunc        func(ctx context.Context, link *ExternalLink) error
	DeleteLinkFunc        func(ctx context.Context, linkID uuid.UUID) error
	GetLinkByIDFunc       func(ctx context.Context, linkID uuid.UUID) (*ExternalLink, error)
}

func (m *MockRepository) GetCompany(ctx context.Context) (*CompanyProfile, error) {
	if m.GetCompanyFunc != nil {
		return m.GetCompanyFunc(ctx)
	}
	return nil, nil
}

func (m *MockRepository) CreateCompany(ctx context.Context, comp *CompanyProfile) error {
	if m.CreateCompanyFunc != nil {
		return m.CreateCompanyFunc(ctx, comp)
	}
	return nil
}

func (m *MockRepository) UpdateCompany(ctx context.Context, comp *CompanyProfile) error {
	if m.UpdateCompanyFunc != nil {
		return m.UpdateCompanyFunc(ctx, comp)
	}
	return nil
}

func (m *MockRepository) GetLinksByCompany(ctx context.Context, companyID uuid.UUID) ([]ExternalLink, error) {
	if m.GetLinksByCompanyFunc != nil {
		return m.GetLinksByCompanyFunc(ctx, companyID)
	}
	return nil, nil
}

func (m *MockRepository) CreateLink(ctx context.Context, link *ExternalLink) error {
	if m.CreateLinkFunc != nil {
		return m.CreateLinkFunc(ctx, link)
	}
	return nil
}

func (m *MockRepository) UpdateLink(ctx context.Context, link *ExternalLink) error {
	if m.UpdateLinkFunc != nil {
		return m.UpdateLinkFunc(ctx, link)
	}
	return nil
}

func (m *MockRepository) DeleteLink(ctx context.Context, linkID uuid.UUID) error {
	if m.DeleteLinkFunc != nil {
		return m.DeleteLinkFunc(ctx, linkID)
	}
	return nil
}

func (m *MockRepository) GetLinkByID(ctx context.Context, linkID uuid.UUID) (*ExternalLink, error) {
	if m.GetLinkByIDFunc != nil {
		return m.GetLinkByIDFunc(ctx, linkID)
	}
	return nil, nil
}

// MockService is a mock implementation of Service interface
type MockService struct {
	GetCompanyProfileFunc    func(ctx context.Context) (*CompanyProfile, error)
	CreateCompanyProfileFunc func(ctx context.Context, req *CreateCompanyRequest) (*CompanyProfile, error)
	UpdateCompanyProfileFunc func(ctx context.Context, req *UpdateCompanyRequest) (*CompanyProfile, error)
	AddExternalLinkFunc      func(ctx context.Context, companyID uuid.UUID, req *CreateLinkRequest) (*ExternalLink, error)
	UpdateExternalLinkFunc   func(ctx context.Context, linkID uuid.UUID, req *UpdateLinkRequest) (*ExternalLink, error)
	DeleteExternalLinkFunc   func(ctx context.Context, linkID uuid.UUID) error
}

func (m *MockService) GetCompanyProfile(ctx context.Context) (*CompanyProfile, error) {
	if m.GetCompanyProfileFunc != nil {
		return m.GetCompanyProfileFunc(ctx)
	}
	return nil, nil
}

func (m *MockService) CreateCompanyProfile(ctx context.Context, req *CreateCompanyRequest) (*CompanyProfile, error) {
	if m.CreateCompanyProfileFunc != nil {
		return m.CreateCompanyProfileFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockService) UpdateCompanyProfile(ctx context.Context, req *UpdateCompanyRequest) (*CompanyProfile, error) {
	if m.UpdateCompanyProfileFunc != nil {
		return m.UpdateCompanyProfileFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockService) AddExternalLink(ctx context.Context, companyID uuid.UUID, req *CreateLinkRequest) (*ExternalLink, error) {
	if m.AddExternalLinkFunc != nil {
		return m.AddExternalLinkFunc(ctx, companyID, req)
	}
	return nil, nil
}

func (m *MockService) UpdateExternalLink(ctx context.Context, linkID uuid.UUID, req *UpdateLinkRequest) (*ExternalLink, error) {
	if m.UpdateExternalLinkFunc != nil {
		return m.UpdateExternalLinkFunc(ctx, linkID, req)
	}
	return nil, nil
}

func (m *MockService) DeleteExternalLink(ctx context.Context, linkID uuid.UUID) error {
	if m.DeleteExternalLinkFunc != nil {
		return m.DeleteExternalLinkFunc(ctx, linkID)
	}
	return nil
}
