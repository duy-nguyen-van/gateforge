package repositories

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/db"
	"github.com/gateforge-iam/gateforge-iam/internal/errors"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// sqliteTestSchema mirrors production tables without PostgreSQL-specific column types.
var sqliteTestSchema = []string{
	`CREATE TABLE users (
		id TEXT PRIMARY KEY,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		first_name TEXT,
		last_name TEXT,
		email TEXT NOT NULL,
		email_lower TEXT NOT NULL,
		email_verified INTEGER DEFAULT 0,
		status TEXT,
		is_platform_admin INTEGER NOT NULL DEFAULT 0
	)`,
	`CREATE UNIQUE INDEX idx_users_email_lower ON users(email_lower)`,
	`CREATE TABLE password_credentials (
		id TEXT PRIMARY KEY,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		user_id TEXT NOT NULL,
		password_hash TEXT NOT NULL
	)`,
	`CREATE UNIQUE INDEX idx_password_credentials_user_id ON password_credentials(user_id)`,
	`CREATE TABLE tenants (
		id TEXT PRIMARY KEY,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		name TEXT,
		domain TEXT
	)`,
	`CREATE TABLE tenant_memberships (
		id TEXT PRIMARY KEY,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		user_id TEXT NOT NULL,
		tenant_id TEXT NOT NULL,
		role TEXT NOT NULL,
		status TEXT NOT NULL
	)`,
	`CREATE UNIQUE INDEX idx_tenant_memberships_user_tenant ON tenant_memberships(user_id, tenant_id)`,
	`CREATE TABLE tenant_identity_providers (
		id TEXT PRIMARY KEY,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		tenant_id TEXT NOT NULL,
		provider TEXT NOT NULL,
		enabled INTEGER NOT NULL DEFAULT 0,
		oauth_client_id TEXT,
		oauth_client_secret_encrypted TEXT
	)`,
	`CREATE UNIQUE INDEX idx_tenant_identity_providers_tenant_provider ON tenant_identity_providers(tenant_id, provider)`,
	`CREATE TABLE clients (
		id TEXT PRIMARY KEY,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		tenant_id TEXT NOT NULL,
		client_id TEXT NOT NULL,
		client_secret TEXT,
		name TEXT,
		redirect_uris TEXT,
		grant_types TEXT,
		scopes TEXT,
		is_public INTEGER DEFAULT 0
	)`,
	`CREATE UNIQUE INDEX idx_clients_tenant_client_id ON clients(tenant_id, client_id)`,
	`CREATE TABLE sessions (
		id TEXT PRIMARY KEY,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		user_id TEXT NOT NULL,
		tenant_id TEXT NOT NULL,
		ip_address TEXT,
		user_agent TEXT,
		expires_at DATETIME
	)`,
	`CREATE TABLE refresh_tokens (
		id TEXT PRIMARY KEY,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		tenant_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		oauth_client_id TEXT NOT NULL,
		token_hash TEXT NOT NULL,
		revoked INTEGER DEFAULT 0,
		expires_at DATETIME NOT NULL,
		client_record_id TEXT
	)`,
	`CREATE UNIQUE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash)`,
	`CREATE TABLE authorization_codes (
		id TEXT PRIMARY KEY,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		code TEXT NOT NULL UNIQUE,
		tenant_id TEXT NOT NULL,
		oauth_client_id TEXT NOT NULL,
		user_id TEXT NOT NULL,
		scope TEXT,
		redirect_uri TEXT,
		code_challenge TEXT,
		code_challenge_method TEXT,
		nonce TEXT,
		expires_at DATETIME NOT NULL,
		client_record_id TEXT
	)`,
	`CREATE TABLE audit_logs (
		id TEXT PRIMARY KEY,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		tenant_id TEXT,
		action TEXT NOT NULL,
		result TEXT NOT NULL,
		actor_type TEXT NOT NULL,
		actor_id TEXT,
		resource_type TEXT,
		resource_id TEXT,
		resource_name TEXT,
		ip_address TEXT,
		user_agent TEXT,
		request_id TEXT,
		correlation_id TEXT,
		old_value TEXT,
		new_value TEXT
	)`,
	`CREATE TABLE federated_identities (
		id TEXT PRIMARY KEY,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		user_id TEXT NOT NULL,
		provider TEXT NOT NULL,
		subject TEXT NOT NULL,
		email_at_link TEXT
	)`,
	`CREATE UNIQUE INDEX idx_federated_identities_provider_subject ON federated_identities(provider, subject)`,
	`CREATE TABLE user_mfa_totps (
		id TEXT PRIMARY KEY,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		user_id TEXT NOT NULL,
		secret_encrypted TEXT NOT NULL,
		enabled INTEGER NOT NULL DEFAULT 0,
		verified_at DATETIME
	)`,
	`CREATE TABLE user_mfa_recovery_codes (
		id TEXT PRIMARY KEY,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		user_id TEXT NOT NULL,
		code_hash TEXT NOT NULL,
		used_at DATETIME
	)`,
	`CREATE TABLE webauthn_credentials (
		id TEXT PRIMARY KEY,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at DATETIME,
		user_id TEXT NOT NULL,
		credential_id TEXT NOT NULL,
		public_key TEXT NOT NULL,
		sign_count INTEGER NOT NULL,
		device_name TEXT
	)`,
	`CREATE UNIQUE INDEX idx_webauthn_credentials_credential_id ON webauthn_credentials(credential_id)`,
}

