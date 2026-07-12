package models

import "github.com/gateforge-iam/gateforge-iam/internal/constants"

// TenantMembership links a global user to a tenant with role and status.
type TenantMembership struct {
	BaseModel
	UserID   string                           `gorm:"column:user_id;type:uuid;not null;uniqueIndex:idx_tenant_memberships_user_tenant"`
	TenantID string                           `gorm:"column:tenant_id;type:uuid;not null;uniqueIndex:idx_tenant_memberships_user_tenant"`
	Role     constants.TenantMembershipRole   `gorm:"column:role;type:varchar(32);not null"`
	Status   constants.TenantMembershipStatus `gorm:"column:status;type:varchar(32);not null"`

	User   *User   `gorm:"foreignKey:UserID"`
	Tenant *Tenant `gorm:"foreignKey:TenantID"`
}

func (TenantMembership) TableName() string {
	return "tenant_memberships"
}
