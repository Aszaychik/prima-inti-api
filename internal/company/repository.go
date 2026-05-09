package company

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetCompany(ctx context.Context) (*CompanyProfile, error)
	CreateCompany(ctx context.Context, comp *CompanyProfile) error
	UpdateCompany(ctx context.Context, comp *CompanyProfile) error

	GetLinksByCompany(ctx context.Context, companyID uuid.UUID) ([]ExternalLink, error)
	CreateLink(ctx context.Context, link *ExternalLink) error
	UpdateLink(ctx context.Context, link *ExternalLink) error
	DeleteLink(ctx context.Context, linkID uuid.UUID) error
	GetLinkByID(ctx context.Context, linkID uuid.UUID) (*ExternalLink, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetCompany(ctx context.Context) (*CompanyProfile, error) {
	var company CompanyProfile
	err := r.db.WithContext(ctx).Preload("ExternalLinks").First(&company).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}

func (r *repository) CreateCompany(ctx context.Context, comp *CompanyProfile) error {
	return r.db.WithContext(ctx).Create(comp).Error
}

func (r *repository) UpdateCompany(ctx context.Context, comp *CompanyProfile) error {
	return r.db.WithContext(ctx).Save(comp).Error
}

func (r *repository) GetLinksByCompany(ctx context.Context, companyID uuid.UUID) ([]ExternalLink, error) {
	var links []ExternalLink
	err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Find(&links).Error
	return links, err
}

func (r *repository) CreateLink(ctx context.Context, link *ExternalLink) error {
	return r.db.WithContext(ctx).Create(link).Error
}

func (r *repository) UpdateLink(ctx context.Context, link *ExternalLink) error {
	return r.db.WithContext(ctx).Save(link).Error
}

func (r *repository) DeleteLink(ctx context.Context, linkID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&ExternalLink{}, "id = ?", linkID).Error
}

func (r *repository) GetLinkByID(ctx context.Context, linkID uuid.UUID) (*ExternalLink, error) {
	var link ExternalLink
	err := r.db.WithContext(ctx).First(&link, "id = ?", linkID).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}
