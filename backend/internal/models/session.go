package models

import "time"

// Session represents an authenticated browser or API session.
type Session struct {
	BaseModel
	UserID    string     `gorm:"column:user_id;type:uuid;not null;index"`
	TenantID  string     `gorm:"column:tenant_id;type:uuid;not null;index"`
	IPAddress string     `gorm:"column:ip_address;type:varchar(50)"`
	UserAgent string     `gorm:"column:user_agent;type:text"`
	ExpiresAt *time.Time `gorm:"column:expires_at;type:timestamptz"`

	User   *User   `gorm:"foreignKey:UserID"`
	Tenant *Tenant `gorm:"foreignKey:TenantID"`
}

func (Session) TableName() string {
	return "sessions"
}
