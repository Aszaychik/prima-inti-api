package company

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	apperrors "github.com/aszaychik/prima-inti-api/internal/errors"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetCompanyProfile godoc
// @Summary Get company profile
// @Description Get the company profile (public)
// @Tags company
// @Accept json
// @Produce json
// @Success 200 {object} errors.Response{success=bool,data=CompanyResponse} "Company profile retrieved"
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Company profile not found"
// @Failure 500 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Internal server error"
// @Router /api/v1/company-profile [get]
func (h *Handler) GetCompanyProfile(c *gin.Context) {
	profile, err := h.service.GetCompanyProfile(c.Request.Context())
	if err != nil {
		apiErr := apperrors.InternalServerError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	if profile == nil {
		apiErr := apperrors.NotFound("company profile not found")
		c.JSON(apiErr.Status, apiErr)
		return
	}
	c.JSON(http.StatusOK, toCompanyResponse(profile))
}

// CreateCompanyProfile godoc
// @Summary Create company profile (admin only)
// @Description Create the company profile (only admin)
// @Tags company
// @Accept json
// @Produce json
// @Param body body CreateCompanyRequest true "Company data"
// @Security BearerAuth
// @Success 201 {object} errors.Response{success=bool,data=CompanyResponse} "Company profile created"
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Invalid input"
// @Failure 409 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Company profile already exists"
// @Failure 500 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Internal server error"
// @Router /api/v1/company-profile [post]
func (h *Handler) CreateCompanyProfile(c *gin.Context) {
	var req CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := apperrors.FromGinValidation(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	profile, err := h.service.CreateCompanyProfile(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case ErrCompanyAlreadyExists:
			apiErr := apperrors.Conflict(err.Error())
			c.JSON(apiErr.Status, apiErr)
		default:
			apiErr := apperrors.InternalServerError(err)
			c.JSON(apiErr.Status, apiErr)
		}
		return
	}
	c.JSON(http.StatusCreated, toCompanyResponse(profile))
}

// UpdateCompanyProfile godoc
// @Summary Update company profile (admin only)
// @Description Update company profile fields (partial update allowed)
// @Tags company
// @Accept json
// @Produce json
// @Param body body UpdateCompanyRequest true "Update fields"
// @Security BearerAuth
// @Success 200 {object} errors.Response{success=bool,data=CompanyResponse} "Company profile updated"
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Invalid input"
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Company profile not found"
// @Failure 500 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Internal server error"
// @Router /api/v1/company-profile [put]
func (h *Handler) UpdateCompanyProfile(c *gin.Context) {
	var req UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := apperrors.FromGinValidation(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	profile, err := h.service.UpdateCompanyProfile(c.Request.Context(), &req)
	if err != nil {
		switch err {
		case ErrCompanyNotFound:
			apiErr := apperrors.NotFound(err.Error())
			c.JSON(apiErr.Status, apiErr)
		default:
			apiErr := apperrors.InternalServerError(err)
			c.JSON(apiErr.Status, apiErr)
		}
		return
	}
	c.JSON(http.StatusOK, toCompanyResponse(profile))
}

// AddExternalLink godoc
// @Summary Add external link to company profile (admin only)
// @Description Add a new external link (e.g., social media, website)
// @Tags company
// @Accept json
// @Produce json
// @Param body body CreateLinkRequest true "Link data"
// @Security BearerAuth
// @Success 201 {object} errors.Response{success=bool,data=ExternalLinkResp} "External link created"
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Invalid input"
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Company profile not found"
// @Failure 500 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Internal server error"
// @Router /api/v1/company-profile/links [post]
func (h *Handler) AddExternalLink(c *gin.Context) {
	var req CreateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := apperrors.FromGinValidation(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	profile, err := h.service.GetCompanyProfile(c.Request.Context())
	if err != nil || profile == nil {
		apiErr := apperrors.NotFound("company profile not found")
		c.JSON(apiErr.Status, apiErr)
		return
	}

	link, err := h.service.AddExternalLink(c.Request.Context(), profile.ID, &req)
	if err != nil {
		apiErr := apperrors.InternalServerError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}
	c.JSON(http.StatusCreated, toLinkResponse(link))
}

// UpdateExternalLink godoc
// @Summary Update an external link (admin only)
// @Description Update an existing external link by its ID
// @Tags company
// @Accept json
// @Produce json
// @Param linkId path string true "Link ID (UUID)"
// @Param body body UpdateLinkRequest true "Update fields"
// @Security BearerAuth
// @Success 200 {object} errors.Response{success=bool,data=ExternalLinkResp} "External link updated"
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Invalid ID or input"
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo} "External link not found"
// @Failure 500 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Internal server error"
// @Router /api/v1/company-profile/links/{linkId} [put]
func (h *Handler) UpdateExternalLink(c *gin.Context) {
	linkID, err := uuid.Parse(c.Param("linkId"))
	if err != nil {
		apiErr := apperrors.BadRequest("invalid link id")
		c.JSON(apiErr.Status, apiErr)
		return
	}

	var req UpdateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := apperrors.FromGinValidation(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	link, err := h.service.UpdateExternalLink(c.Request.Context(), linkID, &req)
	if err != nil {
		switch err {
		case ErrLinkNotFound:
			apiErr := apperrors.NotFound(err.Error())
			c.JSON(apiErr.Status, apiErr)
		default:
			apiErr := apperrors.InternalServerError(err)
			c.JSON(apiErr.Status, apiErr)
		}
		return
	}
	c.JSON(http.StatusOK, toLinkResponse(link))
}

// DeleteExternalLink godoc
// @Summary Delete an external link (admin only)
// @Description Delete an external link by its ID
// @Tags company
// @Produce json
// @Param linkId path string true "Link ID (UUID)"
// @Security BearerAuth
// @Success 204 "No Content"
// @Failure 400 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Invalid link ID"
// @Failure 404 {object} errors.Response{success=bool,error=errors.ErrorInfo} "External link not found"
// @Failure 500 {object} errors.Response{success=bool,error=errors.ErrorInfo} "Internal server error"
// @Router /api/v1/company-profile/links/{linkId} [delete]
func (h *Handler) DeleteExternalLink(c *gin.Context) {
	linkID, err := uuid.Parse(c.Param("linkId"))
	if err != nil {
		apiErr := apperrors.BadRequest("invalid link id")
		c.JSON(apiErr.Status, apiErr)
		return
	}

	err = h.service.DeleteExternalLink(c.Request.Context(), linkID)
	if err != nil {
		switch err {
		case ErrLinkNotFound:
			apiErr := apperrors.NotFound(err.Error())
			c.JSON(apiErr.Status, apiErr)
		default:
			apiErr := apperrors.InternalServerError(err)
			c.JSON(apiErr.Status, apiErr)
		}
		return
	}
	c.Status(http.StatusNoContent)
}

// ----- helpers (same as before) -----
func toCompanyResponse(p *CompanyProfile) CompanyResponse {
	resp := CompanyResponse{
		ID:        p.ID,
		Name:      p.Name,
		Phone:     p.Phone,
		Email:     p.Email,
		Address:   p.Address,
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
	}
	if len(p.ExternalLinks) > 0 {
		links := make([]ExternalLinkResp, len(p.ExternalLinks))
		for i, l := range p.ExternalLinks {
			links[i] = toLinkResponse(&l)
		}
		resp.ExternalLinks = links
	}
	return resp
}

func toLinkResponse(l *ExternalLink) ExternalLinkResp {
	return ExternalLinkResp{
		ID:       l.ID,
		Platform: l.Platform,
		URL:      l.URL,
	}
}
