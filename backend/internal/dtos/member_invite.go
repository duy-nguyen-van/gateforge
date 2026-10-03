package dtos

// MemberInvitePreview is the public view of an organization invite.
type MemberInvitePreview struct {
	Email            string `json:"email"`
	OrganizationName string `json:"organization_name"`
	Role             string `json:"role"`
}

// AcceptMemberInviteRequest creates an account when needed and signs the invited person in.
type AcceptMemberInviteRequest struct {
	Token     string `json:"token" validate:"required"`
	Password  string `json:"password" validate:"required,min=12,max=128"`
	FirstName string `json:"first_name,omitempty" validate:"omitempty,max=100"`
	LastName  string `json:"last_name,omitempty" validate:"omitempty,max=100"`
}
