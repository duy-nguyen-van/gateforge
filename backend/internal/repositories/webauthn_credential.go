package repositories

import (
	"context"
	"encoding/base64"
	"encoding/json"
	stderrors "errors"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/db"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"github.com/go-webauthn/webauthn/webauthn"
	"gorm.io/gorm"
)

// WebauthnCredentialRepository persists passkey credential records.
type WebauthnCredentialRepository interface {
	Create(ctx context.Context, row *models.WebauthnCredential) error
	ListByUserID(ctx context.Context, userID string) ([]models.WebauthnCredential, error)
	ListByUserIDPaginated(ctx context.Context, userID string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.WebauthnCredential], error)
	GetByCredentialID(ctx context.Context, credentialID string) (*models.WebauthnCredential, error)
	UpdateCredentialJSON(ctx context.Context, id string, credentialJSON string, signCount int64) error
	DeleteAllByUserID(ctx context.Context, userID string) (int64, error)
}

type webauthnCredentialRepository struct {
	db *db.PostgresDB
}

func ProvideWebauthnCredentialRepository(db *db.PostgresDB) WebauthnCredentialRepository {
	return &webauthnCredentialRepository{db: db}
}

func (r *webauthnCredentialRepository) Create(ctx context.Context, row *models.WebauthnCredential) error {
	return Create(ctx, r.db, row, DBOp{Operation: "create_webauthn_credential", Resource: "webauthn_credential"}, "Failed to create WebAuthn credential")
}

func (r *webauthnCredentialRepository) ListByUserID(ctx context.Context, userID string) ([]models.WebauthnCredential, error) {
	var rows []models.WebauthnCredential
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at ASC").Find(&rows).Error; err != nil {
		return nil, errors.DatabaseError("Failed to list WebAuthn credentials", err).
			WithOperation("list_webauthn_credentials").
			WithResource("webauthn_credential")
	}
	return rows, nil
}

func (r *webauthnCredentialRepository) ListByUserIDPaginated(ctx context.Context, userID string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.WebauthnCredential], error) {
	query := r.db.WithContext(ctx).Model(&models.WebauthnCredential{}).Where("user_id = ?", userID)
	return Paginate[models.WebauthnCredential](ctx, query, pr, PaginateOptions{
		OrderBy:             "created_at ASC",
		MaxPageSize:         constants.MaxPageSize,
		CountFailureMessage: "Failed to count WebAuthn credentials",
		FindFailureMessage:  "Failed to list WebAuthn credentials",
		DBOp:                DBOp{Operation: "list_webauthn_credentials", Resource: "webauthn_credential"},
	})
}

func (r *webauthnCredentialRepository) GetByCredentialID(ctx context.Context, credentialID string) (*models.WebauthnCredential, error) {
	var row models.WebauthnCredential
	err := r.db.WithContext(ctx).Where("credential_id = ?", credentialID).First(&row).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFoundError("WebAuthn credential", err).
				WithOperation("get_webauthn_credential").
				WithResource("webauthn_credential")
		}
		return nil, errors.DatabaseError("Failed to get WebAuthn credential", err).
			WithOperation("get_webauthn_credential").
			WithResource("webauthn_credential")
	}
	return &row, nil
}

func (r *webauthnCredentialRepository) DeleteAllByUserID(ctx context.Context, userID string) (int64, error) {
	res := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&models.WebauthnCredential{})
	if res.Error != nil {
		return 0, errors.DatabaseError("Failed to delete WebAuthn credentials", res.Error).
			WithOperation("delete_webauthn_credentials").
			WithResource("webauthn_credential")
	}
	return res.RowsAffected, nil
}

func (r *webauthnCredentialRepository) UpdateCredentialJSON(ctx context.Context, id string, credentialJSON string, signCount int64) error {
	res := r.db.WithContext(ctx).Model(&models.WebauthnCredential{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"public_key": credentialJSON,
			"sign_count": signCount,
		})
	if res.Error != nil {
		return errors.DatabaseError("Failed to update WebAuthn credential", res.Error).
			WithOperation("update_webauthn_credential").
			WithResource("webauthn_credential")
	}
	return nil
}

// CredentialIDString encodes raw credential ID for DB unique index.
func CredentialIDString(id []byte) string {
	return base64.RawURLEncoding.EncodeToString(id)
}

// MarshalWebauthnCredential serializes a go-webauthn Credential for storage in public_key column.
func MarshalWebauthnCredential(c *webauthn.Credential) (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// UnmarshalWebauthnCredential loads a Credential from DB JSON.
func UnmarshalWebauthnCredential(data string) (*webauthn.Credential, error) {
	var c webauthn.Credential
	if err := json.Unmarshal([]byte(data), &c); err != nil {
		return nil, err
	}
	return &c, nil
}
