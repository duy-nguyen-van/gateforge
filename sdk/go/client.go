package gateforge

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

// Option configures a Client.
type Option func(*Client)

// TokenProvider returns a bearer access token for authenticated requests.
type TokenProvider func(ctx context.Context) (string, error)

// Client is the public GateForge IAM SDK entry point.
type Client struct {
	baseURL       string
	httpClient    *http.Client
	userAgent     string
	tokenProvider TokenProvider

	api *openapi.APIClient
}

// NewClient constructs a Client for the given GateForge IAM base URL.
func NewClient(baseURL string, opts ...Option) *Client {
	c := &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{Timeout: 30 * time.Second},
		userAgent:  "gateforge-go-sdk/0.1.0",
	}
	for _, opt := range opts {
		opt(c)
	}
	c.rebuildAPI()
	return c
}

func (c *Client) rebuildAPI() {
	cfg := openapi.NewConfiguration()
	cfg.Servers = openapi.ServerConfigurations{{URL: c.baseURL}}
	cfg.HTTPClient = c.httpClient
	cfg.UserAgent = c.userAgent
	cfg.AddDefaultHeader("User-Agent", c.userAgent)
	c.api = openapi.NewAPIClient(cfg)
}

// BaseURL returns the configured API base URL (without a trailing slash).
func (c *Client) BaseURL() string { return c.baseURL }

// HTTPClient returns the underlying HTTP client.
func (c *Client) HTTPClient() *http.Client { return c.httpClient }

// UserAgent returns the User-Agent header value.
func (c *Client) UserAgent() string { return c.userAgent }

// API exposes the generated OpenAPI service groups.
func (c *Client) API() *openapi.APIClient { return c.api }

// WithHTTPClient overrides the default HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		if hc != nil {
			c.httpClient = hc
		}
	}
}

// WithTokenProvider sets a bearer token provider for authenticated calls.
func WithTokenProvider(tp TokenProvider) Option {
	return func(c *Client) {
		c.tokenProvider = tp
	}
}

// WithUserAgent overrides the default User-Agent string.
func WithUserAgent(ua string) Option {
	return func(c *Client) {
		if ua != "" {
			c.userAgent = ua
		}
	}
}

// withAuth attaches a bearer token from the token provider when configured.
func (c *Client) withAuth(ctx context.Context) context.Context {
	if c.tokenProvider == nil {
		return ctx
	}
	token, err := c.tokenProvider(ctx)
	if err != nil || token == "" {
		return ctx
	}
	return context.WithValue(ctx, openapi.ContextAccessToken, token)
}

// GetHealth calls GET /api/v1/.
func (c *Client) GetHealth(ctx context.Context) (*openapi.HealthEnvelope, *http.Response, error) {
	return c.api.HealthAPI.GetHealth(ctx).Execute()
}

// GetMe calls GET /api/v1/me.
func (c *Client) GetMe(ctx context.Context) (*openapi.UserResponseEnvelope, *http.Response, error) {
	return c.api.UsersAPI.GetMe(c.withAuth(ctx)).Execute()
}

// Login calls POST /api/v1/login.
func (c *Client) Login(ctx context.Context, body openapi.LoginRequest) (*openapi.LoginResultEnvelope, *http.Response, error) {
	return c.api.AuthAPI.Login(ctx).LoginRequest(body).Execute()
}

// RefreshToken calls POST /api/v1/refresh.
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*openapi.LoginResponseEnvelope, *http.Response, error) {
	return c.api.AuthAPI.RefreshToken(ctx).RefreshTokenRequest(openapi.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}).Execute()
}

// Logout calls POST /api/v1/logout.
func (c *Client) Logout(ctx context.Context) (*openapi.MetaOnlyEnvelope, *http.Response, error) {
	return c.api.AuthAPI.Logout(c.withAuth(ctx)).Execute()
}

// GetOpenIDConfiguration fetches OIDC discovery document.
func (c *Client) GetOpenIDConfiguration(ctx context.Context) (*openapi.OpenIDConfigurationResponse, *http.Response, error) {
	return c.api.OIDCAPI.GetOpenIdConfiguration(ctx).Execute()
}

// CreateToken exchanges an authorization code or refresh token at POST /token.
func (c *Client) CreateToken(ctx context.Context, form url.Values) (*openapi.OIDCTokenResponse, *http.Response, error) {
	req := c.api.OIDCAPI.CreateToken(ctx)
	if v := form.Get("grant_type"); v != "" {
		req = req.GrantType(v)
	}
	if v := form.Get("code"); v != "" {
		req = req.Code(v)
	}
	if v := form.Get("redirect_uri"); v != "" {
		req = req.RedirectUri(v)
	}
	if v := form.Get("client_id"); v != "" {
		req = req.ClientId(v)
	}
	if v := form.Get("client_secret"); v != "" {
		req = req.ClientSecret(v)
	}
	if v := form.Get("code_verifier"); v != "" {
		req = req.CodeVerifier(v)
	}
	if v := form.Get("refresh_token"); v != "" {
		req = req.RefreshToken(v)
	}
	if v := form.Get("scope"); v != "" {
		req = req.Scope(v)
	}
	return req.Execute()
}

// ExchangeAuthorizationCode is a convenience wrapper for the authorization_code grant.
func (c *Client) ExchangeAuthorizationCode(ctx context.Context, clientID, redirectURI, code, codeVerifier, clientSecret string) (*openapi.OIDCTokenResponse, *http.Response, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("client_id", clientID)
	form.Set("redirect_uri", redirectURI)
	form.Set("code", code)
	if codeVerifier != "" {
		form.Set("code_verifier", codeVerifier)
	}
	if clientSecret != "" {
		form.Set("client_secret", clientSecret)
	}
	return c.CreateToken(ctx, form)
}

// RefreshOIDCToken refreshes OIDC tokens via the refresh_token grant.
func (c *Client) RefreshOIDCToken(ctx context.Context, clientID, refreshToken, clientSecret string) (*openapi.OIDCTokenResponse, *http.Response, error) {
	form := url.Values{}
	form.Set("grant_type", "refresh_token")
	form.Set("client_id", clientID)
	form.Set("refresh_token", refreshToken)
	if clientSecret != "" {
		form.Set("client_secret", clientSecret)
	}
	return c.CreateToken(ctx, form)
}

// GetUserInfo calls GET /userinfo with an OIDC access token.
func (c *Client) GetUserInfo(ctx context.Context, accessToken string) (*openapi.UserInfoResponse, *http.Response, error) {
	ctx = context.WithValue(ctx, openapi.ContextAccessToken, accessToken)
	return c.api.OIDCAPI.GetUserInfo(ctx).Execute()
}

// AsOpenAPIError unwraps a generated client error when present.
func AsOpenAPIError(err error) (*openapi.GenericOpenAPIError, bool) {
	if err == nil {
		return nil, false
	}
	var oe *openapi.GenericOpenAPIError
	if errors.As(err, &oe) {
		return oe, true
	}
	return nil, false
}
