package services

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/crypto"
	"github.com/gateforge-iam/gateforge-iam/internal/db"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestFederationCallbackURL(t *testing.T) {
	require.Equal(t, "http://localhost:8080/oidc/federation/google/callback",
		FederationCallbackURL(&config.Config{AppBaseURL: "http://localhost:8080"}, "google"))
	require.Equal(t, "http://localhost:3000/oidc/federation/google/callback",
		FederationCallbackURL(&config.Config{}, "google"))
}

func TestFederationEncryptionKey(t *testing.T) {
	cfg := testConfig()
	require.Equal(t, cfg.MFAEncryptionKey, federationEncryptionKey(cfg))
	cfg.MFAEncryptionKey = ""
	require.Equal(t, cfg.JWTSecret, federationEncryptionKey(cfg))
}

func newFederationTestService(providers map[string]FederationIdentityProvider, tip *fedTipStub) *federationService {
	return &federationService{
		cfg:            testConfig(),
		cache:          newMemCache(),
		userRepo:       newUserTestRepo(),
		membershipRepo: &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}},
		fedRepo:        newFedIdentityTestRepo(),
		tipRepo:        tip,
		providers:      providers,
		audit:          &auditCapture{},
	}
}

type fedTipStub struct {
	configured map[string]bool
	enabled    map[string]bool
}

func (s *fedTipStub) key(tenantID, provider string) string { return tenantID + ":" + provider }

func (s *fedTipStub) IsProviderEnabled(_ context.Context, tenantID, provider string) (bool, error) {
	if s.enabled == nil {
		return false, nil
	}
	return s.enabled[s.key(tenantID, provider)], nil
}

func (s *fedTipStub) IsProviderConfigured(_ context.Context, tenantID, provider string) (bool, error) {
	if s.configured == nil {
		return false, nil
	}
	return s.configured[s.key(tenantID, provider)], nil
}

func (s *fedTipStub) GetByTenantAndProvider(context.Context, string, string) (*models.TenantIdentityProvider, error) {
	return nil, errors.NotFoundError("tip", nil)
}

func (s *fedTipStub) SetProviderEnabled(context.Context, string, string, bool) error { return nil }

func (s *fedTipStub) UpdateProvider(context.Context, string, string, repositories.TenantIdentityProviderPatch) (*models.TenantIdentityProvider, error) {
	return nil, nil
}

func (s *fedTipStub) ListByTenant(context.Context, string) ([]models.TenantIdentityProvider, error) {
	return nil, nil
}

func TestFederationService_BuildAuthorizeRedirectURL(t *testing.T) {
	provider := &mockFedProvider{id: "google", name: "Google", configured: true}
	tip := &fedTipStub{
		configured: map[string]bool{"tenant-1:google": true},
		enabled:    map[string]bool{"tenant-1:google": true},
	}
	svc := newFederationTestService(map[string]FederationIdentityProvider{"google": provider}, tip)

	url, err := svc.BuildAuthorizeRedirectURL(context.Background(), "google", "/dashboard", "tenant-1")
	require.NoError(t, err)
	require.Contains(t, url, "https://idp.example/authorize")

	_, err = svc.BuildAuthorizeRedirectURL(context.Background(), "", "/dashboard", "tenant-1")
	require.Error(t, err)

	_, err = svc.BuildAuthorizeRedirectURL(context.Background(), "unknown", "/dashboard", "tenant-1")
	require.Error(t, err)

	tip.enabled["tenant-1:google"] = false
	_, err = svc.BuildAuthorizeRedirectURL(context.Background(), "google", "/dashboard", "tenant-1")
	require.Error(t, err)
}

