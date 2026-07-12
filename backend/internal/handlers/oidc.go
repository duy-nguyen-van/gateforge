package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"
	"github.com/gateforge-iam/gateforge-iam/internal/services"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// OIDCHandler serves OIDC/OAuth2 Phase 2 endpoints (authorize, token, userinfo, discovery, JWKS).
type OIDCHandler struct {
	BaseHandler
	cfg            *config.Config
	oidc           *auth.OIDCSigner
	oidcService    services.OIDCService
	sessionService services.SessionService
	userService    services.UserService
	federationSvc  services.FederationService
	auditService   services.AuditService
	clientRepo     repositories.ClientRepository
	validator      *validator.Validate
}

// ProvideOIDCHandler wires OIDC HTTP handlers.
func ProvideOIDCHandler(
	cfg *config.Config,
	oidc *auth.OIDCSigner,
	oidcService services.OIDCService,
	sessionService services.SessionService,
	userService services.UserService,
	federationSvc services.FederationService,
	auditService services.AuditService,
	clientRepo repositories.ClientRepository,
	validator *validator.Validate,
) *OIDCHandler {
	return &OIDCHandler{
		BaseHandler:    *NewBaseHandler(),
		cfg:            cfg,
		oidc:           oidc,
		oidcService:    oidcService,
		sessionService: sessionService,
		userService:    userService,
		federationSvc:  federationSvc,
		auditService:   auditService,
		clientRepo:     clientRepo,
		validator:      validator,
	}
}

// LoginOIDC godoc
// @Summary Login (OIDC browser flow) and continue /authorize
// @Tags OIDC
// @Accept json
// @Produce json
// @Param body body dtos.LoginRequest true "Credentials + return_to"
// @Success 302 {string} string "Redirect to return_to with session cookie set"
// @Router /oidc/login [post]
func (h *OIDCHandler) Login(c echo.Context) error {
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

	returnTo := strings.TrimSpace(req.ReturnTo)
	if returnTo == "" {
		return h.HandleError(c, errors.ValidationError("return_to is required for OIDC login", nil))
	}
	if !isSafeAuthorizeReturnURL(c, h.cfg.AppBaseURL, returnTo) {
		return h.HandleError(c, errors.ValidationError("Invalid return_to (must be /authorize on this host)", nil))
	}

	u, err := h.userService.AuthenticateUser(c.Request().Context(), &req)
	if err != nil {
		return h.HandleError(c, err)
	}

	tenantID := h.resolveTenantFromReturnTo(c, returnTo)
	if tenantID == "" {
		tenantID = h.cfg.DefaultTenantID
	}

	sid, ttl, err := h.sessionService.Create(c.Request().Context(), u.ID, tenantID, c.RealIP(), c.Request().UserAgent(), req.RememberMe)
	if err != nil {
		return h.HandleError(c, err)
	}
	cookie := &http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.cfg.AppEnv == config.EnvironmentProduction,
		MaxAge:   int(ttl.Seconds()),
	}
	c.SetCookie(cookie)
	h.auditService.Record(c.Request().Context(), domains.AuditRecordParams{
		Action:       constants.AuditActionOIDCLogin,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      u.ID,
		TenantID:     tenantID,
		ResourceType: constants.AuditResourceTypeUser,
		ResourceID:   u.ID,
	})
	return c.Redirect(http.StatusFound, returnTo)
}

// FederationOAuthStart starts upstream OAuth for a registered provider (e.g. google, microsoft).
// @Summary Start federated sign-in (OIDC browser flow)
// @Tags OIDC
// @Param provider path string true "Provider id (e.g. google)"
// @Param return_to query string true "URL to return to after login (must be /authorize on this app)"
// @Success 302 {string} string "Redirect to identity provider"
// @Router /oidc/federation/{provider}/start [get]
func (h *OIDCHandler) FederationOAuthStart(c echo.Context) error {
	providerID := c.Param("provider")
	returnTo := strings.TrimSpace(c.QueryParam("return_to"))
	if returnTo == "" {
		return h.HandleError(c, errors.ValidationError("return_to is required", nil))
	}
	if !isSafeFederationReturnURL(c, h.cfg, returnTo) {
		return h.HandleError(c, errors.ValidationError("Invalid return_to (must be /authorize or SPA federation complete on this app)", nil))
	}

	tenantID := h.resolveTenantFromReturnTo(c, returnTo)
	if tenantID == "" {
		tenantID = h.cfg.DefaultTenantID
	}
	url, err := h.federationSvc.BuildAuthorizeRedirectURL(c.Request().Context(), providerID, returnTo, tenantID)
	if err != nil {
		return h.HandleError(c, err)
	}
	return c.Redirect(http.StatusFound, url)
}

