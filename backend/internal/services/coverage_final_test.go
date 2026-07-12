package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/lib/pq"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/crypto"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"github.com/stretchr/testify/require"
)

// --- broken repo stubs for error-path coverage ---

type errUserCountRepo struct {
	*userTestRepo
	countErr error
}

func (r *errUserCountRepo) Count(context.Context) (int64, error) {
	return 0, r.countErr
}

type errMFACountRepo struct {
	*mfaTOTPTestRepo
	countErr error
}

func (r *errMFACountRepo) CountEnabled(context.Context) (int64, error) {
	return 0, r.countErr
}

type errSessionCountRepo struct {
	*sessionTestRepo
	countErr error
}

func (r *errSessionCountRepo) CountActive(context.Context) (int64, error) {
	return 0, r.countErr
}

type errTenantCountRepo struct {
	*adminTenantTestRepo
	countErr error
}

func (r *errTenantCountRepo) CountUsersByTenantID(context.Context, string) (int64, error) {
	return 0, r.countErr
}

type errAuditListRepo struct {
	listErr error
}

func (r *errAuditListRepo) Create(context.Context, *models.AuditLog) error { return nil }
func (r *errAuditListRepo) List(context.Context, repositories.AuditLogListFilters, *dtos.PageableRequest) (*dtos.DataResponse[models.AuditLog], error) {
	return nil, r.listErr
}
func (r *errAuditListRepo) Count(context.Context, repositories.AuditLogListFilters) (int64, error) {
	return 0, nil
}

type errMembershipDeleteRepo struct {
	*stubMembershipRepo
	deleteErr error
}

func (r *errMembershipDeleteRepo) Delete(context.Context, string, string) error {
	return r.deleteErr
}

type errAuthCodeRepo struct {
	*authCodeTestRepo
	createErr error
	deleteErr error
}

func (r *errAuthCodeRepo) Create(ctx context.Context, row *models.AuthorizationCode) error {
	if r.createErr != nil {
		return r.createErr
	}
	return r.authCodeTestRepo.Create(ctx, row)
}

func (r *errAuthCodeRepo) DeleteByCode(ctx context.Context, code string) error {
	if r.deleteErr != nil {
		return r.deleteErr
	}
	return r.authCodeTestRepo.DeleteByCode(ctx, code)
}

type errRefreshCreateRepo struct {
	*refreshTokenTestRepo
	createErr error
}

func (r *errRefreshCreateRepo) Create(ctx context.Context, rt *models.RefreshToken) error {
	if r.createErr != nil {
		return r.createErr
	}
	return r.refreshTokenTestRepo.Create(ctx, rt)
}

type errUserGetRepo struct {
	*userTestRepo
	getErr error
}

func (r *errUserGetRepo) GetByEmailLower(ctx context.Context, email string) (*models.User, error) {
	if r.getErr != nil {
		return nil, r.getErr
	}
	return r.userTestRepo.GetByEmailLower(ctx, email)
}

type errMFATOTPActiveRepo struct {
	*mfaTOTPTestRepo
	activeErr error
}

func (r *errMFATOTPActiveRepo) GetActiveByUserID(ctx context.Context, userID string) (*models.UserMFATOTP, error) {
	if r.activeErr != nil {
		return nil, r.activeErr
	}
	return r.mfaTOTPTestRepo.GetActiveByUserID(ctx, userID)
}

type fedProviderRedirectErr struct {
	mockFedProvider
	redirectErr error
}

func (m *fedProviderRedirectErr) AuthorizeRedirectURL(context.Context, string, string, string) (string, error) {
	return "", m.redirectErr
}

type fedProviderConfigErr struct {
	mockFedProvider
	configErr error
}

func (m *fedProviderConfigErr) OAuthConfiguredForTenant(context.Context, string) (bool, error) {
	return false, m.configErr
}

// --- admin coverage ---

func TestAdminService_GetStats_ZeroUsers(t *testing.T) {
	svc := &adminService{users: newUserTestRepo(), mfaTOTP: newMFATOTPTestRepo(), sessions: newSessionTestRepo()}
	stats, err := svc.GetStats(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0.0, stats.MFAEnabledPercent)
}

func TestAdminService_GetStats_Errors(t *testing.T) {
	svc := &adminService{users: &errUserCountRepo{userTestRepo: newUserTestRepo(), countErr: context.Canceled}}
	_, err := svc.GetStats(context.Background())
	require.Error(t, err)

	svc = &adminService{users: newUserTestRepo(), mfaTOTP: &errMFACountRepo{mfaTOTPTestRepo: newMFATOTPTestRepo(), countErr: context.Canceled}}
	_, err = svc.GetStats(context.Background())
	require.Error(t, err)

	svc = &adminService{users: newUserTestRepo(), mfaTOTP: newMFATOTPTestRepo(), sessions: &errSessionCountRepo{sessionTestRepo: newSessionTestRepo(), countErr: context.Canceled}}
	_, err = svc.GetStats(context.Background())
	require.Error(t, err)
}

func TestAdminService_ListTenants_CountError(t *testing.T) {
	tenantRepo := &errTenantCountRepo{adminTenantTestRepo: newAdminTenantTestRepo()}
	tenantRepo.tenants["t1"] = &models.Tenant{BaseModel: models.BaseModel{ID: "t1"}, Name: "A"}
	tenantRepo.countErr = context.Canceled
	svc := &adminService{tenants: tenantRepo}
	_, _, err := svc.ListTenants(context.Background(), dtos.NewPageableRequest())
	require.Error(t, err)
}

