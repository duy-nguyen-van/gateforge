package models

import (
	"strings"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"

	"gorm.io/gorm"
)

// User represents a global identity (one row per email).
type User struct {
	BaseModel
	FirstName       string               `gorm:"column:first_name"`
	LastName        string               `gorm:"column:last_name"`
	Email           string               `gorm:"column:email;not null"`
	EmailLower      string               `gorm:"column:email_lower;not null;uniqueIndex:idx_users_email_lower"`
	EmailVerified   bool                 `gorm:"column:email_verified"`
	Status          constants.UserStatus `gorm:"column:status"`
	IsPlatformAdmin bool                 `gorm:"column:is_platform_admin;not null;default:false"`

	Memberships []TenantMembership `gorm:"foreignKey:UserID"`

	// Optional authenticators (use what your product enables).
	PasswordCredential   *PasswordCredential   `gorm:"foreignKey:UserID"`
	WebauthnCredentials  []WebauthnCredential  `gorm:"foreignKey:UserID"`
	UserMFATOTP          *UserMFATOTP          `gorm:"foreignKey:UserID"`
	UserMFARecoveryCodes []UserMFARecoveryCode `gorm:"foreignKey:UserID"`
}

// BeforeSave normalizes EmailLower for case-insensitive global uniqueness.
func (u *User) BeforeSave(tx *gorm.DB) error {
	u.EmailLower = strings.ToLower(strings.TrimSpace(u.Email))
	return nil
}

func (User) TableName() string {
	return "users"
}