// FederationOAuthCallback completes upstream OAuth and sets iam_session.
// @Summary Federated OAuth callback
// @Tags OIDC
// @Param provider path string true "Provider id (e.g. google)"
// @Param code query string true "Authorization code"
// @Param state query string true "State"
// @Success 302 {string} string "Redirect to stored return_to"
// @Router /oidc/federation/{provider}/callback [get]
func (h *OIDCHandler) FederationOAuthCallback(c echo.Context) error {
	providerID := c.Param("provider")

	if errParam := strings.TrimSpace(c.QueryParam("error")); errParam != "" {
		desc := strings.TrimSpace(c.QueryParam("error_description"))
		msg := "Sign-in was cancelled or denied"
		if desc != "" {
			msg = desc
		}
		return h.HandleError(c, errors.ValidationError(msg, nil))
	}

	code := c.QueryParam("code")
	state := c.QueryParam("state")

	u, tenantID, returnTo, err := h.federationSvc.CompleteOAuthLogin(c.Request().Context(), providerID, code, state)
	if err != nil {
		return h.HandleError(c, err)
	}

	sid, ttl, err := h.sessionService.Create(c.Request().Context(), u.ID, tenantID, c.RealIP(), c.Request().UserAgent(), false)
	if err != nil {
		return h.HandleError(c, err)
	}
	cookie := &http.Cookie{
		Name:     constants.SessionCookieName,
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   h.cfg.AppEnv == config.EnvironmentProduction,
		MaxAge:   int(ttl.Seconds()),
	}
	c.SetCookie(cookie)
	return c.Redirect(http.StatusFound, returnTo)
}

func (h *OIDCHandler) resolveTenantFromReturnTo(c echo.Context, returnTo string) string {
	u, err := url.Parse(returnTo)
	if err != nil {
		return ""
	}
	clientID := strings.TrimSpace(u.Query().Get("client_id"))
	if clientID == "" {
		return ""
	}
	client, err := h.clientRepo.GetByClientID(c.Request().Context(), clientID)
	if err != nil || client == nil {
		return ""
	}
	return client.TenantID
}

const federationCompletePath = "/login/federation/complete"

func isSafeFederationReturnURL(c echo.Context, cfg *config.Config, raw string) bool {
	if isSafeAuthorizeReturnURL(c, cfg.AppBaseURL, raw) {
		return true
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return false
	}
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	if u.Path != federationCompletePath {
		return false
	}
	if !u.IsAbs() {
		return true
	}
	return federationReturnHostAllowed(c, cfg, u.Host)
}

func federationReturnHostAllowed(c echo.Context, cfg *config.Config, host string) bool {
	if host == c.Request().Host {
		return true
	}
	for _, base := range []string{cfg.AppBaseURL, cfg.OIDCLoginPageURL} {
		base = strings.TrimSuffix(strings.TrimSpace(base), "/")
		if base == "" {
			continue
		}
		baseURL, err := url.Parse(base)
		if err != nil {
			continue
		}
		if baseURL.Host == host {
			return true
		}
	}
	return false
}

func isSafeAuthorizeReturnURL(c echo.Context, appBaseURL, raw string) bool {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return false
	}
	u, err := url.Parse(raw)
	if err != nil {
		return false
	}
	if !u.IsAbs() {
		return strings.HasPrefix(u.Path, "/authorize")
	}
	base := strings.TrimSuffix(appBaseURL, "/")
	if base != "" {
		baseURL, err := url.Parse(base)
		if err != nil {
			return false
		}
		return u.Scheme == baseURL.Scheme && u.Host == baseURL.Host && strings.HasPrefix(u.Path, "/authorize")
	}
	return u.Host == c.Request().Host && strings.HasPrefix(u.Path, "/authorize")
}

