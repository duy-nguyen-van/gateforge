package services

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/integration/email"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type stubInviteRepo struct {
	byToken   map[string]*models.TenantInvite
	byEmail   map[string]*models.TenantInvite
	created   []*models.TenantInvite
	createErr error
	saveErr   error
	getErr    error
	findErr   error
}

func newStubInviteRepo() *stubInviteRepo {
	return &stubInviteRepo{
		byToken: map[string]*models.TenantInvite{},
		byEmail: map[string]*models.TenantInvite{},
	}
}

func (r *stubInviteRepo) Create(_ context.Context, invite *models.TenantInvite) error {
	if r.createErr != nil {
		return r.createErr
	}
	r.created = append(r.created, invite)
	r.byToken[invite.TokenHash] = invite
	r.byEmail[invite.EmailLower+"|"+invite.TenantID] = invite
	return nil
}

func (r *stubInviteRepo) Save(_ context.Context, invite *models.TenantInvite) error {
	if r.saveErr != nil {
		return r.saveErr
	}
	r.byToken[invite.TokenHash] = invite
	r.byEmail[invite.EmailLower+"|"+invite.TenantID] = invite
	return nil
}

func (r *stubInviteRepo) GetByTokenHash(_ context.Context, tokenHash string) (*models.TenantInvite, error) {
	if r.getErr != nil {
		return nil, r.getErr
	}
	invite, ok := r.byToken[tokenHash]
	if !ok {
		return nil, errors.NotFoundError("Tenant invite", nil)
	}
	return invite, nil
}

func (r *stubInviteRepo) FindPendingByEmailAndTenant(_ context.Context, emailLower, tenantID string) (*models.TenantInvite, error) {
	if r.findErr != nil {
		return nil, r.findErr
	}
	invite := r.byEmail[emailLower+"|"+tenantID]
	if invite == nil || invite.Status != constants.TenantInviteStatusPending {
		return nil, nil
	}
	return invite, nil
}

func (r *stubInviteRepo) GetByID(_ context.Context, id string) (*models.TenantInvite, error) {
	if r.getErr != nil {
		return nil, r.getErr
	}
	for _, invite := range r.byEmail {
		if invite.ID == id {
			return invite, nil
		}
	}
	return nil, errors.NotFoundError("Tenant invite", nil)
}

func (r *stubInviteRepo) ListPending(_ context.Context, pr *dtos.PageableRequest) (*dtos.DataResponse[models.TenantInvite], error) {
	if r.getErr != nil {
		return nil, r.getErr
	}
	rows := make([]models.TenantInvite, 0)
	for _, invite := range r.byEmail {
		if invite.Status == constants.TenantInviteStatusPending {
			rows = append(rows, *invite)
		}
	}
	page, pageable := dtos.PaginateSlice(rows, pr)
	return &dtos.DataResponse[models.TenantInvite]{Data: page, Pageable: pageable}, nil
}

func (r *stubInviteRepo) ListPendingByTenant(_ context.Context, tenantID string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.TenantInvite], error) {
	if r.getErr != nil {
		return nil, r.getErr
	}
	rows := make([]models.TenantInvite, 0)
	for _, invite := range r.byEmail {
		if invite.TenantID == tenantID && invite.Status == constants.TenantInviteStatusPending {
			rows = append(rows, *invite)
		}
	}
	page, pageable := dtos.PaginateSlice(rows, pr)
	return &dtos.DataResponse[models.TenantInvite]{Data: page, Pageable: pageable}, nil
}

type inviteAuthStub struct {
	failAuth bool
	issueErr error
	issued   *dtos.LoginResponse
}

func (s *inviteAuthStub) AuthenticateUser(_ context.Context, _ *dtos.LoginRequest) (*models.User, error) {
	if s.failAuth {
		return nil, errors.UnauthorizedError("Invalid email or password", nil)
	}
	return &models.User{}, nil
}

func (s *inviteAuthStub) IssueTokensForUser(_ context.Context, _ *models.User, tenantID string) (*dtos.LoginResponse, error) {
	if s.issueErr != nil {
		return nil, s.issueErr
	}
	s.issued = &dtos.LoginResponse{AccessToken: "access", ActiveTenantID: tenantID}
	return s.issued, nil
}

