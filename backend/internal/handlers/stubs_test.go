package handlers

import (
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"
	"github.com/gateforge-iam/gateforge-iam/internal/services"

	"github.com/go-playground/validator/v10"
)

type stubAuditService struct{}

func (stubAuditService) Record(context.Context, domains.AuditRecordParams) {}

type stubAdminAuditService struct{}

func (stubAdminAuditService) GetStats(context.Context) (*dtos.AdminStatsResponse, error) {
	return &dtos.AdminStatsResponse{}, nil
}
func (stubAdminAuditService) ListUsers(context.Context, string, string, *dtos.PageableRequest) ([]*dtos.AdminUserResponse, *dtos.Pageable, error) {
	return nil, nil, nil
}
func (stubAdminAuditService) ListTenants(context.Context, *dtos.PageableRequest) ([]*dtos.AdminTenantResponse, *dtos.Pageable, error) {
	return nil, nil, nil
}
func (stubAdminAuditService) GetTenantByID(context.Context, string) (*dtos.AdminTenantResponse, error) {
	return nil, nil
}
func (stubAdminAuditService) CreateTenant(context.Context, *dtos.AdminCreateTenantRequest) (*dtos.AdminTenantResponse, error) {
	return nil, nil
}
func (stubAdminAuditService) UpdateTenant(context.Context, string, *dtos.AdminUpdateTenantRequest) (*dtos.AdminTenantResponse, error) {
	return nil, nil
}
func (stubAdminAuditService) DeleteTenant(context.Context, string) error { return nil }
func (stubAdminAuditService) ListTenantMembers(context.Context, string, *dtos.PageableRequest) ([]*dtos.AdminTenantMemberResponse, *dtos.Pageable, error) {
	return nil, nil, nil
}
func (stubAdminAuditService) ListClients(context.Context, string, *dtos.PageableRequest) ([]*dtos.AdminClientResponse, *dtos.Pageable, error) {
	return nil, nil, nil
}
func (stubAdminAuditService) GetClientByID(context.Context, string) (*dtos.AdminClientResponse, error) {
	return nil, nil
}
func (stubAdminAuditService) CreateClient(context.Context, *dtos.AdminCreateClientRequest) (*dtos.AdminCreateClientResponse, error) {
	return nil, nil
}
func (stubAdminAuditService) UpdateClient(context.Context, string, *dtos.AdminUpdateClientRequest) (*dtos.AdminClientResponse, error) {
	return nil, nil
}
func (stubAdminAuditService) DeleteClient(context.Context, string) error { return nil }
func (stubAdminAuditService) ListIdentityProviders(context.Context, string, *dtos.PageableRequest) ([]*dtos.AdminIdentityProviderResponse, *dtos.Pageable, error) {
	return nil, nil, nil
}
func (stubAdminAuditService) ConfigureIdentityProvider(context.Context, string, string, *dtos.PatchIdentityProviderRequest, constants.AuditActorType) error {
	return nil
}
func (stubAdminAuditService) AddMemberByEmail(context.Context, string, string, string) error {
	return nil
}
func (stubAdminAuditService) RemoveMember(context.Context, string, string) error { return nil }
func (stubAdminAuditService) ListAuditLogs(_ context.Context, _ dtos.AdminAuditLogListParams, _ *dtos.PageableRequest) ([]*dtos.AdminAuditLogResponse, *dtos.Pageable, error) {
	return []*dtos.AdminAuditLogResponse{
		{ID: "00000000-0000-4000-8000-000000000010", Action: constants.AuditActionAuthLogin, Result: string(constants.AuditResultSuccess), ActorType: string(constants.AuditActorTypeUser)},
	}, &dtos.Pageable{Page: 1, PageSize: 20, Total: 1}, nil
}
func (stubAdminAuditService) GetUserByID(context.Context, string) (*dtos.AdminUserDetailResponse, error) {
	return nil, nil
}
func (stubAdminAuditService) DisableUser(context.Context, string, string) error     { return nil }
func (stubAdminAuditService) ForceLogoutUser(context.Context, string, string) error { return nil }
func (stubAdminAuditService) ResetPasskeys(context.Context, string, string) error   { return nil }
func (stubAdminAuditService) ResetMFA(context.Context, string, string) error        { return nil }
func (stubAdminAuditService) GetClientUsage(context.Context, string) (*dtos.AdminClientUsageResponse, error) {
	return nil, nil
}
func (stubAdminAuditService) ListLoginHistory(_ context.Context, _ dtos.AdminLoginHistoryListParams, _ *dtos.PageableRequest) ([]*dtos.AdminAuditLogResponse, *dtos.Pageable, error) {
	return nil, nil, nil
}

