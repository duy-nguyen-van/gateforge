package models

import "time"

// UserMFATOTP stores encrypted TOTP secret per user (at most one active row).
type UserMFATOTP struct {
	BaseModel
	UserID          string     `gorm:"column:user_id;type:uuid;not null;index"`
	SecretEncrypted string     `gorm:"column:secret_encrypted;type:text;not null"`
	Enabled         bool       `gorm:"column:enabled;not null;default:false"`
	VerifiedAt      *time.Time `gorm:"column:verified_at;type:timestamptz"`

	User *User `gorm:"foreignKey:UserID"`
}

func (UserMFATOTP) TableName() string {
	return "user_mfa_totps"
}