type inviteMFAStub struct {
	active    bool
	activeErr error
	ticket    string
	ticketErr error
}

func (s *inviteMFAStub) HasActiveMFA(context.Context, string) (bool, error) {
	return s.active, s.activeErr
}

func (s *inviteMFAStub) CreateLoginTicket(context.Context, auth.MFAPendingPayload) (string, int64, error) {
	if s.ticketErr != nil {
		return "", 0, s.ticketErr
	}
	return s.ticket, 600, nil
}

type failAudit struct{}

func (failAudit) Record(context.Context, domains.AuditRecordParams) {}
func (failAudit) RecordRequired(context.Context, domains.AuditRecordParams) error {
	return errors.InternalError("audit failed", nil)
}
func (failAudit) Shutdown(context.Context) error { return nil }

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

func TestMemberInviteService_PreviewAndAccept_Errors(t *testing.T) {
	raw := "invite-token"
	fresh := func() (*memberInviteService, *stubInviteRepo, *userTestRepo, *stubMembershipRepo) {
		users := newUserTestRepo()
		invites := newStubInviteRepo()
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
		svc := &memberInviteService{invites: invites, users: users, memberships: memberships, auth: &inviteAuthStub{}, audit: &auditCapture{}}
		return svc, invites, users, memberships
	}

	t.Run("empty token", func(t *testing.T) {
		svc, _, _, _ := fresh()
		_, err := svc.Preview(context.Background(), "  ")
		require.Equal(t, errors.ErrorTypeValidation, errors.GetAppError(err).Type)
	})

	t.Run("unknown token", func(t *testing.T) {
		svc, _, _, _ := fresh()
		_, err := svc.Preview(context.Background(), "missing")
		require.Equal(t, errors.ErrorTypeForbidden, errors.GetAppError(err).Type)
	})

	t.Run("expired invite", func(t *testing.T) {
		svc, invites, _, _ := fresh()
		invites.byToken[auth.HashOpaqueToken(raw)].ExpiresAt = time.Now().UTC().Add(-time.Minute)
		_, err := svc.Accept(context.Background(), &dtos.AcceptMemberInviteRequest{Token: raw, Password: "long-secure-passphrase", FirstName: "Ada"})
		require.Equal(t, errors.ErrorTypeForbidden, errors.GetAppError(err).Type)
	})

	t.Run("lookup error", func(t *testing.T) {
		svc, invites, _, _ := fresh()
		invites.getErr = errors.InternalError("db down", nil)
		_, err := svc.Preview(context.Background(), raw)
		require.Equal(t, errors.ErrorTypeInternal, errors.GetAppError(err).Type)
	})

	t.Run("nil request and short password", func(t *testing.T) {
		svc, _, _, _ := fresh()
		_, err := svc.Accept(context.Background(), nil)
		require.Equal(t, errors.ErrorTypeValidation, errors.GetAppError(err).Type)
		_, err = svc.Accept(context.Background(), &dtos.AcceptMemberInviteRequest{Token: raw, Password: "short"})
		require.Equal(t, errors.ErrorTypeValidation, errors.GetAppError(err).Type)
	})

	t.Run("missing first name", func(t *testing.T) {
		svc, _, _, _ := fresh()
		_, err := svc.Accept(context.Background(), &dtos.AcceptMemberInviteRequest{Token: raw, Password: "long-secure-passphrase"})
		require.Equal(t, errors.ErrorTypeValidation, errors.GetAppError(err).Type)
	})

	t.Run("user lookup error", func(t *testing.T) {
		svc, _, users, _ := fresh()
		users.getErr = errors.InternalError("db down", nil)
		_, err := svc.Accept(context.Background(), &dtos.AcceptMemberInviteRequest{Token: raw, Password: "long-secure-passphrase", FirstName: "Ada"})
		require.Equal(t, errors.ErrorTypeInternal, errors.GetAppError(err).Type)
	})

	t.Run("create user error", func(t *testing.T) {
		svc, _, users, _ := fresh()
		users.createErr = errors.InternalError("insert failed", nil)
		_, err := svc.Accept(context.Background(), &dtos.AcceptMemberInviteRequest{Token: raw, Password: "long-secure-passphrase", FirstName: "Ada"})
		require.Equal(t, errors.ErrorTypeInternal, errors.GetAppError(err).Type)
	})

	t.Run("wrong password", func(t *testing.T) {
		svc, _, users, memberships := fresh()
		users.seed("new@example.com", "long-secure-passphrase")
		svc.auth = &inviteAuthStub{failAuth: true}
		_, err := svc.Accept(context.Background(), &dtos.AcceptMemberInviteRequest{Token: raw, Password: "long-secure-passphrase", FirstName: "Ada"})
		require.Equal(t, errors.ErrorTypeUnauthorized, errors.GetAppError(err).Type)
		require.Empty(t, memberships.byUser)
	})

	t.Run("already verified skips update", func(t *testing.T) {
		svc, _, users, _ := fresh()
		u := users.seed("new@example.com", "long-secure-passphrase")
		u.EmailVerified = true
		users.markErr = errors.InternalError("should not mark", nil)
		out, err := svc.Accept(context.Background(), &dtos.AcceptMemberInviteRequest{Token: raw, Password: "long-secure-passphrase"})
		require.NoError(t, err)
		require.Equal(t, "access", out.Login.AccessToken)
		require.True(t, u.EmailVerified)
	})

	t.Run("mark verified error", func(t *testing.T) {
		svc, _, users, _ := fresh()
		users.seed("new@example.com", "long-secure-passphrase")
		users.markErr = errors.InternalError("mark failed", nil)
		_, err := svc.Accept(context.Background(), &dtos.AcceptMemberInviteRequest{Token: raw, Password: "long-secure-passphrase"})
		require.Equal(t, errors.ErrorTypeInternal, errors.GetAppError(err).Type)
	})

	t.Run("membership lookup and create errors", func(t *testing.T) {
		svc, _, _, memberships := fresh()
		memberships.existsErr = errors.InternalError("membership lookup", nil)
		_, err := svc.Accept(context.Background(), &dtos.AcceptMemberInviteRequest{Token: raw, Password: "long-secure-passphrase", FirstName: "Ada"})
		require.Equal(t, errors.ErrorTypeInternal, errors.GetAppError(err).Type)

		memberships.existsErr = nil
		memberships.createErr = errors.InternalError("membership create", nil)
		_, err = svc.Accept(context.Background(), &dtos.AcceptMemberInviteRequest{Token: raw, Password: "long-secure-passphrase", FirstName: "Ada"})
		require.Equal(t, errors.ErrorTypeInternal, errors.GetAppError(err).Type)
	})

	t.Run("already a member", func(t *testing.T) {
		svc, _, users, memberships := fresh()
		u := users.seed("new@example.com", "long-secure-passphrase")
		memberships.active[u.ID] = map[string]bool{"tenant-invite": true}
		_, err := svc.Accept(context.Background(), &dtos.AcceptMemberInviteRequest{Token: raw, Password: "long-secure-passphrase"})
		require.NoError(t, err)
		require.Empty(t, memberships.byUser[u.ID])
	})

	t.Run("save invite error", func(t *testing.T) {
		svc, invites, _, _ := fresh()
		invites.saveErr = errors.InternalError("save failed", nil)
		_, err := svc.Accept(context.Background(), &dtos.AcceptMemberInviteRequest{Token: raw, Password: "long-secure-passphrase", FirstName: "Ada"})
		require.Equal(t, errors.ErrorTypeInternal, errors.GetAppError(err).Type)
	})

	t.Run("mfa required", func(t *testing.T) {
		svc, _, users, _ := fresh()
		users.seed("new@example.com", "long-secure-passphrase")
		svc.mfa = &inviteMFAStub{active: true, ticket: "ticket-1"}
		out, err := svc.Accept(context.Background(), &dtos.AcceptMemberInviteRequest{Token: raw, Password: "long-secure-passphrase"})
		require.NoError(t, err)
		require.Nil(t, out.Login)
		require.True(t, out.MFA.MfaRequired)
		require.Equal(t, "ticket-1", out.MFA.MfaTicket)
	})

	t.Run("mfa lookup error", func(t *testing.T) {
		svc, _, users, _ := fresh()
		users.seed("new@example.com", "long-secure-passphrase")
		svc.mfa = &inviteMFAStub{activeErr: errors.InternalError("mfa lookup", nil)}
		_, err := svc.Accept(context.Background(), &dtos.AcceptMemberInviteRequest{Token: raw, Password: "long-secure-passphrase"})
		require.Equal(t, errors.ErrorTypeInternal, errors.GetAppError(err).Type)
	})

	t.Run("mfa ticket error", func(t *testing.T) {
		svc, _, users, _ := fresh()
		users.seed("new@example.com", "long-secure-passphrase")
		svc.mfa = &inviteMFAStub{active: true, ticketErr: errors.InternalError("ticket", nil)}
		_, err := svc.Accept(context.Background(), &dtos.AcceptMemberInviteRequest{Token: raw, Password: "long-secure-passphrase"})
		require.Equal(t, errors.ErrorTypeInternal, errors.GetAppError(err).Type)
	})

	t.Run("issue tokens error", func(t *testing.T) {
		svc, _, _, _ := fresh()
		svc.auth = &inviteAuthStub{issueErr: errors.InternalError("tokens", nil)}
		_, err := svc.Accept(context.Background(), &dtos.AcceptMemberInviteRequest{Token: raw, Password: "long-secure-passphrase", FirstName: "Ada"})
		require.Equal(t, errors.ErrorTypeInternal, errors.GetAppError(err).Type)
	})
}