type stubAuthUserService struct {
	registerErr   error
	authErr       error
	completeLogin *dtos.LoginResponse
	selection     *dtos.TenantSelectionResponse
	completeErr   error
	revokeErr     error
	hasUser       *models.User
	mfaOnLogin    bool
}

func (s *stubAuthUserService) testUser() *models.User {
	if s.hasUser != nil {
		return s.hasUser
	}
	return &models.User{
		BaseModel: models.NewBaseModel(),
		Email:     "user@example.com",
		FirstName: "Test",
		LastName:  "User",
	}
}

func (s *stubAuthUserService) Register(_ context.Context, _ *dtos.RegisterRequest, _ string) (*models.User, error) {
	if s.registerErr != nil {
		return nil, s.registerErr
	}
	return s.testUser(), nil
}
func (s *stubAuthUserService) AuthenticateUser(_ context.Context, _ *dtos.LoginRequest) (*models.User, error) {
	if s.authErr != nil {
		return nil, s.authErr
	}
	return s.testUser(), nil
}
func (s *stubAuthUserService) CompleteAuth(_ context.Context, u *models.User, _ services.TenantResolveInput) (*dtos.LoginResponse, *dtos.TenantSelectionResponse, error) {
	if s.completeErr != nil {
		return nil, nil, s.completeErr
	}
	if s.selection != nil {
		return nil, s.selection, nil
	}
	if s.completeLogin != nil {
		return s.completeLogin, nil, nil
	}
	return &dtos.LoginResponse{
		AccessToken:    "access-token",
		RefreshToken:   "refresh-token",
		TokenType:      "Bearer",
		ExpiresIn:      3600,
		ActiveTenantID: testTenantID,
	}, nil, nil
}
func (s *stubAuthUserService) SelectTenant(_ context.Context, _ *dtos.TenantSelectRequest) (*dtos.LoginResponse, error) {
	return &dtos.LoginResponse{AccessToken: "access-token", TokenType: "Bearer", ActiveTenantID: testTenantID}, nil
}
func (s *stubAuthUserService) SwitchTenant(_ context.Context, _, _ string) (*dtos.LoginResponse, error) {
	return &dtos.LoginResponse{AccessToken: "switched-token", TokenType: "Bearer", ActiveTenantID: testTenantID}, nil
}
func (s *stubAuthUserService) IssueTokensForUser(_ context.Context, _ *models.User, _ string) (*dtos.LoginResponse, error) {
	return &dtos.LoginResponse{AccessToken: "mfa-access", TokenType: "Bearer", ActiveTenantID: testTenantID}, nil
}
func (s *stubAuthUserService) Login(_ context.Context, _ *dtos.LoginRequest, _ string) (*dtos.LoginResponse, *dtos.TenantSelectionResponse, error) {
	return nil, nil, nil
}
func (s *stubAuthUserService) Refresh(_ context.Context, _ *dtos.RefreshTokenRequest) (*dtos.LoginResponse, error) {
	return &dtos.LoginResponse{AccessToken: "refreshed-access", RefreshToken: "refreshed-refresh", TokenType: "Bearer"}, nil
}
func (s *stubAuthUserService) GetOneByID(_ context.Context, _ string) (*models.User, error) {
	return s.testUser(), nil
}
func (s *stubAuthUserService) UpdateProfile(_ context.Context, _ string, req *dtos.UpdateProfileRequest) (*models.User, error) {
	u := s.testUser()
	if req.FirstName != nil {
		u.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		u.LastName = *req.LastName
	}
	return u, nil
}
func (s *stubAuthUserService) ListMemberships(_ context.Context, _ string) ([]dtos.TenantSummary, error) {
	return []dtos.TenantSummary{{ID: testTenantID, Name: "Default", Role: "member"}}, nil
}
func (s *stubAuthUserService) ListMembershipsPaginated(_ context.Context, _ string, _ *dtos.PageableRequest) ([]dtos.TenantSummary, *dtos.Pageable, error) {
	return []dtos.TenantSummary{{ID: testTenantID, Name: "Default", Role: "member"}}, &dtos.Pageable{Page: 1, PageSize: 20, Total: 1}, nil
}
func (s *stubAuthUserService) RevokeAllRefreshTokensForUser(_ context.Context, _ string) error {
	return s.revokeErr
}

type stubAuthSessionService struct {
	session     *models.Session
	getSessionErr error
	createSID     string
	createTTL     time.Duration
	createErr     error
	invalidateErr error
}

