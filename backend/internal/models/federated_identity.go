package models

// FederatedIdentity links an upstream IdP subject to a global user.
type FederatedIdentity struct {
	BaseModel
	UserID      string `gorm:"column:user_id;type:uuid;not null;index"`
	Provider    string `gorm:"column:provider;type:varchar(32);not null;uniqueIndex:idx_federated_identities_provider_subject"`
	Subject     string `gorm:"column:subject;type:text;not null;uniqueIndex:idx_federated_identities_provider_subject"`
	EmailAtLink string `gorm:"column:email_at_link;type:text"`

	User *User `gorm:"foreignKey:UserID"`
}

func (FederatedIdentity) TableName() string {
	return "federated_identities"
}
