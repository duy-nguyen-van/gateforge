package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// AuthHandler handles registration, login, and current user.
type AuthHandler struct {
	BaseHandler
	userService       services.UserService
	sessionService    services.SessionService
	mfaService        services.MFAService
	federationService services.FederationService
	auditService      services.AuditService
	cfg               *config.Config
	validator         *validator.Validate
}

// ProvideAuthHandler wires Phase 1 identity endpoints and Phase 3 SSO session cookie on login.
func ProvideAuthHandler(
	userService services.UserService,
	sessionService services.SessionService,
	mfaService services.MFAService,
	federationService services.FederationService,
	auditService services.AuditService,
	cfg *config.Config,
	validator *validator.Validate,
) *AuthHandler {
	return &AuthHandler{
		BaseHandler:       *NewBaseHandler(),
		userService:       userService,
		sessionService:    sessionService,
		mfaService:        mfaService,
		federationService: federationService,
		auditService:      auditService,
		cfg:               cfg,
		validator:         validator,
	}
}

// Register godoc
// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dtos.RegisterRequest true "Credentials"
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.UserResponse}
// @Router /register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req dtos.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return h.HandleError(c, errors.ValidationError("Invalid request body", err))
	}
	if err := h.validator.Struct(req); err != nil {
		fieldErrors := errors.ParseValidationErrors(err)
		if len(fieldErrors) > 0 {
			return h.HandleError(c, errors.ValidationErrorWithDetails("Validation failed", err, fieldErrors))
		}
		return h.HandleError(c, errors.ValidationError("Validation failed", err))
	}

	user, err := h.userService.Register(c.Request().Context(), &req, c.Request().Host)
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "User registered successfully", dtos.NewUserResponse(user), nil)
}

// LoginAPI godoc
// @Summary Login (dashboard/API) with email and password
// @Description Sets the same iam_session cookie as OIDC so /authorize recognizes the browser without a second login. If the user has MFA (TOTP) enabled, `data` is MFALoginChallengeResponse (mfa_ticket) instead of LoginResponse; complete login with POST /mfa/challenge/verify. When the user belongs to multiple tenants and no tenant context is provided, `data` is TenantSelectionResponse.
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dtos.LoginRequest true "Credentials"
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.LoginResponse} "When MFA is off: LoginResponse. When MFA is on: same envelope but `data` is MFALoginChallengeResponse (see top-level description)."
// @Failure 400 {object} object{meta=dtos.Meta}
// @Failure 401 {object} object{meta=dtos.Meta}
// @Router /login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req dtos.LoginRequest
	if err := c.Bind(&req); err != nil {
		return h.HandleError(c, errors.ValidationError("Invalid request body", err))
	}
	if err := h.validator.Struct(req); err != nil {
		fieldErrors := errors.ParseValidationErrors(err)
		if len(fieldErrors) > 0 {
			return h.HandleError(c, errors.ValidationErrorWithDetails("Validation failed", err, fieldErrors))
		}
		return h.HandleError(c, errors.ValidationError("Validation failed", err))
	}

	u, err := h.userService.AuthenticateUser(c.Request().Context(), &req)
	if err != nil {
		return h.HandleError(c, err)
	}

	ctx := c.Request().Context()
	hasMFA, err := h.mfaService.HasActiveMFA(ctx, u.ID)
	if err != nil {
		return h.HandleError(c, err)
	}
	if hasMFA {
		tenantID := req.TenantID
		if tenantID == "" {
			tenantID = h.cfg.DefaultTenantID
		}
		ticket, exp, err := h.mfaService.CreateLoginTicket(ctx, auth.MFAPendingPayload{
			UserID:     u.ID,
			TenantID:   tenantID,
			RememberMe: req.RememberMe,
			ReturnTo:   req.ReturnTo,
		})
		if err != nil {
			return h.HandleError(c, err)
		}
		return h.SuccessResponse(c, "MFA required", dtos.MFALoginChallengeResponse{
			MfaRequired: true,
			MfaTicket:   ticket,
			ExpiresIn:   exp,
		}, nil)
	}

	loginOut, selectionOut, err := h.userService.CompleteAuth(ctx, u, services.TenantResolveInput{
		Host:          c.Request().Host,
		TenantIDParam: req.TenantID,
		UserID:        u.ID,
	})
	if err != nil {
		return h.HandleError(c, err)
	}
	if selectionOut != nil {
		return h.SuccessResponse(c, "Tenant selection required", selectionOut, nil)
	}

	sid, ttl, err := h.sessionService.Create(ctx, u.ID, loginOut.ActiveTenantID, c.RealIP(), c.Request().UserAgent(), req.RememberMe)
	if err != nil {
		return h.HandleError(c, err)
	}
	h.setSessionCookie(c, sid, ttl)

	return h.SuccessResponse(c, "Login successful", loginOut, nil)
}

