package models

import "time"

// RefreshToken stores OAuth2 refresh tokens.
// TokenHash is a one-way hash of the issued refresh token value.
// OAuthClientID is the public OAuth2 client_id; ClientRecordID optionally references clients.id.
type RefreshToken struct {
	HardDeleteModel
	TenantID       string    `gorm:"column:tenant_id;type:uuid;not null;index"`
	UserID         string    `gorm:"column:user_id;type:uuid;not null;index:idx_refresh_tokens_user_oauth_client_revoked"`
	OAuthClientID  string    `gorm:"column:oauth_client_id;type:varchar(255);not null;index:idx_refresh_tokens_user_oauth_client_revoked"`
	TokenHash      string    `gorm:"column:token_hash;type:text;not null;uniqueIndex"`
	Revoked        bool      `gorm:"column:revoked;default:false;index:idx_refresh_tokens_user_oauth_client_revoked"`
	ExpiresAt      time.Time `gorm:"column:expires_at;type:timestamptz;not null"`
	ClientRecordID *string   `gorm:"column:client_record_id;type:uuid;index"`

	User   *User   `gorm:"foreignKey:UserID"`
	Tenant *Tenant `gorm:"foreignKey:TenantID"`
	Client *Client `gorm:"foreignKey:ClientRecordID"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
