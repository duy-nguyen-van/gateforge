# GateForge SDKs

Official client SDKs for [GateForge IAM](../README.md), generated from the canonical OpenAPI contract at [`api/openapi.yaml`](../api/openapi.yaml).

## Layout

| Path | Purpose |
|------|---------|
| [`config/`](config/) | OpenAPI Generator configs (`go.yaml`, `typescript.yaml`) |
| [`go/`](go/) | Handwritten Go SDK (`package gateforge`) + generated `openapi/` client |
| [`typescript/`](typescript/) | Handwritten TypeScript SDK (`@gateforge/sdk`) + generated `src/generated/` |

## Generate clients

Requires Docker. From the monorepo root:

```bash
make sdk-generate
```

This runs [OpenAPI Generator](https://openapi-generator.tech/) `v7.19.0`:

1. **Go** → `sdk/go/openapi`
2. **TypeScript (fetch)** → `sdk/typescript/src/generated` (then NodeNext import patch)

Drift check:

```bash
make sdk-check
```

## Verify

```bash
make sdk-go-test    # unit tests (httptest)
make sdk-ts-build   # npm ci && build
cd sdk/typescript && npm test
```

Live integration (optional):

```bash
GATEFORGE_BASE_URL=http://localhost:3000 go test -tags=integration ./sdk/go
```

## Auth guidance

| Concern | Go SDK | TypeScript SDK |
|---------|--------|----------------|
| Bearer JWT | `WithTokenProvider` | `getAccessToken` / `refreshAccessToken` |
| OIDC code + PKCE | `PKCEGenerate`, `BuildAuthorizeURL`, `ExchangeAuthorizationCode` | `createPKCE`, `buildAuthorizeUrl`, `parseCallbackParams` |
| Token introspection (RFC 7662) | `IntrospectToken` (confidential client) | `client.oidc.introspectToken` (confidential client) |
| Browser cookies / CSRF | App-owned | `credentials: 'include'`, `prefetchCsrfToken` |
| WebAuthn ceremonies | Raw start/finish types only | `registerPasskey` / `loginWithPasskey` via `@simplewebauthn/browser` |
| Token storage | App-owned | Default memory store; persistence must be explicit |

Never put admin API keys or client secrets in browser bundles. Prefer in-memory access tokens; do not log refresh tokens.

## Releases (v0.x)

Publish only from explicit tags — never from every `main` build.

| SDK | Tag pattern | Consumer install |
|-----|-------------|------------------|
| Go | `sdk/go/v0.1.0` | `go get github.com/gateforge-iam/gateforge-iam/sdk/go@sdk/go/v0.1.0` |
| TypeScript | `sdk/typescript/v0.1.0` | `npm install @gateforge/sdk@0.1.0` (CI uses npm trusted publishing / OIDC) |

Workflow: [`.github/workflows/sdk-release.yml`](../.github/workflows/sdk-release.yml).

## Package docs

- [Go SDK](go/README.md)
- [TypeScript SDK](typescript/README.md)
- [API contract](../api/README.md)
