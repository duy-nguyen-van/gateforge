import { ApiError } from '@/api/types'

/** Surface API, browser WebAuthn, and network errors for passkey UI. */
export function formatWebAuthnError(err: unknown, fallback: string): string {
  if (err instanceof ApiError) {
    return err.message
  }
  if (err instanceof Error && err.message) {
    return err.message
  }
  return fallback
}

/** Unwrap legacy go-webauthn `{ publicKey }` envelopes if present. */
export function unwrapWebAuthnOptions<T>(options: T): T {
  if (options && typeof options === 'object' && 'publicKey' in options) {
    const wrapped = options as T & { publicKey?: T }
    if (wrapped.publicKey) {
      return wrapped.publicKey
    }
  }
  return options
}