// Authorize handles GET /authorize (authorization code + PKCE). Authenticates via iam_session cookie after POST /oidc/login. If unauthenticated, redirects to the login page with return_to=/authorize...
// @Summary Authorize endpoint (authorization code + PKCE)
// @Tags OIDC
// @Accept json
// @Produce json
// @Param Cookie header string false "Session cookie from browser login (iam_session)"
// @Param state query string false "State parameter"
// @Param response_type query string false "Response type (code)"
// @Param client_id query string false "Client ID"
// @Param redirect_uri query string false "Redirect URI"
// @Param scope query string false "Scope"
// @Param nonce query string false "Nonce"
// @Param code_challenge query string false "Code challenge"
// @Param code_challenge_method query string false "Code challenge method"
// @Success 302 {string} string "Redirect to client or login"
// @Router /authorize [get]
func (h *OIDCHandler) Authorize(c echo.Context) error {
	userID, ok := h.resolveUserIDForAuthorize(c)
	if !ok {
		return c.Redirect(http.StatusFound, h.loginRedirectURL(c))
	}

	authzQuery := dtos.NewAuthorizeQueryFromURLValues(c.Request().URL.Query())
	loc, oerr := h.oidcService.Authorize(c.Request().Context(), userID, &authzQuery)
	if oerr != nil {
		if oerr.RedirectTo != "" {
			return c.Redirect(http.StatusFound, oerr.RedirectTo)
		}
		code := oerr.Code
		if code == "" {
			code = constants.OAuthInvalidRequest
		}
		status := http.StatusBadRequest
		if oerr.Code == constants.OAuthServerError {
			status = http.StatusInternalServerError
		}
		return c.JSON(status, map[string]string{
			"error":             code,
			"error_description": oerr.Description,
		})
	}
	return c.Redirect(http.StatusFound, loc)
}

func (h *OIDCHandler) resolveUserIDForAuthorize(c echo.Context) (string, bool) {
	if ck, err := c.Cookie(constants.SessionCookieName); err == nil && ck != nil && ck.Value != "" {
		if uid, err := h.sessionService.GetUserID(c.Request().Context(), ck.Value); err == nil && uid != "" {
			return uid, true
		}
	}
	return "", false
}

func (h *OIDCHandler) loginRedirectURL(c echo.Context) string {
	base := strings.TrimSuffix(h.oidcLoginPageURL(), "/")
	q := url.Values{}
	q.Set("return_to", h.absoluteAuthorizeURL(c))
	return base + "?" + q.Encode()
}

func (h *OIDCHandler) oidcLoginPageURL() string {
	if h.cfg.OIDCLoginPageURL != "" {
		return h.cfg.OIDCLoginPageURL
	}
	if b := strings.TrimSuffix(h.cfg.AppBaseURL, "/"); b != "" {
		return b + "/login"
	}
	return "/login"
}

func (h *OIDCHandler) absoluteAuthorizeURL(c echo.Context) string {
	issuer := strings.TrimSuffix(h.oidcService.OpenIDIssuer(), "/")
	if u, err := url.Parse(issuer); err == nil && u.Scheme != "" && u.Host != "" {
		return issuer + c.Request().URL.RequestURI()
	}
	scheme := c.Scheme()
	if c.Request().Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return scheme + "://" + c.Request().Host + c.Request().URL.RequestURI()
}

