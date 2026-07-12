package services

import (
	"context"
	"testing"

	"github.com/descope/virtualwebauthn"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/domains"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/stretchr/testify/require"
)

func webauthnVirtualRP(cfg *config.Config) virtualwebauthn.RelyingParty {
	return virtualwebauthn.RelyingParty{
		Name:   cfg.WebauthnRPDisplayName,
		ID:     cfg.WebauthnRPID,
		Origin: cfg.WebauthnRPOrigins[0],
	}
}

func TestWebauthnService_FullRegisterLoginFlow(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("passkey@example.com", "secret")
	creds := newWebauthnCredTestRepo()
	cfg := testConfig()
	svc := newWebauthnTestService(t, users, creds)
	ctx := context.Background()
	rp := webauthnVirtualRP(cfg)

	opts, token, err := svc.RegisterStart(ctx, u.ID, "Laptop")
	require.NoError(t, err)
	parsed, err := virtualwebauthn.ParseAttestationOptions(string(opts))
	require.NoError(t, err)
	authenticator := virtualwebauthn.NewAuthenticator()
	credential := virtualwebauthn.NewCredential(virtualwebauthn.KeyTypeEC2)
	attestation := virtualwebauthn.CreateAttestationResponse(rp, authenticator, credential, *parsed)
	require.NoError(t, svc.RegisterFinish(ctx, u.ID, token, []byte(attestation)))
	require.Len(t, creds.byUser[u.ID], 1)
	require.Equal(t, "Laptop", creds.byUser[u.ID][0].DeviceName)

	waUser := &domains.WebAuthnUser{User: u}
	authenticator.Options.UserHandle = waUser.WebAuthnID()
	authenticator.AddCredential(credential)

	loginOpts, loginToken, err := svc.LoginStart(ctx, "passkey@example.com")
	require.NoError(t, err)
	loginParsed, err := virtualwebauthn.ParseAssertionOptions(string(loginOpts))
	require.NoError(t, err)
	assertion := virtualwebauthn.CreateAssertionResponse(rp, authenticator, credential, *loginParsed)
	got, err := svc.LoginFinish(ctx, "passkey@example.com", loginToken, []byte(assertion))
	require.NoError(t, err)
	require.Equal(t, u.ID, got.ID)
}

func TestWebauthnService_RegisterFinish_DefaultDeviceName(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("default-device@example.com", "secret")
	creds := newWebauthnCredTestRepo()
	cfg := testConfig()
	svc := newWebauthnTestService(t, users, creds)
	ctx := context.Background()
	rp := webauthnVirtualRP(cfg)

	opts, token, err := svc.RegisterStart(ctx, u.ID, "   ")
	require.NoError(t, err)
	parsed, err := virtualwebauthn.ParseAttestationOptions(string(opts))
	require.NoError(t, err)
	authenticator := virtualwebauthn.NewAuthenticator()
	credential := virtualwebauthn.NewCredential(virtualwebauthn.KeyTypeEC2)
	attestation := virtualwebauthn.CreateAttestationResponse(rp, authenticator, credential, *parsed)
	require.NoError(t, svc.RegisterFinish(ctx, u.ID, token, []byte(attestation)))
	require.Equal(t, "Passkey", creds.byUser[u.ID][0].DeviceName)
}

func TestWebauthnService_LoginFinish_UnknownUser(t *testing.T) {
	users := newUserTestRepo()
	creds := newWebauthnCredTestRepo()
	svc := newWebauthnTestService(t, users, creds)
	ctx := context.Background()

	opts, token, err := svc.LoginStart(ctx, "ghost@example.com")
	require.NoError(t, err)

	_, err = svc.LoginFinish(ctx, "ghost@example.com", token, opts)
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.Equal(t, errors.ErrorTypeUnauthorized, appErr.Type)
}

func TestWebauthnService_LoginFinish_ValidationFailed(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("bad-assert@example.com", "secret")
	creds := newWebauthnCredTestRepo()
	cfg := testConfig()
	svc := newWebauthnTestService(t, users, creds)
	ctx := context.Background()
	rp := webauthnVirtualRP(cfg)

	opts, token, err := svc.RegisterStart(ctx, u.ID, "Key")
	require.NoError(t, err)
	parsed, err := virtualwebauthn.ParseAttestationOptions(string(opts))
	require.NoError(t, err)
	authenticator := virtualwebauthn.NewAuthenticator()
	credential := virtualwebauthn.NewCredential(virtualwebauthn.KeyTypeEC2)
	attestation := virtualwebauthn.CreateAttestationResponse(rp, authenticator, credential, *parsed)
	require.NoError(t, svc.RegisterFinish(ctx, u.ID, token, []byte(attestation)))

	loginOpts, loginToken, err := svc.LoginStart(ctx, "bad-assert@example.com")
	require.NoError(t, err)
	_, err = virtualwebauthn.ParseAssertionOptions(string(loginOpts))
	require.NoError(t, err)

	_, err = svc.LoginFinish(ctx, "bad-assert@example.com", loginToken, []byte(`{"type":"public-key","id":"x","response":{}}`))
	require.Error(t, err)
}

