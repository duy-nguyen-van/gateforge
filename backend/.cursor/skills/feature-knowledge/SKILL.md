---
name: feature-knowledge
description: >-
  Load IAM feature documentation before editing OIDC, SSO, federation, WebAuthn,
  MFA, sessions, tenants, or admin auth. Use when implementing or refactoring
  identity flows in iam-backend.
---

# IAM feature knowledge

## When to use

You are about to change handlers, services, auth, or migrations for identity features — read docs **before** coding.

## Read first

1. [docs/README.md](../../../docs/README.md) — hub, DB matrix, Redis keys, route surfaces
2. Feature guide for your area:

| Area | Feature doc | Testing doc |
|------|-------------|-------------|
| OIDC / PKCE | [docs/features/OIDC.md](../../../docs/features/OIDC.md) | [docs/testing/OIDC_CURL.md](../../../docs/testing/OIDC_CURL.md) |
| SSO cookie | [docs/features/SSO_SESSION.md](../../../docs/features/SSO_SESSION.md) | [docs/testing/SSO_CURL.md](../../../docs/testing/SSO_CURL.md) |
| Google federation | [docs/features/FEDERATION.md](../../../docs/features/FEDERATION.md) | [docs/testing/FEDERATION_CURL.md](../../../docs/testing/FEDERATION_CURL.md) |
| Passkeys + MFA | [docs/features/PASSKEY_MFA.md](../../../docs/features/PASSKEY_MFA.md) | [docs/testing/PASSKEY_MFA_CURL.md](../../../docs/testing/PASSKEY_MFA_CURL.md) |
| Multi-tenant | [docs/features/MULTI_TENANT.md](../../../docs/features/MULTI_TENANT.md) | login/OIDC testing docs |
| Authorization | [docs/features/AUTHORIZATION.md](../../../docs/features/AUTHORIZATION.md) | admin console |

## After changes

Update the feature doc + testing doc + swag (for `/api/v1`) in the same change.
