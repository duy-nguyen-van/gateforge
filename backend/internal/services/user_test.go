package services

import (
	"context"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/stretchr/testify/require"
)

func newUserTestService(t *testing.T, users *userTestRepo, memberships *stubMembershipRepo) UserService {
	t.Helper()
	cfg := testConfig()
	tokenSvc, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	require.NoError(t, err)
	tenantCtx := ProvideTenantContextService(cfg, &stubClientRepo{}, &stubTenantRepo{}, memberships)
	audit := &auditCapture{}
	return ProvideUserService(users, memberships, newRefreshTokenTestRepo(), tenantCtx, cfg, tokenSvc, audit)
}

func TestUserService_Register(t *testing.T) {
	users := newUserTestRepo()
	memberships := &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}}
	svc := newUserTestService(t, users, memberships)

	u, err := svc.Register(context.Background(), &dtos.RegisterRequest{
		Email:     "new@example.com",
		Password:  "password123",
		FirstName: "New",
		LastName:  "User",
		TenantID:  testConfig().DefaultTenantID,
	}, "")
	require.NoError(t, err)
	require.Equal(t, "new@example.com", u.Email)
	require.Len(t, users.created, 1)
	require.Len(t, memberships.byUser[u.ID], 1)
	require.Equal(t, testConfig().DefaultTenantID, memberships.byUser[u.ID][0].TenantID)
}

func TestUserService_Register_DuplicateEmail(t *testing.T) {
	users := newUserTestRepo()
	users.seed("taken@example.com", "password123")
	memberships := &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}}
	svc := newUserTestService(t, users, memberships)

	_, err := svc.Register(context.Background(), &dtos.RegisterRequest{
		Email:    "taken@example.com",
		Password: "password123",
	}, "")
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.Equal(t, errors.ErrorTypeConflict, appErr.Type)
}

func TestUserService_AuthenticateUser(t *testing.T) {
	users := newUserTestRepo()
	users.seed("alice@example.com", "secret123")
	memberships := &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}}
	svc := newUserTestService(t, users, memberships)

	u, err := svc.AuthenticateUser(context.Background(), &dtos.LoginRequest{
		Email:    "alice@example.com",
		Password: "secret123",
	})
	require.NoError(t, err)
	require.Equal(t, "alice@example.com", u.Email)

	_, err = svc.AuthenticateUser(context.Background(), &dtos.LoginRequest{
		Email:    "alice@example.com",
		Password: "wrong",
	})
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.Equal(t, errors.ErrorTypeUnauthorized, appErr.Type)

	_, err = svc.AuthenticateUser(context.Background(), &dtos.LoginRequest{
		Email:    "missing@example.com",
		Password: "secret123",
	})
	require.Error(t, err)
}

func TestUserService_AuthenticateUser_InactiveAndNoPassword(t *testing.T) {
	users := newUserTestRepo()
	disabled := users.seed("disabled@example.com", "secret123")
	disabled.Status = constants.UserStatusDisabled

	noPass := &models.User{
		BaseModel: models.NewBaseModel(),
		Email:     "nopass@example.com",
		Status:    constants.UserStatusActive,
	}
	noPass.EmailLower = "nopass@example.com"
	users.byEmail["nopass@example.com"] = noPass
	users.users[noPass.ID] = noPass

	memberships := &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}}
	svc := newUserTestService(t, users, memberships)

	_, err := svc.AuthenticateUser(context.Background(), &dtos.LoginRequest{Email: "disabled@example.com", Password: "secret123"})
	require.Error(t, err)

	_, err = svc.AuthenticateUser(context.Background(), &dtos.LoginRequest{Email: "nopass@example.com", Password: "secret123"})
	require.Error(t, err)
}