func TestWebauthnService_RegisterFinish_UserNotFoundAfterSession(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("gone@example.com", "secret")
	creds := newWebauthnCredTestRepo()
	cfg := testConfig()
	svc := newWebauthnTestService(t, users, creds)
	ctx := context.Background()
	rp := webauthnVirtualRP(cfg)

	opts, token, err := svc.RegisterStart(ctx, u.ID, "Key")
	require.NoError(t, err)
	delete(users.users, u.ID)

	parsed, err := virtualwebauthn.ParseAttestationOptions(string(opts))
	require.NoError(t, err)
	authenticator := virtualwebauthn.NewAuthenticator()
	credential := virtualwebauthn.NewCredential(virtualwebauthn.KeyTypeEC2)
	attestation := virtualwebauthn.CreateAttestationResponse(rp, authenticator, credential, *parsed)
	err = svc.RegisterFinish(ctx, u.ID, token, []byte(attestation))
	require.Error(t, err)
}

func TestWebauthnService_ListCredentials_Error(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("listerr@example.com", "secret")
	creds := &webauthnCredTestRepo{byUser: map[string][]models.WebauthnCredential{}}
	svc := newWebauthnTestService(t, users, creds)

	_, _, err := svc.ListCredentials(context.Background(), u.ID, nil)
	require.NoError(t, err)
}

func TestWebauthnService_LoadUserForWebAuthn_ListError(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("loaderr@example.com", "secret")
	creds := &errWebauthnCredRepo{listErr: context.Canceled}
	cfg := testConfig()
	wa, err := auth.ProvideWebAuthn(cfg)
	require.NoError(t, err)
	svc := ProvideWebauthnService(cfg, wa, users, creds, auth.NewEphemeralStore(newMemCache(), cfg), &auditCapture{})

	_, _, err = svc.RegisterStart(context.Background(), u.ID, "Key")
	require.Error(t, err)
}

func TestWebauthnService_LoginStart_CredListError(t *testing.T) {
	users := newUserTestRepo()
	users.seed("loginlist@example.com", "secret")
	creds := &errWebauthnCredRepo{listErr: context.Canceled}
	cfg := testConfig()
	wa, err := auth.ProvideWebAuthn(cfg)
	require.NoError(t, err)
	svc := ProvideWebauthnService(cfg, wa, users, creds, auth.NewEphemeralStore(newMemCache(), cfg), &auditCapture{})

	_, _, err = svc.LoginStart(context.Background(), "loginlist@example.com")
	require.Error(t, err)
}

type errWebauthnCredRepo struct {
	listErr error
}

func (r *errWebauthnCredRepo) Create(context.Context, *models.WebauthnCredential) error { return nil }
func (r *errWebauthnCredRepo) ListByUserID(context.Context, string) ([]models.WebauthnCredential, error) {
	return nil, r.listErr
}
func (r *errWebauthnCredRepo) ListByUserIDPaginated(context.Context, string, *dtos.PageableRequest) (*dtos.DataResponse[models.WebauthnCredential], error) {
	return nil, r.listErr
}
func (r *errWebauthnCredRepo) GetByCredentialID(context.Context, string) (*models.WebauthnCredential, error) {
	return nil, context.Canceled
}
func (r *errWebauthnCredRepo) UpdateCredentialJSON(context.Context, string, string, int64) error {
	return nil
}
func (r *errWebauthnCredRepo) DeleteAllByUserID(context.Context, string) (int64, error) {
	return 0, nil
}

type trackingWebauthnCredRepo struct {
	*webauthnCredTestRepo
	updateErr error
}

func (r *trackingWebauthnCredRepo) UpdateCredentialJSON(ctx context.Context, id, jsonStr string, signCount int64) error {
	if r.updateErr != nil {
		return r.updateErr
	}
	return r.webauthnCredTestRepo.UpdateCredentialJSON(ctx, id, jsonStr, signCount)
}

func TestWebauthnService_RegisterFinish_VerificationFailed(t *testing.T) {
	users := newUserTestRepo()
	u1 := users.seed("user1@example.com", "secret")
	u2 := users.seed("user2@example.com", "secret")
	creds := newWebauthnCredTestRepo()
	cfg := testConfig()
	svc := newWebauthnTestService(t, users, creds)
	ctx := context.Background()
	rp := webauthnVirtualRP(cfg)

	opts, token, err := svc.RegisterStart(ctx, u1.ID, "Key")
	require.NoError(t, err)
	parsed, err := virtualwebauthn.ParseAttestationOptions(string(opts))
	require.NoError(t, err)
	authenticator := virtualwebauthn.NewAuthenticator()
	credential := virtualwebauthn.NewCredential(virtualwebauthn.KeyTypeEC2)
	attestation := virtualwebauthn.CreateAttestationResponse(rp, authenticator, credential, *parsed)

	// Use session token from u1's registration but finish for u2 with same attestation.
	err = svc.RegisterFinish(ctx, u2.ID, token, []byte(attestation))
	require.Error(t, err)
}

