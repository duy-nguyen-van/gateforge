package models

import (
	"github.com/lib/pq"
)

// Consent records user consent to client scopes (OAuth).
// OAuthClientID is the public OAuth2 client_id string (matches clients.client_id).
type Consent struct {
	BaseModel
	TenantID       string         `gorm:"column:tenant_id;type:uuid;not null;uniqueIndex:idx_consents_tenant_user_oauth_client"`
	UserID         string         `gorm:"column:user_id;type:uuid;not null;uniqueIndex:idx_consents_tenant_user_oauth_client"`
	OAuthClientID  string         `gorm:"column:oauth_client_id;type:varchar(255);not null;uniqueIndex:idx_consents_tenant_user_oauth_client"`
	Scopes         pq.StringArray `gorm:"column:scopes;type:text[]"`
	Granted        bool           `gorm:"column:granted;default:true"`
	ClientRecordID *string        `gorm:"column:client_record_id;type:uuid;index"` // optional FK to clients.id

	User   *User   `gorm:"foreignKey:UserID"`
	Tenant *Tenant `gorm:"foreignKey:TenantID"`
	Client *Client `gorm:"foreignKey:ClientRecordID"`
}

func (Consent) TableName() string {
	return "consents"
}