func migrateTestSchema(gormDB *gorm.DB) error {
	for _, stmt := range sqliteTestSchema {
		if err := gormDB.Exec(stmt).Error; err != nil {
			return err
		}
	}
	return nil
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func newTestDB(t *testing.T) *db.PostgresDB {
	t.Helper()
	gormDB, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)
	require.NoError(t, migrateTestSchema(gormDB))
	return &db.PostgresDB{DB: gormDB}
}

func closedTestDB(t *testing.T) *db.PostgresDB {
	t.Helper()
	pg := newTestDB(t)
	sqlDB, err := pg.DB.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())
	return pg
}

func testCtx() context.Context {
	return context.Background()
}

func requireNotFound(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, errors.ErrorTypeNotFound, appErr.Type)
}

func requireDatabaseErr(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	appErr := errors.GetAppError(err)
	require.NotNil(t, appErr)
	require.Equal(t, errors.ErrorTypeDatabase, appErr.Type)
}

func requireDatabaseOrClosedErr(t *testing.T, err error) {
	t.Helper()
	require.Error(t, err)
	if appErr := errors.GetAppError(err); appErr != nil {
		require.Equal(t, errors.ErrorTypeDatabase, appErr.Type)
	}
}

func seedTenant(t *testing.T, pg *db.PostgresDB, name, domain string) *models.Tenant {
	t.Helper()
	tenant := &models.Tenant{
		BaseModel: models.NewBaseModel(),
		Name:      name,
		Domain:    domain,
	}
	require.NoError(t, pg.WithContext(testCtx()).Create(tenant).Error)
	return tenant
}

func seedUser(t *testing.T, pg *db.PostgresDB, email string) *models.User {
	t.Helper()
	user := &models.User{
		BaseModel:  models.NewBaseModel(),
		FirstName:  "Test",
		LastName:   "User",
		Email:      email,
		Status:     constants.UserStatusActive,
	}
	require.NoError(t, pg.WithContext(testCtx()).Create(user).Error)
	return user
}

func seedUserWithPassword(t *testing.T, pg *db.PostgresDB, email, hash string) *models.User {
	t.Helper()
	repo := ProvideUserRepository(pg)
	user := &models.User{
		BaseModel: models.NewBaseModel(),
		FirstName: "Test",
		LastName:  "User",
		Email:     email,
		Status:    constants.UserStatusActive,
	}
	require.NoError(t, repo.CreateWithPasswordHash(testCtx(), user, hash))
	return user
}

func seedMembership(t *testing.T, pg *db.PostgresDB, userID, tenantID string) *models.TenantMembership {
	t.Helper()
	m := &models.TenantMembership{
		BaseModel: models.NewBaseModel(),
		UserID:    userID,
		TenantID:  tenantID,
		Role:      constants.TenantMembershipRoleMember,
		Status:    constants.TenantMembershipStatusActive,
	}
	require.NoError(t, pg.WithContext(testCtx()).Create(m).Error)
	return m
}

func seedClient(t *testing.T, pg *db.PostgresDB, tenantID, clientID string) *models.Client {
	t.Helper()
	client := &models.Client{
		BaseModel: models.NewBaseModel(),
		TenantID:  tenantID,
		ClientID:  clientID,
		Name:      "Test Client",
		IsPublic:  true,
	}
	require.NoError(t, pg.WithContext(testCtx()).Create(client).Error)
	return client
}

func futureTime() time.Time {
	return time.Now().UTC().Add(time.Hour)
}

func pastTime() time.Time {
	return time.Now().UTC().Add(-time.Hour)
}