func TestAdminService_ListAuditLogs_Error(t *testing.T) {
	svc := &adminService{auditLogs: &errAuditListRepo{listErr: context.Canceled}}
	_, _, err := svc.ListAuditLogs(context.Background(), dtos.AdminAuditLogListParams{}, dtos.NewPageableRequest())
	require.Error(t, err)
}

func TestAdminService_ConfigureIdentityProvider_Validation(t *testing.T) {
	svc := newAdminIdpTestService(&adminIdpTipStub{})
	err := svc.ConfigureIdentityProvider(context.Background(), "t1", "unknown-provider", &dtos.PatchIdentityProviderRequest{}, constants.AuditActorTypeUser)
	require.Error(t, err)

	enabled := true
	err = svc.ConfigureIdentityProvider(context.Background(), "t1", constants.IdentityProviderGoogle, &dtos.PatchIdentityProviderRequest{Enabled: &enabled}, constants.AuditActorTypeUser)
	require.Error(t, err)
}

type errAdminIdpTipStub struct {
	adminIdpTipStub
	getErr error
}

func (s *errAdminIdpTipStub) GetByTenantAndProvider(context.Context, string, string) (*models.TenantIdentityProvider, error) {
	return nil, s.getErr
}

func TestAdminService_ListIdentityProviders_RepoError(t *testing.T) {
	svc := &adminService{cfg: testConfig(), tipRepo: &errAdminIdpTipStub{getErr: context.Canceled}}
	_, _, err := svc.ListIdentityProviders(context.Background(), "t1", dtos.NewPageableRequest())
	require.Error(t, err)
}

func TestAdminService_AddMemberByEmail_Errors(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	users := newUserTestRepo()
	memberships := &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}}
	svc := &adminService{tenants: tenantRepo, users: users, memberships: memberships}

	err := svc.AddMemberByEmail(context.Background(), "missing-tenant", "a@example.com", "")
	require.Error(t, err)

	tenantID := "tenant-1"
	tenantRepo.tenants[tenantID] = &models.Tenant{BaseModel: models.BaseModel{ID: tenantID}}
	err = svc.AddMemberByEmail(context.Background(), tenantID, "missing@example.com", "")
	require.Error(t, err)
}

func TestAdminService_RemoveMember_Error(t *testing.T) {
	memberships := &errMembershipDeleteRepo{
		stubMembershipRepo: &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}},
		deleteErr:          context.Canceled,
	}
	svc := &adminService{memberships: memberships, audit: &auditCapture{}}
	require.Error(t, svc.RemoveMember(context.Background(), "t1", "u1"))
}

func TestAdminService_GetClientByID_NotFound(t *testing.T) {
	svc := newAdminClientTestService(newAdminTenantTestRepo(), newAdminClientTestRepo())
	_, err := svc.GetClientByID(context.Background(), "missing")
	require.Error(t, err)
}

func TestAdminService_CreateClient_Validation(t *testing.T) {
	svc := newAdminClientTestService(newAdminTenantTestRepo(), newAdminClientTestRepo())
	_, err := svc.CreateClient(context.Background(), nil)
	require.Error(t, err)

	tenantID := "t1"
	tenantRepo := newAdminTenantTestRepo()
	tenantRepo.tenants[tenantID] = &models.Tenant{BaseModel: models.BaseModel{ID: tenantID}}
	svc = newAdminClientTestService(tenantRepo, newAdminClientTestRepo())
	_, err = svc.CreateClient(context.Background(), &dtos.AdminCreateClientRequest{TenantID: tenantID, Name: "  "})
	require.Error(t, err)
}

func TestAdminService_UpdateClient_Validation(t *testing.T) {
	svc := newAdminClientTestService(newAdminTenantTestRepo(), newAdminClientTestRepo())
	_, err := svc.UpdateClient(context.Background(), "id", nil)
	require.Error(t, err)
	_, err = svc.UpdateClient(context.Background(), "id", &dtos.AdminUpdateClientRequest{})
	require.Error(t, err)
}

func TestAdminService_GetTenantByID_NotFound(t *testing.T) {
	svc := &adminService{tenants: newAdminTenantTestRepo()}
	_, err := svc.GetTenantByID(context.Background(), "missing")
	require.Error(t, err)
}

func TestAdminService_GetUserByID_PasskeyListError(t *testing.T) {
	targetID := "user-1"
	users := &adminUserTestRepo{users: map[string]*models.User{
		targetID: {BaseModel: models.BaseModel{ID: targetID}, Email: "u@example.com", Status: constants.UserStatusActive},
	}}
	webauthn := &errWebauthnCredRepo{listErr: context.Canceled}
	svc := &adminService{
		users: users, webauthnCreds: webauthn, mfaTOTP: &adminMfaTOTPStub{},
		sessions: newSessionTestRepo(), memberships: &stubMembershipRepo{active: map[string]map[string]bool{}},
	}
	_, err := svc.GetUserByID(context.Background(), targetID)
	require.Error(t, err)
}

type countingAuditLogRepo struct {
	stubAuditLogRepo
	count int64
}

func (r *countingAuditLogRepo) Count(context.Context, repositories.AuditLogListFilters) (int64, error) {
	return r.count, nil
}

