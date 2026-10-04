package services

import (
	"context"
	stderrors "errors"
	"net/url"
	"strings"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/logger"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/monitoring"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

const memberInviteTTL = 7 * 24 * time.Hour

// MemberInviteAcceptResult is either a session or an MFA challenge. UserID is for the session cookie and is not returned to the client.
type MemberInviteAcceptResult struct {
	Login  *dtos.LoginResponse
	MFA    *dtos.MFALoginChallengeResponse
	UserID string
}

// MemberInviteService turns an emailed link into an account and an organization membership.
type MemberInviteService interface {
	Preview(ctx context.Context, token string) (*dtos.MemberInvitePreview, error)
	Accept(ctx context.Context, req *dtos.AcceptMemberInviteRequest) (*MemberInviteAcceptResult, error)
}

type memberInviteAuth interface {
	AuthenticateUser(ctx context.Context, req *dtos.LoginRequest) (*models.User, error)
	IssueTokensForUser(ctx context.Context, u *models.User, tenantID string) (*dtos.LoginResponse, error)
}

type memberInviteMFA interface {
	HasActiveMFA(ctx context.Context, userID string) (bool, error)
	CreateLoginTicket(ctx context.Context, payload auth.MFAPendingPayload) (string, int64, error)
}

type memberInviteService struct {
	invites     repositories.TenantInviteRepository
	users       repositories.UserRepository
	memberships repositories.TenantMembershipRepository
	auth        memberInviteAuth
	mfa         memberInviteMFA
	audit       AuditService
}

// ProvideMemberInviteService wires invite preview and acceptance.
func ProvideMemberInviteService(
	invites repositories.TenantInviteRepository,
	users repositories.UserRepository,
	memberships repositories.TenantMembershipRepository,
	userSvc UserService,
	mfa MFAService,
	audit AuditService,
) MemberInviteService {
	return &memberInviteService{
		invites:     invites,
		users:       users,
		memberships: memberships,
		auth:        userSvc,
		mfa:         mfa,
		audit:       audit,
	}
}

func (s *memberInviteService) Preview(ctx context.Context, token string) (*dtos.MemberInvitePreview, error) {
	return monitoring.Observe(ctx, adminTracer, "MemberInviteService.Preview", nil,
		func(ctx context.Context) (*dtos.MemberInvitePreview, error) {
			invite, err := s.pendingInvite(ctx, token)
			if err != nil {
				return nil, err
			}
			preview := &dtos.MemberInvitePreview{
				Email: invite.Email,
				Role:  invite.Role,
			}
			if invite.Tenant != nil {
				preview.OrganizationName = invite.Tenant.Name
			}
			return preview, nil
		})
}

func (s *memberInviteService) Accept(ctx context.Context, req *dtos.AcceptMemberInviteRequest) (*MemberInviteAcceptResult, error) {
	return monitoring.Observe(ctx, adminTracer, "MemberInviteService.Accept", nil,
		func(ctx context.Context) (*MemberInviteAcceptResult, error) {
			return s.accept(ctx, req)
		})
}

func (s *memberInviteService) accept(ctx context.Context, req *dtos.AcceptMemberInviteRequest) (*MemberInviteAcceptResult, error) {
	if req == nil {
		return nil, errors.ValidationError("Invite token is required", nil)
	}
	invite, err := s.pendingInvite(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	if err := ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	user, err := s.users.GetByEmailLower(ctx, invite.EmailLower)
	if err != nil {
		var appErr *errors.AppError
		if !stderrors.As(err, &appErr) || appErr.Type != errors.ErrorTypeNotFound {
			return nil, err
		}
		user, err = s.createInvitedUser(ctx, invite, req)
		if err != nil {
			return nil, err
		}
	} else if _, err := s.auth.AuthenticateUser(ctx, &dtos.LoginRequest{
		Email:    invite.Email,
		Password: req.Password,
	}); err != nil {
		return nil, err
	} else if err := s.markInviteEmailVerified(ctx, user); err != nil {
		return nil, err
	}

	if err := s.ensureActiveMembership(ctx, user.ID, invite); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	invite.Status = constants.TenantInviteStatusAccepted
	invite.AcceptedUserID = &user.ID
	invite.ExpiresAt = now
	if err := s.invites.Save(ctx, invite); err != nil {
		return nil, err
	}
	if s.audit != nil {
		s.audit.Record(ctx, domains.AuditRecordParams{
			Action:       constants.AuditActionAuthMemberInviteAccept,
			Result:       constants.AuditResultSuccess,
			ActorType:    constants.AuditActorTypeUser,
			ActorID:      user.ID,
			TenantID:     invite.TenantID,
			ResourceType: constants.AuditResourceTypeMembership,
			ResourceID:   invite.ID,
		})
	}
	if s.mfa != nil {
		active, err := s.mfa.HasActiveMFA(ctx, user.ID)
		if err != nil {
			return nil, err
		}
		if active {
			ticket, exp, err := s.mfa.CreateLoginTicket(ctx, auth.MFAPendingPayload{
				UserID:   user.ID,
				TenantID: invite.TenantID,
			})
			if err != nil {
				return nil, err
			}
			return &MemberInviteAcceptResult{
				UserID: user.ID,
				MFA: &dtos.MFALoginChallengeResponse{
					MfaRequired: true,
					MfaTicket:   ticket,
					ExpiresIn:   exp,
				},
			}, nil
		}
	}
	login, err := s.auth.IssueTokensForUser(ctx, user, invite.TenantID)
	if err != nil {
		return nil, err
	}
	return &MemberInviteAcceptResult{Login: login, UserID: user.ID}, nil
}

func (s *memberInviteService) pendingInvite(ctx context.Context, token string) (*models.TenantInvite, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, errors.ValidationError("Invite token is required", nil)
	}
	invite, err := s.invites.GetByTokenHash(ctx, auth.HashOpaqueToken(token))
	if err != nil {
		var appErr *errors.AppError
		if stderrors.As(err, &appErr) && appErr.Type == errors.ErrorTypeNotFound {
			return nil, errors.ForbiddenError("Invite is invalid or expired", err)
		}
		return nil, err
	}
	if !repositories.PendingInviteFresh(invite, time.Now().UTC()) {
		return nil, errors.ForbiddenError("Invite is invalid or expired", nil)
	}
	return invite, nil
}

func (s *memberInviteService) createInvitedUser(ctx context.Context, invite *models.TenantInvite, req *dtos.AcceptMemberInviteRequest) (*models.User, error) {
	first := strings.TrimSpace(req.FirstName)
	if first == "" {
		return nil, errors.ValidationError("First name is required", nil)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), constants.BcryptCost)
	if err != nil {
		return nil, errors.InternalError("Failed to hash password", err)
	}
	user := &models.User{
		BaseModel:     models.NewBaseModel(),
		FirstName:     first,
		LastName:      strings.TrimSpace(req.LastName),
		Email:         invite.Email,
		EmailLower:    invite.EmailLower,
		EmailVerified: true,
		Status:        constants.UserStatusActive,
	}
	if err := s.users.CreateWithPasswordHash(ctx, user, string(hash)); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *memberInviteService) markInviteEmailVerified(ctx context.Context, user *models.User) error {
	if user.EmailVerified {
		return nil
	}
	if err := s.users.MarkEmailVerified(ctx, user.ID); err != nil {
		return err
	}
	user.EmailVerified = true
	return nil
}

