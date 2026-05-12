package brand

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

// ListBrands godoc
// @Summary List all brands
// @Tags brands
// @Produce json
// @Success 200 {object} errors.Response{success=bool,data=[]BrandResponse}
// @Failure 500 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/brands [get]
func (h *Handler) ListBrands(c *gin.Context) {
	brands, err := h.service.List(c.Request.Context())
	if err != nil {
		apiErr := apperrors.InternalServerError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	resp := make([]BrandResponse, len(brands))
	for i, b := range brands {
		resp[i] = toBrandResponse(&b)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// GetBrand godoc
// @Summary Get brand by ID
// @Tags brands
// @Produce json
// @Param id path string true "Brand ID"
// @Success 200 {object} errors.Response{success=bool,data=BrandResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/brands/{id} [get]
func (h *Handler) GetBrand(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiErr := apperrors.BadRequest("invalid brand id")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	b, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		apiErr := apperrors.NotFound("brand not found")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": toBrandResponse(b)})
}

// CreateBrand godoc
// @Summary Create brand (admin only)
// @Tags brands
// @Accept json
// @Produce json
// @Param body body CreateBrandRequest true "Brand data"
// @Security BearerAuth
// @Success 201 {object} errors.Response{success=bool,data=BrandResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 409 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/brands [post]
func (h *Handler) CreateBrand(c *gin.Context) {
	var req CreateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := apperrors.FromGinValidation(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	b, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if err == ErrBrandAlreadyExists {
			apiErr := apperrors.Conflict(err.Error())
			c.JSON(apiErr.Status, apiErr)
			return
		}
		apiErr := apperrors.InternalServerError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": toBrandResponse(b)})
}

// UpdateBrand godoc
// @Summary Update brand (admin only)
// @Tags brands
// @Accept json
// @Produce json
// @Param id path string true "Brand ID"
// @Param body body UpdateBrandRequest true "Update fields"
// @Security BearerAuth
// @Success 200 {object} errors.Response{success=bool,data=BrandResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/brands/{id} [put]
func (h *Handler) UpdateBrand(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiErr := apperrors.BadRequest("invalid brand id")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	var req UpdateBrandRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := apperrors.FromGinValidation(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	b, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		if err == ErrBrandNotFound {
			apiErr := apperrors.NotFound(err.Error())
			c.JSON(apiErr.Status, apiErr)
			return
		}
		apiErr := apperrors.InternalServerError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": toBrandResponse(b)})
}

// DeleteBrand godoc
// @Summary Delete brand (admin only)
// @Tags brands
// @Produce json
// @Param id path string true "Brand ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/brands/{id} [delete]
func (h *Handler) DeleteBrand(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiErr := apperrors.BadRequest("invalid brand id")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		if err == ErrBrandNotFound {
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

func toBrandResponse(b *Brand) BrandResponse {
	return BrandResponse{
		ID:        b.ID,
		Name:      b.Name,
		LogoURL:   b.LogoURL,
		CreatedAt: b.CreatedAt.Format(time.RFC3339),
		UpdatedAt: b.UpdatedAt.Format(time.RFC3339),
	}
}
