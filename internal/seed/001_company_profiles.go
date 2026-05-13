package seed

import (
	"log/slog"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/aszaychik/prima-inti-api/internal/company"
)

func init() {
	Register("001_company", seedCompany)
}

func seedCompany(db *gorm.DB) error {
	slog.Info("Seeding company profile")

	var count int64
	db.Model(&company.CompanyProfile{}).Count(&count)
	if count > 0 {
		slog.Info("Company profile already exists, skipping")
		return nil
	}

	companyID := uuid.New()
	comp := company.CompanyProfile{
		ID:      companyID,
		Name:    "CV PRIMA INTI VANINDO",
		LogoURL: "https://i.imgpeek.com/3I37ao3Cb5nx",
		Phone:   "+6289519750202",
		Email:   "aszaychik@gmail.com",
		Address: "Jl. Rajawali No.27A Ds, Punggul, Kec. Gedangan, Kabupaten Sidoarjo, Jawa Timur 61254",
	}
	if err := db.Create(&comp).Error; err != nil {
		return err
	}

	links := []company.ExternalLink{
		{
			ID:        uuid.New(),
			CompanyID: companyID,
			Platform:  "tokopedia",
			URL:       "https://www.tokopedia.com/prima-automation",
		},
		{
			ID:        uuid.New(),
			CompanyID: companyID,
			Platform:  "shopee",
			URL:       "https://shopee.co.id/primainti123",
		},
	}
	for _, link := range links {
		if err := db.Create(&link).Error; err != nil {
			return err
		}
	}
	return nil
}
