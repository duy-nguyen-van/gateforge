package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestAuthHandler_Register_success(t *testing.T) {
	h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
	c, rec := newJSONContext(http.MethodPost, "/register", `{"email":"user@example.com","password":"secretpass"}`)

	err := h.Register(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "user@example.com")
}

func TestAuthHandler_Register_validationError(t *testing.T) {
	h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
	c, rec := newJSONContext(http.MethodPost, "/register", `{"email":"not-an-email","password":"short"}`)

	err := h.Register(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAuthHandler_Register_invalidBody(t *testing.T) {
	h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
	c, rec := newJSONContext(http.MethodPost, "/register", `{`)

	err := h.Register(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAuthHandler_Login_success(t *testing.T) {
	h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{createSID: "sess-1"}, &stubAuthMFAService{}, &stubAuthFederationService{})
	body := `{"email":"user@example.com","password":"secretpass","remember_me":true}`
	c, rec := newJSONContext(http.MethodPost, "/login", body)

	err := h.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "access-token")
	cookie := rec.Result().Cookies()
	require.Len(t, cookie, 1)
	require.Equal(t, constants.SessionCookieName, cookie[0].Name)
	require.Equal(t, "sess-1", cookie[0].Value)
}

func TestAuthHandler_Login_mfaRequired(t *testing.T) {
	h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{hasMFA: true, ticket: "ticket-1"}, &stubAuthFederationService{})
	c, rec := newJSONContext(http.MethodPost, "/login", `{"email":"user@example.com","password":"secretpass"}`)

	err := h.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"mfa_required":true`)
	require.Contains(t, rec.Body.String(), "ticket-1")
}

func TestAuthHandler_Login_tenantSelection(t *testing.T) {
	h := authHandler(
		&stubAuthUserService{selection: &dtos.TenantSelectionResponse{SelectionRequired: true, SelectionToken: "sel-token"}},
		&stubAuthSessionService{},
		&stubAuthMFAService{},
		&stubAuthFederationService{},
	)
	c, rec := newJSONContext(http.MethodPost, "/login", `{"email":"user@example.com","password":"secretpass"}`)

	err := h.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"selection_required":true`)
}

