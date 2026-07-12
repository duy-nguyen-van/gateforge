package services

import (
	"context"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/crypto"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/stretchr/testify/require"
)

func newMFATestService(t *testing.T) (MFAService, *mfaTOTPTestRepo, *mfaRecoveryTestRepo, *auth.EphemeralStore) {
	t.Helper()
	cfg := testConfig()
	totpRepo := newMFATOTPTestRepo()
	recoveryRepo := newMFARecoveryTestRepo()
	ephemeral := auth.NewEphemeralStore(newMemCache(), cfg)
	svc := ProvideMFAService(cfg, totpRepo, recoveryRepo, ephemeral, &auditCapture{})
	return svc, totpRepo, recoveryRepo, ephemeral
}

func TestMFAService_SetupAndVerifyTOTP(t *testing.T) {
	svc, totpRepo, _, _ := newMFATestService(t)
	userID := "user-mfa"

	secret, uri, err := svc.SetupTOTP(context.Background(), userID, "user@example.com")
	require.NoError(t, err)
	require.NotEmpty(t, secret)
	require.Contains(t, uri, "otpauth://")
	require.NotNil(t, totpRepo.byUser[userID])

	code, err := totp.GenerateCode(secret, time.Now())
	require.NoError(t, err)
	require.NoError(t, svc.VerifyTOTPEnrollment(context.Background(), userID, code))

	err = svc.VerifyTOTPEnrollment(context.Background(), userID, code)
	require.Error(t, err)

	secret2, _, err := svc.SetupTOTP(context.Background(), userID, "user@example.com")
	require.Error(t, err) // already enabled
	_ = secret2
}

func TestMFAService_HasActiveMFA(t *testing.T) {
	svc, totpRepo, _, _ := newMFATestService(t)
	active, err := svc.HasActiveMFA(context.Background(), "none")
	require.NoError(t, err)
	require.False(t, active)

	totpRepo.byUser["u1"] = &models.UserMFATOTP{UserID: "u1", Enabled: true}
	active, err = svc.HasActiveMFA(context.Background(), "u1")
	require.NoError(t, err)
	require.True(t, active)
}

func TestMFAService_RegenerateRecoveryCodes(t *testing.T) {
	svc, totpRepo, recoveryRepo, _ := newMFATestService(t)
	userID := "user-recovery"

	_, err := svc.RegenerateRecoveryCodes(context.Background(), userID)
	require.Error(t, err)

	cfg := testConfig()
	enc, err := crypto.EncryptMFASecret(cfg.MFAEncryptionKey, "JBSWY3DPEHPK3PXP")
	require.NoError(t, err)
	totpRepo.byUser[userID] = &models.UserMFATOTP{UserID: userID, SecretEncrypted: enc, Enabled: true}

	codes, err := svc.RegenerateRecoveryCodes(context.Background(), userID)
	require.NoError(t, err)
	require.Len(t, codes, cfg.MFARecoveryCodeCount)
	require.Len(t, recoveryRepo.byUser[userID], cfg.MFARecoveryCodeCount)
}

func TestMFAService_CreateLoginTicketAndVerify(t *testing.T) {
	svc, totpRepo, recoveryRepo, _ := newMFATestService(t)
	userID := "user-login-mfa"
	cfg := testConfig()

	secret := "JBSWY3DPEHPK3PXP"
	enc, err := crypto.EncryptMFASecret(cfg.MFAEncryptionKey, secret)
	require.NoError(t, err)
	totpRepo.byUser[userID] = &models.UserMFATOTP{UserID: userID, SecretEncrypted: enc, Enabled: true}

	ticket, exp, err := svc.CreateLoginTicket(context.Background(), auth.MFAPendingPayload{
		UserID:   userID,
		TenantID: "tenant-1",
	})
	require.NoError(t, err)
	require.NotEmpty(t, ticket)
	require.Greater(t, exp, int64(0))

	code, err := totp.GenerateCode(secret, time.Now())
	require.NoError(t, err)
	payload, err := svc.VerifyLoginChallenge(context.Background(), ticket, code)
	require.NoError(t, err)
	require.Equal(t, userID, payload.UserID)

	ticket2, _, err := svc.CreateLoginTicket(context.Background(), auth.MFAPendingPayload{UserID: userID, TenantID: "tenant-1"})
	require.NoError(t, err)
	plainCode := "ABCDE"
	hash := auth.HashOpaqueToken(plainCode)
	recoveryRepo.byUser[userID] = []models.UserMFARecoveryCode{{BaseModel: models.NewBaseModel(), UserID: userID, CodeHash: hash}}
	payload, err = svc.VerifyLoginChallenge(context.Background(), ticket2, plainCode)
	require.NoError(t, err)
	require.Equal(t, userID, payload.UserID)

	_, err = svc.VerifyLoginChallenge(context.Background(), "bad-ticket", code)
	require.Error(t, err)

	ticket3, _, err := svc.CreateLoginTicket(context.Background(), auth.MFAPendingPayload{UserID: userID, TenantID: "tenant-1"})
	require.NoError(t, err)
	_, err = svc.VerifyLoginChallenge(context.Background(), ticket3, "")
	require.Error(t, err)
}

func TestMFAService_VerifyLoginChallenge_InvalidCode(t *testing.T) {
	svc, totpRepo, _, _ := newMFATestService(t)
	userID := "user-bad-code"
	enc, _ := crypto.EncryptMFASecret(testConfig().MFAEncryptionKey, "JBSWY3DPEHPK3PXP")
	totpRepo.byUser[userID] = &models.UserMFATOTP{UserID: userID, SecretEncrypted: enc, Enabled: true}

	ticket, _, err := svc.CreateLoginTicket(context.Background(), auth.MFAPendingPayload{UserID: userID, TenantID: "t1"})
	require.NoError(t, err)

	_, err = svc.VerifyLoginChallenge(context.Background(), ticket, "000000")
	require.Error(t, err)
}

func TestMFAService_mfaKeyFallback(t *testing.T) {
	cfg := testConfig()
	cfg.MFAEncryptionKey = ""
	svc := ProvideMFAService(cfg, newMFATOTPTestRepo(), newMFARecoveryTestRepo(), auth.NewEphemeralStore(newMemCache(), cfg), &auditCapture{})
	require.Equal(t, cfg.JWTSecret, svc.(*mfaService).mfaKey())
}

func TestRandomRecoveryCode(t *testing.T) {
	code, err := randomRecoveryCode()
	require.NoError(t, err)
	require.Len(t, code, 10)
}
