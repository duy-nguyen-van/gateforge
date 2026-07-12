package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/services"

	"github.com/labstack/echo/v4"
)

type stubTenantIdpAdminService struct {
	services.AdminService
	lastTenantID   string
	lastProviderID string
	lastReq        *dtos.PatchIdentityProviderRequest
	err            error
}

func (s *stubTenantIdpAdminService) ConfigureIdentityProvider(_ context.Context, tenantID, providerID string, req *dtos.PatchIdentityProviderRequest, _ constants.AuditActorType) error {
	s.lastTenantID = tenantID
	s.lastProviderID = providerID
	s.lastReq = req
	return s.err
}

func TestTenantIdentityAdminHandler_PatchIdentityProvider_Unauthorized(t *testing.T) {
	e := echo.New()
	cfg := &config.Config{AdminAPIKey: "secret-key"}
	svc := &stubTenantIdpAdminService{}
	h := ProvideTenantIdentityAdminHandler(cfg, svc)

	req := httptest.NewRequest(http.MethodPatch, "/internal/tenants/t1/identity-providers/google", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/internal/tenants/:tenantId/identity-providers/:provider")
	c.SetParamNames("tenantId", "provider")
	c.SetParamValues("t1", "google")

	err := h.PatchIdentityProvider(c)
	if err == nil {
		t.Fatal("expected error")
	}
	var he *echo.HTTPError
	if !errors.As(err, &he) || he.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %v", err)
	}
}

func TestTenantIdentityAdminHandler_PatchIdentityProvider_Success(t *testing.T) {
	e := echo.New()
	cfg := &config.Config{AdminAPIKey: "secret-key"}
	svc := &stubTenantIdpAdminService{}
	h := ProvideTenantIdentityAdminHandler(cfg, svc)

	body, _ := json.Marshal(map[string]any{
		"enabled":             true,
		"oauth_client_id":     "cid",
		"oauth_client_secret": "sec",
	})
	req := httptest.NewRequest(http.MethodPatch, "/internal/tenants/t1/identity-providers/google", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Admin-API-Key", "secret-key")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/internal/tenants/:tenantId/identity-providers/:provider")
	c.SetParamNames("tenantId", "provider")
	c.SetParamValues("t1", "google")

	if err := h.PatchIdentityProvider(c); err != nil {
		t.Fatalf("PatchIdentityProvider: %v", err)
	}
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d", rec.Code)
	}
	if svc.lastTenantID != "t1" || svc.lastProviderID != "google" {
		t.Fatalf("tenant=%q provider=%q", svc.lastTenantID, svc.lastProviderID)
	}
	if svc.lastReq == nil || svc.lastReq.Enabled == nil || !*svc.lastReq.Enabled || svc.lastReq.OAuthClientID != "cid" {
		t.Fatalf("unexpected req: %+v", svc.lastReq)
	}
}

func TestTenantIdentityAdminHandler_PatchIdentityProvider_DisabledWhenNoKey(t *testing.T) {
	e := echo.New()
	cfg := &config.Config{AdminAPIKey: ""}
	svc := &stubTenantIdpAdminService{}
	h := ProvideTenantIdentityAdminHandler(cfg, svc)

	req := httptest.NewRequest(http.MethodPatch, "/internal/tenants/t1/identity-providers/google", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/internal/tenants/:tenantId/identity-providers/:provider")
	c.SetParamNames("tenantId", "provider")
	c.SetParamValues("t1", "google")

	err := h.PatchIdentityProvider(c)
	if err == nil {
		t.Fatal("expected error")
	}
	var he *echo.HTTPError
	if !errors.As(err, &he) || he.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %v", err)
	}
}