// Token handles POST /token (authorization_code grant). Accepts application/x-www-form-urlencoded or application/json.
// @Summary Token endpoint (authorization_code grant)
// @Tags OIDC
// @Accept json
// @Accept x-www-form-urlencoded
// @Produce json
// @Success 200 {object} domains.OIDCTokenResponse
// @Router /token [post]
func (h *OIDCHandler) Token(c echo.Context) error {
	var form url.Values
	ct := c.Request().Header.Get("Content-Type")
	if strings.HasPrefix(ct, "application/json") {
		var body struct {
			GrantType    string `json:"grant_type"`
			Code         string `json:"code"`
			RedirectURI  string `json:"redirect_uri"`
			ClientID     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
			CodeVerifier string `json:"code_verifier"`
		}
		if err := c.Bind(&body); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": constants.OAuthInvalidRequest, "error_description": "invalid JSON body"})
		}
		form = url.Values{}
		form.Set("grant_type", body.GrantType)
		form.Set("code", body.Code)
		form.Set("redirect_uri", body.RedirectURI)
		form.Set("client_id", body.ClientID)
		form.Set("client_secret", body.ClientSecret)
		form.Set("code_verifier", body.CodeVerifier)
	} else {
		if err := c.Request().ParseForm(); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": constants.OAuthInvalidRequest, "error_description": "invalid form body"})
		}
		form = c.Request().Form
	}

	clientID, clientSecret := "", ""
	if u, p, ok := c.Request().BasicAuth(); ok {
		clientID, clientSecret = u, p
	}
	if clientID == "" {
		clientID = form.Get("client_id")
	}
	if clientSecret == "" {
		clientSecret = form.Get("client_secret")
	}

	out, terr := h.oidcService.AuthorizationCodeToken(c.Request().Context(), "", clientID, clientSecret, form)
	if terr != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": terr.Code, "error_description": terr.Error()})
	}
	return c.JSON(http.StatusOK, out)
}

// UserInfo handles GET /userinfo (Bearer access token from token endpoint, RS256).
// @Summary Userinfo endpoint (Bearer access token from token endpoint, RS256)
// @Tags OIDC
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization Bearer token"
// @Success 200 {object} map[string]interface{}
// @Router /userinfo [get]
func (h *OIDCHandler) UserInfo(c echo.Context) error {
	authz := c.Request().Header.Get("Authorization")
	if authz == "" || !strings.HasPrefix(authz, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": constants.OAuthInvalidToken, "error_description": "Bearer access token required"})
	}
	raw := strings.TrimPrefix(authz, "Bearer ")
	out, terr := h.oidcService.UserInfo(c.Request().Context(), raw)
	if terr != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": terr.Code, "error_description": terr.Error()})
	}
	return c.JSON(http.StatusOK, out)
}

// JWKS handles GET /.well-known/jwks.json.
// @Summary JWKS endpoint (JSON Web Key Set)
// @Tags OIDC
// @Accept json
// @Produce json
// @Success 200 {string} string "JWKS JSON"
// @Router /.well-known/jwks.json [get]
func (h *OIDCHandler) JWKS(c echo.Context) error {
	b, err := h.oidc.MarshalJWKS()
	if err != nil {
		return h.InternalErrorResponse(c, "Failed to marshal JWKS", err)
	}
	return c.Blob(http.StatusOK, "application/json", b)
}

// OpenIDConfiguration handles GET /.well-known/openid-configuration.
// @Summary OpenID Configuration endpoint
// @Tags OIDC
// @Accept json
// @Produce json
// @Success 200 {object} dtos.OpenIDConfigurationResponse
// @Router /.well-known/openid-configuration [get]
func (h *OIDCHandler) OpenIDConfiguration(c echo.Context) error {
	base := strings.TrimSuffix(h.oidcService.OpenIDIssuer(), "/")
	u, err := url.Parse(base)
	if err != nil || u.Scheme == "" {
		base = "http://localhost:3000"
	}

	response := dtos.OpenIDConfigurationResponse{
		Issuer:                            base,
		AuthorizationEndpoint:             base + "/authorize",
		TokenEndpoint:                     base + "/token",
		UserinfoEndpoint:                  base + "/userinfo",
		JWKSURI:                           base + "/.well-known/jwks.json",
		ResponseTypesSupported:            []string{"code"},
		SubjectTypesSupported:             []string{"public"},
		IDTokenSigningAlgValuesSupported:  []string{"RS256"},
		ScopesSupported:                   []string{"openid", "email", "profile"},
		TokenEndpointAuthMethodsSupported: []string{"none", "client_secret_post", "client_secret_basic"},
		CodeChallengeMethodsSupported:     []string{"S256"},
		GrantTypesSupported:               []string{"authorization_code"},
	}
	return c.JSON(http.StatusOK, response)
}
