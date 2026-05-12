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
}