func TestAdminService_GetClientUsage_AuditCounts(t *testing.T) {
	clientID := "c1"
	clientRepo := newAdminClientTestRepo()
	clientRepo.clients[clientID] = &models.Client{BaseModel: models.BaseModel{ID: clientID}, ClientID: "app"}
	auditRepo := &countingAuditLogRepo{count: 3}
	svc := &adminService{clients: clientRepo, refreshTokens: newRefreshTokenTestRepo(), auditLogs: auditRepo}
	usage, err := svc.GetClientUsage(context.Background(), clientID)
	require.NoError(t, err)
	require.Equal(t, int64(3), usage.AuthorizeEvents30d)
	require.Equal(t, int64(3), usage.TokenIssueEvents30d)
}

func TestAdminService_ResetMFA_UserNotFound(t *testing.T) {
	svc, _, _, _, _ := newAdminUserTestService(&adminUserTestRepo{users: map[string]*models.User{}})
	require.Error(t, svc.ResetMFA(context.Background(), "actor", "missing"))
}

func TestAdminService_ResetPasskeys_UserNotFound(t *testing.T) {
	svc, _, _, _, _ := newAdminUserTestService(&adminUserTestRepo{users: map[string]*models.User{}})
	require.Error(t, svc.ResetPasskeys(context.Background(), "actor", "missing"))
}

// --- OIDC coverage ---

func TestOIDCService_Authorize_ConfidentialBadPKCEMethod(t *testing.T) {
	svc, _, clients, _, users, memberships := newOIDCTestService(t)
	u := users.seed("conf@example.com", "secret")
	memberships.active[u.ID] = map[string]bool{"tenant-1": true}
	clients.byClientID["conf"] = &models.Client{
		BaseModel: models.BaseModel{ID: "c1"}, TenantID: "tenant-1", ClientID: "conf",
		ClientSecret: "sec", RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes: pq.StringArray{"authorization_code"}, Scopes: pq.StringArray{"openid"},
	}
	_, err := svc.Authorize(context.Background(), u.ID, &dtos.AuthorizeQuery{
		ClientID: "conf", RedirectURI: "http://localhost/callback", ResponseType: "code",
		CodeChallenge: "abc", CodeChallengeMethod: "plain",
	})
	require.NotNil(t, err)
}

func TestOIDCService_Authorize_AuthCodeCreateError(t *testing.T) {
	cfg := testConfig()
	signer, err := auth.ProvideOIDCSigner(cfg)
	require.NoError(t, err)
	clients := &stubClientRepo{byClientID: map[string]*models.Client{}}
	authCodes := &errAuthCodeRepo{authCodeTestRepo: newAuthCodeTestRepo(), createErr: context.Canceled}
	users := newUserTestRepo()
	memberships := &stubMembershipRepo{active: map[string]map[string]bool{}}
	svc := ProvideOIDCService(cfg, signer, clients, authCodes, users, newRefreshTokenTestRepo(), memberships, &auditCapture{})

	u := users.seed("fail@example.com", "secret")
	memberships.active[u.ID] = map[string]bool{"tenant-1": true}
	clients.byClientID["app"] = &models.Client{
		BaseModel: models.BaseModel{ID: "c1"}, TenantID: "tenant-1", ClientID: "app",
		IsPublic: true, RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes: pq.StringArray{"authorization_code"}, Scopes: pq.StringArray{"openid"},
	}
	_, oauthErr := svc.Authorize(context.Background(), u.ID, &dtos.AuthorizeQuery{
		ClientID: "app", RedirectURI: "http://localhost/callback", ResponseType: "code",
		CodeChallenge: pkceChallenge("v"), CodeChallengeMethod: "S256",
	})
	require.NotNil(t, oauthErr)
}

func TestOIDCService_AuthorizationCodeToken_ClientValidation(t *testing.T) {
	svc, _, clients, authCodes, users, _ := newOIDCTestService(t)
	u := users.seed("clientval@example.com", "secret")
	clientRecordID := models.NewBaseModel().ID
	clients.byClientID["app"] = &models.Client{
		BaseModel: models.BaseModel{ID: clientRecordID}, TenantID: "t1", ClientID: "app",
		ClientSecret: "sec", RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes: pq.StringArray{"authorization_code"}, Scopes: pq.StringArray{"openid"},
	}
	code := "code-client-val"
	authCodes.byCode[code] = &models.AuthorizationCode{
		Code: code, TenantID: "t1", OAuthClientID: "app", UserID: u.ID,
		RedirectURI: "http://localhost/callback", ExpiresAt: time.Now().UTC().Add(time.Minute),
		ClientRecordID: &clientRecordID,
	}
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", "http://localhost/callback")

	_, tokenErr := svc.AuthorizationCodeToken(context.Background(), "t1", "", "wrong-secret", form)
	require.NotNil(t, tokenErr)

	form.Set("client_id", "other-app")
	_, tokenErr = svc.AuthorizationCodeToken(context.Background(), "t1", "other-app", "sec", form)
	require.NotNil(t, tokenErr)
}

func TestOIDCService_AuthorizationCodeToken_InactiveUser(t *testing.T) {
	svc, _, clients, authCodes, users, _ := newOIDCTestService(t)
	u := users.seed("inactive-oidc@example.com", "secret")
	u.Status = constants.UserStatusDisabled
	clientRecordID := models.NewBaseModel().ID
	clients.byClientID["app"] = &models.Client{
		BaseModel: models.BaseModel{ID: clientRecordID}, TenantID: "t1", ClientID: "app",
		ClientSecret: "sec", RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes: pq.StringArray{"authorization_code"}, Scopes: pq.StringArray{"openid"},
	}
	code := "code-inactive"
	authCodes.byCode[code] = &models.AuthorizationCode{
		Code: code, TenantID: "t1", OAuthClientID: "app", UserID: u.ID,
		Scope: "openid", RedirectURI: "http://localhost/callback",
		ExpiresAt: time.Now().UTC().Add(time.Minute), ClientRecordID: &clientRecordID,
	}
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", "http://localhost/callback")
	form.Set("client_id", "app")

	_, tokenErr := svc.AuthorizationCodeToken(context.Background(), "t1", "app", "sec", form)
	require.NotNil(t, tokenErr)
	require.Equal(t, constants.OAuthInvalidGrant, tokenErr.Code)
}

