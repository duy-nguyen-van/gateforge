package services

import (
	"context"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/stretchr/testify/require"
)

func TestSessionService_BrowserSessionTTL(t *testing.T) {
	cfg := testConfig()
	svc := ProvideSessionService(cfg, newSessionTestRepo(), &auditCapture{})

	require.Equal(t, cfg.SSOSessionTTL, svc.BrowserSessionTTL(false))
	require.Equal(t, cfg.SSOSessionRememberTTL, svc.BrowserSessionTTL(true))

	cfg.SSOSessionTTL = 0
	cfg.SSOSessionRememberTTL = 0
	svc2 := ProvideSessionService(cfg, newSessionTestRepo(), &auditCapture{})
	require.Equal(t, cfg.JWTRefreshTTL, svc2.BrowserSessionTTL(false))
	require.Equal(t, 720*time.Hour, svc2.BrowserSessionTTL(true))
}

func TestSessionService_CreateGetInvalidate(t *testing.T) {
	repo := newSessionTestRepo()
	audit := &auditCapture{}
	svc := ProvideSessionService(testConfig(), repo, audit)

	id, ttl, err := svc.Create(context.Background(), "user-1", "tenant-1", "127.0.0.1", "agent", false)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	require.Greater(t, ttl, time.Duration(0))
	require.Equal(t, constants.AuditActionSessionCreate, audit.params[0].Action)

	sess, err := svc.GetSession(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, "user-1", sess.UserID)

	userID, err := svc.GetUserID(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, "user-1", userID)

	_, err = svc.GetSession(context.Background(), "")
	require.Error(t, err)

	_, err = svc.GetUserID(context.Background(), "missing")
	require.Error(t, err)

	require.NoError(t, svc.Invalidate(context.Background(), id))
	require.NoError(t, svc.Invalidate(context.Background(), ""))
	_, err = svc.GetSession(context.Background(), id)
	require.Error(t, err)
}

func TestSessionService_InvalidateAllForUser(t *testing.T) {
	repo := newSessionTestRepo()
	audit := &auditCapture{}
	svc := ProvideSessionService(testConfig(), repo, audit)

	_, _, err := svc.Create(context.Background(), "user-1", "tenant-1", "", "", true)
	require.NoError(t, err)
	_, _, err = svc.Create(context.Background(), "user-1", "tenant-2", "", "", false)
	require.NoError(t, err)
	_, _, err = svc.Create(context.Background(), "user-2", "tenant-1", "", "", false)
	require.NoError(t, err)

	require.NoError(t, svc.InvalidateAllForUser(context.Background(), ""))
	require.NoError(t, svc.InvalidateAllForUser(context.Background(), "user-1"))
	require.Equal(t, constants.AuditActionSessionRevokeAll, audit.params[len(audit.params)-1].Action)

	_, err = svc.GetSession(context.Background(), "")
	require.Error(t, err)
}

func TestSessionService_CreateRepoError(t *testing.T) {
	repo := &sessionTestRepo{sessions: map[string]*models.Session{}}
	repo.sessions = nil // force panic? no - use broken repo
	broken := &brokenSessionRepo{}
	svc := ProvideSessionService(testConfig(), broken, &auditCapture{})
	_, _, err := svc.Create(context.Background(), "u", "t", "", "", false)
	require.Error(t, err)
}

type brokenSessionRepo struct{}

func (brokenSessionRepo) Create(context.Context, *models.Session) error {
	return context.Canceled
}
func (brokenSessionRepo) GetValidByID(context.Context, string) (*models.Session, error) {
	return nil, context.Canceled
}
func (brokenSessionRepo) DeleteByID(context.Context, string) error { return context.Canceled }
func (brokenSessionRepo) DeleteAllByUserID(context.Context, string) error {
	return context.Canceled
}
func (brokenSessionRepo) CountActive(context.Context) (int64, error) { return 0, nil }
func (brokenSessionRepo) CountActiveByUserID(context.Context, string) (int64, error) {
	return 0, nil
}
