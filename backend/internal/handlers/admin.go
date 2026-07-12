package handlers

import (
	"net/http"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// AdminHandler serves JWT-protected platform admin APIs for the console SPA.
type AdminHandler struct {
	BaseHandler
	adminService services.AdminService
	validator    *validator.Validate
}

// ProvideAdminHandler wires admin console endpoints.
func ProvideAdminHandler(adminService services.AdminService, validator *validator.Validate) *AdminHandler {
	return &AdminHandler{
		BaseHandler:  *NewBaseHandler(),
		adminService: adminService,
		validator:    validator,
	}
}

// GetStats godoc
// @Summary Platform admin dashboard stats
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.AdminStatsResponse}
// @Router /admin/stats [get]
func (h *AdminHandler) GetStats(c echo.Context) error {
	stats, err := h.adminService.GetStats(c.Request().Context())
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "Admin stats retrieved successfully", stats, nil)
}

// ListUsers godoc
// @Summary List users (platform admin)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param tenant_id query string false "Filter by tenant UUID"
// @Param search query string false "Search email or name"
// @Success 200 {object} object{meta=dtos.Meta,data=[]dtos.AdminUserResponse}
// @Router /admin/users [get]
func (h *AdminHandler) ListUsers(c echo.Context) error {
	users, pageable, err := h.adminService.ListUsers(
		c.Request().Context(),
		c.QueryParam("tenant_id"),
		c.QueryParam("search"),
		pageableFromQuery(c),
	)
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "Users retrieved successfully", users, pageable)
}

// ListTenants godoc
// @Summary List tenants (platform admin)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} object{meta=dtos.Meta,data=[]dtos.AdminTenantResponse}
// @Router /admin/tenants [get]
func (h *AdminHandler) ListTenants(c echo.Context) error {
	tenants, pageable, err := h.adminService.ListTenants(c.Request().Context(), pageableFromQuery(c))
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "Tenants retrieved successfully", tenants, pageable)
}

// CreateTenant godoc
// @Summary Create a tenant (platform admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dtos.AdminCreateTenantRequest true "Tenant name and optional domain"
// @Success 201 {object} object{meta=dtos.Meta,data=dtos.AdminTenantResponse}
// @Router /admin/tenants [post]
func (h *AdminHandler) CreateTenant(c echo.Context) error {
	var body dtos.AdminCreateTenantRequest
	if err := c.Bind(&body); err != nil {
		return h.HandleError(c, errors.ValidationError("Invalid request body", err))
	}
	if err := h.validator.Struct(body); err != nil {
		return h.HandleError(c, errors.ValidationError("Validation failed", err))
	}
	tenant, err := h.adminService.CreateTenant(c.Request().Context(), &body)
	if err != nil {
		return h.HandleError(c, err)
	}
	res := &dtos.BaseResponse[*dtos.AdminTenantResponse]{}
	res.Meta = dtos.GetMeta(c, constants.Success, http.StatusCreated)
	res.Meta.Message = "Tenant created successfully"
	res.Data = tenant
	return res.JSON(c)
}

// GetTenant godoc
// @Summary Get a tenant by ID (platform admin)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param tenantId path string true "Tenant UUID"
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.AdminTenantResponse}
// @Router /admin/tenants/{tenantId} [get]
func (h *AdminHandler) GetTenant(c echo.Context) error {
	tenantID := c.Param("tenantId")
	if tenantID == "" {
		return h.HandleError(c, errors.ValidationError("tenantId is required", nil))
	}
	tenant, err := h.adminService.GetTenantByID(c.Request().Context(), tenantID)
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "Tenant retrieved successfully", tenant, nil)
}

// UpdateTenant godoc
// @Summary Update a tenant (platform admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tenantId path string true "Tenant UUID"
// @Param body body dtos.AdminUpdateTenantRequest true "Fields to update"
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.AdminTenantResponse}
// @Router /admin/tenants/{tenantId} [patch]
func (h *AdminHandler) UpdateTenant(c echo.Context) error {
	tenantID := c.Param("tenantId")
	if tenantID == "" {
		return h.HandleError(c, errors.ValidationError("tenantId is required", nil))
	}
	var body dtos.AdminUpdateTenantRequest
	if err := c.Bind(&body); err != nil {
		return h.HandleError(c, errors.ValidationError("Invalid request body", err))
	}
	if err := h.validator.Struct(body); err != nil {
		return h.HandleError(c, errors.ValidationError("Validation failed", err))
	}
	tenant, err := h.adminService.UpdateTenant(c.Request().Context(), tenantID, &body)
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "Tenant updated successfully", tenant, nil)
}

