package dtos

import (
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func TestNewAdminAuditLogResponse_Nil(t *testing.T) {
	require.Nil(t, NewAdminAuditLogResponse(nil))
}

func TestNewAdminAuditLogResponse(t *testing.T) {
	createdAt := time.Date(2024, 4, 12, 11, 0, 0, 0, time.UTC)
	logID := uuid.Must(uuid.NewV7()).String()
	tenantID := uuid.Must(uuid.NewV7()).String()
	actorID := "user-123"
	resourceType := "client"
	resourceID := uuid.Must(uuid.NewV7()).String()
	resourceName := "My Client"
	ip := "192.168.1.1"
	ua := "Mozilla/5.0"
	requestID := "req-abc"
	correlationID := "corr-xyz"

	row := &models.AuditLog{
		BaseModel:     models.BaseModel{ID: logID, CreatedAt: createdAt},
		TenantID:      &tenantID,
		Action:        "client.create",
		Result:        "success",
		ActorType:     "user",
		ActorID:       &actorID,
		ResourceType:  &resourceType,
		ResourceID:    &resourceID,
		ResourceName:  &resourceName,
		IPAddress:     &ip,
		UserAgent:     &ua,
		RequestID:     &requestID,
		CorrelationID: &correlationID,
		OldValue:      datatypes.JSON(`{"name":"old"}`),
		NewValue:      datatypes.JSON(`{"name":"new"}`),
	}

	resp := NewAdminAuditLogResponse(row)
	require.Equal(t, logID, resp.ID)
	require.Equal(t, createdAt, resp.CreatedAt)
	require.Equal(t, &tenantID, resp.TenantID)
	require.Equal(t, "client.create", resp.Action)
	require.Equal(t, "success", resp.Result)
	require.Equal(t, "user", resp.ActorType)
	require.Equal(t, &actorID, resp.ActorID)
	require.Equal(t, &resourceType, resp.ResourceType)
	require.Equal(t, &resourceID, resp.ResourceID)
	require.Equal(t, &resourceName, resp.ResourceName)
	require.Equal(t, &ip, resp.IPAddress)
	require.Equal(t, &ua, resp.UserAgent)
	require.Equal(t, &requestID, resp.RequestID)
	require.Equal(t, &correlationID, resp.CorrelationID)
	require.Equal(t, map[string]any{"name": "old"}, resp.OldValue)
	require.Equal(t, map[string]any{"name": "new"}, resp.NewValue)
}

func TestNewAdminAuditLogResponse_InvalidJSONValues(t *testing.T) {
	row := &models.AuditLog{
		BaseModel: models.BaseModel{ID: uuid.Must(uuid.NewV7()).String()},
		Action:    "test.action",
		Result:    "failure",
		ActorType: "system",
		OldValue:  datatypes.JSON(`not-json`),
		NewValue:  datatypes.JSON(`{invalid`),
	}

	resp := NewAdminAuditLogResponse(row)
	require.NotNil(t, resp)
	require.Nil(t, resp.OldValue)
	require.Nil(t, resp.NewValue)
}

func TestNewAdminAuditLogResponse_EmptyJSONValues(t *testing.T) {
	row := &models.AuditLog{
		BaseModel: models.BaseModel{ID: uuid.Must(uuid.NewV7()).String()},
		Action:    "test.action",
		Result:    "success",
		ActorType: "user",
	}

	resp := NewAdminAuditLogResponse(row)
	require.NotNil(t, resp)
	require.Nil(t, resp.OldValue)
	require.Nil(t, resp.NewValue)
}
