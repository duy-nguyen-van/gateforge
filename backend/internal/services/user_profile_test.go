package services

import (
	"context"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"github.com/stretchr/testify/require"
)

type profileTestUserRepo struct {
	users map[string]*models.User
}

func newProfileTestUserRepo() *profileTestUserRepo {
	return &profileTestUserRepo{users: map[string]*models.User{}}
}

func (r *profileTestUserRepo) CreateWithPasswordHash(_ context.Context, _ *models.User, _ string) error {
	return nil
}

func (r *profileTestUserRepo) CreateUserOnly(_ context.Context, _ *models.User) error {
	return nil
}

func (r *profileTestUserRepo) GetOneByID(_ context.Context, id string) (*models.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, errors.NotFoundError("User", nil)
	}
	return u, nil
}

func (r *profileTestUserRepo) GetByEmailLower(_ context.Context, _ string) (*models.User, error) {
	return nil, errors.NotFoundError("User", nil)
}

func (r *profileTestUserRepo) Count(_ context.Context) (int64, error) {
	return int64(len(r.users)), nil
}

func (r *profileTestUserRepo) CountPlatformAdmins(_ context.Context) (int64, error) {
	return 0, nil
}

func (r *profileTestUserRepo) SetPlatformAdmin(_ context.Context, _ string, _ bool) error {
	return nil
}

func (r *profileTestUserRepo) UpdateStatus(_ context.Context, _ string, _ constants.UserStatus) error {
	return nil
}

func (r *profileTestUserRepo) UpdateProfile(_ context.Context, userID string, patch repositories.UserProfilePatch) (*models.User, error) {
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

func (r *profileTestUserRepo) List(_ context.Context, _, _ string, _ *dtos.PageableRequest) (*dtos.DataResponse[models.User], error) {
	return nil, nil
}

func TestUserService_UpdateProfile(t *testing.T) {
	repo := newProfileTestUserRepo()
	user := &models.User{
		BaseModel: models.NewBaseModel(),
		FirstName: "Ada",
		LastName:  "Lovelace",
		Email:     "ada@example.com",
		Status:    constants.UserStatusActive,
	}
	repo.users[user.ID] = user

	svc := &userService{userRepo: repo}

	firstName := "Grace"
	lastName := "Hopper"
	updated, err := svc.UpdateProfile(context.Background(), user.ID, &dtos.UpdateProfileRequest{
		FirstName: &firstName,
		LastName:  &lastName,
	})
	require.NoError(t, err)
	require.Equal(t, "Grace", updated.FirstName)
	require.Equal(t, "Hopper", updated.LastName)

	_, err = svc.UpdateProfile(context.Background(), user.ID, &dtos.UpdateProfileRequest{})
	require.Error(t, err)

	empty := "   "
	_, err = svc.UpdateProfile(context.Background(), user.ID, &dtos.UpdateProfileRequest{FirstName: &empty})
	require.Error(t, err)
}
