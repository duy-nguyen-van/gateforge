---
name: debug
description: "Debug systematically with root cause analysis before fixes. Use for bugs, build failures, unexpected behavior, auth/OIDC/MFA issues."
---

# Debugging (frontend)

## Core Principle
**NO FIXES WITHOUT ROOT CAUSE INVESTIGATION FIRST**

## When to Use
- Auth/login/MFA/WebAuthn flow bugs
- Token refresh or session lost on reload
- OIDC redirect loops or CSRF failures
- Build/lint failures (`make check`)
- Admin console data not loading

## Framework

### Phase 1: Root Cause
1. Read console errors and network failures
2. When did it start? (`git log`, recent env changes)
3. Scope: all users vs admin vs specific route
4. Form 2–3 hypotheses

### Phase 2: Evidence
- Browser Network tab: URL, status, cookies, CSRF header
- Application tab: sessionStorage/localStorage keys
- `make build` for TypeScript errors
- Backend running on `VITE_API_PROXY_TARGET` (default `:3000`)

### IAM references (backend)
| Doc | Use for |
|-----|---------|
| [backend/docs/README.md](../../../backend/docs/README.md) | Hub |
| [backend/docs/features/OIDC.md](../../../backend/docs/features/OIDC.md) | Authorize, token |
| [backend/docs/features/SSO_SESSION.md](../../../backend/docs/features/SSO_SESSION.md) | SSO cookie |
| [backend/docs/features/FEDERATION.md](../../../backend/docs/features/FEDERATION.md) | Google login |
| [backend/docs/features/PASSKEY_MFA.md](../../../backend/docs/features/PASSKEY_MFA.md) | Passkeys, MFA |

### Phase 3: Fix & Verify
- Minimal fix at root layer (auth/api/feature/route)
- `make check`
- Manual browser verification

## Rules
- Never guess — use network/console evidence
- Fix root cause, not symptoms
- Check Vite proxy before blaming backend

See `.cursor/rules/debugging-guide.mdc`.