func TestOIDCService_AuthorizationCodeToken_DeleteCodeError(t *testing.T) {
	cfg := testConfig()
	signer, err := auth.ProvideOIDCSigner(cfg)
	require.NoError(t, err)
	clients := &stubClientRepo{byClientID: map[string]*models.Client{}}
	authCodes := &errAuthCodeRepo{authCodeTestRepo: newAuthCodeTestRepo(), deleteErr: context.Canceled}
	users := newUserTestRepo()
	svc := ProvideOIDCService(cfg, signer, clients, authCodes, users, newRefreshTokenTestRepo(), &stubMembershipRepo{active: map[string]map[string]bool{}}, &auditCapture{})

	u := users.seed("del@example.com", "secret")
	clientRecordID := models.NewBaseModel().ID
	clients.byClientID["app"] = &models.Client{
		BaseModel: models.BaseModel{ID: clientRecordID}, TenantID: "t1", ClientID: "app",
		IsPublic: true, RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes: pq.StringArray{"authorization_code"}, Scopes: pq.StringArray{"openid"},
	}
	code := "code-del-err"
	authCodes.byCode[code] = &models.AuthorizationCode{
		Code: code, TenantID: "t1", OAuthClientID: "app", UserID: u.ID,
		RedirectURI: "http://localhost/callback", ExpiresAt: time.Now().UTC().Add(time.Minute),
		ClientRecordID: &clientRecordID,
	}
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", "http://localhost/callback")
	form.Set("client_id", "app")

	_, tokenErr := svc.AuthorizationCodeToken(context.Background(), "t1", "app", "", form)
	require.NotNil(t, tokenErr)
}

func TestOIDCService_AuthorizationCodeToken_PersistRefreshError(t *testing.T) {
	cfg := testConfig()
	signer, err := auth.ProvideOIDCSigner(cfg)
	require.NoError(t, err)
	clients := &stubClientRepo{byClientID: map[string]*models.Client{}}
	authCodes := newAuthCodeTestRepo()
	refresh := &errRefreshCreateRepo{refreshTokenTestRepo: newRefreshTokenTestRepo(), createErr: context.Canceled}
	users := newUserTestRepo()
	svc := ProvideOIDCService(cfg, signer, clients, authCodes, users, refresh, &stubMembershipRepo{active: map[string]map[string]bool{}}, &auditCapture{})

	u := users.seed("refresh@example.com", "secret")
	clientRecordID := models.NewBaseModel().ID
	clients.byClientID["app"] = &models.Client{
		BaseModel: models.BaseModel{ID: clientRecordID}, TenantID: "t1", ClientID: "app",
		IsPublic: true, RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes: pq.StringArray{"authorization_code"}, Scopes: pq.StringArray{"openid"},
	}
	code := "code-refresh-err"
	authCodes.byCode[code] = &models.AuthorizationCode{
		Code: code, TenantID: "t1", OAuthClientID: "app", UserID: u.ID,
		Scope: "openid", RedirectURI: "http://localhost/callback",
		ExpiresAt: time.Now().UTC().Add(time.Minute), ClientRecordID: &clientRecordID,
	}
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", "http://localhost/callback")
	form.Set("client_id", "app")

	_, tokenErr := svc.AuthorizationCodeToken(context.Background(), "t1", "app", "", form)
	require.NotNil(t, tokenErr)
}

func TestOIDCService_UserInfo_UserNotFound(t *testing.T) {
	svc, signer, clients, _, _, _ := newOIDCTestService(t)
	clientID := "app"
	clients.byClientID[clientID] = &models.Client{ClientID: clientID}
	access, _, err := signer.SignAccessTokenOIDC("missing-user", clientID, "openid", clientID)
	require.NoError(t, err)
	_, tokenErr := svc.UserInfo(context.Background(), access)
	require.NotNil(t, tokenErr)
}

// --- MFA coverage ---

func TestMFAService_SetupTOTP_HasActiveError(t *testing.T) {
	totpRepo := &errMFATOTPActiveRepo{mfaTOTPTestRepo: newMFATOTPTestRepo(), activeErr: context.Canceled}
	cfg := testConfig()
	svc := ProvideMFAService(cfg, totpRepo, newMFARecoveryTestRepo(), auth.NewEphemeralStore(newMemCache(), cfg), &auditCapture{})
	_, _, err := svc.SetupTOTP(context.Background(), "u1", "u@example.com")
	require.Error(t, err)
}

func TestMFAService_VerifyTOTPEnrollment_InvalidCode(t *testing.T) {
	svc, totpRepo, _, _ := newMFATestService(t)
	userID := "enroll-bad"
	secret, _, err := svc.SetupTOTP(context.Background(), userID, "u@example.com")
	require.NoError(t, err)
	require.NotNil(t, totpRepo.byUser[userID])
	err = svc.VerifyTOTPEnrollment(context.Background(), userID, "000000")
	require.Error(t, err)
	_ = secret
}

