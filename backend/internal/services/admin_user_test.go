package services

import (
	"context"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"github.com/stretchr/testify/require"
)

type adminUserTestRepo struct {
	users map[string]*models.User
}

func (r *adminUserTestRepo) CreateWithPasswordHash(context.Context, *models.User, string) error {
	return nil
}
func (r *adminUserTestRepo) CreateUserOnly(context.Context, *models.User) error { return nil }
func (r *adminUserTestRepo) GetOneByID(_ context.Context, id string) (*models.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, errors.NotFoundError("User", nil)
	}
	return u, nil
}
func (r *adminUserTestRepo) GetByEmailLower(context.Context, string) (*models.User, error) {
	return nil, errors.NotFoundError("User", nil)
}
func (r *adminUserTestRepo) Count(context.Context) (int64, error) { return 0, nil }
func (r *adminUserTestRepo) CountPlatformAdmins(context.Context) (int64, error) {
	var n int64
	for _, u := range r.users {
		if u.IsPlatformAdmin {
			n++
		}
	}
	return n, nil
}
func (r *adminUserTestRepo) SetPlatformAdmin(context.Context, string, bool) error { return nil }
func (r *adminUserTestRepo) UpdateStatus(_ context.Context, userID string, status constants.UserStatus) error {
	u, ok := r.users[userID]
	if !ok {
		return errors.NotFoundError("User", nil)
	}
	u.Status = status
	return nil
}
func (r *adminUserTestRepo) UpdateProfile(_ context.Context, userID string, patch repositories.UserProfilePatch) (*models.User, error) {
	u, ok := r.users[userID]
	if !ok {
		return nil, errors.NotFoundError("User", nil)
	}
	if patch.FirstName != nil {
		u.FirstName = *patch.FirstName
	}
	if patch.LastName != nil {
		u.LastName = *patch.LastName
	}
	return u, nil
}
func (r *adminUserTestRepo) List(context.Context, string, string, *dtos.PageableRequest) (*dtos.DataResponse[models.User], error) {
	return nil, nil
}

type adminSessionStub struct {
	invalidated []string
}

func (s *adminSessionStub) Create(context.Context, string, string, string, string, bool) (string, time.Duration, error) {
	return "", 0, nil
}
func (s *adminSessionStub) GetSession(context.Context, string) (*models.Session, error) {
	return nil, nil
}
func (s *adminSessionStub) GetUserID(context.Context, string) (string, error) { return "", nil }
func (s *adminSessionStub) Invalidate(context.Context, string) error          { return nil }
func (s *adminSessionStub) InvalidateAllForUser(_ context.Context, userID string) error {
	s.invalidated = append(s.invalidated, userID)
	return nil
}
func (s *adminSessionStub) BrowserSessionTTL(bool) time.Duration { return time.Hour }

type adminUserSvcStub struct {
	revoked []string
}

func (s *adminUserSvcStub) Register(context.Context, *dtos.RegisterRequest, string) (*models.User, error) {
	return nil, nil
}
func (s *adminUserSvcStub) AuthenticateUser(context.Context, *dtos.LoginRequest) (*models.User, error) {
	return nil, nil
}
func (s *adminUserSvcStub) CompleteAuth(context.Context, *models.User, TenantResolveInput) (*dtos.LoginResponse, *dtos.TenantSelectionResponse, error) {
	return nil, nil, nil
}
func (s *adminUserSvcStub) SelectTenant(context.Context, *dtos.TenantSelectRequest) (*dtos.LoginResponse, error) {
	return nil, nil
}
func (s *adminUserSvcStub) SwitchTenant(context.Context, string, string) (*dtos.LoginResponse, error) {
	return nil, nil
}
func (s *adminUserSvcStub) IssueTokensForUser(context.Context, *models.User, string) (*dtos.LoginResponse, error) {
	return nil, nil
}
func (s *adminUserSvcStub) Login(context.Context, *dtos.LoginRequest, string) (*dtos.LoginResponse, *dtos.TenantSelectionResponse, error) {
	return nil, nil, nil
}
func (s *adminUserSvcStub) Refresh(context.Context, *dtos.RefreshTokenRequest) (*dtos.LoginResponse, error) {
	return nil, nil
}
func (s *adminUserSvcStub) GetOneByID(context.Context, string) (*models.User, error) { return nil, nil }
func (s *adminUserSvcStub) UpdateProfile(context.Context, string, *dtos.UpdateProfileRequest) (*models.User, error) {
	return nil, nil
}
func (s *adminUserSvcStub) ListMemberships(context.Context, string) ([]dtos.TenantSummary, error) {
	return nil, nil
}
func (s *adminUserSvcStub) ListMembershipsPaginated(context.Context, string, *dtos.PageableRequest) ([]dtos.TenantSummary, *dtos.Pageable, error) {
	return nil, nil, nil
}
func (s *adminUserSvcStub) RevokeAllRefreshTokensForUser(_ context.Context, userID string) error {
	s.revoked = append(s.revoked, userID)
	return nil
}