func (s *memberInviteService) ensureActiveMembership(ctx context.Context, userID string, invite *models.TenantInvite) error {
	ok, err := s.memberships.ExistsActive(ctx, userID, invite.TenantID)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	membership := &models.TenantMembership{
		BaseModel: models.NewBaseModel(),
		UserID:    userID,
		TenantID:  invite.TenantID,
		Role:      constants.TenantMembershipRole(invite.Role),
		Status:    constants.TenantMembershipStatusActive,
	}
	return s.memberships.Create(ctx, membership)
}

func (s *adminService) inviteNewMember(ctx context.Context, tenant *models.Tenant, emailAddr, role string) error {
	if s.invites == nil {
		return errors.NotFoundError("User", nil)
	}
	emailAddr = strings.TrimSpace(emailAddr)
	emailLower := strings.ToLower(emailAddr)
	raw, hash, err := auth.NewOpaqueRefreshToken()
	if err != nil {
		return errors.InternalError("Failed to create invite token", err)
	}
	invite, err := s.invites.FindPendingByEmailAndTenant(ctx, emailLower, tenant.ID)
	if err != nil {
		return err
	}
	expires := time.Now().UTC().Add(memberInviteTTL)
	if invite == nil {
		invite = &models.TenantInvite{
			BaseModel:  models.NewBaseModel(),
			Email:      emailAddr,
			EmailLower: emailLower,
			TenantID:   tenant.ID,
			Role:       role,
			TokenHash:  hash,
			Status:     constants.TenantInviteStatusPending,
			ExpiresAt:  expires,
		}
		if err := s.invites.Create(ctx, invite); err != nil {
			return err
		}
	} else {
		invite.Email = emailAddr
		invite.Role = role
		invite.TokenHash = hash
		invite.ExpiresAt = expires
		if err := s.invites.Save(ctx, invite); err != nil {
			return err
		}
	}
	if err := s.audit.RecordRequired(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionAdminMemberAdd,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		TenantID:     tenant.ID,
		ResourceType: constants.AuditResourceTypeMembership,
		ResourceID:   invite.ID,
		NewValue:     map[string]any{"email_invited": true, "role": role},
	}); err != nil {
		return err
	}
	orgName := ""
	if tenant != nil {
		orgName = tenant.Name
	}
	if err := s.mail.SendMemberInviteEmail(ctx, emailAddr, orgName, role, s.memberInviteURL(raw), invite.ID+"/"+hash[:16]); err != nil {
		logger.From(ctx).Error("failed to send member invite email",
			zap.String("tenant_id", tenant.ID),
			zap.String("invite_id", invite.ID),
			zap.String("recipient", maskEmail(emailAddr)),
			zap.Error(err),
		)
	}
	return nil
}