func TestUserService_LoginAndRefresh(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("login@example.com", "secret123")
	tenantID := "tenant-login"
	memberships := &stubMembershipRepo{
		byUser: map[string][]models.TenantMembership{
			u.ID: {{UserID: u.ID, TenantID: tenantID, Status: constants.TenantMembershipStatusActive}},
		},
		active: map[string]map[string]bool{u.ID: {tenantID: true}},
	}
	cfg := testConfig()
	tokenSvc, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	require.NoError(t, err)
	refreshRepo := newRefreshTokenTestRepo()
	tenantCtx := ProvideTenantContextService(cfg, &stubClientRepo{}, &stubTenantRepo{}, memberships)
	svc := ProvideUserService(users, memberships, refreshRepo, tenantCtx, cfg, tokenSvc, &auditCapture{})

	loginResp, sel, err := svc.Login(context.Background(), &dtos.LoginRequest{
		Email:    "login@example.com",
		Password: "secret123",
		TenantID: tenantID,
	}, "")
	require.NoError(t, err)
	require.Nil(t, sel)
	require.NotEmpty(t, loginResp.AccessToken)
	require.NotEmpty(t, loginResp.RefreshToken)
	require.Equal(t, tenantID, loginResp.ActiveTenantID)

	refreshed, err := svc.Refresh(context.Background(), &dtos.RefreshTokenRequest{RefreshToken: loginResp.RefreshToken})
	require.NoError(t, err)
	require.NotEmpty(t, refreshed.AccessToken)
	require.NotEqual(t, loginResp.RefreshToken, refreshed.RefreshToken)
}

func TestUserService_Refresh_InvalidToken(t *testing.T) {
	users := newUserTestRepo()
	memberships := &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}}
	cfg := testConfig()
	tokenSvc, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	require.NoError(t, err)
	svc := ProvideUserService(users, memberships, newRefreshTokenTestRepo(), ProvideTenantContextService(cfg, &stubClientRepo{}, &stubTenantRepo{}, memberships), cfg, tokenSvc, &auditCapture{})

	_, err = svc.Refresh(context.Background(), &dtos.RefreshTokenRequest{RefreshToken: ""})
	require.Error(t, err)

	_, err = svc.Refresh(context.Background(), &dtos.RefreshTokenRequest{RefreshToken: "not-a-real-token"})
	require.Error(t, err)
}

func TestUserService_CompleteAuth_TenantPicker(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("multi@example.com", "secret123")
	memberships := &stubMembershipRepo{
		byUser: map[string][]models.TenantMembership{
			u.ID: {
				{UserID: u.ID, TenantID: "tenant-a", Status: constants.TenantMembershipStatusActive},
				{UserID: u.ID, TenantID: "tenant-b", Status: constants.TenantMembershipStatusActive},
			},
		},
		active: map[string]map[string]bool{u.ID: {"tenant-a": true, "tenant-b": true}},
	}
	cfg := testConfig()
	cfg.DefaultTenantID = ""
	tokenSvc, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	require.NoError(t, err)
	svc := ProvideUserService(users, memberships, newRefreshTokenTestRepo(), ProvideTenantContextService(cfg, &stubClientRepo{}, &stubTenantRepo{}, memberships), cfg, tokenSvc, &auditCapture{})

	loginResp, sel, err := svc.CompleteAuth(context.Background(), u, TenantResolveInput{UserID: u.ID})
	require.NoError(t, err)
	require.Nil(t, loginResp)
	require.NotNil(t, sel)
	require.True(t, sel.SelectionRequired)
	require.Len(t, sel.Tenants, 2)
	require.NotEmpty(t, sel.SelectionToken)
}

func TestUserService_SelectTenantAndSwitchTenant(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("picker@example.com", "secret123")
	memberships := &stubMembershipRepo{
		byUser: map[string][]models.TenantMembership{
			u.ID: {
				{UserID: u.ID, TenantID: "tenant-a", Status: constants.TenantMembershipStatusActive},
				{UserID: u.ID, TenantID: "tenant-b", Status: constants.TenantMembershipStatusActive},
			},
		},
		active: map[string]map[string]bool{u.ID: {"tenant-a": true, "tenant-b": true}},
	}
	cfg := testConfig()
	tokenSvc, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	require.NoError(t, err)
	audit := &auditCapture{}
	svc := ProvideUserService(users, memberships, newRefreshTokenTestRepo(), ProvideTenantContextService(cfg, &stubClientRepo{}, &stubTenantRepo{}, memberships), cfg, tokenSvc, audit)

	token, _, err := tokenSvc.SignSelectionToken(u.ID)
	require.NoError(t, err)

	_, err = svc.SelectTenant(context.Background(), &dtos.TenantSelectRequest{
		SelectionToken: token,
		TenantID:       "tenant-a",
	})
	require.NoError(t, err)
	require.Equal(t, constants.AuditActionTenantSelect, audit.params[len(audit.params)-1].Action)

	switched, err := svc.SwitchTenant(context.Background(), u.ID, "tenant-b")
	require.NoError(t, err)
	require.Equal(t, "tenant-b", switched.ActiveTenantID)
	require.Equal(t, constants.AuditActionTenantSwitch, audit.params[len(audit.params)-1].Action)
}

