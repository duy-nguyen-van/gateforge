package constants

type AuditResult string

const (
	AuditResultSuccess AuditResult = "success"
	AuditResultFailure AuditResult = "failure"
	AuditResultDenied  AuditResult = "denied"
)

type AuditActorType string

const (
	AuditActorTypeUser        AuditActorType = "user"
	AuditActorTypeSystem      AuditActorType = "system"
	AuditActorTypeOAuthClient AuditActorType = "oauth_client"
	AuditActorTypeAdminAPIKey AuditActorType = "admin_api_key"
)

type AuditResourceType string

const (
	AuditResourceTypeUser               AuditResourceType = "user"
	AuditResourceTypeTenant             AuditResourceType = "tenant"
	AuditResourceTypeClient             AuditResourceType = "client"
	AuditResourceTypeMembership         AuditResourceType = "membership"
	AuditResourceTypeIdentityProvider   AuditResourceType = "identity_provider"
	AuditResourceTypeSession            AuditResourceType = "session"
	AuditResourceTypeWebauthnCredential AuditResourceType = "webauthn_credential"
)

// Audit actions (dotted namespace).
const (
	AuditActionAuthRegister          = "auth.register"
	AuditActionAuthLogin             = "auth.login"
	AuditActionAuthLogout            = "auth.logout"
	AuditActionAuthRefresh           = "auth.refresh"
	AuditActionTenantSelect          = "tenant.select"
	AuditActionTenantSwitch          = "tenant.switch"
	AuditActionOIDCLogin             = "oidc.login"
	AuditActionOIDCAuthorize         = "oidc.authorize"
	AuditActionOIDCTokenIssue        = "oidc.token.issue"
	AuditActionFederationStart       = "federation.start"
	AuditActionFederationLogin       = "federation.login"
	AuditActionMFATOTPSetup          = "mfa.totp.setup"
	AuditActionMFATOTPEnable         = "mfa.totp.enable"
	AuditActionMFARecoveryRegenerate = "mfa.recovery_codes.regenerate"
	AuditActionMFAChallengeVerify    = "mfa.challenge.verify"
	AuditActionWebauthnRegister      = "webauthn.register"
	AuditActionWebauthnLogin         = "webauthn.login"
	AuditActionSessionCreate         = "session.create"
	AuditActionSessionRevokeAll      = "session.revoke_all"
	AuditActionAdminMemberAdd        = "admin.member.add"
	AuditActionAdminMemberRemove     = "admin.member.remove"
	AuditActionAdminTenantCreate     = "admin.tenant.create"
	AuditActionAdminTenantUpdate     = "admin.tenant.update"
	AuditActionAdminTenantDelete     = "admin.tenant.delete"
	AuditActionAdminClientCreate     = "admin.client.create"
	AuditActionAdminClientUpdate     = "admin.client.update"
	AuditActionAdminClientDelete     = "admin.client.delete"
	AuditActionAdminIDPPatch         = "admin.idp.patch"
	AuditActionAdminUserDisable      = "admin.user.disable"
	AuditActionAdminUserForceLogout  = "admin.user.force_logout"
	AuditActionAdminPasskeyReset     = "admin.passkey.reset"
	AuditActionAdminMFAReset         = "admin.mfa.reset"
)
