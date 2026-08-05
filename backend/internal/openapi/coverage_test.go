package openapicov_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// Public routes that must appear in the canonical OpenAPI contract.
var requiredPublicRoutes = []struct {
	method string
	path   string
}{
	{"get", "/.well-known/jwks.json"},
	{"get", "/.well-known/openid-configuration"},
	{"get", "/authorize"},
	{"post", "/oidc/login"},
	{"get", "/oidc/federation/{provider}/start"},
	{"get", "/oidc/federation/{provider}/callback"},
	{"post", "/token"},
	{"post", "/introspect"},
	{"get", "/userinfo"},
	{"get", "/api/v1/"},
	{"get", "/api/v1/health/database"},
	{"get", "/api/v1/health/metrics"},
	{"post", "/api/v1/register"},
	{"post", "/api/v1/login"},
	{"post", "/api/v1/login/session"},
	{"get", "/api/v1/federation/providers"},
	{"post", "/api/v1/refresh"},
	{"post", "/api/v1/webauthn/login/start"},
	{"post", "/api/v1/webauthn/login/finish"},
	{"post", "/api/v1/mfa/challenge/verify"},
	{"post", "/api/v1/logout"},
	{"get", "/api/v1/me"},
	{"patch", "/api/v1/me"},
	{"get", "/api/v1/me/tenants"},
	{"post", "/api/v1/tenants/select"},
	{"post", "/api/v1/tenants/switch"},
	{"get", "/api/v1/webauthn/credentials"},
	{"post", "/api/v1/webauthn/register/start"},
	{"post", "/api/v1/webauthn/register/finish"},
	{"post", "/api/v1/mfa/totp/setup"},
	{"post", "/api/v1/mfa/totp/verify"},
	{"post", "/api/v1/mfa/recovery-codes"},
	{"get", "/api/v1/admin/stats"},
	{"get", "/api/v1/admin/users"},
	{"get", "/api/v1/admin/users/{userId}"},
	{"post", "/api/v1/admin/users/{userId}/disable"},
	{"post", "/api/v1/admin/users/{userId}/force-logout"},
	{"post", "/api/v1/admin/users/{userId}/reset-passkey"},
	{"post", "/api/v1/admin/users/{userId}/reset-mfa"},
	{"get", "/api/v1/admin/tenants"},
	{"post", "/api/v1/admin/tenants"},
	{"get", "/api/v1/admin/tenants/{tenantId}"},
	{"patch", "/api/v1/admin/tenants/{tenantId}"},
	{"delete", "/api/v1/admin/tenants/{tenantId}"},
	{"get", "/api/v1/admin/clients"},
	{"get", "/api/v1/admin/clients/{clientId}/usage"},
	{"post", "/api/v1/admin/clients"},
	{"get", "/api/v1/admin/clients/{clientId}"},
	{"patch", "/api/v1/admin/clients/{clientId}"},
	{"delete", "/api/v1/admin/clients/{clientId}"},
	{"get", "/api/v1/admin/audit-logs"},
	{"get", "/api/v1/admin/login-history"},
	{"get", "/api/v1/admin/tenants/{tenantId}/identity-providers"},
	{"patch", "/api/v1/admin/tenants/{tenantId}/identity-providers/{provider}"},
	{"get", "/api/v1/admin/tenants/{tenantId}/members"},
	{"post", "/api/v1/admin/tenants/{tenantId}/members"},
	{"delete", "/api/v1/admin/tenants/{tenantId}/members/{userId}"},
}

func TestOpenAPICoversPublicRoutes(t *testing.T) {
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	specPath := filepath.Clean(filepath.Join(filepath.Dir(thisFile), "..", "..", "..", "api", "openapi.yaml"))
	raw, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("read openapi: %v (path %s)", err, specPath)
	}
	spec := string(raw)

	for _, route := range requiredPublicRoutes {
		key := `"` + route.path + `":`
		idx := strings.Index(spec, key)
		if idx < 0 {
			t.Errorf("missing path %s", route.path)
			continue
		}
		// Scan until the next path key at the same indent level ("/...":).
		rest := spec[idx:]
		next := strings.Index(rest[1:], "\n  \"/")
		chunk := rest
		if next >= 0 {
			chunk = rest[:next+1]
		}
		if !strings.Contains(chunk, "\n    "+route.method+":") {
			t.Errorf("missing method %s on path %s", route.method, route.path)
		}
	}

	if strings.Contains(spec, `/api/v1/internal/tenants`) {
		t.Error("internal admin-key route must not appear in public OpenAPI")
	}
	for _, id := range []string{"getMe", "createToken", "introspectToken", "login", "getAdminStats", "getJwks"} {
		if !strings.Contains(spec, "operationId: "+id) {
			t.Errorf("expected operationId %s", id)
		}
	}
}