func (s *stubAuthSessionService) Create(_ context.Context, _, _, _, _ string, _ bool) (string, time.Duration, error) {
	if s.createErr != nil {
		return "", 0, s.createErr
	}
	sid := s.createSID
	if sid == "" {
		sid = "session-id-123"
	}
	ttl := s.createTTL
	if ttl == 0 {
		ttl = time.Hour
	}
	return sid, ttl, nil
}
func (s *stubAuthSessionService) GetSession(_ context.Context, _ string) (*models.Session, error) {
	if s.getSessionErr != nil {
		return nil, s.getSessionErr
	}
	if s.session != nil {
		return s.session, nil
	}
	return &models.Session{BaseModel: models.NewBaseModel(), UserID: testUserID, TenantID: testTenantID}, nil
}
func (s *stubAuthSessionService) GetUserID(_ context.Context, _ string) (string, error) {
	return testUserID, nil
}
func (s *stubAuthSessionService) Invalidate(context.Context, string) error { return nil }
func (s *stubAuthSessionService) InvalidateAllForUser(_ context.Context, _ string) error {
	return s.invalidateErr
}
func (s *stubAuthSessionService) BrowserSessionTTL(bool) time.Duration { return time.Hour }

type stubAuthMFAService struct {
	hasMFA         bool
	hasMFAErr      error
	ticket         string
	ticketExp      int64
	ticketErr      error
	verifyChallengeErr error
}

func (s *stubAuthMFAService) HasActiveMFA(_ context.Context, _ string) (bool, error) {
	return s.hasMFA, s.hasMFAErr
}
func (s *stubAuthMFAService) CreateLoginTicket(_ context.Context, _ auth.MFAPendingPayload) (string, int64, error) {
	if s.ticketErr != nil {
		return "", 0, s.ticketErr
	}
	ticket := s.ticket
	if ticket == "" {
		ticket = "mfa-ticket-123"
	}
	exp := s.ticketExp
	if exp == 0 {
		exp = 600
	}
	return ticket, exp, nil
}
func (s *stubAuthMFAService) SetupTOTP(_ context.Context, _, _ string) (string, string, error) {
	return "SECRET", "otpauth://totp/test", nil
}
func (s *stubAuthMFAService) VerifyTOTPEnrollment(context.Context, string, string) error { return nil }
func (s *stubAuthMFAService) RegenerateRecoveryCodes(_ context.Context, _ string) ([]string, error) {
	return []string{"code-1", "code-2"}, nil
}
func (s *stubAuthMFAService) VerifyLoginChallenge(_ context.Context, _, _ string) (*auth.MFAPendingPayload, error) {
	if s.verifyChallengeErr != nil {
		return nil, s.verifyChallengeErr
	}
	return &auth.MFAPendingPayload{UserID: testUserID, TenantID: testTenantID}, nil
}

type stubAuthFederationService struct {
	providers []services.ProviderAvailability
	listErr   error
}

func (s *stubAuthFederationService) BuildAuthorizeRedirectURL(context.Context, string, string, string) (string, error) {
	return "", nil
}
func (s *stubAuthFederationService) CompleteOAuthLogin(context.Context, string, string, string) (*models.User, string, string, error) {
	return nil, "", "", nil
}
func (s *stubAuthFederationService) ListAvailableProviders(_ context.Context, _ string) ([]services.ProviderAvailability, error) {
	if s.listErr != nil {
		return nil, s.listErr
	}
	if s.providers != nil {
		return s.providers, nil
	}
	return []services.ProviderAvailability{{Provider: "google", Name: "Google"}}, nil
}

type stubOIDCService struct {
	authorizeLoc string
	authorizeErr *domains.OAuthRedirectError
	tokenResp    *domains.OIDCTokenResponse
	tokenErr     *domains.OAuthTokenError
	userInfo     map[string]any
	userInfoErr  *domains.OAuthTokenError
	issuer       string
}

