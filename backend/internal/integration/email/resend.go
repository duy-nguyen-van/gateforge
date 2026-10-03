package email

import (
	"context"
	"fmt"
	"strings"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"

	"github.com/resend/resend-go/v3"
)

type ResendSender struct {
	client *resend.Client
	from   string
}

func NewResendSender(cfg config.Config) (*ResendSender, error) {
	if strings.TrimSpace(cfg.ResendAPIKey) == "" {
		return nil, errors.InternalError("Resend API key is required", fmt.Errorf("RESEND_API_KEY is empty")).
			WithOperation("initialize_resend").
			WithResource("resend")
	}
	from := cfg.EmailFrom
	if from == "" {
		from = "GateForge <onboarding@resend.dev>"
	}
	return &ResendSender{client: resend.NewClient(cfg.ResendAPIKey), from: from}, nil
}

func (s *ResendSender) SendEmail(ctx context.Context, request EmailRequest) (*EmailResponse, error) {
	params := &resend.SendEmailRequest{
		From:    s.from,
		To:      request.To,
		Subject: request.Subject,
		Html:    request.HTMLBody,
		Text:    request.TextBody,
	}
	var sent *resend.SendEmailResponse
	var err error
	if request.IdempotencyKey != "" {
		sent, err = s.client.Emails.SendWithOptions(ctx, params, &resend.SendEmailOptions{IdempotencyKey: request.IdempotencyKey})
	} else {
		sent, err = s.client.Emails.SendWithContext(ctx, params)
	}
	if err != nil {
		return nil, errors.ExternalServiceError("Failed to send email via Resend", err).
			WithOperation("send_email").
			WithResource("resend")
	}
	id := ""
	if sent != nil {
		id = sent.Id
	}
	return &EmailResponse{MessageID: id, Provider: constants.EmailProviderResend, Status: "sent"}, nil
}

func (s *ResendSender) SendRawEmail(ctx context.Context, rawData []byte) (*EmailResponse, error) {
	return nil, errors.InternalError("Resend does not support raw MIME send", nil).
		WithOperation("send_raw_email").
		WithResource("resend")
}
