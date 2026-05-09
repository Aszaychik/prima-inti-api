package company

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	apperrors "github.com/aszaychik/prima-inti-api/internal/errors"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(apperrors.ErrorHandler()) // optional, but matches real setup
	return router
}

func TestHandler_GetCompanyProfile_Success(t *testing.T) {
	mockSvc := &MockService{}
	expectedCompany := &CompanyProfile{
		ID:      uuid.New(),
		Name:    "Test Co",
		Phone:   "123",
		Email:   "test@co.com",
		Address: "123 St",
	}
	mockSvc.GetCompanyProfileFunc = func(ctx context.Context) (*CompanyProfile, error) {
		return expectedCompany, nil
	}
	handler := NewHandler(mockSvc)

	router := setupTestRouter()
	router.GET("/company-profile", handler.GetCompanyProfile)

	req := httptest.NewRequest(http.MethodGet, "/company-profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp CompanyResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	assert.Equal(t, expectedCompany.Name, resp.Name)
}

func TestHandler_GetCompanyProfile_NotFound(t *testing.T) {
	mockSvc := &MockService{}
	mockSvc.GetCompanyProfileFunc = func(ctx context.Context) (*CompanyProfile, error) {
		return nil, nil
	}
	handler := NewHandler(mockSvc)

	router := setupTestRouter()
	router.GET("/company-profile", handler.GetCompanyProfile)

	req := httptest.NewRequest(http.MethodGet, "/company-profile", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	var errResp apperrors.APIError
	json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.Equal(t, "company profile not found", errResp.Message)
}

func TestHandler_CreateCompanyProfile_Success(t *testing.T) {
	mockSvc := &MockService{}
	mockSvc.CreateCompanyProfileFunc = func(ctx context.Context, req *CreateCompanyRequest) (*CompanyProfile, error) {
		return &CompanyProfile{
			ID:      uuid.New(),
			Name:    req.Name,
			Phone:   req.Phone,
			Email:   req.Email,
			Address: req.Address,
		}, nil
	}
	handler := NewHandler(mockSvc)

	router := setupTestRouter()
	router.POST("/company-profile", handler.CreateCompanyProfile)

	body := CreateCompanyRequest{
		Name:    "New Co",
		Phone:   "555",
		Email:   "new@co.com",
		Address: "Addr",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/company-profile", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp CompanyResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, body.Name, resp.Name)
}

func TestHandler_CreateCompanyProfile_ValidationError(t *testing.T) {
	mockSvc := &MockService{}
	handler := NewHandler(mockSvc)

	router := setupTestRouter()
	router.POST("/company-profile", handler.CreateCompanyProfile)

	// Missing required fields
	body := map[string]string{"name": "only name"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/company-profile", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var errResp apperrors.APIError
	json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.Equal(t, "Validation failed", errResp.Message)
}

func TestHandler_AddExternalLink_Success(t *testing.T) {
	mockSvc := &MockService{}
	companyID := uuid.New()
	mockSvc.GetCompanyProfileFunc = func(ctx context.Context) (*CompanyProfile, error) {
		return &CompanyProfile{ID: companyID}, nil
	}
	mockSvc.AddExternalLinkFunc = func(ctx context.Context, cID uuid.UUID, req *CreateLinkRequest) (*ExternalLink, error) {
		return &ExternalLink{
			ID:        uuid.New(),
			CompanyID: cID,
			Platform:  req.Platform,
			URL:       req.URL,
		}, nil
	}
	handler := NewHandler(mockSvc)

	router := setupTestRouter()
	router.POST("/links", handler.AddExternalLink)

	reqBody := CreateLinkRequest{
		Platform: "linkedin",
		URL:      "https://linkedin.com/company/test",
	}
	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/links", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp ExternalLinkResp
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, reqBody.Platform, resp.Platform)
}

func TestHandler_DeleteExternalLink_Success(t *testing.T) {
	mockSvc := &MockService{}
	mockSvc.DeleteExternalLinkFunc = func(ctx context.Context, linkID uuid.UUID) error {
		return nil
	}
	handler := NewHandler(mockSvc)

	router := setupTestRouter()
	router.DELETE("/links/:linkId", handler.DeleteExternalLink)

	linkID := uuid.New()
	req := httptest.NewRequest(http.MethodDelete, "/links/"+linkID.String(), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}
