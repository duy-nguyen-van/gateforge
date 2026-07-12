package services

import (
	"context"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/stretchr/testify/require"
)

func TestWebauthnService_LoginStart_WithCredentials(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("creds@example.com", "secret")
	creds := newWebauthnCredTestRepo()
	creds.byUser[u.ID] = []models.WebauthnCredential{{
		BaseModel: models.NewBaseModel(), UserID: u.ID,
		CredentialID: "fake-cred", PublicKey: `{}`, DeviceName: "Key",
	}}
	svc := newWebauthnTestService(t, users, creds)

	opts, token, err := svc.LoginStart(context.Background(), "creds@example.com")
	require.NoError(t, err)
	require.NotEmpty(t, opts)
	require.NotEmpty(t, token)
}

func TestWebauthnService_RegisterStart_UserNotFound(t *testing.T) {
	users := newUserTestRepo()
	svc := newWebauthnTestService(t, users, newWebauthnCredTestRepo())

	_, _, err := svc.RegisterStart(context.Background(), "missing-user", "Device")
	require.Error(t, err)
}

func TestWebauthnService_LoginFinish_InactiveUser(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("inactive@example.com", "secret")
	u.Status = constants.UserStatusDisabled
	creds := newWebauthnCredTestRepo()
	creds.byUser[u.ID] = []models.WebauthnCredential{{BaseModel: models.NewBaseModel(), UserID: u.ID, CredentialID: "c1"}}
	svc := newWebauthnTestService(t, users, creds)

	_, err := svc.LoginFinish(context.Background(), "inactive@example.com", "token", []byte(`{}`))
	require.Error(t, err)
}

func TestWebauthnService_RegisterFinish_InvalidCredentialJSON(t *testing.T) {
	users := newUserTestRepo()
	u := users.seed("wa@example.com", "secret")
	svc := newWebauthnTestService(t, users, newWebauthnCredTestRepo())

	opts, token, err := svc.RegisterStart(context.Background(), u.ID, "Laptop")
	require.NoError(t, err)
	require.NotEmpty(t, opts)

	err = svc.RegisterFinish(context.Background(), u.ID, token, []byte(`not-json`))
	require.Error(t, err)
}
