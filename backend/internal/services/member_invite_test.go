package services

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/integration/email"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type stubInviteRepo struct {
	byToken map[string]*models.TenantInvite
	byEmail map[string]*models.TenantInvite
	created []*models.TenantInvite
}

func newStubInviteRepo() *stubInviteRepo {
	return &stubInviteRepo{
		byToken: map[string]*models.TenantInvite{},
		byEmail: map[string]*models.TenantInvite{},
	}
}

func (r *stubInviteRepo) Create(_ context.Context, invite *models.TenantInvite) error {
	r.created = append(r.created, invite)
	r.byToken[invite.TokenHash] = invite
	r.byEmail[invite.EmailLower+"|"+invite.TenantID] = invite
	return nil
}

func (r *stubInviteRepo) Save(_ context.Context, invite *models.TenantInvite) error {
	r.byToken[invite.TokenHash] = invite
	r.byEmail[invite.EmailLower+"|"+invite.TenantID] = invite
	return nil
}

func (r *stubInviteRepo) GetByTokenHash(_ context.Context, tokenHash string) (*models.TenantInvite, error) {
	invite, ok := r.byToken[tokenHash]
	if !ok {
		return nil, errors.NotFoundError("Tenant invite", nil)
	}
	return invite, nil
}

func (r *stubInviteRepo) FindPendingByEmailAndTenant(_ context.Context, emailLower, tenantID string) (*models.TenantInvite, error) {
	invite := r.byEmail[emailLower+"|"+tenantID]
	if invite == nil || invite.Status != constants.TenantInviteStatusPending {
		return nil, nil
	}
	return invite, nil
}

type inviteAuthStub struct {
	failAuth bool
	issued   *dtos.LoginResponse
}

func (s *inviteAuthStub) AuthenticateUser(_ context.Context, _ *dtos.LoginRequest) (*models.User, error) {
	if s.failAuth {
		return nil, errors.UnauthorizedError("Invalid email or password", nil)
	}
	return &models.User{}, nil
}

func (s *inviteAuthStub) IssueTokensForUser(_ context.Context, _ *models.User, tenantID string) (*dtos.LoginResponse, error) {
	s.issued = &dtos.LoginResponse{AccessToken: "access", ActiveTenantID: tenantID}
	return s.issued, nil
}

func TestAdminService_AddMemberByEmail_InvitesUnknownEmail(t *testing.T) {
	tenantRepo := newAdminTenantTestRepo()
	tenantID := "tenant-invite"
	tenantRepo.tenants[tenantID] = &models.Tenant{BaseModel: models.BaseModel{ID: tenantID}, Name: "Acme"}
	users := newUserTestRepo()
	invites := newStubInviteRepo()
	sender := new(MockEmailSender)
	sender.On("SendEmail", mock.Anything, mock.MatchedBy(func(req email.EmailRequest) bool {
		return req.TemplateID == email.TemplateMemberInvite &&
			req.To[0] == "new@example.com" &&
			strings.Contains(req.HTMLBody, "invite=") &&
			strings.Contains(req.TextBody, "Acme")
	})).Return(&email.EmailResponse{Status: "sent"}, nil)

	svc := &adminService{
		cfg:         &config.Config{OIDCLoginPageURL: "http://localhost:5173/login"},
		tenants:     tenantRepo,
		users:       users,
		memberships: &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}},
		audit:       &auditCapture{},
		mail:        ProvideEmailService(sender),
		invites:     invites,
	}
	require.NoError(t, svc.AddMemberByEmail(context.Background(), tenantID, "new@example.com", ""))
	require.Len(t, invites.created, 1)
	require.Equal(t, "member", invites.created[0].Role)
	require.Equal(t, constants.TenantInviteStatusPending, invites.created[0].Status)
	sender.AssertExpectations(t)
}

func TestMemberInviteService_Accept_CreatesUser(t *testing.T) {
	users := newUserTestRepo()
	invites := newStubInviteRepo()
	raw := "invite-token"
	invite := &models.TenantInvite{
		BaseModel:  models.NewBaseModel(),
		Email:      "new@example.com",
		EmailLower: "new@example.com",
		TenantID:   "tenant-invite",
		Role:       "member",
		TokenHash:  auth.HashOpaqueToken(raw),
		Status:     constants.TenantInviteStatusPending,
		ExpiresAt:  time.Now().UTC().Add(time.Hour),
		Tenant:     &models.Tenant{Name: "Acme"},
	}
	require.NoError(t, invites.Create(context.Background(), invite))
	memberships := &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}}
	authStub := &inviteAuthStub{}
	svc := &memberInviteService{
		invites:     invites,
		users:       users,
		memberships: memberships,
		auth:        authStub,
		audit:       &auditCapture{},
	}

	preview, err := svc.Preview(context.Background(), raw)
	require.NoError(t, err)
	require.Equal(t, "new@example.com", preview.Email)
	require.Equal(t, "Acme", preview.OrganizationName)

	out, err := svc.Accept(context.Background(), &dtos.AcceptMemberInviteRequest{
		Token:     raw,
		Password:  "long-secure-passphrase",
		FirstName: "Ada",
		LastName:  "Lovelace",
	})
	require.NoError(t, err)
	require.Equal(t, "access", out.Login.AccessToken)
	require.Equal(t, "tenant-invite", out.Login.ActiveTenantID)
	require.Equal(t, constants.TenantInviteStatusAccepted, invite.Status)
	created := users.byEmail["new@example.com"]
	require.NotNil(t, created)
	require.Equal(t, "Ada", created.FirstName)
	require.True(t, created.EmailVerified)
	require.Len(t, memberships.byUser[created.ID], 1)
}

func TestMemberInviteService_Accept_ExistingUser(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("member@example.com", "long-secure-passphrase")
	invites := newStubInviteRepo()
	raw := "invite-token"
	invite := &models.TenantInvite{
		BaseModel:  models.NewBaseModel(),
		Email:      "member@example.com",
		EmailLower: "member@example.com",
		TenantID:   "tenant-invite",
		Role:       "admin",
		TokenHash:  auth.HashOpaqueToken(raw),
		Status:     constants.TenantInviteStatusPending,
		ExpiresAt:  time.Now().UTC().Add(time.Hour),
	}
	require.NoError(t, invites.Create(context.Background(), invite))
	memberships := &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}}
	svc := &memberInviteService{
		invites:     invites,
		users:       users,
		memberships: memberships,
		auth:        &inviteAuthStub{},
		audit:       &auditCapture{},
	}

	_, err := svc.Accept(context.Background(), &dtos.AcceptMemberInviteRequest{
		Token:    raw,
		Password: "long-secure-passphrase",
	})
	require.NoError(t, err)
	require.Len(t, memberships.byUser[u.ID], 1)
	require.Equal(t, constants.TenantMembershipRoleAdmin, memberships.byUser[u.ID][0].Role)
	require.True(t, u.EmailVerified)
}
