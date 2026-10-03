package monitoring

import (
	"net/http"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedactedHeaders(t *testing.T) {
	h := http.Header{}
	h.Set("Authorization", "Bearer secret")
	h.Set("Cookie", "session=x")
	h.Set("X-CSRF-Token", "tok")
	h.Set("X-Admin-API-Key", "admin-key")
	h.Set("Content-Type", "application/json")

	out := RedactedHeaders(h)
	require.NotNil(t, out)
	assert.Equal(t, "[redacted]", out.Get("Authorization"))
	assert.Equal(t, "[redacted]", out.Get("Cookie"))
	assert.Equal(t, "[redacted]", out.Get("X-CSRF-Token"))
	assert.Equal(t, "[redacted]", out.Get("X-Admin-API-Key"))
	assert.Equal(t, "application/json", out.Get("Content-Type"))
	assert.Equal(t, "Bearer secret", h.Get("Authorization"))
}

func TestRedactedHeadersNil(t *testing.T) {
	assert.Nil(t, RedactedHeaders(nil))
}

func TestSkipSensitiveRequestPath(t *testing.T) {
	t.Parallel()

	sensitive := []string{
		"/token",
		"/authorize",
		"/introspect",
		"/userinfo",
		"/oidc/login",
		"/oidc/federation/google/callback",
		"/api/v1/login",
		"/api/v1/register",
		"/api/v1/refresh",
		"/api/v1/logout",
		"/api/v1/mfa/totp/verify",
		"/api/v1/webauthn/login/finish",
		"/api/v1/tenants/select",
		"/api/v1/tenants/switch",
	}
	for _, path := range sensitive {
		assert.True(t, SkipSensitiveRequestPath(path), path)
	}

	assert.False(t, SkipSensitiveRequestPath("/api/v1/me"))
	assert.False(t, SkipSensitiveRequestPath("/health"))
	assert.False(t, SkipSensitiveRequestPath("/api/v1/admin/users"))
}

func TestRedactSentryRequest(t *testing.T) {
	req := &sentry.Request{
		URL:     "https://iam.example.com/api/v1/me",
		Data:    `{"ok":true}`,
		Cookies: "iam_session=secret",
		Headers: map[string]string{
			"authorization":   "Bearer secret",
			"Cookie":          "iam_session=x",
			"X-CSRF-Token":    "tok",
			"X-Admin-API-Key": "admin-key",
			"Content-Type":    "application/json",
		},
	}
	RedactSentryRequest(req)
	assert.Equal(t, "[redacted]", req.Headers["Authorization"])
	assert.Equal(t, "[redacted]", req.Headers["Cookie"])
	assert.Equal(t, "[redacted]", req.Headers["X-Csrf-Token"])
	assert.Equal(t, "[redacted]", req.Headers["X-Admin-Api-Key"])
	assert.Equal(t, "application/json", req.Headers["Content-Type"])
	assert.Equal(t, "[redacted]", req.Cookies)
	assert.Equal(t, `{"ok":true}`, req.Data)
}

func TestRedactSentryRequest_LoginStripsBody(t *testing.T) {
	req := &sentry.Request{
		URL:  "https://iam.example.com/api/v1/login",
		Data: `{"password":"secret"}`,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	RedactSentryRequest(req)
	assert.Empty(t, req.Data)
}

func TestRedactSentryRequest_TokenStripsBody(t *testing.T) {
	req := &sentry.Request{
		URL:  "https://iam.example.com/token?grant_type=authorization_code",
		Data: `client_secret=secret`,
	}
	RedactSentryRequest(req)
	assert.Empty(t, req.Data)
}

func TestRedactSentryRequestNil(t *testing.T) {
	RedactSentryRequest(nil)
}

func TestSetScopeData(t *testing.T) {
	scope := sentry.NewScope()
	SetScopeData(scope, "k", nil)
	SetScopeData(scope, "k", "v")
}
