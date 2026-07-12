import type {
  AuthenticationResponseJSON,
  PublicKeyCredentialCreationOptionsJSON,
  PublicKeyCredentialRequestOptionsJSON,
  RegistrationResponseJSON,
} from '@simplewebauthn/browser'

import { apiUrl } from '@/lib/utils'
import {
  ApiError,
  type ApiEnvelope,
  type LoginRequest,
  type LoginResponse,
  type LoginResult,
  type MFAChallengeVerifyRequest,
  type MFARecoveryCodesResponse,
  type MFATOTPSetupResponse,
  type MFATOTPVerifyRequest,
  type RefreshTokenRequest,
  type RegisterRequest,
  type UpdateProfileRequest,
  type UserResponse,
  type WebauthnCredentialResponse,
  type WebauthnLoginFinishRequest,
  type WebauthnLoginStartRequest,
  type WebauthnLoginStartResponse,
  type WebauthnRegisterFinishRequest,
  type WebauthnRegisterStartRequest,
  type WebauthnRegisterStartResponse,
  type AdminStatsResponse,
  type AdminUserResponse,
  type AdminTenantResponse,
  type AdminCreateTenantRequest,
  type AdminUpdateTenantRequest,
  type AdminTenantMemberResponse,
  type AdminClientResponse,
  type AdminCreateClientRequest,
  type AdminCreateClientResponse,
  type AdminUpdateClientRequest,
  type AdminIdentityProviderResponse,
  type PatchIdentityProviderRequest,
  type PublicFederationProviderResponse,
  type TenantSelectRequest,
  type TenantSummary,
  type TenantSwitchRequest,
  type AdminListParams,
  type AdminAuditLogListParams,
  type AdminLoginHistoryListParams,
  type AdminAuditLogResponse,
  type AdminUserDetailResponse,
  type AdminClientUsageResponse,
  type AdminAddMemberRequest,
  type PaginationParams,
} from '@/api/types'
import { getAccessToken, refreshTokens } from '@/auth/token-store'

export type RequestOptions = {
  auth?: boolean
  credentials?: RequestCredentials
  csrfToken?: string
  redirect?: RequestRedirect
}

async function parseEnvelope<T>(response: Response): Promise<ApiEnvelope<T>> {
  const body = (await response.json()) as ApiEnvelope<T>
  if (!response.ok) {
    throw new ApiError(
      body.meta?.message ?? 'Request failed',
      response.status,
      body.meta?.error_code,
    )
  }
  return body
}

export async function apiFetch<T>(
  path: string,
  init: RequestInit = {},
  options: RequestOptions = {},
): Promise<ApiEnvelope<T>> {
  const headers = new Headers(init.headers)
  if (!headers.has('Content-Type') && init.body) {
    headers.set('Content-Type', 'application/json')
  }

  if (options.auth) {
    const token = getAccessToken()
    if (token) {
      headers.set('Authorization', `Bearer ${token}`)
    }
  }

  if (options.csrfToken) {
    headers.set('X-CSRF-Token', options.csrfToken)
  }

  const doFetch = async (): Promise<Response> =>
    fetch(apiUrl(path), {
      ...init,
      headers,
      credentials: options.credentials ?? 'include',
      redirect: options.redirect ?? 'follow',
    })

  let response = await doFetch()

  if (response.status === 401 && options.auth) {
    const refreshed = await refreshTokens()
    if (refreshed) {
      headers.set('Authorization', `Bearer ${getAccessToken()}`)
      response = await doFetch()
    }
  }

  return parseEnvelope<T>(response)
}

async function parseErrorResponse(response: Response): Promise<ApiError> {
  let message = 'Request failed'
  let errorCode: string | undefined
  try {
    const body = (await response.json()) as ApiEnvelope<unknown>
    message = body.meta?.message ?? message
    errorCode = body.meta?.error_code
  } catch {
    // ignore parse errors
  }
  return new ApiError(message, response.status, errorCode)
}

export async function apiFetchNoContent(
  path: string,
  init: RequestInit = {},
  options: RequestOptions = {},
): Promise<void> {
  const headers = new Headers(init.headers)
  if (!headers.has('Content-Type') && init.body) {
    headers.set('Content-Type', 'application/json')
  }

  if (options.auth) {
    const token = getAccessToken()
    if (token) {
      headers.set('Authorization', `Bearer ${token}`)
    }
  }

  if (options.csrfToken) {
    headers.set('X-CSRF-Token', options.csrfToken)
  }

  const doFetch = async (): Promise<Response> =>
    fetch(apiUrl(path), {
      ...init,
      headers,
      credentials: options.credentials ?? 'include',
      redirect: options.redirect ?? 'follow',
    })

  let response = await doFetch()

  if (response.status === 401 && options.auth) {
    const refreshed = await refreshTokens()
    if (refreshed) {
      headers.set('Authorization', `Bearer ${getAccessToken()}`)
      response = await doFetch()
    }
  }

  if (!response.ok) {
    throw await parseErrorResponse(response)
  }
}

