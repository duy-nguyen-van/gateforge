package repositories

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/lib/pq"

	"github.com/gateforge-iam/gateforge-iam/internal/db"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"

	"gorm.io/gorm"
)

// ClientPatch updates mutable OAuth client fields.
type ClientPatch struct {
	Name         *string
	RedirectUris *[]string
	GrantTypes   *[]string
	Scopes       *[]string
	IsPublic     *bool
	ClientSecret *string
}

// ClientRepository loads OAuth2 clients for OIDC.
type ClientRepository interface {
	GetByTenantAndClientID(ctx context.Context, tenantID, clientID string) (*models.Client, error)
	GetByClientID(ctx context.Context, clientID string) (*models.Client, error)
	GetByID(ctx context.Context, id string) (*models.Client, error)
	List(ctx context.Context, tenantID string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.Client], error)
	Create(ctx context.Context, client *models.Client) error
	Update(ctx context.Context, id string, patch ClientPatch) (*models.Client, error)
	Delete(ctx context.Context, id string) error
	ClientIDTaken(ctx context.Context, tenantID, clientID, excludeID string) (bool, error)
}

type clientRepository struct {
	db *db.PostgresDB
}

// ProvideClientRepository wires client persistence.
func ProvideClientRepository(db *db.PostgresDB) ClientRepository {
	return &clientRepository{db: db}
}

func (r *clientRepository) GetByTenantAndClientID(ctx context.Context, tenantID, clientID string) (*models.Client, error) {
	var c models.Client
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND client_id = ?", tenantID, clientID).
		First(&c).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFoundError("OAuth client", err).
				WithOperation("get_client").
				WithResource("client")
		}
		return nil, errors.DatabaseError("Failed to load OAuth client", err).
			WithOperation("get_client").
			WithResource("client")
	}
	return &c, nil
}

func (r *clientRepository) GetByClientID(ctx context.Context, clientID string) (*models.Client, error) {
	var c models.Client
	err := r.db.WithContext(ctx).
		Where("client_id = ?", clientID).
		First(&c).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFoundError("OAuth client", err).
				WithOperation("get_client").
				WithResource("client")
		}
		return nil, errors.DatabaseError("Failed to load OAuth client", err).
			WithOperation("get_client").
			WithResource("client")
	}
	return &c, nil
}

func (r *clientRepository) GetByID(ctx context.Context, id string) (*models.Client, error) {
	var c models.Client
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&c).Error
	if err != nil {
		if stderrors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.NotFoundError("OAuth client", err).
				WithOperation("get_client_by_id").
				WithResource("client")
		}
		return nil, errors.DatabaseError("Failed to load OAuth client", err).
			WithOperation("get_client_by_id").
			WithResource("client")
	}
	return &c, nil
}

func (r *clientRepository) List(ctx context.Context, tenantID string, pr *dtos.PageableRequest) (*dtos.DataResponse[models.Client], error) {
	query := r.db.WithContext(ctx).Model(&models.Client{})
	if tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	return Paginate[models.Client](ctx, query, pr, PaginateOptions{
		OrderBy:             "created_at DESC",
		CountFailureMessage: "Failed to count clients",
		FindFailureMessage:  "Failed to list clients",
		DBOp:                DBOp{Operation: "list_clients", Resource: "client"},
	})
}

func (r *clientRepository) Create(ctx context.Context, client *models.Client) error {
	return Create(ctx, r.db, client, DBOp{Operation: "create_client", Resource: "client"}, "Failed to create OAuth client")
}

func (r *clientRepository) Update(ctx context.Context, id string, patch ClientPatch) (*models.Client, error) {
	updates := map[string]any{}
	if patch.Name != nil {
		updates["name"] = *patch.Name
	}
	if patch.RedirectUris != nil {
		updates["redirect_uris"] = pq.StringArray(*patch.RedirectUris)
	}
	if patch.GrantTypes != nil {
		updates["grant_types"] = pq.StringArray(*patch.GrantTypes)
	}
	if patch.Scopes != nil {
		updates["scopes"] = pq.StringArray(*patch.Scopes)
	}
	if patch.IsPublic != nil {
		updates["is_public"] = *patch.IsPublic
	}
	if patch.ClientSecret != nil {
		updates["client_secret"] = *patch.ClientSecret
	}
	if len(updates) == 0 {
		return r.GetByID(ctx, id)
	}

	res := r.db.WithContext(ctx).Model(&models.Client{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return nil, errors.DatabaseError("Failed to update OAuth client", res.Error).
			WithOperation("update_client").
			WithResource("client")
	}
	if res.RowsAffected == 0 {
		return nil, errors.NotFoundError("OAuth client", nil).
			WithOperation("update_client").
			WithResource("client")
	}
	return r.GetByID(ctx, id)
}

func (r *clientRepository) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&models.Client{}, "id = ?", id)
	if res.Error != nil {
		return errors.DatabaseError("Failed to delete OAuth client", res.Error).
			WithOperation("delete_client").
			WithResource("client")
	}
	if res.RowsAffected == 0 {
		return errors.NotFoundError("OAuth client", nil).
			WithOperation("delete_client").
			WithResource("client")
	}
	return nil
}

func (r *clientRepository) ClientIDTaken(ctx context.Context, tenantID, clientID, excludeID string) (bool, error) {
	clientID = strings.TrimSpace(clientID)
	if clientID == "" {
		return false, nil
	}
	query := r.db.WithContext(ctx).Model(&models.Client{}).
		Where("tenant_id = ? AND client_id = ?", tenantID, clientID)
	if excludeID != "" {
		query = query.Where("id != ?", excludeID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, errors.DatabaseError("Failed to check OAuth client ID", err).
			WithOperation("check_client_id").
			WithResource("client")
	}
	return count > 0, nil
}