func TestAuthHandler_Login_authFailure(t *testing.T) {
	h := authHandler(&stubAuthUserService{authErr: errUnauthorized()}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
	c, rec := newJSONContext(http.MethodPost, "/login", `{"email":"user@example.com","password":"wrong"}`)

	err := h.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthHandler_Refresh_success(t *testing.T) {
	h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
	c, rec := newJSONContext(http.MethodPost, "/refresh", `{"refresh_token":"rt-123"}`)

	err := h.Refresh(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "refreshed-access")
}

func TestAuthHandler_Refresh_validationError(t *testing.T) {
	h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
	c, rec := newJSONContext(http.MethodPost, "/refresh", `{}`)

	err := h.Refresh(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAuthHandler_Me_success(t *testing.T) {
	h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{hasMFA: true}, &stubAuthFederationService{})
	c, rec := newJSONContext(http.MethodGet, "/me", "")
	c.Set(auth.EchoContextUserIDKey, testUserID)
	c.Set(auth.EchoContextTenantIDKey, testTenantID)

	err := h.Me(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), testTenantID)
	require.Contains(t, rec.Body.String(), `"mfa_enabled":true`)
}

func TestAuthHandler_Me_unauthorized(t *testing.T) {
	h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
	c, rec := newJSONContext(http.MethodGet, "/me", "")

	err := h.Me(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthHandler_UpdateMe_success(t *testing.T) {
	h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
	c, rec := newJSONContext(http.MethodPatch, "/me", `{"first_name":"Updated"}`)
	c.Set(auth.EchoContextUserIDKey, testUserID)

	err := h.UpdateMe(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "Profile updated successfully")
}

func TestAuthHandler_UpdateMe_unauthorized(t *testing.T) {
	h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
	c, rec := newJSONContext(http.MethodPatch, "/me", `{"first_name":"Updated"}`)

	err := h.UpdateMe(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthHandler_Logout_success(t *testing.T) {
	cfg := handlerTestConfig()
	h := ProvideAuthHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{}, stubAuditService{}, cfg, validator.New())
	c, rec := newJSONContext(http.MethodPost, "/logout", "")
	c.Set(auth.EchoContextUserIDKey, testUserID)
	c.Request().AddCookie(&http.Cookie{Name: constants.SessionCookieName, Value: "existing-session"})

	err := h.Logout(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	cookies := rec.Result().Cookies()
	require.NotEmpty(t, cookies)
	found := false
	for _, ck := range cookies {
		if ck.Name == constants.SessionCookieName {
			found = true
			require.Equal(t, "", ck.Value)
			require.Equal(t, -1, ck.MaxAge)
		}
	}
	require.True(t, found)
}

func TestAuthHandler_Logout_unauthorized(t *testing.T) {
	h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
	c, rec := newJSONContext(http.MethodPost, "/logout", "")

	err := h.Logout(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthHandler_ExchangeSession_noCookie(t *testing.T) {
	h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
	c, rec := newJSONContext(http.MethodPost, "/login/session", "")

	err := h.ExchangeSession(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthHandler_ExchangeSession_success(t *testing.T) {
	h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
	c, rec := newJSONContext(http.MethodPost, "/login/session", "")
	c.Request().AddCookie(&http.Cookie{Name: constants.SessionCookieName, Value: "sess-abc"})

	err := h.ExchangeSession(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "access-token")
}

func TestAuthHandler_ExchangeSession_tenantSelection(t *testing.T) {
	h := authHandler(
		&stubAuthUserService{selection: &dtos.TenantSelectionResponse{SelectionRequired: true}},
		&stubAuthSessionService{},
		&stubAuthMFAService{},
		&stubAuthFederationService{},
	)
	c, rec := newJSONContext(http.MethodPost, "/login/session", "")
	c.Request().AddCookie(&http.Cookie{Name: constants.SessionCookieName, Value: "sess-abc"})

	err := h.ExchangeSession(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"selection_required":true`)
}

func TestAuthHandler_ListFederationProviders_success(t *testing.T) {
	h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/federation/providers?page=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.ListFederationProviders(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "Google")
}

func TestAuthHandler_ListFederationProviders_defaultTenant(t *testing.T) {
	h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{listErr: errors.InternalError("boom", nil)})
	c, rec := newJSONContext(http.MethodGet, "/federation/providers", "")

	err := h.ListFederationProviders(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestAuthHandler_ListMyTenants_success(t *testing.T) {
	h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
	c, rec := newJSONContext(http.MethodGet, "/me/tenants?page=1", "")
	c.Set(auth.EchoContextUserIDKey, testUserID)

	err := h.ListMyTenants(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "Default")
}

func TestAuthHandler_SelectTenant_success(t *testing.T) {
	h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
	body := `{"selection_token":"sel-token","tenant_id":"` + testTenantID + `"}`
	c, rec := newJSONContext(http.MethodPost, "/tenants/select", body)

	err := h.SelectTenant(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "access-token")
}

func TestAuthHandler_SwitchTenant_success(t *testing.T) {
	h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
	c, rec := newJSONContext(http.MethodPost, "/tenants/switch", `{"tenant_id":"`+testTenantID+`"}`)
	c.Set(auth.EchoContextUserIDKey, testUserID)

	err := h.SwitchTenant(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "switched-token")
}

func TestAuthHandler_setSessionCookie_productionSecure(t *testing.T) {
	cfg := handlerTestConfig()
	cfg.AppEnv = config.EnvironmentProduction
	h := ProvideAuthHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{}, stubAuditService{}, cfg, validator.New())
	c, rec := newJSONContext(http.MethodPost, "/login", `{"email":"user@example.com","password":"secretpass"}`)

	err := h.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	cookies := rec.Result().Cookies()
	require.NotEmpty(t, cookies)
	require.True(t, cookies[0].Secure)
}

func newJSONContext(method, path, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	var req *http.Request
	if strings.TrimSpace(body) == "" {
		req = httptest.NewRequest(method, path, nil)
	} else {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec), rec
}