func TestMFAService_TryRecoveryCode_NoMatch(t *testing.T) {
	svc, totpRepo, recoveryRepo, _ := newMFATestService(t)
	userID := "recovery-miss"
	cfg := testConfig()
	enc, _ := crypto.EncryptMFASecret(cfg.MFAEncryptionKey, "JBSWY3DPEHPK3PXP")
	totpRepo.byUser[userID] = &models.UserMFATOTP{UserID: userID, SecretEncrypted: enc, Enabled: true}
	recoveryRepo.byUser[userID] = []models.UserMFARecoveryCode{{BaseModel: models.NewBaseModel(), UserID: userID, CodeHash: "other"}}

	ticket, _, err := svc.CreateLoginTicket(context.Background(), auth.MFAPendingPayload{UserID: userID, TenantID: "t1"})
	require.NoError(t, err)
	_, err = svc.VerifyLoginChallenge(context.Background(), ticket, "ZZZZZZZZZZ")
	require.Error(t, err)
}

func TestMFAService_RegenerateRecoveryCodes_DefaultCount(t *testing.T) {
	cfg := testConfig()
	cfg.MFARecoveryCodeCount = 0
	totpRepo := newMFATOTPTestRepo()
	userID := "default-count"
	enc, _ := crypto.EncryptMFASecret(cfg.MFAEncryptionKey, "JBSWY3DPEHPK3PXP")
	totpRepo.byUser[userID] = &models.UserMFATOTP{UserID: userID, SecretEncrypted: enc, Enabled: true}
	svc := ProvideMFAService(cfg, totpRepo, newMFARecoveryTestRepo(), auth.NewEphemeralStore(newMemCache(), cfg), &auditCapture{})
	codes, err := svc.RegenerateRecoveryCodes(context.Background(), userID)
	require.NoError(t, err)
	require.Len(t, codes, 10)
}

// --- user coverage ---

func TestUserService_AuthenticateUser_RepoError(t *testing.T) {
	users := &errUserGetRepo{userTestRepo: newUserTestRepo(), getErr: context.Canceled}
	cfg := testConfig()
	tokenSvc, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	require.NoError(t, err)
	memberships := &stubMembershipRepo{active: map[string]map[string]bool{}}
	svc := ProvideUserService(users, memberships, newRefreshTokenTestRepo(), ProvideTenantContextService(cfg, &stubClientRepo{}, &stubTenantRepo{}, memberships), cfg, tokenSvc, &auditCapture{})
	_, err = svc.AuthenticateUser(context.Background(), &dtos.LoginRequest{Email: "x@example.com", Password: "p"})
	require.Error(t, err)
}

type errMembershipCreateRepo struct {
	*stubMembershipRepo
	createErr error
}

func (r *errMembershipCreateRepo) Create(context.Context, *models.TenantMembership) error {
	return r.createErr
}

func TestUserService_Register_MembershipCreateError(t *testing.T) {
	users := newUserTestRepo()
	memberships := &errMembershipCreateRepo{
		stubMembershipRepo: &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}},
		createErr:          context.Canceled,
	}
	cfg := testConfig()
	tokenSvc, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	require.NoError(t, err)
	tenantCtx := ProvideTenantContextService(cfg, &stubClientRepo{}, &stubTenantRepo{}, memberships)
	svc := ProvideUserService(users, memberships, newRefreshTokenTestRepo(), tenantCtx, cfg, tokenSvc, &auditCapture{})

	_, err = svc.Register(context.Background(), &dtos.RegisterRequest{
		Email: "new@example.com", Password: "password123", TenantID: cfg.DefaultTenantID,
	}, "")
	require.Error(t, err)
}

type errMembershipExistsRepo struct {
	existsErr error
}

func (r *errMembershipExistsRepo) Create(context.Context, *models.TenantMembership) error { return nil }
func (r *errMembershipExistsRepo) GetActive(context.Context, string, string) (*models.TenantMembership, error) {
	return nil, nil
}
func (r *errMembershipExistsRepo) ExistsActive(context.Context, string, string) (bool, error) {
	return false, r.existsErr
}
func (r *errMembershipExistsRepo) ListByUserID(context.Context, string) ([]models.TenantMembership, error) {
	return nil, nil
}
func (r *errMembershipExistsRepo) ListByUserIDPaginated(context.Context, string, *dtos.PageableRequest) (*dtos.DataResponse[models.TenantMembership], error) {
	return nil, nil
}
func (r *errMembershipExistsRepo) ListByTenantID(context.Context, string) ([]models.TenantMembership, error) {
	return nil, nil
}
func (r *errMembershipExistsRepo) ListByTenantIDPaginated(context.Context, string, *dtos.PageableRequest) (*dtos.DataResponse[models.TenantMembership], error) {
	return nil, nil
}
func (r *errMembershipExistsRepo) CountByTenantID(context.Context, string) (int64, error) {
	return 0, nil
}
func (r *errMembershipExistsRepo) Delete(context.Context, string, string) error { return nil }

// --- federation coverage ---

func TestFederationService_BuildAuthorizeProviderRedirectError(t *testing.T) {
	provider := &fedProviderRedirectErr{
		mockFedProvider: mockFedProvider{id: "google", configured: true},
		redirectErr:     errors.ExternalServiceError("redirect failed", nil),
	}
	tip := &fedTipStub{configured: map[string]bool{"tenant-1:google": true}, enabled: map[string]bool{"tenant-1:google": true}}
	svc := newFederationTestService(map[string]FederationIdentityProvider{"google": provider}, tip)
	_, err := svc.BuildAuthorizeRedirectURL(context.Background(), "google", "/", "tenant-1")
	require.Error(t, err)
}

