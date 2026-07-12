package services

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/crypto"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/stretchr/testify/require"
)

func TestOIDCFederationProvider_AuthorizeRedirectURL(t *testing.T) {
	var issuer string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{
			"issuer":%q,
			"authorization_endpoint":%q,
			"token_endpoint":%q,
			"jwks_uri":%q
		}`, issuer, issuer+"/auth", issuer+"/token", issuer+"/jwks")
	}))
	t.Cleanup(srv.Close)
	issuer = srv.URL

	spec := constants.IdentityProviderSpec{
		ID: "mock-idp", DisplayName: "Mock", IssuerURL: issuer, Scopes: []string{"openid"},
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

	url, err := provider.AuthorizeRedirectURL(context.Background(), "tenant-1", "state-1", "nonce-1")
	require.NoError(t, err)
	require.Contains(t, url, "client_id=cid")
	require.Contains(t, url, "state=state-1")
}
