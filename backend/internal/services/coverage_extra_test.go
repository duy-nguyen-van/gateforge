package services

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/lib/pq"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/crypto"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/stretchr/testify/require"
)

func TestAdminService_ConfigureIdentityProvider_Success(t *testing.T) {
	tip := &adminIdpTipStub{}
	svc := newAdminIdpTestService(tip)
	enabled := true
	secret := "oauth-secret-value"
	err := svc.ConfigureIdentityProvider(context.Background(), "tenant-1", constants.IdentityProviderGoogle, &dtos.PatchIdentityProviderRequest{
		Enabled:           &enabled,
		OAuthClientID:     "google-client",
		OAuthClientSecret: secret,
	}, constants.AuditActorTypeUser)
	require.NoError(t, err)
	row := tip.byKey[tip.key("tenant-1", constants.IdentityProviderGoogle)]
	require.True(t, row.Enabled)
	require.Equal(t, "google-client", row.OAuthClientID)
}

func TestAdminService_ConfigureIdentityProvider_NilBody(t *testing.T) {
	svc := newAdminIdpTestService(&adminIdpTipStub{})
	err := svc.ConfigureIdentityProvider(context.Background(), "tenant-1", constants.IdentityProviderGoogle, nil, constants.AuditActorTypeUser)
	require.Error(t, err)
}

func TestAdminClient_NormalizeHelpers(t *testing.T) {
	uris, err := normalizeRedirectUris([]string{" http://localhost/callback ", ""})
	require.NoError(t, err)
	require.Equal(t, []string{"http://localhost/callback"}, uris)

	_, err = normalizeRedirectUris([]string{"  ", ""})
	require.Error(t, err)

	grants := normalizeGrantTypes([]string{"", "authorization_code", "authorization_code"})
	require.Equal(t, defaultClientGrantTypes, grants)

	scopes := normalizeScopes([]string{"", "openid", "openid"})
	require.Equal(t, []string{"openid"}, scopes)

	id, err := generateClientIdentifier()
	require.NoError(t, err)
	require.NotEmpty(t, id)

	secret, err := generateClientSecret()
	require.NoError(t, err)
	require.NotEmpty(t, secret)
}

func TestAdminService_UpdateClient_AllFields(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	clientRepo := newAdminClientTestRepo()
	id := "client-all"
	clientRepo.clients[id] = &models.Client{
		BaseModel:    models.BaseModel{ID: id},
		TenantID:     "tenant-1",
		ClientID:     "app",
		Name:         "App",
		RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes:   pq.StringArray{"authorization_code"},
		Scopes:       pq.StringArray{"openid"},
		IsPublic:     true,
	}
	svc := newAdminClientTestService(tenantRepo, clientRepo)

	name := "Renamed"
	uris := []string{"http://localhost/new"}
	grants := []string{"authorization_code", "refresh_token"}
	scopes := []string{"openid", "profile"}
	isPublic := false
	updated, err := svc.UpdateClient(context.Background(), id, &dtos.AdminUpdateClientRequest{
		Name:         &name,
		RedirectUris: uris,
		GrantTypes:   grants,
		Scopes:       scopes,
		IsPublic:     &isPublic,
	})
	require.NoError(t, err)
	require.Equal(t, "Renamed", updated.Name)
	require.False(t, updated.IsPublic)
}

func TestAdminService_DeleteTenant_Success(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	tenantID := "tenant-del"
	tenantRepo.tenants[tenantID] = &models.Tenant{BaseModel: models.BaseModel{ID: tenantID}, Name: "Delete Me"}
	svc := ProvideAdminService(testConfig(), nil, tenantRepo, nil, nil, nil, nil, nil, nil, nil, nil, nil, &auditCapture{}, nil, nil)

	require.NoError(t, svc.DeleteTenant(context.Background(), tenantID))
	_, ok := tenantRepo.tenants[tenantID]
	require.False(t, ok)
}

