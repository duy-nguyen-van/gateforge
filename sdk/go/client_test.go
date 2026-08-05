package gateforge

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/sdk/go/openapi"
)

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient("https://iam.example.com/")
	if c.BaseURL() != "https://iam.example.com" {
		t.Fatalf("base URL = %q", c.BaseURL())
	}
	if c.UserAgent() == "" {
		t.Fatal("expected user agent")
	}
	if c.API() == nil {
		t.Fatal("expected generated API client")
	}
}

func TestBuildAuthorizeURL_AndPKCE(t *testing.T) {
	pkce, err := PKCEGenerate()
	if err != nil {
		t.Fatal(err)
	}
	if pkce.Method != "S256" || pkce.Verifier == "" || pkce.Challenge == "" {
		t.Fatalf("invalid pkce: %+v", pkce)
	}

	u, err := BuildAuthorizeURL("https://iam.example.com", AuthorizeParams{
		ClientID:      "app",
		RedirectURI:   "https://app.example.com/cb",
		Scope:         "openid profile",
		State:         "s1",
		Nonce:         "n1",
		CodeChallenge: pkce.Challenge,
	})
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := url.Parse(u)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Path != "/authorize" {
		t.Fatalf("path = %s", parsed.Path)
	}
	q := parsed.Query()
	if q.Get("client_id") != "app" || q.Get("code_challenge_method") != "S256" {
		t.Fatalf("query = %v", q)
	}
}

func TestClient_GetHealth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"meta": map[string]any{"message": "ok", "code": 200},
			"data": map[string]any{"status": "healthy", "service": "gateforge-iam", "version": "0.1.0", "timestamp": "2026-01-01T00:00:00Z"},
		})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, WithHTTPClient(srv.Client()))
	env, resp, err := c.GetHealth(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("status = %d", resp.StatusCode)
	}
	if env == nil {
		t.Fatal("nil envelope")
	}
	data := env.GetData()
	if data.GetStatus() != "healthy" {
		t.Fatalf("envelope = %+v", env)
	}
}

func TestClient_GetMe_SendsBearer(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("Authorization = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"meta": map[string]any{"message": "ok", "code": 200},
			"data": map[string]any{
				"id": "u1", "email": "a@b.co", "first_name": "A", "last_name": "B",
				"email_verified": true, "is_platform_admin": false, "mfa_enabled": false,
				"created_at": "2026-01-01T00:00:00Z", "updated_at": "2026-01-01T00:00:00Z",
			},
		})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, WithHTTPClient(srv.Client()), WithTokenProvider(func(context.Context) (string, error) {
		return "test-token", nil
	}))
	env, _, err := c.GetMe(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	data := env.GetData()
	if data.GetEmail() != "a@b.co" {
		t.Fatalf("email = %s", data.GetEmail())
	}
}

func TestClient_Login_ErrorEnvelope(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"meta": map[string]any{"error_code": "UNAUTHORIZED", "message": "bad credentials", "code": 401},
			"data": nil,
		})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, WithHTTPClient(srv.Client()))
	_, resp, err := c.Login(context.Background(), openapi.LoginRequest{Email: "a@b.co", Password: "x"})
	if err == nil {
		t.Fatal("expected error")
	}
	if resp == nil || resp.StatusCode != 401 {
		t.Fatalf("resp = %+v", resp)
	}
	if !strings.Contains(err.Error(), "401") && !strings.Contains(strings.ToLower(err.Error()), "unauthorized") {
		// Generated client returns GenericOpenAPIError; ensure non-nil failure.
		t.Logf("error: %v", err)
	}
}

func TestClient_ExchangeAuthorizationCode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/token" || r.Method != http.MethodPost {
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("grant_type") != "authorization_code" || r.Form.Get("code") != "abc" {
			t.Fatalf("form = %v content-type=%s", r.Form, r.Header.Get("Content-Type"))
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "at", "token_type": "Bearer", "expires_in": 3600, "id_token": "id",
		})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, WithHTTPClient(srv.Client()))
	tok, _, err := c.ExchangeAuthorizationCode(context.Background(), "client", "https://app/cb", "abc", "verifier", "")
	if err != nil {
		t.Fatal(err)
	}
	if tok.GetAccessToken() != "at" {
		t.Fatalf("token = %+v", tok)
	}
}

func TestClient_IntrospectToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/introspect" || r.Method != http.MethodPost {
			t.Fatalf("%s %s", r.Method, r.URL.Path)
		}
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if r.Form.Get("token") != "access-jwt" || r.Form.Get("token_type_hint") != "access_token" {
			t.Fatalf("form = %v", r.Form)
		}
		if r.Form.Get("client_id") != "rs-client" || r.Form.Get("client_secret") != "rs-secret" {
			t.Fatalf("client auth form = %v", r.Form)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"active": true, "token_type": "access_token", "sub": "user-1", "client_id": "app",
		})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, WithHTTPClient(srv.Client()))
	out, _, err := c.IntrospectToken(context.Background(), "rs-client", "rs-secret", "access-jwt", "access_token")
	if err != nil {
		t.Fatal(err)
	}
	if !out.GetActive() || out.GetSub() != "user-1" || out.GetTokenType() != "access_token" {
		t.Fatalf("introspect = %+v", out)
	}
}
