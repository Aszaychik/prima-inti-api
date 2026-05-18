package product

import (
	"github.com/google/uuid"
)

type ProductFilters struct {
	BrandID      *uuid.UUID
	BrandName    *string
	CategoryID   *uuid.UUID
	CategoryName *string
	SeriesID     *uuid.UUID
	SeriesName   *string
	Model        *string
}
