package monitoring

import (
	"net/http"
	"strings"

	"github.com/getsentry/sentry-go"
)

var redactedHeaderKeys = []string{
	"Authorization",
	"Cookie",
	"Set-Cookie",
	"X-CSRF-Token",
	"X-Admin-API-Key",
}

// SetScopeData attaches data to a Sentry scope for error events.
// sentry-go v0.46 removed SetExtra; strings use tags, other values use contexts.
func SetScopeData(scope *sentry.Scope, key string, value interface{}) {
	if value == nil {
		return
	}
	if s, ok := value.(string); ok {
		scope.SetTag(key, s)
		return
	}
	scope.SetContext(key, sentry.Context{"value": value})
}

// RedactedHeaders copies h without auth, cookie, CSRF, or admin-key values.
func RedactedHeaders(h http.Header) http.Header {
	if h == nil {
		return nil
	}
	out := h.Clone()
	for _, key := range redactedHeaderKeys {
		if len(out.Values(key)) == 0 {
			continue
		}
		out.Del(key)
		out.Set(key, "[redacted]")
	}
	return out
}

// RedactSentryRequest strips secrets from sentryecho's request payload before send.
func RedactSentryRequest(req *sentry.Request) {
	if req == nil {
		return
	}
	if req.Headers != nil {
		h := make(http.Header, len(req.Headers))
		for k, v := range req.Headers {
			h.Set(k, v)
		}
		redacted := RedactedHeaders(h)
		out := make(map[string]string, len(redacted))
		for k := range redacted {
			out[k] = redacted.Get(k)
		}
		req.Headers = out
	}
	if req.Cookies != "" {
		req.Cookies = "[redacted]"
	}
	if skipSentryRequestBody(req.URL) {
		req.Data = ""
	}
}

func skipSentryRequestBody(rawURL string) bool {
	path := rawURL
	if i := strings.IndexByte(path, '?'); i >= 0 {
		path = path[:i]
	}
	if i := strings.Index(path, "://"); i >= 0 {
		rest := path[i+3:]
		if slash := strings.IndexByte(rest, '/'); slash >= 0 {
			path = rest[slash:]
		} else {
			path = "/"
		}
	}
	return SkipSensitiveRequestPath(path)
}

// SkipSensitiveRequestPath reports whether the HTTP path may carry secrets in the body.
func SkipSensitiveRequestPath(path string) bool {
	switch path {
	case "/token", "/authorize", "/introspect", "/userinfo":
		return true
	}
	if strings.HasPrefix(path, "/oidc/") {
		return true
	}
	if !strings.HasPrefix(path, "/api/v1/") {
		return false
	}
	rest := strings.TrimPrefix(path, "/api/v1")
	switch {
	case rest == "/login", strings.HasPrefix(rest, "/login/"):
		return true
	case rest == "/register", rest == "/refresh", rest == "/logout":
		return true
	case strings.HasPrefix(rest, "/mfa/"):
		return true
	case strings.HasPrefix(rest, "/webauthn/"):
		return true
	case rest == "/tenants/select", rest == "/tenants/switch":
		return true
	default:
		return false
	}
}
