package category

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

// ListCategories godoc
// @Summary List all categories
// @Tags categories
// @Produce json
// @Success 200 {object} errors.Response{success=bool,data=[]CategoryResponse}
// @Failure 500 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/categories [get]
func (h *Handler) ListCategories(c *gin.Context) {
	cats, err := h.service.List(c.Request.Context())
	if err != nil {
		apiErr := apperrors.InternalServerError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	resp := make([]CategoryResponse, len(cats))
	for i, cat := range cats {
		resp[i] = toCategoryResponse(&cat)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// GetCategory godoc
// @Summary Get category by ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} errors.Response{success=bool,data=CategoryResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/categories/{id} [get]
func (h *Handler) GetCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiErr := apperrors.BadRequest("invalid category id")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	cat, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		apiErr := apperrors.NotFound("category not found")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": toCategoryResponse(cat)})
}

// CreateCategory godoc
// @Summary Create category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Param body body CreateCategoryRequest true "Category data"
// @Security BearerAuth
// @Success 201 {object} errors.Response{success=bool,data=CategoryResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 409 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/categories [post]
func (h *Handler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := apperrors.FromGinValidation(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	cat, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if err == ErrCategoryAlreadyExists {
			apiErr := apperrors.Conflict(err.Error())
			c.JSON(apiErr.Status, apiErr)
			return
		}
		apiErr := apperrors.InternalServerError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": toCategoryResponse(cat)})
}

// UpdateCategory godoc
// @Summary Update category (admin only)
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param body body UpdateCategoryRequest true "Update fields"
// @Security BearerAuth
// @Success 200 {object} errors.Response{success=bool,data=CategoryResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/categories/{id} [put]
func (h *Handler) UpdateCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiErr := apperrors.BadRequest("invalid category id")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := apperrors.FromGinValidation(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	cat, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		if err == ErrCategoryNotFound {
			apiErr := apperrors.NotFound(err.Error())
			c.JSON(apiErr.Status, apiErr)
			return
		}
		apiErr := apperrors.InternalServerError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": toCategoryResponse(cat)})
}

// DeleteCategory godoc
// @Summary Delete category (admin only)
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/categories/{id} [delete]
func (h *Handler) DeleteCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiErr := apperrors.BadRequest("invalid category id")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		if err == ErrCategoryNotFound {
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

func toCategoryResponse(cat *Category) CategoryResponse {
	return CategoryResponse{
		ID:          cat.ID,
		Name:        cat.Name,
		Description: cat.Description,
		CreatedAt:   cat.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   cat.UpdatedAt.Format(time.RFC3339),
	}
}
