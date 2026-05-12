package series

import "github.com/google/uuid"

type CreateSeriesRequest struct {
	BrandID uuid.UUID `json:"brand_id" binding:"required"`
	Name    string    `json:"name" binding:"required"`
}

type UpdateSeriesRequest struct {
	Name *string `json:"name,omitempty"`
}

type SeriesResponse struct {
	ID        uuid.UUID `json:"id"`
	BrandID   uuid.UUID `json:"brand_id"`
	Name      string    `json:"name"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}