func TestUserService_ListMemberships(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("member@example.com", "secret123")
	memberships := &stubMembershipRepo{
		byUser: map[string][]models.TenantMembership{
			u.ID: {{
				UserID:   u.ID,
				TenantID: "tenant-a",
				Role:     constants.TenantMembershipRoleMember,
				Tenant:   &models.Tenant{Name: "Acme", Domain: "acme"},
			}},
		},
		active: map[string]map[string]bool{u.ID: {"tenant-a": true}},
	}
	svc := newUserTestService(t, users, memberships)

	summaries, err := svc.ListMemberships(context.Background(), u.ID)
	require.NoError(t, err)
	require.Len(t, summaries, 1)
	require.Equal(t, "Acme", summaries[0].Name)

	page, pageable, err := svc.ListMembershipsPaginated(context.Background(), u.ID, dtos.NewPageableRequest())
	require.NoError(t, err)
	require.Len(t, page, 1)
	require.Equal(t, int64(1), pageable.Total)
}

func TestUserService_RevokeAllRefreshTokensForUser(t *testing.T) {
	users := newUserTestRepo()
	memberships := &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}}
	refreshRepo := newRefreshTokenTestRepo()
	cfg := testConfig()
	tokenSvc, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	require.NoError(t, err)
	svc := ProvideUserService(users, memberships, refreshRepo, ProvideTenantContextService(cfg, &stubClientRepo{}, &stubTenantRepo{}, memberships), cfg, tokenSvc, &auditCapture{})

	require.NoError(t, svc.RevokeAllRefreshTokensForUser(context.Background(), ""))
	require.NoError(t, svc.RevokeAllRefreshTokensForUser(context.Background(), "user-1"))
	require.Equal(t, []string{"user-1"}, refreshRepo.revoked)
}

func TestUserService_GetOneByID(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("get@example.com", "secret123")
	memberships := &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}}
	svc := newUserTestService(t, users, memberships)

	got, err := svc.GetOneByID(context.Background(), u.ID)
	require.NoError(t, err)
	require.Equal(t, u.ID, got.ID)
}

func TestTenantSummariesFromMemberships(t *testing.T) {
	rows := []models.TenantMembership{{
		TenantID: "t1",
		Role:     constants.TenantMembershipRoleAdmin,
		Tenant:   &models.Tenant{Name: "Org", Domain: "org"},
	}}
	out := tenantSummariesFromMemberships(rows)
	require.Equal(t, "Org", out[0].Name)
	require.Equal(t, "admin", out[0].Role)
}

func TestUserService_Refresh_InactiveUser(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("inactive@example.com", "secret123")
	u.Status = constants.UserStatusDisabled
	memberships := &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}}
	cfg := testConfig()
	tokenSvc, err := auth.NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
	require.NoError(t, err)
	refreshRepo := newRefreshTokenTestRepo()
	raw, hash, err := auth.NewOpaqueRefreshToken()
	require.NoError(t, err)
	require.NoError(t, refreshRepo.Create(context.Background(), &models.RefreshToken{
		UserID:    u.ID,
		TenantID:  "tenant-1",
		TokenHash: hash,
		ExpiresAt: time.Now().UTC().Add(time.Hour),
	}))
	svc := ProvideUserService(users, memberships, refreshRepo, ProvideTenantContextService(cfg, &stubClientRepo{}, &stubTenantRepo{}, memberships), cfg, tokenSvc, &auditCapture{})

	_, err = svc.Refresh(context.Background(), &dtos.RefreshTokenRequest{RefreshToken: raw})
	require.Error(t, err)
}