func TestFederationService_BuildAuthorizeConfiguredCheckError(t *testing.T) {
	provider := &fedProviderConfigErr{mockFedProvider: mockFedProvider{id: "google"}}
	svc := newFederationTestService(map[string]FederationIdentityProvider{"google": provider}, &fedTipStub{})
	_, err := svc.BuildAuthorizeRedirectURL(context.Background(), "google", "/", "tenant-1")
	require.Error(t, err)
}

func TestFederationService_ListAvailableProviders_SkipsUnconfigured(t *testing.T) {
	provider := &mockFedProvider{id: "google", name: "Google", configured: false}
	tip := &fedTipStub{configured: map[string]bool{"tenant-1:google": false}}
	svc := newFederationTestService(map[string]FederationIdentityProvider{"google": provider}, tip)
	out, err := svc.ListAvailableProviders(context.Background(), "tenant-1")
	require.NoError(t, err)
	require.Empty(t, out)
}

func TestFederationService_ResolveFedRepoError(t *testing.T) {
	fedRepo := &errFedIdentityRepo{getErr: context.Canceled}
	svc := &federationService{userRepo: newUserTestRepo(), fedRepo: fedRepo, membershipRepo: &stubMembershipRepo{active: map[string]map[string]bool{}}}
	_, err := svc.federationResolveOrProvisionUser(context.Background(), "t1", "google", &domains.OIDCUserClaims{Sub: "s", Email: "e@example.com"})
	require.Error(t, err)
}

type errFedIdentityRepo struct {
	getErr error
}

func (r *errFedIdentityRepo) GetByProviderSubject(context.Context, string, string) (*models.FederatedIdentity, error) {
	return nil, r.getErr
}
func (r *errFedIdentityRepo) ListByUserID(context.Context, string) ([]models.FederatedIdentity, error) {
	return nil, nil
}
func (r *errFedIdentityRepo) Create(context.Context, *models.FederatedIdentity) error { return nil }

// --- platform admin bootstrap ---

func TestPlatformAdminBootstrap_CountError(t *testing.T) {
	repo := &errBootstrapCountRepo{stubBootstrapUserRepo: &stubBootstrapUserRepo{users: map[string]*models.User{}, byEmail: map[string]*models.User{}, promoted: map[string]bool{}}}
	b := ProvidePlatformAdminBootstrap(&config.Config{BootstrapAdminEmail: "a@example.com", BootstrapAdminPassword: "password123"}, repo, &stubBootstrapMembershipRepo{})
	require.Error(t, b.Run(context.Background()))
}

type errBootstrapCountRepo struct {
	*stubBootstrapUserRepo
}

func (r *errBootstrapCountRepo) CountPlatformAdmins(context.Context) (int64, error) {
	return 0, context.Canceled
}

func TestPlatformAdminBootstrap_GetUserNonNotFoundError(t *testing.T) {
	repo := &errBootstrapGetRepo{stubBootstrapUserRepo: &stubBootstrapUserRepo{users: map[string]*models.User{}, byEmail: map[string]*models.User{}, promoted: map[string]bool{}}}
	b := ProvidePlatformAdminBootstrap(&config.Config{
		DefaultTenantID: "t1", BootstrapAdminEmail: "a@example.com", BootstrapAdminPassword: "password123",
	}, repo, &stubBootstrapMembershipRepo{})
	require.Error(t, b.Run(context.Background()))
}

type errBootstrapGetRepo struct {
	*stubBootstrapUserRepo
}

func (r *errBootstrapGetRepo) GetByEmailLower(context.Context, string) (*models.User, error) {
	return nil, context.Canceled
}

// --- session coverage ---

func TestSessionService_Invalidate_Error(t *testing.T) {
	svc := ProvideSessionService(testConfig(), brokenSessionRepo{}, &auditCapture{})
	require.Error(t, svc.Invalidate(context.Background(), "sess-1"))
}

// --- tenant context ---

func TestTenantContextService_ValidateMembership_Error(t *testing.T) {
	memberships := &errMembershipExistsRepo{existsErr: context.Canceled}
	svc := ProvideTenantContextService(testConfig(), &stubClientRepo{}, &stubTenantRepo{}, memberships)
	require.Error(t, svc.ValidateMembership(context.Background(), "u1", "t1"))
}

func TestUserService_Login_Success(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("loginfn@example.com", "secret123")
	tenantID := "tenant-login-fn"
	memberships := &stubMembershipRepo{
		byUser: map[string][]models.TenantMembership{
			u.ID: {{UserID: u.ID, TenantID: tenantID, Status: constants.TenantMembershipStatusActive}},
		},
		active: map[string]map[string]bool{u.ID: {tenantID: true}},
	}
	svc := newUserTestService(t, users, memberships)
	resp, sel, err := svc.Login(context.Background(), &dtos.LoginRequest{Email: "loginfn@example.com", Password: "secret123", TenantID: tenantID}, "")
	require.NoError(t, err)
	require.Nil(t, sel)
	require.NotEmpty(t, resp.AccessToken)
}

func TestUserService_IssueAccessAndRefresh_Error(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("issue@example.com", "secret")
	memberships := &stubMembershipRepo{active: map[string]map[string]bool{}}
	refresh := &errRefreshCreateRepo{refreshTokenTestRepo: newRefreshTokenTestRepo(), createErr: context.Canceled}
	cfg := testConfig()
	tokenSvc, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	require.NoError(t, err)
	svc := ProvideUserService(users, memberships, refresh, ProvideTenantContextService(cfg, &stubClientRepo{}, &stubTenantRepo{}, memberships), cfg, tokenSvc, &auditCapture{})
	_, err = svc.IssueTokensForUser(context.Background(), u, "tenant-1")
	require.Error(t, err)
}

