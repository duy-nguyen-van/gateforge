import {
  startAuthentication as sbStartAuth,
  startRegistration as sbStartReg,
  type AuthenticationResponseJSON,
  type PublicKeyCredentialCreationOptionsJSON,
  type PublicKeyCredentialRequestOptionsJSON,
  type RegistrationResponseJSON,
} from '@simplewebauthn/browser'

import type { GateForgeClient } from './client.js'
import type { EmptyDataEnvelope, LoginResultEnvelope } from './generated/src/index.js'

export type RegistrationOptions = PublicKeyCredentialCreationOptionsJSON
export type AuthenticationOptions = PublicKeyCredentialRequestOptionsJSON

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

/**
 * Start WebAuthn registration (passkey create) via `@simplewebauthn/browser`.
 */
export async function startRegistration(
  options: RegistrationOptions,
): Promise<RegistrationResponseJSON> {
  return sbStartReg({ optionsJSON: options })
}

/**
 * Start WebAuthn authentication (passkey get) via `@simplewebauthn/browser`.
 */
export async function startAuthentication(
  options: AuthenticationOptions,
): Promise<AuthenticationResponseJSON> {
  return sbStartAuth({ optionsJSON: options })
}

export type LoginWithPasskeyOptions = {
  tenantId?: string
  rememberMe?: boolean
  returnTo?: string
}

/**
 * Full passkey registration: start API → browser create → finish API.
 * Requires an authenticated GateForge session / access token.
 */
export async function registerPasskey(
  client: GateForgeClient,
  deviceName?: string,
): Promise<EmptyDataEnvelope> {
  const start = await client.webauthn.startWebauthnRegister({
    webauthnRegisterStartRequest: deviceName
      ? { device_name: deviceName }
      : undefined,
  })

  const options = unwrapWebAuthnOptions(
    start.data?.options,
  ) as RegistrationOptions | undefined
  if (!options) {
    throw new Error('webauthn: register/start missing options')
  }
  const sessionToken = start.data?.session_token
  if (!sessionToken) {
    throw new Error('webauthn: register/start missing session_token')
  }

  const credential = await startRegistration(options)
  return client.webauthn.finishWebauthnRegister({
    webauthnRegisterFinishRequest: {
      session_token: sessionToken,
      credential: credential as unknown as { [key: string]: any },
    },
  })
}

/**
 * Full passkey login: start API → browser get → finish API.
 */
export async function loginWithPasskey(
  client: GateForgeClient,
  email: string,
  opts?: LoginWithPasskeyOptions,
): Promise<LoginResultEnvelope> {
  const start = await client.webauthn.startWebauthnLogin({
    webauthnLoginStartRequest: {
      email,
      tenant_id: opts?.tenantId,
    },
  })

  const options = unwrapWebAuthnOptions(
    start.data?.options,
  ) as AuthenticationOptions | undefined
  if (!options) {
    throw new Error('webauthn: login/start missing options')
  }
  const sessionToken = start.data?.session_token
  if (!sessionToken) {
    throw new Error('webauthn: login/start missing session_token')
  }

  const credential = await startAuthentication(options)
  return client.webauthn.finishWebauthnLogin({
    webauthnLoginFinishRequest: {
      email,
      session_token: sessionToken,
      credential: credential as unknown as { [key: string]: any },
      tenant_id: opts?.tenantId,
      remember_me: opts?.rememberMe,
      return_to: opts?.returnTo,
    },
  })
}
