package services

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/logger"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"go.uber.org/zap"
)

// WebauthnService implements passkey registration and login.
type WebauthnService interface {
	ListCredentials(ctx context.Context, userID string, pr *dtos.PageableRequest) ([]models.WebauthnCredential, *dtos.Pageable, error)
	RegisterStart(ctx context.Context, userID, deviceName string) (options json.RawMessage, sessionToken string, err error)
	RegisterFinish(ctx context.Context, userID, sessionToken string, credentialJSON []byte) error
	LoginStart(ctx context.Context, email string) (options json.RawMessage, sessionToken string, err error)
	LoginFinish(ctx context.Context, email, sessionToken string, credentialJSON []byte) (*models.User, error)
}

type webauthnService struct {
	cfg       *config.Config
	wa        *webauthn.WebAuthn
	userRepo  repositories.UserRepository
	credRepo  repositories.WebauthnCredentialRepository
	ephemeral *auth.EphemeralStore
	audit     AuditService
}

func ProvideWebauthnService(
	cfg *config.Config,
	wa *webauthn.WebAuthn,
	userRepo repositories.UserRepository,
	credRepo repositories.WebauthnCredentialRepository,
	ephemeral *auth.EphemeralStore,
	audit AuditService,
) WebauthnService {
	return &webauthnService{
		cfg:       cfg,
		wa:        wa,
		userRepo:  userRepo,
		credRepo:  credRepo,
		ephemeral: ephemeral,
		audit:     audit,
	}
}

type webauthnRegEnvelope struct {
	Session    webauthn.SessionData `json:"session"`
	DeviceName string               `json:"device_name"`
}

type webauthnLoginEnvelope struct {
	Session webauthn.SessionData `json:"session"`
}

func (s *webauthnService) loadUserForWebAuthn(ctx context.Context, userID string) (*models.User, error) {
	u, err := s.userRepo.GetOneByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	creds, err := s.credRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	u.WebauthnCredentials = creds
	return u, nil
}

func (s *webauthnService) ListCredentials(ctx context.Context, userID string, pr *dtos.PageableRequest) ([]models.WebauthnCredential, *dtos.Pageable, error) {
	result, err := s.credRepo.ListByUserIDPaginated(ctx, userID, pr)
	if err != nil {
		return nil, nil, err
	}
	return result.Data, result.Pageable, nil
}

func (s *webauthnService) RegisterStart(ctx context.Context, userID, deviceName string) (json.RawMessage, string, error) {
	u, err := s.loadUserForWebAuthn(ctx, userID)
	if err != nil {
		return nil, "", err
	}
	waUser := &domains.WebAuthnUser{User: u}
	var exclude []protocol.CredentialDescriptor
	for _, c := range waUser.WebAuthnCredentials() {
		exclude = append(exclude, c.Descriptor())
	}
	creation, session, err := s.wa.BeginRegistration(waUser, webauthn.WithExclusions(exclude))
	if err != nil {
		logger.Log.Error("webauthn register start", zap.Error(err))
		return nil, "", errors.ValidationError("WebAuthn registration failed to start", err).
			WithOperation("webauthn_register_start").
			WithResource("webauthn")
	}
	env := webauthnRegEnvelope{Session: *session, DeviceName: strings.TrimSpace(deviceName)}
	token, err := s.ephemeral.PutWebauthnRegistrationSession(ctx, &env)
	if err != nil {
		return nil, "", err
	}
	// Marshal PublicKeyCredentialCreationOptions only — @simplewebauthn/browser expects the
	// unwrapped options JSON, not go-webauthn's CredentialCreation { publicKey, mediation } envelope.
	raw, err := json.Marshal(creation.Response)
	if err != nil {
		return nil, "", errors.InternalError("Failed to marshal WebAuthn options", err).
			WithOperation("webauthn_register_start").
			WithResource("webauthn")
	}
	return raw, token, nil
}

