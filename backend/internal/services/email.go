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
		return s.sendTemplate(ctx, "send_welcome_email", userEmail, email.TemplateWelcome, "", email.Vars{
			"UserName": userName,
		})
	})
}

// SendPasswordResetEmail sends a password reset email. resetURL is the full link.
func (s EmailService) SendPasswordResetEmail(ctx context.Context, userEmail, resetURL string) error {
	return monitoring.ObserveErr(ctx, emailTracer, "EmailService.SendPasswordResetEmail", nil, func(ctx context.Context) error {
		return s.sendTemplate(ctx, "send_password_reset_email", userEmail, email.TemplatePasswordReset, "", email.Vars{
			"CTAURL": resetURL,
		})
	})
}

// SendNotificationEmail sends a notification email
func (s *EmailService) SendNotificationEmail(ctx context.Context, userEmail, subject, message string) error {
	return monitoring.ObserveErr(ctx, emailTracer, "EmailService.SendNotificationEmail", nil, func(ctx context.Context) error {
		return s.sendTemplate(ctx, "send_notification_email", userEmail, email.TemplateNotification, "", email.Vars{
			"Subject": subject,
			"Message": message,
		})
	})
}

// SendMemberInviteEmail asks someone without an account to sign in and create one.
func (s EmailService) SendMemberInviteEmail(ctx context.Context, userEmail, orgName, role, signInURL, inviteID string) error {
	if s.emailSender == nil {
		return nil
	}
	return monitoring.ObserveErr(ctx, emailTracer, "EmailService.SendMemberInviteEmail", nil, func(ctx context.Context) error {
		return s.sendTemplate(ctx, "send_member_invite_email", userEmail, email.TemplateMemberInvite, "member-invite/"+inviteID, email.Vars{
			"OrgName": orgName,
			"Role":    role,
			"CTAURL":  signInURL,
		})
	})
}

// SendMemberAddedEmail tells an existing user they now belong to an organization.
// A nil sender is a no-op so tests that do not wire mail still succeed.
func (s EmailService) SendMemberAddedEmail(ctx context.Context, userEmail, orgName, role, signInURL, membershipID string) error {
	if s.emailSender == nil {
		return nil
	}
	return monitoring.ObserveErr(ctx, emailTracer, "EmailService.SendMemberAddedEmail", nil, func(ctx context.Context) error {
		return s.sendTemplate(ctx, "send_member_added_email", userEmail, email.TemplateMemberAdded, "member-added/"+membershipID, email.Vars{
			"OrgName": orgName,
			"Role":    role,
			"CTAURL":  signInURL,
		})
	})
}

func (s EmailService) sendTemplate(ctx context.Context, operation, to, name, idempotencyKey string, vars email.Vars) error {
	msg, err := email.Render(name, vars)
	if err != nil {
		return errors.InternalError("Failed to render email template", err).
			WithOperation(operation).
			WithResource("email").
			WithContext("template", name)
	}
	_, err = s.emailSender.SendEmail(ctx, email.EmailRequest{
		To:             []string{to},
		Subject:        msg.Subject,
		TextBody:       msg.TextBody,
		HTMLBody:       msg.HTMLBody,
		TemplateID:     name,
		TemplateData:   vars.Map(),
		IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		return errors.ExternalServiceError("Failed to send email", err).
			WithOperation(operation).
			WithResource("email").
			WithContext("template", name)
	}
	return nil
}
