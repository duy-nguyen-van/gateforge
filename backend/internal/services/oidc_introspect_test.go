package services

import (
	"context"
	"testing"
	"time"

	"github.com/lib/pq"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/stretchr/testify/require"
)

func newOIDCIntrospectFixture(t *testing.T) (OIDCService, *auth.OIDCSigner, *stubClientRepo, *refreshTokenTestRepo, *models.User) {
	t.Helper()
	cfg := testConfig()
	signer, err := auth.ProvideOIDCSigner(cfg)
	require.NoError(t, err)
	clients := &stubClientRepo{byClientID: map[string]*models.Client{}}
	users := newUserTestRepo()
	refresh := newRefreshTokenTestRepo()
	memberships := &stubMembershipRepo{active: map[string]map[string]bool{}}
	svc := ProvideOIDCService(cfg, signer, clients, newAuthCodeTestRepo(), users, refresh, memberships, &auditCapture{})
	u := users.seed("introspect@example.com", "secret")
	return svc, signer, clients, refresh, u
}

func seedConfidentialClient(clients *stubClientRepo, clientID, secret string) {
	clients.byClientID[clientID] = &models.Client{
		BaseModel:    models.NewBaseModel(),
		TenantID:     "tenant-1",
		ClientID:     clientID,
		ClientSecret: secret,
		IsPublic:     false,
		RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes:   pq.StringArray{"authorization_code"},
		Scopes:       pq.StringArray{"openid", "profile"},
	}
}

func TestOIDCService_Introspect_AccessToken(t *testing.T) {
	svc, signer, clients, _, u := newOIDCIntrospectFixture(t)
	seedConfidentialClient(clients, "rs-client", "rs-secret")

	access, _, err := signer.SignAccessTokenOIDC(u.ID, "app", "openid profile", "app")
	require.NoError(t, err)

	out, terr := svc.Introspect(context.Background(), "rs-client", "rs-secret", access, "access_token")
	require.Nil(t, terr)
	require.True(t, out.Active)
	require.Equal(t, "access_token", out.TokenType)
	require.Equal(t, u.ID, out.Sub)
	require.Equal(t, "app", out.ClientID)
	require.Equal(t, "openid profile", out.Scope)
	require.NotZero(t, out.Exp)
	require.NotZero(t, out.Iat)
}

func TestOIDCService_Introspect_RefreshToken(t *testing.T) {
	svc, _, clients, refresh, u := newOIDCIntrospectFixture(t)
	seedConfidentialClient(clients, "rs-client", "rs-secret")

	raw, hash, err := auth.NewOpaqueRefreshToken()
	require.NoError(t, err)
	require.NoError(t, refresh.Create(context.Background(), &models.RefreshToken{
		TenantID:      "tenant-1",
		UserID:        u.ID,
		OAuthClientID: "app",
		TokenHash:     hash,
		ExpiresAt:     time.Now().UTC().Add(time.Hour),
	}))

	out, terr := svc.Introspect(context.Background(), "rs-client", "rs-secret", raw, "refresh_token")
	require.Nil(t, terr)
	require.True(t, out.Active)
	require.Equal(t, "refresh_token", out.TokenType)
	require.Equal(t, u.ID, out.Sub)
	require.Equal(t, "app", out.ClientID)
	require.Empty(t, out.Scope)
}

func TestOIDCService_Introspect_Inactive(t *testing.T) {
	svc, _, clients, _, _ := newOIDCIntrospectFixture(t)
	seedConfidentialClient(clients, "rs-client", "rs-secret")

	out, terr := svc.Introspect(context.Background(), "rs-client", "rs-secret", "not-a-real-token", "")
	require.Nil(t, terr)
	require.False(t, out.Active)
	require.Empty(t, out.Sub)
	require.Empty(t, out.ClientID)
}