func TestWebauthnService_LoginFinish_CredentialUserMismatch(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("mismatch@example.com", "secret")
	other := users.seed("other@example.com", "secret")
	creds := newWebauthnCredTestRepo()
	cfg := testConfig()
	svc := newWebauthnTestService(t, users, creds)
	ctx := context.Background()
	rp := webauthnVirtualRP(cfg)

	opts, token, err := svc.RegisterStart(ctx, u.ID, "Key")
	require.NoError(t, err)
	parsed, err := virtualwebauthn.ParseAttestationOptions(string(opts))
	require.NoError(t, err)
	authenticator := virtualwebauthn.NewAuthenticator()
	credential := virtualwebauthn.NewCredential(virtualwebauthn.KeyTypeEC2)
	attestation := virtualwebauthn.CreateAttestationResponse(rp, authenticator, credential, *parsed)
	require.NoError(t, svc.RegisterFinish(ctx, u.ID, token, []byte(attestation)))

	// Point stored credential at a different user.
	creds.byUser[u.ID][0].UserID = other.ID

	waUser := &domains.WebAuthnUser{User: u}
	u.WebauthnCredentials = creds.byUser[u.ID]
	authenticator.Options.UserHandle = waUser.WebAuthnID()
	authenticator.AddCredential(credential)

	loginOpts, loginToken, err := svc.LoginStart(ctx, "mismatch@example.com")
	require.NoError(t, err)
	loginParsed, err := virtualwebauthn.ParseAssertionOptions(string(loginOpts))
	require.NoError(t, err)
	assertion := virtualwebauthn.CreateAssertionResponse(rp, authenticator, credential, *loginParsed)
	_, err = svc.LoginFinish(ctx, "mismatch@example.com", loginToken, []byte(assertion))
	require.Error(t, err)
}

func TestWebauthnService_LoginFinish_UpdateCredentialError(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("updateerr@example.com", "secret")
	creds := &trackingWebauthnCredRepo{webauthnCredTestRepo: newWebauthnCredTestRepo(), updateErr: context.Canceled}
	cfg := testConfig()
	wa, err := auth.ProvideWebAuthn(cfg)
	require.NoError(t, err)
	svc := ProvideWebauthnService(cfg, wa, users, creds, auth.NewEphemeralStore(newMemCache(), cfg), &auditCapture{})
	ctx := context.Background()
	rp := webauthnVirtualRP(cfg)

	opts, token, err := svc.RegisterStart(ctx, u.ID, "Key")
	require.NoError(t, err)
	parsed, err := virtualwebauthn.ParseAttestationOptions(string(opts))
	require.NoError(t, err)
	authenticator := virtualwebauthn.NewAuthenticator()
	credential := virtualwebauthn.NewCredential(virtualwebauthn.KeyTypeEC2)
	attestation := virtualwebauthn.CreateAttestationResponse(rp, authenticator, credential, *parsed)
	require.NoError(t, svc.RegisterFinish(ctx, u.ID, token, []byte(attestation)))

	waUser := &domains.WebAuthnUser{User: u}
	authenticator.Options.UserHandle = waUser.WebAuthnID()
	authenticator.AddCredential(credential)
	loginOpts, loginToken, err := svc.LoginStart(ctx, "updateerr@example.com")
	require.NoError(t, err)
	loginParsed, err := virtualwebauthn.ParseAssertionOptions(string(loginOpts))
	require.NoError(t, err)
	assertion := virtualwebauthn.CreateAssertionResponse(rp, authenticator, credential, *loginParsed)
	_, err = svc.LoginFinish(ctx, "updateerr@example.com", loginToken, []byte(assertion))
	require.Error(t, err)
}

func TestWebauthnService_ListCredentials_RepoError(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("listcred@example.com", "secret")
	creds := &errWebauthnCredRepo{listErr: context.Canceled}
	cfg := testConfig()
	wa, err := auth.ProvideWebAuthn(cfg)
	require.NoError(t, err)
	svc := ProvideWebauthnService(cfg, wa, users, creds, auth.NewEphemeralStore(newMemCache(), cfg), &auditCapture{})
	_, _, err = svc.ListCredentials(context.Background(), u.ID, dtos.NewPageableRequest())
	require.Error(t, err)
}