func (s *stubOIDCService) Authorize(_ context.Context, _ string, _ *dtos.AuthorizeQuery) (string, *domains.OAuthRedirectError) {
	if s.authorizeErr != nil {
		return "", s.authorizeErr
	}
	if s.authorizeLoc != "" {
		return s.authorizeLoc, nil
	}
	return "http://client.example/callback?code=abc&state=xyz", nil
}
func (s *stubOIDCService) AuthorizationCodeToken(_ context.Context, _, _, _ string, _ url.Values) (*domains.OIDCTokenResponse, *domains.OAuthTokenError) {
	if s.tokenErr != nil {
		return nil, s.tokenErr
	}
	if s.tokenResp != nil {
		return s.tokenResp, nil
	}
	return &domains.OIDCTokenResponse{AccessToken: "oidc-access", TokenType: "Bearer", ExpiresIn: 3600}, nil
}
func (s *stubOIDCService) UserInfo(_ context.Context, _ string) (map[string]any, *domains.OAuthTokenError) {
	if s.userInfoErr != nil {
		return nil, s.userInfoErr
	}
	if s.userInfo != nil {
		return s.userInfo, nil
	}
	return map[string]any{"sub": testUserID, "email": "user@example.com"}, nil
}
func (s *stubOIDCService) OpenIDIssuer() string {
	if s.issuer != "" {
		return s.issuer
	}
	return "http://localhost:8080"
}

type stubOIDCFederationService struct {
	redirectURL string
	redirectErr error
	user        *models.User
	tenantID    string
	returnTo    string
	completeErr error
}

func (s *stubOIDCFederationService) BuildAuthorizeRedirectURL(_ context.Context, _, _, _ string) (string, error) {
	if s.redirectErr != nil {
		return "", s.redirectErr
	}
	if s.redirectURL != "" {
		return s.redirectURL, nil
	}
	return "https://accounts.google.com/o/oauth2/auth", nil
}
func (s *stubOIDCFederationService) CompleteOAuthLogin(_ context.Context, _, _, _ string) (*models.User, string, string, error) {
	if s.completeErr != nil {
		return nil, "", "", s.completeErr
	}
	u := s.user
	if u == nil {
		u = &models.User{BaseModel: models.NewBaseModel(), Email: "fed@example.com"}
	}
	tenantID := s.tenantID
	if tenantID == "" {
		tenantID = testTenantID
	}
	returnTo := s.returnTo
	if returnTo == "" {
		returnTo = "/login/federation/complete"
	}
	return u, tenantID, returnTo, nil
}
func (s *stubOIDCFederationService) ListAvailableProviders(context.Context, string) ([]services.ProviderAvailability, error) {
	return nil, nil
}

type stubClientRepo struct {
	client *models.Client
	err    error
}

func (s *stubClientRepo) GetByTenantAndClientID(context.Context, string, string) (*models.Client, error) {
	return nil, nil
}
func (s *stubClientRepo) GetByClientID(_ context.Context, clientID string) (*models.Client, error) {
	if s.err != nil {
		return nil, s.err
	}
	if s.client != nil {
		return s.client, nil
	}
	return &models.Client{BaseModel: models.NewBaseModel(), ClientID: clientID, TenantID: testTenantID}, nil
}
func (s *stubClientRepo) GetByID(context.Context, string) (*models.Client, error) { return nil, nil }
func (s *stubClientRepo) List(context.Context, string, *dtos.PageableRequest) (*dtos.DataResponse[models.Client], error) {
	return nil, nil
}
func (s *stubClientRepo) Create(context.Context, *models.Client) error { return nil }
func (s *stubClientRepo) Update(context.Context, string, repositories.ClientPatch) (*models.Client, error) {
	return nil, nil
}
func (s *stubClientRepo) Delete(context.Context, string) error { return nil }
func (s *stubClientRepo) ClientIDTaken(context.Context, string, string, string) (bool, error) {
	return false, nil
}

type stubWebauthnService struct {
	credentials []models.WebauthnCredential
	listErr     error
	options     json.RawMessage
	token       string
	startErr    error
	finishErr   error
	loginUser   *models.User
}

func (s *stubWebauthnService) ListCredentials(_ context.Context, _ string, _ *dtos.PageableRequest) ([]models.WebauthnCredential, *dtos.Pageable, error) {
	if s.listErr != nil {
		return nil, nil, s.listErr
	}
	if s.credentials != nil {
		return s.credentials, &dtos.Pageable{Page: 1, PageSize: 20, Total: int64(len(s.credentials))}, nil
	}
	return []models.WebauthnCredential{{BaseModel: models.NewBaseModel(), DeviceName: "MacBook"}}, &dtos.Pageable{Page: 1, PageSize: 20, Total: 1}, nil
}
func (s *stubWebauthnService) RegisterStart(_ context.Context, _, _ string) (json.RawMessage, string, error) {
	if s.startErr != nil {
		return nil, "", s.startErr
	}
	opts := s.options
	if opts == nil {
		opts = json.RawMessage(`{"challenge":"abc"}`)
	}
	token := s.token
	if token == "" {
		token = "reg-token"
	}
	return opts, token, nil
}
func (s *stubWebauthnService) RegisterFinish(context.Context, string, string, []byte) error {
	return s.finishErr
}
func (s *stubWebauthnService) LoginStart(_ context.Context, _ string) (json.RawMessage, string, error) {
	if s.startErr != nil {
		return nil, "", s.startErr
	}
	return json.RawMessage(`{"challenge":"login"}`), "login-token", nil
}
func (s *stubWebauthnService) LoginFinish(_ context.Context, _, _ string, _ []byte) (*models.User, error) {
	if s.finishErr != nil {
		return nil, s.finishErr
	}
	if s.loginUser != nil {
		return s.loginUser, nil
	}
	return &models.User{BaseModel: models.NewBaseModel(), Email: "user@example.com"}, nil
}

