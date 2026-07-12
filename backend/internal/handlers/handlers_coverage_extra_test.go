package handlers

import (
	"context"
	"encoding/base64"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type getUserErrUserService struct {
	stubAuthUserService
}

func (s *getUserErrUserService) GetOneByID(context.Context, string) (*models.User, error) {
	return nil, errors.NotFoundError("User", nil)
}

type listMembershipsErrUserService struct {
	stubAuthUserService
}

func (s *listMembershipsErrUserService) ListMemberships(context.Context, string) ([]dtos.TenantSummary, error) {
	return nil, errors.InternalError("memberships failed", nil)
}

type refreshErrUserService struct {
	stubAuthUserService
}

func (s *refreshErrUserService) Refresh(context.Context, *dtos.RefreshTokenRequest) (*dtos.LoginResponse, error) {
	return nil, errors.UnauthorizedError("invalid refresh token", nil)
}

type issueTokensErrUserService struct {
	stubAuthUserService
}

func (s *issueTokensErrUserService) IssueTokensForUser(context.Context, *models.User, string) (*dtos.LoginResponse, error) {
	return nil, errors.InternalError("token issue failed", nil)
}

type setupTOTPErrMFAService struct {
	stubAuthMFAService
}

func (s *setupTOTPErrMFAService) SetupTOTP(context.Context, string, string) (string, string, error) {
	return "", "", errors.InternalError("setup failed", nil)
}

type recoveryCodesErrMFAService struct {
	stubAuthMFAService
}

func (s *recoveryCodesErrMFAService) RegenerateRecoveryCodes(context.Context, string) ([]string, error) {
	return nil, errors.ForbiddenError("MFA not enabled", nil)
}

func TestAuthHandler_remainingPaths(t *testing.T) {
	t.Run("Login complete auth error", func(t *testing.T) {
		h := authHandler(
			&stubAuthUserService{completeErr: errors.InternalError("complete failed", nil)},
			&stubAuthSessionService{},
			&stubAuthMFAService{},
			&stubAuthFederationService{},
		)
		c, rec := newJSONContext(http.MethodPost, "/", `{"email":"user@example.com","password":"secretpass"}`)
		require.NoError(t, h.Login(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Login HasActiveMFA error", func(t *testing.T) {
		h := authHandler(
			&stubAuthUserService{},
			&stubAuthSessionService{},
			&stubAuthMFAService{hasMFAErr: errors.InternalError("mfa check failed", nil)},
			&stubAuthFederationService{},
		)
		c, rec := newJSONContext(http.MethodPost, "/", `{"email":"user@example.com","password":"secretpass"}`)
		require.NoError(t, h.Login(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("ExchangeSession GetOneByID error", func(t *testing.T) {
		h := authHandler(&getUserErrUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodPost, "/", "")
		c.Request().AddCookie(&http.Cookie{Name: constants.SessionCookieName, Value: "sess"})
		require.NoError(t, h.ExchangeSession(c))
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("ExchangeSession CompleteAuth error", func(t *testing.T) {
		h := authHandler(
			&stubAuthUserService{completeErr: errors.InternalError("complete failed", nil)},
			&stubAuthSessionService{},
			&stubAuthMFAService{},
			&stubAuthFederationService{},
		)
		c, rec := newJSONContext(http.MethodPost, "/", "")
		c.Request().AddCookie(&http.Cookie{Name: constants.SessionCookieName, Value: "sess"})
		require.NoError(t, h.ExchangeSession(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("Refresh service error", func(t *testing.T) {
		h := authHandler(&refreshErrUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodPost, "/", `{"refresh_token":"bad"}`)
		require.NoError(t, h.Refresh(c))
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Me GetOneByID error", func(t *testing.T) {
		h := authHandler(&getUserErrUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodGet, "/", "")
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.Me(c))
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("Me ListMemberships error", func(t *testing.T) {
		h := authHandler(&listMembershipsErrUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodGet, "/", "")
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.Me(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("UpdateMe build response error", func(t *testing.T) {
		h := authHandler(&listMembershipsErrUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodPatch, "/", `{"first_name":"A"}`)
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.UpdateMe(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("ListMyTenants service error", func(t *testing.T) {
		user := &stubAuthUserService{}
		h := authHandler(user, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		h.userService = &listMembershipsPaginatedErrUserService{stubAuthUserService: user}
		c, rec := newJSONContext(http.MethodGet, "/", "")
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.ListMyTenants(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("SwitchTenant service error", func(t *testing.T) {
		h := authHandler(&switchTenantErrUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodPost, "/", `{"tenant_id":"`+testTenantID+`"}`)
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.SwitchTenant(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("ExchangeSession empty cookie value", func(t *testing.T) {
		h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodPost, "/", "")
		c.Request().AddCookie(&http.Cookie{Name: constants.SessionCookieName, Value: "   "})
		require.NoError(t, h.ExchangeSession(c))
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("Login MFA with explicit tenant", func(t *testing.T) {
		h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{hasMFA: true, ticket: "mfa-t"}, &stubAuthFederationService{})
		body := `{"email":"user@example.com","password":"secretpass","tenant_id":"` + testTenantID + `"}`
		c, rec := newJSONContext(http.MethodPost, "/", body)
		require.NoError(t, h.Login(c))
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "mfa-t")
	})

	t.Run("SelectTenant bind error", func(t *testing.T) {
		h := authHandler(&stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, &stubAuthFederationService{})
		c, rec := newJSONContext(http.MethodPost, "/", `{`)
		require.NoError(t, h.SelectTenant(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

type listMembershipsPaginatedErrUserService struct {
	*stubAuthUserService
}

func (s *listMembershipsPaginatedErrUserService) ListMembershipsPaginated(context.Context, string, *dtos.PageableRequest) ([]dtos.TenantSummary, *dtos.Pageable, error) {
	return nil, nil, errors.InternalError("list failed", nil)
}

type switchTenantErrUserService struct {
	stubAuthUserService
}

func (s *switchTenantErrUserService) SwitchTenant(context.Context, string, string) (*dtos.LoginResponse, error) {
	return nil, errors.InternalError("switch failed", nil)
}

func TestMFAHandler_remainingPaths(t *testing.T) {
	t.Run("TOTPSetup user lookup error", func(t *testing.T) {
		h := ProvideMFAHandler(&stubAuthMFAService{}, &getUserErrUserService{}, &stubAuthSessionService{}, handlerTestConfig(), validator.New())
		c, rec := newJSONContext(http.MethodPost, "/", "")
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.TOTPSetup(c))
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("TOTPSetup service error", func(t *testing.T) {
		h := ProvideMFAHandler(&setupTOTPErrMFAService{}, &stubAuthUserService{}, &stubAuthSessionService{}, handlerTestConfig(), validator.New())
		c, rec := newJSONContext(http.MethodPost, "/", "")
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.TOTPSetup(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("TOTPVerifyEnrollment bind error", func(t *testing.T) {
		h := newMFAHandler()
		c, rec := newJSONContext(http.MethodPost, "/", `{`)
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.TOTPVerifyEnrollment(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("RecoveryCodes service error", func(t *testing.T) {
		h := ProvideMFAHandler(&recoveryCodesErrMFAService{}, &stubAuthUserService{}, &stubAuthSessionService{}, handlerTestConfig(), validator.New())
		c, rec := newJSONContext(http.MethodPost, "/", "")
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.RecoveryCodes(c))
		require.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("ChallengeVerify user lookup error", func(t *testing.T) {
		h := ProvideMFAHandler(&stubAuthMFAService{}, &getUserErrUserService{}, &stubAuthSessionService{}, handlerTestConfig(), validator.New())
		c, rec := newJSONContext(http.MethodPost, "/", `{"mfa_ticket":"t","code":"123456"}`)
		require.NoError(t, h.ChallengeVerify(c))
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("ChallengeVerify IssueTokens error", func(t *testing.T) {
		h := ProvideMFAHandler(&stubAuthMFAService{}, &issueTokensErrUserService{}, &stubAuthSessionService{}, handlerTestConfig(), validator.New())
		c, rec := newJSONContext(http.MethodPost, "/", `{"mfa_ticket":"t","code":"123456"}`)
		require.NoError(t, h.ChallengeVerify(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("ChallengeVerify session create error", func(t *testing.T) {
		h := ProvideMFAHandler(&stubAuthMFAService{}, &stubAuthUserService{}, &stubAuthSessionService{createErr: errors.InternalError("session failed", nil)}, handlerTestConfig(), validator.New())
		c, rec := newJSONContext(http.MethodPost, "/", `{"mfa_ticket":"t","code":"123456"}`)
		require.NoError(t, h.ChallengeVerify(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("validationErr plain error", func(t *testing.T) {
		h := newMFAHandler()
		c, rec := newJSONContext(http.MethodPost, "/", "")
		require.NoError(t, h.validationErr(c, stderrors.New("plain validation")))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestWebauthnHandler_remainingPaths(t *testing.T) {
	t.Run("RegisterStart unauthorized", func(t *testing.T) {
		h := newWebauthnHandler(nil, nil)
		c, rec := newJSONContext(http.MethodPost, "/", `{}`)
		require.NoError(t, h.RegisterStart(c))
		require.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("RegisterStart bind error", func(t *testing.T) {
		h := newWebauthnHandler(nil, nil)
		c, rec := newJSONContext(http.MethodPost, "/", `{`)
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.RegisterStart(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("RegisterStart service error", func(t *testing.T) {
		h := newWebauthnHandler(&stubWebauthnService{startErr: errors.InternalError("start failed", nil)}, nil)
		c, rec := newJSONContext(http.MethodPost, "/", `{}`)
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.RegisterStart(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("RegisterFinish validation error", func(t *testing.T) {
		h := newWebauthnHandler(nil, nil)
		c, rec := newJSONContext(http.MethodPost, "/", `{}`)
		c.Set(auth.EchoContextUserIDKey, testUserID)
		require.NoError(t, h.RegisterFinish(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("LoginStart bind error", func(t *testing.T) {
		h := newWebauthnHandler(nil, nil)
		c, rec := newJSONContext(http.MethodPost, "/", `{`)
		require.NoError(t, h.LoginStart(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("LoginFinish MFA ticket error", func(t *testing.T) {
		h := newWebauthnHandler(nil, &stubAuthMFAService{hasMFA: true, ticketErr: errors.InternalError("ticket failed", nil)})
		body := `{"email":"user@example.com","session_token":"t","credential":{"id":"c"}}`
		c, rec := newJSONContext(http.MethodPost, "/", body)
		require.NoError(t, h.LoginFinish(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("LoginFinish complete auth error", func(t *testing.T) {
		userSvc := &stubAuthUserService{completeErr: errors.InternalError("complete failed", nil)}
		h := ProvideWebauthnHandler(&stubWebauthnService{}, userSvc, &stubAuthSessionService{}, &stubAuthMFAService{}, handlerTestConfig(), validator.New())
		body := `{"email":"user@example.com","session_token":"t","credential":{"id":"c"}}`
		c, rec := newJSONContext(http.MethodPost, "/", body)
		require.NoError(t, h.LoginFinish(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("LoginFinish session create error", func(t *testing.T) {
		h := ProvideWebauthnHandler(&stubWebauthnService{}, &stubAuthUserService{}, &stubAuthSessionService{createErr: errors.InternalError("session failed", nil)}, &stubAuthMFAService{}, handlerTestConfig(), validator.New())
		body := `{"email":"user@example.com","session_token":"t","credential":{"id":"c"}}`
		c, rec := newJSONContext(http.MethodPost, "/", body)
		require.NoError(t, h.LoginFinish(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("LoginFinish production secure cookie", func(t *testing.T) {
		cfg := handlerTestConfig()
		cfg.AppEnv = config.EnvironmentProduction
		h := ProvideWebauthnHandler(&stubWebauthnService{}, &stubAuthUserService{}, &stubAuthSessionService{}, &stubAuthMFAService{}, cfg, validator.New())
		body := `{"email":"user@example.com","session_token":"t","credential":{"id":"c"}}`
		c, rec := newJSONContext(http.MethodPost, "/", body)
		require.NoError(t, h.LoginFinish(c))
		require.Equal(t, http.StatusOK, rec.Code)
		cookies := rec.Result().Cookies()
		require.NotEmpty(t, cookies)
		require.True(t, cookies[0].Secure)
	})

	t.Run("validationErr plain error", func(t *testing.T) {
		h := newWebauthnHandler(nil, nil)
		c, rec := newJSONContext(http.MethodPost, "/", "")
		require.NoError(t, h.validationErr(c, stderrors.New("plain validation")))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestOIDCHandler_remainingPaths(t *testing.T) {
	t.Run("Login bind error", func(t *testing.T) {
		h := newOIDCHandler(t, nil, nil, nil)
		c, rec := newJSONContext(http.MethodPost, "/", `{`)
		require.NoError(t, h.Login(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Token basic auth credentials", func(t *testing.T) {
		h := newOIDCHandler(t, nil, nil, nil)
		e := echo.New()
		form := strings.NewReader("grant_type=authorization_code&code=code-1")
		req := httptest.NewRequest(http.MethodPost, "/token", form)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		req.SetBasicAuth("client-id", "client-secret")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		require.NoError(t, h.Token(c))
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("FederationOAuthCallback session error", func(t *testing.T) {
		h := newOIDCHandler(t, nil, &stubOIDCFederationService{}, nil)
		h.sessionService = &stubAuthSessionService{createErr: errors.InternalError("session failed", nil)}
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/cb?code=abc&state=xyz", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("provider")
		c.SetParamValues("google")
		require.NoError(t, h.FederationOAuthCallback(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("isSafeAuthorizeReturnURL relative path", func(t *testing.T) {
		e := echo.New()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
		require.True(t, isSafeAuthorizeReturnURL(c, "http://localhost:8080", "/authorize?client_id=app"))
	})

	t.Run("isSafeAuthorizeReturnURL bad base URL", func(t *testing.T) {
		e := echo.New()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
		require.False(t, isSafeAuthorizeReturnURL(c, "://bad", "http://localhost:8080/authorize"))
	})

	t.Run("isSafeFederationReturnURL empty", func(t *testing.T) {
		cfg := handlerTestConfig()
		e := echo.New()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
		require.False(t, isSafeFederationReturnURL(c, cfg, ""))
		require.False(t, isSafeFederationReturnURL(c, cfg, "   "))
	})

	t.Run("federationReturnHostAllowed request host", func(t *testing.T) {
		cfg := handlerTestConfig()
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/", nil)
		c := e.NewContext(req, httptest.NewRecorder())
		require.True(t, federationReturnHostAllowed(c, cfg, "localhost:8080"))
	})

	t.Run("federationReturnHostAllowed app base URL", func(t *testing.T) {
		cfg := handlerTestConfig()
		cfg.AppBaseURL = "http://localhost:8080/"
		e := echo.New()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
		require.True(t, federationReturnHostAllowed(c, cfg, "localhost:8080"))
	})

	t.Run("isSafeFederationReturnURL bad URL parse", func(t *testing.T) {
		cfg := handlerTestConfig()
		e := echo.New()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
		require.False(t, isSafeFederationReturnURL(c, cfg, "://bad"))
	})

	t.Run("isSafeAuthorizeReturnURL empty raw", func(t *testing.T) {
		e := echo.New()
		c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
		require.False(t, isSafeAuthorizeReturnURL(c, "http://localhost:8080", ""))
	})

	t.Run("resolveUserIDForAuthorize invalid session", func(t *testing.T) {
		h := newOIDCHandler(t, nil, nil, nil)
		h.sessionService = &sessionGetUserIDErrService{}
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/authorize", nil)
		req.AddCookie(&http.Cookie{Name: constants.SessionCookieName, Value: "bad-session"})
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		uid, ok := h.resolveUserIDForAuthorize(c)
		require.False(t, ok)
		require.Empty(t, uid)
	})
}

type sessionGetUserIDErrService struct {
	stubAuthSessionService
}

func (s *sessionGetUserIDErrService) GetUserID(context.Context, string) (string, error) {
	return "", errors.UnauthorizedError("invalid session", nil)
}

type getTenantErrAdminService struct {
	stubAdminTenantService
}

func (s *getTenantErrAdminService) GetTenantByID(context.Context, string) (*dtos.AdminTenantResponse, error) {
	return nil, errors.NotFoundError("Tenant", nil)
}

func TestAdminHandler_remainingPaths(t *testing.T) {
	t.Run("ListAuditLogs success with date filters", func(t *testing.T) {
		h := ProvideAdminHandler(stubAdminAuditService{}, validator.New())
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/admin/audit-logs?from=2024-01-01T00:00:00Z&to=2024-12-31T23:59:59Z&tenant_id="+testTenantID, nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		require.NoError(t, h.ListAuditLogs(c))
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), constants.AuditActionAuthLogin)
	})

	t.Run("PatchIdentityProvider missing tenantId", func(t *testing.T) {
		h := ProvideAdminHandler(stubAdminListService{}, validator.New())
		c, rec := newJSONContext(http.MethodPatch, "/", `{"enabled":true}`)
		c.SetParamNames("tenantId", "provider")
		c.SetParamValues("", "google")
		require.NoError(t, h.PatchIdentityProvider(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("GetTenant service error", func(t *testing.T) {
		h := ProvideAdminHandler(&getTenantErrAdminService{}, validator.New())
		c, rec := newJSONContext(http.MethodGet, "/", "")
		c.SetParamNames("tenantId")
		c.SetParamValues(testTenantID)
		require.NoError(t, h.GetTenant(c))
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("UpdateClient validation error", func(t *testing.T) {
		h := ProvideAdminHandler(stubAdminClientService{}, validator.New())
		c, rec := newJSONContext(http.MethodPatch, "/", `{"redirect_uris":["not-a-url"]}`)
		c.SetParamNames("clientId")
		c.SetParamValues(testClientPK)
		require.NoError(t, h.UpdateClient(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("UpdateClient service error", func(t *testing.T) {
		h := ProvideAdminHandler(errAdminService{err: errors.InternalError("update failed", nil)}, validator.New())
		c, rec := newJSONContext(http.MethodPatch, "/", `{"name":"Renamed"}`)
		c.SetParamNames("clientId")
		c.SetParamValues(testClientPK)
		require.NoError(t, h.UpdateClient(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("AddMember bind error", func(t *testing.T) {
		h := ProvideAdminHandler(stubAdminListService{}, validator.New())
		c, rec := newJSONContext(http.MethodPost, "/", `{`)
		c.SetParamNames("tenantId")
		c.SetParamValues(testTenantID)
		require.NoError(t, h.AddMember(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("DeleteTenant service error", func(t *testing.T) {
		h := ProvideAdminHandler(errAdminService{err: errors.InternalError("delete failed", nil)}, validator.New())
		c, rec := newJSONContext(http.MethodDelete, "/", "")
		c.SetParamNames("tenantId")
		c.SetParamValues(testTenantID)
		require.NoError(t, h.DeleteTenant(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("GetClient service error", func(t *testing.T) {
		h := ProvideAdminHandler(errAdminService{err: errors.InternalError("get failed", nil)}, validator.New())
		c, rec := newJSONContext(http.MethodGet, "/", "")
		c.SetParamNames("clientId")
		c.SetParamValues(testClientPK)
		require.NoError(t, h.GetClient(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("ForceLogoutUser missing userId", func(t *testing.T) {
		h := ProvideAdminHandler(stubAdminListService{}, validator.New())
		c, rec := newJSONContext(http.MethodPost, "/", "")
		c.SetParamNames("userId")
		c.SetParamValues("")
		c.Set(auth.EchoContextUserIDKey, "actor")
		require.NoError(t, h.ForceLogoutUser(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("CreateTenant service error", func(t *testing.T) {
		h := ProvideAdminHandler(errAdminService{err: errors.InternalError("create failed", nil)}, validator.New())
		c, rec := newJSONContext(http.MethodPost, "/", `{"name":"Acme"}`)
		require.NoError(t, h.CreateTenant(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("UpdateTenant validation error", func(t *testing.T) {
		h := ProvideAdminHandler(stubAdminTenantService{}, validator.New())
		c, rec := newJSONContext(http.MethodPatch, "/", `{"name":""}`)
		c.SetParamNames("tenantId")
		c.SetParamValues(testTenantID)
		require.NoError(t, h.UpdateTenant(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("ListTenantMembers service error", func(t *testing.T) {
		h := ProvideAdminHandler(errAdminService{err: errors.InternalError("list failed", nil)}, validator.New())
		c, rec := newJSONContext(http.MethodGet, "/", "")
		c.SetParamNames("tenantId")
		c.SetParamValues(testTenantID)
		require.NoError(t, h.ListTenantMembers(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("GetUser service error", func(t *testing.T) {
		h := ProvideAdminHandler(&stubAdminUserService{}, validator.New())
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("userId")
		c.SetParamValues("missing")
		require.NoError(t, h.GetUser(c))
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("DisableUser service error", func(t *testing.T) {
		h := ProvideAdminHandler(&stubAdminUserService{}, validator.New())
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("userId")
		c.SetParamValues("missing")
		c.Set(auth.EchoContextUserIDKey, "actor")
		require.NoError(t, h.DisableUser(c))
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestTenantIdentityAdminHandler_remainingPaths(t *testing.T) {
	cfg := handlerTestConfig()
	cfg.AdminAPIKey = "secret"

	t.Run("missing tenantId", func(t *testing.T) {
		h := ProvideTenantIdentityAdminHandler(cfg, stubAdminListService{})
		e := echo.New()
		req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"enabled":true}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("X-Admin-API-Key", "secret")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("tenantId", "provider")
		c.SetParamValues("", "google")
		require.NoError(t, h.PatchIdentityProvider(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("missing provider", func(t *testing.T) {
		h := ProvideTenantIdentityAdminHandler(cfg, stubAdminListService{})
		e := echo.New()
		req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"enabled":true}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("X-Admin-API-Key", "secret")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("tenantId", "provider")
		c.SetParamValues(testTenantID, "")
		require.NoError(t, h.PatchIdentityProvider(c))
		require.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("service error", func(t *testing.T) {
		h := ProvideTenantIdentityAdminHandler(cfg, errAdminService{err: errors.InternalError("configure failed", nil)})
		e := echo.New()
		req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"enabled":true}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Set("X-Admin-API-Key", "secret")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("tenantId", "provider")
		c.SetParamValues(testTenantID, "google")
		require.NoError(t, h.PatchIdentityProvider(c))
		require.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestActorUserIDFromContext_success(t *testing.T) {
	e := echo.New()
	c := e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	c.Set(auth.EchoContextUserIDKey, testUserID)
	uid, err := actorUserIDFromContext(c)
	require.NoError(t, err)
	require.Equal(t, testUserID, uid)
}

func TestLoginHistoryFiltersFromQuery_invalidTo(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/?to=not-a-date", nil)
	c := e.NewContext(req, httptest.NewRecorder())
	_, err := loginHistoryFiltersFromQuery(c)
	require.Error(t, err)
}

func TestOIDCHandler_Token_basicAuthOnlyClientID(t *testing.T) {
	h := newOIDCHandler(t, nil, nil, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/token", strings.NewReader("grant_type=authorization_code&code=c1"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("only-client-id:")))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	require.NoError(t, h.Token(c))
	require.Equal(t, http.StatusOK, rec.Code)
}