type adminMfaTOTPStub struct {
	disabled []string
	active   map[string]bool
}

func (s *adminMfaTOTPStub) GetByUserID(context.Context, string) (*models.UserMFATOTP, error) {
	return nil, nil
}
func (s *adminMfaTOTPStub) GetActiveByUserID(_ context.Context, userID string) (*models.UserMFATOTP, error) {
	if s.active != nil && s.active[userID] {
		return &models.UserMFATOTP{Enabled: true}, nil
	}
	return nil, nil
}
func (s *adminMfaTOTPStub) UpsertPending(context.Context, *models.UserMFATOTP) error { return nil }
func (s *adminMfaTOTPStub) MarkVerifiedAndEnabled(context.Context, string) error     { return nil }
func (s *adminMfaTOTPStub) Disable(_ context.Context, userID string) error {
	s.disabled = append(s.disabled, userID)
	return nil
}
func (s *adminMfaTOTPStub) CountEnabled(context.Context) (int64, error) { return 0, nil }

type adminMfaRecoveryStub struct {
	cleared []string
}

func (s *adminMfaRecoveryStub) ReplaceAllForUser(_ context.Context, userID string, _ []*models.UserMFARecoveryCode) error {
	s.cleared = append(s.cleared, userID)
	return nil
}
func (s *adminMfaRecoveryStub) FindUnusedByUserID(context.Context, string) ([]models.UserMFARecoveryCode, error) {
	return nil, nil
}
func (s *adminMfaRecoveryStub) MarkUsed(context.Context, string) error { return nil }

type adminWebauthnStub struct {
	deleted []string
}

func (s *adminWebauthnStub) Create(context.Context, *models.WebauthnCredential) error { return nil }
func (s *adminWebauthnStub) ListByUserID(context.Context, string) ([]models.WebauthnCredential, error) {
	return []models.WebauthnCredential{{BaseModel: models.NewBaseModel()}}, nil
}
func (s *adminWebauthnStub) ListByUserIDPaginated(_ context.Context, _ string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.WebauthnCredential], error) {
	rows := []models.WebauthnCredential{{BaseModel: models.NewBaseModel()}}
	page, pageable := dtos.PaginateSlice(rows, pr)
	return &dtos.DataResponse[models.WebauthnCredential]{Data: page, Pageable: pageable}, nil
}
func (s *adminWebauthnStub) GetByCredentialID(context.Context, string) (*models.WebauthnCredential, error) {
	return nil, nil
}
func (s *adminWebauthnStub) UpdateCredentialJSON(context.Context, string, string, int64) error {
	return nil
}
func (s *adminWebauthnStub) DeleteAllByUserID(_ context.Context, userID string) (int64, error) {
	s.deleted = append(s.deleted, userID)
	return 1, nil
}

type adminAuditCapture struct {
	params []domains.AuditRecordParams
}

func (a *adminAuditCapture) Record(_ context.Context, p domains.AuditRecordParams) {
	a.params = append(a.params, p)
}

