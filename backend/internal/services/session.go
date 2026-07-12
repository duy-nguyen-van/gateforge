package services

import (
	"context"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	apperrors "github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/logger"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"go.uber.org/zap"
)

// SessionService creates and resolves browser sessions shared by the dashboard API and OIDC /authorize.
type SessionService interface {
	Create(ctx context.Context, userID, tenantID, ip, userAgent string, remember bool) (sessionID string, ttl time.Duration, err error)
	GetSession(ctx context.Context, sessionID string) (*models.Session, error)
	GetUserID(ctx context.Context, sessionID string) (userID string, err error)
	Invalidate(ctx context.Context, sessionID string) error
	InvalidateAllForUser(ctx context.Context, userID string) error
	BrowserSessionTTL(remember bool) time.Duration
}

type sessionService struct {
	cfg   *config.Config
	repo  repositories.SessionRepository
	audit AuditService
}

// ProvideSessionService wires session creation for cookie-based login.
func ProvideSessionService(cfg *config.Config, repo repositories.SessionRepository, audit AuditService) SessionService {
	return &sessionService{cfg: cfg, repo: repo, audit: audit}
}

func (s *sessionService) BrowserSessionTTL(remember bool) time.Duration {
	if remember {
		if s.cfg.SSOSessionRememberTTL > 0 {
			return s.cfg.SSOSessionRememberTTL
		}
		return 720 * time.Hour
	}
	if s.cfg.SSOSessionTTL > 0 {
		return s.cfg.SSOSessionTTL
	}
	return s.cfg.JWTRefreshTTL
}

func (s *sessionService) Create(ctx context.Context, userID, tenantID, ip, userAgent string, remember bool) (string, time.Duration, error) {
	ttl := s.BrowserSessionTTL(remember)
	exp := time.Now().UTC().Add(ttl)
	row := &models.Session{
		UserID:    userID,
		TenantID:  tenantID,
		IPAddress: ip,
		UserAgent: userAgent,
		ExpiresAt: &exp,
	}
	if err := s.repo.Create(ctx, row); err != nil {
		logger.Log.Error("session create failed",
			zap.String("operation", "session_create"),
			zap.String("user_id", userID),
			zap.String("tenant_id", tenantID),
			zap.Error(err))
		return "", 0, err
	}
	logger.Log.Info("session created",
		zap.String("operation", "session_create"),
		zap.String("user_id", userID),
		zap.String("tenant_id", tenantID),
		zap.Bool("remember", remember))
	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionSessionCreate,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      userID,
		TenantID:     tenantID,
		ResourceType: constants.AuditResourceTypeSession,
		ResourceID:   row.ID,
		NewValue:     map[string]any{"remember": remember},
	})
	return row.ID, ttl, nil
}

func (s *sessionService) GetSession(ctx context.Context, sessionID string) (*models.Session, error) {
	if sessionID == "" {
		return nil, apperrors.UnauthorizedError("invalid session", nil)
	}
	sess, err := s.repo.GetValidByID(ctx, sessionID)
	if err != nil {
		return nil, apperrors.UnauthorizedError("invalid or expired session", nil)
	}
	return sess, nil
}

func (s *sessionService) GetUserID(ctx context.Context, sessionID string) (string, error) {
	sess, err := s.GetSession(ctx, sessionID)
	if err != nil {
		return "", err
	}
	return sess.UserID, nil
}

func (s *sessionService) Invalidate(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return nil
	}
	if err := s.repo.DeleteByID(ctx, sessionID); err != nil {
		logger.Log.Error("session invalidate failed",
			zap.String("operation", "session_invalidate"),
			zap.Error(err))
		return err
	}
	logger.Log.Info("session invalidated",
		zap.String("operation", "session_invalidate"))
	return nil
}

func (s *sessionService) InvalidateAllForUser(ctx context.Context, userID string) error {
	if userID == "" {
		return nil
	}
	if err := s.repo.DeleteAllByUserID(ctx, userID); err != nil {
		logger.Log.Error("session invalidate all for user failed",
			zap.String("operation", "session_invalidate_all"),
			zap.String("user_id", userID),
			zap.Error(err))
		return err
	}
	logger.Log.Info("all sessions invalidated for user",
		zap.String("operation", "session_invalidate_all"),
		zap.String("user_id", userID))
	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionSessionRevokeAll,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      userID,
		ResourceType: constants.AuditResourceTypeUser,
		ResourceID:   userID,
	})
	return nil
}