func TestFederationService_CompleteOAuthLogin(t *testing.T) {
	provider := &mockFedProvider{
		id:         "google",
		configured: true,
		claims: &domains.OIDCUserClaims{
			Sub:           "google-sub-1",
			Email:         "fed@example.com",
			EmailVerified: true,
			Name:          "Fed User",
		},
	}
	tip := &fedTipStub{
		configured: map[string]bool{"tenant-1:google": true},
		enabled:    map[string]bool{"tenant-1:google": true},
	}
	users := newUserTestRepo()
	existing := users.seed("fed@example.com", "secret")
	memberships := &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{
		existing.ID: {"tenant-1": true},
	}}
	fedRepo := newFedIdentityTestRepo()
	cache := newMemCache()
	svc := &federationService{
		cfg: testConfig(), cache: cache, userRepo: users, membershipRepo: memberships,
		fedRepo: fedRepo, tipRepo: tip, providers: map[string]FederationIdentityProvider{"google": provider},
		audit: &auditCapture{},
	}

	state, err := federationRandomHex(16)
	require.NoError(t, err)
	payload, _ := json.Marshal(oauthStatePayload{ReturnTo: "/home", TenantID: "tenant-1", Nonce: "nonce", ProviderID: "google"})
	require.NoError(t, cache.Set(context.Background(), federationStateCacheKey("google", state), string(payload), 0))

	u, tenantID, returnTo, err := svc.CompleteOAuthLogin(context.Background(), "google", "code-1", state)
	require.NoError(t, err)
	require.Equal(t, existing.ID, u.ID)
	require.Equal(t, "tenant-1", tenantID)
	require.Equal(t, "/home", returnTo)
	require.Len(t, fedRepo.created, 1)

	_, _, _, err = svc.CompleteOAuthLogin(context.Background(), "google", "", state)
	require.Error(t, err)

	_, _, _, err = svc.CompleteOAuthLogin(context.Background(), "google", "code", "bad-state")
	require.Error(t, err)
}

func TestFederationService_ListAvailableProviders(t *testing.T) {
	provider := &mockFedProvider{id: "google", name: "Google", configured: true}
	tip := &fedTipStub{
		configured: map[string]bool{"tenant-1:google": true},
		enabled:    map[string]bool{"tenant-1:google": true},
	}
	svc := newFederationTestService(map[string]FederationIdentityProvider{"google": provider}, tip)

	_, err := svc.ListAvailableProviders(context.Background(), "")
	require.Error(t, err)

	out, err := svc.ListAvailableProviders(context.Background(), "tenant-1")
	require.NoError(t, err)
	require.Len(t, out, 1)
	require.Equal(t, "google", out[0].Provider)
}

func TestFederationService_ResolveExistingFederatedUser(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("linked@example.com", "secret")
	fedRepo := newFedIdentityTestRepo()
	fi := &models.FederatedIdentity{UserID: u.ID, Provider: "google", Subject: "sub-1", User: u}
	fedRepo.byProviderSub[fedRepo.key("google", "sub-1")] = fi
	memberships := &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}}
	svc := &federationService{userRepo: users, fedRepo: fedRepo, membershipRepo: memberships}

	got, err := svc.federationResolveOrProvisionUser(context.Background(), "tenant-1", "google", &domains.OIDCUserClaims{Sub: "sub-1"})
	require.NoError(t, err)
	require.Equal(t, u.ID, got.ID)
}

func TestFederationService_EnsureMembershipCreates(t *testing.T) {
	memberships := &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}}
	svc := &federationService{membershipRepo: memberships}

	require.NoError(t, svc.ensureMembership(context.Background(), "user-1", "tenant-1"))
	require.True(t, memberships.active["user-1"]["tenant-1"])
}

func TestFederationService_LinkExistingEmailUser(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("existing@example.com", "secret")
	fedRepo := newFedIdentityTestRepo()
	memberships := &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}}
	svc := &federationService{userRepo: users, fedRepo: fedRepo, membershipRepo: memberships}

	got, err := svc.federationResolveOrProvisionUser(context.Background(), "tenant-1", "google", &domains.OIDCUserClaims{
		Sub: "sub-new", Email: "existing@example.com", EmailVerified: true,
	})
	require.NoError(t, err)
	require.Equal(t, u.ID, got.ID)
	require.Len(t, fedRepo.created, 1)
}

