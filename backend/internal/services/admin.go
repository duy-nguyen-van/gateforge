package services

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/crypto"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"
)

// AdminService exposes platform admin read/update operations for the console.
type AdminService interface {
	GetStats(ctx context.Context) (*dtos.AdminStatsResponse, error)
	ListUsers(ctx context.Context, tenantID, search string, pr *dtos.PageableRequest) ([]*dtos.AdminUserResponse, *dtos.Pageable, error)
	GetUserByID(ctx context.Context, userID string) (*dtos.AdminUserDetailResponse, error)
	DisableUser(ctx context.Context, actorUserID, targetUserID string) error
	ForceLogoutUser(ctx context.Context, actorUserID, targetUserID string) error
	ResetPasskeys(ctx context.Context, actorUserID, targetUserID string) error
	ResetMFA(ctx context.Context, actorUserID, targetUserID string) error
	ListTenants(ctx context.Context, pr *dtos.PageableRequest) ([]*dtos.AdminTenantResponse, *dtos.Pageable, error)
	GetTenantByID(ctx context.Context, tenantID string) (*dtos.AdminTenantResponse, error)
	CreateTenant(ctx context.Context, req *dtos.AdminCreateTenantRequest) (*dtos.AdminTenantResponse, error)
	UpdateTenant(ctx context.Context, tenantID string, req *dtos.AdminUpdateTenantRequest) (*dtos.AdminTenantResponse, error)
	DeleteTenant(ctx context.Context, tenantID string) error
	ListTenantMembers(ctx context.Context, tenantID string, pr *dtos.PageableRequest) ([]*dtos.AdminTenantMemberResponse, *dtos.Pageable, error)
	ListClients(ctx context.Context, tenantID string, pr *dtos.PageableRequest) ([]*dtos.AdminClientResponse, *dtos.Pageable, error)
	GetClientByID(ctx context.Context, clientID string) (*dtos.AdminClientResponse, error)
	CreateClient(ctx context.Context, req *dtos.AdminCreateClientRequest) (*dtos.AdminCreateClientResponse, error)
	UpdateClient(ctx context.Context, clientID string, req *dtos.AdminUpdateClientRequest) (*dtos.AdminClientResponse, error)
	DeleteClient(ctx context.Context, clientID string) error
	GetClientUsage(ctx context.Context, clientID string) (*dtos.AdminClientUsageResponse, error)
	ListIdentityProviders(ctx context.Context, tenantID string, pr *dtos.PageableRequest) ([]*dtos.AdminIdentityProviderResponse, *dtos.Pageable, error)
	ConfigureIdentityProvider(ctx context.Context, tenantID, providerID string, req *dtos.PatchIdentityProviderRequest, actorType constants.AuditActorType) error
	AddMemberByEmail(ctx context.Context, tenantID, email, role string) error
	RemoveMember(ctx context.Context, tenantID, userID string) error
	ListAuditLogs(ctx context.Context, filters dtos.AdminAuditLogListParams, pr *dtos.PageableRequest) ([]*dtos.AdminAuditLogResponse, *dtos.Pageable, error)
	ListLoginHistory(ctx context.Context, filters dtos.AdminLoginHistoryListParams, pr *dtos.PageableRequest) ([]*dtos.AdminAuditLogResponse, *dtos.Pageable, error)
}

type adminService struct {
	cfg           *config.Config
	users         repositories.UserRepository
	tenants       repositories.TenantRepository
	clients       repositories.ClientRepository
	sessions      repositories.SessionRepository
	mfaTOTP       repositories.UserMFATOTPRepository
	mfaRecovery   repositories.UserMFARecoveryCodeRepository
	tipRepo       repositories.TenantIdentityProviderRepository
	memberships   repositories.TenantMembershipRepository
	auditLogs     repositories.AuditLogRepository
	refreshTokens repositories.RefreshTokenRepository
	webauthnCreds repositories.WebauthnCredentialRepository
	audit         AuditService
	sessionSvc    SessionService
	userSvc       UserService
}

// ProvideAdminService wires admin console operations.
func ProvideAdminService(
	cfg *config.Config,
	users repositories.UserRepository,
	tenants repositories.TenantRepository,
	clients repositories.ClientRepository,
	sessions repositories.SessionRepository,
	mfaTOTP repositories.UserMFATOTPRepository,
	mfaRecovery repositories.UserMFARecoveryCodeRepository,
	tipRepo repositories.TenantIdentityProviderRepository,
	memberships repositories.TenantMembershipRepository,
	auditLogs repositories.AuditLogRepository,
	refreshTokens repositories.RefreshTokenRepository,
	webauthnCreds repositories.WebauthnCredentialRepository,
	audit AuditService,
	sessionSvc SessionService,
	userSvc UserService,
) AdminService {
	return &adminService{
		cfg:           cfg,
		users:         users,
		tenants:       tenants,
		clients:       clients,
		sessions:      sessions,
		mfaTOTP:       mfaTOTP,
		mfaRecovery:   mfaRecovery,
		tipRepo:       tipRepo,
		memberships:   memberships,
		auditLogs:     auditLogs,
		refreshTokens: refreshTokens,
		webauthnCreds: webauthnCreds,
		audit:         audit,
		sessionSvc:    sessionSvc,
		userSvc:       userSvc,
	}
}

