package dtos

import (
	"encoding/json"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/models"
)

// AdminAuditLogResponse is one audit log row for the admin console.
type AdminAuditLogResponse struct {
	ID            string         `json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	TenantID      *string        `json:"tenant_id,omitempty"`
	Action        string         `json:"action"`
	Result        string         `json:"result"`
	ActorType     string         `json:"actor_type"`
	ActorID       *string        `json:"actor_id,omitempty"`
	ResourceType  *string        `json:"resource_type,omitempty"`
	ResourceID    *string        `json:"resource_id,omitempty"`
	ResourceName  *string        `json:"resource_name,omitempty"`
	IPAddress     *string        `json:"ip_address,omitempty"`
	UserAgent     *string        `json:"user_agent,omitempty"`
	RequestID     *string        `json:"request_id,omitempty"`
	CorrelationID *string        `json:"correlation_id,omitempty"`
	OldValue      map[string]any `json:"old_value,omitempty"`
	NewValue      map[string]any `json:"new_value,omitempty"`
}

// NewAdminAuditLogResponse maps a persistence model to an API response.
func NewAdminAuditLogResponse(row *models.AuditLog) *AdminAuditLogResponse {
	if row == nil {
		return nil
	}
	resp := &AdminAuditLogResponse{
		ID:            row.ID,
		CreatedAt:     row.CreatedAt,
		TenantID:      row.TenantID,
		Action:        row.Action,
		Result:        row.Result,
		ActorType:     row.ActorType,
		ActorID:       row.ActorID,
		ResourceType:  row.ResourceType,
		ResourceID:    row.ResourceID,
		ResourceName:  row.ResourceName,
		IPAddress:     row.IPAddress,
		UserAgent:     row.UserAgent,
		RequestID:     row.RequestID,
		CorrelationID: row.CorrelationID,
	}
	if len(row.OldValue) > 0 {
		var old map[string]any
		if err := json.Unmarshal(row.OldValue, &old); err == nil {
			resp.OldValue = old
		}
	}
	if len(row.NewValue) > 0 {
		var nv map[string]any
		if err := json.Unmarshal(row.NewValue, &nv); err == nil {
			resp.NewValue = nv
		}
	}
	return resp
}

// AdminAuditLogListParams filters audit log listing.
type AdminAuditLogListParams struct {
	TenantID string
	Action   string
	Result   string
	ActorID  string
	From     *time.Time
	To       *time.Time
}

// AdminLoginHistoryListParams filters login history listing.
type AdminLoginHistoryListParams struct {
	TenantID string
	ActorID  string
	Result   string
	From     *time.Time
	To       *time.Time
}
