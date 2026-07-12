package models

// TenantIdentityProvider toggles an external login provider per tenant.
type TenantIdentityProvider struct {
	BaseModel
	TenantID                   string `gorm:"column:tenant_id;type:uuid;not null;uniqueIndex:idx_tenant_identity_providers_tenant_provider"`
	Provider                   string `gorm:"column:provider;type:varchar(32);not null;uniqueIndex:idx_tenant_identity_providers_tenant_provider"`
	Enabled                    bool   `gorm:"column:enabled;not null"`
	OAuthClientID              string `gorm:"column:oauth_client_id;type:varchar(255)"`
	OAuthClientSecretEncrypted string `gorm:"column:oauth_client_secret_encrypted;type:text"`

	Tenant *Tenant `gorm:"foreignKey:TenantID"`
}

func (TenantIdentityProvider) TableName() string {
	return "tenant_identity_providers"
}
