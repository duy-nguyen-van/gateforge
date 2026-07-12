---
name: iam-api-manual-testing
description: >-
  Manual OIDC, auth, federation, SSO session, MFA/WebAuthn, and authorization testing for
  this repo using documented curl flows and Swagger. Use when the user asks to test
  OIDC, OAuth, login, token, federation, Google IdP, session/SSO, passkey, MFA,
  or to verify behavior with curl/Postman.
disable-model-invocation: true
---

# IAM API manual testing

## When to use

User wants step-by-step manual verification, curl examples, or to align implementation with documented IAM behavior.

## Route surfaces (know which base path)

| Surface | Base | Key paths |
|---------|------|-----------|
| **OIDC / OAuth2** | `/` (root) | `/.well-known/jwks.json`, `/.well-known/openid-configuration`, `/authorize`, `/token`, `/userinfo`, `/oidc/login`, `/oidc/federation/:provider/start`, `/oidc/federation/:provider/callback` |
| **App API** | `/api/v1` | `/register`, `/login`, `/refresh`, `/logout`, `/me`, `/webauthn/*`, `/mfa/*`, `/admin/*`, `/internal/tenants/:tenantId/identity-providers/google` |
| **Swagger** | `/swagger/*` | Non-prod only, basic auth (`router.go`) |

swag `@BasePath` is `/api/v1` — **OIDC root routes are not in `docs/swagger.yaml`**; use feature + testing docs below.

## Primary references

| Concept | Feature doc | curl / Postman |
|---------|-------------|----------------|
| Hub (DB, Redis) | [docs/README.md](../../../docs/README.md) | — |
| OIDC + PKCE | [docs/features/OIDC.md](../../../docs/features/OIDC.md) | [docs/testing/OIDC_CURL.md](../../../docs/testing/OIDC_CURL.md) |
| SSO session | [docs/features/SSO_SESSION.md](../../../docs/features/SSO_SESSION.md) | [docs/testing/SSO_CURL.md](../../../docs/testing/SSO_CURL.md) |
| Google federation | [docs/features/FEDERATION.md](../../../docs/features/FEDERATION.md) | [docs/testing/FEDERATION_CURL.md](../../../docs/testing/FEDERATION_CURL.md) |
| Passkeys + MFA | [docs/features/PASSKEY_MFA.md](../../../docs/features/PASSKEY_MFA.md) | [docs/testing/PASSKEY_MFA_CURL.md](../../../docs/testing/PASSKEY_MFA_CURL.md) |
| Multi-tenant | [docs/features/MULTI_TENANT.md](../../../docs/features/MULTI_TENANT.md) | login/OIDC testing docs |
| Authorization | [docs/features/AUTHORIZATION.md](../../../docs/features/AUTHORIZATION.md) | admin API |
| `/api/v1` shapes | `docs/swagger.yaml` or `/swagger` | — |

## Redis quick ref

| Key | Purpose |
|-----|---------|
| `iam:webauthn:reg:{token}` | Passkey registration ceremony |
| `iam:webauthn:login:{token}` | Passkey login ceremony |
| `iam:mfa:pending:{ticket}` | Post-login MFA step-up |
| `oidc_federation_state:{state}` | Google OAuth state |

## CSRF notes (`internal/middlewares/csrf.go`)

**Skipped** (no `X-CSRF-Token` required):
- `/api/v1/register`, `/login`, `/refresh`, `/logout`
- `/api/v1/webauthn/*`, `/api/v1/mfa/*`
- Root OIDC: `/token`, `/authorize`, `/userinfo`, `/.well-known/*`, `/oidc/federation/*`

**Not skipped** (CSRF applies):
- `/oidc/login` — browser OIDC password form; needs CSRF token from `X-CSRF-Token` header / cookie

## Environment

- Local server: `make up` (from repo root; uses `cmd/server`).
- Dependencies: `make container-up` then `make migrate-up` as needed (`Makefile`, `cmd/server/.env`).
- Env reference: `cmd/server/.env.example`, `examples/env/`.

## Workflow

1. Confirm scenario (authorization code, SSO session, federation callback, WebAuthn, MFA, admin IdP toggle, etc.).
2. Read the matching **feature doc**, then follow the **testing doc** with the correct path prefix (root vs `/api/v1`).
3. Follow documented parameter names (`client_id`, `redirect_uri`, `scopes`, `state`, `nonce`).
4. If behavior diverges from the doc, treat the doc as the product contract unless the user asks to change it — update **code + feature doc + testing doc + swag** (for `/api/v1`) in the same change.

## Optional

- `docs/postman/IAM_OIDC.postman_collection.json` for non-curl workflows.
