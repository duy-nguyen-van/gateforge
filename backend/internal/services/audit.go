package services

import (
	"context"
	"encoding/json"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/logger"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/monitoring"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"
	"github.com/gateforge-iam/gateforge-iam/internal/request"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/datatypes"
)

// AuditService persists security and admin audit events.
type AuditService interface {
	Record(ctx context.Context, params domains.AuditRecordParams)
	RecordRequired(ctx context.Context, params domains.AuditRecordParams) error
	Shutdown(ctx context.Context) error
}

type auditService struct {
	repo repositories.AuditLogRepository
	ch   chan *models.AuditLog
	stop chan struct{}
	done chan struct{}
}

// ProvideAuditService wires audit log recording.
func ProvideAuditService(repo repositories.AuditLogRepository, cfg *config.Config, lc fx.Lifecycle) AuditService {
	buffer := 1024
	if cfg != nil && cfg.AuditAsyncBuffer > 0 {
		buffer = cfg.AuditAsyncBuffer
	}
	s := &auditService{
		repo: repo,
		ch:   make(chan *models.AuditLog, buffer),
		stop: make(chan struct{}),
		done: make(chan struct{}),
	}
	go s.loop()
	if lc != nil {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return s.Shutdown(ctx)
			},
		})
	}
	return s
}

func (s *auditService) Shutdown(ctx context.Context) error {
	select {
	case <-s.stop:
		return nil
	default:
		close(s.stop)
	}
	select {
	case <-s.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (s *auditService) loop() {
	defer close(s.done)
	for {
		select {
		case <-s.stop:
			for {
				select {
				case row := <-s.ch:
					_ = s.repo.Create(context.Background(), row)
				default:
					return
				}
			}
		case row := <-s.ch:
			_ = s.repo.Create(context.Background(), row)
		}
	}
}

func (s *auditService) Record(ctx context.Context, params domains.AuditRecordParams) {
	row := s.build(ctx, params)
	select {
	case s.ch <- row:
	default:
		monitoring.AddAuditDropped(ctx)
		logger.From(ctx).Warn("audit log dropped", zap.String("action", params.Action))
	}
}

func (s *auditService) RecordRequired(ctx context.Context, params domains.AuditRecordParams) error {
	row := s.build(ctx, params)
	if err := s.repo.Create(ctx, row); err != nil {
		return errors.InternalError("Failed to write audit log", err).
			WithOperation("audit_record").
			WithResource("audit_log")
	}
	return nil
}

func (s *auditService) build(ctx context.Context, params domains.AuditRecordParams) *models.AuditLog {
	log := &models.AuditLog{
		BaseModel: models.NewBaseModel(),
		Action:    params.Action,
		Result:    string(params.Result),
		ActorType: string(params.ActorType),
	}

	if ac, ok := request.AuditContextFromContext(ctx); ok {
		if params.ActorID == "" {
			params.ActorID = ac.ActorID
		}
		if params.ActorType == "" && ac.ActorType != "" {
			params.ActorType = constants.AuditActorType(ac.ActorType)
		}
		if params.TenantID == "" {
			params.TenantID = ac.TenantID
		}
		if ac.IPAddress != "" {
			ip := ac.IPAddress
			log.IPAddress = &ip
		}
		if ac.UserAgent != "" {
			ua := ac.UserAgent
			log.UserAgent = &ua
		}
		if ac.RequestID != "" {
			rid := ac.RequestID
			log.RequestID = &rid
		}
		if ac.CorrelationID != "" {
			cid := ac.CorrelationID
			log.CorrelationID = &cid
		}
	}

	if params.ActorType != "" {
		log.ActorType = string(params.ActorType)
	}
	if params.ActorID != "" {
		aid := params.ActorID
		log.ActorID = &aid
	}
	if params.TenantID != "" {
		tid := params.TenantID
		log.TenantID = &tid
	}
	if params.ResourceType != "" {
		rt := string(params.ResourceType)
		log.ResourceType = &rt
	}
	if params.ResourceID != "" {
		rid := params.ResourceID
		log.ResourceID = &rid
	}
	if params.ResourceName != "" {
		rn := params.ResourceName
		log.ResourceName = &rn
	}
	if params.OldValue != nil {
		if b, err := json.Marshal(params.OldValue); err == nil {
			log.OldValue = datatypes.JSON(b)
		}
	}
	if params.NewValue != nil {
		if b, err := json.Marshal(params.NewValue); err == nil {
			log.NewValue = datatypes.JSON(b)
		}
	}
	return log
}