// ExchangeSession godoc
// @Summary Exchange iam_session cookie for API tokens (after federated sign-in)
// @Description Reads the browser iam_session cookie set by federation callback and returns JWT tokens for the SPA dashboard. Used by GET /login/federation/complete.
// @Tags Auth
// @Produce json
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.LoginResponse} "When tenant is resolved. When multiple tenants apply, `data` is TenantSelectionResponse."
// @Failure 401 {object} object{meta=dtos.Meta}
// @Router /login/session [post]
func (h *AuthHandler) ExchangeSession(c echo.Context) error {
	ck, err := c.Cookie(constants.SessionCookieName)
	if err != nil || ck == nil || strings.TrimSpace(ck.Value) == "" {
		return h.UnauthorizedErrorResponse(c, "No active session")
	}

	ctx := c.Request().Context()
	sess, err := h.sessionService.GetSession(ctx, ck.Value)
	if err != nil {
		return h.HandleError(c, err)
	}

	u, err := h.userService.GetOneByID(ctx, sess.UserID)
	if err != nil {
		return h.HandleError(c, err)
	}

	loginOut, selectionOut, err := h.userService.CompleteAuth(ctx, u, services.TenantResolveInput{
		Host:          c.Request().Host,
		TenantIDParam: sess.TenantID,
		UserID:        u.ID,
	})
	if err != nil {
		return h.HandleError(c, err)
	}
	if selectionOut != nil {
		return h.SuccessResponse(c, "Tenant selection required", selectionOut, nil)
	}

	return h.SuccessResponse(c, "Session exchanged successfully", loginOut, nil)
}

// ListFederationProviders godoc
// @Summary List enabled federation sign-in methods for a tenant
// @Tags Auth
// @Produce json
// @Param tenant_id query string false "Tenant UUID (defaults to DEFAULT_TENANT_ID)"
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} object{meta=dtos.Meta,data=[]dtos.PublicFederationProviderResponse}
// @Router /federation/providers [get]
func (h *AuthHandler) ListFederationProviders(c echo.Context) error {
	tenantID := strings.TrimSpace(c.QueryParam("tenant_id"))
	if tenantID == "" {
		tenantID = h.cfg.DefaultTenantID
	}
	providers, err := h.federationService.ListAvailableProviders(c.Request().Context(), tenantID)
	if err != nil {
		return h.HandleError(c, err)
	}
	out := make([]dtos.PublicFederationProviderResponse, 0, len(providers))
	for _, p := range providers {
		out = append(out, dtos.PublicFederationProviderResponse{Provider: p.Provider, Name: p.Name})
	}
	page, pageable := dtos.PaginateSlice(out, pageableFromQuery(c))
	return h.SuccessResponse(c, "Federation providers retrieved successfully", page, pageable)
}

// Logout godoc
// @Summary Log out everywhere for the current user (SSO session + all refresh tokens)
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} object{meta=dtos.Meta}
// @Router /logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	uid, ok := c.Get(auth.EchoContextUserIDKey).(string)
	if !ok || uid == "" {
		return h.UnauthorizedErrorResponse(c, "User not authenticated")
	}
	ctx := c.Request().Context()
	if err := h.sessionService.InvalidateAllForUser(ctx, uid); err != nil {
		return h.HandleError(c, err)
	}
	if err := h.userService.RevokeAllRefreshTokensForUser(ctx, uid); err != nil {
		return h.HandleError(c, err)
	}
	h.auditService.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionAuthLogout,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      uid,
		ResourceType: constants.AuditResourceTypeUser,
		ResourceID:   uid,
	})
	h.clearSessionCookie(c)
	return h.SuccessResponse(c, "Logged out successfully", map[string]any{}, nil)
}

func (h *AuthHandler) clearSessionCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.cfg.AppEnv == config.EnvironmentProduction,
		MaxAge:   -1,
	})
}

func (h *AuthHandler) setSessionCookie(c echo.Context, sid string, ttl time.Duration) {
	c.SetCookie(&http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.cfg.AppEnv == config.EnvironmentProduction,
		MaxAge:   int(ttl.Seconds()),
	})
}

// Refresh godoc
// @Summary Exchange refresh token for a new access token (and rotated refresh token)
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dtos.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.LoginResponse}
// @Router /refresh [post]
func (h *AuthHandler) Refresh(c echo.Context) error {
	var req dtos.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return h.HandleError(c, errors.ValidationError("Invalid request body", err))
	}
	if err := h.validator.Struct(req); err != nil {
		fieldErrors := errors.ParseValidationErrors(err)
		if len(fieldErrors) > 0 {
			return h.HandleError(c, errors.ValidationErrorWithDetails("Validation failed", err, fieldErrors))
		}
		return h.HandleError(c, errors.ValidationError("Validation failed", err))
	}

	out, err := h.userService.Refresh(c.Request().Context(), &req)
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "Token refreshed successfully", out, nil)
}

