package main

import (
	"log/slog"
	"os"

	"github.com/aszaychik/prima-inti-api/internal/config"
	"github.com/aszaychik/prima-inti-api/internal/db"
	"github.com/aszaychik/prima-inti-api/internal/seed"
)

func main() {
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
	defer sqlDB.Close()

	// Run all seeders in order (files are sorted by prefix)
	if err := seed.Run(database); err != nil {
		slog.Error("Seeding failed", "err", err)
		os.Exit(1)
	}

	slog.Info("Seeding completed successfully")
}
