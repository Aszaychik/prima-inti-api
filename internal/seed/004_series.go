package seed

import (
	"github.com/aszaychik/prima-inti-api/internal/series"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func init() {
	Register("004_series", seedSeries)
}

func seedSeries(db *gorm.DB) error {
	brandID := uuid.MustParse("58c78c4e-4350-4513-aadb-791a34472f5f")
	seriesList := []series.Series{
		{
			ID:      uuid.MustParse("0280b56c-4d55-44ad-972c-c16f8cf2b38c"),
			BrandID: brandID,
			Name:    "FD5 Series",
		},
	}
	for _, s := range seriesList {
		var existing series.Series
		err := db.Where("brand_id = ? AND name = ?", s.BrandID, s.Name).First(&existing).Error
		if err == nil {
			continue
		}
		if err := db.Create(&s).Error; err != nil {
			return err
		}
	}
	return nil
}
