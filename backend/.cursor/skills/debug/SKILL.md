---
name: debug
description: "Debug systematically with root cause analysis before fixes. Use for bugs, test failures, unexpected behavior, performance issues, log analysis, CI/CD failures, database diagnostics."
---

# Debugging (backend)

## Core Principle
**NO FIXES WITHOUT ROOT CAUSE INVESTIGATION FIRST**

## When to Use
- Test failures, OIDC/federation/MFA flow bugs
- Session/cache inconsistencies
- Migration or DB errors
- CI failures (`make lint`, `make tests`)

## Framework

### Phase 1: Root Cause
1. Read error messages and stack traces
2. When did it start? (`git log`, recent deploy/config)
3. Scope: all users vs one tenant vs one endpoint
4. Form 2–3 hypotheses

### Phase 2: Evidence
- `go test -race ./internal/<package>`
- `make test-specific TEST=...`
- Manual repro: flow docs below — **use correct path prefix** (root OIDC vs `/api/v1`)
- Postgres: `make migrate-status`; Redis: session/MFA keys
- Logs: global Zap `logger.Sugar`; correlation ID via `internal/request/context.go`

### IAM flow docs
| Doc | Paths |
|-----|-------|
| [docs/README.md](../../../docs/README.md) | Hub, DB, Redis |
| [docs/features/OIDC.md](../../../docs/features/OIDC.md) | `/authorize`, `/token`, `/userinfo` |
| [docs/features/FEDERATION.md](../../../docs/features/FEDERATION.md) | `/oidc/federation/*` |
| [docs/features/SSO_SESSION.md](../../../docs/features/SSO_SESSION.md) | `iam_session`, login + authorize |
| [docs/features/PASSKEY_MFA.md](../../../docs/features/PASSKEY_MFA.md) | `/api/v1/webauthn/*`, `/api/v1/mfa/*` |

### Phase 3: Fix & Verify
- Minimal fix at root layer
- Regression test in colocated `*_test.go`
- `make tests`

## Rules
- Never guess — use evidence
- Fix root cause, not symptoms
- Always write a regression test when fixing bugs

See `.cursor/rules/debugging-guide.mdc`.
