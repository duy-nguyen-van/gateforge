package repositories

import (
	"encoding/json"
	"testing"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/stretchr/testify/require"
)

func TestWebauthnCredentialRepository_CRUD(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideWebauthnCredentialRepository(pg)
	ctx := testCtx()
	user := seedUser(t, pg, "passkey@test.com")

	row := &models.WebauthnCredential{
		BaseModel:    models.NewBaseModel(),
		UserID:       user.ID,
		CredentialID: CredentialIDString([]byte("cred-id-bytes")),
		PublicKey:    `{"id":"cred"}`,
		SignCount:    1,
		DeviceName:   "MacBook",
	}
	require.NoError(t, repo.Create(ctx, row))

	got, err := repo.GetByCredentialID(ctx, row.CredentialID)
	require.NoError(t, err)
	require.Equal(t, row.ID, got.ID)

	list, err := repo.ListByUserID(ctx, user.ID)
	require.NoError(t, err)
	require.Len(t, list, 1)

	paged, err := repo.ListByUserIDPaginated(ctx, user.ID, &dtos.PageableRequest{Page: 1, PageSize: 10})
	require.NoError(t, err)
	require.Len(t, paged.Data, 1)

	require.NoError(t, repo.UpdateCredentialJSON(ctx, row.ID, `{"id":"updated"}`, 42))
	got, err = repo.GetByCredentialID(ctx, row.CredentialID)
	require.NoError(t, err)
	require.Equal(t, int64(42), got.SignCount)

	deleted, err := repo.DeleteAllByUserID(ctx, user.ID)
	require.NoError(t, err)
	require.Equal(t, int64(1), deleted)
}

func TestWebauthnCredentialRepository_GetByCredentialID_NotFound(t *testing.T) {
	pg := newTestDB(t)
	repo := ProvideWebauthnCredentialRepository(pg)

	_, err := repo.GetByCredentialID(testCtx(), "missing")
	requireNotFound(t, err)
}

func TestWebauthnCredentialRepository_Create_DatabaseError(t *testing.T) {
	pg := closedTestDB(t)
	repo := ProvideWebauthnCredentialRepository(pg)
	row := &models.WebauthnCredential{
		BaseModel: models.NewBaseModel(), UserID: "u", CredentialID: "c", PublicKey: "{}",
	}

	err := repo.Create(testCtx(), row)
	requireDatabaseErr(t, err)
}

func TestCredentialIDString(t *testing.T) {
	raw := []byte{1, 2, 3}
	encoded := CredentialIDString(raw)
	require.NotEmpty(t, encoded)
	require.Equal(t, encoded, CredentialIDString(raw))
}

func TestMarshalUnmarshalWebauthnCredential(t *testing.T) {
	cred := &webauthn.Credential{
		ID:        []byte("id"),
		PublicKey: []byte("pk"),
	}
	data, err := MarshalWebauthnCredential(cred)
	require.NoError(t, err)

	loaded, err := UnmarshalWebauthnCredential(data)
	require.NoError(t, err)
	require.Equal(t, cred.ID, loaded.ID)

	_, err = UnmarshalWebauthnCredential("not-json")
	require.Error(t, err)

	_, err = MarshalWebauthnCredential(&webauthn.Credential{})
	require.NoError(t, err)
	// Ensure marshal produces valid JSON
	var tmp json.RawMessage
	require.NoError(t, json.Unmarshal([]byte(data), &tmp))
}
