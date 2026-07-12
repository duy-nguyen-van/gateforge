package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/crypto"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/logger"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
)

// MFAService handles TOTP enrollment, recovery codes, and MFA-after-login verification.
type MFAService interface {
	HasActiveMFA(ctx context.Context, userID string) (bool, error)
	CreateLoginTicket(ctx context.Context, p auth.MFAPendingPayload) (ticket string, expiresSec int64, err error)
	SetupTOTP(ctx context.Context, userID, accountEmail string) (secret string, otpauthURI string, err error)
	VerifyTOTPEnrollment(ctx context.Context, userID, code string) error
	RegenerateRecoveryCodes(ctx context.Context, userID string) (plainCodes []string, err error)
	VerifyLoginChallenge(ctx context.Context, ticket, code string) (*auth.MFAPendingPayload, error)
}

type mfaService struct {
	cfg          *config.Config
	totpRepo     repositories.UserMFATOTPRepository
	recoveryRepo repositories.UserMFARecoveryCodeRepository
	ephemeral    *auth.EphemeralStore
	audit        AuditService
}

func ProvideMFAService(
	cfg *config.Config,
	totpRepo repositories.UserMFATOTPRepository,
	recoveryRepo repositories.UserMFARecoveryCodeRepository,
	ephemeral *auth.EphemeralStore,
	audit AuditService,
) MFAService {
	return &mfaService{
		cfg:          cfg,
		totpRepo:     totpRepo,
		recoveryRepo: recoveryRepo,
		ephemeral:    ephemeral,
		audit:        audit,
	}
}

func (s *mfaService) mfaKey() string {
	if s.cfg.MFAEncryptionKey != "" {
		return s.cfg.MFAEncryptionKey
	}
	return s.cfg.JWTSecret
}

func (s *mfaService) HasActiveMFA(ctx context.Context, userID string) (bool, error) {
	row, err := s.totpRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return false, err
	}
	return row != nil, nil
}

func (s *mfaService) CreateLoginTicket(ctx context.Context, p auth.MFAPendingPayload) (string, int64, error) {
	ticket, err := s.ephemeral.PutMFAPending(ctx, p)
	if err != nil {
		return "", 0, err
	}
	return ticket, int64(s.cfg.MFAPendingTicketTTL.Seconds()), nil
}

func (s *mfaService) SetupTOTP(ctx context.Context, userID, accountEmail string) (secret string, otpauthURI string, err error) {
	active, err := s.HasActiveMFA(ctx, userID)
	if err != nil {
		return "", "", err
	}
	if active {
		return "", "", errors.ValidationError("TOTP is already enabled", nil).
			WithOperation("mfa_totp_setup").
			WithResource("user_mfa_totp")
	}
	issuer := s.cfg.AppName
	if issuer == "" {
		issuer = s.cfg.WebauthnRPDisplayName
	}
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: accountEmail,
		Period:      30,
		SecretSize:  20,
	})
	if err != nil {
		return "", "", errors.InternalError("Failed to generate TOTP secret", err).
			WithOperation("mfa_totp_setup").
			WithResource("user_mfa_totp")
	}
	enc, err := crypto.EncryptMFASecret(s.mfaKey(), key.Secret())
	if err != nil {
		return "", "", errors.InternalError("Failed to encrypt TOTP secret", err).
			WithOperation("mfa_totp_setup").
			WithResource("user_mfa_totp")
	}
	row := &models.UserMFATOTP{
		BaseModel:       models.NewBaseModel(),
		UserID:          userID,
		SecretEncrypted: enc,
		Enabled:         false,
	}
	if err := s.totpRepo.UpsertPending(ctx, row); err != nil {
		return "", "", err
	}
	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionMFATOTPSetup,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      userID,
		ResourceType: constants.AuditResourceTypeUser,
		ResourceID:   userID,
	})
	return key.Secret(), key.String(), nil
}

func (s *mfaService) VerifyTOTPEnrollment(ctx context.Context, userID, code string) error {
	row, err := s.totpRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if row.Enabled {
		return errors.ValidationError("TOTP already verified", nil).
			WithOperation("mfa_totp_verify").
			WithResource("user_mfa_totp")
	}
	plain, err := crypto.DecryptMFASecret(s.mfaKey(), row.SecretEncrypted)
	if err != nil {
		return errors.InternalError("Failed to decrypt TOTP secret", err).
			WithOperation("mfa_totp_verify").
			WithResource("user_mfa_totp")
	}
	if !totp.Validate(code, plain) {
		return errors.ValidationError("Invalid verification code", nil).
			WithOperation("mfa_totp_verify").
			WithResource("user_mfa_totp")
	}
	if err := s.totpRepo.MarkVerifiedAndEnabled(ctx, userID); err != nil {
		return err
	}
	logger.Log.Info("mfa totp enabled",
		zap.String("operation", "mfa_totp_verify"),
		zap.String("user_id", userID))
	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionMFATOTPEnable,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      userID,
		ResourceType: constants.AuditResourceTypeUser,
		ResourceID:   userID,
	})
	return nil
}

