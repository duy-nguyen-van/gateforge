package services

import (
	"context"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	apperrors "github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"
)

var loginHistoryActions = []string{
	constants.AuditActionAuthLogin,
	constants.AuditActionWebauthnLogin,
	constants.AuditActionFederationLogin,
	constants.AuditActionOIDCLogin,
}

func (s *adminService) GetUserByID(ctx context.Context, userID string) (*dtos.AdminUserDetailResponse, error) {
	u, err := s.users.GetOneByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	mfaEnabled := false
	if totp, err := s.mfaTOTP.GetActiveByUserID(ctx, userID); err == nil && totp != nil && totp.Enabled {
		mfaEnabled = true
	}

	passkeys, err := s.webauthnCreds.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	activeSessions, err := s.sessions.CountActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	memberships, err := s.memberships.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	membershipOut := make([]dtos.AdminUserMembershipResponse, 0, len(memberships))
	for _, m := range memberships {
		entry := dtos.AdminUserMembershipResponse{
			TenantID: m.TenantID,
			Role:     string(m.Role),
			Status:   string(m.Status),
		}
		if tenant, err := s.tenants.GetByID(ctx, m.TenantID); err == nil && tenant != nil {
			entry.TenantName = tenant.Name
		}
		membershipOut = append(membershipOut, entry)
	}

	return &dtos.AdminUserDetailResponse{
		ID:              u.ID,
		Email:           u.Email,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		Status:          string(u.Status),
		MFAEnabled:      mfaEnabled,
		IsPlatformAdmin: u.IsPlatformAdmin,
		PasskeyCount:    len(passkeys),
		ActiveSessions:  activeSessions,
		Memberships:     membershipOut,
		CreatedAt:       u.CreatedAt,
	}, nil
}

func (s *adminService) DisableUser(ctx context.Context, actorUserID, targetUserID string) error {
	if actorUserID == targetUserID {
		return apperrors.ForbiddenError("You cannot disable your own account", nil)
	}

	target, err := s.users.GetOneByID(ctx, targetUserID)
	if err != nil {
		return err
	}

	if target.Status == constants.UserStatusDisabled {
		return nil
	}

	if target.IsPlatformAdmin {
		adminCount, err := s.users.CountPlatformAdmins(ctx)
		if err != nil {
			return err
		}
		if adminCount <= 1 {
			return apperrors.ForbiddenError("Cannot disable the last platform admin", nil)
		}
	}

	oldStatus := target.Status
	if err := s.users.UpdateStatus(ctx, targetUserID, constants.UserStatusDisabled); err != nil {
		return err
	}
	if err := s.sessionSvc.InvalidateAllForUser(ctx, targetUserID); err != nil {
		return err
	}
	if err := s.userSvc.RevokeAllRefreshTokensForUser(ctx, targetUserID); err != nil {
		return err
	}

	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionAdminUserDisable,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      actorUserID,
		ResourceType: constants.AuditResourceTypeUser,
		ResourceID:   targetUserID,
		ResourceName: target.Email,
		OldValue:     map[string]any{"status": string(oldStatus)},
		NewValue:     map[string]any{"status": string(constants.UserStatusDisabled)},
	})
	return nil
}

func (s *adminService) ForceLogoutUser(ctx context.Context, actorUserID, targetUserID string) error {
	target, err := s.users.GetOneByID(ctx, targetUserID)
	if err != nil {
		return err
	}

	if err := s.sessionSvc.InvalidateAllForUser(ctx, targetUserID); err != nil {
		return err
	}
	if err := s.userSvc.RevokeAllRefreshTokensForUser(ctx, targetUserID); err != nil {
		return err
	}

	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionAdminUserForceLogout,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      actorUserID,
		ResourceType: constants.AuditResourceTypeUser,
		ResourceID:   targetUserID,
		ResourceName: target.Email,
	})
	return nil
}