func TestFederationService_CompleteOAuthLogin_ErrorPaths(t *testing.T) {
	provider := &mockFedProvider{id: "google", configured: true, exchangeErr: errors.ExternalServiceError("exchange failed", nil)}
	tip := &fedTipStub{configured: map[string]bool{"tenant-1:google": true}, enabled: map[string]bool{"tenant-1:google": true}}
	cache := newMemCache()
	svc := &federationService{
		cfg: testConfig(), cache: cache, userRepo: newUserTestRepo(),
		membershipRepo: &stubMembershipRepo{active: map[string]map[string]bool{}},
		fedRepo:        newFedIdentityTestRepo(), tipRepo: tip,
		providers: map[string]FederationIdentityProvider{"google": provider},
		audit:     &auditCapture{},
	}

	state, err := federationRandomHex(8)
	require.NoError(t, err)
	payload, _ := json.Marshal(oauthStatePayload{TenantID: "tenant-1", ProviderID: "microsoft", Nonce: "n"})
	require.NoError(t, cache.Set(context.Background(), federationStateCacheKey("google", state), string(payload), 0))
	_, _, _, err = svc.CompleteOAuthLogin(context.Background(), "google", "code", state)
	require.Error(t, err)

	tip.enabled["tenant-1:google"] = false
	payload2, _ := json.Marshal(oauthStatePayload{TenantID: "tenant-1", ProviderID: "google", Nonce: "n"})
	state2, _ := federationRandomHex(8)
	require.NoError(t, cache.Set(context.Background(), federationStateCacheKey("google", state2), string(payload2), 0))
	_, _, _, err = svc.CompleteOAuthLogin(context.Background(), "google", "code", state2)
	require.Error(t, err)

	payload3, _ := json.Marshal(oauthStatePayload{TenantID: "tenant-1", ProviderID: "google", Nonce: "n"})
	state3, _ := federationRandomHex(8)
	tip.enabled["tenant-1:google"] = true
	require.NoError(t, cache.Set(context.Background(), federationStateCacheKey("google", state3), string(payload3), 0))
	_, _, _, err = svc.CompleteOAuthLogin(context.Background(), "google", "code", state3)
	require.Error(t, err)
}

func TestFederationService_BuildAuthorizeNotConfigured(t *testing.T) {
	provider := &mockFedProvider{id: "google", configured: false}
	tip := &fedTipStub{configured: map[string]bool{"tenant-1:google": false}}
	svc := newFederationTestService(map[string]FederationIdentityProvider{"google": provider}, tip)
	_, err := svc.BuildAuthorizeRedirectURL(context.Background(), "google", "/", "tenant-1")
	require.Error(t, err)
}

func TestUserService_SelectTenant_InvalidToken(t *testing.T) {
	users := newUserTestRepo()
	memberships := &stubMembershipRepo{active: map[string]map[string]bool{}}
	svc := newUserTestService(t, users, memberships)

	_, err := svc.SelectTenant(context.Background(), &dtos.TenantSelectRequest{
		SelectionToken: "invalid",
		TenantID:       "tenant-1",
	})
	require.Error(t, err)
}

func TestAdminIdpTipStub_UpdateProviderPreservesExisting(t *testing.T) {
	enc, err := crypto.EncryptMFASecret(testConfig().MFAEncryptionKey, "secret")
	require.NoError(t, err)
	tip := &adminIdpTipStub{byKey: map[string]*models.TenantIdentityProvider{
		"tenant-1:google": {TenantID: "tenant-1", Provider: constants.IdentityProviderGoogle, OAuthClientID: "existing", OAuthClientSecretEncrypted: enc, Enabled: true},
	}}
	svc := newAdminIdpTestService(tip)
	disabled := false
	err = svc.ConfigureIdentityProvider(context.Background(), "tenant-1", constants.IdentityProviderGoogle, &dtos.PatchIdentityProviderRequest{
		Enabled: &disabled,
	}, constants.AuditActorTypeUser)
	require.NoError(t, err)
	require.False(t, tip.byKey["tenant-1:google"].Enabled)
}

func TestAdminService_ListIdentityProviders_WithExistingTip(t *testing.T) {
	enc, _ := crypto.EncryptMFASecret(testConfig().MFAEncryptionKey, "secret")
	tip := &adminIdpTipStub{byKey: map[string]*models.TenantIdentityProvider{
		"tenant-1:" + constants.IdentityProviderGoogle: {
			TenantID: "tenant-1", Provider: constants.IdentityProviderGoogle,
			OAuthClientID: "cid", OAuthClientSecretEncrypted: enc, Enabled: true,
		},
	}}
	svc := newAdminIdpTestService(tip)
	rows, _, err := svc.ListIdentityProviders(context.Background(), "tenant-1", dtos.NewPageableRequest())
	require.NoError(t, err)
	require.True(t, rows[0].Configured)
}

func TestFederationService_ResolveInactiveUser(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("inactive@example.com", "secret")
	u.Status = constants.UserStatusDisabled
	fedRepo := newFedIdentityTestRepo()
	fi := &models.FederatedIdentity{UserID: u.ID, Provider: "google", Subject: "sub-1", User: u}
	fedRepo.byProviderSub[fedRepo.key("google", "sub-1")] = fi
	svc := &federationService{userRepo: users, fedRepo: fedRepo, membershipRepo: &stubMembershipRepo{active: map[string]map[string]bool{}}}

	_, err := svc.federationResolveOrProvisionUser(context.Background(), "tenant-1", "google", &domains.OIDCUserClaims{Sub: "sub-1"})
	require.Error(t, err)
}

func TestAdminService_UpdateTenant_AllFields(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	tenantID := "tenant-upd"
	tenantRepo.tenants[tenantID] = &models.Tenant{BaseModel: models.BaseModel{ID: tenantID}, Name: "Old", Domain: "old.example.com"}
	svc := ProvideAdminService(testConfig(), nil, tenantRepo, nil, nil, nil, nil, nil, nil, nil, nil, nil, &auditCapture{}, nil, nil)

	name := "New Name"
	domain := "new.example.com"
	updated, err := svc.UpdateTenant(context.Background(), tenantID, &dtos.AdminUpdateTenantRequest{
		Name: &name, Domain: &domain,
	})
	require.NoError(t, err)
	require.Equal(t, "New Name", updated.Name)
}

