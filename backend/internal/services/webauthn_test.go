package services

import (
	"context"
	"testing"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/stretchr/testify/require"
)

type webauthnCredTestRepo struct {
	byUser map[string][]models.WebauthnCredential
}

func newWebauthnCredTestRepo() *webauthnCredTestRepo {
	return &webauthnCredTestRepo{byUser: map[string][]models.WebauthnCredential{}}
}

func (r *webauthnCredTestRepo) Create(_ context.Context, row *models.WebauthnCredential) error {
	r.byUser[row.UserID] = append(r.byUser[row.UserID], *row)
	return nil
}

func (r *webauthnCredTestRepo) ListByUserID(_ context.Context, userID string) ([]models.WebauthnCredential, error) {
	return r.byUser[userID], nil
}

func (r *webauthnCredTestRepo) ListByUserIDPaginated(_ context.Context, userID string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.WebauthnCredential], error) {
	rows := r.byUser[userID]
	page, pageable := dtos.PaginateSlice(rows, pr)
	return &dtos.DataResponse[models.WebauthnCredential]{Data: page, Pageable: pageable}, nil
}

func (r *webauthnCredTestRepo) GetByCredentialID(_ context.Context, credID string) (*models.WebauthnCredential, error) {
	for _, rows := range r.byUser {
		for i := range rows {
			if rows[i].CredentialID == credID {
				return &rows[i], nil
			}
		}
	}
	return nil, context.Canceled
}

func (r *webauthnCredTestRepo) UpdateCredentialJSON(_ context.Context, id, jsonStr string, signCount int64) error {
	return nil
}

func (r *webauthnCredTestRepo) DeleteAllByUserID(_ context.Context, userID string) (int64, error) {
	n := int64(len(r.byUser[userID]))
	delete(r.byUser, userID)
	return n, nil
}

func newWebauthnTestService(t *testing.T, users *userTestRepo, creds *webauthnCredTestRepo) WebauthnService {
	t.Helper()
	cfg := testConfig()
	wa, err := auth.ProvideWebAuthn(cfg)
	require.NoError(t, err)
	return ProvideWebauthnService(cfg, wa, users, creds, auth.NewEphemeralStore(newMemCache(), cfg), &auditCapture{})
}

func TestWebauthnService_ListCredentials(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("wa@example.com", "secret")
	creds := newWebauthnCredTestRepo()
	creds.byUser[u.ID] = []models.WebauthnCredential{{BaseModel: models.NewBaseModel(), UserID: u.ID, DeviceName: "Mac"}}
	svc := newWebauthnTestService(t, users, creds)

	rows, pageable, err := svc.ListCredentials(context.Background(), u.ID, dtos.NewPageableRequest())
	require.NoError(t, err)
	require.Len(t, rows, 1)
	require.Equal(t, int64(1), pageable.Total)
}

func TestWebauthnService_RegisterStart(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("wa@example.com", "secret")
	svc := newWebauthnTestService(t, users, newWebauthnCredTestRepo())

	opts, token, err := svc.RegisterStart(context.Background(), u.ID, "Laptop")
	require.NoError(t, err)
	require.NotEmpty(t, opts)
	require.NotEmpty(t, token)
}

func TestWebauthnService_LoginStart_UnknownUser(t *testing.T) {
	users := newUserTestRepo()
	svc := newWebauthnTestService(t, users, newWebauthnCredTestRepo())

	opts, token, err := svc.LoginStart(context.Background(), "unknown@example.com")
	require.NoError(t, err)
	require.NotEmpty(t, opts)
	require.NotEmpty(t, token)
}

func TestWebauthnService_LoginStart_InactiveUser(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("inactive@example.com", "secret")
	u.Status = constants.UserStatusDisabled
	svc := newWebauthnTestService(t, users, newWebauthnCredTestRepo())

	_, token, err := svc.LoginStart(context.Background(), "inactive@example.com")
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestWebauthnService_LoginFinish_InvalidSession(t *testing.T) {
	users := newUserTestRepo()
	users.seed("wa@example.com", "secret")
	svc := newWebauthnTestService(t, users, newWebauthnCredTestRepo())

	_, err := svc.LoginFinish(context.Background(), "wa@example.com", "bad-token", []byte(`{}`))
	require.Error(t, err)
}

func TestWebauthnService_RegisterFinish_InvalidSession(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("wa@example.com", "secret")
	svc := newWebauthnTestService(t, users, newWebauthnCredTestRepo())

	err := svc.RegisterFinish(context.Background(), u.ID, "bad-token", []byte(`{}`))
	require.Error(t, err)
}

func TestWebauthnService_LoginFinish_NoCredentials(t *testing.T) {
	users := newUserTestRepo()
	users.seed("wa@example.com", "secret")
	cfg := testConfig()
	wa, err := webauthn.New(&webauthn.Config{
		RPID: cfg.WebauthnRPID, RPDisplayName: cfg.WebauthnRPDisplayName, RPOrigins: cfg.WebauthnRPOrigins,
	})
	require.NoError(t, err)
	ephemeral := auth.NewEphemeralStore(newMemCache(), cfg)
	_, session, err := wa.BeginDiscoverableLogin()
	require.NoError(t, err)
	env := webauthnLoginEnvelope{Session: *session}
	token, err := ephemeral.PutWebauthnLoginSession(context.Background(), &env)
	require.NoError(t, err)

	svc := ProvideWebauthnService(cfg, wa, users, newWebauthnCredTestRepo(), ephemeral, &auditCapture{})
	_, err = svc.LoginFinish(context.Background(), "wa@example.com", token, []byte(`{}`))
	require.Error(t, err)
}