func TestUserService_BuildTenantSelection_ListError(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("selerr@example.com", "secret")
	memberships := &errMembershipListRepo{
		stubMembershipRepo: &stubMembershipRepo{
			byUser: map[string][]models.TenantMembership{
				u.ID: {
					{UserID: u.ID, TenantID: "t1", Status: constants.TenantMembershipStatusActive},
					{UserID: u.ID, TenantID: "t2", Status: constants.TenantMembershipStatusActive},
				},
			},
			active: map[string]map[string]bool{u.ID: {"t1": true, "t2": true}},
		},
		listErr: context.Canceled,
	}
	cfg := testConfig()
	cfg.DefaultTenantID = ""
	tokenSvc, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	require.NoError(t, err)
	svc := ProvideUserService(users, memberships, newRefreshTokenTestRepo(), ProvideTenantContextService(cfg, &stubClientRepo{}, &stubTenantRepo{}, memberships), cfg, tokenSvc, &auditCapture{})
	_, sel, err := svc.CompleteAuth(context.Background(), u, TenantResolveInput{UserID: u.ID})
	require.Error(t, err)
	require.Nil(t, sel)
}

type errMembershipListRepo struct {
	*stubMembershipRepo
	listErr error
}

func (r *errMembershipListRepo) ListByUserID(ctx context.Context, userID string) ([]models.TenantMembership, error) {
	if r.listErr != nil {
		return nil, r.listErr
	}
	return r.stubMembershipRepo.ListByUserID(ctx, userID)
}
func (r *errMembershipListRepo) ListByUserIDPaginated(ctx context.Context, userID string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.TenantMembership], error) {
	if r.listErr != nil {
		return nil, r.listErr
	}
	return r.stubMembershipRepo.ListByUserIDPaginated(ctx, userID, pr)
}

func TestAdminService_DisableUser_SessionInvalidateError(t *testing.T) {
	actorID := "actor-1"
	targetID := "target-1"
	users := &adminUserTestRepo{users: map[string]*models.User{
		targetID: {BaseModel: models.BaseModel{ID: targetID}, Email: "t@example.com", Status: constants.UserStatusActive},
	}}
	sessions := &adminSessionStub{}
	svc, _, _, _, _ := newAdminUserTestService(users)
	svc.sessions = nil
	svc.sessionSvc = &errSessionInvalidateStub{}
	err := svc.DisableUser(context.Background(), actorID, targetID)
	require.Error(t, err)
	_ = sessions
}

type errSessionInvalidateStub struct {
	adminSessionStub
}

func (errSessionInvalidateStub) InvalidateAllForUser(context.Context, string) error {
	return context.Canceled
}

func TestAdminService_GetClientUsage_CountError(t *testing.T) {
	clientID := "c1"
	clientRepo := newAdminClientTestRepo()
	clientRepo.clients[clientID] = &models.Client{BaseModel: models.BaseModel{ID: clientID}, ClientID: "app"}
	auditRepo := &errAuditCountRepo{}
	svc := &adminService{clients: clientRepo, refreshTokens: newRefreshTokenTestRepo(), auditLogs: auditRepo}
	_, err := svc.GetClientUsage(context.Background(), clientID)
	require.Error(t, err)
}

type errAuditCountRepo struct {
	stubAuditLogRepo
}

func (errAuditCountRepo) Count(context.Context, repositories.AuditLogListFilters) (int64, error) {
	return 0, context.Canceled
}

func TestFederationOIDC_AuthorizeRedirectURL_Success(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	issuer := srv.URL
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"issuer":%q,"authorization_endpoint":%q,"token_endpoint":%q,"jwks_uri":%q}`,
			issuer, issuer+"/auth", issuer+"/token", issuer+"/jwks")
	})
	mux.HandleFunc("/jwks", func(w http.ResponseWriter, r *http.Request) {
		n := base64.RawURLEncoding.EncodeToString(key.N.Bytes())
		_, _ = fmt.Fprintf(w, `{"keys":[{"kty":"RSA","kid":"test","n":%q,"e":"AQAB"}]}`, n)
	})
	spec := constants.IdentityProviderSpec{ID: "mock-idp", DisplayName: "Mock", IssuerURL: issuer, Scopes: []string{"openid"}}
	enc, err := crypto.EncryptMFASecret(testConfig().MFAEncryptionKey, "client-secret")
	require.NoError(t, err)
	tip := &adminIdpTipStub{byKey: map[string]*models.TenantIdentityProvider{
		"tenant-1:mock-idp": {TenantID: "tenant-1", Provider: "mock-idp", OAuthClientID: "cid", OAuthClientSecretEncrypted: enc},
	}}
	provider := newOIDCFederationProvider(testConfig(), tip, spec).(*oidcFederationProvider)
	url, err := provider.AuthorizeRedirectURL(context.Background(), "tenant-1", "state-1", "nonce-1")
	require.NoError(t, err)
	require.Contains(t, url, "client_id=cid")
}

func TestFederationOIDC_ExchangeAuthorizationCode_NoIDToken(t *testing.T) {
	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	issuer := srv.URL
	mux.HandleFunc("/.well-known/openid-configuration", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, `{"issuer":%q,"authorization_endpoint":%q,"token_endpoint":%q,"jwks_uri":%q}`,
			issuer, issuer+"/auth", issuer+"/token", issuer+"/jwks")
	})
	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"access_token":"at","token_type":"Bearer"}`)
	})
	spec := constants.IdentityProviderSpec{ID: "mock-idp", DisplayName: "Mock", IssuerURL: issuer, Scopes: []string{"openid"}}
	enc, err := crypto.EncryptMFASecret(testConfig().MFAEncryptionKey, "secret")
	require.NoError(t, err)
	tip := &adminIdpTipStub{byKey: map[string]*models.TenantIdentityProvider{
		"tenant-1:mock-idp": {TenantID: "tenant-1", Provider: "mock-idp", OAuthClientID: "cid", OAuthClientSecretEncrypted: enc},
	}}
	provider := newOIDCFederationProvider(testConfig(), tip, spec).(*oidcFederationProvider)
	_, err = provider.ExchangeAuthorizationCode(context.Background(), "tenant-1", "code", "nonce")
	require.Error(t, err)
}