func TestAdminService_DisableUser_AlreadyDisabled(t *testing.T) {
	targetID := "user-disabled"
	users := &adminUserTestRepo{users: map[string]*models.User{
		targetID: {BaseModel: models.BaseModel{ID: targetID}, Status: constants.UserStatusDisabled},
	}}
	svc, _, _, _, _ := newAdminUserTestService(users)
	require.NoError(t, svc.DisableUser(context.Background(), "actor", targetID))
}

func TestUserService_UpdateProfile_Inactive(t *testing.T) {
	repo := newProfileTestUserRepo()
	user := &models.User{
		BaseModel: models.NewBaseModel(), Status: constants.UserStatusDisabled,
	}
	repo.users[user.ID] = user
	svc := &userService{userRepo: repo}
	name := "X"
	_, err := svc.UpdateProfile(context.Background(), user.ID, &dtos.UpdateProfileRequest{FirstName: &name})
	require.Error(t, err)
}

func TestUserService_SwitchTenant_NoMembership(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("switch@example.com", "secret")
	memberships := &stubMembershipRepo{active: map[string]map[string]bool{}}
	svc := newUserTestService(t, users, memberships)
	_, err := svc.SwitchTenant(context.Background(), u.ID, "missing-tenant")
	require.Error(t, err)
}

func TestOIDCService_AuthorizationCodeToken_NoOpenIDScope(t *testing.T) {
	svc, _, clients, authCodes, users, _ := newOIDCTestService(t)
	u := users.seed("noid@example.com", "secret")
	clientRecordID := models.NewBaseModel().ID
	clients.byClientID["app"] = &models.Client{
		BaseModel: models.BaseModel{ID: clientRecordID}, TenantID: "t1", ClientID: "app",
		ClientSecret: "sec", RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes: pq.StringArray{"authorization_code"}, Scopes: pq.StringArray{"profile"},
	}
	code := "code-profile-only"
	authCodes.byCode[code] = &models.AuthorizationCode{
		Code: code, TenantID: "t1", OAuthClientID: "app", UserID: u.ID,
		Scope: "profile", RedirectURI: "http://localhost/callback",
		ExpiresAt: time.Now().UTC().Add(time.Minute), ClientRecordID: &clientRecordID,
	}
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", "http://localhost/callback")
	form.Set("client_id", "app")

	tokens, tokenErr := svc.AuthorizationCodeToken(context.Background(), "t1", "app", "sec", form)
	require.Nil(t, tokenErr)
	require.Empty(t, tokens.IDToken)
	require.NotEmpty(t, tokens.AccessToken)
}

func TestFederationService_CompleteOAuthLogin_InvalidPayload(t *testing.T) {
	provider := &mockFedProvider{id: "google", configured: true}
	cache := newMemCache()
	svc := &federationService{
		cfg: testConfig(), cache: cache,
		providers: map[string]FederationIdentityProvider{"google": provider},
		tipRepo:   &fedTipStub{enabled: map[string]bool{"tenant-1:google": true}},
	}
	state := "state-invalid-json"
	require.NoError(t, cache.Set(context.Background(), federationStateCacheKey("google", state), "not-json", 0))
	_, _, _, err := svc.CompleteOAuthLogin(context.Background(), "google", "code", state)
	require.Error(t, err)
}

func TestFederationService_ProviderByID(t *testing.T) {
	svc := ProvideFederationService(testConfig(), newMemCache(), newUserTestRepo(), &stubMembershipRepo{active: map[string]map[string]bool{}}, newFedIdentityTestRepo(), &fedTipStub{}, nil, &auditCapture{}).(*federationService)
	_, err := svc.providerByID("")
	require.Error(t, err)
	_, err = svc.providerByID("nope")
	require.Error(t, err)
	p, err := svc.providerByID(constants.IdentityProviderGoogle)
	require.NoError(t, err)
	require.Equal(t, constants.IdentityProviderGoogle, p.ID())
}

func TestMFAService_ValidateTOTPForLogin_NoActiveRow(t *testing.T) {
	svc, _, _, _ := newMFATestService(t)
	ticket, _, err := svc.CreateLoginTicket(context.Background(), auth.MFAPendingPayload{UserID: "no-mfa", TenantID: "t1"})
	require.NoError(t, err)
	_, err = svc.VerifyLoginChallenge(context.Background(), ticket, "123456")
	require.Error(t, err)
}

func TestSessionService_InvalidateAllForUser_Error(t *testing.T) {
	svc := ProvideSessionService(testConfig(), brokenSessionRepo{}, &auditCapture{})
	require.Error(t, svc.InvalidateAllForUser(context.Background(), "user-1"))
}
