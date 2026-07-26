import type { GateForgeClient } from './client.js'

/** PKCE verifier/challenge pair (S256). */
export type PKCE = {
  verifier: string
  challenge: string
  method: 'S256'
}

/** Query parameters for an OIDC authorization request. */
export type AuthorizeParams = {
  clientId: string
  redirectUri: string
  scope?: string
  state?: string
  nonce?: string
  codeChallenge?: string
  codeChallengeMethod?: string
  responseType?: string
  /** Additional query parameters (e.g. login_hint, prompt). */
  extra?: Record<string, string>
}

/** Parsed OAuth/OIDC redirect callback query parameters. */
export type CallbackParams = {
  code?: string
  state?: string
  error?: string
  error_description?: string
}

/** Pluggable token storage (access is typically memory-only). */
export type TokenStore = {
  getAccessToken(): string | null | undefined
  getRefreshToken(): string | null | undefined
  setTokens(tokens: {
    access_token: string
    refresh_token?: string | null
  }): void
  clear(): void
}

function base64UrlEncode(bytes: Uint8Array): string {
  let binary = ''
  for (const b of bytes) {
    binary += String.fromCharCode(b)
  }
  return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '')
}

function randomUrlSafeString(byteLength = 32): string {
  const buf = new Uint8Array(byteLength)
  crypto.getRandomValues(buf)
  return base64UrlEncode(buf)
}

/**
 * Create a Proof Key for Code Exchange (PKCE) S256 pair.
 * Works in browsers and Node 20+ (Web Crypto).
 */
export async function createPKCE(): Promise<PKCE> {
  const verifier = randomUrlSafeString(32)
  const digest = await crypto.subtle.digest(
    'SHA-256',
    new TextEncoder().encode(verifier),
  )
  const challenge = base64UrlEncode(new Uint8Array(digest))
  return { verifier, challenge, method: 'S256' }
}

/** Generate a cryptographically random OAuth `state` value. */
export function generateState(): string {
  return randomUrlSafeString(16)
}

/** Generate a cryptographically random OIDC `nonce` value. */
export function generateNonce(): string {
  return randomUrlSafeString(16)
}

/**
 * Build an OIDC/OAuth2 authorize URL.
 * `issuerOrAuthorizeUrl` may be an issuer base (appends `/authorize`) or a
 * full authorize endpoint URL.
 */
export function buildAuthorizeUrl(
  issuerOrAuthorizeUrl: string,
  params: AuthorizeParams,
): string {
  const base = issuerOrAuthorizeUrl.replace(/\/+$/, '')
  if (!base) {
    throw new Error('auth: issuer or authorize URL is required')
  }
  if (!params.clientId) {
    throw new Error('auth: clientId is required')
  }
  if (!params.redirectUri) {
    throw new Error('auth: redirectUri is required')
  }

  const authorizeUrl = base.endsWith('/authorize') ? base : `${base}/authorize`
  const url = new URL(authorizeUrl)
  url.searchParams.set('client_id', params.clientId)
  url.searchParams.set('redirect_uri', params.redirectUri)
  url.searchParams.set('response_type', params.responseType ?? 'code')
  if (params.scope) {
    url.searchParams.set('scope', params.scope)
  }
  if (params.state) {
    url.searchParams.set('state', params.state)
  }
  if (params.nonce) {
    url.searchParams.set('nonce', params.nonce)
  }
  if (params.codeChallenge) {
    url.searchParams.set('code_challenge', params.codeChallenge)
    url.searchParams.set(
      'code_challenge_method',
      params.codeChallengeMethod ?? 'S256',
    )
  }
  if (params.extra) {
    for (const [key, value] of Object.entries(params.extra)) {
      if (key) {
        url.searchParams.set(key, value)
      }
    }
  }
  return url.toString()
}

/** Parse `code` / `state` / `error` from an OAuth redirect URL. */
export function parseCallbackParams(url: string | URL): CallbackParams {
  const parsed = typeof url === 'string' ? new URL(url) : url
  const get = (key: string): string | undefined => {
    const value = parsed.searchParams.get(key)
    return value == null || value === '' ? undefined : value
  }
  return {
    code: get('code'),
    state: get('state'),
    error: get('error'),
    error_description: get('error_description'),
  }
}

/**
 * Prefetch CSRF token via GET `/.well-known/openid-configuration`.
 * GateForge sets `X-CSRF-Token` on this response for browser OIDC POSTs.
 */
export async function prefetchCsrfToken(
  client: GateForgeClient,
): Promise<string | undefined> {
  const response = await client.fetch(
    client.url('/.well-known/openid-configuration'),
    { credentials: client.credentials },
  )
  return response.headers.get('X-CSRF-Token') ?? undefined
}

/**
 * Build the browser redirect URL to start federated sign-in.
 * Hits `GET /oidc/federation/{provider}/start?return_to=...`.
 */
export function federationStartUrl(
  baseUrl: string,
  provider: string,
  returnTo: string,
): string {
  const base = baseUrl.replace(/\/+$/, '')
  const url = new URL(
    `${base}/oidc/federation/${encodeURIComponent(provider)}/start`,
  )
  url.searchParams.set('return_to', returnTo)
  return url.toString()
}

/** In-memory token store (default — do not put access tokens in localStorage). */
export class MemoryTokenStore implements TokenStore {
  private accessToken: string | null = null
  private refreshToken: string | null = null

  getAccessToken(): string | null {
    return this.accessToken
  }

  getRefreshToken(): string | null {
    return this.refreshToken
  }

  setTokens(tokens: {
    access_token: string
    refresh_token?: string | null
  }): void {
    this.accessToken = tokens.access_token
    this.refreshToken =
      tokens.refresh_token === undefined
        ? this.refreshToken
        : tokens.refresh_token
  }

  clear(): void {
    this.accessToken = null
    this.refreshToken = null
  }
}
