package models

import (
	"gorm.io/datatypes"
)

// AuditLog is an append-only security and admin action record.
type AuditLog struct {
	BaseModel
	TenantID      *string        `gorm:"column:tenant_id;type:uuid;index:idx_audit_logs_tenant_created_at,priority:1"`
	Action        string         `gorm:"column:action;type:varchar(100);not null;index:idx_audit_logs_action_created_at,priority:1"`
	Result        string         `gorm:"column:result;type:varchar(20);not null"`
	ActorType     string         `gorm:"column:actor_type;type:varchar(50);not null"`
	ActorID       *string        `gorm:"column:actor_id;type:text;index:idx_audit_logs_actor_created_at,priority:1"`
	ResourceType  *string        `gorm:"column:resource_type;type:varchar(50);index:idx_audit_logs_resource,priority:1"`
	ResourceID    *string        `gorm:"column:resource_id;type:uuid;index:idx_audit_logs_resource,priority:2"`
	ResourceName  *string        `gorm:"column:resource_name;type:varchar(255)"`
	IPAddress     *string        `gorm:"column:ip_address;type:inet"`
	UserAgent     *string        `gorm:"column:user_agent;type:text"`
	RequestID     *string        `gorm:"column:request_id;type:text"`
	CorrelationID *string        `gorm:"column:correlation_id;type:text"`
	OldValue      datatypes.JSON `gorm:"column:old_value;type:jsonb"`
	NewValue      datatypes.JSON `gorm:"column:new_value;type:jsonb"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
