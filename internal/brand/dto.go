package brand

import "github.com/google/uuid"

type CreateBrandRequest struct {
	Name    string `json:"name" binding:"required"`
	LogoURL string `json:"logo_url"`
}

type UpdateBrandRequest struct {
	Name    *string `json:"name,omitempty"`
	LogoURL *string `json:"logo_url,omitempty"`
}

type BrandResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	LogoURL   string    `json:"logo_url"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}
