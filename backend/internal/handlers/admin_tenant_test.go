package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/dtos"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type stubAdminTenantService struct {
	stubAdminAuditService
}

func (stubAdminTenantService) GetTenantByID(_ context.Context, tenantID string) (*dtos.AdminTenantResponse, error) {
	return &dtos.AdminTenantResponse{ID: tenantID, Name: "Acme"}, nil
}

func (stubAdminTenantService) CreateTenant(_ context.Context, req *dtos.AdminCreateTenantRequest) (*dtos.AdminTenantResponse, error) {
	return &dtos.AdminTenantResponse{Name: req.Name, Domain: req.Domain}, nil
}

func (stubAdminTenantService) UpdateTenant(_ context.Context, tenantID string, req *dtos.AdminUpdateTenantRequest) (*dtos.AdminTenantResponse, error) {
	resp := &dtos.AdminTenantResponse{ID: tenantID}
	if req.Name != nil {
		resp.Name = *req.Name
	}
	return resp, nil
}

func (stubAdminTenantService) DeleteTenant(context.Context, string) error { return nil }

func (stubAdminTenantService) ListTenantMembers(_ context.Context, tenantID string, _ *dtos.PageableRequest) ([]*dtos.AdminTenantMemberResponse, *dtos.Pageable, error) {
	return []*dtos.AdminTenantMemberResponse{
		{UserID: "u1", Email: "user@example.com", Role: "member", Status: "active"},
	}, &dtos.Pageable{Page: 1, PageSize: 20, Total: 1}, nil
}

func TestAdminHandler_CreateTenant(t *testing.T) {
	h := ProvideAdminHandler(stubAdminTenantService{}, validator.New())
	e := echo.New()
	body := `{"name":"Acme Corp","domain":"acme"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/tenants", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.CreateTenant(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.Contains(t, rec.Body.String(), "Acme Corp")
}

func TestAdminHandler_GetTenant(t *testing.T) {
	h := ProvideAdminHandler(stubAdminTenantService{}, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/tenants/tenant-1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("tenantId")
	c.SetParamValues("tenant-1")

	err := h.GetTenant(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "Acme")
}

func TestAdminHandler_ListTenantMembers(t *testing.T) {
	h := ProvideAdminHandler(stubAdminTenantService{}, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/tenants/tenant-1/members", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("tenantId")
	c.SetParamValues("tenant-1")

	err := h.ListTenantMembers(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "user@example.com")
}

func TestAdminHandler_DeleteTenant(t *testing.T) {
	h := ProvideAdminHandler(stubAdminTenantService{}, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/admin/tenants/tenant-1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("tenantId")
	c.SetParamValues("tenant-1")

	err := h.DeleteTenant(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestAdminHandler_UpdateTenant(t *testing.T) {
	h := ProvideAdminHandler(stubAdminTenantService{}, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/admin/tenants/tenant-1", strings.NewReader(`{"name":"Updated"}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("tenantId")
	c.SetParamValues("tenant-1")

	err := h.UpdateTenant(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "Updated")
}