type stubAdminListService struct {
	stubAdminAuditService
}

func (stubAdminListService) GetStats(context.Context) (*dtos.AdminStatsResponse, error) {
	return &dtos.AdminStatsResponse{TotalUsers: 10, ActiveSessions: 5}, nil
}
func (stubAdminListService) ListUsers(_ context.Context, _, _ string, _ *dtos.PageableRequest) ([]*dtos.AdminUserResponse, *dtos.Pageable, error) {
	return []*dtos.AdminUserResponse{{ID: testUserID, Email: "admin@example.com"}}, &dtos.Pageable{Page: 1, PageSize: 20, Total: 1}, nil
}
func (stubAdminListService) ListTenants(context.Context, *dtos.PageableRequest) ([]*dtos.AdminTenantResponse, *dtos.Pageable, error) {
	return []*dtos.AdminTenantResponse{{ID: testTenantID, Name: "Acme"}}, &dtos.Pageable{Page: 1, PageSize: 20, Total: 1}, nil
}
func (stubAdminListService) ListClients(_ context.Context, _ string, _ *dtos.PageableRequest) ([]*dtos.AdminClientResponse, *dtos.Pageable, error) {
	return []*dtos.AdminClientResponse{{ID: testClientPK, ClientID: "spa-app"}}, &dtos.Pageable{Page: 1, PageSize: 20, Total: 1}, nil
}
func (stubAdminListService) ListIdentityProviders(_ context.Context, _ string, _ *dtos.PageableRequest) ([]*dtos.AdminIdentityProviderResponse, *dtos.Pageable, error) {
	return []*dtos.AdminIdentityProviderResponse{{Provider: "google", Enabled: true}}, &dtos.Pageable{Page: 1, PageSize: 20, Total: 1}, nil
}
func (stubAdminListService) AddMemberByEmail(context.Context, string, string, string) error { return nil }
func (stubAdminListService) RemoveMember(context.Context, string, string) error               { return nil }
func (stubAdminListService) GetClientUsage(_ context.Context, _ string) (*dtos.AdminClientUsageResponse, error) {
	return &dtos.AdminClientUsageResponse{ActiveRefreshTokens: 5}, nil
}
func (stubAdminListService) ListLoginHistory(_ context.Context, _ dtos.AdminLoginHistoryListParams, _ *dtos.PageableRequest) ([]*dtos.AdminAuditLogResponse, *dtos.Pageable, error) {
	return []*dtos.AdminAuditLogResponse{{ID: "00000000-0000-4000-8000-000000000011", Action: constants.AuditActionAuthLogin}}, &dtos.Pageable{Page: 1, PageSize: 20, Total: 1}, nil
}

type stubAdminUserActionsService struct {
	stubAdminAuditService
	forceLogoutTarget string
	resetMFATarget    string
	resetPasskeyTarget string
}

func (s *stubAdminUserActionsService) ForceLogoutUser(_ context.Context, _, target string) error {
	s.forceLogoutTarget = target
	return nil
}
func (s *stubAdminUserActionsService) ResetMFA(_ context.Context, _, target string) error {
	s.resetMFATarget = target
	return nil
}
func (s *stubAdminUserActionsService) ResetPasskeys(_ context.Context, _, target string) error {
	s.resetPasskeyTarget = target
	return nil
}
func (s *stubAdminUserActionsService) GetClientUsage(_ context.Context, clientID string) (*dtos.AdminClientUsageResponse, error) {
	return &dtos.AdminClientUsageResponse{ClientID: clientID, ActiveRefreshTokens: 2}, nil
}

func authHandler(user services.UserService, session services.SessionService, mfa services.MFAService, fed services.FederationService) *AuthHandler {
	return ProvideAuthHandler(user, session, mfa, fed, stubAuditService{}, handlerTestConfig(), validator.New())
}

func errUnauthorized() error {
	return errors.UnauthorizedError("invalid credentials", nil)
}
