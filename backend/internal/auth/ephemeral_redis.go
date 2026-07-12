package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/cache"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"

	"github.com/google/uuid"
)

const (
	redisKeyWebauthnReg   = "iam:webauthn:reg:%s"
	redisKeyWebauthnLogin = "iam:webauthn:login:%s"
	redisKeyMFAPending    = "iam:mfa:pending:%s"
)

// EphemeralStore keeps short-lived WebAuthn session data and MFA pending tickets in Redis.
type EphemeralStore struct {
	cache cache.Cache
	cfg   *config.Config
}

func NewEphemeralStore(c cache.Cache, cfg *config.Config) *EphemeralStore {
	return &EphemeralStore{cache: c, cfg: cfg}
}

func (s *EphemeralStore) PutWebauthnRegistrationSession(ctx context.Context, data any) (token string, err error) {
	return s.putJSON(ctx, redisKeyWebauthnReg, data, s.cfg.WebauthnSessionTTL)
}

func (s *EphemeralStore) TakeWebauthnRegistrationSession(ctx context.Context, token string, dest any) error {
	return s.takeJSON(ctx, fmt.Sprintf(redisKeyWebauthnReg, token), dest)
}

func (s *EphemeralStore) PutWebauthnLoginSession(ctx context.Context, data any) (token string, err error) {
	return s.putJSON(ctx, redisKeyWebauthnLogin, data, s.cfg.WebauthnSessionTTL)
}

func (s *EphemeralStore) TakeWebauthnLoginSession(ctx context.Context, token string, dest any) error {
	return s.takeJSON(ctx, fmt.Sprintf(redisKeyWebauthnLogin, token), dest)
}

// MFAPendingPayload is stored until the user completes MFA after password or passkey.
type MFAPendingPayload struct {
	UserID     string `json:"user_id"`
	TenantID   string `json:"tenant_id"`
	RememberMe bool   `json:"remember_me"`
	ReturnTo   string `json:"return_to,omitempty"`
}

func (s *EphemeralStore) PutMFAPending(ctx context.Context, p MFAPendingPayload) (ticket string, err error) {
	ticket = uuid.NewString()
	key := fmt.Sprintf(redisKeyMFAPending, ticket)
	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	if err := s.cache.Set(ctx, key, string(b), s.cfg.MFAPendingTicketTTL); err != nil {
		return "", err
	}
	return ticket, nil
}

func (s *EphemeralStore) TakeMFAPending(ctx context.Context, ticket string) (*MFAPendingPayload, error) {
	key := fmt.Sprintf(redisKeyMFAPending, ticket)
	raw, err := s.cache.Get(ctx, key)
	if err != nil {
		if appErr := errors.GetAppError(err); appErr != nil && appErr.Type == errors.ErrorTypeNotFound {
			return nil, errors.UnauthorizedError("Invalid or expired MFA ticket", nil).
				WithOperation("mfa_pending").
				WithResource("mfa_ticket")
		}
		return nil, err
	}
	_ = s.cache.Delete(ctx, key) // single-use; ignore delete error
	var p MFAPendingPayload
	if err := json.Unmarshal([]byte(raw), &p); err != nil {
		return nil, errors.ValidationError("Invalid MFA ticket payload", err)
	}
	return &p, nil
}

func (s *EphemeralStore) putJSON(ctx context.Context, keyFmt string, v any, ttl time.Duration) (token string, err error) {
	token = uuid.NewString()
	key := fmt.Sprintf(keyFmt, token)
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	if err := s.cache.Set(ctx, key, string(b), ttl); err != nil {
		return "", err
	}
	return token, nil
}

func (s *EphemeralStore) takeJSON(ctx context.Context, key string, dest any) error {
	raw, err := s.cache.Get(ctx, key)
	if err != nil {
		if appErr := errors.GetAppError(err); appErr != nil && appErr.Type == errors.ErrorTypeNotFound {
			return errors.UnauthorizedError("Invalid or expired WebAuthn session", nil).
				WithOperation("webauthn_session").
				WithResource("session")
		}
		return err
	}
	_ = s.cache.Delete(ctx, key)
	if err := json.Unmarshal([]byte(raw), dest); err != nil {
		return errors.ValidationError("Invalid WebAuthn session payload", err)
	}
	return nil
}
