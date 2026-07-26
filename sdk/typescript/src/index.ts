export {
  GateForgeClient,
  type GateForgeClientOptions,
  type TokenProvider,
  type RefreshAccessToken,
} from './client.js'
export { APIError, GateForgeError, type Meta, type Envelope } from './errors.js'
export {
  createPKCE,
  buildAuthorizeUrl,
  generateState,
  generateNonce,
  parseCallbackParams,
  prefetchCsrfToken,
  federationStartUrl,
  MemoryTokenStore,
  type PKCE,
  type AuthorizeParams,
  type CallbackParams,
  type TokenStore,
} from './auth.js'
export {
  startRegistration,
  startAuthentication,
  registerPasskey,
  loginWithPasskey,
  unwrapWebAuthnOptions,
  type RegistrationOptions,
  type AuthenticationOptions,
  type LoginWithPasskeyOptions,
} from './webauthn.js'

/** Generated OpenAPI models, runtime helpers, and `*Api` classes. */
export * from './generated/index.js'
