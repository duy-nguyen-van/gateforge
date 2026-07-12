package services

import (
	"context"
	"encoding/json"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/logger"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"
	"github.com/gateforge-iam/gateforge-iam/internal/request"

	"go.uber.org/zap"
	"gorm.io/datatypes"
)

// AuditService persists security and admin audit events.
type AuditService interface {
	Record(ctx context.Context, params domains.AuditRecordParams)
}

type auditService struct {
	repo repositories.AuditLogRepository
}

// ProvideAuditService wires audit log recording.
func ProvideAuditService(repo repositories.AuditLogRepository) AuditService {
	return &auditService{repo: repo}
}

func (s *auditService) Record(ctx context.Context, params domains.AuditRecordParams) {
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

	if err := s.repo.Create(ctx, log); err != nil {
		logger.Log.Warn("audit log persist failed",
			zap.String("action", params.Action),
			zap.String("result", string(params.Result)),
			zap.Error(err),
		)
	}
}
