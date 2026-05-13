package company

import "github.com/google/uuid"

// CreateCompanyRequest
type CreateCompanyRequest struct {
	Name    string `json:"name" binding:"required"`
	LogoURL string `json:"logo_url"`
	Phone   string `json:"phone" binding:"required"`
	Email   string `json:"email" binding:"required,email"`
	Address string `json:"address" binding:"required"`
}

// UpdateCompanyRequest
type UpdateCompanyRequest struct {
	Name    *string `json:"name,omitempty"`
	LogoURL *string `json:"logo_url,omitempty"`
	Phone   *string `json:"phone,omitempty"`
	Email   *string `json:"email,omitempty"`
	Address *string `json:"address,omitempty"`
}

// CompanyResponse
type CompanyResponse struct {
	ID            uuid.UUID          `json:"id"`
	Name          string             `json:"name"`
	LogoURL       string             `json:"logo_url"`
	Phone         string             `json:"phone"`
	Email         string             `json:"email"`
	Address       string             `json:"address"`
	ExternalLinks []ExternalLinkResp `json:"external_links,omitempty"`
	CreatedAt     string             `json:"created_at"`
	UpdatedAt     string             `json:"updated_at"`
}

type ExternalLinkResp struct {
	ID       uuid.UUID `json:"id"`
	Platform string    `json:"platform"`
	URL      string    `json:"url"`
}

// CreateLinkRequest
type CreateLinkRequest struct {
	Platform string `json:"platform" binding:"required"`
	URL      string `json:"url" binding:"required,url"`
}

// UpdateLinkRequest
type UpdateLinkRequest struct {
	Platform *string `json:"platform,omitempty"`
	URL      *string `json:"url,omitempty"`
}
