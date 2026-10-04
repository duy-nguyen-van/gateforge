package services

import (
	"context"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/monitoring"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"go.uber.org/fx"
)

// RetentionService deletes expired authorization codes, sessions, and stale refresh tokens.
type RetentionService interface {
	PurgeOnce(ctx context.Context) error
}

type retentionService struct {
	cfg  *config.Config
	repo repositories.RetentionRepository
}

func ProvideRetentionService(cfg *config.Config, repo repositories.RetentionRepository, lc fx.Lifecycle) RetentionService {
	svc := &retentionService{cfg: cfg, repo: repo}
	if lc == nil || cfg == nil {
		return svc
	}
	ctx, cancel := context.WithCancel(context.Background())
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go svc.loop(ctx)
			return nil
		},
		OnStop: func(context.Context) error {
			cancel()
			return nil
		},
	})
	return svc
}

func (s *retentionService) loop(ctx context.Context) {
	interval := time.Hour
	if s.cfg != nil && s.cfg.RetentionInterval > 0 {
		interval = s.cfg.RetentionInterval
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	_ = s.PurgeOnce(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = s.PurgeOnce(ctx)
		}
	}
}

func (s *retentionService) PurgeOnce(ctx context.Context) error {
	if s.repo == nil {
		return nil
	}
	days := 7
	if s.cfg != nil && s.cfg.RetentionRefreshDays > 0 {
		days = s.cfg.RetentionRefreshDays
	}
	before := time.Now().UTC().Add(-time.Duration(days) * 24 * time.Hour)
	locked, counts, err := s.repo.Purge(ctx, before, 1000)
	if err != nil || !locked {
		return err
	}
	monitoring.AddRetentionRows(ctx, "authorization_codes", counts.AuthorizationCodes)
	monitoring.AddRetentionRows(ctx, "sessions", counts.Sessions)
	monitoring.AddRetentionRows(ctx, "refresh_tokens", counts.RefreshTokens)
	return nil
}
