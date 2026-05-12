package seed

import (
	"github.com/aszaychik/prima-inti-api/internal/brand"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func init() {
	Register("003_brands", seedBrands)
}

func seedBrands(db *gorm.DB) error {
	brands := []brand.Brand{
		{
			ID:      uuid.MustParse("58c78c4e-4350-4513-aadb-791a34472f5f"),
			Name:    "Kinco",
			LogoURL: "https://cdn.brandfetch.io/idtzjAhCEM/w/145/h/46/theme/dark/logo.png?c=1dxbfHSJFAPEGdCLU4o5B",
		},
	}
	for _, b := range brands {
		var existing brand.Brand
		err := db.Where("name = ?", b.Name).First(&existing).Error
		if err == nil {
			continue
		}
		if err := db.Create(&b).Error; err != nil {
			return err
		}
	}
	return nil
}