// DeleteTenant godoc
// @Summary Delete a tenant (platform admin)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param tenantId path string true "Tenant UUID"
// @Success 204
// @Router /admin/tenants/{tenantId} [delete]
func (h *AdminHandler) DeleteTenant(c echo.Context) error {
	tenantID := c.Param("tenantId")
	if tenantID == "" {
		return h.HandleError(c, errors.ValidationError("tenantId is required", nil))
	}
	if err := h.adminService.DeleteTenant(c.Request().Context(), tenantID); err != nil {
		return h.HandleError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// ListTenantMembers godoc
// @Summary List members of a tenant (platform admin)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param tenantId path string true "Tenant UUID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} object{meta=dtos.Meta,data=[]dtos.AdminTenantMemberResponse}
// @Router /admin/tenants/{tenantId}/members [get]
func (h *AdminHandler) ListTenantMembers(c echo.Context) error {
	tenantID := c.Param("tenantId")
	if tenantID == "" {
		return h.HandleError(c, errors.ValidationError("tenantId is required", nil))
	}
	members, pageable, err := h.adminService.ListTenantMembers(c.Request().Context(), tenantID, pageableFromQuery(c))
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "Tenant members retrieved successfully", members, pageable)
}

// ListClients godoc
// @Summary List OAuth clients (platform admin)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param tenant_id query string false "Filter by tenant UUID"
// @Success 200 {object} object{meta=dtos.Meta,data=[]dtos.AdminClientResponse}
// @Router /admin/clients [get]
func (h *AdminHandler) ListClients(c echo.Context) error {
	clients, pageable, err := h.adminService.ListClients(
		c.Request().Context(),
		c.QueryParam("tenant_id"),
		pageableFromQuery(c),
	)
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "Clients retrieved successfully", clients, pageable)
}

// CreateClient godoc
// @Summary Create an OAuth client (platform admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dtos.AdminCreateClientRequest true "OAuth client registration"
// @Success 201 {object} object{meta=dtos.Meta,data=dtos.AdminCreateClientResponse}
// @Router /admin/clients [post]
func (h *AdminHandler) CreateClient(c echo.Context) error {
	var body dtos.AdminCreateClientRequest
	if err := c.Bind(&body); err != nil {
		return h.HandleError(c, errors.ValidationError("Invalid request body", err))
	}
	if err := h.validator.Struct(body); err != nil {
		return h.HandleError(c, errors.ValidationError("Validation failed", err))
	}
	client, err := h.adminService.CreateClient(c.Request().Context(), &body)
	if err != nil {
		return h.HandleError(c, err)
	}
	res := &dtos.BaseResponse[*dtos.AdminCreateClientResponse]{}
	res.Meta = dtos.GetMeta(c, constants.Success, http.StatusCreated)
	res.Meta.Message = "OAuth client created successfully"
	res.Data = client
	return res.JSON(c)
}

// GetClient godoc
// @Summary Get an OAuth client by ID (platform admin)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param clientId path string true "Client record UUID"
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.AdminClientResponse}
// @Router /admin/clients/{clientId} [get]
func (h *AdminHandler) GetClient(c echo.Context) error {
	clientID := c.Param("clientId")
	if clientID == "" {
		return h.HandleError(c, errors.ValidationError("clientId is required", nil))
	}
	client, err := h.adminService.GetClientByID(c.Request().Context(), clientID)
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "OAuth client retrieved successfully", client, nil)
}

// UpdateClient godoc
// @Summary Update an OAuth client (platform admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param clientId path string true "Client record UUID"
// @Param body body dtos.AdminUpdateClientRequest true "Fields to update"
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.AdminClientResponse}
// @Router /admin/clients/{clientId} [patch]
func (h *AdminHandler) UpdateClient(c echo.Context) error {
	clientID := c.Param("clientId")
	if clientID == "" {
		return h.HandleError(c, errors.ValidationError("clientId is required", nil))
	}
	var body dtos.AdminUpdateClientRequest
	if err := c.Bind(&body); err != nil {
		return h.HandleError(c, errors.ValidationError("Invalid request body", err))
	}
	if err := h.validator.Struct(body); err != nil {
		return h.HandleError(c, errors.ValidationError("Validation failed", err))
	}
	client, err := h.adminService.UpdateClient(c.Request().Context(), clientID, &body)
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "OAuth client updated successfully", client, nil)
}