func TestFederationService_ProvisionNewFederatedUser(t *testing.T) {
	gormDB := federationSQLiteDB(t)
	svc := &federationService{
		pg:             &db.PostgresDB{DB: gormDB},
		userRepo:       newUserTestRepo(),
		fedRepo:        newFedIdentityTestRepo(),
		membershipRepo: &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}},
	}

	u, err := svc.federationResolveOrProvisionUser(context.Background(), "tenant-1", "google", &domains.OIDCUserClaims{
		Sub: "brand-new-sub", Email: "brand-new@example.com", EmailVerified: true, Name: "Brand New",
	})
	require.NoError(t, err)
	require.Equal(t, "brand-new@example.com", u.Email)
}

func federationSQLiteDB(t *testing.T) *gorm.DB {
	t.Helper()
	gormDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	for _, stmt := range []string{
		`CREATE TABLE users (id TEXT PRIMARY KEY, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME, first_name TEXT, last_name TEXT, email TEXT, email_lower TEXT, email_verified INTEGER, status TEXT, is_platform_admin INTEGER DEFAULT 0)`,
		`CREATE TABLE federated_identities (id TEXT PRIMARY KEY, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME, user_id TEXT, provider TEXT, subject TEXT, email_at_link TEXT)`,
		`CREATE TABLE tenant_memberships (id TEXT PRIMARY KEY, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME, user_id TEXT, tenant_id TEXT, role TEXT, status TEXT)`,
	} {
		require.NoError(t, gormDB.Exec(stmt).Error)
	}
	return gormDB
}

func TestFederationIsNotFoundAppError(t *testing.T) {
	require.True(t, federationIsNotFoundAppError(errors.NotFoundError("x", nil)))
	require.False(t, federationIsNotFoundAppError(errors.InternalError("x", nil)))
}

func TestOIDCFederationProvider_Basics(t *testing.T) {
	tip := &fedTipStub{configured: map[string]bool{"tenant-1:google": true}}
	provider := newOIDCFederationProvider(testConfig(), tip, constants.SupportedIdentityProviders[0]).(*oidcFederationProvider)
	require.Equal(t, constants.IdentityProviderGoogle, provider.ID())
	require.NotEmpty(t, provider.DisplayName())

	ok, err := provider.OAuthConfiguredForTenant(context.Background(), "tenant-1")
	require.NoError(t, err)
	require.True(t, ok)
}

func TestOIDCFederationProvider_LoadTenantOAuthErrors(t *testing.T) {
	tip := &adminIdpTipStub{}
	spec := constants.IdentityProviderSpec{
		ID: "test-idp", DisplayName: "Test", IssuerURL: "http://127.0.0.1:1", Scopes: []string{"openid"},
	}
	provider := newOIDCFederationProvider(testConfig(), tip, spec).(*oidcFederationProvider)

	_, _, err := provider.loadTenantOAuth(context.Background(), "tenant-1")
	require.Error(t, err)

	enc, err := crypto.EncryptMFASecret(testConfig().MFAEncryptionKey, "secret")
	require.NoError(t, err)
	tip.byKey = map[string]*models.TenantIdentityProvider{
		tip.key("tenant-1", "test-idp"): {
			TenantID:                   "tenant-1",
			Provider:                   "test-idp",
			OAuthClientID:              "client-id",
			OAuthClientSecretEncrypted: enc,
		},
	}
	_, _, err = provider.loadTenantOAuth(context.Background(), "tenant-1")
	require.Error(t, err)
}

func TestProvideFederationService(t *testing.T) {
	svc := ProvideFederationService(testConfig(), newMemCache(), newUserTestRepo(), &stubMembershipRepo{active: map[string]map[string]bool{}}, newFedIdentityTestRepo(), &fedTipStub{}, nil, &auditCapture{})
	require.NotNil(t, svc)
}
