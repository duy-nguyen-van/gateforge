package models

import "time"

// UserMFARecoveryCode is a hashed one-time recovery code.
type UserMFARecoveryCode struct {
	BaseModel
	UserID   string     `gorm:"column:user_id;type:uuid;not null;index"`
	CodeHash string     `gorm:"column:code_hash;type:varchar(128);not null"`
	UsedAt   *time.Time `gorm:"column:used_at;type:timestamptz"`

	User *User `gorm:"foreignKey:UserID"`
}

func (UserMFARecoveryCode) TableName() string {
	return "user_mfa_recovery_codes"
}
