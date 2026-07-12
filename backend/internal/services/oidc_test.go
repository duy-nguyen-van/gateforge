package services

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"net/url"
	"testing"
	"time"

	"github.com/lib/pq"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/stretchr/testify/require"
)

func newOIDCTestService(t *testing.T) (OIDCService, *auth.OIDCSigner, *stubClientRepo, *authCodeTestRepo, *userTestRepo, *stubMembershipRepo) {
	t.Helper()
	cfg := testConfig()
	signer, err := auth.ProvideOIDCSigner(cfg)
	require.NoError(t, err)
	clients := &stubClientRepo{byClientID: map[string]*models.Client{}}
	authCodes := newAuthCodeTestRepo()
	users := newUserTestRepo()
	memberships := &stubMembershipRepo{active: map[string]map[string]bool{}}
	svc := ProvideOIDCService(cfg, signer, clients, authCodes, users, newRefreshTokenTestRepo(), memberships, &auditCapture{})
	return svc, signer, clients, authCodes, users, memberships
}

func pkceChallenge(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func TestOIDCService_OpenIDIssuer(t *testing.T) {
	svc, _, _, _, _, _ := newOIDCTestService(t)
	require.Equal(t, testConfig().AppBaseURL, svc.OpenIDIssuer())
}

func TestOIDCService_Authorize_Success(t *testing.T) {
	svc, _, clients, authCodes, users, memberships := newOIDCTestService(t)
	u := users.seed("oidc@example.com", "secret")
	tenantID := "tenant-oidc"
	memberships.active[u.ID] = map[string]bool{tenantID: true}
	clientID := "app-client"
	clients.byClientID[clientID] = &models.Client{
		BaseModel:    models.NewBaseModel(),
		TenantID:     tenantID,
		ClientID:     clientID,
		IsPublic:     true,
		RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes:   pq.StringArray{"authorization_code"},
		Scopes:       pq.StringArray{"openid", "profile", "email"},
	}
	verifier := "verifier-1234567890"
	challenge := pkceChallenge(verifier)

	redirect, oauthErr := svc.Authorize(context.Background(), u.ID, &dtos.AuthorizeQuery{
		ResponseType:        "code",
		ClientID:            clientID,
		RedirectURI:         "http://localhost/callback",
		Scope:               "openid profile",
		State:               "state-1",
		CodeChallenge:       challenge,
		CodeChallengeMethod: "S256",
	})
	require.Nil(t, oauthErr)
	require.Contains(t, redirect, "code=")
	require.Contains(t, redirect, "state=state-1")
	require.Len(t, authCodes.byCode, 1)
}

func TestOIDCService_Authorize_Errors(t *testing.T) {
	svc, _, clients, _, users, memberships := newOIDCTestService(t)
	u := users.seed("oidc@example.com", "secret")
	clients.byClientID["app"] = &models.Client{
		BaseModel:    models.NewBaseModel(),
		TenantID:     "tenant-1",
		ClientID:     "app",
		RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes:   pq.StringArray{"authorization_code"},
		Scopes:       pq.StringArray{"openid"},
		IsPublic:     true,
	}

	_, err := svc.Authorize(context.Background(), u.ID, nil)
	require.NotNil(t, err)

	_, err = svc.Authorize(context.Background(), u.ID, &dtos.AuthorizeQuery{ClientID: "missing"})
	require.NotNil(t, err)

	memberships.active[u.ID] = map[string]bool{"tenant-1": true}
	_, err = svc.Authorize(context.Background(), u.ID, &dtos.AuthorizeQuery{
		ClientID: "app", RedirectURI: "http://evil/callback", ResponseType: "code",
	})
	require.NotNil(t, err)

	_, err = svc.Authorize(context.Background(), u.ID, &dtos.AuthorizeQuery{
		ClientID: "app", RedirectURI: "http://localhost/callback", ResponseType: "token",
	})
	require.NotNil(t, err)
}

func TestOIDCService_AuthorizationCodeToken(t *testing.T) {
	svc, signer, clients, authCodes, users, memberships := newOIDCTestService(t)
	u := users.seed("token@example.com", "secret")
	tenantID := "tenant-token"
	memberships.active[u.ID] = map[string]bool{tenantID: true}
	clientRecordID := models.NewBaseModel().ID
	clientSecret := "client-secret"
	clients.byClientID["confidential"] = &models.Client{
		BaseModel:    models.BaseModel{ID: clientRecordID},
		TenantID:     tenantID,
		ClientID:     "confidential",
		ClientSecret: clientSecret,
		IsPublic:     false,
		RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes:   pq.StringArray{"authorization_code"},
		Scopes:       pq.StringArray{"openid", "email", "profile"},
	}

	code := "auth-code-1"
	verifier := "pkce-verifier-value-123"
	challenge := pkceChallenge(verifier)
	authCodes.byCode[code] = &models.AuthorizationCode{
		Code:                code,
		TenantID:            tenantID,
		OAuthClientID:       "confidential",
		UserID:              u.ID,
		Scope:               "openid email profile",
		RedirectURI:         "http://localhost/callback",
		CodeChallenge:       challenge,
		CodeChallengeMethod: "S256",
		ExpiresAt:           time.Now().UTC().Add(5 * time.Minute),
		ClientRecordID:      &clientRecordID,
	}

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", "http://localhost/callback")
	form.Set("client_id", "confidential")
	form.Set("code_verifier", verifier)

	tokens, tokenErr := svc.AuthorizationCodeToken(context.Background(), tenantID, "confidential", clientSecret, form)
	require.Nil(t, tokenErr)
	require.NotEmpty(t, tokens.AccessToken)
	require.NotEmpty(t, tokens.RefreshToken)
	require.NotEmpty(t, tokens.IDToken)

	claims, err := signer.ParseAccessTokenOIDC(tokens.AccessToken)
	require.NoError(t, err)
	require.Equal(t, u.ID, claims.Subject)

	_, tokenErr = svc.AuthorizationCodeToken(context.Background(), tenantID, "confidential", clientSecret, url.Values{"grant_type": {"refresh_token"}})
	require.NotNil(t, tokenErr)
}

func TestOIDCService_UserInfo(t *testing.T) {
	svc, signer, clients, _, users, _ := newOIDCTestService(t)
	u := users.seed("info@example.com", "secret")
	u.FirstName = "Ada"
	u.LastName = "Lovelace"
	u.EmailVerified = true
	clientID := "app"
	clients.byClientID[clientID] = &models.Client{ClientID: clientID}

	access, _, err := signer.SignAccessTokenOIDC(u.ID, clientID, "openid email profile", clientID)
	require.NoError(t, err)

	out, tokenErr := svc.UserInfo(context.Background(), access)
	require.Nil(t, tokenErr)
	require.Equal(t, u.ID, out["sub"])
	require.Equal(t, u.Email, out["email"])
	require.Equal(t, "Ada Lovelace", out["name"])

	_, tokenErr = svc.UserInfo(context.Background(), "invalid-token")
	require.NotNil(t, tokenErr)
}

func TestOIDCHelpers(t *testing.T) {
	client := &models.Client{
		RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes:   pq.StringArray{"authorization_code"},
		Scopes:       pq.StringArray{"openid", "profile"},
	}
	require.True(t, redirectAllowed(client, "http://localhost/callback"))
	require.False(t, redirectAllowed(client, ""))
	require.True(t, grantAllowed(client, "authorization_code"))
	require.Equal(t, "", validateScopes(client, "openid profile"))
	require.Contains(t, validateScopes(client, "admin"), "scope not allowed")

	verifier := "abc"
	challenge := pkceChallenge(verifier)
	require.True(t, verifyPKCE(verifier, challenge, "S256"))
	require.False(t, verifyPKCE("wrong", challenge, "S256"))

	require.True(t, scopeIncludes("openid profile", "openid"))
	require.False(t, scopeIncludes("profile", "openid"))

	_, err := oauthRedirectErr("http://localhost/cb", "st", constants.OAuthInvalidRequest, "bad")
	require.NotNil(t, err)
	_, err = oauthRedirectErr("", "st", constants.OAuthInvalidRequest, "bad")
	require.NotNil(t, err)

	public := &models.Client{IsPublic: true}
	require.NotNil(t, authorizePKCEError(public, "", "", "http://localhost/cb", "st"))
	require.NotNil(t, authorizePKCEError(public, "ch", "plain", "http://localhost/cb", "st"))
	require.Nil(t, authorizePKCEError(public, "ch", "S256", "http://localhost/cb", "st"))
}
