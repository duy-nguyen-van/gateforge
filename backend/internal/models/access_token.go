package models

import "time"

// AccessToken stores metadata for opaque (reference) OAuth2 access tokens.
// TokenHash is a one-way hash (e.g. SHA-256 hex) of the issued token — never store raw tokens.
// OAuthClientID is the public client_id string; ClientRecordID optionally references clients.id.
type AccessToken struct {
	HardDeleteModel
	TenantID       string     `gorm:"column:tenant_id;type:uuid;not null;index"`
	UserID         *string    `gorm:"column:user_id;type:uuid;index"`
	OAuthClientID  *string    `gorm:"column:oauth_client_id;type:varchar(255);index"`
	TokenHash      string     `gorm:"column:token_hash;type:text;not null;uniqueIndex"`
	ExpiresAt      *time.Time `gorm:"column:expires_at;type:timestamptz"`
	ClientRecordID *string    `gorm:"column:client_record_id;type:uuid;index"`

	User   *User   `gorm:"foreignKey:UserID"`
	Tenant *Tenant `gorm:"foreignKey:TenantID"`
	Client *Client `gorm:"foreignKey:ClientRecordID"`
}

func (AccessToken) TableName() string {
	return "access_tokens"
}
