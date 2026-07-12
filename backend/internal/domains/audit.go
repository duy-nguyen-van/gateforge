package domains

import "github.com/gateforge-iam/gateforge-iam/internal/constants"

// AuditRecordParams describes one audit event to persist.
type AuditRecordParams struct {
	Action       string
	Result       constants.AuditResult
	ActorType    constants.AuditActorType
	ActorID      string
	TenantID     string
	ResourceType constants.AuditResourceType
	ResourceID   string
	ResourceName string
	OldValue     any
	NewValue     any
}
