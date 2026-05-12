package seed

import (
	"github.com/aszaychik/prima-inti-api/internal/category"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func init() {
	Register("002_categories", seedCategories)
}

func seedCategories(db *gorm.DB) error {
	cats := []category.Category{
		{
			ID:          uuid.MustParse("d9adbfcb-0ebb-4c9c-937e-43dd18fccaac"),
			Name:        "Servo Driver",
			Description: "High-performance servo drives for precise motion control.",
		},
		{
			ID:          uuid.MustParse("dd13305f-6b40-4b97-ac49-4d5bb35c472b"),
			Name:        "Servo Motor",
			Description: "Compact and efficient servo motors for industrial automation.",
		},
	}
	for _, c := range cats {
		var existing category.Category
		err := db.Where("name = ?", c.Name).First(&existing).Error
		if err == nil {
			continue
		}
		if err := db.Create(&c).Error; err != nil {
			return err
		}
	}
	return nil
}
