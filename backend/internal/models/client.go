package models

import (
	"github.com/lib/pq"
)

// Client is an OAuth2 / OpenID client registered for a tenant.
// ClientID is the OAuth2 client_id string exposed to clients; BaseModel.ID is the internal primary key.
type Client struct {
	BaseModel
	TenantID     string         `gorm:"column:tenant_id;type:uuid;not null;uniqueIndex:idx_clients_tenant_client_id"`
	ClientID     string         `gorm:"column:client_id;type:varchar(255);not null;uniqueIndex:idx_clients_tenant_client_id"`
	ClientSecret string         `gorm:"column:client_secret;type:varchar(255)"`
	Name         string         `gorm:"column:name;type:varchar(255)"`
	RedirectUris pq.StringArray `gorm:"column:redirect_uris;type:text[]"`
	GrantTypes   pq.StringArray `gorm:"column:grant_types;type:text[]"`
	Scopes       pq.StringArray `gorm:"column:scopes;type:text[]"`
	IsPublic     bool           `gorm:"column:is_public;default:false"`

	Tenant *Tenant `gorm:"foreignKey:TenantID"`
}

func (Client) TableName() string {
	return "clients"
}
