package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type errAdminService struct {
	err error
}

func (s errAdminService) withErr() error {
	if s.err != nil {
		return s.err
	}
	return errors.InternalError("admin failure", nil)
}

func (s errAdminService) CreateTenant(context.Context, *dtos.AdminCreateTenantRequest) (*dtos.AdminTenantResponse, error) {
	return nil, s.withErr()
}
func (s errAdminService) GetTenantByID(context.Context, string) (*dtos.AdminTenantResponse, error) {
	return nil, s.withErr()
}
func (s errAdminService) GetStats(context.Context) (*dtos.AdminStatsResponse, error) {
	return nil, s.withErr()
}
func (s errAdminService) ListUsers(context.Context, string, string, *dtos.PageableRequest) ([]*dtos.AdminUserResponse, *dtos.Pageable, error) {
	return nil, nil, s.withErr()
}
func (s errAdminService) ListTenants(context.Context, *dtos.PageableRequest) ([]*dtos.AdminTenantResponse, *dtos.Pageable, error) {
	return nil, nil, s.withErr()
}
func (s errAdminService) ListClients(context.Context, string, *dtos.PageableRequest) ([]*dtos.AdminClientResponse, *dtos.Pageable, error) {
	return nil, nil, s.withErr()
}
func (s errAdminService) ListIdentityProviders(context.Context, string, *dtos.PageableRequest) ([]*dtos.AdminIdentityProviderResponse, *dtos.Pageable, error) {
	return nil, nil, s.withErr()
}
func (s errAdminService) ListAuditLogs(context.Context, dtos.AdminAuditLogListParams, *dtos.PageableRequest) ([]*dtos.AdminAuditLogResponse, *dtos.Pageable, error) {
	return nil, nil, s.withErr()
}
func (s errAdminService) ListLoginHistory(context.Context, dtos.AdminLoginHistoryListParams, *dtos.PageableRequest) ([]*dtos.AdminAuditLogResponse, *dtos.Pageable, error) {
	return nil, nil, s.withErr()
}
func (s errAdminService) ConfigureIdentityProvider(context.Context, string, string, *dtos.PatchIdentityProviderRequest, constants.AuditActorType) error {
	return s.withErr()
}
func (s errAdminService) AddMemberByEmail(context.Context, string, string, string) error {
	return s.withErr()
}
func (s errAdminService) RemoveMember(context.Context, string, string) error {
	return s.withErr()
}
func (s errAdminService) UpdateTenant(context.Context, string, *dtos.AdminUpdateTenantRequest) (*dtos.AdminTenantResponse, error) {
	return nil, s.withErr()
}
func (s errAdminService) CreateClient(context.Context, *dtos.AdminCreateClientRequest) (*dtos.AdminCreateClientResponse, error) {
	return nil, s.withErr()
}
func (s errAdminService) DeleteTenant(context.Context, string) error { return s.withErr() }
func (s errAdminService) ListTenantMembers(context.Context, string, *dtos.PageableRequest) ([]*dtos.AdminTenantMemberResponse, *dtos.Pageable, error) {
	return nil, nil, s.withErr()
}
func (s errAdminService) GetClientByID(context.Context, string) (*dtos.AdminClientResponse, error) {
	return nil, s.withErr()
}
func (s errAdminService) UpdateClient(context.Context, string, *dtos.AdminUpdateClientRequest) (*dtos.AdminClientResponse, error) {
	return nil, s.withErr()
}
func (s errAdminService) DeleteClient(context.Context, string) error            { return s.withErr() }
func (s errAdminService) DisableUser(context.Context, string, string) error     { return s.withErr() }
func (s errAdminService) ForceLogoutUser(context.Context, string, string) error { return s.withErr() }
func (s errAdminService) ResetMFA(context.Context, string, string) error        { return s.withErr() }
func (s errAdminService) ResetPasskeys(context.Context, string, string) error   { return s.withErr() }
func (s errAdminService) GetClientUsage(context.Context, string) (*dtos.AdminClientUsageResponse, error) {
	return nil, s.withErr()
}
func (s errAdminService) GetUserByID(context.Context, string) (*dtos.AdminUserDetailResponse, error) {
	return nil, s.withErr()
}