func TestAdminService_InviteNewMember_Branches(t *testing.T) {
	tenantID := "tenant-invite"
	newSvc := func(invites *stubInviteRepo, sender *MockEmailSender) *adminService {
		tenantRepo := newAdminTenantTestRepo()
		tenantRepo.tenants[tenantID] = &models.Tenant{BaseModel: models.BaseModel{ID: tenantID}, Name: "Acme"}
		return &adminService{
			cfg:         &config.Config{OIDCLoginPageURL: "http://localhost:5173/login?return_to=/console"},
			tenants:     tenantRepo,
			users:       newUserTestRepo(),
			memberships: &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}},
			audit:       &auditCapture{},
			mail:        ProvideEmailService(sender),
			invites:     invites,
		}
	}
	matchInvite := mock.MatchedBy(func(req email.EmailRequest) bool {
		return req.TemplateID == email.TemplateMemberInvite && strings.Contains(req.HTMLBody, "invite=")
	})

	t.Run("refreshes pending invite", func(t *testing.T) {
		invites := newStubInviteRepo()
		sender := new(MockEmailSender)
		sender.On("SendEmail", mock.Anything, matchInvite).Return(&email.EmailResponse{Status: "sent"}, nil).Twice()
		svc := newSvc(invites, sender)
		require.NoError(t, svc.AddMemberByEmail(context.Background(), tenantID, " New@Example.com ", "admin"))
		firstHash := invites.created[0].TokenHash
		require.NoError(t, svc.AddMemberByEmail(context.Background(), tenantID, "new@example.com", "member"))
		require.Len(t, invites.created, 1)
		require.NotEqual(t, firstHash, invites.created[0].TokenHash)
		require.Equal(t, "member", invites.created[0].Role)
		require.Equal(t, "new@example.com", invites.created[0].EmailLower)
		sender.AssertExpectations(t)
	})

	t.Run("send failure keeps invite", func(t *testing.T) {
		invites := newStubInviteRepo()
		sender := new(MockEmailSender)
		sender.On("SendEmail", mock.Anything, matchInvite).Return(nil, errors.InternalError("smtp", nil))
		svc := newSvc(invites, sender)
		require.NoError(t, svc.AddMemberByEmail(context.Background(), tenantID, "new@example.com", ""))
		require.Equal(t, constants.TenantInviteStatusPending, invites.created[0].Status)
	})

	t.Run("repo and audit errors", func(t *testing.T) {
		invites := newStubInviteRepo()
		invites.findErr = errors.InternalError("find failed", nil)
		svc := newSvc(invites, new(MockEmailSender))
		err := svc.AddMemberByEmail(context.Background(), tenantID, "new@example.com", "")
		require.Equal(t, errors.ErrorTypeInternal, errors.GetAppError(err).Type)

		invites.findErr = nil
		invites.createErr = errors.InternalError("create failed", nil)
		err = svc.AddMemberByEmail(context.Background(), tenantID, "new@example.com", "")
		require.Equal(t, errors.ErrorTypeInternal, errors.GetAppError(err).Type)

		invites.createErr = nil
		svc.audit = failAudit{}
		err = svc.AddMemberByEmail(context.Background(), tenantID, "new@example.com", "")
		require.Equal(t, errors.ErrorTypeInternal, errors.GetAppError(err).Type)
	})

	t.Run("invite url", func(t *testing.T) {
		svc := &adminService{cfg: &config.Config{}}
		got := svc.memberInviteURL("raw token")
		require.Contains(t, got, "http://localhost:5173/login")
		require.Contains(t, got, "invite=")

		svc.cfg.OIDCLoginPageURL = "http://bad host"
		require.Equal(t, "http://bad host", svc.memberInviteURL("tok"))
	})
}

