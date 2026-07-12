package request

import "context"

type ctxKeyAudit struct{}

// AuditContext carries HTTP and actor metadata for audit log recording.
type AuditContext struct {
	IPAddress     string
	UserAgent     string
	RequestID     string
	CorrelationID string
	ActorType     string
	ActorID       string
	TenantID      string
}

// AuditContextFromContext retrieves audit metadata from the context.
func AuditContextFromContext(ctx context.Context) (AuditContext, bool) {
	val, ok := ctx.Value(ctxKeyAudit{}).(AuditContext)
	return val, ok
}

// NewAuditContextContext stores audit metadata on the context.
func NewAuditContextContext(ctx context.Context, ac AuditContext) context.Context {
	return context.WithValue(ctx, ctxKeyAudit{}, ac)
}