// DeleteClient godoc
// @Summary Delete an OAuth client (platform admin)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param clientId path string true "Client record UUID"
// @Success 204
// @Router /admin/clients/{clientId} [delete]
func (h *AdminHandler) DeleteClient(c echo.Context) error {
	clientID := c.Param("clientId")
	if clientID == "" {
		return h.HandleError(c, errors.ValidationError("clientId is required", nil))
	}
	if err := h.adminService.DeleteClient(c.Request().Context(), clientID); err != nil {
		return h.HandleError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// ListIdentityProviders godoc
// @Summary List identity providers for a tenant
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param tenantId path string true "Tenant UUID"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} object{meta=dtos.Meta,data=[]dtos.AdminIdentityProviderResponse}
// @Router /admin/tenants/{tenantId}/identity-providers [get]
func (h *AdminHandler) ListIdentityProviders(c echo.Context) error {
	tenantID := c.Param("tenantId")
	if tenantID == "" {
		return h.HandleError(c, errors.ValidationError("tenantId is required", nil))
	}
	providers, pageable, err := h.adminService.ListIdentityProviders(c.Request().Context(), tenantID, pageableFromQuery(c))
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "Identity providers retrieved successfully", providers, pageable)
}

// PatchIdentityProvider godoc
// @Summary Configure an upstream identity provider for a tenant (platform admin)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tenantId path string true "Tenant UUID"
// @Param provider path string true "Provider id (e.g. google)"
// @Param body body dtos.PatchIdentityProviderRequest true "OAuth credentials and enabled flag"
// @Success 204
// @Router /admin/tenants/{tenantId}/identity-providers/{provider} [patch]
func (h *AdminHandler) PatchIdentityProvider(c echo.Context) error {
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
		constants.AuditActorTypeUser,
	); err != nil {
		return h.HandleError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// AddMember godoc
// @Summary Add an existing user to a tenant
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tenantId path string true "Tenant UUID"
// @Param body body dtos.AdminAddMemberRequest true "User email"
// @Success 204
// @Router /admin/tenants/{tenantId}/members [post]
func (h *AdminHandler) AddMember(c echo.Context) error {
	tenantID := c.Param("tenantId")
	var body dtos.AdminAddMemberRequest
	if err := c.Bind(&body); err != nil {
		return h.HandleError(c, errors.ValidationError("Invalid request body", err))
	}
	if err := h.validator.Struct(body); err != nil {
		return h.HandleError(c, errors.ValidationError("Validation failed", err))
	}
	if err := h.adminService.AddMemberByEmail(c.Request().Context(), tenantID, body.Email, body.Role); err != nil {
		return h.HandleError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// RemoveMember godoc
// @Summary Remove a user from a tenant
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param tenantId path string true "Tenant UUID"
// @Param userId path string true "User UUID"
// @Success 204
// @Router /admin/tenants/{tenantId}/members/{userId} [delete]
func (h *AdminHandler) RemoveMember(c echo.Context) error {
	tenantID := c.Param("tenantId")
	userID := c.Param("userId")
	if err := h.adminService.RemoveMember(c.Request().Context(), tenantID, userID); err != nil {
		return h.HandleError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// ListAuditLogs godoc
// @Summary List platform audit logs
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param tenant_id query string false "Filter by tenant UUID"
// @Param action query string false "Filter by action (exact or prefix with trailing dot)"
// @Param result query string false "Filter by result (success, failure, denied)"
// @Param actor_id query string false "Filter by actor id"
// @Param from query string false "ISO8601 lower bound on created_at"
// @Param to query string false "ISO8601 upper bound on created_at"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} object{meta=dtos.Meta,data=[]dtos.AdminAuditLogResponse}
// @Router /admin/audit-logs [get]
func (h *AdminHandler) ListAuditLogs(c echo.Context) error {
	filters := dtos.AdminAuditLogListParams{
		TenantID: c.QueryParam("tenant_id"),
		Action:   c.QueryParam("action"),
		Result:   c.QueryParam("result"),
		ActorID:  c.QueryParam("actor_id"),
	}
	if from := c.QueryParam("from"); from != "" {
		t, err := time.Parse(time.RFC3339, from)
		if err != nil {
			return h.HandleError(c, errors.ValidationError("Invalid from timestamp (use RFC3339)", err))
		}
		filters.From = &t
	}
	if to := c.QueryParam("to"); to != "" {
		t, err := time.Parse(time.RFC3339, to)
		if err != nil {
			return h.HandleError(c, errors.ValidationError("Invalid to timestamp (use RFC3339)", err))
		}
		filters.To = &t
	}

	rows, pageable, err := h.adminService.ListAuditLogs(c.Request().Context(), filters, pageableFromQuery(c))
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "Audit logs retrieved", rows, pageable)
}
