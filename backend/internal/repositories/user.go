package repositories

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/db"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"gorm.io/gorm"
)

// UserProfilePatch updates mutable profile fields for the current user.
type UserProfilePatch struct {
	FirstName *string
	LastName  *string
}

// UserRepository defines user persistence for Phase 1 identity.
type UserRepository interface {
	CreateWithPasswordHash(ctx context.Context, user *models.User, passwordHash string) error
	// CreateUserOnly inserts a user without password_credentials (e.g. federated signup).
	CreateUserOnly(ctx context.Context, user *models.User) error
	GetOneByID(ctx context.Context, id string) (*models.User, error)
	GetByEmailLower(ctx context.Context, emailLower string) (*models.User, error)
	Count(ctx context.Context) (int64, error)
	CountPlatformAdmins(ctx context.Context) (int64, error)
	SetPlatformAdmin(ctx context.Context, userID string, isAdmin bool) error
	UpdateStatus(ctx context.Context, userID string, status constants.UserStatus) error
	UpdateProfile(ctx context.Context, userID string, patch UserProfilePatch) (*models.User, error)
	List(ctx context.Context, tenantID, search string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.User], error)
}

type userRepository struct {
	db *db.PostgresDB
}

// ProvideUserRepository creates a new user repository.
func ProvideUserRepository(db *db.PostgresDB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateWithPasswordHash(ctx context.Context, user *models.User, passwordHash string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return errors.DatabaseError("Failed to create user", err).
				WithOperation("create_user").
				WithResource("user")
		}
		pc := &models.PasswordCredential{
			BaseModel:    models.NewBaseModel(),
			UserID:       user.ID,
			PasswordHash: passwordHash,
		}
		if err := tx.Create(pc).Error; err != nil {
			return errors.DatabaseError("Failed to create password credential", err).
				WithOperation("create_user").
				WithResource("password_credential")
		}
		return nil
	})
}

func (r *userRepository) CreateUserOnly(ctx context.Context, user *models.User) error {
	return Create(ctx, r.db, user, DBOp{Operation: "create_user", Resource: "user"}, "Failed to create user")
}

func (r *userRepository) GetOneByID(ctx context.Context, id string) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFoundError("User", err).
				WithOperation("get_user").
				WithResource("user").
				WithContext("user_id", id)
		}
		return nil, errors.DatabaseError("Failed to get user by ID", err).
			WithOperation("get_user_by_id").
			WithResource("user").
			WithContext("user_id", id)
	}
	return &u, nil
}

func (r *userRepository) GetByEmailLower(ctx context.Context, emailLower string) (*models.User, error) {
	var u models.User
	err := r.db.WithContext(ctx).
		Preload("PasswordCredential").
		Where("email_lower = ?", emailLower).
		First(&u).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFoundError("User", err).
				WithOperation("get_user_by_email").
				WithResource("user")
		}
		return nil, errors.DatabaseError("Failed to get user by email", err).
			WithOperation("get_user_by_email").
			WithResource("user")
	}
	return &u, nil
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Count(&count).Error; err != nil {
		return 0, errors.DatabaseError("Failed to count users", err).
			WithOperation("count_users").
			WithResource("user")
	}
	return count, nil
}

func (r *userRepository) CountPlatformAdmins(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.User{}).Where("is_platform_admin = ?", true).Count(&count).Error; err != nil {
		return 0, errors.DatabaseError("Failed to count platform admins", err).
			WithOperation("count_platform_admins").
			WithResource("user")
	}
	return count, nil
}

func (r *userRepository) UpdateStatus(ctx context.Context, userID string, status constants.UserStatus) error {
	res := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("status", status)
	if res.Error != nil {
		return errors.DatabaseError("Failed to update user status", res.Error).
			WithOperation("update_user_status").
			WithResource("user").
			WithContext("user_id", userID)
	}
	if res.RowsAffected == 0 {
		return errors.NotFoundError("User", nil).
			WithOperation("update_user_status").
			WithResource("user").
			WithContext("user_id", userID)
	}
	return nil
}

func (r *userRepository) UpdateProfile(ctx context.Context, userID string, patch UserProfilePatch) (*models.User, error) {
	updates := map[string]any{}
	if patch.FirstName != nil {
		updates["first_name"] = *patch.FirstName
	}
	if patch.LastName != nil {
		updates["last_name"] = *patch.LastName
	}
	if len(updates) == 0 {
		return r.GetOneByID(ctx, userID)
	}

	res := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Updates(updates)
	if res.Error != nil {
		return nil, errors.DatabaseError("Failed to update user profile", res.Error).
			WithOperation("update_user_profile").
			WithResource("user").
			WithContext("user_id", userID)
	}
	if res.RowsAffected == 0 {
		return nil, errors.NotFoundError("User", nil).
			WithOperation("update_user_profile").
			WithResource("user").
			WithContext("user_id", userID)
	}
	return r.GetOneByID(ctx, userID)
}

func (r *userRepository) SetPlatformAdmin(ctx context.Context, userID string, isAdmin bool) error {
	res := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", userID).Update("is_platform_admin", isAdmin)
	if res.Error != nil {
		return errors.DatabaseError("Failed to update platform admin flag", res.Error).
			WithOperation("set_platform_admin").
			WithResource("user").
			WithContext("user_id", userID)
	}
	if res.RowsAffected == 0 {
		return errors.NotFoundError("User", nil).
			WithOperation("set_platform_admin").
			WithResource("user").
			WithContext("user_id", userID)
	}
	return nil
}

func (r *userRepository) List(ctx context.Context, tenantID, search string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.User], error) {
	query := r.db.WithContext(ctx).Model(&models.User{})
	if tenantID != "" {
		query = query.Where(
			"id IN (SELECT user_id FROM tenant_memberships WHERE tenant_id = ? AND status = ? AND deleted_at IS NULL)",
			tenantID, "active",
		)
	}
	if search = strings.TrimSpace(search); search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"email_lower LIKE ? OR first_name ILIKE ? OR last_name ILIKE ?",
			like, like, like,
		)
	}

	return PaginateWithFind[models.User](ctx, query, query.Preload("UserMFATOTP"), pr, PaginateOptions{
		OrderBy:             "created_at DESC",
		CountFailureMessage: "Failed to count users",
		FindFailureMessage:  "Failed to list users",
		DBOp:                DBOp{Operation: "list_users", Resource: "user"},
	})
}
