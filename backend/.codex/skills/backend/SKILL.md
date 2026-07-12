---
name: backend
description: Use for backend work in this repository: Go, Echo, Uber FX, PostgreSQL, Redis, Atlas migrations, OIDC/OAuth2, sessions, federation, MFA/WebAuthn, OpenAPI, testing, debugging, or code review.
---

# GateForge IAM Backend

Use this skill when working under `backend/` or when the user asks about the IAM backend.

## First Reads

- For code discovery, prefer `codebase-memory-mcp` graph tools before grep.
- Read [backend/AGENTS.md](../../../AGENTS.md) for the current repo map, layers, route surfaces, and commands.
- For OIDC, SSO, federation, WebAuthn, MFA, tenant, or admin-auth changes, read `backend/docs/README.md` and the relevant `backend/docs/features/` guide before editing.

## Feature Docs Map

- OIDC: `docs/features/OIDC.md`; manual flow: `docs/testing/OIDC_CURL.md`
- SSO session: `docs/features/SSO_SESSION.md`; manual flow: `docs/testing/SSO_CURL.md`
- Federation: `docs/features/FEDERATION.md`; manual flow: `docs/testing/FEDERATION_CURL.md`
- Passkey/MFA: `docs/features/PASSKEY_MFA.md`; manual flow: `docs/testing/PASSKEY_MFA_CURL.md`
- Multi-tenant: `docs/features/MULTI_TENANT.md`
- Authorization/admin roles: `docs/features/AUTHORIZATION.md`

## Working Rules

- Keep handlers thin: bind/validate DTOs, call services, return JSON via `BaseHandler` / `ErrorHandler`.
- Put business rules in `internal/services/`; services accept `context.Context` and typed params, not `echo.Context`.
- Repositories use GORM and `internal/models`; they must not import handlers or services.
- New routes go in `cmd/server/routes/router.go`; choose root OIDC paths vs `/api/v1` deliberately.
- Update `internal/middlewares/csrf.go` only when a route truly belongs in the CSRF skip list.
- Use `internal/errors.AppError` for safe client errors unless extending a file with an explicit local exception.
- Never hardcode secrets or log tokens, MFA seeds, session IDs, credentials, or private keys.

## HTTP And Docs

- `/api/v1` routes are covered by swag OpenAPI. When routes, handlers, DTOs, status codes, params, tags, security, or annotations change, run `make swagger-load` and keep `docs/swagger.yaml`, `docs/swagger.json`, and `docs/docs.go` together.
- Root OIDC routes (`/authorize`, `/token`, `/userinfo`, `/.well-known/*`, `/oidc/*`) are outside `@BasePath`; update manual docs in `docs/features/` or `docs/testing/` when behavior changes.

## Tests And Checks

- For changed Go behavior, add/update the smallest meaningful test: success plus representative error path.
- Unit tests stay fast; Docker/Postgres tests require `//go:build integration`.
- Use table-driven tests and existing `testify/require` patterns.
- Run the narrowest useful check first, then broader checks when relevant:
  - `go test -race ./internal/<package>`
  - `make tests`
  - `make lint`
- After migration SQL changes in `cmd/migrations/sql/`, run `make migrate-hash`; do not hand-edit `atlas.sum`.

## Review Checklist

Check in this order: correctness, auth/security, error leaks, layer boundaries, concurrency/cache races, N+1 or extra Redis calls, tests, OpenAPI/manual docs drift.

## Git

Use conventional commits when asked to commit: `<type>(<scope>): <imperative summary>` under 72 chars, lowercase summary, no trailing period.
