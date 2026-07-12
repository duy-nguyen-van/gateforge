package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func newOIDCHandler(t *testing.T, oidcSvc *stubOIDCService, fed *stubOIDCFederationService, clientRepo *stubClientRepo) *OIDCHandler {
	t.Helper()
	cfg := handlerTestConfig()
	signer, err := auth.ProvideOIDCSigner(cfg)
	require.NoError(t, err)
	if oidcSvc == nil {
		oidcSvc = &stubOIDCService{}
	}
	if fed == nil {
		fed = &stubOIDCFederationService{}
	}
	if clientRepo == nil {
		clientRepo = &stubClientRepo{}
	}
	return ProvideOIDCHandler(cfg, signer, oidcSvc, &stubAuthSessionService{}, &stubAuthUserService{}, fed, stubAuditService{}, clientRepo, validator.New())
}

func TestOIDCHandler_Login_success(t *testing.T) {
	h := newOIDCHandler(t, nil, nil, nil)
	body := `{"email":"user@example.com","password":"secretpass","return_to":"/authorize?client_id=app"}`
	c, rec := newJSONContext(http.MethodPost, "/oidc/login", body)

	err := h.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, rec.Code)
	require.Contains(t, rec.Header().Get("Location"), "/authorize")
	cookies := rec.Result().Cookies()
	require.NotEmpty(t, cookies)
	require.Equal(t, constants.SessionCookieName, cookies[0].Name)
}

