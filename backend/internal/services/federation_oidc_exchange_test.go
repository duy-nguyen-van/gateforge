package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/crypto"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestOIDCFederationProvider_ExchangeAuthorizationCode(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	var issuer string
	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	issuer = srv.URL

	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{
			"issuer":%q,
			"authorization_endpoint":%q,
			"token_endpoint":%q,
			"jwks_uri":%q
		}`, issuer, issuer+"/auth", issuer+"/token", issuer+"/jwks")
	})
	mux.HandleFunc("/jwks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		n := base64.RawURLEncoding.EncodeToString(key.N.Bytes())
		_, _ = fmt.Fprintf(w, `{"keys":[{"kty":"RSA","kid":"test","n":%q,"e":"AQAB"}]}`, n)
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		idt := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"iss": issuer, "sub": "fed-sub-1", "aud": "cid",
			"nonce": "nonce-1", "email": "fed@example.com", "email_verified": true,
			"exp": now.Add(time.Hour).Unix(), "iat": now.Unix(),
		})
		idt.Header["kid"] = "test"
		raw, err := idt.SignedString(key)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"access_token":"at","token_type":"Bearer","id_token":%q}`, raw)
	})

	spec := constants.IdentityProviderSpec{
		ID: "mock-idp", DisplayName: "Mock", IssuerURL: issuer,
		Scopes: []string{"openid", "email"}, RequireVerifiedEmail: true,
	}
	enc, err := crypto.EncryptMFASecret(testConfig().MFAEncryptionKey, "client-secret")
	require.NoError(t, err)
	tip := &adminIdpTipStub{byKey: map[string]*models.TenantIdentityProvider{
		"tenant-1:mock-idp": {
			TenantID: "tenant-1", Provider: "mock-idp",
			OAuthClientID: "cid", OAuthClientSecretEncrypted: enc,
		},
	}}
	provider := newOIDCFederationProvider(testConfig(), tip, spec).(*oidcFederationProvider)

	claims, err := provider.ExchangeAuthorizationCode(context.Background(), "tenant-1", "auth-code", "nonce-1")
	require.NoError(t, err)
	require.Equal(t, "fed-sub-1", claims.Sub)
	require.Equal(t, "fed@example.com", claims.Email)
}
