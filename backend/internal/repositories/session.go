package repositories

import (
	"context"
	stderrors "errors"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/db"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"gorm.io/gorm"
)

// SessionRepository persists browser sessions for the OIDC authorization-code + cookie login flow.
type SessionRepository interface {
	Create(ctx context.Context, s *models.Session) error
	GetValidByID(ctx context.Context, id string) (*models.Session, error)
	DeleteByID(ctx context.Context, id string) error
	DeleteAllByUserID(ctx context.Context, userID string) error
	CountActive(ctx context.Context) (int64, error)
	CountActiveByUserID(ctx context.Context, userID string) (int64, error)
}

type sessionRepository struct {
	db *db.PostgresDB
}

// ProvideSessionRepository wires session persistence.
func ProvideSessionRepository(db *db.PostgresDB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, s *models.Session) error {
	return Create(ctx, r.db, s, DBOp{Operation: "create_session", Resource: "session"}, "Failed to create session")
}

func (r *sessionRepository) GetValidByID(ctx context.Context, id string) (*models.Session, error) {
	var s models.Session
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&s).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFoundError("Session", err).
				WithOperation("get_session").
				WithResource("session")
		}
		return nil, errors.DatabaseError("Failed to load session", err).
			WithOperation("get_session").
			WithResource("session")
	}
	if s.ExpiresAt != nil && time.Now().UTC().After(*s.ExpiresAt) {
		return nil, errors.NotFoundError("Session", stderrors.New("session expired")).
			WithOperation("get_session").
			WithResource("session")
	}
	return &s, nil
}

func (r *sessionRepository) DeleteByID(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Session{})
	if res.Error != nil {
		return errors.DatabaseError("Failed to delete session", res.Error).
			WithOperation("delete_session").
			WithResource("session")
	}
	return nil
}

func (r *sessionRepository) DeleteAllByUserID(ctx context.Context, userID string) error {
	res := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.Session{})
	if res.Error != nil {
		return errors.DatabaseError("Failed to delete user sessions", res.Error).
			WithOperation("delete_sessions_by_user").
			WithResource("session")
	}
	return nil
}

func (r *sessionRepository) CountActive(ctx context.Context) (int64, error) {
	var count int64
	now := time.Now().UTC()
	err := r.db.WithContext(ctx).Model(&models.Session{}).
		Where("expires_at IS NULL OR expires_at > ?", now).
		Count(&count).Error
	if err != nil {
		return 0, errors.DatabaseError("Failed to count active sessions", err).
			WithOperation("count_active_sessions").
			WithResource("session")
	}
	return count, nil
}

func (r *sessionRepository) CountActiveByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64
	now := time.Now().UTC()
	err := r.db.WithContext(ctx).Model(&models.Session{}).
		Where("user_id = ? AND (expires_at IS NULL OR expires_at > ?)", userID, now).
		Count(&count).Error
	if err != nil {
		return 0, errors.DatabaseError("Failed to count active sessions for user", err).
			WithOperation("count_active_sessions_user").
			WithResource("session")
	}
	return count, nil
}