func (s *adminService) GetStats(ctx context.Context) (*dtos.AdminStatsResponse, error) {
	totalUsers, err := s.users.Count(ctx)
	if err != nil {
		return nil, err
	}
	mfaEnabled, err := s.mfaTOTP.CountEnabled(ctx)
	if err != nil {
		return nil, err
	}
	activeSessions, err := s.sessions.CountActive(ctx)
	if err != nil {
		return nil, err
	}

	var pct float64
	if totalUsers > 0 {
		pct = math.Round(float64(mfaEnabled)/float64(totalUsers)*1000) / 10
	}

	return &dtos.AdminStatsResponse{
		TotalUsers:        totalUsers,
		MFAEnabledCount:   mfaEnabled,
		MFAEnabledPercent: pct,
		ActiveSessions:    activeSessions,
	}, nil
}

func (s *adminService) ListUsers(ctx context.Context, tenantID, search string, pr *dtos.PageableRequest) ([]*dtos.AdminUserResponse, *dtos.Pageable, error) {
	result, err := s.users.List(ctx, tenantID, search, pr)
	if err != nil {
		return nil, nil, err
	}
	out := make([]*dtos.AdminUserResponse, 0, len(result.Data))
	for i := range result.Data {
		out = append(out, dtos.NewAdminUserResponse(&result.Data[i], tenantID))
	}
	return out, result.Pageable, nil
}

func (s *adminService) ListTenants(ctx context.Context, pr *dtos.PageableRequest) ([]*dtos.AdminTenantResponse, *dtos.Pageable, error) {
	result, err := s.tenants.List(ctx, pr)
	if err != nil {
		return nil, nil, err
	}
	out := make([]*dtos.AdminTenantResponse, 0, len(result.Data))
	for i := range result.Data {
		userCount, err := s.tenants.CountUsersByTenantID(ctx, result.Data[i].ID)
		if err != nil {
			return nil, nil, err
		}
		out = append(out, dtos.NewAdminTenantResponse(&result.Data[i], userCount))
	}
	return out, result.Pageable, nil
}

func (s *adminService) ListClients(ctx context.Context, tenantID string, pr *dtos.PageableRequest) ([]*dtos.AdminClientResponse, *dtos.Pageable, error) {
	result, err := s.clients.List(ctx, tenantID, pr)
	if err != nil {
		return nil, nil, err
	}
	out := make([]*dtos.AdminClientResponse, 0, len(result.Data))
	for i := range result.Data {
		out = append(out, dtos.NewAdminClientResponse(&result.Data[i]))
	}
	return out, result.Pageable, nil
}

func (s *adminService) ListIdentityProviders(ctx context.Context, tenantID string, pr *dtos.PageableRequest) ([]*dtos.AdminIdentityProviderResponse, *dtos.Pageable, error) {
	out := make([]*dtos.AdminIdentityProviderResponse, 0, len(constants.SupportedIdentityProviders))
	for _, spec := range constants.SupportedIdentityProviders {
		tip, err := s.tipRepo.GetByTenantAndProvider(ctx, tenantID, spec.ID)
		if err != nil {
			if appErr := errors.GetAppError(err); appErr != nil && appErr.Type == errors.ErrorTypeNotFound {
				out = append(out, s.newIdentityProviderResponse(spec, tenantID, nil))
				continue
			}
			return nil, nil, err
		}
		out = append(out, s.newIdentityProviderResponse(spec, tenantID, tip))
	}
	page, pageable := dtos.PaginateSlice(out, pr)
	return page, pageable, nil
}

func (s *adminService) newIdentityProviderResponse(spec constants.IdentityProviderSpec, tenantID string, tip *models.TenantIdentityProvider) *dtos.AdminIdentityProviderResponse {
	resp := &dtos.AdminIdentityProviderResponse{
		Provider:        spec.ID,
		Name:            spec.DisplayName,
		TenantID:        tenantID,
		RedirectURI:     FederationCallbackURL(s.cfg, spec.ID),
		SetupConsoleURL: spec.SetupConsoleURL,
	}
	if tip != nil {
		resp.Enabled = tip.Enabled
		resp.OAuthClientID = strings.TrimSpace(tip.OAuthClientID)
		resp.OAuthClientSecretSet = strings.TrimSpace(tip.OAuthClientSecretEncrypted) != ""
		resp.Configured = resp.OAuthClientID != "" && resp.OAuthClientSecretSet
	}
	return resp
}

