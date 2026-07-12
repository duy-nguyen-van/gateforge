import type { LoginResponse } from '@/api/types'

const REFRESH_SESSION_KEY = 'iam_refresh_token'
const REFRESH_LOCAL_KEY = 'iam_refresh_token'

let accessToken: string | null = null
let refreshToken: string | null = null

type TokenListener = (tokens: LoginResponse | null) => void
const listeners = new Set<TokenListener>()

function readStoredRefreshToken(): string | null {
  return sessionStorage.getItem(REFRESH_SESSION_KEY) ?? localStorage.getItem(REFRESH_LOCAL_KEY)
}

function persistRefreshToken(token: string | null, rememberMe = false) {
  sessionStorage.removeItem(REFRESH_SESSION_KEY)
  localStorage.removeItem(REFRESH_LOCAL_KEY)
  if (!token) {
    return
  }
  if (rememberMe) {
    localStorage.setItem(REFRESH_LOCAL_KEY, token)
  } else {
    sessionStorage.setItem(REFRESH_SESSION_KEY, token)
  }
}

function notify() {
  const snapshot =
    accessToken && refreshToken
      ? {
          access_token: accessToken,
          refresh_token: refreshToken,
          token_type: 'Bearer',
          expires_in: 0,
          refresh_expires_in: 0,
        }
      : null
  listeners.forEach((listener) => listener(snapshot))
}

export function getAccessToken(): string | null {
  return accessToken
}

export function getRefreshToken(): string | null {
  return refreshToken
}

export function hydrateRefreshTokenFromStorage(): string | null {
  const stored = readStoredRefreshToken()
  if (stored) {
    refreshToken = stored
  }
  return stored
}

export function setTokens(tokens: LoginResponse | null, rememberMe = false) {
  if (!tokens) {
    accessToken = null
    refreshToken = null
    persistRefreshToken(null)
  } else {
    accessToken = tokens.access_token
    refreshToken = tokens.refresh_token
    persistRefreshToken(tokens.refresh_token, rememberMe)
  }
  notify()
}

export function clearTokens() {
  setTokens(null)
}

export function subscribeTokens(listener: TokenListener) {
  listeners.add(listener)
  return () => listeners.delete(listener)
}

let refreshPromise: Promise<boolean> | null = null

export async function refreshTokens(): Promise<boolean> {
  if (!refreshToken) {
    refreshToken = readStoredRefreshToken()
  }
  if (!refreshToken) {
    return false
  }

  if (refreshPromise) {
    return refreshPromise
  }

  refreshPromise = (async () => {
    try {
      const { refreshUser } = await import('@/api/client')
      const envelope = await refreshUser({ refresh_token: refreshToken! })
      const rememberMe = localStorage.getItem(REFRESH_LOCAL_KEY) !== null
      setTokens(envelope.data, rememberMe)
      return true
    } catch {
      clearTokens()
      return false
    } finally {
      refreshPromise = null
    }
  })()

  return refreshPromise
}
