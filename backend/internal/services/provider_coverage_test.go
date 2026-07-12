package services

import (
	"context"
	"errors"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/integration/email"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"
	"github.com/gateforge-iam/gateforge-iam/internal/request"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/stretchr/testify/require"
)

type stubEmailSender struct {
	err error
}

func (s stubEmailSender) SendEmail(context.Context, email.EmailRequest) (*email.EmailResponse, error) {
	if s.err != nil {
		return &email.EmailResponse{Status: "failed", Provider: "stub"}, s.err
	}
	return &email.EmailResponse{MessageID: "msg-1", Status: "sent", Provider: "stub"}, nil
}

func (s stubEmailSender) SendRawEmail(context.Context, []byte) (*email.EmailResponse, error) {
	return s.SendEmail(context.Background(), email.EmailRequest{})
}

func TestProvideEmailService(t *testing.T) {
	svc := ProvideEmailService(stubEmailSender{})
	require.NotNil(t, svc.emailSender)
}

func TestEmailService_SendWelcomeEmail_Success(t *testing.T) {
	testutil.InitLogger()
	svc := ProvideEmailService(stubEmailSender{})
	require.NoError(t, svc.SendWelcomeEmail(context.Background(), "user@example.com", "Jane"))
}

func TestEmailService_SendWelcomeEmail_Error(t *testing.T) {
	testutil.InitLogger()
	svc := ProvideEmailService(stubEmailSender{err: errors.New("send failed")})
	err := svc.SendWelcomeEmail(context.Background(), "user@example.com", "Jane")
	require.Error(t, err)
}

func TestEmailService_SendPasswordResetEmail_Success(t *testing.T) {
	testutil.InitLogger()
	svc := ProvideEmailService(stubEmailSender{})
	require.NoError(t, svc.SendPasswordResetEmail(context.Background(), "user@example.com", "token-123"))
}

func TestEmailService_SendPasswordResetEmail_Error(t *testing.T) {
	testutil.InitLogger()
	svc := ProvideEmailService(stubEmailSender{err: errors.New("send failed")})
	require.Error(t, svc.SendPasswordResetEmail(context.Background(), "user@example.com", "token"))
}

func TestEmailService_SendNotificationEmail_Success(t *testing.T) {
	testutil.InitLogger()
	svc := ProvideEmailService(stubEmailSender{})
	require.NoError(t, svc.SendNotificationEmail(context.Background(), "user@example.com", "Subject", "Body"))
}

func TestEmailService_SendNotificationEmail_Error(t *testing.T) {
	testutil.InitLogger()
	svc := ProvideEmailService(stubEmailSender{err: errors.New("send failed")})
	require.Error(t, svc.SendNotificationEmail(context.Background(), "user@example.com", "Subject", "Body"))
}

type failingAuditRepo struct{}

func (failingAuditRepo) Create(context.Context, *models.AuditLog) error {
	return errors.New("db down")
}

func (failingAuditRepo) List(context.Context, repositories.AuditLogListFilters, *dtos.PageableRequest) (*dtos.DataResponse[models.AuditLog], error) {
	return nil, nil
}

func (failingAuditRepo) Count(context.Context, repositories.AuditLogListFilters) (int64, error) {
	return 0, nil
}

func TestProvideAuditService_RecordWithContext(t *testing.T) {
	testutil.InitLogger()
	svc := ProvideAuditService(failingAuditRepo{})

	ctx := request.NewAuditContextContext(context.Background(), request.AuditContext{
		ActorID:       "actor-1",
		ActorType:     string(constants.AuditActorTypeUser),
		TenantID:      "tenant-1",
		IPAddress:     "127.0.0.1",
		UserAgent:     "test-agent",
		RequestID:     "req-1",
		CorrelationID: "corr-1",
	})

	svc.Record(ctx, domains.AuditRecordParams{
		Action:       "login",
		Result:       constants.AuditResultSuccess,
		ResourceType: constants.AuditResourceTypeUser,
		ResourceID:   "user-1",
		ResourceName: "Jane",
		OldValue:     map[string]string{"status": "pending"},
		NewValue:     map[string]string{"status": "active"},
	})
}

func TestProvideAuditService_RecordCreateError(t *testing.T) {
	testutil.InitLogger()
	svc := ProvideAuditService(failingAuditRepo{})
	svc.Record(context.Background(), domains.AuditRecordParams{
		Action: "logout",
		Result: constants.AuditResultFailure,
	})
}
