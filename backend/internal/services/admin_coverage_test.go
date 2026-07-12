package services

import (
	"context"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/stretchr/testify/require"
)

func TestAdminService_GetStats(t *testing.T) {
	users := newUserTestRepo()
	users.seed("a@example.com", "secret")
	users.seed("b@example.com", "secret")
	mfa := newMFATOTPTestRepo()
	mfa.byUser["u1"] = &models.UserMFATOTP{Enabled: true}
	sessions := newSessionTestRepo()
	require.NoError(t, sessions.Create(context.Background(), &models.Session{UserID: "u1", TenantID: "t1"}))

	svc := &adminService{
		users:    users,
		mfaTOTP:  mfa,
		sessions: sessions,
	}

	stats, err := svc.GetStats(context.Background())
	require.NoError(t, err)
	require.Equal(t, int64(2), stats.TotalUsers)
	require.Equal(t, int64(1), stats.MFAEnabledCount)
	require.Equal(t, 50.0, stats.MFAEnabledPercent)
	require.Equal(t, int64(1), stats.ActiveSessions)
}

func TestAdminService_ListUsers(t *testing.T) {
	users := newUserTestRepo()
	users.seed("list@example.com", "secret")
	svc := &adminService{users: users}

	rows, pageable, err := svc.ListUsers(context.Background(), "", "", dtos.NewPageableRequest())
	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.Equal(t, int64(1), pageable.Total)
}

func TestAdminService_ListTenantsAndClients(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	tenantID := "tenant-list"
	tenantRepo.tenants[tenantID] = &models.Tenant{BaseModel: models.BaseModel{ID: tenantID}, Name: "Acme"}
	clientRepo := newAdminClientTestRepo()
	clientRepo.clients["c1"] = &models.Client{BaseModel: models.BaseModel{ID: "c1"}, TenantID: tenantID, ClientID: "app", Name: "App"}

	svc := &adminService{tenants: tenantRepo, clients: clientRepo, cfg: testConfig()}

	tenants, _, err := svc.ListTenants(context.Background(), dtos.NewPageableRequest())
	require.NoError(t, err)
	require.Len(t, tenants, 1)

	clients, _, err := svc.ListClients(context.Background(), tenantID, dtos.NewPageableRequest())
	require.NoError(t, err)
	require.Len(t, clients, 1)
}

func TestAdminService_ListIdentityProviders(t *testing.T) {
	svc := newAdminIdpTestService(&adminIdpTipStub{})
	rows, pageable, err := svc.ListIdentityProviders(context.Background(), "tenant-1", dtos.NewPageableRequest())
	require.NoError(t, err)
	require.NotEmpty(t, rows)
	require.NotNil(t, pageable)
	require.NotEmpty(t, rows[0].RedirectURI)
}

func TestAdminService_GetTenantByID(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	tenantID := "tenant-get"
	tenantRepo.tenants[tenantID] = &models.Tenant{BaseModel: models.BaseModel{ID: tenantID}, Name: "Acme"}
	svc := &adminService{tenants: tenantRepo}

	tenant, err := svc.GetTenantByID(context.Background(), tenantID)
	require.NoError(t, err)
	require.Equal(t, "Acme", tenant.Name)
}

func TestAdminService_GetUserByID(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("detail@example.com", "secret")
	memberships := &stubMembershipRepo{
		byUser: map[string][]models.TenantMembership{
			u.ID: {{UserID: u.ID, TenantID: "tenant-1", Role: constants.TenantMembershipRoleMember, Status: constants.TenantMembershipStatusActive}},
		},
	}
	tenantRepo := newAdminTenantTestRepo()
	tenantRepo.tenants["tenant-1"] = &models.Tenant{BaseModel: models.BaseModel{ID: "tenant-1"}, Name: "Acme"}
	mfa := newMFATOTPTestRepo()
	mfa.byUser[u.ID] = &models.UserMFATOTP{Enabled: true}
	sessions := newSessionTestRepo()
	require.NoError(t, sessions.Create(context.Background(), &models.Session{UserID: u.ID, TenantID: "tenant-1"}))
	webauthn := newWebauthnCredTestRepo()
	webauthn.byUser[u.ID] = []models.WebauthnCredential{{BaseModel: models.NewBaseModel()}}

	svc := &adminService{
		users: users, memberships: memberships, tenants: tenantRepo,
		mfaTOTP: mfa, sessions: sessions, webauthnCreds: webauthn,
	}

	detail, err := svc.GetUserByID(context.Background(), u.ID)
	require.NoError(t, err)
	require.Equal(t, u.Email, detail.Email)
	require.True(t, detail.MFAEnabled)
	require.Equal(t, 1, detail.PasskeyCount)
	require.Equal(t, int64(1), detail.ActiveSessions)
	require.Len(t, detail.Memberships, 1)
	require.Equal(t, "Acme", detail.Memberships[0].TenantName)
}

func TestAdminService_RemoveMember(t *testing.T) {
	memberships := &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}}
	audit := &auditCapture{}
	svc := &adminService{memberships: memberships, audit: audit}

	require.NoError(t, svc.RemoveMember(context.Background(), "tenant-1", "user-1"))
	require.Equal(t, constants.AuditActionAdminMemberRemove, audit.params[0].Action)
}

func TestAdminService_GetClientUsage(t *testing.T) {
	clientRepo := newAdminClientTestRepo()
	clientID := "client-usage"
	clientRepo.clients[clientID] = &models.Client{
		BaseModel: models.BaseModel{ID: clientID},
		ClientID:  "my-app",
	}
	refreshRepo := newRefreshTokenTestRepo()
	auditRepo := &stubAuditLogRepo{}

	svc := &adminService{clients: clientRepo, refreshTokens: refreshRepo, auditLogs: auditRepo}

	usage, err := svc.GetClientUsage(context.Background(), clientID)
	require.NoError(t, err)
	require.Equal(t, "my-app", usage.ClientID)
	require.Equal(t, int64(1), usage.TotalRefreshTokens)
}

func TestAdminService_AddMemberByEmail_Success(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	tenantID := "tenant-add"
	tenantRepo.tenants[tenantID] = &models.Tenant{BaseModel: models.BaseModel{ID: tenantID}, Name: "Acme"}
	users := newUserTestRepo()
	u := users.seed("member@example.com", "secret")
	memberships := &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}}
	audit := &auditCapture{}
	svc := &adminService{tenants: tenantRepo, users: users, memberships: memberships, audit: audit}

	require.NoError(t, svc.AddMemberByEmail(context.Background(), tenantID, "member@example.com", ""))
	require.Len(t, memberships.byUser[u.ID], 1)
	require.Equal(t, constants.AuditActionAdminMemberAdd, audit.params[0].Action)
}
