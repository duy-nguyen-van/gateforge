package models

import "time"

// TenantInvite holds an organization seat for an email that may not have an account yet.
type TenantInvite struct {
	BaseModel
	Email          string    `gorm:"column:email;not null"`
	EmailLower     string    `gorm:"column:email_lower;not null"`
	TenantID       string    `gorm:"column:tenant_id;type:uuid;not null"`
	Role           string    `gorm:"column:role;type:varchar(32);not null"`
	TokenHash      string    `gorm:"column:token_hash;not null"`
	Status         string    `gorm:"column:status;type:varchar(32);not null"`
	ExpiresAt      time.Time `gorm:"column:expires_at;not null"`
	AcceptedUserID *string   `gorm:"column:accepted_user_id;type:uuid"`

	Tenant *Tenant `gorm:"foreignKey:TenantID"`
}

func (TenantInvite) TableName() string {
	return "tenant_invites"
}
