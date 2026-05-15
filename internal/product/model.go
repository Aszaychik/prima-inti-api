package product

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BrandID     uuid.UUID  `gorm:"type:uuid;not null" json:"brand_id"`
	CategoryID  uuid.UUID  `gorm:"type:uuid;not null" json:"category_id"`
	SeriesID    *uuid.UUID `gorm:"type:uuid" json:"series_id,omitempty"`
	Model       string     `gorm:"not null" json:"model"`
	Description string     `json:"description"`
	Price       *float64   `gorm:"type:decimal(10,2)" json:"price,omitempty"`
	Stock       int        `gorm:"not null;default:0" json:"stock"`
	ImageURL    string     `json:"image_url"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relationships (for preloading)
	Brand    *Brand    `gorm:"foreignKey:BrandID"`
	Category *Category `gorm:"foreignKey:CategoryID"`
	Series   *Series   `gorm:"foreignKey:SeriesID"`
}

// Define minimal structs for preloading (or import from their packages if no cycle)
type Brand struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Category struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Series struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