func TestOIDCHandler_Login_missingReturnTo(t *testing.T) {
	h := newOIDCHandler(t, nil, nil, nil)
	c, rec := newJSONContext(http.MethodPost, "/oidc/login", `{"email":"user@example.com","password":"secretpass"}`)

	err := h.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestOIDCHandler_Login_invalidReturnTo(t *testing.T) {
	h := newOIDCHandler(t, nil, nil, nil)
	body := `{"email":"user@example.com","password":"secretpass","return_to":"https://evil.example/steal"}`
	c, rec := newJSONContext(http.MethodPost, "/oidc/login", body)

	err := h.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestOIDCHandler_FederationOAuthStart_success(t *testing.T) {
	h := newOIDCHandler(t, nil, &stubOIDCFederationService{redirectURL: "https://accounts.google.com/o/oauth2/auth"}, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/oidc/federation/google/start?return_to=/authorize?client_id=app", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("provider")
	c.SetParamValues("google")

	err := h.FederationOAuthStart(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, rec.Code)
	require.Contains(t, rec.Header().Get("Location"), "accounts.google.com")
}

func TestOIDCHandler_FederationOAuthStart_missingReturnTo(t *testing.T) {
	h := newOIDCHandler(t, nil, nil, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/oidc/federation/google/start", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("provider")
	c.SetParamValues("google")

	err := h.FederationOAuthStart(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestOIDCHandler_FederationOAuthStart_federationCompletePath(t *testing.T) {
	h := newOIDCHandler(t, nil, nil, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/oidc/federation/google/start?return_to=/login/federation/complete", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("provider")
	c.SetParamValues("google")

	err := h.FederationOAuthStart(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, rec.Code)
}

func TestOIDCHandler_FederationOAuthCallback_providerError(t *testing.T) {
	h := newOIDCHandler(t, nil, nil, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/oidc/federation/google/callback?error=access_denied&error_description=User+cancelled", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("provider")
	c.SetParamValues("google")

	err := h.FederationOAuthCallback(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), constants.ValidationError)
}

func TestOIDCHandler_FederationOAuthCallback_success(t *testing.T) {
	h := newOIDCHandler(t, nil, &stubOIDCFederationService{returnTo: "/login/federation/complete"}, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/oidc/federation/google/callback?code=abc&state=xyz", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("provider")
	c.SetParamValues("google")

	err := h.FederationOAuthCallback(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, rec.Code)
	require.Contains(t, rec.Header().Get("Location"), "/login/federation/complete")
}

func TestOIDCHandler_Authorize_redirectToLogin(t *testing.T) {
	h := newOIDCHandler(t, nil, nil, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/authorize?client_id=app&response_type=code", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Authorize(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, rec.Code)
	require.Contains(t, rec.Header().Get("Location"), "return_to=")
}

func TestOIDCHandler_Authorize_success(t *testing.T) {
	h := newOIDCHandler(t, &stubOIDCService{authorizeLoc: "http://client.example/callback?code=abc"}, nil, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/authorize?client_id=app&response_type=code", nil)
	req.AddCookie(&http.Cookie{Name: constants.SessionCookieName, Value: "sess-1"})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Authorize(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, rec.Code)
	require.Contains(t, rec.Header().Get("Location"), "client.example")
}

func TestOIDCHandler_Authorize_oauthRedirectError(t *testing.T) {
	h := newOIDCHandler(t, &stubOIDCService{
		authorizeErr: &domains.OAuthRedirectError{RedirectTo: "http://client.example/callback?error=invalid_scope"},
	}, nil, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/authorize", nil)
	req.AddCookie(&http.Cookie{Name: constants.SessionCookieName, Value: "sess-1"})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Authorize(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusFound, rec.Code)
}

func TestOIDCHandler_Authorize_oauthJSONError(t *testing.T) {
	h := newOIDCHandler(t, &stubOIDCService{
		authorizeErr: &domains.OAuthRedirectError{Code: constants.OAuthInvalidRequest, Description: "bad request"},
	}, nil, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/authorize", nil)
	req.AddCookie(&http.Cookie{Name: constants.SessionCookieName, Value: "sess-1"})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Authorize(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), constants.OAuthInvalidRequest)
}

func TestOIDCHandler_Authorize_serverError(t *testing.T) {
	h := newOIDCHandler(t, &stubOIDCService{
		authorizeErr: &domains.OAuthRedirectError{Code: constants.OAuthServerError, Description: "boom"},
	}, nil, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/authorize", nil)
	req.AddCookie(&http.Cookie{Name: constants.SessionCookieName, Value: "sess-1"})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Authorize(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestOIDCHandler_Token_jsonSuccess(t *testing.T) {
	h := newOIDCHandler(t, &stubOIDCService{tokenResp: &domains.OIDCTokenResponse{AccessToken: "at-1", TokenType: "Bearer"}}, nil, nil)
	body := `{"grant_type":"authorization_code","code":"code-1","client_id":"app","code_verifier":"verifier"}`
	c, rec := newJSONContext(http.MethodPost, "/token", body)

	err := h.Token(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "at-1")
}

func TestOIDCHandler_Token_formSuccess(t *testing.T) {
	h := newOIDCHandler(t, nil, nil, nil)
	e := echo.New()
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", "code-1")
	form.Set("client_id", "app")
	req := httptest.NewRequest(http.MethodPost, "/token", strings.NewReader(form.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.Token(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestOIDCHandler_Token_invalidJSON(t *testing.T) {
	h := newOIDCHandler(t, nil, nil, nil)
	c, rec := newJSONContext(http.MethodPost, "/token", `{`)

	err := h.Token(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestOIDCHandler_Token_error(t *testing.T) {
	h := newOIDCHandler(t, &stubOIDCService{
		tokenErr: &domains.OAuthTokenError{Code: constants.OAuthInvalidGrant, Description: "bad code"},
	}, nil, nil)
	c, rec := newJSONContext(http.MethodPost, "/token", `{"grant_type":"authorization_code","code":"bad"}`)

	err := h.Token(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), constants.OAuthInvalidGrant)
}

func TestOIDCHandler_UserInfo_missingBearer(t *testing.T) {
	h := newOIDCHandler(t, nil, nil, nil)
	c, rec := newJSONContext(http.MethodGet, "/userinfo", "")

	err := h.UserInfo(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestOIDCHandler_UserInfo_success(t *testing.T) {
	h := newOIDCHandler(t, &stubOIDCService{userInfo: map[string]any{"sub": testUserID}}, nil, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/userinfo", nil)
	req.Header.Set("Authorization", "Bearer access-token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.UserInfo(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), testUserID)
}

func TestOIDCHandler_UserInfo_invalidToken(t *testing.T) {
	h := newOIDCHandler(t, &stubOIDCService{
		userInfoErr: &domains.OAuthTokenError{Code: constants.OAuthInvalidToken, Description: "expired"},
	}, nil, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/userinfo", nil)
	req.Header.Set("Authorization", "Bearer bad-token")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.UserInfo(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestOIDCHandler_JWKS(t *testing.T) {
	h := newOIDCHandler(t, nil, nil, nil)
	c, rec := newJSONContext(http.MethodGet, "/.well-known/jwks.json", "")

	err := h.JWKS(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Header().Get(echo.HeaderContentType), "application/json")
	require.Contains(t, rec.Body.String(), `"keys"`)
}

func TestOIDCHandler_OpenIDConfiguration(t *testing.T) {
	h := newOIDCHandler(t, &stubOIDCService{issuer: "http://localhost:8080"}, nil, nil)
	c, rec := newJSONContext(http.MethodGet, "/.well-known/openid-configuration", "")

	err := h.OpenIDConfiguration(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"issuer":"http://localhost:8080"`)
	require.Contains(t, rec.Body.String(), "/authorize")
}

func TestOIDCHandler_resolveTenantFromReturnTo(t *testing.T) {
	h := newOIDCHandler(t, nil, nil, &stubClientRepo{
		client: &models.Client{BaseModel: models.NewBaseModel(), ClientID: "app", TenantID: testTenantID},
	})
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	tenantID := h.resolveTenantFromReturnTo(c, "/authorize?client_id=app")
	require.Equal(t, testTenantID, tenantID)
}

func TestIsSafeAuthorizeReturnURL(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.True(t, isSafeAuthorizeReturnURL(c, "http://localhost:8080", "/authorize?client_id=app"))
	require.True(t, isSafeAuthorizeReturnURL(c, "http://localhost:8080", "http://localhost:8080/authorize?client_id=app"))
	require.False(t, isSafeAuthorizeReturnURL(c, "http://localhost:8080", "http://evil.example/authorize"))
}

func TestIsSafeFederationReturnURL(t *testing.T) {
	cfg := handlerTestConfig()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	require.True(t, isSafeFederationReturnURL(c, cfg, "/login/federation/complete"))
	require.True(t, isSafeFederationReturnURL(c, cfg, "http://localhost:8080/login/federation/complete"))
	require.False(t, isSafeFederationReturnURL(c, cfg, "http://evil.example/login/federation/complete"))
}

func TestOIDCHandler_Login_authFailure(t *testing.T) {
	h := newOIDCHandler(t, nil, nil, nil)
	userSvc := &stubAuthUserService{authErr: errors.UnauthorizedError("bad creds", nil)}
	h.userService = userSvc
	body := `{"email":"user@example.com","password":"wrong","return_to":"/authorize?client_id=app"}`
	c, rec := newJSONContext(http.MethodPost, "/oidc/login", body)

	err := h.Login(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestOIDCHandler_FederationOAuthCallback_serviceError(t *testing.T) {
	h := newOIDCHandler(t, nil, &stubOIDCFederationService{completeErr: errors.ValidationError("bad state", nil)}, nil)
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/oidc/federation/google/callback?code=abc&state=xyz", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("provider")
	c.SetParamValues("google")

	err := h.FederationOAuthCallback(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}