func TestAdminService_ListAndResendTenantInvite(t *testing.T) {
	tenantID := "tenant-invite"
	tenantRepo := newAdminTenantTestRepo()
	tenantRepo.tenants[tenantID] = &models.Tenant{BaseModel: models.BaseModel{ID: tenantID}, Name: "Acme"}
	otherID := "tenant-other"
	tenantRepo.tenants[otherID] = &models.Tenant{BaseModel: models.BaseModel{ID: otherID}, Name: "Other"}
	invites := newStubInviteRepo()
	sender := new(MockEmailSender)
	sender.On("SendEmail", mock.Anything, mock.MatchedBy(func(req email.EmailRequest) bool {
		return req.TemplateID == email.TemplateMemberInvite
	})).Return(&email.EmailResponse{Status: "sent"}, nil).Twice()
	svc := &adminService{
		cfg:         &config.Config{OIDCLoginPageURL: "http://localhost:5173/login"},
		tenants:     tenantRepo,
		users:       newUserTestRepo(),
		memberships: &stubMembershipRepo{byUser: map[string][]models.TenantMembership{}, active: map[string]map[string]bool{}},
		audit:       &auditCapture{},
		mail:        ProvideEmailService(sender),
		invites:     invites,
	}

	require.NoError(t, svc.AddMemberByEmail(context.Background(), tenantID, "new@example.com", "member"))
	rows, page, err := svc.ListTenantInvites(context.Background(), tenantID, &dtos.PageableRequest{Page: 1, PageSize: 20})
	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.Equal(t, "new@example.com", rows[0].Email)
	require.Equal(t, "pending", rows[0].Status)
	require.Equal(t, tenantID, rows[0].TenantID)
	require.Equal(t, int64(1), page.Total)

	invites.created[0].Tenant = &models.Tenant{Name: "Acme"}
	all, _, err := svc.ListInvites(context.Background(), &dtos.PageableRequest{Page: 1, PageSize: 20})
	require.NoError(t, err)
	require.Len(t, all, 1)
	require.Equal(t, "Acme", all[0].TenantName)

	invites.created[0].ExpiresAt = time.Now().UTC().Add(-time.Hour)
	rows, _, err = svc.ListTenantInvites(context.Background(), tenantID, &dtos.PageableRequest{Page: 1, PageSize: 20})
	require.NoError(t, err)
	require.Equal(t, "expired", rows[0].Status)

	require.NoError(t, svc.ResendTenantInvite(context.Background(), tenantID, rows[0].ID))
	require.True(t, invites.created[0].ExpiresAt.After(time.Now().UTC()))
	sender.AssertExpectations(t)

	err = svc.ResendTenantInvite(context.Background(), otherID, rows[0].ID)
	require.Equal(t, errors.ErrorTypeNotFound, errors.GetAppError(err).Type)

	invites.created[0].Status = constants.TenantInviteStatusAccepted
	err = svc.ResendTenantInvite(context.Background(), tenantID, rows[0].ID)
	require.Equal(t, errors.ErrorTypeNotFound, errors.GetAppError(err).Type)

	svc.invites = nil
	empty, emptyPage, err := svc.ListTenantInvites(context.Background(), tenantID, nil)
	require.NoError(t, err)
	require.Empty(t, empty)
	require.Equal(t, int64(0), emptyPage.Total)
	err = svc.ResendTenantInvite(context.Background(), tenantID, rows[0].ID)
	require.Equal(t, errors.ErrorTypeNotFound, errors.GetAppError(err).Type)
}
