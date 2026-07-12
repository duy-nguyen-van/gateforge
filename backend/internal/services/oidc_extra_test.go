package services

import (
	"context"
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

func TestOIDCService_Authorize_NoMembership(t *testing.T) {
	svc, _, clients, _, users, _ := newOIDCTestService(t)
	u := users.seed("nomember@example.com", "secret")
	clients.byClientID["app"] = &models.Client{
		BaseModel: models.BaseModel{ID: "c1"}, TenantID: "tenant-x", ClientID: "app",
		RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes:   pq.StringArray{"authorization_code"},
		Scopes:       pq.StringArray{"openid"},
		IsPublic:     true,
	}
	_, err := svc.Authorize(context.Background(), u.ID, &dtos.AuthorizeQuery{
		ClientID: "app", RedirectURI: "http://localhost/callback", ResponseType: "code",
		CodeChallenge: pkceChallenge("v"), CodeChallengeMethod: "S256",
	})
	require.NotNil(t, err)
}

func TestOIDCService_Authorize_GrantNotAllowed(t *testing.T) {
	svc, _, clients, _, users, memberships := newOIDCTestService(t)
	u := users.seed("grant@example.com", "secret")
	memberships.active[u.ID] = map[string]bool{"tenant-1": true}
	clients.byClientID["app"] = &models.Client{
		BaseModel: models.BaseModel{ID: "c1"}, TenantID: "tenant-1", ClientID: "app",
		RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes:   pq.StringArray{"client_credentials"},
		Scopes:       pq.StringArray{"openid"},
		IsPublic:     true,
	}
	_, err := svc.Authorize(context.Background(), u.ID, &dtos.AuthorizeQuery{
		ClientID: "app", RedirectURI: "http://localhost/callback", ResponseType: "code",
		CodeChallenge: pkceChallenge("v"), CodeChallengeMethod: "S256",
	})
	require.NotNil(t, err)
}

func TestOIDCService_Authorize_InvalidScope(t *testing.T) {
	svc, _, clients, _, users, memberships := newOIDCTestService(t)
	u := users.seed("scope@example.com", "secret")
	memberships.active[u.ID] = map[string]bool{"tenant-1": true}
	clients.byClientID["app"] = &models.Client{
		BaseModel: models.BaseModel{ID: "c1"}, TenantID: "tenant-1", ClientID: "app",
		RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes:   pq.StringArray{"authorization_code"},
		Scopes:       pq.StringArray{"openid"},
		IsPublic:     true,
	}
	_, err := svc.Authorize(context.Background(), u.ID, &dtos.AuthorizeQuery{
		ClientID: "app", RedirectURI: "http://localhost/callback", ResponseType: "code",
		Scope: "admin", CodeChallenge: pkceChallenge("v"), CodeChallengeMethod: "S256",
	})
	require.NotNil(t, err)
}

func TestOIDCService_AuthorizationCodeToken_PublicClient(t *testing.T) {
	svc, _, clients, authCodes, users, memberships := newOIDCTestService(t)
	u := users.seed("public@example.com", "secret")
	tenantID := "tenant-pub"
	memberships.active[u.ID] = map[string]bool{tenantID: true}
	clientRecordID := models.NewBaseModel().ID
	clients.byClientID["public-app"] = &models.Client{
		BaseModel: models.BaseModel{ID: clientRecordID}, TenantID: tenantID, ClientID: "public-app",
		IsPublic: true, RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes: pq.StringArray{"authorization_code"}, Scopes: pq.StringArray{"openid profile"},
	}
	code := "pub-code"
	verifier := "public-verifier-xyz"
	authCodes.byCode[code] = &models.AuthorizationCode{
		Code: code, TenantID: tenantID, OAuthClientID: "public-app", UserID: u.ID,
		Scope: "openid profile", RedirectURI: "http://localhost/callback",
		CodeChallenge: pkceChallenge(verifier), CodeChallengeMethod: "S256",
		ExpiresAt: time.Now().UTC().Add(5 * time.Minute), ClientRecordID: &clientRecordID,
	}
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", "http://localhost/callback")
	form.Set("client_id", "public-app")
	form.Set("code_verifier", verifier)

	tokens, tokenErr := svc.AuthorizationCodeToken(context.Background(), tenantID, "public-app", "", form)
	require.Nil(t, tokenErr)
	require.NotEmpty(t, tokens.AccessToken)
}

func TestOIDCService_AuthorizationCodeToken_InvalidCode(t *testing.T) {
	svc, _, _, _, _, _ := newOIDCTestService(t)
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", "missing")
	form.Set("redirect_uri", "http://localhost/callback")
	form.Set("client_id", "app")

	_, tokenErr := svc.AuthorizationCodeToken(context.Background(), "", "app", "", form)
	require.NotNil(t, tokenErr)
	require.Equal(t, constants.OAuthInvalidGrant, tokenErr.Code)
}

func TestOIDCService_AuthorizationCodeToken_WrongRedirectURI(t *testing.T) {
	svc, _, clients, authCodes, users, _ := newOIDCTestService(t)
	u := users.seed("redirect@example.com", "secret")
	clientRecordID := models.NewBaseModel().ID
	clients.byClientID["app"] = &models.Client{
		BaseModel: models.BaseModel{ID: clientRecordID}, TenantID: "t1", ClientID: "app",
		ClientSecret: "sec", RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes: pq.StringArray{"authorization_code"}, Scopes: pq.StringArray{"openid"},
	}
	code := "code-wrong-redirect"
	authCodes.byCode[code] = &models.AuthorizationCode{
		Code: code, TenantID: "t1", OAuthClientID: "app", UserID: u.ID,
		RedirectURI: "http://localhost/callback", ExpiresAt: time.Now().UTC().Add(time.Minute),
		ClientRecordID: &clientRecordID,
	}
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", "http://evil/callback")
	form.Set("client_id", "app")

	_, tokenErr := svc.AuthorizationCodeToken(context.Background(), "t1", "app", "sec", form)
	require.NotNil(t, tokenErr)
}

func TestOIDCService_AuthorizationCodeToken_InvalidPKCE(t *testing.T) {
	svc, _, clients, authCodes, users, _ := newOIDCTestService(t)
	u := users.seed("pkce@example.com", "secret")
	clientRecordID := models.NewBaseModel().ID
	clients.byClientID["app"] = &models.Client{
		BaseModel: models.BaseModel{ID: clientRecordID}, TenantID: "t1", ClientID: "app",
		IsPublic: true, RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes: pq.StringArray{"authorization_code"}, Scopes: pq.StringArray{"openid"},
	}
	code := "code-bad-pkce"
	authCodes.byCode[code] = &models.AuthorizationCode{
		Code: code, TenantID: "t1", OAuthClientID: "app", UserID: u.ID,
		RedirectURI: "http://localhost/callback", CodeChallenge: pkceChallenge("good"),
		CodeChallengeMethod: "S256", ExpiresAt: time.Now().UTC().Add(time.Minute),
		ClientRecordID: &clientRecordID,
	}
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", "http://localhost/callback")
	form.Set("client_id", "app")
	form.Set("code_verifier", "bad-verifier")

	_, tokenErr := svc.AuthorizationCodeToken(context.Background(), "t1", "app", "", form)
	require.NotNil(t, tokenErr)
}

func TestOIDCService_UserInfo_ProfileOnlyScope(t *testing.T) {
	svc, signer, clients, _, users, _ := newOIDCTestService(t)
	u := users.seed("profile@example.com", "secret")
	u.FirstName = "Test"
	clients.byClientID["app"] = &models.Client{ClientID: "app"}
	access, _, err := signer.SignAccessTokenOIDC(u.ID, "app", "profile", "app")
	require.NoError(t, err)
	out, tokenErr := svc.UserInfo(context.Background(), access)
	require.Nil(t, tokenErr)
	require.Equal(t, u.ID, out["sub"])
	_, hasEmail := out["email"]
	require.False(t, hasEmail)
}

func TestOIDCService_OpenIDIssuer_Default(t *testing.T) {
	cfg := testConfig()
	cfg.AppBaseURL = ""
	signer, err := auth.ProvideOIDCSigner(cfg)
	require.NoError(t, err)
	svc := ProvideOIDCService(cfg, signer, &stubClientRepo{byClientID: map[string]*models.Client{}}, newAuthCodeTestRepo(), newUserTestRepo(), newRefreshTokenTestRepo(), &stubMembershipRepo{active: map[string]map[string]bool{}}, &auditCapture{})
	require.Equal(t, "http://localhost:3000", svc.OpenIDIssuer())
}
