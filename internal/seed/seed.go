package seed

import (
	"log/slog"

	"gorm.io/gorm"
)

type SeederFunc func(*gorm.DB) error

var seeders []struct {
	Name string
	Run  SeederFunc
}

func Register(name string, fn SeederFunc) {
	seeders = append(seeders, struct {
		Name string
		Run  SeederFunc
	}{Name: name, Run: fn})
}

func Run(db *gorm.DB) error {
	for _, s := range seeders {
		slog.Info("Running seeder", "name", s.Name)
		if err := s.Run(db); err != nil {
			slog.Error("Seeder failed", "name", s.Name, "err", err)
			return err
		}
	}
	return nil
}
