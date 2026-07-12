package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestAdminHandler_GetStats(t *testing.T) {
	h := ProvideAdminHandler(stubAdminListService{}, validator.New())
	c, rec := newJSONContext(http.MethodGet, "/admin/stats", "")

	err := h.GetStats(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"total_users":10`)
}

func TestAdminHandler_ListUsers(t *testing.T) {
	h := ProvideAdminHandler(stubAdminListService{}, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/users?tenant_id="+testTenantID+"&search=admin", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.ListUsers(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "admin@example.com")
}

func TestAdminHandler_ListTenants(t *testing.T) {
	h := ProvideAdminHandler(stubAdminListService{}, validator.New())
	c, rec := newJSONContext(http.MethodGet, "/admin/tenants?page=1", "")

	err := h.ListTenants(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "Acme")
}

func TestAdminHandler_ListClients(t *testing.T) {
	h := ProvideAdminHandler(stubAdminListService{}, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/clients?tenant_id="+testTenantID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.ListClients(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "spa-app")
}

func TestAdminHandler_ListIdentityProviders(t *testing.T) {
	h := ProvideAdminHandler(stubAdminListService{}, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/tenants/"+testTenantID+"/identity-providers", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("tenantId")
	c.SetParamValues(testTenantID)

	err := h.ListIdentityProviders(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"provider":"google"`)
}

func TestAdminHandler_ListIdentityProviders_missingTenantID(t *testing.T) {
	h := ProvideAdminHandler(stubAdminListService{}, validator.New())
	c, rec := newJSONContext(http.MethodGet, "/admin/tenants//identity-providers", "")
	c.SetParamNames("tenantId")
	c.SetParamValues("")

	err := h.ListIdentityProviders(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAdminHandler_PatchIdentityProvider(t *testing.T) {
	svc := stubAdminListService{}
	h := ProvideAdminHandler(svc, validator.New())
	body := `{"enabled":true,"oauth_client_id":"cid","oauth_client_secret":"sec"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/admin/tenants/"+testTenantID+"/identity-providers/google", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("tenantId", "provider")
	c.SetParamValues(testTenantID, "google")

	err := h.PatchIdentityProvider(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestAdminHandler_AddMember(t *testing.T) {
	h := ProvideAdminHandler(stubAdminListService{}, validator.New())
	body := `{"email":"member@example.com","role":"member"}`
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/admin/tenants/"+testTenantID+"/members", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("tenantId")
	c.SetParamValues(testTenantID)

	err := h.AddMember(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestAdminHandler_RemoveMember(t *testing.T) {
	h := ProvideAdminHandler(stubAdminListService{}, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/admin/tenants/"+testTenantID+"/members/"+testUserID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("tenantId", "userId")
	c.SetParamValues(testTenantID, testUserID)

	err := h.RemoveMember(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, rec.Code)
}

func TestAdminHandler_ListAuditLogs_invalidFrom(t *testing.T) {
	h := ProvideAdminHandler(stubAdminAuditService{}, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/audit-logs?from=not-a-date", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.ListAuditLogs(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAdminHandler_ListAuditLogs_invalidTo(t *testing.T) {
	h := ProvideAdminHandler(stubAdminAuditService{}, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/audit-logs?to=bad", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.ListAuditLogs(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAdminHandler_ForceLogoutUser(t *testing.T) {
	svc := &stubAdminUserActionsService{}
	h := ProvideAdminHandler(svc, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/admin/users/"+testUserID+"/force-logout", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("userId")
	c.SetParamValues(testUserID)
	c.Set(auth.EchoContextUserIDKey, "actor-1")

	err := h.ForceLogoutUser(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, rec.Code)
	require.Equal(t, testUserID, svc.forceLogoutTarget)
}

func TestAdminHandler_ResetUserMFA(t *testing.T) {
	svc := &stubAdminUserActionsService{}
	h := ProvideAdminHandler(svc, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/admin/users/"+testUserID+"/reset-mfa", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("userId")
	c.SetParamValues(testUserID)
	c.Set(auth.EchoContextUserIDKey, "actor-1")

	err := h.ResetUserMFA(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, rec.Code)
	require.Equal(t, testUserID, svc.resetMFATarget)
}

func TestAdminHandler_ResetUserPasskeys(t *testing.T) {
	svc := &stubAdminUserActionsService{}
	h := ProvideAdminHandler(svc, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/admin/users/"+testUserID+"/reset-passkey", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("userId")
	c.SetParamValues(testUserID)
	c.Set(auth.EchoContextUserIDKey, "actor-1")

	err := h.ResetUserPasskeys(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, rec.Code)
	require.Equal(t, testUserID, svc.resetPasskeyTarget)
}

func TestAdminHandler_GetClientUsage(t *testing.T) {
	h := ProvideAdminHandler(&stubAdminUserActionsService{}, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/clients/"+testClientPK+"/usage", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("clientId")
	c.SetParamValues(testClientPK)

	err := h.GetClientUsage(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), testClientPK)
}

func TestAdminHandler_ListLoginHistory(t *testing.T) {
	h := ProvideAdminHandler(stubAdminListService{}, validator.New())
	c, rec := newJSONContext(http.MethodGet, "/admin/login-history?page=1", "")

	err := h.ListLoginHistory(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), constants.AuditActionAuthLogin)
}

func TestAdminHandler_ListLoginHistory_invalidFrom(t *testing.T) {
	h := ProvideAdminHandler(stubAdminListService{}, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/admin/login-history?from=invalid", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.ListLoginHistory(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAdminHandler_DisableUser_missingActor(t *testing.T) {
	h := ProvideAdminHandler(stubAdminAuditService{}, validator.New())
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/admin/users/"+testUserID+"/disable", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("userId")
	c.SetParamValues(testUserID)

	err := h.DisableUser(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAdminHandler_GetUser_missingUserID(t *testing.T) {
	h := ProvideAdminHandler(stubAdminAuditService{}, validator.New())
	c, rec := newJSONContext(http.MethodGet, "/admin/users/", "")
	c.SetParamNames("userId")
	c.SetParamValues("")

	err := h.GetUser(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}
