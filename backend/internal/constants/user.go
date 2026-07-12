package constants

type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusDisabled UserStatus = "disabled"
	UserStatusLocked   UserStatus = "locked"
)