func (s *adminService) ResetMFA(ctx context.Context, actorUserID, targetUserID string) error {
	target, err := s.users.GetOneByID(ctx, targetUserID)
	if err != nil {
		return err
	}

	mfaEnabled := false
	if totp, err := s.mfaTOTP.GetActiveByUserID(ctx, targetUserID); err == nil && totp != nil && totp.Enabled {
		mfaEnabled = true
	}

	if err := s.mfaTOTP.Disable(ctx, targetUserID); err != nil {
		return err
	}
	if err := s.mfaRecovery.ReplaceAllForUser(ctx, targetUserID, nil); err != nil {
		return err
	}

	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionAdminMFAReset,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      actorUserID,
		ResourceType: constants.AuditResourceTypeUser,
		ResourceID:   targetUserID,
		ResourceName: target.Email,
		OldValue:     map[string]any{"mfa_enabled": mfaEnabled},
		NewValue:     map[string]any{"mfa_enabled": false},
	})
	return nil
}

func (s *adminService) ResetPasskeys(ctx context.Context, actorUserID, targetUserID string) error {
	target, err := s.users.GetOneByID(ctx, targetUserID)
	if err != nil {
		return err
	}

	passkeys, err := s.webauthnCreds.ListByUserID(ctx, targetUserID)
	if err != nil {
		return err
	}
	count := len(passkeys)

	deleted, err := s.webauthnCreds.DeleteAllByUserID(ctx, targetUserID)
	if err != nil {
		return err
	}

	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionAdminPasskeyReset,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      actorUserID,
		ResourceType: constants.AuditResourceTypeUser,
		ResourceID:   targetUserID,
		ResourceName: target.Email,
		OldValue:     map[string]any{"passkey_count": count},
		NewValue:     map[string]any{"deleted_count": deleted},
	})
	return nil
}

func (s *adminService) GetClientUsage(ctx context.Context, clientID string) (*dtos.AdminClientUsageResponse, error) {
	client, err := s.clients.GetByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	usage, err := s.refreshTokens.UsageByClientRecordID(ctx, client.ID)
	if err != nil {
		return nil, err
	}

	since := time.Now().UTC().Add(-30 * 24 * time.Hour)
	baseFilters := repositories.AuditLogListFilters{
		ResourceType: string(constants.AuditResourceTypeClient),
		ResourceName: client.ClientID,
		From:         &since,
	}

	authorizeCount, err := s.auditLogs.Count(ctx, repositories.AuditLogListFilters{
		ResourceType: baseFilters.ResourceType,
		ResourceName: baseFilters.ResourceName,
		From:         baseFilters.From,
		Action:       constants.AuditActionOIDCAuthorize,
	})
	if err != nil {
		return nil, err
	}

	tokenIssueCount, err := s.auditLogs.Count(ctx, repositories.AuditLogListFilters{
		ResourceType: baseFilters.ResourceType,
		ResourceName: baseFilters.ResourceName,
		From:         baseFilters.From,
		Action:       constants.AuditActionOIDCTokenIssue,
	})
	if err != nil {
		return nil, err
	}

	return &dtos.AdminClientUsageResponse{
		ClientID:            client.ClientID,
		TotalRefreshTokens:  usage.TotalIssued,
		ActiveRefreshTokens: usage.ActiveCount,
		LastTokenIssuedAt:   usage.LastIssuedAt,
		AuthorizeEvents30d:  authorizeCount,
		TokenIssueEvents30d: tokenIssueCount,
	}, nil
}

func (s *adminService) ListLoginHistory(ctx context.Context, filters dtos.AdminLoginHistoryListParams, pr *dtos.PageableRequest) ([]*dtos.AdminAuditLogResponse, *dtos.Pageable, error) {
	result, err := s.auditLogs.List(ctx, repositories.AuditLogListFilters{
		TenantID:  filters.TenantID,
		ActionsIn: loginHistoryActions,
		Result:    filters.Result,
		ActorID:   filters.ActorID,
		From:      filters.From,
		To:        filters.To,
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
