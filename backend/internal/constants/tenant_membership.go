package constants

type TenantMembershipRole string

const (
	TenantMembershipRoleMember TenantMembershipRole = "member"
	TenantMembershipRoleAdmin  TenantMembershipRole = "admin"
)

type TenantMembershipStatus string

const (
	TenantMembershipStatusActive    TenantMembershipStatus = "active"
	TenantMembershipStatusInvited   TenantMembershipStatus = "invited"
	TenantMembershipStatusSuspended TenantMembershipStatus = "suspended"
)
