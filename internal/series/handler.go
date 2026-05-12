package series

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

// ListSeries godoc
// @Summary List all series
// @Tags series
// @Produce json
// @Success 200 {object} errors.Response{success=bool,data=[]SeriesResponse}
// @Failure 500 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/series [get]
func (h *Handler) ListSeries(c *gin.Context) {
	series, err := h.service.List(c.Request.Context())
	if err != nil {
		apiErr := apperrors.InternalServerError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	resp := make([]SeriesResponse, len(series))
	for i, s := range series {
		resp[i] = toSeriesResponse(&s)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// GetSeries godoc
// @Summary Get series by ID
// @Tags series
// @Produce json
// @Param id path string true "Series ID"
// @Success 200 {object} errors.Response{success=bool,data=SeriesResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/series/{id} [get]
func (h *Handler) GetSeries(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiErr := apperrors.BadRequest("invalid series id")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	s, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		apiErr := apperrors.NotFound("series not found")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": toSeriesResponse(s)})
}

// ListSeriesByBrand godoc
// @Summary List series by brand ID
// @Tags series
// @Produce json
// @Param brandId path string true "Brand ID"
// @Success 200 {object} errors.Response{success=bool,data=[]SeriesResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/brands/{brandId}/series [get]
func (h *Handler) ListSeriesByBrand(c *gin.Context) {
	brandID, err := uuid.Parse(c.Param("brandId"))
	if err != nil {
		apiErr := apperrors.BadRequest("invalid brand id")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	series, err := h.service.ListByBrand(c.Request.Context(), brandID)
	if err != nil {
		apiErr := apperrors.InternalServerError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	resp := make([]SeriesResponse, len(series))
	for i, s := range series {
		resp[i] = toSeriesResponse(&s)
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// CreateSeries godoc
// @Summary Create series (admin only)
// @Tags series
// @Accept json
// @Produce json
// @Param body body CreateSeriesRequest true "Series data"
// @Security BearerAuth
// @Success 201 {object} errors.Response{success=bool,data=SeriesResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 409 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/series [post]
func (h *Handler) CreateSeries(c *gin.Context) {
	var req CreateSeriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := apperrors.FromGinValidation(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	s, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if err == ErrSeriesAlreadyExists {
			apiErr := apperrors.Conflict(err.Error())
			c.JSON(apiErr.Status, apiErr)
			return
		}
		apiErr := apperrors.InternalServerError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": toSeriesResponse(s)})
}

// UpdateSeries godoc
// @Summary Update series (admin only)
// @Tags series
// @Accept json
// @Produce json
// @Param id path string true "Series ID"
// @Param body body UpdateSeriesRequest true "Update fields"
// @Security BearerAuth
// @Success 200 {object} errors.Response{success=bool,data=SeriesResponse}
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/series/{id} [put]
func (h *Handler) UpdateSeries(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiErr := apperrors.BadRequest("invalid series id")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	var req UpdateSeriesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := apperrors.FromGinValidation(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	s, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		if err == ErrSeriesNotFound {
			apiErr := apperrors.NotFound(err.Error())
			c.JSON(apiErr.Status, apiErr)
			return
		}
		apiErr := apperrors.InternalServerError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": toSeriesResponse(s)})
}

// DeleteSeries godoc
// @Summary Delete series (admin only)
// @Tags series
// @Produce json
// @Param id path string true "Series ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo}
// @Router /api/v1/series/{id} [delete]
func (h *Handler) DeleteSeries(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		apiErr := apperrors.BadRequest("invalid series id")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		if err == ErrSeriesNotFound {
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

func toSeriesResponse(s *Series) SeriesResponse {
	return SeriesResponse{
		ID:        s.ID,
		BrandID:   s.BrandID,
		Name:      s.Name,
		CreatedAt: s.CreatedAt.Format(time.RFC3339),
		UpdatedAt: s.UpdatedAt.Format(time.RFC3339),
	}
}
