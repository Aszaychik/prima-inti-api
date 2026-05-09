package company

import (
	"time"

	"github.com/google/uuid"
)

type CompanyProfile struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Phone     string    `gorm:"not null" json:"phone"`
	Email     string    `gorm:"not null" json:"email"`
	Address   string    `gorm:"not null" json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	ExternalLinks []ExternalLink `gorm:"foreignKey:CompanyID" json:"external_links,omitempty"`
}

func (CompanyProfile) TableName() string {
	return "company_profiles"
}

type ExternalLink struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CompanyID uuid.UUID `gorm:"type:uuid;not null" json:"company_id"`
	Platform  string    `gorm:"not null" json:"platform"`
	URL       string    `gorm:"not null" json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ExternalLink) TableName() string {
	return "external_links"
}