func (s *webauthnService) RegisterFinish(ctx context.Context, userID, sessionToken string, credentialJSON []byte) error {
	var env webauthnRegEnvelope
	if err := s.ephemeral.TakeWebauthnRegistrationSession(ctx, sessionToken, &env); err != nil {
		return err
	}
	u, err := s.loadUserForWebAuthn(ctx, userID)
	if err != nil {
		return err
	}
	waUser := &domains.WebAuthnUser{User: u}
	parsed, err := protocol.ParseCredentialCreationResponseBytes(credentialJSON)
	if err != nil {
		return errors.ValidationError("Invalid WebAuthn credential", err).
			WithOperation("webauthn_register_finish").
			WithResource("webauthn")
	}
	cred, err := s.wa.CreateCredential(waUser, env.Session, parsed)
	if err != nil {
		logger.Log.Warn("webauthn register finish verify failed", zap.Error(err))
		return errors.ValidationError("WebAuthn registration verification failed", err).
			WithOperation("webauthn_register_finish").
			WithResource("webauthn")
	}
	jsonStr, err := repositories.MarshalWebauthnCredential(cred)
	if err != nil {
		return errors.InternalError("Failed to serialize credential", err).
			WithOperation("webauthn_register_finish").
			WithResource("webauthn")
	}
	device := env.DeviceName
	if device == "" {
		device = "Passkey"
	}
	row := &models.WebauthnCredential{
		BaseModel:    models.NewBaseModel(),
		UserID:       u.ID,
		CredentialID: repositories.CredentialIDString(cred.ID),
		PublicKey:    jsonStr,
		SignCount:    int64(cred.Authenticator.SignCount),
		DeviceName:   device,
	}
	if err := s.credRepo.Create(ctx, row); err != nil {
		return err
	}
	logger.Log.Info("webauthn credential registered",
		zap.String("operation", "webauthn_register_finish"),
		zap.String("user_id", userID))
	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionWebauthnRegister,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      userID,
		ResourceType: constants.AuditResourceTypeWebauthnCredential,
		ResourceID:   row.ID,
		ResourceName: device,
	})
	return nil
}

func (s *webauthnService) LoginStart(ctx context.Context, email string) (json.RawMessage, string, error) {
	emailLower := strings.ToLower(strings.TrimSpace(email))
	u, err := s.userRepo.GetByEmailLower(ctx, emailLower)
	if err != nil {
		if appErr := errors.GetAppError(err); appErr != nil && appErr.Type == errors.ErrorTypeNotFound {
			return s.fakeLoginOptions(ctx)
		}
		return nil, "", err
	}
	if u.Status != constants.UserStatusActive {
		return s.fakeLoginOptions(ctx)
	}
	creds, err := s.credRepo.ListByUserID(ctx, u.ID)
	if err != nil {
		return nil, "", err
	}
	u.WebauthnCredentials = creds
	if len(u.WebauthnCredentials) == 0 {
		return s.fakeLoginOptions(ctx)
	}
	waUser := &domains.WebAuthnUser{User: u}
	assertion, session, err := s.wa.BeginLogin(waUser)
	if err != nil {
		logger.Log.Warn("webauthn login start", zap.Error(err))
		return s.fakeLoginOptions(ctx)
	}
	env := webauthnLoginEnvelope{Session: *session}
	token, err := s.ephemeral.PutWebauthnLoginSession(ctx, &env)
	if err != nil {
		return nil, "", err
	}
	raw, err := json.Marshal(assertion.Response)
	if err != nil {
		return nil, "", errors.InternalError("Failed to marshal WebAuthn options", err).
			WithOperation("webauthn_login_start").
			WithResource("webauthn")
	}
	return raw, token, nil
}

