package handlers

import (
	"crypto/subtle"
	"net/http"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/services"

	"github.com/labstack/echo/v4"
)

// TenantIdentityAdminHandler exposes operator APIs for tenant IdP settings (no end-user JWT yet).
type TenantIdentityAdminHandler struct {
	BaseHandler
	cfg          *config.Config
	adminService services.AdminService
}

// ProvideTenantIdentityAdminHandler wires admin routes for tenant identity providers.
func ProvideTenantIdentityAdminHandler(
	cfg *config.Config,
	adminService services.AdminService,
) *TenantIdentityAdminHandler {
	return &TenantIdentityAdminHandler{
		BaseHandler:  *NewBaseHandler(),
		cfg:          cfg,
		adminService: adminService,
	}
}

// PatchIdentityProvider configures an upstream identity provider for a tenant (X-Admin-API-Key).
// @Summary Configure an upstream identity provider for a tenant
// @Tags internal
// @Param tenantId path string true "Tenant UUID"
// @Param provider path string true "Provider id (e.g. google)"
// @Param X-Admin-API-Key header string true "Admin API key"
// @Param body body dtos.PatchIdentityProviderRequest true "OAuth credentials and enabled flag"
// @Success 204
// @Router /internal/tenants/{tenantId}/identity-providers/{provider} [patch]
func (h *TenantIdentityAdminHandler) PatchIdentityProvider(c echo.Context) error {
	if h.cfg.AdminAPIKey == "" {
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	}
	sent := []byte(c.Request().Header.Get("X-Admin-API-Key"))
	expected := []byte(h.cfg.AdminAPIKey)
	if len(sent) == 0 || subtle.ConstantTimeCompare(sent, expected) != 1 {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid admin key")
	}

	tenantID := c.Param("tenantId")
	if tenantID == "" {
		return h.HandleError(c, errors.ValidationError("tenantId is required", nil))
	}
	providerID := c.Param("provider")
	if providerID == "" {
		return h.HandleError(c, errors.ValidationError("provider is required", nil))
	}

	var body dtos.PatchIdentityProviderRequest
	if err := c.Bind(&body); err != nil {
		return h.HandleError(c, errors.ValidationError("Invalid request body", err))
	}

	if err := h.adminService.ConfigureIdentityProvider(
		c.Request().Context(),
		tenantID,
		providerID,
		&body,
		constants.AuditActorTypeAdminAPIKey,
	); err != nil {
		return h.HandleError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}
