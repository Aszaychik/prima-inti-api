package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/aszaychik/prima-inti-api/internal/company"
	"github.com/aszaychik/prima-inti-api/internal/config"
	"github.com/aszaychik/prima-inti-api/internal/db"
)

func main() {
	flag.Parse()

	slog.Info("Starting seeder...")

	cfg, err := config.LoadConfig("")
	if err != nil {
		slog.Error("Failed to load configuration", "err", err)
		os.Exit(1)
	}

	database, err := db.NewPostgresDBFromDatabaseConfig(cfg.Database)
	if err != nil {
		slog.Error("Failed to connect to database", "err", err)
		os.Exit(1)
	}

	sqlDB, err := database.DB()
	if err != nil {
		slog.Error("Failed to get database instance", "err", err)
		os.Exit(1)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			slog.Warn("Failed to close database connection", "err", err)
		}
	}()

	if err := seedCompanyProfile(database); err != nil {
		slog.Error("Seeding failed", "err", err)
		os.Exit(1)
	}

	slog.Info("Seeding completed successfully")
}

func seedCompanyProfile(db *gorm.DB) error {
	slog.Info("Seeding company profile...")

	// Check if company profile already exists
	var existing company.CompanyProfile
	result := db.First(&existing)
	if result.Error == nil {
		slog.Info("Company profile already exists, skipping seed")
		return nil
	}
	if result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	// Create company profile
	companyID := uuid.New()
	comp := company.CompanyProfile{
		ID:      companyID,
		Name:    "CV PRIMA INTI VANINDO",
		Phone:   "+6289519750202",
		Email:   "aszaychik@gmail.com",
		Address: "Jl. Rajawali No.27A Ds, Punggul, Kec. Gedangan, Kabupaten Sidoarjo, Jawa Timur 61254",
	}
	if err := db.Create(&comp).Error; err != nil {
		return err
	}
	slog.Info("Company profile created", "id", companyID)

	// Create external links
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
		slog.Info("External link created", "platform", link.Platform, "url", link.URL)
	}

	return nil
}
