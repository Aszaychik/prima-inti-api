package series

import (
	"time"

	"github.com/google/uuid"
)

type Series struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	BrandID   uuid.UUID `gorm:"type:uuid;not null" json:"brand_id"`
	Name      string    `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
