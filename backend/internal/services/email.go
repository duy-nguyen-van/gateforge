package services

import (
	"context"
	"fmt"

	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/integration/email"
	"github.com/gateforge-iam/gateforge-iam/internal/logger"

	"go.uber.org/zap"
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
	message := &email.EmailRequest{
		To:       []string{userEmail},
		Subject:  "Welcome to My Echo App!",
		TextBody: fmt.Sprintf("Hello %s,\n\nWelcome to My Echo App! We're excited to have you on board.\n\nBest regards,\nThe Team", userName),
		HTMLBody: fmt.Sprintf(`
			<html>
				<body>
					<h2>Welcome to My Echo App!</h2>
					<p>Hello %s,</p>
					<p>Welcome to My Echo App! We're excited to have you on board.</p>
					<p>Best regards,<br>The Team</p>
				</body>
			</html>
		`, userName),
	}

	_, err := s.emailSender.SendEmail(ctx, *message)

	if err != nil {
		logger.Log.Error("email welcome send failed",
			zap.String("operation", "send_welcome_email"),
			zap.String("user_email", userEmail),
			zap.Error(err))
		return errors.ExternalServiceError("Failed to send welcome email", err).
			WithOperation("send_welcome_email").
			WithResource("email").
			WithContext("user_email", userEmail)
	}

	logger.Log.Info("email welcome sent",
		zap.String("operation", "send_welcome_email"),
		zap.String("user_email", userEmail))
	return nil
}

// SendPasswordResetEmail sends a password reset email
func (s *EmailService) SendPasswordResetEmail(ctx context.Context, userEmail, resetToken string) error {
	resetURL := fmt.Sprintf("https://yourapp.com/reset-password?token=%s", resetToken)

	message := &email.EmailRequest{
		To:       []string{userEmail},
		Subject:  "Password Reset Request",
		TextBody: fmt.Sprintf("You requested a password reset. Click the link below to reset your password:\n\n%s\n\nIf you didn't request this, please ignore this email.", resetURL),
		HTMLBody: fmt.Sprintf(`
			<html>
				<body>
					<h2>Password Reset Request</h2>
					<p>You requested a password reset. Click the link below to reset your password:</p>
					<p><a href="%s">Reset Password</a></p>
					<p>If you didn't request this, please ignore this email.</p>
				</body>
			</html>
		`, resetURL),
	}

	_, err := s.emailSender.SendEmail(ctx, *message)

	if err != nil {
		logger.Log.Error("email password reset send failed",
			zap.String("operation", "send_password_reset_email"),
			zap.String("user_email", userEmail),
			zap.Error(err))
		return errors.ExternalServiceError("Failed to send password reset email", err).
			WithOperation("send_password_reset_email").
			WithResource("email").
			WithContext("user_email", userEmail)
	}

	logger.Log.Info("email password reset sent",
		zap.String("operation", "send_password_reset_email"),
		zap.String("user_email", userEmail))
	return nil
}

// SendNotificationEmail sends a notification email
func (s *EmailService) SendNotificationEmail(ctx context.Context, userEmail, subject, message string) error {
	emailMessage := &email.EmailRequest{
		To:       []string{userEmail},
		Subject:  subject,
		TextBody: message,
		HTMLBody: fmt.Sprintf(`
			<html>
				<body>
					<h2>%s</h2>
					<p>%s</p>
				</body>
			</html>
		`, subject, message),
	}

	_, err := s.emailSender.SendEmail(ctx, *emailMessage)

	if err != nil {
		logger.Log.Error("email notification send failed",
			zap.String("operation", "send_notification_email"),
			zap.String("user_email", userEmail),
			zap.String("subject", subject),
			zap.Error(err))
		return errors.ExternalServiceError("Failed to send notification email", err).
			WithOperation("send_notification_email").
			WithResource("email").
			WithContext("user_email", userEmail)
	}

	logger.Log.Info("email notification sent",
		zap.String("operation", "send_notification_email"),
		zap.String("user_email", userEmail),
		zap.String("subject", subject))
	return nil
}
