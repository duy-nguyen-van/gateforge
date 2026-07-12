package services

import (
	"context"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"github.com/stretchr/testify/require"
)

type stubBootstrapUserRepo struct {
	platformAdminCount int64
	users              map[string]*models.User
	byEmail            map[string]*models.User
	created            []*models.User
	promoted           map[string]bool
}

func (s *stubBootstrapUserRepo) CreateWithPasswordHash(ctx context.Context, user *models.User, passwordHash string) error {
	s.created = append(s.created, user)
	s.users[user.ID] = user
	s.byEmail[user.EmailLower] = user
	return nil
}

func (s *stubBootstrapUserRepo) CreateUserOnly(ctx context.Context, user *models.User) error {
	return nil
}

func (s *stubBootstrapUserRepo) GetOneByID(ctx context.Context, id string) (*models.User, error) {
	u, ok := s.users[id]
	if !ok {
		return nil, errors.NotFoundError("User", nil)
	}
	return u, nil
}

func (s *stubBootstrapUserRepo) GetByEmailLower(ctx context.Context, emailLower string) (*models.User, error) {
	u, ok := s.byEmail[emailLower]
	if !ok {
		return nil, errors.NotFoundError("User", nil)
	}
	return u, nil
}

func (s *stubBootstrapUserRepo) Count(ctx context.Context) (int64, error) {
	return int64(len(s.users)), nil
}

func (s *stubBootstrapUserRepo) CountPlatformAdmins(ctx context.Context) (int64, error) {
	return s.platformAdminCount, nil
}

func (s *stubBootstrapUserRepo) SetPlatformAdmin(ctx context.Context, userID string, isAdmin bool) error {
	u, ok := s.users[userID]
	if !ok {
		return errors.NotFoundError("User", nil)
	}
	u.IsPlatformAdmin = isAdmin
	s.promoted[userID] = isAdmin
	if isAdmin {
		s.platformAdminCount++
	}
	return nil
}

func (s *stubBootstrapUserRepo) UpdateStatus(ctx context.Context, userID string, status constants.UserStatus) error {
	u, ok := s.users[userID]
	if !ok {
		return errors.NotFoundError("User", nil)
	}
	u.Status = status
	return nil
}

func (s *stubBootstrapUserRepo) UpdateProfile(_ context.Context, userID string, patch repositories.UserProfilePatch) (*models.User, error) {
	u, ok := s.users[userID]
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

func (s *stubBootstrapUserRepo) List(ctx context.Context, tenantID, search string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.User], error) {
	return nil, nil
}

type stubBootstrapMembershipRepo struct {
	created []*models.TenantMembership
}

func (s *stubBootstrapMembershipRepo) Create(ctx context.Context, m *models.TenantMembership) error {
	s.created = append(s.created, m)
	return nil
}

func (s *stubBootstrapMembershipRepo) GetActive(ctx context.Context, userID, tenantID string) (*models.TenantMembership, error) {
	return nil, errors.NotFoundError("Tenant membership", nil)
}

func (s *stubBootstrapMembershipRepo) ExistsActive(ctx context.Context, userID, tenantID string) (bool, error) {
	return false, nil
}

func (s *stubBootstrapMembershipRepo) ListByUserID(ctx context.Context, userID string) ([]models.TenantMembership, error) {
	return nil, nil
}

func (s *stubBootstrapMembershipRepo) ListByUserIDPaginated(context.Context, string, *dtos.PageableRequest) (*dtos.DataResponse[models.TenantMembership], error) {
	return &dtos.DataResponse[models.TenantMembership]{Pageable: &dtos.Pageable{}}, nil
}

func (s *stubBootstrapMembershipRepo) ListByTenantID(ctx context.Context, tenantID string) ([]models.TenantMembership, error) {
	return nil, nil
}

func (s *stubBootstrapMembershipRepo) ListByTenantIDPaginated(context.Context, string, *dtos.PageableRequest) (*dtos.DataResponse[models.TenantMembership], error) {
	return &dtos.DataResponse[models.TenantMembership]{Pageable: &dtos.Pageable{}}, nil
}

func (s *stubBootstrapMembershipRepo) CountByTenantID(ctx context.Context, tenantID string) (int64, error) {
	return 0, nil
}

func (s *stubBootstrapMembershipRepo) Delete(ctx context.Context, userID, tenantID string) error {
	return nil
}

func TestPlatformAdminBootstrap_ShortPassword(t *testing.T) {
	repo := &stubBootstrapUserRepo{users: map[string]*models.User{}, byEmail: map[string]*models.User{}, promoted: map[string]bool{}}
	b := ProvidePlatformAdminBootstrap(&config.Config{
		BootstrapAdminEmail:    "admin@example.com",
		BootstrapAdminPassword: "short",
	}, repo, &stubBootstrapMembershipRepo{})

	err := b.Run(context.Background())
	require.Error(t, err)
}

func TestPlatformAdminBootstrap_NoCredentialsLogsAndSkips(t *testing.T) {
	repo := &stubBootstrapUserRepo{users: map[string]*models.User{}, byEmail: map[string]*models.User{}, promoted: map[string]bool{}}
	b := ProvidePlatformAdminBootstrap(&config.Config{}, repo, &stubBootstrapMembershipRepo{})
	require.NoError(t, b.Run(context.Background()))
}