// fakeLoginOptions returns a syntactically valid assertion when the user is unknown or has no passkeys,
// to reduce email enumeration and timing leaks vs absent users.
func (s *webauthnService) fakeLoginOptions(ctx context.Context) (json.RawMessage, string, error) {
	assertion, session, err := s.wa.BeginDiscoverableLogin()
	if err != nil {
		return nil, "", errors.InternalError("WebAuthn login failed to start", err).
			WithOperation("webauthn_login_start").
			WithResource("webauthn")
	}
	env := webauthnLoginEnvelope{Session: *session}
	token, err := s.ephemeral.PutWebauthnLoginSession(ctx, &env)
	if err != nil {
		return nil, "", err
	}
	raw, err := json.Marshal(assertion.Response)
	if err != nil {
		return nil, "", err
	}
	return raw, token, nil
}

func (s *webauthnService) LoginFinish(ctx context.Context, email, sessionToken string, credentialJSON []byte) (*models.User, error) {
	var env webauthnLoginEnvelope
	if err := s.ephemeral.TakeWebauthnLoginSession(ctx, sessionToken, &env); err != nil {
		return nil, err
	}
	emailLower := strings.ToLower(strings.TrimSpace(email))
	u, err := s.userRepo.GetByEmailLower(ctx, emailLower)
	if err != nil {
		return nil, errors.UnauthorizedError("WebAuthn authentication failed", nil).
			WithOperation("webauthn_login_finish").
			WithResource("webauthn")
	}
	if u.Status != constants.UserStatusActive {
		return nil, errors.UnauthorizedError("WebAuthn authentication failed", nil).
			WithOperation("webauthn_login_finish").
			WithResource("webauthn")
	}
	creds, err := s.credRepo.ListByUserID(ctx, u.ID)
	if err != nil {
		return nil, err
	}
	u.WebauthnCredentials = creds
	if len(u.WebauthnCredentials) == 0 {
		return nil, errors.UnauthorizedError("WebAuthn authentication failed", nil).
			WithOperation("webauthn_login_finish").
			WithResource("webauthn")
	}
	waUser := &domains.WebAuthnUser{User: u}
	parsed, err := protocol.ParseCredentialRequestResponseBytes(credentialJSON)
	if err != nil {
		return nil, errors.ValidationError("Invalid WebAuthn assertion", err).
			WithOperation("webauthn_login_finish").
			WithResource("webauthn")
	}
	cred, err := s.wa.ValidateLogin(waUser, env.Session, parsed)
	if err != nil {
		s.audit.Record(ctx, domains.AuditRecordParams{
			Action:    constants.AuditActionWebauthnLogin,
			Result:    constants.AuditResultFailure,
			ActorType: constants.AuditActorTypeUser,
		})
		return nil, errors.UnauthorizedError("WebAuthn authentication failed", nil).
			WithOperation("webauthn_login_finish").
			WithResource("webauthn")
	}
	credID := repositories.CredentialIDString(cred.ID)
	row, err := s.credRepo.GetByCredentialID(ctx, credID)
	if err != nil {
		return nil, errors.UnauthorizedError("WebAuthn authentication failed", nil).
			WithOperation("webauthn_login_finish").
			WithResource("webauthn")
	}
	if row.UserID != u.ID {
		return nil, errors.UnauthorizedError("WebAuthn authentication failed", nil).
			WithOperation("webauthn_login_finish").
			WithResource("webauthn")
	}
	updatedJSON, err := repositories.MarshalWebauthnCredential(cred)
	if err != nil {
		return nil, err
	}
	if err := s.credRepo.UpdateCredentialJSON(ctx, row.ID, updatedJSON, int64(cred.Authenticator.SignCount)); err != nil {
		return nil, err
	}
	logger.Log.Info("webauthn login success",
		zap.String("operation", "webauthn_login_finish"),
		zap.String("user_id", u.ID))
	s.audit.Record(ctx, domains.AuditRecordParams{
		Action:       constants.AuditActionWebauthnLogin,
		Result:       constants.AuditResultSuccess,
		ActorType:    constants.AuditActorTypeUser,
		ActorID:      u.ID,
		ResourceType: constants.AuditResourceTypeUser,
		ResourceID:   u.ID,
	})
	return u, nil
}
