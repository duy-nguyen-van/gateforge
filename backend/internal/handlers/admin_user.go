package handlers

import (
	"net/http"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"

	"github.com/labstack/echo/v4"
)

func actorUserIDFromContext(c echo.Context) (string, error) {
	v := c.Get(auth.EchoContextUserIDKey)
	if v == nil {
		return "", errors.UnauthorizedError("Authentication required", nil)
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return "", errors.UnauthorizedError("Authentication required", nil)
	}
	return s, nil
}

func loginHistoryFiltersFromQuery(c echo.Context) (dtos.AdminLoginHistoryListParams, error) {
	filters := dtos.AdminLoginHistoryListParams{
		TenantID: c.QueryParam("tenant_id"),
		ActorID:  c.QueryParam("actor_id"),
		Result:   c.QueryParam("result"),
	}
	if from := c.QueryParam("from"); from != "" {
		t, err := time.Parse(time.RFC3339, from)
		if err != nil {
			return filters, errors.ValidationError("Invalid from timestamp (use RFC3339)", err)
		}
		filters.From = &t
	}
	if to := c.QueryParam("to"); to != "" {
		t, err := time.Parse(time.RFC3339, to)
		if err != nil {
			return filters, errors.ValidationError("Invalid to timestamp (use RFC3339)", err)
		}
		filters.To = &t
	}
	return filters, nil
}

// GetUser godoc
// @Summary Get user details (platform admin)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User UUID"
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.AdminUserDetailResponse}
// @Router /admin/users/{userId} [get]
func (h *AdminHandler) GetUser(c echo.Context) error {
	userID := c.Param("userId")
	if userID == "" {
		return h.HandleError(c, errors.ValidationError("userId is required", nil))
	}
	user, err := h.adminService.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "User retrieved successfully", user, nil)
}

// DisableUser godoc
// @Summary Disable a user account (platform admin)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User UUID"
// @Success 204
// @Router /admin/users/{userId}/disable [post]
func (h *AdminHandler) DisableUser(c echo.Context) error {
	actorID, err := actorUserIDFromContext(c)
	if err != nil {
		return h.HandleError(c, err)
	}
	targetID := c.Param("userId")
	if targetID == "" {
		return h.HandleError(c, errors.ValidationError("userId is required", nil))
	}
	if err := h.adminService.DisableUser(c.Request().Context(), actorID, targetID); err != nil {
		return h.HandleError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// ForceLogoutUser godoc
// @Summary Force logout a user (platform admin)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User UUID"
// @Success 204
// @Router /admin/users/{userId}/force-logout [post]
func (h *AdminHandler) ForceLogoutUser(c echo.Context) error {
	actorID, err := actorUserIDFromContext(c)
	if err != nil {
		return h.HandleError(c, err)
	}
	targetID := c.Param("userId")
	if targetID == "" {
		return h.HandleError(c, errors.ValidationError("userId is required", nil))
	}
	if err := h.adminService.ForceLogoutUser(c.Request().Context(), actorID, targetID); err != nil {
		return h.HandleError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// ResetUserMFA godoc
// @Summary Reset MFA (TOTP + recovery codes) for a user (platform admin)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User UUID"
// @Success 204
// @Router /admin/users/{userId}/reset-mfa [post]
func (h *AdminHandler) ResetUserMFA(c echo.Context) error {
	actorID, err := actorUserIDFromContext(c)
	if err != nil {
		return h.HandleError(c, err)
	}
	targetID := c.Param("userId")
	if targetID == "" {
		return h.HandleError(c, errors.ValidationError("userId is required", nil))
	}
	if err := h.adminService.ResetMFA(c.Request().Context(), actorID, targetID); err != nil {
		return h.HandleError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// ResetUserPasskeys godoc
// @Summary Reset all passkeys for a user (platform admin)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param userId path string true "User UUID"
// @Success 204
// @Router /admin/users/{userId}/reset-passkey [post]
func (h *AdminHandler) ResetUserPasskeys(c echo.Context) error {
	actorID, err := actorUserIDFromContext(c)
	if err != nil {
		return h.HandleError(c, err)
	}
	targetID := c.Param("userId")
	if targetID == "" {
		return h.HandleError(c, errors.ValidationError("userId is required", nil))
	}
	if err := h.adminService.ResetPasskeys(c.Request().Context(), actorID, targetID); err != nil {
		return h.HandleError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// GetClientUsage godoc
// @Summary Get OAuth client usage metrics (platform admin)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param clientId path string true "Client record UUID"
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.AdminClientUsageResponse}
// @Router /admin/clients/{clientId}/usage [get]
func (h *AdminHandler) GetClientUsage(c echo.Context) error {
	clientID := c.Param("clientId")
	if clientID == "" {
		return h.HandleError(c, errors.ValidationError("clientId is required", nil))
	}
	usage, err := h.adminService.GetClientUsage(c.Request().Context(), clientID)
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "Client usage retrieved successfully", usage, nil)
}

// ListLoginHistory godoc
// @Summary List login history events (platform admin)
// @Tags Admin
// @Produce json
// @Security BearerAuth
// @Param tenant_id query string false "Filter by tenant UUID"
// @Param result query string false "Filter by result (success, failure, denied)"
// @Param actor_id query string false "Filter by actor id"
// @Param from query string false "ISO8601 lower bound on created_at"
// @Param to query string false "ISO8601 upper bound on created_at"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} object{meta=dtos.Meta,data=[]dtos.AdminAuditLogResponse}
// @Router /admin/login-history [get]
func (h *AdminHandler) ListLoginHistory(c echo.Context) error {
	filters, err := loginHistoryFiltersFromQuery(c)
	if err != nil {
		return h.HandleError(c, err)
	}
	rows, pageable, err := h.adminService.ListLoginHistory(c.Request().Context(), filters, pageableFromQuery(c))
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "Login history retrieved", rows, pageable)
}
