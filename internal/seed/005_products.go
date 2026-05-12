package seed

import (
	"fmt"

	"github.com/aszaychik/prima-inti-api/internal/product"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func init() {
	Register("005_products", seedProducts)
}

func seedProducts(db *gorm.DB) error {
	brandID := uuid.MustParse("58c78c4e-4350-4513-aadb-791a34472f5f")         // Kinco brand (already seeded in 003_brands.go)
	driveCategoryID := uuid.MustParse("d9adbfcb-0ebb-4c9c-937e-43dd18fccaac") // Servo Driver category (already seeded in 002_categories.go)
	motorCategoryID := uuid.MustParse("dd13305f-6b40-4b97-ac49-4d5bb35c472b") // Servo Motor category (already seeded in 002_categories.go)
	seriesID := uuid.MustParse("0280b56c-4d55-44ad-972c-c16f8cf2b38c")        // FD5 Series (already seeded in 004_series.go)

	imageURL := "https://en.kinco.cn/userfiles/images/2022/09/14/2022091413566676.png"

	// ----- 1. Servo Drives (FD5) -----
	driveModels := []struct {
		Model       string
		Power       string
		Description string
	}{
		// FD425
		{"FD425-PF-000", "200W", "Pulse interface"},
		{"FD425-EF-000", "200W", "EtherCAT interface"},
		{"FD425-CF-000", "200W", "CANopen interface"},
		{"FD425-LF-000", "200W", "Modbus RTU interface"},
		{"FD425-PA-000", "400W", "Pulse interface"},
		{"FD425-EA-000", "400W", "EtherCAT interface"},
		{"FD425-CA-000", "400W", "CANopen interface"},
		{"FD425-LA-000", "400W", "Modbus RTU interface"},
		// FD435
		{"FD435-PA-000", "400W", "Pulse interface, three-phase"},
		{"FD435-EA-000", "400W", "EtherCAT interface, three-phase"},
		{"FD435-CA-000", "400W", "CANopen interface, three-phase"},
		{"FD435-LA-000", "400W", "Modbus RTU interface, three-phase"},
		{"FD435-PF-000", "750W", "Pulse interface, three-phase"},
		{"FD435-EF-000", "750W", "EtherCAT interface, three-phase"},
		{"FD435-CF-000", "750W", "CANopen interface, three-phase"},
		{"FD435-LF-000", "750W", "Modbus RTU interface, three-phase"},
		// FD625
		{"FD625-A-000", "1kW", "Standard interface"},
	}

	for _, dm := range driveModels {
		if exists(db, dm.Model, brandID) {
			continue
		}
		prod := product.Product{
			ID:          uuid.New(),
			BrandID:     brandID,
			CategoryID:  driveCategoryID,
			SeriesID:    &seriesID,
			Model:       dm.Model,
			Description: fmt.Sprintf("Kinco FD5 series %s servo drive. %s.", dm.Power, dm.Description),
			Price:       nil,
			Stock:       100,
			ImageURL:    imageURL,
		}
		if err := db.Create(&prod).Error; err != nil {
			return err
		}
	}

	// ----- 2. Servo Motors (SMC series) -----
	motorModels := []struct {
		Model       string
		Power       string
		Torque      string
		Description string
	}{
		// 60 flange
		{"SMC60S-0020-30K-KLSU", "200W", "0.64 Nm", "60mm flange, low inertia"},
		{"SMC60S-0040-30K-KLSU", "400W", "1.27 Nm", "60mm flange, low inertia"},
		{"SMC60S-0075-30K-KLSU", "750W", "2.39 Nm", "60mm flange, low inertia"},
		// 80 flange
		{"SMC80S-0075-30K-KLSU", "750W", "2.39 Nm", "80mm flange, low inertia"},
		// 130 flange (medium inertia)
		{"SMC130D-0100-20MAK-SLSP", "1.0kW", "4.77 Nm", "130mm flange, medium inertia, 2000rpm"},
		{"SMC130D-0200-20MAK-SLSP", "2.0kW", "9.55 Nm", "130mm flange, medium inertia, 2000rpm"},
		{"SMC130D-0300-20MAK-SLSP", "3.0kW", "14.3 Nm", "130mm flange, medium inertia, 2000rpm"},
	}

	for _, mm := range motorModels {
		if exists(db, mm.Model, brandID) {
			continue
		}
		prod := product.Product{
			ID:          uuid.New(),
			BrandID:     brandID,
			CategoryID:  motorCategoryID,
			SeriesID:    &seriesID,
			Model:       mm.Model,
			Description: fmt.Sprintf("Kinco %s servo motor. Power: %s, Torque: %s. %s", mm.Model, mm.Power, mm.Torque, mm.Description),
			Price:       nil,
			Stock:       100,
			ImageURL:    imageURL,
		}
		if err := db.Create(&prod).Error; err != nil {
			return err
		}
	}

	return nil
}

// helper to avoid duplicate products
func exists(db *gorm.DB, model string, brandID uuid.UUID) bool {
	var count int64
	db.Model(&product.Product{}).Where("model = ? AND brand_id = ?", model, brandID).Count(&count)
	return count > 0
}
