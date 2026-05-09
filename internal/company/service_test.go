package company

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestService_GetCompanyProfile_Exists(t *testing.T) {
	mockRepo := &MockRepository{}
	expectedCompany := &CompanyProfile{
		ID:      uuid.New(),
		Name:    "Existing Co",
		Phone:   "123",
		Email:   "exist@test.com",
		Address: "Somewhere",
	}
	mockRepo.GetCompanyFunc = func(ctx context.Context) (*CompanyProfile, error) {
		return expectedCompany, nil
	}
	svc := NewService(mockRepo)

	company, err := svc.GetCompanyProfile(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedCompany, company)
}

func TestService_GetCompanyProfile_NotFound(t *testing.T) {
	mockRepo := &MockRepository{}
	mockRepo.GetCompanyFunc = func(ctx context.Context) (*CompanyProfile, error) {
		return nil, nil
	}
	svc := NewService(mockRepo)

	company, err := svc.GetCompanyProfile(context.Background())
	assert.NoError(t, err)
	assert.Nil(t, company)
}

func TestService_CreateCompanyProfile_FirstTime(t *testing.T) {
	mockRepo := &MockRepository{}
	mockRepo.GetCompanyFunc = func(ctx context.Context) (*CompanyProfile, error) {
		return nil, nil // no existing company
	}
	var createdCompany *CompanyProfile
	mockRepo.CreateCompanyFunc = func(ctx context.Context, comp *CompanyProfile) error {
		createdCompany = comp
		return nil
	}
	svc := NewService(mockRepo)

	req := &CreateCompanyRequest{
		Name:    "New Co",
		Phone:   "999",
		Email:   "new@test.com",
		Address: "New Address",
	}
	company, err := svc.CreateCompanyProfile(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, company)
	assert.Equal(t, req.Name, company.Name)
	assert.Equal(t, req.Phone, company.Phone)
	assert.NotNil(t, createdCompany)
}

func TestService_CreateCompanyProfile_AlreadyExists(t *testing.T) {
	mockRepo := &MockRepository{}
	mockRepo.GetCompanyFunc = func(ctx context.Context) (*CompanyProfile, error) {
		return &CompanyProfile{ID: uuid.New()}, nil
	}
	svc := NewService(mockRepo)

	req := &CreateCompanyRequest{Name: "Should Fail"}
	_, err := svc.CreateCompanyProfile(context.Background(), req)
	assert.ErrorIs(t, err, ErrCompanyAlreadyExists)
}

func TestService_UpdateCompanyProfile_NotFound(t *testing.T) {
	mockRepo := &MockRepository{}
	mockRepo.GetCompanyFunc = func(ctx context.Context) (*CompanyProfile, error) {
		return nil, nil
	}
	svc := NewService(mockRepo)

	req := &UpdateCompanyRequest{Name: ptrString("Updated")}
	_, err := svc.UpdateCompanyProfile(context.Background(), req)
	assert.ErrorIs(t, err, ErrCompanyNotFound)
}

func TestService_AddExternalLink_Success(t *testing.T) {
	mockRepo := &MockRepository{}
	companyID := uuid.New()
	var createdLink *ExternalLink
	mockRepo.CreateLinkFunc = func(ctx context.Context, link *ExternalLink) error {
		createdLink = link
		return nil
	}
	svc := NewService(mockRepo)

	req := &CreateLinkRequest{
		Platform: "github",
		URL:      "https://github.com/test",
	}
	link, err := svc.AddExternalLink(context.Background(), companyID, req)
	assert.NoError(t, err)
	assert.Equal(t, companyID, link.CompanyID)
	assert.Equal(t, req.Platform, link.Platform)
	assert.Equal(t, req.URL, link.URL)
	assert.NotNil(t, createdLink)
}

func TestService_UpdateExternalLink_NotFound(t *testing.T) {
	mockRepo := &MockRepository{}
	mockRepo.GetLinkByIDFunc = func(ctx context.Context, linkID uuid.UUID) (*ExternalLink, error) {
		return nil, errors.New("not found")
	}
	svc := NewService(mockRepo)

	req := &UpdateLinkRequest{Platform: ptrString("newplatform")}
	_, err := svc.UpdateExternalLink(context.Background(), uuid.New(), req)
	assert.ErrorIs(t, err, ErrLinkNotFound)
}

func TestService_DeleteExternalLink_Success(t *testing.T) {
	mockRepo := &MockRepository{}
	deleted := false
	mockRepo.DeleteLinkFunc = func(ctx context.Context, linkID uuid.UUID) error {
		deleted = true
		return nil
	}
	svc := NewService(mockRepo)

	err := svc.DeleteExternalLink(context.Background(), uuid.New())
	assert.NoError(t, err)
	assert.True(t, deleted)
}

// helper
func ptrString(s string) *string {
	return &s
}