func (s *adminService) ListInvites(ctx context.Context, pr *dtos.PageableRequest) ([]*dtos.AdminTenantInviteResponse, *dtos.Pageable, error) {
	return monitoring.Observe2(ctx, adminTracer, "AdminService.ListInvites", nil,
		func(ctx context.Context) ([]*dtos.AdminTenantInviteResponse, *dtos.Pageable, error) {
			if s.invites == nil {
				page := &dtos.Pageable{Page: 1, PageSize: constants.DefaultPageSize, Total: 0}
				return []*dtos.AdminTenantInviteResponse{}, page, nil
			}
			listed, err := s.invites.ListPending(ctx, pr)
			if err != nil {
				return nil, nil, err
			}
			rows, pageable := mapAdminInvites(listed)
			return rows, pageable, nil
		})
}

func (s *adminService) ListTenantInvites(ctx context.Context, tenantID string, pr *dtos.PageableRequest) ([]*dtos.AdminTenantInviteResponse, *dtos.Pageable, error) {
	return monitoring.Observe2(ctx, adminTracer, "AdminService.ListTenantInvites",
		[]attribute.KeyValue{attribute.String("tenant_id", tenantID)},
		func(ctx context.Context) ([]*dtos.AdminTenantInviteResponse, *dtos.Pageable, error) {
			return s.listTenantInvites(ctx, tenantID, pr)
		})
}

func (s *adminService) listTenantInvites(ctx context.Context, tenantID string, pr *dtos.PageableRequest) ([]*dtos.AdminTenantInviteResponse, *dtos.Pageable, error) {
	if _, err := s.tenants.GetByID(ctx, tenantID); err != nil {
		return nil, nil, err
	}
	if s.invites == nil {
		page := &dtos.Pageable{Page: 1, PageSize: constants.DefaultPageSize, Total: 0}
		return []*dtos.AdminTenantInviteResponse{}, page, nil
	}
	page, err := s.invites.ListPendingByTenant(ctx, tenantID, pr)
	if err != nil {
		return nil, nil, err
	}
	rows, pageable := mapAdminInvites(page)
	return rows, pageable, nil
}

func mapAdminInvites(page *dtos.DataResponse[models.TenantInvite]) ([]*dtos.AdminTenantInviteResponse, *dtos.Pageable) {
	now := time.Now().UTC()
	out := make([]*dtos.AdminTenantInviteResponse, 0, len(page.Data))
	for i := range page.Data {
		invite := &page.Data[i]
		status := invite.Status
		if status == constants.TenantInviteStatusPending && !now.Before(invite.ExpiresAt) {
			status = "expired"
		}
		tenantName := ""
		if invite.Tenant != nil {
			tenantName = invite.Tenant.Name
		}
		out = append(out, &dtos.AdminTenantInviteResponse{
			ID:         invite.ID,
			Email:      invite.Email,
			Role:       invite.Role,
			Status:     status,
			TenantID:   invite.TenantID,
			TenantName: tenantName,
			CreatedAt:  invite.CreatedAt,
			ExpiresAt:  invite.ExpiresAt,
		})
	}
	return out, page.Pageable
}

func (s *adminService) ResendTenantInvite(ctx context.Context, tenantID, inviteID string) error {
	return monitoring.ObserveErr(ctx, adminTracer, "AdminService.ResendTenantInvite",
		[]attribute.KeyValue{
			attribute.String("tenant_id", tenantID),
			attribute.String("invite_id", inviteID),
		},
		func(ctx context.Context) error {
			return s.resendTenantInvite(ctx, tenantID, inviteID)
		})
}

func (s *adminService) resendTenantInvite(ctx context.Context, tenantID, inviteID string) error {
	if s.invites == nil {
		return errors.NotFoundError("Tenant invite", nil)
	}
	tenant, err := s.tenants.GetByID(ctx, tenantID)
	if err != nil {
		return err
	}
	invite, err := s.invites.GetByID(ctx, inviteID)
	if err != nil {
		return err
	}
	if invite.TenantID != tenantID || invite.Status != constants.TenantInviteStatusPending {
		return errors.NotFoundError("Tenant invite", nil)
	}
	return s.inviteNewMember(ctx, tenant, invite.Email, invite.Role)
}

func (s *adminService) memberInviteURL(token string) string {
	base := s.memberSignInURL()
	if base == "" {
		base = "http://localhost:5173/login"
	}
	parsed, err := url.Parse(base)
	if err != nil {
		return base
	}
	query := parsed.Query()
	query.Set("invite", token)
	parsed.RawQuery = query.Encode()
	return parsed.String()
}