func (s *adminService) ConfigureIdentityProvider(ctx context.Context, tenantID, providerID string, req *dtos.PatchIdentityProviderRequest, actorType constants.AuditActorType) error {
	if req == nil {
		return errors.ValidationError("Request body is required", nil)
	}
	spec, ok := constants.IdentityProviderByID(providerID)
	if !ok {
		return errors.ValidationError("Unsupported identity provider", nil)
	}

	existing, existingErr := s.tipRepo.GetByTenantAndProvider(ctx, tenantID, spec.ID)
	hasExisting := existingErr == nil && existing != nil

	clientID := strings.TrimSpace(req.OAuthClientID)
	if clientID == "" && hasExisting {
		clientID = strings.TrimSpace(existing.OAuthClientID)
	}
	secretSet := hasExisting && strings.TrimSpace(existing.OAuthClientSecretEncrypted) != ""
	if strings.TrimSpace(req.OAuthClientSecret) != "" {
		secretSet = true
	}

	enable := false
	if req.Enabled != nil {
		enable = *req.Enabled
	} else if hasExisting {
		enable = existing.Enabled
	}
	if enable && (clientID == "" || !secretSet) {
		return errors.ValidationError(
			fmt.Sprintf("%s OAuth client ID and secret are required before enabling", spec.DisplayName),
			nil,
		)
	}

	patch := repositories.TenantIdentityProviderPatch{Enabled: req.Enabled}
	if strings.TrimSpace(req.OAuthClientID) != "" {
		patch.OAuthClientID = strings.TrimSpace(req.OAuthClientID)
	}
	if plain := strings.TrimSpace(req.OAuthClientSecret); plain != "" {
		enc, err := crypto.EncryptMFASecret(federationEncryptionKey(s.cfg), plain)
		if err != nil {
			return errors.InternalError(fmt.Sprintf("Failed to encrypt %s OAuth secret", spec.DisplayName), err)
		}
		patch.OAuthClientSecretEncrypted = enc
	}

	prevEnabled := false
	if hasExisting {
		prevEnabled = existing.Enabled
	}
	updated, err := s.tipRepo.UpdateProvider(ctx, tenantID, spec.ID, patch)
	if err != nil {
		return err
	}

	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionAdminIDPPatch,
		Result:       constants.AuditResultSuccess,
		ActorType:    actorType,
		TenantID:     tenantID,
		ResourceType: constants.AuditResourceTypeIdentityProvider,
		ResourceName: spec.ID,
		OldValue:     map[string]any{"enabled": prevEnabled},
		NewValue: map[string]any{
			"enabled":                 updated.Enabled,
			"oauth_client_id_set":     strings.TrimSpace(updated.OAuthClientID) != "",
			"oauth_client_secret_set": strings.TrimSpace(updated.OAuthClientSecretEncrypted) != "",
		},
	})
	return nil
}

func (s *adminService) AddMemberByEmail(ctx context.Context, tenantID, email, role string) error {
	if _, err := s.tenants.GetByID(ctx, tenantID); err != nil {
		return err
	}
	u, err := s.users.GetByEmailLower(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		return err
	}
	if role == "" {
		role = string(constants.TenantMembershipRoleMember)
	}
	ok, err := s.memberships.ExistsActive(ctx, u.ID, tenantID)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	membership := &models.TenantMembership{
		BaseModel: models.NewBaseModel(),
		UserID:    u.ID,
		TenantID:  tenantID,
		Role:      constants.TenantMembershipRole(role),
		Status:    constants.TenantMembershipStatusActive,
	}
	if err := s.memberships.Create(ctx, membership); err != nil {
		return err
	}
	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionAdminMemberAdd,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		TenantID:     tenantID,
		ResourceType: constants.AuditResourceTypeMembership,
		ResourceID:   membership.ID,
		ResourceName: u.ID,
		NewValue:     map[string]any{"user_id": u.ID, "role": role},
	})
	return nil
}

func (s *adminService) RemoveMember(ctx context.Context, tenantID, userID string) error {
	if err := s.memberships.Delete(ctx, userID, tenantID); err != nil {
		return err
	}
	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionAdminMemberRemove,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		TenantID:     tenantID,
		ResourceType: constants.AuditResourceTypeMembership,
		ResourceName: userID,
		OldValue:     map[string]any{"user_id": userID},
	})
	return nil
}

func (s *adminService) ListAuditLogs(ctx context.Context, filters dtos.AdminAuditLogListParams, pr *dtos.PageableRequest) ([]*dtos.AdminAuditLogResponse, *dtos.Pageable, error) {
	result, err := s.auditLogs.List(ctx, repositories.AuditLogListFilters{
		TenantID: filters.TenantID,
		Action:   filters.Action,
		Result:   filters.Result,
		ActorID:  filters.ActorID,
		From:     filters.From,
		To:       filters.To,
	}, pr)
	if err != nil {
		return nil, nil, err
	}
	out := make([]*dtos.AdminAuditLogResponse, 0, len(result.Data))
	for i := range result.Data {
		out = append(out, dtos.NewAdminAuditLogResponse(&result.Data[i]))
	}
	return out, result.Pageable, nil
}