export async function prefetchCsrfToken(): Promise<string | undefined> {
  const response = await fetch(apiUrl('/.well-known/openid-configuration'), {
    credentials: 'include',
  })
  return response.headers.get('X-CSRF-Token') ?? undefined
}

export async function registerUser(body: RegisterRequest) {
  return apiFetch<UserResponse>('/api/v1/register', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export async function loginUser(body: LoginRequest) {
  return apiFetch<LoginResult>('/api/v1/login', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export async function refreshUser(body: RefreshTokenRequest) {
  return apiFetch<LoginResponse>('/api/v1/refresh', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export async function logoutUser() {
  return apiFetch<Record<string, never>>(
    '/api/v1/logout',
    { method: 'POST', body: JSON.stringify({}) },
    { auth: true },
  )
}

export async function getMe() {
  return apiFetch<UserResponse>('/api/v1/me', { method: 'GET' }, { auth: true })
}

export async function updateProfile(body: UpdateProfileRequest) {
  return apiFetch<UserResponse>(
    '/api/v1/me',
    { method: 'PATCH', body: JSON.stringify(body) },
    { auth: true },
  )
}

export async function listMyTenants(params: PaginationParams = {}) {
  return apiFetch<TenantSummary[]>(
    `/api/v1/me/tenants${buildPaginationQuery(params)}`,
    { method: 'GET' },
    { auth: true },
  )
}

export async function selectTenant(body: TenantSelectRequest) {
  return apiFetch<LoginResponse>('/api/v1/tenants/select', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export async function switchTenant(body: TenantSwitchRequest) {
  return apiFetch<LoginResponse>(
    '/api/v1/tenants/switch',
    { method: 'POST', body: JSON.stringify(body) },
    { auth: true },
  )
}

export async function verifyMfaChallenge(body: MFAChallengeVerifyRequest) {
  return apiFetch<LoginResponse>('/api/v1/mfa/challenge/verify', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export async function setupTotp() {
  return apiFetch<MFATOTPSetupResponse>(
    '/api/v1/mfa/totp/setup',
    { method: 'POST', body: JSON.stringify({}) },
    { auth: true },
  )
}

export async function verifyTotpEnrollment(body: MFATOTPVerifyRequest) {
  return apiFetch<Record<string, never>>('/api/v1/mfa/totp/verify', {
    method: 'POST',
    body: JSON.stringify(body),
  }, { auth: true })
}

export async function generateRecoveryCodes() {
  return apiFetch<MFARecoveryCodesResponse>(
    '/api/v1/mfa/recovery-codes',
    { method: 'POST', body: JSON.stringify({}) },
    { auth: true },
  )
}

export async function oidcLogin(body: LoginRequest, csrfToken: string) {
  const response = await fetch(apiUrl('/oidc/login'), {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-CSRF-Token': csrfToken,
    },
    credentials: 'include',
    redirect: 'manual',
    body: JSON.stringify(body),
  })

  if (response.status >= 300 && response.status < 400) {
    const location = response.headers.get('Location')
    if (location) {
      window.location.href = location
      return
    }
  }

  if (!response.ok) {
    let message = 'OIDC login failed'
    try {
      const bodyJson = (await response.json()) as ApiEnvelope<unknown>
      message = bodyJson.meta?.message ?? message
    } catch {
      // ignore parse errors
    }
    throw new ApiError(message, response.status)
  }
}

export function federationStartUrl(provider: string, returnTo: string) {
  const params = new URLSearchParams({ return_to: returnTo })
  return apiUrl(`/oidc/federation/${encodeURIComponent(provider)}/start?${params.toString()}`)
}

export function federationCompleteReturnTo() {
  return `${window.location.origin}/login/federation/complete`
}

export async function exchangeSessionLogin() {
  return apiFetch<LoginResult>('/api/v1/login/session', { method: 'POST' }, { credentials: 'include' })
}

export async function listWebauthnCredentials(params: PaginationParams = {}) {
  return apiFetch<WebauthnCredentialResponse[]>(
    `/api/v1/webauthn/credentials${buildPaginationQuery(params)}`,
    { method: 'GET' },
    { auth: true },
  )
}

export async function webauthnRegisterStart(body: WebauthnRegisterStartRequest) {
  return apiFetch<WebauthnRegisterStartResponse>(
    '/api/v1/webauthn/register/start',
    { method: 'POST', body: JSON.stringify(body) },
    { auth: true },
  )
}

export async function webauthnRegisterFinish(body: WebauthnRegisterFinishRequest) {
  return apiFetch<Record<string, never>>(
    '/api/v1/webauthn/register/finish',
    { method: 'POST', body: JSON.stringify(body) },
    { auth: true },
  )
}

export async function webauthnLoginStart(body: WebauthnLoginStartRequest) {
  return apiFetch<WebauthnLoginStartResponse>('/api/v1/webauthn/login/start', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

export async function webauthnLoginFinish(body: WebauthnLoginFinishRequest) {
  return apiFetch<LoginResult>('/api/v1/webauthn/login/finish', {
    method: 'POST',
    body: JSON.stringify(body),
  })
}

function buildPaginationQuery(params: PaginationParams): string {
  const search = new URLSearchParams()
  if (params.page !== undefined) search.set('page', String(params.page))
  if (params.page_size !== undefined) search.set('page_size', String(params.page_size))
  const qs = search.toString()
  return qs ? `?${qs}` : ''
}

function buildLoginHistoryQuery(params: AdminLoginHistoryListParams): string {
  const search = new URLSearchParams()
  if (params.page !== undefined) search.set('page', String(params.page))
  if (params.page_size !== undefined) search.set('page_size', String(params.page_size))
  if (params.tenant_id) search.set('tenant_id', params.tenant_id)
  if (params.result) search.set('result', params.result)
  if (params.actor_id) search.set('actor_id', params.actor_id)
  const qs = search.toString()
  return qs ? `?${qs}` : ''
}

function buildAuditLogQuery(params: AdminAuditLogListParams): string {
  const search = new URLSearchParams()
  if (params.page !== undefined) search.set('page', String(params.page))
  if (params.page_size !== undefined) search.set('page_size', String(params.page_size))
  if (params.tenant_id) search.set('tenant_id', params.tenant_id)
  if (params.action) search.set('action', params.action)
  if (params.result) search.set('result', params.result)
  if (params.actor_id) search.set('actor_id', params.actor_id)
  const qs = search.toString()
  return qs ? `?${qs}` : ''
}

function buildQuery(params: AdminListParams): string {
  const search = new URLSearchParams()
  if (params.page !== undefined) search.set('page', String(params.page))
  if (params.page_size !== undefined) search.set('page_size', String(params.page_size))
  if (params.tenant_id) search.set('tenant_id', params.tenant_id)
  if (params.search) search.set('search', params.search)
  const qs = search.toString()
  return qs ? `?${qs}` : ''
}

export async function getAdminStats() {
  return apiFetch<AdminStatsResponse>('/api/v1/admin/stats', { method: 'GET' }, { auth: true })
}

export async function listAdminUsers(params: AdminListParams = {}) {
  return apiFetch<AdminUserResponse[]>(
    `/api/v1/admin/users${buildQuery(params)}`,
    { method: 'GET' },
    { auth: true },
  )
}

export async function getAdminUser(userId: string) {
  return apiFetch<AdminUserDetailResponse>(
    `/api/v1/admin/users/${encodeURIComponent(userId)}`,
    { method: 'GET' },
    { auth: true },
  )
}

export async function disableAdminUser(userId: string) {
  await apiFetchNoContent(
    `/api/v1/admin/users/${encodeURIComponent(userId)}/disable`,
    { method: 'POST' },
    { auth: true },
  )
}

export async function forceLogoutAdminUser(userId: string) {
  await apiFetchNoContent(
    `/api/v1/admin/users/${encodeURIComponent(userId)}/force-logout`,
    { method: 'POST' },
    { auth: true },
  )
}

export async function resetAdminUserPasskeys(userId: string) {
  await apiFetchNoContent(
    `/api/v1/admin/users/${encodeURIComponent(userId)}/reset-passkey`,
    { method: 'POST' },
    { auth: true },
  )
}

export async function resetAdminUserMFA(userId: string) {
  await apiFetchNoContent(
    `/api/v1/admin/users/${encodeURIComponent(userId)}/reset-mfa`,
    { method: 'POST' },
    { auth: true },
  )
}

export async function listAdminTenants(params: AdminListParams = {}) {
  return apiFetch<AdminTenantResponse[]>(
    `/api/v1/admin/tenants${buildQuery(params)}`,
    { method: 'GET' },
    { auth: true },
  )
}

export async function getAdminTenant(tenantId: string) {
  return apiFetch<AdminTenantResponse>(
    `/api/v1/admin/tenants/${encodeURIComponent(tenantId)}`,
    { method: 'GET' },
    { auth: true },
  )
}

export async function createAdminTenant(body: AdminCreateTenantRequest) {
  return apiFetch<AdminTenantResponse>(
    '/api/v1/admin/tenants',
    { method: 'POST', body: JSON.stringify(body) },
    { auth: true },
  )
}

export async function updateAdminTenant(tenantId: string, body: AdminUpdateTenantRequest) {
  return apiFetch<AdminTenantResponse>(
    `/api/v1/admin/tenants/${encodeURIComponent(tenantId)}`,
    { method: 'PATCH', body: JSON.stringify(body) },
    { auth: true },
  )
}

export async function deleteAdminTenant(tenantId: string) {
  await apiFetchNoContent(
    `/api/v1/admin/tenants/${encodeURIComponent(tenantId)}`,
    { method: 'DELETE' },
    { auth: true },
  )
}

export async function listTenantMembers(tenantId: string, params: PaginationParams = {}) {
  return apiFetch<AdminTenantMemberResponse[]>(
    `/api/v1/admin/tenants/${encodeURIComponent(tenantId)}/members${buildPaginationQuery(params)}`,
    { method: 'GET' },
    { auth: true },
  )
}

export async function listAdminClients(params: AdminListParams = {}) {
  return apiFetch<AdminClientResponse[]>(
    `/api/v1/admin/clients${buildQuery(params)}`,
    { method: 'GET' },
    { auth: true },
  )
}

export async function getAdminClient(clientId: string) {
  return apiFetch<AdminClientResponse>(
    `/api/v1/admin/clients/${encodeURIComponent(clientId)}`,
    { method: 'GET' },
    { auth: true },
  )
}

export async function createAdminClient(body: AdminCreateClientRequest) {
  return apiFetch<AdminCreateClientResponse>(
    '/api/v1/admin/clients',
    { method: 'POST', body: JSON.stringify(body) },
    { auth: true },
  )
}

export async function updateAdminClient(clientId: string, body: AdminUpdateClientRequest) {
  return apiFetch<AdminClientResponse>(
    `/api/v1/admin/clients/${encodeURIComponent(clientId)}`,
    { method: 'PATCH', body: JSON.stringify(body) },
    { auth: true },
  )
}

export async function deleteAdminClient(clientId: string) {
  await apiFetchNoContent(
    `/api/v1/admin/clients/${encodeURIComponent(clientId)}`,
    { method: 'DELETE' },
    { auth: true },
  )
}

export async function getAdminClientUsage(clientId: string) {
  return apiFetch<AdminClientUsageResponse>(
    `/api/v1/admin/clients/${encodeURIComponent(clientId)}/usage`,
    { method: 'GET' },
    { auth: true },
  )
}

export async function listAdminAuditLogs(params: AdminAuditLogListParams = {}) {
  return apiFetch<AdminAuditLogResponse[]>(
    `/api/v1/admin/audit-logs${buildAuditLogQuery(params)}`,
    { method: 'GET' },
    { auth: true },
  )
}

export async function listAdminLoginHistory(params: AdminLoginHistoryListParams = {}) {
  return apiFetch<AdminAuditLogResponse[]>(
    `/api/v1/admin/login-history${buildLoginHistoryQuery(params)}`,
    { method: 'GET' },
    { auth: true },
  )
}

export async function listTenantIdentityProviders(tenantId: string, params: PaginationParams = {}) {
  return apiFetch<AdminIdentityProviderResponse[]>(
    `/api/v1/admin/tenants/${encodeURIComponent(tenantId)}/identity-providers${buildPaginationQuery(params)}`,
    { method: 'GET' },
    { auth: true },
  )
}

export async function addTenantMember(tenantId: string, body: AdminAddMemberRequest) {
  await apiFetchNoContent(
    `/api/v1/admin/tenants/${encodeURIComponent(tenantId)}/members`,
    { method: 'POST', body: JSON.stringify(body) },
    { auth: true },
  )
}

export async function removeTenantMember(tenantId: string, userId: string) {
  await apiFetchNoContent(
    `/api/v1/admin/tenants/${encodeURIComponent(tenantId)}/members/${encodeURIComponent(userId)}`,
    { method: 'DELETE' },
    { auth: true },
  )
}

export async function patchIdentityProvider(tenantId: string, provider: string, body: PatchIdentityProviderRequest) {
  await apiFetchNoContent(
    `/api/v1/admin/tenants/${encodeURIComponent(tenantId)}/identity-providers/${encodeURIComponent(provider)}`,
    { method: 'PATCH', body: JSON.stringify(body) },
    { auth: true },
  )
}

export async function listFederationProviders(tenantId: string) {
  const params = new URLSearchParams({ tenant_id: tenantId })
  return apiFetch<PublicFederationProviderResponse[]>(
    `/api/v1/federation/providers?${params.toString()}`,
    { method: 'GET' },
  )
}

export type {
  AuthenticationResponseJSON,
  PublicKeyCredentialCreationOptionsJSON,
  PublicKeyCredentialRequestOptionsJSON,
  RegistrationResponseJSON,
}
