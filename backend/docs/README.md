# GateForge IAM — documentation hub

Feature guides, database reference, and manual testing for the iam-backend identity service.

**Before changing OIDC, SSO, federation, WebAuthn, MFA, sessions, tenants, or admin auth:** read this hub and the relevant [features/](features/) guide. After code changes, update the feature doc, matching [testing/](testing/) guide, and `docs/swagger.yaml` (for `/api/v1` routes).

## Feature index

| Feature | Guide | Manual testing |
|---------|-------|----------------|
| OIDC (authorization code + PKCE) | [features/OIDC.md](features/OIDC.md) | [testing/OIDC_CURL.md](testing/OIDC_CURL.md) |
| SSO browser session (`iam_session`) | [features/SSO_SESSION.md](features/SSO_SESSION.md) | [testing/SSO_CURL.md](testing/SSO_CURL.md) |
| Google federation (upstream OAuth) | [features/FEDERATION.md](features/FEDERATION.md) | [testing/FEDERATION_CURL.md](testing/FEDERATION_CURL.md) |
| Passkeys + TOTP MFA | [features/PASSKEY_MFA.md](features/PASSKEY_MFA.md) | [testing/PASSKEY_MFA_CURL.md](testing/PASSKEY_MFA_CURL.md) |
| Multi-tenant memberships | [features/MULTI_TENANT.md](features/MULTI_TENANT.md) | (see login/OIDC testing docs) |
| Authorization (platform admin, tenant roles) | [features/AUTHORIZATION.md](features/AUTHORIZATION.md) | (admin console + JWT) |

Postman: [postman/IAM_OIDC.postman_collection.json](postman/IAM_OIDC.postman_collection.json)

OpenAPI (`/api/v1` only): [swagger.yaml](swagger.yaml) or `/swagger` (non-prod, basic auth).

## Route surfaces

Registered in `cmd/server/routes/router.go`. Two URL namespaces:

| Surface | Base | Examples | Auth |
|---------|------|----------|------|
| **OIDC / OAuth2** | `/` (root) | `/.well-known/openid-configuration`, `/authorize`, `/token`, `/userinfo`, `/oidc/login`, `/oidc/federation/:provider/*` | Browser `iam_session` for authorize; Bearer for userinfo |
| **App API** | `/api/v1` | `/register`, `/login`, `/refresh`, `/logout`, `/me`, `/webauthn/*`, `/mfa/*`, `/tenants/*`, `/admin/*` | Public, JWT Bearer, or `X-Admin-API-Key` (internal IdP toggle) |

Root OIDC routes are **outside** swag `@BasePath` — document them in [features/](features/), not only in OpenAPI.

## Feature → persistence matrix

| Feature | PostgreSQL | Redis |
|---------|------------|-------|
| OIDC | `clients`, `authorization_codes`, `access_tokens`, `refresh_tokens`, `consents` | — |
| SSO session | `sessions`, `refresh_tokens` (revoke on logout) | — |
| Federation | `tenant_identity_providers`, `federated_identities`, `users`, `tenant_memberships` | `oidc_federation_state:{state}` (10m) |
| Passkey | `webauthn_credentials` | `iam:webauthn:reg:{token}`, `iam:webauthn:login:{token}` |
| MFA | `user_mfa_totps`, `user_mfa_recovery_codes` | `iam:mfa:pending:{ticket}` |
| Multi-tenant | `tenant_memberships`, `tenants`, `users` | — |
| Authorization | `users.is_platform_admin`, `tenant_memberships.role` | — |
| Audit logs | `audit_logs` | — |

## PostgreSQL tables

| Table | Purpose | Key columns |
|-------|---------|-------------|
| `tenants` | Organizations | `id`, `name`, `domain` |
| `tenant_memberships` | User ↔ tenant access | `user_id`, `tenant_id`, `role` (`member`, `admin`), `status` |
| `users` | Global identities | `email_lower` (unique), `is_platform_admin`, `status` |
| `password_credentials` | Password hashes | `user_id`, `password_hash` |
| `clients` | OAuth/OIDC clients per tenant | `tenant_id`, `client_id`, `redirect_uris`, `is_public` |
| `authorization_codes` | Auth code + PKCE metadata | `code`, `user_id`, `tenant_id`, `code_challenge`, `nonce`, `scope`, `expires_at` |
| `access_tokens` | Opaque OIDC access token hashes | `token_hash`, `user_id`, `oauth_client_id`, `expires_at` |
| `refresh_tokens` | Refresh token hashes | `token_hash`, `user_id`, `oauth_client_id`, `revoked`, `expires_at` |
| `consents` | User consent per OAuth client | `user_id`, `oauth_client_id`, `scopes`, `granted` |
| `sessions` | Browser SSO sessions | `id` (= cookie value), `user_id`, `tenant_id`, `expires_at` |
| `webauthn_credentials` | Passkeys (user-scoped) | `user_id`, `credential_id`, `public_key`, `sign_count` |
| `user_mfa_totps` | TOTP secrets (encrypted) | `user_id`, `secret_encrypted`, `enabled`, `verified_at` |
| `user_mfa_recovery_codes` | Recovery code hashes | `user_id`, `code_hash`, `used_at` |
| `tenant_identity_providers` | Per-tenant IdP toggles | `tenant_id`, `provider`, `enabled` |
| `federated_identities` | Linked upstream IdP subjects | `user_id`, `provider`, `subject` (unique per provider+subject) |
| `audit_logs` | Immutable security/admin event trail | `action`, `result`, `actor_type`, `actor_id`, `resource_type`, `resource_id`, `old_value`, `new_value` |

Migrations: `cmd/migrations/sql/`. Multi-tenant schema: `20260617120000_multi_tenant_memberships.sql` (global users, no `tenant_id` on `webauthn_credentials` / `user_mfa_totps` / `federated_identities`).

## Redis ephemeral keys

Defined in `internal/auth/ephemeral_redis.go` and `internal/services/federation.go`:

| Key pattern | TTL | Payload | Behavior |
|-------------|-----|---------|----------|
| `iam:webauthn:reg:{token}` | `WEBAUTHN_SESSION_TTL` (default 5m) | WebAuthn registration session JSON | Single-use (load + delete) |
| `iam:webauthn:login:{token}` | `WEBAUTHN_SESSION_TTL` | WebAuthn login session JSON | Single-use |
| `iam:mfa:pending:{ticket}` | `MFA_PENDING_TICKET_TTL` (default 10m) | `MFAPendingPayload` (`user_id`, `tenant_id`, `remember_me`, `return_to`) | Single-use |
| `oidc_federation_state:{state}` | 10m | `return_to`, `tenant_id`, `nonce`, `provider` | Single-use |

Requires `CACHE_PROVIDER=redis` and reachable Redis.

## Code map (by layer)

| Layer | Path |
|-------|------|
| Routes | `cmd/server/routes/router.go` |
| Handlers | `internal/handlers/` (`oidc.go`, `auth.go`, `webauthn.go`, `mfa.go`, `admin.go`) |
| Services | `internal/services/` (`oidc.go`, `session.go`, `federation.go`, `webauthn.go`, `mfa.go`, `user.go`, `admin.go`) |
| Auth primitives | `internal/auth/` (`jwt.go`, `oidc_signer.go`, `ephemeral_redis.go`) |
| Repositories | `internal/repositories/` |
| Models | `internal/models/` |

Agent-oriented layer details: [AGENTS.md](../AGENTS.md).

## Local dev quick start

```bash
make container-up   # Postgres + Redis
make migrate-up
make up             # server :3000
```

Env reference: `cmd/server/.env.example`.
