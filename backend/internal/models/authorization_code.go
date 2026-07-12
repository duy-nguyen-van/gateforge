package models

import "time"

// AuthorizationCode is an OAuth2 authorization code (PKCE optional).
// OAuthClientID is the public OAuth2 client_id string.
type AuthorizationCode struct {
	BaseModel
	Code                string    `gorm:"column:code;type:varchar(255);primaryKey"`
	TenantID            string    `gorm:"column:tenant_id;type:uuid;not null;index"`
	OAuthClientID       string    `gorm:"column:oauth_client_id;type:varchar(255);not null;index"`
	UserID              string    `gorm:"column:user_id;type:uuid;not null;index"`
	Scope               string    `gorm:"column:scope;type:text"`
	RedirectURI         string    `gorm:"column:redirect_uri;type:text"`
	CodeChallenge       string    `gorm:"column:code_challenge;type:text"`
	CodeChallengeMethod string    `gorm:"column:code_challenge_method;type:varchar(10)"`
	Nonce               string    `gorm:"column:nonce;type:text"`
	ExpiresAt           time.Time `gorm:"column:expires_at;type:timestamptz;not null"`
	ClientRecordID      *string   `gorm:"column:client_record_id;type:uuid;index"`

	Tenant *Tenant `gorm:"foreignKey:TenantID"`
	User   *User   `gorm:"foreignKey:UserID"`
	Client *Client `gorm:"foreignKey:ClientRecordID"`
}

func (AuthorizationCode) TableName() string {
	return "authorization_codes"
}