func newAdminUserTestService(users *adminUserTestRepo) (*adminService, *adminSessionStub, *adminUserSvcStub, *adminWebauthnStub, *adminAuditCapture) {
	return newAdminUserTestServiceWithMFA(users, nil, nil)
}

func newAdminUserTestServiceWithMFA(users *adminUserTestRepo, mfaTOTP *adminMfaTOTPStub, mfaRecovery *adminMfaRecoveryStub) (*adminService, *adminSessionStub, *adminUserSvcStub, *adminWebauthnStub, *adminAuditCapture) {
	sessions := &adminSessionStub{}
	userSvc := &adminUserSvcStub{}
	webauthn := &adminWebauthnStub{}
	audit := &adminAuditCapture{}
	svc := &adminService{
		users:         users,
		sessions:      nil,
		mfaTOTP:       mfaTOTP,
		mfaRecovery:   mfaRecovery,
		webauthnCreds: webauthn,
		audit:         audit,
		sessionSvc:    sessions,
		userSvc:       userSvc,
	}
	return svc, sessions, userSvc, webauthn, audit
}

func TestAdminService_ResetMFA(t *testing.T) {
	actorID := "00000000-0000-4000-8000-000000000001"
	targetID := "00000000-0000-4000-8000-000000000002"
	users := &adminUserTestRepo{
		users: map[string]*models.User{
			targetID: {
				BaseModel: models.BaseModel{ID: targetID},
				Email:     "target@example.com",
				Status:    constants.UserStatusActive,
			},
		},
	}
	mfaTOTP := &adminMfaTOTPStub{active: map[string]bool{targetID: true}}
	mfaRecovery := &adminMfaRecoveryStub{}
	svc, _, _, _, audit := newAdminUserTestServiceWithMFA(users, mfaTOTP, mfaRecovery)

	err := svc.ResetMFA(context.Background(), actorID, targetID)
	require.NoError(t, err)
	require.Equal(t, []string{targetID}, mfaTOTP.disabled)
	require.Equal(t, []string{targetID}, mfaRecovery.cleared)
	require.Equal(t, constants.AuditActionAdminMFAReset, audit.params[0].Action)
}

func TestAdminService_DisableUser(t *testing.T) {
	actorID := "00000000-0000-4000-8000-000000000001"
	targetID := "00000000-0000-4000-8000-000000000002"
	users := &adminUserTestRepo{
		users: map[string]*models.User{
			actorID: {
				BaseModel: models.BaseModel{ID: actorID},
				Status:    constants.UserStatusActive,
			},
			targetID: {
				BaseModel: models.BaseModel{ID: targetID},
				Email:     "target@example.com",
				Status:    constants.UserStatusActive,
			},
		},
	}
	svc, sessions, userSvc, _, audit := newAdminUserTestService(users)

	err := svc.DisableUser(context.Background(), actorID, targetID)
	require.NoError(t, err)
	require.Equal(t, constants.UserStatusDisabled, users.users[targetID].Status)
	require.Equal(t, []string{targetID}, sessions.invalidated)
	require.Equal(t, []string{targetID}, userSvc.revoked)
	require.Len(t, audit.params, 1)
	require.Equal(t, constants.AuditActionAdminUserDisable, audit.params[0].Action)
}

func TestAdminService_DisableUser_BlocksSelf(t *testing.T) {
	userID := "00000000-0000-4000-8000-000000000001"
	users := &adminUserTestRepo{
		users: map[string]*models.User{
			userID: {BaseModel: models.BaseModel{ID: userID}, Status: constants.UserStatusActive},
		},
	}
	svc, _, _, _, _ := newAdminUserTestService(users)

	err := svc.DisableUser(context.Background(), userID, userID)
	require.Error(t, err)
}

func TestAdminService_DisableUser_BlocksLastPlatformAdmin(t *testing.T) {
	actorID := "00000000-0000-4000-8000-000000000001"
	targetID := "00000000-0000-4000-8000-000000000002"
	users := &adminUserTestRepo{
		users: map[string]*models.User{
			actorID: {BaseModel: models.BaseModel{ID: actorID}, Status: constants.UserStatusActive},
			targetID: {
				BaseModel:       models.BaseModel{ID: targetID},
				Status:          constants.UserStatusActive,
				IsPlatformAdmin: true,
			},
		},
	}
	svc, _, _, _, _ := newAdminUserTestService(users)

	err := svc.DisableUser(context.Background(), actorID, targetID)
	require.Error(t, err)
}