func (s *mfaService) RegenerateRecoveryCodes(ctx context.Context, userID string) ([]string, error) {
	active, err := s.HasActiveMFA(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !active {
		return nil, errors.ForbiddenError("Enable TOTP before generating recovery codes", nil).
			WithOperation("mfa_recovery_codes").
			WithResource("user_mfa_recovery_code")
	}
	n := s.cfg.MFARecoveryCodeCount
	if n <= 0 {
		n = 10
	}
	plain := make([]string, 0, n)
	rows := make([]*models.UserMFARecoveryCode, 0, n)
	for i := 0; i < n; i++ {
		code, err := randomRecoveryCode()
		if err != nil {
			return nil, err
		}
		plain = append(plain, code)
		h := auth.HashOpaqueToken(code)
		rows = append(rows, &models.UserMFARecoveryCode{
			BaseModel: models.NewBaseModel(),
			UserID:    userID,
			CodeHash:  h,
		})
	}
	if err := s.recoveryRepo.ReplaceAllForUser(ctx, userID, rows); err != nil {
		return nil, err
	}
	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionMFARecoveryRegenerate,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      userID,
		ResourceType: constants.AuditResourceTypeUser,
		ResourceID:   userID,
	})
	return plain, nil
}

func (s *mfaService) VerifyLoginChallenge(ctx context.Context, ticket, code string) (*auth.MFAPendingPayload, error) {
	payload, err := s.ephemeral.TakeMFAPending(ctx, strings.TrimSpace(ticket))
	if err != nil {
		return nil, err
	}
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, errors.ValidationError("Code is required", nil).
			WithOperation("mfa_challenge_verify").
			WithResource("mfa")
	}
	if len(code) == 6 {
		if err := s.validateTOTPForLogin(ctx, payload.UserID, code); err != nil {
			s.audit.Record(ctx, domains.AuditRecordParams{
				Action:    constants.AuditActionMFAChallengeVerify,
				Result:    constants.AuditResultFailure,
				ActorType: constants.AuditActorTypeUser,
				ActorID:   payload.UserID,
				TenantID:  payload.TenantID,
			})
			return nil, err
		}
		s.audit.Record(ctx, domains.AuditRecordParams{
			Action:    constants.AuditActionMFAChallengeVerify,
			Result:    constants.AuditResultSuccess,
			ActorType: constants.AuditActorTypeUser,
			ActorID:   payload.UserID,
			TenantID:  payload.TenantID,
		})
		return payload, nil
	}
	if err := s.tryRecoveryCode(ctx, payload.UserID, code); err != nil {
		s.audit.Record(ctx, domains.AuditRecordParams{
			Action:    constants.AuditActionMFAChallengeVerify,
			Result:    constants.AuditResultFailure,
			ActorType: constants.AuditActorTypeUser,
			ActorID:   payload.UserID,
			TenantID:  payload.TenantID,
		})
		return nil, err
	}
	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:    constants.AuditActionMFAChallengeVerify,
		Result:    constants.AuditResultSuccess,
		ActorType: constants.AuditActorTypeUser,
		ActorID:   payload.UserID,
		TenantID:  payload.TenantID,
		NewValue:  map[string]any{"method": "recovery_code"},
	})
	return payload, nil
}

func (s *mfaService) validateTOTPForLogin(ctx context.Context, userID, code string) error {
	row, err := s.totpRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if row == nil {
		return errors.UnauthorizedError("Invalid verification code", nil).
			WithOperation("mfa_challenge_verify").
			WithResource("mfa")
	}
	plain, err := crypto.DecryptMFASecret(s.mfaKey(), row.SecretEncrypted)
	if err != nil {
		return errors.InternalError("Failed to decrypt TOTP secret", err).
			WithOperation("mfa_challenge_verify").
			WithResource("user_mfa_totp")
	}
	if !totp.Validate(code, plain) {
		return errors.UnauthorizedError("Invalid verification code", nil).
			WithOperation("mfa_challenge_verify").
			WithResource("mfa")
	}
	return nil
}

func (s *mfaService) tryRecoveryCode(ctx context.Context, userID, code string) error {
	hash := auth.HashOpaqueToken(code)
	rows, err := s.recoveryRepo.FindUnusedByUserID(ctx, userID)
	if err != nil {
		return err
	}
	for _, r := range rows {
		if r.CodeHash == hash {
			if err := s.recoveryRepo.MarkUsed(ctx, r.ID); err != nil {
				return err
			}
			return nil
		}
	}
	return errors.UnauthorizedError("Invalid verification code", nil).
		WithOperation("mfa_challenge_verify").
		WithResource("mfa")
}

func randomRecoveryCode() (string, error) {
	b := make([]byte, 5)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return strings.ToUpper(hex.EncodeToString(b)), nil
}
