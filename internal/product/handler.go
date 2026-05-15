package product

import (
	"net/http"
	"time"

	apperrors "github.com/aszaychik/prima-inti-api/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// ListProducts godoc
// @Summary List all products
// @Tags products
// @Produce json
// @Success 200 {object} errors.Response{success=bool,data=[]ProductResponse}
// @Router /api/v1/products [get]
func (h *Handler) ListProducts(c *gin.Context) {
	products, err := h.service.List(c.Request.Context())
	if err != nil {
		apiErr := apperrors.InternalServerError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	resp := make([]ProductResponse, len(products))
	for i, p := range products {
		resp[i] = toProductResponse(&p)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// GetProduct godoc
// @Summary Get product by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} errors.Response{success=bool,data=ProductResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/products/{id} [get]
func (h *Handler) GetProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiErr := apperrors.BadRequest("invalid product id")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	p, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		apiErr := apperrors.NotFound("product not found")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": toProductResponse(p)})
}

// ListProductsByBrand godoc
// @Summary List products by brand ID
// @Tags products
// @Produce json
// @Param brandId path string true "Brand ID"
// @Success 200 {object} errors.Response{success=bool,data=[]ProductResponse}
// @Router /api/v1/brands/{brandId}/products [get]
func (h *Handler) ListProductsByBrand(c *gin.Context) {
	brandID, err := uuid.Parse(c.Param("brandId"))
	if err != nil {
		apiErr := apperrors.BadRequest("invalid brand id")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	products, err := h.service.ListByBrand(c.Request.Context(), brandID)
	if err != nil {
		apiErr := apperrors.InternalServerError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	resp := make([]ProductResponse, len(products))
	for i, p := range products {
		resp[i] = toProductResponse(&p)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// ListProductsByCategory godoc
// @Summary List products by category ID
// @Tags products
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} errors.Response{success=bool,data=[]ProductResponse}
// @Router /api/v1/categories/{id}/products [get]
func (h *Handler) ListProductsByCategory(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiErr := apperrors.BadRequest("invalid category id")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	products, err := h.service.ListByCategory(c.Request.Context(), categoryID)
	if err != nil {
		apiErr := apperrors.InternalServerError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	resp := make([]ProductResponse, len(products))
	for i, p := range products {
		resp[i] = toProductResponse(&p)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// ListProductsBySeries godoc
// @Summary List products by series ID
// @Tags products
// @Produce json
// @Param id path string true "Series ID"
// @Success 200 {object} errors.Response{success=bool,data=[]ProductResponse}
// @Router /api/v1/series/{id}/products [get]
func (h *Handler) ListProductsBySeries(c *gin.Context) {
	seriesID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiErr := apperrors.BadRequest("invalid series id")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	products, err := h.service.ListBySeries(c.Request.Context(), seriesID)
	if err != nil {
		apiErr := apperrors.InternalServerError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	resp := make([]ProductResponse, len(products))
	for i, p := range products {
		resp[i] = toProductResponse(&p)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// CreateProduct godoc
// @Summary Create product (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Param body body CreateProductRequest true "Product data"
// @Security BearerAuth
// @Success 201 {object} errors.Response{success=bool,data=ProductResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/products [post]
func (h *Handler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := apperrors.FromGinValidation(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	p, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case ErrBrandNotFound, ErrCategoryNotFound, ErrSeriesNotFound:
			apiErr := apperrors.NotFound(err.Error())
			c.JSON(apiErr.Status, apiErr)
		default:
			apiErr := apperrors.InternalServerError(err)
			c.JSON(apiErr.Status, apiErr)
		}
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": toProductResponse(p)})
}

// UpdateProduct godoc
// @Summary Update product (admin only)
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param body body UpdateProductRequest true "Update fields"
// @Security BearerAuth
// @Success 200 {object} errors.Response{success=bool,data=ProductResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/products/{id} [put]
func (h *Handler) UpdateProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiErr := apperrors.BadRequest("invalid product id")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := apperrors.FromGinValidation(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	p, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		switch err {
		case ErrProductNotFound:
			apiErr := apperrors.NotFound(err.Error())
			c.JSON(apiErr.Status, apiErr)
		case ErrBrandNotFound, ErrCategoryNotFound, ErrSeriesNotFound:
			apiErr := apperrors.NotFound(err.Error())
			c.JSON(apiErr.Status, apiErr)
		default:
			apiErr := apperrors.InternalServerError(err)
			c.JSON(apiErr.Status, apiErr)
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": toProductResponse(p)})
}

// DeleteProduct godoc
// @Summary Delete product (admin only)
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/products/{id} [delete]
func (h *Handler) DeleteProduct(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiErr := apperrors.BadRequest("invalid product id")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		if err == ErrProductNotFound {
			apiErr := apperrors.NotFound(err.Error())
			c.JSON(apiErr.Status, apiErr)
			return
		}
		apiErr := apperrors.InternalServerError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	c.Status(http.StatusNoContent)
}

func toProductResponse(p *Product) ProductResponse {
	resp := ProductResponse{
		ID: p.ID.String(),
		Brand: BrandInfo{
			ID:   p.BrandID.String(),
			Name: "",
		},
		Category: CategoryInfo{
			ID:   p.CategoryID.String(),
			Name: "",
		},
		Model:       p.Model,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		ImageURL:    p.ImageURL,
		CreatedAt:   p.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   p.UpdatedAt.Format(time.RFC3339),
	}

	if p.Brand != nil {
		resp.Brand.Name = p.Brand.Name
	}
	if p.Category != nil {
		resp.Category.Name = p.Category.Name
	}
	if p.Series != nil && p.SeriesID != nil {
		resp.Series = &SeriesInfo{
			ID:   p.Series.ID.String(),
			Name: p.Series.Name,
		}
	} else if p.SeriesID != nil {
		resp.Series = &SeriesInfo{
			ID:   p.SeriesID.String(),
			Name: "",
		}
	}
	return resp
}