func TestPlatformAdminBootstrap_ExistingMembershipSkipped(t *testing.T) {
	existing := &models.User{
		BaseModel:  models.BaseModel{ID: "user-1"},
		Email:      "admin@example.com",
		EmailLower: "admin@example.com",
	}
	repo := &stubBootstrapUserRepo{
		users:    map[string]*models.User{"user-1": existing},
		byEmail:  map[string]*models.User{"admin@example.com": existing},
		promoted: map[string]bool{},
	}
	memberships := &stubBootstrapMembershipStub{exists: true}
	b := ProvidePlatformAdminBootstrap(&config.Config{
		DefaultTenantID:        "tenant-1",
		BootstrapAdminEmail:    "admin@example.com",
		BootstrapAdminPassword: "password123",
	}, repo, memberships)

	require.NoError(t, b.Run(context.Background()))
	require.Empty(t, memberships.created)
}

type stubBootstrapMembershipStub struct {
	exists  bool
	created []*models.TenantMembership
}

func (s *stubBootstrapMembershipStub) Create(ctx context.Context, m *models.TenantMembership) error {
	s.created = append(s.created, m)
	return nil
}
func (s *stubBootstrapMembershipStub) GetActive(context.Context, string, string) (*models.TenantMembership, error) {
	return nil, errors.NotFoundError("membership", nil)
}
func (s *stubBootstrapMembershipStub) ExistsActive(context.Context, string, string) (bool, error) {
	return s.exists, nil
}
func (s *stubBootstrapMembershipStub) ListByUserID(context.Context, string) ([]models.TenantMembership, error) {
	return nil, nil
}
func (s *stubBootstrapMembershipStub) ListByUserIDPaginated(context.Context, string, *dtos.PageableRequest) (*dtos.DataResponse[models.TenantMembership], error) {
	return &dtos.DataResponse[models.TenantMembership]{Pageable: &dtos.Pageable{}}, nil
}
func (s *stubBootstrapMembershipStub) ListByTenantID(context.Context, string) ([]models.TenantMembership, error) {
	return nil, nil
}
func (s *stubBootstrapMembershipStub) ListByTenantIDPaginated(context.Context, string, *dtos.PageableRequest) (*dtos.DataResponse[models.TenantMembership], error) {
	return &dtos.DataResponse[models.TenantMembership]{Pageable: &dtos.Pageable{}}, nil
}
func (s *stubBootstrapMembershipStub) CountByTenantID(context.Context, string) (int64, error) {
	return 0, nil
}
func (s *stubBootstrapMembershipStub) Delete(context.Context, string, string) error { return nil }

func TestPlatformAdminBootstrap_SkipsWhenAdminsExist(t *testing.T) {
	repo := &stubBootstrapUserRepo{platformAdminCount: 1, users: map[string]*models.User{}, byEmail: map[string]*models.User{}, promoted: map[string]bool{}}
	memberships := &stubBootstrapMembershipRepo{}
	b := ProvidePlatformAdminBootstrap(&config.Config{
		BootstrapAdminEmail:    "admin@example.com",
		BootstrapAdminPassword: "password123",
	}, repo, memberships)

	require.NoError(t, b.Run(context.Background()))
	require.Empty(t, repo.created)
}

func TestPlatformAdminBootstrap_CreatesNewAdmin(t *testing.T) {
	repo := &stubBootstrapUserRepo{users: map[string]*models.User{}, byEmail: map[string]*models.User{}, promoted: map[string]bool{}}
	memberships := &stubBootstrapMembershipRepo{}
	b := ProvidePlatformAdminBootstrap(&config.Config{
		DefaultTenantID:        "tenant-1",
		BootstrapAdminEmail:    "admin@example.com",
		BootstrapAdminPassword: "password123",
	}, repo, memberships)

	require.NoError(t, b.Run(context.Background()))
	require.Len(t, repo.created, 1)
	require.True(t, repo.created[0].IsPlatformAdmin)
	require.Equal(t, constants.UserStatusActive, repo.created[0].Status)
	require.Len(t, memberships.created, 1)
	require.Equal(t, "tenant-1", memberships.created[0].TenantID)
}

func TestPlatformAdminBootstrap_PromotesExistingUser(t *testing.T) {
	existing := &models.User{
		BaseModel:  models.BaseModel{ID: "user-1"},
		Email:      "admin@example.com",
		EmailLower: "admin@example.com",
	}
	repo := &stubBootstrapUserRepo{
		users:    map[string]*models.User{"user-1": existing},
		byEmail:  map[string]*models.User{"admin@example.com": existing},
		promoted: map[string]bool{},
	}
	memberships := &stubBootstrapMembershipRepo{}
	b := ProvidePlatformAdminBootstrap(&config.Config{
		DefaultTenantID:        "tenant-1",
		BootstrapAdminEmail:    "admin@example.com",
		BootstrapAdminPassword: "password123",
	}, repo, memberships)

	require.NoError(t, b.Run(context.Background()))
	require.Empty(t, repo.created)
	require.True(t, repo.promoted["user-1"])
	require.True(t, existing.IsPlatformAdmin)
	require.Len(t, memberships.created, 1)
}
