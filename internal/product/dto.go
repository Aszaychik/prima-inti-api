package product

import "github.com/google/uuid"

type CreateProductRequest struct {
	BrandID     uuid.UUID  `json:"brand_id" binding:"required"`
	CategoryID  uuid.UUID  `json:"category_id" binding:"required"`
	SeriesID    *uuid.UUID `json:"series_id,omitempty"`
	Model       string     `json:"model" binding:"required"`
	Description string     `json:"description"`
	Price       *float64   `json:"price,omitempty"`
	Stock       int        `json:"stock"`
	ImageURL    string     `json:"image_url"`
}

type UpdateProductRequest struct {
	BrandID     *uuid.UUID `json:"brand_id,omitempty"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	SeriesID    *uuid.UUID `json:"series_id,omitempty"`
	Model       *string    `json:"model,omitempty"`
	Description *string    `json:"description,omitempty"`
	Price       *float64   `json:"price,omitempty"`
	Stock       *int       `json:"stock,omitempty"`
	ImageURL    *string    `json:"image_url,omitempty"`
}

type ProductResponse struct {
	ID          string       `json:"id"`
	Brand       BrandInfo    `json:"brand"`
	Category    CategoryInfo `json:"category"`
	Series      *SeriesInfo  `json:"series,omitempty"`
	Model       string       `json:"model"`
	Description string       `json:"description"`
	Price       *float64     `json:"price,omitempty"`
	Stock       int          `json:"stock"`
	ImageURL    string       `json:"image_url"`
	CreatedAt   string       `json:"created_at"`
	UpdatedAt   string       `json:"updated_at"`
}

type BrandInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CategoryInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SeriesInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
