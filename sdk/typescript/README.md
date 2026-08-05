# @gateforge/sdk

Official TypeScript SDK for GateForge IAM. Wraps the generated OpenAPI fetch client with browser auth helpers (PKCE, CSRF, federation) and WebAuthn/passkey flows.

## Install

```bash
npm install @gateforge/sdk
```

## GateForgeClient

```ts
import { GateForgeClient, MemoryTokenStore } from '@gateforge/sdk'

const tokens = new MemoryTokenStore()

const client = new GateForgeClient({
  baseUrl: 'https://iam.example.com',
  credentials: 'include',
  getAccessToken: () => tokens.getAccessToken(),
  refreshAccessToken: async () => {
    const refresh = tokens.getRefreshToken()
    if (!refresh) return null
    const res = await client.auth.refreshToken({
      refreshTokenRequest: { refresh_token: refresh },
    })
    tokens.setTokens(res.data!)
    return res.data?.access_token
  },
  onUnauthorized: () => tokens.clear(),
})

const health = await client.health.getHealth()
const me = await client.users.getMe()
```

Generated API groups: `auth`, `users`, `admin`, `health`, `mfa`, `webauthn`, `oidc`, `federation`, `tenants`.

## OIDC PKCE (browser)

```ts
import {
  GateForgeClient,
  createPKCE,
  buildAuthorizeUrl,
  generateState,
  generateNonce,
  parseCallbackParams,
  prefetchCsrfToken,
} from '@gateforge/sdk'

const client = new GateForgeClient({ baseUrl: 'https://iam.example.com' })

const pkce = await createPKCE()
const state = generateState()
const nonce = generateNonce()
sessionStorage.setItem('pkce_verifier', pkce.verifier)
sessionStorage.setItem('oauth_state', state)

const authorizeUrl = buildAuthorizeUrl('https://iam.example.com', {
  clientId: 'my-app',
  redirectUri: 'https://app.example.com/callback',
  scope: 'openid profile',
  state,
  nonce,
  codeChallenge: pkce.challenge,
})
window.location.assign(authorizeUrl)

// On the callback page:
const { code, state: returnedState, error } = parseCallbackParams(window.location.href)
if (error) throw new Error(error)
if (returnedState !== sessionStorage.getItem('oauth_state')) throw new Error('state mismatch')

const token = await client.oidc.createToken({
  grantType: 'authorization_code',
  code: code!,
  redirectUri: 'https://app.example.com/callback',
  clientId: 'my-app',
  codeVerifier: sessionStorage.getItem('pkce_verifier')!,
})
```

## Token introspection (RFC 7662)

Server-side / confidential clients only — never ship `client_secret` in a browser bundle.

```ts
import { GateForgeClient } from '@gateforge/sdk'

const client = new GateForgeClient({ baseUrl: 'https://iam.example.com' })

const result = await client.oidc.introspectToken({
  token: accessToken,
  tokenTypeHint: 'access_token',
  clientId: 'rs-client',
  clientSecret: 'rs-secret',
})

if (result.active) {
  console.log(result.sub, result.scope, result.token_type)
}
```

For browser `POST /oidc/login`, prefetch CSRF and pass it via the generated client's `apiKey` (or a `headers` default):

```ts
const csrf = await prefetchCsrfToken(client)
// e.g. new Configuration({ apiKey: async () => csrf! }) or headers: { 'X-CSRF-Token': csrf! }
```

Federation start URL:

```ts
import { federationStartUrl } from '@gateforge/sdk'

window.location.assign(
  federationStartUrl('https://iam.example.com', 'google', 'https://app.example.com/login/federation/complete'),
)
```

## Passkeys (WebAuthn)

```ts
import { GateForgeClient, registerPasskey, loginWithPasskey } from '@gateforge/sdk'

const client = new GateForgeClient({
  baseUrl: 'https://iam.example.com',
  getAccessToken: () => sessionStorage.getItem('access_token'),
})

// Authenticated user registers a passkey
await registerPasskey(client, 'MacBook Touch ID')

// Passwordless sign-in
const result = await loginWithPasskey(client, 'user@example.com', {
  rememberMe: true,
})
```

Low-level browser helpers (`startRegistration` / `startAuthentication`) wrap `@simplewebauthn/browser` with `optionsJSON`.

## Scripts

```bash
npm run build
npm run typecheck
npm test
```

## Regenerate OpenAPI client

From the monorepo root (requires Docker + `api/openapi.yaml`):

```bash
make sdk-generate
```

Generation runs `sdk/typescript/scripts/patch-generated.mjs` so the fetch client typechecks under NodeNext.
