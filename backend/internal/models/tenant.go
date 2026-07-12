package models

// Tenant represents a tenant (organization) in the system.
type Tenant struct {
	BaseModel
	Name   string `gorm:"column:name;type:varchar(255)"`
	Domain string `gorm:"column:domain;type:varchar(255)"`
}

func (Tenant) TableName() string {
	return "tenants"
}