func TestAdminService_ForceLogoutUser(t *testing.T) {
	actorID := "00000000-0000-4000-8000-000000000001"
	targetID := "00000000-0000-4000-8000-000000000002"
	users := &adminUserTestRepo{
		users: map[string]*models.User{
			targetID: {
				BaseModel: models.BaseModel{ID: targetID},
				Email:     "target@example.com",
				Status:    constants.UserStatusActive,
			},
		},
	}
	svc, sessions, userSvc, _, audit := newAdminUserTestService(users)

	err := svc.ForceLogoutUser(context.Background(), actorID, targetID)
	require.NoError(t, err)
	require.Equal(t, constants.UserStatusActive, users.users[targetID].Status)
	require.Equal(t, []string{targetID}, sessions.invalidated)
	require.Equal(t, []string{targetID}, userSvc.revoked)
	require.Equal(t, constants.AuditActionAdminUserForceLogout, audit.params[0].Action)
}

func TestAdminService_ResetPasskeys(t *testing.T) {
	actorID := "00000000-0000-4000-8000-000000000001"
	targetID := "00000000-0000-4000-8000-000000000002"
	users := &adminUserTestRepo{
		users: map[string]*models.User{
			targetID: {
				BaseModel: models.BaseModel{ID: targetID},
				Email:     "target@example.com",
				Status:    constants.UserStatusActive,
			},
		},
	}
	svc, _, _, webauthn, audit := newAdminUserTestService(users)

	err := svc.ResetPasskeys(context.Background(), actorID, targetID)
	require.NoError(t, err)
	require.Equal(t, []string{targetID}, webauthn.deleted)
	require.Equal(t, constants.AuditActionAdminPasskeyReset, audit.params[0].Action)
}

func TestAdminService_ListLoginHistory(t *testing.T) {
	repo := &stubAuditLogRepo{
		list: &dtos.DataResponse[models.AuditLog]{
			Data: []models.AuditLog{
				{BaseModel: models.NewBaseModel(), Action: constants.AuditActionAuthLogin, Result: string(constants.AuditResultSuccess), ActorType: string(constants.AuditActorTypeUser)},
			},
			Pageable: &dtos.Pageable{Page: 1, PageSize: 20, Total: 1},
		},
	}
	svc := &adminService{auditLogs: repo}

	rows, pageable, err := svc.ListLoginHistory(context.Background(), dtos.AdminLoginHistoryListParams{}, dtos.NewPageableRequest())
	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.Equal(t, int64(1), pageable.Total)
}

type adminAuditLogRepoWithFilters struct {
	lastFilters repositories.AuditLogListFilters
}

func (r *adminAuditLogRepoWithFilters) Create(context.Context, *models.AuditLog) error { return nil }
func (r *adminAuditLogRepoWithFilters) List(_ context.Context, filters repositories.AuditLogListFilters, _ *dtos.PageableRequest) (*dtos.DataResponse[models.AuditLog], error) {
	r.lastFilters = filters
	return &dtos.DataResponse[models.AuditLog]{Pageable: &dtos.Pageable{}}, nil
}
func (r *adminAuditLogRepoWithFilters) Count(context.Context, repositories.AuditLogListFilters) (int64, error) {
	return 0, nil
}

func TestAdminService_ListLoginHistory_UsesLoginActions(t *testing.T) {
	repo := &adminAuditLogRepoWithFilters{}
	svc := &adminService{auditLogs: repo}

	_, _, err := svc.ListLoginHistory(context.Background(), dtos.AdminLoginHistoryListParams{}, dtos.NewPageableRequest())
	require.NoError(t, err)
	require.Equal(t, loginHistoryActions, repo.lastFilters.ActionsIn)
}