// Me godoc
// @Summary Current user profile
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.UserResponse}
// @Router /me [get]
func (h *AuthHandler) Me(c echo.Context) error {
	uid, ok := c.Get(auth.EchoContextUserIDKey).(string)
	if !ok || uid == "" {
		return h.UnauthorizedErrorResponse(c, "User not authenticated")
	}
	resp, err := h.buildCurrentUserResponse(c, uid)
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "User retrieved successfully", resp, nil)
}

// UpdateMe godoc
// @Summary Update current user profile
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dtos.UpdateProfileRequest true "Profile fields"
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.UserResponse}
// @Router /me [patch]
func (h *AuthHandler) UpdateMe(c echo.Context) error {
	uid, ok := c.Get(auth.EchoContextUserIDKey).(string)
	if !ok || uid == "" {
		return h.UnauthorizedErrorResponse(c, "User not authenticated")
	}
	var req dtos.UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return h.HandleError(c, errors.ValidationError("Invalid request body", err))
	}
	if err := h.validator.Struct(req); err != nil {
		fieldErrors := errors.ParseValidationErrors(err)
		if len(fieldErrors) > 0 {
			return h.HandleError(c, errors.ValidationErrorWithDetails("Validation failed", err, fieldErrors))
		}
		return h.HandleError(c, errors.ValidationError("Validation failed", err))
	}
	if _, err := h.userService.UpdateProfile(c.Request().Context(), uid, &req); err != nil {
		return h.HandleError(c, err)
	}
	resp, err := h.buildCurrentUserResponse(c, uid)
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "Profile updated successfully", resp, nil)
}

func (h *AuthHandler) buildCurrentUserResponse(c echo.Context, uid string) (*dtos.UserResponse, error) {
	activeTenant, _ := c.Get(auth.EchoContextTenantIDKey).(string)
	ctx := c.Request().Context()
	user, err := h.userService.GetOneByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	mfaEnabled, err := h.mfaService.HasActiveMFA(ctx, uid)
	if err != nil {
		return nil, err
	}
	tenants, err := h.userService.ListMemberships(ctx, uid)
	if err != nil {
		return nil, err
	}
	resp := dtos.NewUserResponse(user)
	resp.MFAEnabled = mfaEnabled
	resp.ActiveTenantID = activeTenant
	resp.Tenants = tenants
	return resp, nil
}

// ListMyTenants godoc
// @Summary List tenants the current user can access
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Success 200 {object} object{meta=dtos.Meta,data=[]dtos.TenantSummary}
// @Router /me/tenants [get]
func (h *AuthHandler) ListMyTenants(c echo.Context) error {
	uid, ok := c.Get(auth.EchoContextUserIDKey).(string)
	if !ok || uid == "" {
		return h.UnauthorizedErrorResponse(c, "User not authenticated")
	}
	tenants, pageable, err := h.userService.ListMembershipsPaginated(c.Request().Context(), uid, pageableFromQuery(c))
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "Tenants retrieved successfully", tenants, pageable)
}

// SelectTenant godoc
// @Summary Complete login by selecting a tenant
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dtos.TenantSelectRequest true "Selection token and tenant"
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.LoginResponse}
// @Router /tenants/select [post]
func (h *AuthHandler) SelectTenant(c echo.Context) error {
	var req dtos.TenantSelectRequest
	if err := c.Bind(&req); err != nil {
		return h.HandleError(c, errors.ValidationError("Invalid request body", err))
	}
	if err := h.validator.Struct(req); err != nil {
		return h.HandleError(c, errors.ValidationError("Validation failed", err))
	}
	out, err := h.userService.SelectTenant(c.Request().Context(), &req)
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "Tenant selected", out, nil)
}

// SwitchTenant godoc
// @Summary Switch active tenant and re-issue tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dtos.TenantSwitchRequest true "Target tenant"
// @Success 200 {object} object{meta=dtos.Meta,data=dtos.LoginResponse}
// @Router /tenants/switch [post]
func (h *AuthHandler) SwitchTenant(c echo.Context) error {
	uid, ok := c.Get(auth.EchoContextUserIDKey).(string)
	if !ok || uid == "" {
		return h.UnauthorizedErrorResponse(c, "User not authenticated")
	}
	var req dtos.TenantSwitchRequest
	if err := c.Bind(&req); err != nil {
		return h.HandleError(c, errors.ValidationError("Invalid request body", err))
	}
	if err := h.validator.Struct(req); err != nil {
		return h.HandleError(c, errors.ValidationError("Validation failed", err))
	}
	out, err := h.userService.SwitchTenant(c.Request().Context(), uid, req.TenantID)
	if err != nil {
		return h.HandleError(c, err)
	}
	return h.SuccessResponse(c, "Tenant switched", out, nil)
}