func TestAdminService_UpdateClient_PublicToConfidential(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	clientRepo := newAdminClientTestRepo()
	id := "client-pub"
	clientRepo.clients[id] = &models.Client{
		BaseModel: models.BaseModel{ID: id}, TenantID: "tenant-1", ClientID: "pub-app",
		Name: "Pub", IsPublic: true, RedirectUris: pq.StringArray{"http://localhost/callback"},
		GrantTypes: pq.StringArray{"authorization_code"}, Scopes: pq.StringArray{"openid"},
	}
	svc := newAdminClientTestService(tenantRepo, clientRepo)
	isPublic := false
	updated, err := svc.UpdateClient(context.Background(), id, &dtos.AdminUpdateClientRequest{IsPublic: &isPublic})
	require.NoError(t, err)
	require.False(t, updated.IsPublic)
	require.NotEmpty(t, clientRepo.clients[id].ClientSecret)
}

func TestAdminService_UpdateClient_EmptyName(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	clientRepo := newAdminClientTestRepo()
	id := "client-name"
	clientRepo.clients[id] = &models.Client{BaseModel: models.BaseModel{ID: id}, TenantID: "t1", ClientID: "app", Name: "App"}
	svc := newAdminClientTestService(tenantRepo, clientRepo)
	empty := "  "
	_, err := svc.UpdateClient(context.Background(), id, &dtos.AdminUpdateClientRequest{Name: &empty})
	require.Error(t, err)
}

func TestUserService_ListMemberships_Error(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("listm@example.com", "secret")
	memberships := &errMembershipListRepo{
		stubMembershipRepo: &stubMembershipRepo{active: map[string]map[string]bool{}},
		listErr:            context.Canceled,
	}
	cfg := testConfig()
	tokenSvc, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	require.NoError(t, err)
	svc := ProvideUserService(users, memberships, newRefreshTokenTestRepo(), ProvideTenantContextService(cfg, &stubClientRepo{}, &stubTenantRepo{}, memberships), cfg, tokenSvc, &auditCapture{})
	_, err = svc.ListMemberships(context.Background(), u.ID)
	require.Error(t, err)
}

func TestFederationService_EnsureMembership_Error(t *testing.T) {
	memberships := &errMembershipCreateRepo{
		stubMembershipRepo: &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}},
		createErr:          context.Canceled,
	}
	svc := &federationService{membershipRepo: memberships}
	require.Error(t, svc.ensureMembership(context.Background(), "u1", "t1"))
}

func TestPlatformAdminBootstrap_EnsureMembershipError(t *testing.T) {
	repo := &stubBootstrapUserRepo{users: map[string]*models.User{}, byEmail: map[string]*models.User{}, promoted: map[string]bool{}}
	memberships := &errBootstrapMembershipCreateRepo{}
	b := ProvidePlatformAdminBootstrap(&config.Config{
		DefaultTenantID: "t1", BootstrapAdminEmail: "admin@example.com", BootstrapAdminPassword: "password123",
	}, repo, memberships)
	require.Error(t, b.Run(context.Background()))
}

type errBootstrapMembershipCreateRepo struct {
	stubBootstrapMembershipRepo
}

func (errBootstrapMembershipCreateRepo) Create(context.Context, *models.TenantMembership) error {
	return context.Canceled
}

func TestMFAService_CreateLoginTicket_PutError(t *testing.T) {
	cfg := testConfig()
	brokenCache := &brokenMemCache{}
	svc := ProvideMFAService(cfg, newMFATOTPTestRepo(), newMFARecoveryTestRepo(), auth.NewEphemeralStore(brokenCache, cfg), &auditCapture{})
	_, _, err := svc.CreateLoginTicket(context.Background(), auth.MFAPendingPayload{UserID: "u1", TenantID: "t1"})
	require.Error(t, err)
}

type brokenMemCache struct{}

func (brokenMemCache) Get(context.Context, string) (string, error) { return "", context.Canceled }
func (brokenMemCache) Set(context.Context, string, string, time.Duration) error {
	return context.Canceled
}
func (brokenMemCache) Delete(context.Context, string) error         { return nil }
func (brokenMemCache) Exists(context.Context, string) (bool, error) { return false, nil }
func (brokenMemCache) Close() error                                 { return nil }

func TestOIDCService_VerifyPKCE_EmptyVerifier(t *testing.T) {
	require.False(t, verifyPKCE("", "challenge", "S256"))
	require.False(t, verifyPKCE("v", "challenge", "plain"))
}