func TestOIDCService_Introspect_RevokedRefresh(t *testing.T) {
	svc, _, clients, refresh, u := newOIDCIntrospectFixture(t)
	seedConfidentialClient(clients, "rs-client", "rs-secret")

	raw, hash, err := auth.NewOpaqueRefreshToken()
	require.NoError(t, err)
	rt := &models.RefreshToken{
		TenantID:      "tenant-1",
		UserID:        u.ID,
		OAuthClientID: "app",
		TokenHash:     hash,
		ExpiresAt:     time.Now().UTC().Add(time.Hour),
	}
	require.NoError(t, refresh.Create(context.Background(), rt))
	require.NoError(t, refresh.RevokeByID(context.Background(), rt.ID))

	out, terr := svc.Introspect(context.Background(), "rs-client", "rs-secret", raw, "refresh_token")
	require.Nil(t, terr)
	require.False(t, out.Active)
}

func TestOIDCService_Introspect_AuthErrors(t *testing.T) {
	svc, signer, clients, _, u := newOIDCIntrospectFixture(t)
	seedConfidentialClient(clients, "rs-client", "rs-secret")
	clients.byClientID["public-app"] = &models.Client{
		BaseModel: models.NewBaseModel(),
		ClientID:  "public-app",
		IsPublic:  true,
	}

	access, _, err := signer.SignAccessTokenOIDC(u.ID, "app", "openid", "app")
	require.NoError(t, err)

	tests := []struct {
		name         string
		clientID     string
		clientSecret string
		token        string
		wantCode     string
	}{
		{name: "missing token", clientID: "rs-client", clientSecret: "rs-secret", token: "", wantCode: constants.OAuthInvalidRequest},
		{name: "missing client", clientID: "", clientSecret: "rs-secret", token: access, wantCode: constants.OAuthInvalidClient},
		{name: "unknown client", clientID: "missing", clientSecret: "x", token: access, wantCode: constants.OAuthInvalidClient},
		{name: "wrong secret", clientID: "rs-client", clientSecret: "wrong", token: access, wantCode: constants.OAuthInvalidClient},
		{name: "public client", clientID: "public-app", clientSecret: "", token: access, wantCode: constants.OAuthInvalidClient},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, terr := svc.Introspect(context.Background(), tt.clientID, tt.clientSecret, tt.token, "")
			require.Nil(t, out)
			require.NotNil(t, terr)
			require.Equal(t, tt.wantCode, terr.Code)
		})
	}
}

func TestOIDCService_Introspect_HintOrdering(t *testing.T) {
	svc, signer, clients, refresh, u := newOIDCIntrospectFixture(t)
	seedConfidentialClient(clients, "rs-client", "rs-secret")

	access, _, err := signer.SignAccessTokenOIDC(u.ID, "app", "openid", "app")
	require.NoError(t, err)
	raw, hash, err := auth.NewOpaqueRefreshToken()
	require.NoError(t, err)
	require.NoError(t, refresh.Create(context.Background(), &models.RefreshToken{
		TenantID:      "tenant-1",
		UserID:        u.ID,
		OAuthClientID: "app",
		TokenHash:     hash,
		ExpiresAt:     time.Now().UTC().Add(time.Hour),
	}))

	// Wrong hint still finds the other token type.
	out, terr := svc.Introspect(context.Background(), "rs-client", "rs-secret", access, "refresh_token")
	require.Nil(t, terr)
	require.True(t, out.Active)
	require.Equal(t, "access_token", out.TokenType)

	out, terr = svc.Introspect(context.Background(), "rs-client", "rs-secret", raw, "access_token")
	require.Nil(t, terr)
	require.True(t, out.Active)
	require.Equal(t, "refresh_token", out.TokenType)

	// No hint: JWT-shaped tries access first.
	out, terr = svc.Introspect(context.Background(), "rs-client", "rs-secret", access, "")
	require.Nil(t, terr)
	require.True(t, out.Active)
	require.Equal(t, "access_token", out.TokenType)
}

func TestLooksLikeJWT(t *testing.T) {
	require.True(t, looksLikeJWT("aaa.bbb.ccc"))
	require.False(t, looksLikeJWT("opaque-token"))
	require.False(t, looksLikeJWT("aaa.bbb"))
	require.False(t, looksLikeJWT(".."))
}
