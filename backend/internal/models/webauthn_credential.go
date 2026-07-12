package models

// WebauthnCredential stores a registered WebAuthn passkey / credential for a user.
type WebauthnCredential struct {
	BaseModel
	UserID       string `gorm:"column:user_id;type:uuid;not null;index"`
	CredentialID string `gorm:"column:credential_id;type:text;not null;uniqueIndex"`
	PublicKey    string `gorm:"column:public_key;type:text;not null"`
	SignCount    int64  `gorm:"column:sign_count;not null"`
	DeviceName   string `gorm:"column:device_name;type:varchar(255)"`

	User *User `gorm:"foreignKey:UserID"`
}

func (WebauthnCredential) TableName() string {
	return "webauthn_credentials"
}
