package models

// PasswordCredential stores a password hash for local email/password authentication.
type PasswordCredential struct {
	BaseModel
	UserID       string `gorm:"column:user_id;type:uuid;not null;uniqueIndex"`
	PasswordHash string `gorm:"column:password_hash;not null"`
	User         *User  `gorm:"foreignKey:UserID"`
}

func (PasswordCredential) TableName() string {
	return "password_credentials"
}
