package services

import (
	"context"

	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/integration/email"
	"github.com/gateforge-iam/gateforge-iam/internal/monitoring"
)

// EmailService handles email business logic
type EmailService struct {
	emailSender email.EmailSender
}

// NewEmailService creates a new email service
func ProvideEmailService(emailSender email.EmailSender) EmailService {
	return EmailService{
		emailSender: emailSender,
	}
}

// SendWelcomeEmail sends a welcome email to a new user
func (s *EmailService) SendWelcomeEmail(ctx context.Context, userEmail, userName string) error {
	return monitoring.ObserveErr(ctx, emailTracer, "EmailService.SendWelcomeEmail", nil, func(ctx context.Context) error {
		return s.sendTemplate(ctx, "send_welcome_email", userEmail, email.TemplateWelcome, email.Vars{
			"UserName": userName,
		})
	})
}

// SendPasswordResetEmail sends a password reset email. resetURL is the full link.
func (s EmailService) SendPasswordResetEmail(ctx context.Context, userEmail, resetURL string) error {
	return monitoring.ObserveErr(ctx, emailTracer, "EmailService.SendPasswordResetEmail", nil, func(ctx context.Context) error {
		return s.sendTemplate(ctx, "send_password_reset_email", userEmail, email.TemplatePasswordReset, email.Vars{
			"CTAURL": resetURL,
		})
	})
}

// SendNotificationEmail sends a notification email
func (s *EmailService) SendNotificationEmail(ctx context.Context, userEmail, subject, message string) error {
	return monitoring.ObserveErr(ctx, emailTracer, "EmailService.SendNotificationEmail", nil, func(ctx context.Context) error {
		return s.sendTemplate(ctx, "send_notification_email", userEmail, email.TemplateNotification, email.Vars{
			"Subject": subject,
			"Message": message,
		})
	})
}

func (s EmailService) sendTemplate(ctx context.Context, operation, to, name string, vars email.Vars) error {
	msg, err := email.Render(name, vars)
	if err != nil {
		return errors.InternalError("Failed to render email template", err).
			WithOperation(operation).
			WithResource("email").
			WithContext("template", name)
	}
	_, err = s.emailSender.SendEmail(ctx, email.EmailRequest{
		To:           []string{to},
		Subject:      msg.Subject,
		TextBody:     msg.TextBody,
		HTMLBody:     msg.HTMLBody,
		TemplateID:   name,
		TemplateData: vars.Map(),
	})
	if err != nil {
		return errors.ExternalServiceError("Failed to send email", err).
			WithOperation(operation).
			WithResource("email").
			WithContext("template", name)
	}
	return nil
}