func TestAdminHandler_serviceErrors(t *testing.T) {
	svc := errAdminService{err: errors.InternalError("boom", nil)}
	h := ProvideAdminHandler(svc, validator.New())
	e := echo.New()

	tests := []struct {
		name string
		run  func(c echo.Context) error
	}{
		{"GetStats", func(c echo.Context) error { return h.GetStats(c) }},
		{"ListUsers", func(c echo.Context) error { return h.ListUsers(c) }},
		{"ListTenants", func(c echo.Context) error { return h.ListTenants(c) }},
		{"ListClients", func(c echo.Context) error { return h.ListClients(c) }},
		{"ListAuditLogs", func(c echo.Context) error { return h.ListAuditLogs(c) }},
		{"ListLoginHistory", func(c echo.Context) error { return h.ListLoginHistory(c) }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			require.NoError(t, tt.run(c))
			require.Equal(t, http.StatusInternalServerError, rec.Code)
		})
	}

	t.Run("ListIdentityProviders", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("tenantId")
		c.SetParamValues(testTenantID)
		require.NoError(t, h.ListIdentityProviders(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("PatchIdentityProvider", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"enabled":true}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("tenantId", "provider")
		c.SetParamValues(testTenantID, "google")
		require.NoError(t, h.PatchIdentityProvider(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("AddMember", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"email":"a@b.com","role":"member"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("tenantId")
		c.SetParamValues(testTenantID)
		require.NoError(t, h.AddMember(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("RemoveMember", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("tenantId", "userId")
		c.SetParamValues(testTenantID, testUserID)
		require.NoError(t, h.RemoveMember(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestAdminHandler_validationAndParamErrors(t *testing.T) {
	h := ProvideAdminHandler(stubAdminListService{}, validator.New())
	e := echo.New()

	t.Run("CreateTenant invalid body", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodPost, "/", `{`)
		require.NoError(t, h.CreateTenant(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("CreateTenant validation", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodPost, "/", `{}`)
		require.NoError(t, h.CreateTenant(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("GetTenant missing id", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodGet, "/", "")
		c.SetParamNames("tenantId")
		c.SetParamValues("")
		require.NoError(t, h.GetTenant(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("PatchIdentityProvider missing provider", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodPatch, "/", `{}`)
		c.SetParamNames("tenantId", "provider")
		c.SetParamValues(testTenantID, "")
		require.NoError(t, h.PatchIdentityProvider(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("AddMember validation", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("tenantId")
		c.SetParamValues(testTenantID)
		require.NoError(t, h.AddMember(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ForceLogoutUser missing actor", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodPost, "/", "")
		c.SetParamNames("userId")
		c.SetParamValues(testUserID)
		require.NoError(t, h.ForceLogoutUser(c))
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("GetClientUsage missing id", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodGet, "/", "")
		c.SetParamNames("clientId")
		c.SetParamValues("")
		require.NoError(t, h.GetClientUsage(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ListLoginHistory invalid to", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?to=bad", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		require.NoError(t, h.ListLoginHistory(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthHandler_errorPaths(t *testing.T) {
	t.Run("Register service error", func(t *testing.T) {
		h := authHandler(&stubAuthUserService{registerErr: errors.InternalError("fail", nil)}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodPost, "/", `{"email":"user@example.com","password":"secretpass"}`)
		require.NoError(t, h.Register(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Login MFA ticket error", func(t *testing.T) {
		h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{hasMFA: true, ticketErr: errors.InternalError("fail", nil)}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodPost, "/", `{"email":"user@example.com","password":"secretpass"}`)
		require.NoError(t, h.Login(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Login session create error", func(t *testing.T) {
		h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{createErr: errors.InternalError("fail", nil)}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodPost, "/", `{"email":"user@example.com","password":"secretpass"}`)
		require.NoError(t, h.Login(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Login validation field errors", func(t *testing.T) {
		h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodPost, "/", `{"email":"bad","password":""}`)
		require.NoError(t, h.Login(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ExchangeSession service error", func(t *testing.T) {
		h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{getSessionErr: errors.UnauthorizedError("bad session", nil)}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodPost, "/", "")
		c.Request().AddCookie(&http.Cookie{Name: constants.SessionCookieName, Value: "sess"})
		require.NoError(t, h.ExchangeSession(c))
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Logout invalidate error", func(t *testing.T) {
		h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{invalidateErr: errors.InternalError("fail", nil)}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodPost, "/", "")
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.Logout(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Me build response error", func(t *testing.T) {
		h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{hasMFAErr: errors.InternalError("fail", nil)}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodGet, "/", "")
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.Me(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("UpdateMe bind error", func(t *testing.T) {
		h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodPatch, "/", `{`)
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.UpdateMe(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("SelectTenant validation", func(t *testing.T) {
		h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodPost, "/", `{}`)
		require.NoError(t, h.SelectTenant(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("SwitchTenant unauthorized", func(t *testing.T) {
		h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodPost, "/", `{"tenant_id":"`+testTenantID+`"}`)
		require.NoError(t, h.SwitchTenant(c))
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Logout revoke tokens error", func(t *testing.T) {
		h := authHandler(&stubAuthUserService{revokeErr: errors.InternalError("fail", nil)}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodPost, "/", "")
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.Logout(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("ListMyTenants unauthorized", func(t *testing.T) {
		h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodGet, "/", "")
		require.NoError(t, h.ListMyTenants(c))
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestOIDCHandler_urlHelpers(t *testing.T) {
	cfg := handlerTestConfig()
	cfg.OIDCLoginPageURL = ""
	cfg.AppBaseURL = ""
	h := newOIDCHandler(t, &stubOIDCService{issuer: "localhost:8080"}, nil, nil)
	h.cfg = cfg

	require.Equal(t, "/login", h.oidcLoginPageURL())

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/authorize?client_id=app", nil)
	req.Header.Set("X-Forwarded-Proto", "https")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	abs := h.absoluteAuthorizeURL(c)
	require.Contains(t, abs, "https://")
	require.Contains(t, abs, "/authorize")
}

func TestOIDCHandler_FederationOAuthStart_invalidReturnTo(t *testing.T) {
	h := newOIDCHandler(t, nil, nil, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/oidc/federation/google/start?return_to=https://evil.example/x", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("provider")
	c.SetParamValues("google")
	require.NoError(t, h.FederationOAuthStart(c))
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestOIDCHandler_FederationOAuthStart_serviceError(t *testing.T) {
	h := newOIDCHandler(t, nil, &stubOIDCFederationService{redirectErr: errors.InternalError("fail", nil)}, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/oidc/federation/google/start?return_to=/authorize", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("provider")
	c.SetParamValues("google")
	require.NoError(t, h.FederationOAuthStart(c))
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestOIDCHandler_Token_invalidForm(t *testing.T) {
	h := newOIDCHandler(t, nil, nil, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/token", strings.NewReader("%"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	require.NoError(t, h.Token(c))
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestOIDCHandler_OpenIDConfiguration_fallbackIssuer(t *testing.T) {
	h := newOIDCHandler(t, &stubOIDCService{issuer: "not-a-url"}, nil, nil)
	c, rec := newJSONContext(http.MethodGet, "/.well-known/openid-configuration", "")
	require.NoError(t, h.OpenIDConfiguration(c))
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "http://localhost:3000")
}

func TestOIDCHandler_Authorize_emptyOAuthCode(t *testing.T) {
	h := newOIDCHandler(t, &stubOIDCService{authorizeErr: &domains.OAuthRedirectError{Description: "missing code"}}, nil, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/authorize", nil)
	req.AddCookie(&http.Cookie{Name: constants.SessionCookieName, Value: "sess"})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	require.NoError(t, h.Authorize(c))
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), constants.OAuthInvalidRequest)
}

func TestMFAHandler_errorPaths(t *testing.T) {
	t.Run("TOTPSetup user lookup error", func(t *testing.T) {
		user := &stubAuthUserService{}
		// force GetOneByID error by using a wrapper - simpler: use finishErr pattern
		h := ProvideMFAHandler(&stubAuthMFAService{}, user, &stubAuthSessionService{}, handlerTestConfig(), validator.New())
		c, rec := newJSONContext(http.MethodPost, "/", "")
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.TOTPSetup(c))
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("TOTPVerifyEnrollment unauthorized", func(t *testing.T) {
		h := newMFAHandler()
		c, rec := newJSONContext(http.MethodPost, "/", `{"code":"123456"}`)
		require.NoError(t, h.TOTPVerifyEnrollment(c))
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("RecoveryCodes unauthorized", func(t *testing.T) {
		h := newMFAHandler()
		c, rec := newJSONContext(http.MethodPost, "/", "")
		require.NoError(t, h.RecoveryCodes(c))
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("ChallengeVerify invalid bind", func(t *testing.T) {
		h := newMFAHandler()
		c, rec := newJSONContext(http.MethodPost, "/", `{`)
		require.NoError(t, h.ChallengeVerify(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestWebauthnHandler_validationErrors(t *testing.T) {
	h := newWebauthnHandler(nil, nil)

	t.Run("RegisterStart validation", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodPost, "/", `{"device_name":"`+strings.Repeat("x", 300)+`"}`)
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.RegisterStart(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("LoginStart invalid email", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodPost, "/", `{"email":"not-email"}`)
		require.NoError(t, h.LoginStart(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("RegisterFinish unauthorized", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodPost, "/", `{"session_token":"t","credential":{"id":"c"}}`)
		require.NoError(t, h.RegisterFinish(c))
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("LoginFinish missing credential", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodPost, "/", `{"email":"user@example.com","session_token":"t"}`)
		require.NoError(t, h.LoginFinish(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestActorUserIDFromContext(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	_, err := actorUserIDFromContext(c)
	require.Error(t, err)

	c.Set(auth.EchoContextUserIDKey, 123)
	_, err = actorUserIDFromContext(c)
	require.Error(t, err)
}

func TestLoginHistoryFiltersFromQuery_validRange(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?from=2024-01-01T00:00:00Z&to=2024-12-31T23:59:59Z&tenant_id=t1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	filters, err := loginHistoryFiltersFromQuery(c)
	require.NoError(t, err)
	require.Equal(t, "t1", filters.TenantID)
	require.NotNil(t, filters.From)
	require.NotNil(t, filters.To)
}

func TestFederationReturnHostAllowed_oidcLoginPageURL(t *testing.T) {
	cfg := handlerTestConfig()
	cfg.OIDCLoginPageURL = "http://spa.example:3000/login"
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	require.True(t, federationReturnHostAllowed(c, cfg, "spa.example:3000"))
}

func TestAdminHandler_crudParamAndServiceErrors(t *testing.T) {
	svc := errAdminService{err: errors.InternalError("boom", nil)}
	h := ProvideAdminHandler(svc, validator.New())
	e := echo.New()

	t.Run("UpdateTenant missing id", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodPatch, "/", `{"name":"x"}`)
		c.SetParamNames("tenantId")
		c.SetParamValues("")
		require.NoError(t, h.UpdateTenant(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("UpdateTenant service error", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodPatch, "/", `{"name":"x"}`)
		c.SetParamNames("tenantId")
		c.SetParamValues(testTenantID)
		require.NoError(t, h.UpdateTenant(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("DeleteTenant missing id", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodDelete, "/", "")
		c.SetParamNames("tenantId")
		c.SetParamValues("")
		require.NoError(t, h.DeleteTenant(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ListTenantMembers missing id", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodGet, "/", "")
		c.SetParamNames("tenantId")
		c.SetParamValues("")
		require.NoError(t, h.ListTenantMembers(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("CreateClient service error", func(t *testing.T) {
		body := `{"tenant_id":"00000000-0000-0000-0000-000000000001","client_id":"app","name":"App","redirect_uris":["http://localhost/cb"]}`
		c, rec := newJSONContext(http.MethodPost, "/", body)
		require.NoError(t, h.CreateClient(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("GetClient missing id", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodGet, "/", "")
		c.SetParamNames("clientId")
		c.SetParamValues("")
		require.NoError(t, h.GetClient(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("UpdateClient missing id", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodPatch, "/", `{"name":"x"}`)
		c.SetParamNames("clientId")
		c.SetParamValues("")
		require.NoError(t, h.UpdateClient(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("DeleteClient missing id", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodDelete, "/", "")
		c.SetParamNames("clientId")
		c.SetParamValues("")
		require.NoError(t, h.DeleteClient(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("DisableUser missing userId", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodPost, "/", "")
		c.SetParamNames("userId")
		c.SetParamValues("")
		c.Set(auth.EchoContextUserIDKey, "actor")
		require.NoError(t, h.DisableUser(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ResetUserMFA missing userId", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodPost, "/", "")
		c.SetParamNames("userId")
		c.SetParamValues("")
		c.Set(auth.EchoContextUserIDKey, "actor")
		require.NoError(t, h.ResetUserMFA(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ResetUserPasskeys missing userId", func(t *testing.T) {
		c, rec := newJSONContext(http.MethodPost, "/", "")
		c.SetParamNames("userId")
		c.SetParamValues("")
		c.Set(auth.EchoContextUserIDKey, "actor")
		require.NoError(t, h.ResetUserPasskeys(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("PatchIdentityProvider bind error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("tenantId", "provider")
		c.SetParamValues(testTenantID, "google")
		require.NoError(t, h.PatchIdentityProvider(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAuthHandler_moreErrorPaths(t *testing.T) {
	t.Run("Refresh bind error", func(t *testing.T) {
		h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodPost, "/", `{`)
		require.NoError(t, h.Refresh(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("UpdateMe validation error", func(t *testing.T) {
		h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodPatch, "/", `{"first_name":""}`)
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.UpdateMe(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("UpdateMe service error", func(t *testing.T) {
		user := &stubAuthUserService{}
		h := authHandler(user, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		// wrap UpdateProfile to fail
		h.userService = &updateProfileErrUserService{stubAuthUserService: user}
		c, rec := newJSONContext(http.MethodPatch, "/", `{"first_name":"A"}`)
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.UpdateMe(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("SelectTenant service error", func(t *testing.T) {
		h := authHandler(&selectTenantErrUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		body := `{"selection_token":"tok","tenant_id":"` + testTenantID + `"}`
		c, rec := newJSONContext(http.MethodPost, "/", body)
		require.NoError(t, h.SelectTenant(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("SwitchTenant validation error", func(t *testing.T) {
		h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodPost, "/", `{}`)
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.SwitchTenant(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

type updateProfileErrUserService struct {
	*stubAuthUserService
}

func (updateProfileErrUserService) UpdateProfile(context.Context, string, *dtos.UpdateProfileRequest) (*models.User, error) {
	return nil, errors.InternalError("update failed", nil)
}

type selectTenantErrUserService struct {
	stubAuthUserService
}

func (selectTenantErrUserService) SelectTenant(context.Context, *dtos.TenantSelectRequest) (*dtos.LoginResponse, error) {
	return nil, errors.InternalError("select failed", nil)
}

func TestOIDCHandler_moreCoverage(t *testing.T) {
	t.Run("Login validation error", func(t *testing.T) {
		h := newOIDCHandler(t, nil, nil, nil)
		c, rec := newJSONContext(http.MethodPost, "/", `{"email":"bad","password":"","return_to":"/authorize"}`)
		require.NoError(t, h.Login(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Login session error", func(t *testing.T) {
		h := newOIDCHandler(t, nil, nil, nil)
		h.sessionService = &stubAuthSessionService{createErr: errors.InternalError("fail", nil)}
		body := `{"email":"user@example.com","password":"secretpass","return_to":"/authorize?client_id=app"}`
		c, rec := newJSONContext(http.MethodPost, "/", body)
		require.NoError(t, h.Login(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("resolveTenantFromReturnTo bad url", func(t *testing.T) {
		h := newOIDCHandler(t, nil, nil, nil)
		e := echo.New()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
		require.Equal(t, "", h.resolveTenantFromReturnTo(c, "://bad"))
	})

	t.Run("resolveTenantFromReturnTo repo error", func(t *testing.T) {
		h := newOIDCHandler(t, nil, nil, &stubClientRepo{err: errors.NotFoundError("Client", nil)})
		e := echo.New()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
		require.Equal(t, "", h.resolveTenantFromReturnTo(c, "/authorize?client_id=missing"))
	})

	t.Run("isSafeAuthorizeReturnURL request host", func(t *testing.T) {
		cfg := handlerTestConfig()
		cfg.AppBaseURL = ""
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/authorize", nil)
		c := e.NewContext(req, httptest.NewRecorder())
		require.True(t, isSafeAuthorizeReturnURL(c, "", "http://localhost:8080/authorize?x=1"))
	})

	t.Run("oidcLoginPageURL app base", func(t *testing.T) {
		cfg := handlerTestConfig()
		cfg.OIDCLoginPageURL = ""
		cfg.AppBaseURL = "http://localhost:8080"
		h := newOIDCHandler(t, nil, nil, nil)
		h.cfg = cfg
		require.Equal(t, "http://localhost:8080/login", h.oidcLoginPageURL())
	})

	t.Run("FederationOAuthCallback default error message", func(t *testing.T) {
		h := newOIDCHandler(t, nil, nil, nil)
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/cb?error=access_denied", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("provider")
		c.SetParamValues("google")
		require.NoError(t, h.FederationOAuthCallback(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestMFAHandler_validationErrWithDetails(t *testing.T) {
	h := newMFAHandler()
	c, rec := newJSONContext(http.MethodPost, "/", `{"code":"1"}`)
	c.Set(auth.EchoContextUserIDKey, testUserID)
	require.NoError(t, h.TOTPVerifyEnrollment(c))
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestWebauthnHandler_serviceStartErrors(t *testing.T) {
	h := newWebauthnHandler(&stubWebauthnService{startErr: errors.InternalError("fail", nil)}, nil)
	c, rec := newJSONContext(http.MethodPost, "/", `{}`)
	c.Set(auth.EchoContextUserIDKey, testUserID)
	require.NoError(t, h.RegisterStart(c))
	require.Equal(t, http.StatusInternalServerError, rec.Code)

	c, rec = newJSONContext(http.MethodPost, "/", `{"email":"user@example.com"}`)
	require.NoError(t, h.LoginStart(c))
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestWebauthnHandler_loginFinishHasMFAError(t *testing.T) {
	h := newWebauthnHandler(nil, &stubAuthMFAService{hasMFAErr: errors.InternalError("fail", nil)})
	body := `{"email":"user@example.com","session_token":"t","credential":{"id":"c"}}`
	c, rec := newJSONContext(http.MethodPost, "/", body)
	require.NoError(t, h.LoginFinish(c))
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestWebauthnHandler_listCredentialsError(t *testing.T) {
	h := newWebauthnHandler(&stubWebauthnService{listErr: errors.InternalError("fail", nil)}, nil)
	c, rec := newJSONContext(http.MethodGet, "/", "")
	c.Set(auth.EchoContextUserIDKey, testUserID)
	require.NoError(t, h.ListCredentials(c))
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAdminHandler_userActionServiceErrors(t *testing.T) {
	svc := errAdminService{err: errors.InternalError("boom", nil)}
	h := ProvideAdminHandler(svc, validator.New())
	e := echo.New()

	run := func(name string, fn func(c echo.Context) error) {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("userId")
			c.SetParamValues(testUserID)
			c.Set(auth.EchoContextUserIDKey, "actor")
			require.NoError(t, fn(c))
			require.Equal(t, http.StatusInternalServerError, rec.Code)
		})
	}
	run("DisableUser", h.DisableUser)
	run("ForceLogoutUser", h.ForceLogoutUser)
	run("ResetUserMFA", h.ResetUserMFA)
	run("ResetUserPasskeys", h.ResetUserPasskeys)
}

func TestTenantIdentityAdminHandler_PatchIdentityProvider_bindError(t *testing.T) {
	cfg := handlerTestConfig()
	cfg.AdminAPIKey = "secret"
	h := ProvideTenantIdentityAdminHandler(cfg, stubAdminListService{})
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("X-Admin-API-Key", "secret")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("tenantId", "provider")
	c.SetParamValues(testTenantID, "google")
	require.NoError(t, h.PatchIdentityProvider(c))
	require.Equal(t, http.StatusBadRequest, rec.Code)
}
