import type {
  AuthenticationResponseJSON,
  PublicKeyCredentialCreationOptionsJSON,
  PublicKeyCredentialRequestOptionsJSON,
  RegistrationResponseJSON,
} from '@simplewebauthn/browser'

/** API types aligned with iam-backend docs/swagger.yaml */

export interface ApiMeta {
  error_code?: string
  message?: string
  code?: number
  page?: number
  page_size?: number
  total?: number
}

export interface ApiEnvelope<T> {
  meta: ApiMeta
  data: T
}

export interface RegisterRequest {
  email: string
  password: string
  first_name?: string
  last_name?: string
}

export interface LoginRequest {
  email: string
  password: string
  remember_me?: boolean
  return_to?: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
  refresh_expires_in: number
  active_tenant_id?: string
}

export interface TenantSummary {
  id: string
  name: string
  domain: string
  role: string
}

export interface TenantSelectionResponse {
  selection_required: boolean
  tenants: TenantSummary[]
  selection_token: string
  expires_in: number
}

export interface TenantSelectRequest {
  selection_token: string
  tenant_id: string
  remember_me?: boolean
}

export interface TenantSwitchRequest {
  tenant_id: string
}

export interface RefreshTokenRequest {
  refresh_token: string
}

export interface MFALoginChallengeResponse {
  mfa_required: boolean
  mfa_ticket: string
  expires_in: number
}

export interface MFAChallengeVerifyRequest {
  mfa_ticket: string
  code: string
}

export interface MFATOTPSetupResponse {
  secret: string
  otpauth_uri: string
}

export interface MFATOTPVerifyRequest {
  code: string
}

export interface MFARecoveryCodesResponse {
  codes: string[]
}

export interface UserResponse {
  id: string
  email: string
  first_name: string
  last_name: string
  email_verified: boolean
  is_platform_admin: boolean
  mfa_enabled: boolean
  active_tenant_id?: string
  tenants?: TenantSummary[]
  created_at: string
  updated_at: string
}

export interface UpdateProfileRequest {
  first_name?: string
  last_name?: string
}

export interface WebauthnCredentialResponse {
  id: string
  device_name: string
  created_at: string
}

export interface WebauthnRegisterStartRequest {
  device_name?: string
}

export interface WebauthnRegisterStartResponse {
  options: PublicKeyCredentialCreationOptionsJSON
  session_token: string
}

export interface WebauthnRegisterFinishRequest {
  session_token: string
  credential: RegistrationResponseJSON
}

export interface WebauthnLoginStartRequest {
  email: string
  tenant_id?: string
}

export interface WebauthnLoginStartResponse {
  options: PublicKeyCredentialRequestOptionsJSON
  session_token: string
}

export interface WebauthnLoginFinishRequest {
  email: string
  tenant_id?: string
  session_token: string
  credential: AuthenticationResponseJSON
  remember_me?: boolean
  return_to?: string
}

export type LoginResult = LoginResponse | MFALoginChallengeResponse | TenantSelectionResponse

export interface AdminStatsResponse {
  total_users: number
  mfa_enabled_count: number
  mfa_enabled_percent: number
  active_sessions: number
}

export interface AdminUserResponse {
  id: string
  email: string
  first_name: string
  last_name: string
  status: string
  tenant_id: string
  mfa_enabled: boolean
  created_at: string
}

export interface AdminUserMembershipResponse {
  tenant_id: string
  tenant_name: string
  role: string
  status: string
}

export interface AdminUserDetailResponse {
  id: string
  email: string
  first_name: string
  last_name: string
  status: string
  mfa_enabled: boolean
  is_platform_admin: boolean
  passkey_count: number
  active_sessions: number
  memberships: AdminUserMembershipResponse[]
  created_at: string
}

export interface AdminClientUsageResponse {
  client_id: string
  total_refresh_tokens: number
  active_refresh_tokens: number
  last_token_issued_at?: string | null
  authorize_events_30d: number
  token_issue_events_30d: number
}

export interface AdminTenantResponse {
  id: string
  name: string
  domain: string
  user_count: number
  created_at: string
}

export interface AdminCreateTenantRequest {
  name: string
  domain?: string
}

export interface AdminUpdateTenantRequest {
  name?: string
  domain?: string
}

export interface AdminTenantMemberResponse {
  user_id: string
  email: string
  first_name: string
  last_name: string
  role: string
  status: string
  joined_at: string
}

export interface AdminClientResponse {
  id: string
  tenant_id: string
  client_id: string
  name: string
  is_public: boolean
  grant_types: string[]
  redirect_uris: string[]
  scopes: string[]
  client_secret_set: boolean
  created_at: string
}

export interface AdminCreateClientRequest {
  tenant_id: string
  client_id?: string
  name: string
  is_public: boolean
  redirect_uris: string[]
  grant_types?: string[]
  scopes?: string[]
}

export interface AdminUpdateClientRequest {
  name?: string
  redirect_uris?: string[]
  grant_types?: string[]
  scopes?: string[]
  is_public?: boolean
  client_secret?: string
}

export interface AdminCreateClientResponse extends AdminClientResponse {
  client_secret?: string
}

export interface AdminIdentityProviderResponse {
  provider: string
  name: string
  enabled: boolean
  configured: boolean
  tenant_id: string
  oauth_client_id?: string
  oauth_client_secret_set: boolean
  redirect_uri?: string
  setup_console_url?: string
}

export interface PatchIdentityProviderRequest {
  enabled?: boolean
  oauth_client_id?: string
  oauth_client_secret?: string
}

export interface PublicFederationProviderResponse {
  provider: string
  name: string
}

export interface AdminListParams {
  page?: number
  page_size?: number
  tenant_id?: string
  search?: string
}

export interface PaginationParams {
  page?: number
  page_size?: number
}

export interface AdminLoginHistoryListParams {
  page?: number
  page_size?: number
  tenant_id?: string
  result?: string
  actor_id?: string
}

export interface AdminAuditLogListParams {
  page?: number
  page_size?: number
  tenant_id?: string
  action?: string
  result?: string
  actor_id?: string
}

export interface AdminAuditLogResponse {
  id: string
  created_at: string
  tenant_id?: string | null
  action: string
  result: string
  actor_type: string
  actor_id?: string | null
  resource_type?: string | null
  resource_id?: string | null
  resource_name?: string | null
  ip_address?: string | null
  user_agent?: string | null
  request_id?: string | null
  correlation_id?: string | null
  old_value?: Record<string, unknown> | null
  new_value?: Record<string, unknown> | null
}

export interface AdminAddMemberRequest {
  email: string
  role?: 'member' | 'admin'
}

export function isMfaChallenge(data: LoginResult): data is MFALoginChallengeResponse {
  return 'mfa_required' in data && data.mfa_required === true
}

export function isTenantSelection(data: LoginResult): data is TenantSelectionResponse {
  return 'selection_required' in data && data.selection_required === true
}

export class ApiError extends Error {
  readonly status: number
  readonly errorCode?: string

  constructor(message: string, status: number, errorCode?: string) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.errorCode = errorCode
  }
}
