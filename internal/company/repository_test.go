package company

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&CompanyProfile{}, &ExternalLink{})
	require.NoError(t, err)

	return db
}

func TestRepository_CreateCompany(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	company := &CompanyProfile{
		ID:      uuid.New(),
		Name:    "Test Company",
		Phone:   "123456789",
		Email:   "test@example.com",
		Address: "123 Test St",
	}

	err := repo.CreateCompany(context.Background(), company)
	assert.NoError(t, err)

	var found CompanyProfile
	err = db.First(&found, "id = ?", company.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, company.Name, found.Name)
}

func TestRepository_GetCompany_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	company, err := repo.GetCompany(context.Background())
	assert.NoError(t, err)
	assert.Nil(t, company)
}

func TestRepository_UpdateCompany(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	company := &CompanyProfile{
		ID:      uuid.New(),
		Name:    "Original",
		Phone:   "111",
		Email:   "old@example.com",
		Address: "Old Address",
	}
	err := repo.CreateCompany(context.Background(), company)
	require.NoError(t, err)

	company.Name = "Updated"
	err = repo.UpdateCompany(context.Background(), company)
	assert.NoError(t, err)

	var updated CompanyProfile
	db.First(&updated, "id = ?", company.ID)
	assert.Equal(t, "Updated", updated.Name)
}

func TestRepository_CreateLink(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	company := &CompanyProfile{
		ID:      uuid.New(),
		Name:    "Link Test Co",
		Phone:   "123",
		Email:   "link@test.com",
		Address: "Link St",
	}
	err := repo.CreateCompany(context.Background(), company)
	require.NoError(t, err)

	link := &ExternalLink{
		ID:        uuid.New(),
		CompanyID: company.ID,
		Platform:  "twitter",
		URL:       "https://twitter.com/test",
	}
	err = repo.CreateLink(context.Background(), link)
	assert.NoError(t, err)

	var found ExternalLink
	db.First(&found, "id = ?", link.ID)
	assert.Equal(t, link.Platform, found.Platform)
}

func TestRepository_DeleteLink(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepository(db)

	company := &CompanyProfile{ID: uuid.New(), Name: "Del", Phone: "1", Email: "del@test.com", Address: "Del St"}
	repo.CreateCompany(context.Background(), company)

	link := &ExternalLink{ID: uuid.New(), CompanyID: company.ID, Platform: "fb", URL: "https://fb.com/test"}
	repo.CreateLink(context.Background(), link)

	err := repo.DeleteLink(context.Background(), link.ID)
	assert.NoError(t, err)

	var count int64
	db.Model(&ExternalLink{}).Where("id = ?", link.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}
